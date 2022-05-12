package websocket

import (
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/degatedev/degatesdk/conf"
	"github.com/degatedev/degatesdk/degate/binance"
	"github.com/degatedev/degatesdk/log"
	"github.com/gorilla/websocket"
)

type Protocol interface {
	HandlerMessage([]byte)
	Stop()
}

type WebSocketProtocol struct {
	Protocol
	config *conf.AppConfig
	conn   *websocket.Conn
	// webSocketProtocol Protocol
	subscribeMessages []*Msg
	closeCh           *chan bool
	isClose           bool
	isDone            bool
	url               string
	retryNum          int
}

type Msg struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

func (c *WebSocketProtocol) Geturl() string {
	return c.config.WebsocketBaseUrl + "ws"
}

func (c *WebSocketProtocol) Connect() {
	var err error
	if c.retryNum == 0 {
		log.Info("Connection with URL:%v", c.url)
		log.Info("Start to connect....")
		c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
		if err != nil {
			c.retryNum = 1
			log.Error("Can't connect to server. Reason: [%v]. Retrying: %v", err, 1)
			c.reconnect()
			return
		}
		c.connectSuccess()
	}
}

func (c *WebSocketProtocol) connectSuccess() {
	log.Info("Connection success with URL:%v", c.url)
	c.isClose = false
	c.isDone = false
	c.retryNum = 0
	c.conn.SetCloseHandler(nil)
	c.conn.SetPingHandler(nil)
	c.conn.SetPongHandler(nil)
	c.readMessage()
	c.subscribeMessage()
	c.pong()
}

func (c *WebSocketProtocol) pong() {
	go func() {
		var failPong = 0
		for {
			time.Sleep(time.Second * time.Duration(conf.PongInterval))
			if c.isClose || c.isDone {
				return
			}
			if c.conn == nil {
				continue
			}
			err := c.conn.WriteControl(websocket.PongMessage, []byte(""), time.Now().Add(time.Second))
			if err == websocket.ErrCloseSent {
				return
			} else if e, ok := err.(net.Error); ok && e.Temporary() {
				continue
			}
			if err != nil {
				failPong++
				log.Error("pong fail num:%v  error:%v", failPong, err)
			} else {
				failPong = 0
				log.Info("pong success")
			}
			if failPong >= 2 {
				c.Close()
				c.reconnect()
				return
			}
		}

	}()
}

func (c *WebSocketProtocol) reconnect() {
	go func() {
		for {
			if c.isDone {
				log.Info("websocket Reconnect stop: due websocket done")
				return
			}

			time.Sleep(time.Second * time.Duration(c.retryNum*conf.ConnectInterval))
			log.Info("Start to connect....")
			conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
			if c.isDone {
				return
			}
			if err == nil {
				c.conn = conn
				c.connectSuccess()
				return
			}
			c.retryNum++
			log.Error("Can't connect to server. Reason: [%v]. Retrying: %v", err, c.retryNum)
			if c.retryNum >= conf.MaxRetryNum {
				c.retryNum = 0
				msg := &binance.ErrorPayload{
					M: "Max reconnect retries reached",
				}
				msg.E = "error"
				cb, _ := json.Marshal(msg)
				c.HandlerMessage(cb)
				// c.webSocketProtocol.HandlerMessage(cb)
				return
			}
		}
	}()
}

func (c *WebSocketProtocol) subscribeMessage() {
	go func() {
		var messages = c.subscribeMessages
		for {
			var failSubscribe []*Msg
			if c.isDone {
				return
			}
			if c.conn != nil {
				for _, m := range messages {
					err := c.conn.WriteJSON(m)
					if err != nil {
						log.Error("error Subscribe %v: %v", strings.Join(m.Params, ","), err)
						failSubscribe = append(failSubscribe, m)
					}
				}
				if len(failSubscribe) == 0 {
					log.Info("All Subscribe Message Finish")
					return
				}
				messages = failSubscribe
			}
			time.Sleep(time.Second * 1)
		}
	}()

}

func (c *WebSocketProtocol) readMessage() {
	ch := make(chan bool, 1)
	c.closeCh = &ch
	go func() {
		for {
			select {
			case <-*c.closeCh:
				log.Info("websocket readMessage finish: due websocket done")
				return
			default:
				if c.isClose || c.conn == nil {
					return
				}
				_, message, err := c.conn.ReadMessage()
				if err != nil {
					log.Error("websocket ReadMessage error: %v", err)
					if e, ok := err.(*websocket.CloseError); ok {
						if e.Code == websocket.CloseAbnormalClosure {
							log.Error("websocket close err:[%v]. Retrying: %v", err, 1)
							c.reconnect()
							return
						}
						msg := &binance.ErrorPayload{
							M: err.Error(),
						}
						msg.E = "error"
						cb, _ := json.Marshal(msg)
						c.HandlerMessage(cb)
						// c.webSocketProtocol.HandlerMessage(cb)
						return
					}
					continue
				}
				if len(message) > 0 {
					log.Info("websocket ReadMessage message:%v", string(message))
					// c.webSocketProtocol.HandlerMessage(message)
					c.HandlerMessage(message)
				}
			}
		}
	}()
	return
}

func (c *WebSocketProtocol) UnSubscribe(params []string, id int) (err error) {
	if c.conn == nil {
		return
	}
	msg := &Msg{
		Method: "UNSUBSCRIBE",
		Params: params,
		ID:     id,
	}
	err = c.conn.WriteJSON(msg)
	if err != nil {
		log.Error("error UnSubscribe %v: %v", strings.Join(params, ","), err)
		return
	}
	return
}

func (c *WebSocketProtocol) Subscribe(params []string, id int) {
	c.url = c.Geturl()
	msg := &Msg{
		Method: "SUBSCRIBE",
		Params: params,
		ID:     id,
	}
	c.subscribeMessages = append(c.subscribeMessages, msg)
	c.Connect()
	return
}

func (c *WebSocketProtocol) SubscribeUser(listenKey string) {
	c.url = c.Geturl() + "/" + listenKey
	c.Connect()
}

func (c *WebSocketProtocol) Close() {
	if !c.isClose {
		c.isClose = true
		*c.closeCh <- true
	}
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Error("error websocket close: %v", err)
		}
		c.conn = nil
	}
}

func (c *WebSocketProtocol) Stop() {
	c.isDone = true
	c.Close()
}
