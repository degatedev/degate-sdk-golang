package model

import (
	"github.com/degatedev/degate-sdk-golang/degate/binance"
)

type TimeData struct {
	Timestamp int64 `json:"timestamp"`
}

type TimeResponse struct {
	binance.Response
	Data *TimeData `json:"data"`
}

type GasFeeResponse struct {
	binance.Response
	Data *GasFees `json:"data"`
}

type ExchangeInfoResponse struct {
	binance.Response
	Data *ExchangeInfo `json:"data"`
}

type TokensResponse struct {
	binance.Response
	Data []*TokenInfo `json:"data"`
}

type AccessTokenParam struct {
	Owner          string `json:"owner" form:"owner" binding:"required"`
	Time           int64  `json:"time" form:"time" binding:"required"`
	EDDSASignature string `json:"eddsa_signature" form:"eddsa_signature"`
	UseTradeKey    bool   `json:"use_trade_key" form:"use_trade_key"`
}

type AccessTokenResponse struct {
	binance.Response
	Data *AccessToken `json:"data"`
}
