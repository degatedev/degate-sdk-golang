package binance

type TimeResponse struct {
	Response
	Data *Time `json:"data"`
}

type Time struct {
	ServerTime int `json:"serverTime"`
}

type PingResponse struct {
	Response
	Data struct{} `json:"data"`
}

type ExchangeInfoResponse struct {
	Response
	Data *ExchangeInfo `json:"data"`
}

type ExchangeInfo struct {
	ChainID               int64              `json:"chain_id"`
	Timezone              string             `json:"timezone"`
	ServerTime            int64              `json:"serverTime"`
	RateLimits            []*RateLimitFilter `json:"rateLimits"`
	MinLimitOrderUSDValue float64            `json:"minLimitOrderUSDValue"`
}

type RateLimitFilter struct {
	Interval      string `json:"interval"`    // : "MINUTE"
	IntervalNum   int    `json:"intervalNum"` // : 1
	Limit         int    `json:"limit"`       // : 2400
	RateLimitType string `json:"rateLimitType"`
}

type GasFeeResponse struct {
	Response
	Data *OffChainFee `json:"data"`
}

type OffChainFee struct {
	UpdateAccountGasFees       []*GasFee `json:"update_account_gas_fees"`
	WithdrawalGasFees          []*GasFee `json:"withdrawal_gas_fees"`
	EstimatedWithdrawalGasFees []*GasFee `json:"estimated_withdrawal_gas_fees"`
	TransferGasFees            []*GasFee `json:"transfer_gas_fees"`
	TransferNoIDGasFees        []*GasFee `json:"transfer_no_id_gas_fees"`
	OrderGasFees               []*GasFee `json:"order_gas_fees"`
	AddPairGasFees             []*GasFee `json:"add_pair_gas_fees"`
	MiningGasFees              []*GasFee `json:"mining_gas_fees"`
	OnChainCancelOrderGasFees  []*GasFee `json:"on_chain_cancel_order_gas_fees"`
	OrderGasMultiple           uint32    `json:"order_gas_multiple"`
}

type GasFee struct {
	Symbol   string `json:"symbol"`
	TokenId  uint64 `json:"token_id"`
	Volume   string `json:"-"`
	Quantity string `json:"quantity"`
	Decimals int32  `json:"decimals"`
}
