package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	// dm "degate-backend/order-book-api/model"
	"testing"
	"time"

	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
	"github.com/degatedev/degate-sdk-golang/log"
)

func TestTime(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.Time()
	if err != nil {
		log.Print("Time error: %v", err)
		return
	}
	if r.Success() {
		log.Print("Time success\n: %v", time.Unix(int64(r.Data.ServerTime), 0))
	} else {
		log.Print("Time fail: %v", lib.String(r))
	}
}

func TestGasFee(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetGasFee()
	if err != nil {
		log.Print("GasFee error: %v", err)
		return
	}
	if r.Success() {
		log.Print("GasFee success\n: %v", lib.String(r.Data))
	} else {
		log.Print("GasFee fail: %v", lib.String(r))
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
		log.Print("GasFee error: %v", err)
		return
	}
	if r.Success() {
		log.Print("GasFee success\n: %v", lib.String(r.Data))
	} else {
		log.Print("GasFee fail: %v", lib.String(r))
	}
}

func TestTokens(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetTokens()
	if err != nil {
		log.Print("Tokens error: %v", err)
		return
	}
	if r.Success() {
		log.Print("Tokens success\n: %v", lib.String(r.Data))
	} else {
		log.Print("Tokens fail: %v", lib.String(r))
	}
}

func TestTokenListById(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetTokenList(&model.TokenListParam{
		Ids: "0,3,4,5",
	})
	if err != nil {
		log.Print("TokenList error: %v", err)
		return
	}
	if r.Success() {
		log.Print("TokenList success\n: %v", lib.String(r.Data))
	} else {
		log.Print("TokenList fail: %v", lib.String(r))
	}
}

func TestTokenListBySymbol(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetTokenList(&model.TokenListParam{
		Symbols: "ETH,USDC",
	})
	if err != nil {
		log.Print("TokenList error: %v", err)
		return
	}
	if r.Success() {
		log.Print("TokenList success\n: %v", lib.String(r.Data))
	} else {
		log.Print("TokenList fail: %v", lib.String(r))
	}
}

func TestExchangeInfo(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetExchangeInfo()
	if err != nil {
		log.Print("ExchangeInfo error: %v", err)
		return
	}
	if r.Success() {
		log.Print("ExchangeInfo success\n: %v", lib.String(r.Data))
	} else {
		log.Print("ExchangeInfo fail: %v", lib.String(r))
	}
}

func TestPairInfo(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetPair(&request.PairInfoRequest{
		Token1: 3,
		Token2: 5,
	})
	if err != nil {
		log.Print("GetPair error: %v", err)
		return
	}
	if r.Success() {
		log.Print("GetPair success\n: %v", lib.String(r.Data))
	} else {
		log.Print("GetPair fail: %v", lib.String(r))
	}
}
