package request

type StorageIdRequest struct {
	AccountId uint32 `json:"account_id" form:"account_id"`
	Owner     string `json:"-"`
	TokenId   uint32 `json:"token_id" form:"token_id"`
	Window    uint32 `json:"window" form:"window"`
}

type NewOrderRequest struct {
	AccountId        uint32 `json:"account_id"`
	StorageId        uint64 `json:"storage_id"`
	OrderID          string `json:"order_id"`
	SellToken        Token  `json:"sell_token"`
	BuyToken         Token  `json:"buy_token"`
	ValidUntil       int64  `json:"valid_until"`
	FillAmountBOrs   bool   `json:"fill_amount_bors"`
	FeeToken         Token  `json:"fee_token"`
	EDDSASignature   string `json:"eddsa_signature"`
	UiReferrerId     uint64 `json:"ui_referrer_id"`
	NewOrderRespType string `json:"new_order_resp_type"`
}

type CancelOrderRequest struct {
	AccountId      uint32 `json:"account_id"`
	OrderId        string `json:"order_id"`
	EDDSASignature string `json:"eddsa_signature"`
}

type OrderDetailRequest struct {
	OrderId string `json:"order_id"`
}

type OrdersRequest struct {
	AccountId uint32 `json:"account_id"`
	Token1    int32  `json:"token_1"`
	Token2    int32  `json:"token_2"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	Side      string `json:"side"`
	Status    string `json:"status"`
	Limit     int64  `json:"limit"`
	Offset    int64  `json:"offset"`
	OrderId   string `json:"orderId"`
}
