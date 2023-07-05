package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"testing"

	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
)

func TestNewOrder(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Symbol: "ETH",
		Id:     0,
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Symbol: "USDC",
		Id:     2,
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.NewOrder(&model.OrderParam{
		Symbol:   "ETHUSDC",
		Side:     "buy",
		Quantity: 0.1,
		Type:     model.OrderTypeMarket,
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestPlaceOrder(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Symbol: "ETH",
		Id:     0,
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Symbol: "USDC",
		Id:     9,
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.PlaceOrder(&model.OrderParam{
		Symbol:   "ETHUSDC",
		Side:     "buy",
		Quantity: 0.1,
		Price:    1900,
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestMarketOrder(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     3,
		Symbol: "USDT",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     32,
		Symbol: "BTC",
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.MarketOrder(&model.OrderParam{
		Symbol: "BTCUSDT",
		Side:   "BUY",
		//Quantity: 0.1,
		QuoteOrderQty: 1000,
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestCancelOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelOrder(&model.CancelOrderParam{
		OrderId: "1773577844888250318179008512000",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestCancelAllOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelOpenOrders(false)
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("success")
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestCancelOrderOnChain(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelOrderOnChain(&model.CancelOrderParam{
		OrderId: "188820328196434284078497267922",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}
