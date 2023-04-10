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
