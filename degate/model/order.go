package model

import (
	"strings"

	"github.com/degatedev/degate-sdk-golang/degate/binance"
)

const (
	OrderStatusOpen      = "open"
	OrderStatusCanceled  = "canceled"
	OrderStatusCompleted = "completed"

	OrderTypeLimit  = "LIMIT"
	OrderTypeMarket = "MARKET"

	OrderSideBuy  = "BUY"
	OrderSideSell = "SELL"
)

type OrderHeader struct {
	Source string `json:"source"`
}

type OrderParam struct {
	Symbol           string
	Side             string
	Type             string
	Quantity         float64
	QuoteOrderQty    float64
	Price            float64
	ValidUntil       int64
	NewOrderRespType string
	NewClientOrderId string `json:"newClientOrderId"`
	Source           string
}

type Order struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}

type OrderUpdate struct {
	OrderUpdateResult
	TransactionTime int64          `json:"transaction_time"`
	ClientOrderId   string         `json:"client_order_id"`
	Trades          []*TradeResult `json:"trades"`
}

func (o Order) GetID() string {
	return o.OrderId
}

type NewOrderAckResponse struct {
	binance.Response
	Data *Order `json:"data"`
}

type NewOrderResponse struct {
	binance.Response
	Data *OrderUpdate `json:"data"`
}

type StorageIdResponse struct {
	binance.Response
	Data StorageIdData `json:"data"`
}

type BatchStorageIdResponse struct {
	binance.Response
	Data []StorageIdData `json:"data"`
}

type StorageIdData struct {
	ID        string `json:"-"`
	StorageId uint64 `json:"storage_id"`
}

type CancelOrderParam struct {
	OrderId           string `json:"orderId"`
	OrigClientOrderId string `json:"origClientOrderId"`
	Fee               string `json:"fee"`
}

type CancelOrderResponse struct {
	binance.Response
	Data *Order `json:"data"`
}

type CancelAllOrderParam struct {
	AccountId      uint64 `json:"account_id" form:"account_id"`
	Timestamp      int64  `json:"timestamp" form:"timestamp"`
	IncludeGrid    bool   `json:"include_grid" form:"include_grid"`
	EDDSASignature string `json:"eddsa_signature" form:"eddsa_signature"`
}

type CancelAllOrderResponse struct {
	binance.Response
	Data *Order `json:"data"`
}

type OrderDetailParam struct {
	OrderId           string
	OrigClientOrderId string
}

type OrderDetailResponse struct {
	binance.Response
	Data *OrderList `json:"data"`
}

type OrdersParam struct {
	Symbol    string
	StartTime int64
	EndTime   int64
	Limit     int64
	OrderId   string
}

type OrdersParamWithTokenId struct {
	TokenId1  int32
	TokenId2  int32
	StartTime int64
	EndTime   int64
	Limit     int64
}

type OrderList struct {
	ID                      uint64                 `json:"id"`
	StorageID               uint32                 `json:"storage_id"`
	OrderID                 string                 `json:"order_id"`
	ClientOrderId           string                 `json:"client_order_id"`
	AccountID               uint64                 `json:"account_id"`
	PairID                  uint64                 `json:"pair_id"`
	SellToken               *binance.ShowTokenData `json:"sell_token"`
	BuyToken                *binance.ShowTokenData `json:"buy_token"`
	GasFeeToken             *binance.ShowTokenData `json:"gas_fee_token"`
	FeeToken                *binance.ShowTokenData `json:"fee_token"`
	MaxFeeBips              int                    `json:"max_fee_bips"`
	FilledSellTokenVolume   string                 `json:"filled_sell_token_volume"`
	FilledGasFeeTokenVolume string                 `json:"filled_gas_fee_token_volume"`
	FilledBuyTokenVolume    string                 `json:"filled_buy_token_volume"`
	FilledFeeTokenVolume    string                 `json:"filled_fee_token_volume"`
	Status                  string                 `json:"status"`
	IsSystemInvolved        bool                   `json:"is_system_involved"`
	FillAmountBOrs          bool                   `json:"fill_amount_bors"`
	ValidUntil              int64                  `json:"valid_until"`
	OrderType               int64                  `json:"order_type"`
	IsBuy                   bool                   `json:"is_buy"`
	Price                   string                 `json:"price"`
	CloseTime               int64                  `json:"close_time"`
	CreateTime              int64                  `json:"create_time"`
	UpdateTime              int64                  `json:"update_time"`
	GasFee                  string                 `json:"gas_fee"`
}

func (o *OrderList) GetSymbol() string {
	if o.IsBuy {
		return o.BuyToken.Symbol + o.SellToken.Symbol
	} else {
		return o.SellToken.Symbol + o.BuyToken.Symbol
	}
}

func (o *OrderList) Completed() bool {
	return strings.ToUpper(o.Status) == "COMPLETED"
}

type OrdersListData struct {
	binance.ListData
	Data []*OrderList `json:"list"`
}

type OrdersResponse struct {
	binance.Response
	Data OrdersListData `json:"data"`
}

type CancelAllOrdersParam struct {
	IncludeGrid bool `json:"includeGrid"`
}
