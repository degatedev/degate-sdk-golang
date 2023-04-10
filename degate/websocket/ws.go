package websocket

// copy from https://github.com/adshao/go-binance/

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint  string
	Proxy     string
	Keepalive bool
	Timeout   time.Duration
}

type SubscribeParam struct {
	Method string   `json:"method" form:"method"`
	Params []string `json:"params" form:"params"`
	ID     uint64   `json:"id" form:"id"`
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

func subscirbe(ctx context.Context, c *websocket.Conn, param *SubscribeParam) error {
	buf, err := json.Marshal(param)
	if err != nil {
		return err
	}
	fmt.Println("param:", string(buf))

	return c.Write(ctx, websocket.MessageText, buf)
}

var wsServe = func(cfg *WsConfig, param *SubscribeParam, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	c, _, err := websocket.Dial(ctx, cfg.Endpoint, nil)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	if param != nil {
		if err = subscirbe(ctx, c, param); err != nil {
			cancel()
			return nil, nil, err
		}
	}

	c.SetReadLimit(655350)
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer close(doneC)
		defer cancel()
		if cfg.Keepalive {
			go keepAlive(ctx, c, cfg.Timeout)
		}
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			_ = c.Close(websocket.StatusNormalClosure, "normal closure")
		}()
		for {
			_, message, readErr := c.Read(ctx)
			if readErr != nil {
				if !silent {
					errHandler(readErr)
				}
				return
			}
			handler(message)
		}
	}()
	return
}

func keepAlive(ctx context.Context, c *websocket.Conn, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}

		err := c.Ping(ctx)
		if err != nil {
			return
		}

		t.Reset(d)
	}
}

func wsUserDataServe() {

}
