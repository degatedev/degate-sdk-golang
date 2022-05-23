package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"testing"

	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
	"github.com/degatedev/degate-sdk-golang/log"
)

func TestNewOrder(t *testing.T) {
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
	r, err := c.NewOrder(&model.OrderParam{
		Symbol:   "ETHUSDC",
		Side:     "SELL",
		Quantity: 1,
		Price:    2290,
		Type:     model.OrderTypeLimit,
	})
	if err != nil {
		log.Print("NewOrder error: %v", err)
		return
	}
	if r.Success() {
		log.Print("NewOrder success \n %v", lib.String(r.Data))
	} else {
		log.Print("NewOrder fail: %v", lib.String(r))
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
		log.Print("NewOrder error: %v", err)
		return
	}
	if r.Success() {
		log.Print("NewOrder success \n %v", lib.String(r.Data))
	} else {
		log.Print("NewOrder fail: %v", lib.String(r))
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
		Side:     "BUY",
		Quantity: 1,
	})
	if err != nil {
		log.Print("MarketOrder error: %v", err)
		return
	}
	if r.Success() {
		log.Print("MarketOrder success \n %v", lib.String(r.Data))
	} else {
		log.Print("MarketOrder fail: %v", lib.String(r))
	}
}

func TestCancelOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelOrder(&model.CancelOrderParam{
		OrderId: "188759499725809833483796742269",
	})
	if err != nil {
		log.Print("CancelOrder error: %v", err)
		return
	}
	if r.Success() {
		log.Print("CancelOrder success\n %v", lib.String(r.Data))
	} else {
		log.Print("CancelOrder fail: %v", lib.String(r))
	}
}

func TestCancelAllOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelAllOrders(false)
	if err != nil {
		log.Print("CancelAllOrders error: %v", err)
		return
	}
	if r.Success() {
		log.Print("CancelAllOrders success")
	} else {
		log.Print("CancelAllOrders fail: %v", lib.String(r))
	}
}

func TestCancelOrderOnChain(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.CancelOrderOnChain(&model.CancelOrderParam{
		OrderId: "188820328196434284078497267922",
	})
	if err != nil {
		log.Print("CancelOrderOnChain error: %v", err)
		return
	}
	if r.Success() {
		log.Print("CancelOrderOnChain success\n %v", lib.String(r.Data))
	} else {
		log.Print("CancelOrderOnChain fail: %v", lib.String(r))
	}
}
