package model

type Balances struct {
	ID                    uint64         `json:"-"`
	AccountID             uint64         `json:"account_id"`
	Token                 *ShowTokenData `json:"token"`
	Balance               string         `json:"balance"`
	FrozenDepositBalance  string         `json:"frozen_deposit_balance"`
	FrozenOrderBalance    string         `json:"frozen_order_balance"`
	FrozenWithdrawBalance string         `json:"frozen_withdraw_balance"`
}

type DepositData struct {
	ID         uint64         `json:"id"`
	AccountID  uint64         `json:"account_id"`
	Owner      string         `json:"owner"`
	Token      *ShowTokenData `json:"token"`
	L2TrxID    string         `json:"l2_trx_id"`
	Status     string         `json:"status"`
	Nonce      uint64         `json:"nonce"`
	CreateTime int64          `json:"create_time"`
	UpdateTime int64          `json:"update_time"`
}

type WithdrawalData struct {
	ID          int64          `json:"-"`
	AccountID   uint64         `json:"account_id"`
	WithdrawID  string         `json:"withdraw_id"`
	ToAccountID int64          `json:"to_account_id"`
	ToAddress   string         `json:"to_address"`
	Token       *ShowTokenData `json:"token"`
	FeeToken    *ShowTokenData `json:"fee_token"`
	TxHash      string         `json:"tx_hash"`
	Status      string         `json:"status"`
	Type        int            `json:"type"`
	CreateTime  int64          `json:"create_time"`
	UpdateTime  int64          `json:"update_time"`
}

type TransfersDataDetail struct {
	ID          int64          `json:"-"`
	AccountID   uint64         `json:"account_id"`
	Address     string         `json:"address"`
	TransferID  string         `json:"transfer_id"`
	ToAccountID int64          `json:"to_account_id"`
	ToAddress   string         `json:"to_address"`
	Token       *ShowTokenData `json:"token"`
	FeeToken    *ShowTokenData `json:"fee_token"`
	CreateTime  int64          `json:"create_time"`
	UpdateTime  int64          `json:"update_time"`
}

type TradeData struct {
	ID                uint64         `json:"id"`
	AccountID         int64          `json:"account_id"`
	TradeId           string         `json:"trade_id"`
	OrderId           string         `json:"order_id"`
	PairId            uint64         `json:"pair_id"`
	IsMaker           bool           `json:"is_maker"`
	FilledSellToken   *ShowTokenData `json:"filled_sell_token"`
	FilledBuyToken    *ShowTokenData `json:"filled_buy_token"`
	FilledFeeToken    *ShowTokenData `json:"filled_fee_token"`
	FilledGasFeeToken *ShowTokenData `json:"filled_gas_fee_token"`
	CreateTime        int64          `json:"create_time"`
	FillAmountBors    bool           `json:"fill_amount_bors"`
	IsBuy             bool           `json:"is_buy"`
	Price string `json:"price"`
}

type TransferData struct {
	Hash         string `json:"hash"`
	OrderID      string `json:"order_id"`
	Status       string `json:"status"`
	IsIdempotent bool   `json:"is_idempotent"`
}

type PriceData struct {
	TokenID  uint64 `json:"token_id"`
	Price    string `json:"price"`
	Decimals int32  `json:"decimals"`
}

type InitiateWithdrawal struct {
	Hash         string `json:"hash"`
	OrderID      string `json:"order_id"`
	Status       string `json:"status"`
	IsIdempotent bool   `json:"is_idempotent"`
}

type OffChainFee struct {
	UpdateAccountGasFees      []*ShowTokenData `json:"update_account_gas_fees"`
	WithdrawalGasFees         []*ShowTokenData `json:"withdrawal_gas_fees"`
	WithdrawalOtherGasFees    []*ShowTokenData `json:"withdrawal_other_gas_fees"`
	TransferGasFees           []*ShowTokenData `json:"transfer_gas_fees"`
	TransferNoIDGasFees       []*ShowTokenData `json:"transfer_no_id_gas_fees"`
	OrderGasFees              []*ShowTokenData `json:"order_gas_fees"`
	AddPairGasFees            []*ShowTokenData `json:"add_pair_gas_fees"`
	MiningGasFees             []*ShowTokenData `json:"mining_gas_fees"`
	OnChainCancelOrderGasFees []*ShowTokenData `json:"on_chain_cancel_order_gas_fees"`
	UpdateAccountGas          string           `json:"update_account_gas"`
	WithdrawalGas             string           `json:"withdrawal_gas"`
	WithdrawalOtherGas        string           `json:"withdrawal_other_gas"`
	TransferGas               string           `json:"transfer_gas"`
	OrderGas                  string           `json:"order_gas"`
	AddPairGas                string           `json:"add_pair_gas"`
	MiningGas                 string           `json:"mining_gas"`
	TransferNoIDGas           string           `json:"transfer_no_id_gas"`
	UpdateAccountGasUsd       string           `json:"update_account_gas_usd"`
	WithdrawalGasUsd          string           `json:"withdrawal_gas_usd"`
	WithdrawalOtherGasUsd     string           `json:"withdrawal_other_gas_usd"`
	TransferGasUsd            string           `json:"transfer_gas_usd"`
	OrderGasUsd               string           `json:"order_gas_usd"`
	AddPairGasUsd             string           `json:"add_pair_gas_usd"`
	MiningGasUsd              string           `json:"mining_gas_usd"`
	TransferNoIDUsd           string           `json:"transfer_no_id_usd"`
	GasPrice                  string           `json:"gas_price"`
	OrderGasMultiple          uint32           `json:"order_gas_multiple"`
}
