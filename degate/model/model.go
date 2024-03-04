package model

import "github.com/degatedev/degate-sdk-golang/degate/binance"

type TokenInfo struct {
	Id              int    `json:"id"`
	Chain           string `json:"chain"`
	Code            string `json:"code"`
	Symbol          string `json:"symbol"`
	Decimals        int32  `json:"decimals"`
	IsTrustedToken  bool   `json:"is_trusted_token"`
	IsQuotableToken bool   `json:"is_quotable_token"`
	IsGasToken      bool   `json:"is_gas_token"`
	IsListToken     bool   `json:"is_list_token"`
	Active          bool   `json:"active"`
	ShowDecimal     int32  `json:"show_decimal"`
	Priority        uint64 `json:"priority"`
}

type Ticker struct {
	BaseTokenID        uint64                 `json:"base_token_id"`
	QuoteTokenID       uint64                 `json:"quote_token_id"`
	PriceChange        string                 `json:"price_change"`
	PriceChangePercent string                 `json:"price_change_percent"`
	WeightedAvgPrice   string                 `json:"weighted_avg_price"`
	PrevClosePrice     string                 `json:"prev_close_price"`
	LastPrice          string                 `json:"last_price"`
	LastQty            string                 `json:"last_qty"`
	BidPrice           string                 `json:"bid_price"`
	BidQty             string                 `json:"bid_qty"`
	AskPrice           string                 `json:"ask_price"`
	AskQty             string                 `json:"ask_qty"`
	OpenPrice          string                 `json:"open_price"`
	HighPrice          string                 `json:"high_price"`
	LowPrice           string                 `json:"low_price"`
	Volume             string                 `json:"volume"`
	QuoteVolume        string                 `json:"quote_volume"`
	OpenTime           int64                  `json:"open_time"`
	CloseTime          int64                  `json:"close_time"`
	FirstId            string                 `json:"first_id"`
	LastId             string                 `json:"last_id"`
	Count              uint64                 `json:"count"`
	WeekHighPrice      string                 `json:"week_high_price"`
	WeekLowPrice       string                 `json:"week_low_price"`
	BaseTokenPrice     string                 `json:"base_token_price"`
	QuoteTokenPrice    string                 `json:"quote_token_price"`
	MakerFee           string                 `json:"maker_fee"`
	TakerFee           string                 `json:"taker_fee"`
	PairId             uint64                 `json:"pair_id"`
	BaseToken          *binance.ShowTokenData `json:"base_token"`
	QuoteToken         *binance.ShowTokenData `json:"quote_token"`
}

func (t *Ticker) ToBookTicker() (bookTicker *BookTicker) {
	bookTicker = &BookTicker{
		BaseTokenID:  t.BaseTokenID,
		QuoteTokenID: t.QuoteTokenID,
		BidPrice:     t.BidPrice,
		BidQty:       t.BidQty,
		AskPrice:     t.AskPrice,
		AskQty:       t.AskQty,
	}
	return
}

type BookTicker struct {
	BaseTokenID  uint64 `json:"base_token_id"`
	QuoteTokenID uint64 `json:"quote_token_id"`
	BidPrice     string `json:"bid_price"`
	BidQty       string `json:"bid_qty"`
	AskPrice     string `json:"ask_price"`
	AskQty       string `json:"ask_qty"`
}

type ExchangeInfo struct {
	ChainID               int64                      `json:"chain_id"`
	ExchangeAddress       string                     `json:"exchange_address"`
	DepositAddress        string                     `json:"deposit_address"`
	WithdrawalsAddress    string                     `json:"withdrawals_address"`
	SpotTradeAddress      string                     `json:"spot_trade_address"`
	OrderCancelAddress    string                     `json:"order_cancel_address"`
	OrderEffectiveDigits  int                        `json:"order_effective_digits"`
	Timezone              string                     `json:"timezone"`
	ServerTime            int64                      `json:"server_time"`
	OrderMaxVolume        string                     `json:"order_max_volume"`
	RateLimits            []*binance.RateLimitFilter `json:"rate_limits"`
	MinLimitOrderUSDValue float64                    `json:"min_limit_order_usd_value"`
}
