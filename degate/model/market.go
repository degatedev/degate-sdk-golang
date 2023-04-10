package model

import (
	"github.com/degatedev/degate-sdk-golang/degate/binance"
)

type KlineParam struct {
	Symbol    string
	Interval  string
	StartTime int64
	EndTime   int64
	Limit     int64
}

type TradeLastedParam struct {
	Symbol string
	Limit  int64
}

type TradeHistoryParam struct {
	Symbol string
	Limit  int64
	FromId string
}

type TradesParam struct {
	Symbol string
	Limit  int64
	FromId string
}

type DepthParam struct {
	Symbol string
	Limit  int64
}

type DepthResponse struct {
	binance.Response
	Data *Depth `json:"data"`
}

type TickerParam struct {
	Symbol string
}

type TickerResponse struct {
	binance.Response
	Data *Ticker
}

type PairPriceParam struct {
	Symbol string
}

type Depth struct {
	LastUpdateID int64      `json:"last_update_id"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
	QuoteTokenID uint64     `json:"quote_token_id"`
	BaseTokenID  uint64     `json:"base_token_id"`
	Depth        int64      `json:"depth"`
	Label        string     `json:"label"`
}

type PairPriceResponse struct {
	binance.Response
	Data []*PairsPricesRes
}

type PairsPricesRes struct {
	PairID  uint64 `json:"pair_id"`
	Price   string `json:"price"`
	Percent string `json:"percent"`
}

type BookTickerParam struct {
	Symbol string
}

type BookTickerResponse struct {
	binance.Response
	Data *BookTickerData
}

type BookTickerData struct {
	LastUpdateID int64      `json:"last_update_id"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}
