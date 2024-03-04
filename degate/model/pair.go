package model

import "github.com/degatedev/degate-sdk-golang/degate/binance"

type PairInfoResponse struct {
	binance.Response
	Data *PairInfo `json:"data"`
}

type PairInfo struct {
	PairID     uint64                 `json:"pair_id"`
	BaseToken  *binance.ShowTokenData `json:"base_token"`
	QuoteToken *binance.ShowTokenData `json:"quote_token"`
	IsStable   bool                   `json:"is_stable"`
}
