package binance

type OrdersResponse struct {
	Response
	Data []*Order `json:"data"`
}

type OrderResponse struct {
	Response
	Data *Order `json:"data"`
}

type OrderCancelResponse struct {
	Response
	Data *Order `json:"data"`
}

type NewOrderResponse struct {
	Response
	Data interface{} `json:"data"`
}

type OrderAck struct {
	Symbol        string `json:"symbol"`
	OrderId       string `json:"orderId"`
	ClientOrderId string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
}

type OrderResult struct {
	OrderAck
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	Fills               []*OrderFill
}

type OrderFill struct {
	Price                 string `json:"price"`
	Qty                   string `json:"qty"`
	Commission            string `json:"commission"`
	CommissionAsset       string `json:"commissionAsset"`
	GasFeeCommission      string `json:"gasFeeCommission"`
	GasFeeCommissionAsset string `json:"gasFeeCommissionAsset"`
}

type Order struct {
	Symbol              string `json:"symbol"`
	OrderId             string `json:"orderId"`
	ClientOrderId       string `json:"clientOrderId"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	Price               string `json:"price"`
	Status              string `json:"status"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	TimeInForce         string `json:"timeInForce"`
	IsWorking           bool   `json:"isWorking"`
}

type TradeResponse struct {
	Response
	Data []*UserTrade `json:"data"`
}

type UserTrade struct {
	Trade
	BuyOrderId  string `json:"buy_order_id"`
	SellOrderId string `json:"sell_order_id"`
}

type Trade struct {
	Symbol          string `json:"symbol"`
	Id              string `json:"id"`
	PairId          uint64 `json:"pair_id"`
	OrderId         string `json:"orderId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
	BaseTokenId     uint32 `json:"base_token_id"`
	QuoteTokenId    uint32 `json:"quote_token_id"`
	AccountId       uint64 `json:"account_id"`
	UserFlag        uint32 `json:"user_flag"`
	GasFee          string `json:"gas_fee"`
	TradeFee        string `json:"trade_fee"`
	Status          int    `json:"status"`
	ClientOrderId   string `json:"client_order_id"`
}
