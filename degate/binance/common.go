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
	ChainID              int64   `json:"chain_id"`
	ExchangeAddress      string  `json:"exchange_address"`
	DepositAddress       string  `json:"deposit_address"`
	WithdrawalsAddress   string  `json:"withdrawals_address"`
	SpotTradeAddress     string  `json:"spot_trade_address"`
	OrderCancelAddress   string  `json:"order_cancel_address"`
	OrderEffectiveDigits int     `json:"order_effective_digits"`
	MinOrderPrice        float64 `json:"min_order_price"`
	MaxFeeBipsMax        int64   `json:"max_fee_bips_max"`
	Timezone             string  `json:"timezone"`
	ServerTime           int64   `json:"serverTime"`
	OrderMaxVolume       string  `json:"order_max_volume"`
}

type GasFeeResponse struct {
	Response
	Data *OffChainFee `json:"data"`
}

type OffChainFee struct {
	UpdateAccountGasFees      []*GasFee `json:"update_account_gas_fees"`
	WithdrawalGasFees         []*GasFee `json:"withdrawal_gas_fees"`
	WithdrawalOtherGasFees    []*GasFee `json:"withdrawal_other_gas_fees"`
	TransferGasFees           []*GasFee `json:"transfer_gas_fees"`
	TransferNoIDGasFees       []*GasFee `json:"transfer_no_id_gas_fees"`
	OrderGasFees              []*GasFee `json:"order_gas_fees"`
	AddPairGasFees            []*GasFee `json:"add_pair_gas_fees"`
	MiningGasFees             []*GasFee `json:"mining_gas_fees"`
	OnChainCancelOrderGasFees []*GasFee `json:"on_chain_cancel_order_gas_fees"`
	OrderGasMultiple          uint32    `json:"order_gas_multiple"`
}

type GasFee struct {
	Symbol   string `json:"symbol"`
	TokenId  uint64 `json:"token_id"`
	Volume   string `json:"-"`
	Quantity string `json:"quantity"`
	Decimals int32  `json:"decimals"`
}
