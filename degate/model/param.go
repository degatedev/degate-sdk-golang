package model

type DGDepositsParam struct {
	AccountId int64  `json:"account_id" form:"account_id"`
	Tokens    string `json:"tokens" form:"tokens"`
	Start     int64  `json:"start"  form:"start"`
	End       int64  `json:"end" form:"end"`
	Status    string `json:"status" form:"status"`
	Limit     int64  `json:"limit" form:"limit"`
	Offset    int64  `json:"offset" form:"offset"`
	TxId      string `json:"txId" form:"txId"`
}

type WithdrawalParam struct {
	AccountId int64  `json:"account_id" form:"account_id"`
	Start     int64  `json:"start"  form:"start"`
	End       int64  `json:"end" form:"end"`
	Status    string `json:"status" form:"status"`
	Tokens    string `json:"tokens" form:"tokens"`
	Limit     int64  `json:"limit" form:"limit"`
	Offset    int64  `json:"offset" form:"offset"`
	Type      string `json:"type" form:"type"`
}

type TradesUserParam struct {
	AccountId int64  `json:"account_id"`
	Token1    int64  `json:"token_1"`
	Token2    int64  `json:"token_2"`
	OrderID   string `json:"order_id"`
	Offset    int64  `json:"offset"`
	Limit     int64  `json:"limit"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	FromId    string `json:"from_id" `
}

type OrderDetailsParam struct {
	OrderID       string `json:"order_id" form:"order_id"`
	ClientOrderID string `json:"client_order_id" form:"client_order_id"`
}

type TokenListParam struct {
	Ids     string `json:"ids" form:"ids"`
	Symbols string `json:"symbols" form:"symbols"`
}

type DGTradesParam struct {
	Token1 uint64 `json:"token_1" form:"token_1"`
	Token2 uint64 `json:"token_2" form:"token_2"`
	Offset int64  `json:"offset" form:"offset"`
	Limit  int64  `json:"limit" form:"limit"`
	Start  int64  `json:"start"  form:"start"`
	End    int64  `json:"end" form:"end"`
	FromId string `json:"from_id" form:"from_id"`
}

type DGDepthParam struct {
	QuoteTokenID uint64 `json:"quote_token_id" form:"quote_token_id"`
	BaseTokenID  uint64 `json:"base_token_id" form:"base_token_id"`
	Size         int64  `json:"size" form:"size"`
	Depth        int64  `json:"depth" form:"depth"`
}

type KlinesParam struct {
	QuoteTokenID uint64 `json:"quote_token_id" form:"quote_token_id"`
	BaseTokenID  uint64 `json:"base_token_id" form:"base_token_id"`
	Start        int64  `json:"start" form:"start"`
	End          int64  `json:"end" form:"end"`
	Granularity  int    `json:"granularity" form:"granularity"`
	Limit        int    `json:"limit" form:"limit"`
}

type DGTickerParam struct {
	QuoteTokenID uint64 `json:"quote_token_id" form:"quote_token_id"`
	BaseTokenID  uint64 `json:"base_token_id" form:"base_token_id"`
}

type DGCancelOrderParam struct {
	AccountId      uint64 `json:"account_id" form:"account_id"`
	OrderId        string `json:"order_id" form:"order_id"`
	ClientOrderId  string `json:"client_order_id" form:"client_order_id"`
	EDDSASignature string `json:"eddsa_signature" form:"eddsa_signature"`
	FeeToken       Token  `json:"fee_token" form:"fee_token"`
	Key            string `json:"-"`
}

type DGCancelAllParam struct {
	AccountId      uint64 `json:"account_id" form:"account_id"`
	EDDSASignature string `json:"eddsa_signature" form:"eddsa_signature"`
	IncludeGrid    bool   `json:"include_grid" form:"include_grid"`
	Timestamp      int64  `json:"timestamp" form:"timestamp"`
	Key            string `json:"-"`
}

type CancelOrderParamWithTokenId struct {
	AccountId      uint64 `json:"account_id" form:"account_id"`
	OrderId        string `json:"order_id" form:"order_id"`
	ClientOrderId  string `json:"client_order_id" form:"client_order_id"`
	EDDSASignature string `json:"eddsa_signature" form:"eddsa_signature"`
	Key            string `json:"-"`
	SellTokenId    uint64 `json:"sell_token_id" form:"sell_token_id"`
	FeeTokenVol    string `json:"fee_token_vol" form:"fee_token_vol"`
	FeeTokenId     uint64 `json:"fee_token_id" form:"fee_token_id"`
}

type DGOrderParam struct {
	AccountId        uint64 `json:"account_id" form:"account_id"`
	OrderID          string `json:"order_id" form:"order_id"`
	StorageId        uint64 `json:"storage_id" form:"storage_id"`
	SellToken        Token  `json:"sell_token" form:"sell_token"`
	BuyToken         Token  `json:"buy_token" form:"buy_token"`
	ValidUntil       int64  `json:"valid_until" form:"valid_until"`
	FeeToken         Token  `json:"fee_token" form:"fee_token"`
	EDDSASignature   string `json:"eddsa_signature" form:"eddsa_signature"`
	UiReferrerId     uint64 `json:"ui_referrer_id" form:"ui_referrer_id" `
	FillAmountBOrs   bool   `json:"fill_amount_bors" form:"fill_amount_bors"`
	NewOrderRespType string `json:"new_order_resp_type"`
	Type             uint8  `json:"type" form:"type"`
	GridOffset       string `json:"grid_offset" form:"grid_offset"`
	OrderOffset      string `json:"order_offset" form:"order_offset"`
	MaxLevel         uint8  `json:"max_level" form:"max_level"`
	Key              string `json:"-"`
	ClientOrderId    string `json:"client_order_id"`
	Price            string `json:"price"`
	MustProfit       bool   `json:"must_profit"`
	PostOnly         bool   `json:"post_only"`
	FeeBips          string `json:"fee_bips"`
}

type Token struct {
	TokenId uint64 `json:"token_id" form:"token_id"`
	Volume  string `json:"volume"  form:"volume"`
}

type DGPairPriceParam struct {
	Pairs string `json:"pairs"`
}

type WithdrawalGasParam struct {
	To      string `json:"to" form:"to" `            // account address
	TokenId uint64 `json:"token_id" form:"token_id"` // token id
}
