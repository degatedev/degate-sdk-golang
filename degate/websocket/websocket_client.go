package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
	"github.com/degatedev/degate-sdk-golang/log"
)

const (
	KlineStreamName       = "%d.%d@kline_%d"
	TradeStreamName       = "%d.%d@trade"
	TickerStreamName      = "%d.%d@ticker"
	BookTickerStreamName  = "%d.%d@bookTicker"
	DepthStreamName       = "%d.%d@depth%d%s"
	DepthUpdateStreamName = "%d.%d@depth%s"
)

type WebSocketClient struct {
	config  *conf.AppConfig
	handler func(string)
	WebSocketProtocol
}

func (c *WebSocketClient) Init(config *conf.AppConfig) {
	c.config = config
	if len(c.config.WebsocketBaseUrl) == 0 {
		c.config.WebsocketBaseUrl = conf.WebsocketBaseUrl
	}
	if conf.Conf == nil {
		conf.Conf = new(conf.Config).Init()
	}
	conf.Conf.AddTokens(config.Tokens)
	c.BaseUrl = c.config.WebsocketBaseUrl
	c.MessageProcess = c.HandlerMessage
	SaveClient(c)
}

func (c *WebSocketClient) HandlerMessage(message []byte) {
	var (
		err      error
		payload  = map[string]interface{}{}
		response []byte
	)
	err = json.Unmarshal(message, &payload)
	if err != nil {
		log.Error("error websocket HandlerMessage: %v", err)
		return
	}
	if payload["e"] == "kline" {
		var (
			klinePayload  = &model.KlinePayload{}
			bklinePayload *binance.KlinePayload
		)
		err = json.Unmarshal(message, klinePayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bklinePayload, err = lib.ConvertKlinePayload(klinePayload); err != nil {
			log.Error("error websocket kline convert: %v", err)
			return
		}
		if response, err = json.Marshal(bklinePayload); err != nil {
			log.Error("error websocket kline json.Marshal: %v", err)
			return
		}
	} else if payload["e"] == "trade" {
		var (
			tradePayload  = &model.TradePayload{}
			bTradePayload *binance.TradePayload
		)
		err = json.Unmarshal(message, tradePayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bTradePayload, err = lib.ConvertTradePayload(tradePayload); err != nil {
			log.Error("error websocket trade convert: %v", err)
			return
		}
		if response, err = json.Marshal(bTradePayload); err != nil {
			log.Error("error websocket trade json.Marshal: %v", err)
			return
		}
	} else if payload["e"] == "24hrTicker" {
		var (
			tickerPayload  = &model.TickerPayload{}
			bTickerPayload *binance.TickerPayload
		)
		err = json.Unmarshal(message, tickerPayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bTickerPayload, err = lib.ConvertTickerPayload(tickerPayload); err != nil {
			log.Error("error websocket ticker convert: %v", err)
			return
		}
		if response, err = json.Marshal(bTickerPayload); err != nil {
			log.Error("error websocket ticker json.Marshal: %v", err)
			return
		}
	} else if payload["e"] == "bookTicker" {
		var (
			bookTickerPayload  = &model.BookTickerPayload{}
			bBookTickerPayload *binance.BookTickerPayload
		)
		err = json.Unmarshal(message, bookTickerPayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bBookTickerPayload, err = lib.ConvertBookTickerPayload(bookTickerPayload); err != nil {
			log.Error("error websocket bookTicker convert: %v", err)
			return
		}
		if response, err = json.Marshal(bBookTickerPayload); err != nil {
			log.Error("error websocket bookTicker json.Marshal: %v", err)
			return
		}
	} else if payload["e"] == "depth" {
		var (
			depthPayload  = &model.DepthPayload{}
			bDepthPayload *binance.DepthPayload
		)
		err = json.Unmarshal(message, depthPayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bDepthPayload, err = lib.ConvertDepthPayload(depthPayload); err != nil {
			log.Error("error websocket depth convert: %v", err)
			return
		}
		if response, err = json.Marshal(bDepthPayload); err != nil {
			log.Error("error websocket depth json.Marshal: %v", err)
			return
		}
	} else if payload["e"] == "depthUpdate" {
		var (
			depthPayload  = &model.DepthUpdatePayload{}
			bDepthPayload *binance.DepthUploadPayload
		)
		err = json.Unmarshal(message, depthPayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bDepthPayload, err = lib.ConvertDepthUploadPayload(depthPayload); err != nil {
			log.Error("error websocket depthUpload convert: %v", err)
			return
		}
		if response, err = json.Marshal(bDepthPayload); err != nil {
			log.Error("error websocket depthUpload json.Marshal: %v", err)
			return
		}
	} else if payload["e"] == "outboundAccountPosition" {
		var (
			outPayload  = &model.OutboundAccountPositionPayload{}
			bOutPayload *binance.OutboundAccountPositionPayload
		)
		err = json.Unmarshal(message, outPayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bOutPayload, err = lib.ConvertOutboundAccountPositionPayload(outPayload); err != nil {
			log.Error("error websocket outboundAccountPosition convert: %v", err)
			return
		}
		if response, err = json.Marshal(bOutPayload); err != nil {
			log.Error("error websocket outboundAccountPosition json.Marshal: %v", err)
			return
		}
	} else if payload["e"] == "balanceUpdate" {
		var (
			balancePayload  = &model.BalanceUpdatePayload{}
			bBalancePayload *binance.BalanceUpdatePayload
		)
		err = json.Unmarshal(message, balancePayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bBalancePayload, err = lib.ConvertBalanceUpdatePayload(balancePayload); err != nil {
			log.Error("error websocket balanceUpdate convert: %v", err)
			return
		}
		if response, err = json.Marshal(bBalancePayload); err != nil {
			log.Error("error websocket balanceUpdate json.Marshal: %v", err)
			return
		}
	} else if payload["e"] == "executionReport" {
		var (
			reportPayload  = &model.ExecutionReportPayload{}
			bReportPayload *binance.ExecutionReportPayload
		)
		err = json.Unmarshal(message, reportPayload)
		if err != nil {
			log.Error("error websocket HandlerMessage: %v", err)
			return
		}
		if bReportPayload, err = lib.ConvertExecutionReportPayload(reportPayload); err != nil {
			log.Error("error websocket executionReport convert: %v", err)
			return
		}
		if response, err = json.Marshal(bReportPayload); err != nil {
			log.Error("error websocket executionReport json.Marshal: %v", err)
			return
		}
	} else {
		response = message
	}

	if len(response) > 0 && c.handler != nil {
		(c.handler)(string(response))
	}
}

func (c *WebSocketClient) CheckToken(ids ...int) {
	var noConfigIds string
	for _, id := range ids {
		if conf.Conf.GetTokenInfoById(id) == nil {
			noConfigIds += strconv.Itoa(id) + ","
		}
	}
	if len(noConfigIds) > 0 {
		log.Info("CheckToken %v", noConfigIds)
		ec := &spot.Client{}
		ec.SetAppConfig(c.config)
		tokens, err := ec.TokenList(&model.TokenListParam{
			Ids: noConfigIds,
		})
		if err == nil && tokens.Success() && len(tokens.Data) > 0 {
			conf.Conf.AddTokens(tokens.Data)
		}
	}
}

func (c *WebSocketClient) SubscribeTrade(param *model.SubscribeTradeParam, handler func(string)) (err error) {
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = errors.New("lack config tokens")
		return
	}
	streamName := fmt.Sprintf(TradeStreamName, baseToken.Id, quoteToken.Id)
	c.handler = handler
	c.Subscribe([]string{streamName}, param.Id)
	return
}

func (c *WebSocketClient) SubscribeKline(param *model.SubscribeKlineParam, handler func(string)) (err error) {
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = errors.New("lack config tokens")
		return
	}
	interval := conf.KlineInterval[param.Interval]
	if interval == 0 {
		err = errors.New("illegal interval")
		return
	}
	streamName := fmt.Sprintf(KlineStreamName, baseToken.Id, quoteToken.Id, interval)
	c.handler = handler
	c.Subscribe([]string{streamName}, param.Id)
	return
}

func (c *WebSocketClient) SubscribeTicker(param *model.SubscribeTickerParam, handler func(string)) (err error) {
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = errors.New("lack config tokens")
		return
	}
	streamName := fmt.Sprintf(TickerStreamName, baseToken.Id, quoteToken.Id)
	c.handler = handler
	c.Subscribe([]string{streamName}, param.Id)
	return
}

func (c *WebSocketClient) SubscribeBookTicker(param *model.SubscribeBookTickerParam, handler func(string)) (err error) {
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = errors.New("lack config tokens")
		return
	}
	streamName := fmt.Sprintf(BookTickerStreamName, baseToken.Id, quoteToken.Id)
	c.handler = handler
	c.Subscribe([]string{streamName}, param.Id)
	return
}

func (c *WebSocketClient) SubscribeDepth(param *model.SubscribeDepthParam, handler func(string)) (err error) {
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = errors.New("lack config tokens")
		return
	}
	if param.Level != 5 && param.Level != 10 && param.Level != 20 {
		err = errors.New("illegal level")
		return
	}
	if param.Speed != 100 && param.Speed != 1000 {
		err = errors.New("illegal speed")
		return
	}
	var streamName string
	if param.Speed == 1000 {
		streamName = fmt.Sprintf(DepthStreamName, baseToken.Id, quoteToken.Id, param.Level, "")
	} else {
		streamName = fmt.Sprintf(DepthStreamName, baseToken.Id, quoteToken.Id, param.Level, "_100ms")
	}
	c.handler = handler
	c.Subscribe([]string{streamName}, param.Id)
	return
}

func (c *WebSocketClient) SubscribeDepthUpdate(param *model.SubscribeDepthUpdateParam, handler func(string)) (err error) {
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = errors.New("lack config tokens")
		return
	}
	if param.Speed != 100 && param.Speed != 1000 {
		err = errors.New("illegal speed")
		return
	}
	var streamName string
	if param.Speed == 1000 {
		streamName = fmt.Sprintf(DepthUpdateStreamName, baseToken.Id, quoteToken.Id, "")
	} else {
		streamName = fmt.Sprintf(DepthUpdateStreamName, baseToken.Id, quoteToken.Id, "_100ms")
	}
	c.handler = handler
	c.Subscribe([]string{streamName}, param.Id)
	return
}

func (c *WebSocketClient) SubscribeUserData(param *model.SubscribeUserDataParam, handler func(string)) {
	c.handler = handler
	c.SubscribeUser(param.ListenKey)
	return
}

func (c *WebSocketClient) ListSubscribes() (err error) {
	return
}
