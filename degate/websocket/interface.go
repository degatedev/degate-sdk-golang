package websocket

import "github.com/degatedev/degate-sdk-golang/degate/model"

type WsClient interface {
	SubscribeKline(param *model.SubscribeKlineParam, handler func(string)) (err error)
	SubscribeTrade(param *model.SubscribeTradeParam, handler func(string)) (err error)
	SubscribeTicker(param *model.SubscribeTickerParam, handler func(string)) (err error)
	SubscribeBookTicker(param *model.SubscribeBookTickerParam, handler func(string)) (err error)
	SubscribeDepth(param *model.SubscribeDepthParam, handler func(string)) (err error)
	SubscribeDepthUpdate(param *model.SubscribeDepthUpdateParam, handler func(string)) (err error)
	SubscribeUserData(param *model.SubscribeUserDataParam, handler func(string))
}
