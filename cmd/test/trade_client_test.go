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
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     8,
		Symbol: "USDC",
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.NewOrder(&model.OrderParam{
		Symbol:   "ETHUSDC",
		Side:     "SELL",
		Quantity: 0.1,
		Price:    5000,
		Type:     model.OrderTypeMarket,
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("success: %v", lib.String(r.Data))
	} else {
		t.Logf("fail: %v", lib.String(r))
	}
}

func TestPlaceOrder(t *testing.T) {
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
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.PlaceOrder(&model.OrderParam{
		Symbol:   "ETHUSDC",
		Side:     "SELL",
		Quantity: 1,
		Price:    5000,
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("success: %v", lib.String(r.Data))
	} else {
		t.Logf("fail: %v", lib.String(r))
	}
}

func TestMarketOrder(t *testing.T) {
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
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.MarketOrder(&model.OrderParam{
		Symbol:   "ETHUSDC",
		Side:     "SELL",
		Quantity: 0.05,
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("success: %v", lib.String(r.Data))
	} else {
		t.Logf("fail: %v", lib.String(r))
	}
}

func TestCancelOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelOrder(&model.CancelOrderParam{
		OrderId: "196120279748861022004248412422382",
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("success: %v", lib.String(r.Data))
	} else {
		t.Logf("fail: %v", lib.String(r))
	}
}

func TestCancelAllOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelAllOrders(false)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("success")
	} else {
		t.Logf("fail: %v", lib.String(r))
	}
}

func TestCancelOrderOnChain(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelOrderOnChain(&model.CancelOrderParam{
		OrderId: "188820328196434284078497267922",
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("success: %v", lib.String(r.Data))
	} else {
		t.Logf("fail: %v", lib.String(r))
	}
}
