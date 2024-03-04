package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"testing"
	"time"

	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
)

func TestTime(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.Time()
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", time.Unix(int64(r.Data.ServerTime), 0))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestGasFee(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GasFee()
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestTradeFee(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetTradeFee(&model.TradeFeeParam{
		Symbol: "ETH",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestTokens(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetTokens()
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestTokenListById(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.TokenList(&model.TokenListParam{
		Ids: "0,3,4,5",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestTokenListBySymbol(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.TokenList(&model.TokenListParam{
		Symbols: "eth,usdc",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestExchangeInfo(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.ExchangeInfo()
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestPairInfo(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetPair(&request.PairInfoRequest{
		Token1: 0,
		Token2: 9,
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}
