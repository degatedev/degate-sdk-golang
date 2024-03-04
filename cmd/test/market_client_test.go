package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"testing"

	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
)

func TestPing(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Ping()
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestDepth(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     2,
		Symbol: "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Depth(&model.DepthParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestTrades(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     2,
		Symbol: "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Trades(&model.TradeLastedParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestTradesHistory(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     2,
		Symbol: "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.TradesHistory(&model.TradeHistoryParam{
		Symbol: "ETHUSDC",
		Limit:  3,
		FromId: "",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestKline(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     2,
		Symbol: "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Klines(&model.KlineParam{
		Symbol: "ETHUSDC",
		Limit:  20,
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestTicker24(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     9,
		Symbol: "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Ticker24(&model.TickerParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestTickerPrice(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     2,
		Symbol: "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.TickerPrice(&model.PairPriceParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestBookTicker(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     2,
		Symbol: "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.BookTicker(&model.BookTickerParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}
