package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"testing"
	"time"

	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
	"github.com/degatedev/degate-sdk-golang/degate/websocket"
	"github.com/degatedev/degate-sdk-golang/log"
)

func sleep() {
	for {
		time.Sleep(time.Minute * 10)
	}
}

func TestKlineWebscoket(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	c := &websocket.WebSocketClient{}
	c.Init(appConfig)
	err := c.SubscribeKline(&model.SubscribeKlineParam{
		Symbol:   "ETHUSDC",
		Interval: "1m",
	}, func(message string) {
		log.Print(message)
	})
	if err != nil {
		return
	}
	sleep()
}

func TestTradeWebscoket(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
		Decimals: 18,
	})
	c := websocket.WebSocketClient{}
	c.Init(appConfig)
	err := c.SubscribeTrade(&model.SubscribeTradeParam{
		Symbol: "ETHUSDC",
	}, func(message string) {
		log.Print(message)
	})
	if err != nil {
		return
	}
	sleep()
}

func TestTickerWebscoket(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	c := websocket.WebSocketClient{}
	c.Init(appConfig)
	err := c.SubscribeTicker(&model.SubscribeTickerParam{
		Symbol: "ETHUSDC",
	}, func(message string) {
		log.Print(message)
	})
	if err != nil {
		return
	}
	sleep()
}

func TestBookTickerWebscoket(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	c := websocket.WebSocketClient{}
	c.Init(appConfig)
	err := c.SubscribeBookTicker(&model.SubscribeBookTickerParam{
		Symbol: "ETHUSDC",
	}, func(message string) {
		log.Print(message)
	})
	if err != nil {
		return
	}
	sleep()
}

func TestDepthWebscoket(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	c := websocket.WebSocketClient{}
	c.Init(appConfig)
	err := c.SubscribeDepth(&model.SubscribeDepthParam{
		Symbol: "ETHUSDC",
		Level:  5,
		Speed:  100,
	}, func(message string) {
		log.Print(message)
	})
	if err != nil {
		return
	}
	sleep()
}

func TestDepthUpdateWebscoket(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
		Decimals: 18,
	})
	c := websocket.WebSocketClient{}
	c.Init(appConfig)
	err := c.SubscribeDepthUpdate(&model.SubscribeDepthUpdateParam{
		Symbol: "ETHUSDC",
		Speed:  100,
	}, func(message string) {
		log.Print(message)
	})
	if err != nil {
		return
	}
	sleep()
}

func TestUserDataWebscoket(t *testing.T) {
	sc := new(spot.Client)
	sc.SetAppConfig(appConfig)
	keyResponse, err := sc.NewListenKey()
	if err != nil {
		return
	}
	if !keyResponse.Success() {
		log.Print("fail get listenKey code:%v message:%v", keyResponse.Code, keyResponse.Message)
		return
	}
	c := websocket.WebSocketClient{}
	c.Init(appConfig)
	c.SubscribeUserData(&model.SubscribeUserDataParam{
		ListenKey: keyResponse.Data.ListenKey,
	}, func(message string) {
		log.Print(message)
	})
	sleep()
}
