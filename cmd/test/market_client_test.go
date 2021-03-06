package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"testing"

	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
	"github.com/degatedev/degate-sdk-golang/log"
)

func TestPing(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Ping()
	if err != nil {
		log.Print("TestPing error: %v", err)
		return
	}
	if response.Success() {
		log.Print("TestPing success \n %v", lib.String(response.Data))
	} else {
		log.Print("TestPing fail: %v", lib.String(response))
	}
}

func TestTrades(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetTrades(&model.TradeLastedParam{
		Symbol: "ETHUSDC",
		Limit:  0,
	})
	if err != nil {
		log.Print("error TestTrades: %v", err)
		return
	}
	if response.Success() {
		log.Print("TestTrades success \n %v", lib.String(response.Data))
	} else {
		log.Print("TestTrades fail: %v", lib.String(response))
	}
}

func TestHistoryTrades(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetHistoryTrades(&model.TradeHistoryParam{
		Symbol: "ETHUSDC",
		Limit:  20,
	})
	if err != nil {
		log.Print("error TradesHistory: %v", err)
		return
	}
	if response.Success() {
		log.Print("TradesHistory success \n %v", lib.String(response.Data))
	} else {
		log.Print("TradesHistory fail: %v", lib.String(response))
	}
}

func TestDepth(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetDepth(&model.DepthParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		log.Print("error Depth: %v", err)
		return
	}
	if response.Success() {
		log.Print("Depth success: %v", lib.String(response.Data))
	} else {
		log.Print("Depth fail: %v", lib.String(response))
	}
}

func TestKline(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetKlines(&model.KlineParam{
		Symbol: "ETHUSDC",
		Limit:  20,
	})
	if err != nil {
		log.Print("error Kline: %v", err)
		return
	}
	if response.Success() {
		log.Print("Kline success: %v", lib.String(response.Data))
	} else {
		log.Print("Kline fail: %v", lib.String(response))
	}
}

func TestTicker(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetTicker(&model.TickerParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		log.Print("error Ticker: %v", err)
		return
	}
	if response.Success() {
		log.Print("Ticker success: %v", lib.String(response.Data))
	} else {
		log.Print("Ticker fail: %v", lib.String(response))
	}
}

func TestTickerPrice(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     9,
		Symbol: "USDT",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetTickerPrice(&model.PairPriceParam{
		Symbol: "ETHUSDC,ETHUSDT",
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success: %v", lib.String(response.Data))
	} else {
		t.Logf("fail: %v", lib.String(response))
	}
}

func TestBookTicker(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetBookTicker(&model.BookTickerParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		log.Print("BookTicker error: %v", err)
		return
	}
	if response.Success() {
		log.Print("BookTicker success: %v", lib.String(response.Data))
	} else {
		log.Print("BookTicker fail: %v", lib.String(response))
	}
}