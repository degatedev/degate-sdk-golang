package model

type TradeResult struct {
	TradeId      string       `json:"trade_id"`
	BuyerIsMaker bool         `json:"buyer_is_maker"`
	TxType       string       `json:"tx_type"`
	FillSA       string       `json:"fills_a"`
	FillSB       string       `json:"fills_b"`
	OrderA       *OrderDetail `json:"order_a"`
	OrderB       *OrderDetail `json:"order_b"`
	PairId       uint64       `json:"pair_id"`
	IsStatble    bool         `json:"is_stable"`
	Price        string       `json:"price"`
	TradeTime    int64        `json:"trade_time"`
}

type OrderDetail struct {
	OrderId                      string `json:"order_id"`
	AccountId                    uint32 `json:"account_id"`
	StorageId                    uint32 `json:"storage_id"`
	TokenS                       uint32 `json:"token_s"`
	TokenB                       uint32 `json:"token_b"`
	AmountS                      string `json:"amount_s"`
	AmountB                      string `json:"amount_b"`
	ValidUntil                   uint64 `json:"valid_until"`
	FillAmountBorS               bool   `json:"fill_amount_bor_s"`
	FeeBips                      int    `json:"fee_bips"`
	FeeTokenId                   uint32 `json:"fee_token_id"`
	Fee                          string `json:"fee"`
	MaxFee                       string `json:"max_fee"`
	TradeFee                     string `json:"trade_fee"`
	UiReferId                    uint64 `json:"ui_refer_id"`
	AutoMarketOrder              int    `json:"auto_market_order"`
	PreAutoMarketOrderStorageId  int    `json:"pre_auto_market_order_storage_id"`
	AutoMarketStorageId          int    `json:"auto_market_storage_id"`
	NextAutoMarketOrderStorageId int    `json:"next_auto_market_order_storage_id"`
	OrderType    uint16 `json:"order_type"`
	MaxLevel     uint16 `json:"max_level"`
	Level        uint16 `json:"level"`
	GridOffset   string `json:"grid_offset"`
	OrderOffset  string `json:"order_offset"`
	StartOrderId string `json:"start_order_id"`
}

type OrderToken struct {
	TokenId      uint32 `json:"token_id"`
	Volume       string `json:"volume"`
	FilledVolume string `json:"filled_volume"`
}

type OrderUpdateResult struct {
	Status         string      `json:"status"`
	AccountId      uint32      `json:"account_id"`
	StorageId      uint32      `json:"storage_id"`
	OrderId        string      `json:"order_id"`
	SellToken      *OrderToken `json:"sell_token"`
	BuyToken       *OrderToken `json:"buy_token"`
	FrozenVolume   string      `json:"frozen_volume"`
	ValidUntil     uint64      `json:"valid_until"`
	FeeBips        int         `json:"fee_bips"`
	TradingFee     string      `json:"trading_fee"`
	Fee            string      `json:"fee"`
	MaxFee         string      `json:"max_fee"`
	FeeTokenId     uint32      `json:"fee_token_id"`
	UiReferId      uint64      `json:"ui_refer_id"`
	IsBuy          bool        `json:"is_buy"`
	FillAmountBors bool        `json:"fill_amount_bors"`
	PairId         uint64      `json:"pair_id"`
	OrderType      uint16      `json:"order_type"`
	GridId         uint32      `json:"grid_id"`
	GridCreateTs   uint32      `json:"grid_create_ts"`
	Level          uint32      `json:"level"`
	IsGridFlip     bool        `json:"is_grid_flip"`
}
