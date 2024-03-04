package model

import "github.com/degatedev/degate-sdk-golang/degate/binance"

type Balances struct {
	ID                    uint64                 `json:"-"`
	AccountID             uint64                 `json:"account_id"`
	Token                 *binance.ShowTokenData `json:"token"`
	Balance               string                 `json:"balance"`
	FrozenDepositBalance  string                 `json:"frozen_deposit_balance"`
	FrozenOrderBalance    string                 `json:"frozen_order_balance"`
	FrozenWithdrawBalance string                 `json:"frozen_withdraw_balance"`
}

type DepositData struct {
	ID         uint64                 `json:"id"`
	AccountID  uint64                 `json:"account_id"`
	Owner      string                 `json:"owner"`
	Token      *binance.ShowTokenData `json:"token"`
	L2TrxID    string                 `json:"l2_trx_id"`
	Status     string                 `json:"status"`
	Nonce      uint64                 `json:"nonce"`
	CreateTime int64                  `json:"create_time"`
	UpdateTime int64                  `json:"update_time"`
}

type WithdrawalData struct {
	ID          int64                  `json:"-"`
	AccountID   uint64                 `json:"account_id"`
	WithdrawID  string                 `json:"withdraw_id"`
	ToAccountID int64                  `json:"to_account_id"`
	ToAddress   string                 `json:"to_address"`
	Token       *binance.ShowTokenData `json:"token"`
	FeeToken    *binance.ShowTokenData `json:"fee_token"`
	TxHash      string                 `json:"tx_hash"`
	Status      string                 `json:"status"`
	Type        int                    `json:"type"`
	CreateTime  int64                  `json:"create_time"`
	UpdateTime  int64                  `json:"update_time"`
}

type TransfersDataDetail struct {
	ID          int64                  `json:"-"`
	AccountID   uint64                 `json:"account_id"`
	Address     string                 `json:"address"`
	TransferID  string                 `json:"transfer_id"`
	ToAccountID int64                  `json:"to_account_id"`
	ToAddress   string                 `json:"to_address"`
	Token       *binance.ShowTokenData `json:"token"`
	FeeToken    *binance.ShowTokenData `json:"fee_token"`
	CreateTime  int64                  `json:"create_time"`
	UpdateTime  int64                  `json:"update_time"`
}

type TradeData struct {
	ID                uint64                 `json:"id"`
	AccountID         int64                  `json:"account_id"`
	TradeId           string                 `json:"trade_id"`
	OrderId           string                 `json:"order_id"`
	PairId            uint64                 `json:"pair_id"`
	IsMaker           bool                   `json:"is_maker"`
	FilledSellToken   *binance.ShowTokenData `json:"filled_sell_token"`
	FilledBuyToken    *binance.ShowTokenData `json:"filled_buy_token"`
	FilledFeeToken    *binance.ShowTokenData `json:"filled_fee_token"`
	FilledGasFeeToken *binance.ShowTokenData `json:"filled_gas_fee_token"`
	CreateTime        int64                  `json:"create_time"`
	FillAmountBors    bool                   `json:"fill_amount_bors"`
	IsBuy             bool                   `json:"is_buy"`
	Price             string                 `json:"price"`
	SellOrderId       string                 `json:"sell_order_id"`
	BuyOrderId        string                 `json:"buy_order_id"`
	R                 int64                  `json:"R"`
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
	Status string `json:"status"`
}

type OffChainFee struct {
	UpdateAccountGasFees      []*binance.ShowTokenData `json:"update_account_gas_fees"`
	WithdrawalGasFees         []*binance.ShowTokenData `json:"withdrawal_gas_fees"`
	WithdrawalOtherGasFees    []*binance.ShowTokenData `json:"withdrawal_other_gas_fees"`
	TransferGasFees           []*binance.ShowTokenData `json:"transfer_gas_fees"`
	TransferNoIDGasFees       []*binance.ShowTokenData `json:"transfer_no_id_gas_fees"`
	OrderGasFees              []*binance.ShowTokenData `json:"order_gas_fees"`
	AddPairGasFees            []*binance.ShowTokenData `json:"add_pair_gas_fees"`
	MiningGasFees             []*binance.ShowTokenData `json:"mining_gas_fees"`
	OnChainCancelOrderGasFees []*binance.ShowTokenData `json:"on_chain_cancel_order_gas_fees"`
	UpdateAccountGas          string                   `json:"update_account_gas"`
	WithdrawalGas             string                   `json:"withdrawal_gas"`
	WithdrawalOtherGas        string                   `json:"withdrawal_other_gas"`
	TransferGas               string                   `json:"transfer_gas"`
	OrderGas                  string                   `json:"order_gas"`
	AddPairGas                string                   `json:"add_pair_gas"`
	MiningGas                 string                   `json:"mining_gas"`
	TransferNoIDGas           string                   `json:"transfer_no_id_gas"`
	UpdateAccountGasUsd       string                   `json:"update_account_gas_usd"`
	WithdrawalGasUsd          string                   `json:"withdrawal_gas_usd"`
	WithdrawalOtherGasUsd     string                   `json:"withdrawal_other_gas_usd"`
	TransferGasUsd            string                   `json:"transfer_gas_usd"`
	OrderGasUsd               string                   `json:"order_gas_usd"`
	AddPairGasUsd             string                   `json:"add_pair_gas_usd"`
	MiningGasUsd              string                   `json:"mining_gas_usd"`
	TransferNoIDUsd           string                   `json:"transfer_no_id_usd"`
	GasPrice                  string                   `json:"gas_price"`
	OrderGasMultiple          uint32                   `json:"order_gas_multiple"`
}

type GasFees struct {
	UpdateAccountGasFees       *GasFee `json:"update_account_gas_fees"`
	WithdrawalGasFees          *GasFee `json:"withdrawal_gas_fees"`
	EstimatedWithdrawalGasFees *GasFee `json:"estimated_withdrawal_gas_fees"`
	TransferGasFees            *GasFee `json:"transfer_gas_fees"`
	TransferNoIDGasFees        *GasFee `json:"transfer_no_id_gas_fees"`
	OrderGasFees               *GasFee `json:"order_gas_fees"`
	AddPairGasFees             *GasFee `json:"add_pair_gas_fees"`
	MiningGasFees              *GasFee `json:"mining_gas_fees"`
	OnChainCancelOrderGasFees  *GasFee `json:"on_chain_cancel_order_gas_fees"`
	DepositFeeConfirmGasFees   *GasFee `json:"deposit_fee_confirm_gas_fees"`
	EthPrice                   string  `json:"eth_price"`
	GasPrice                   string  `json:"gas_price"`
	PriorityFee                string  `json:"priority_fee"`
	EnclaveAddress             string  `json:"enclave_address"`
	CreateTime                 int64   `json:"create_time"`
}

func (g *GasFees) GetTokenIds() (ids map[uint64]uint64) {
	ids = map[uint64]uint64{}
	if g.UpdateAccountGasFees != nil {
		for _, token := range g.UpdateAccountGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.WithdrawalGasFees != nil {
		for _, token := range g.WithdrawalGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.EstimatedWithdrawalGasFees != nil {
		for _, token := range g.EstimatedWithdrawalGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.TransferGasFees != nil {
		for _, token := range g.TransferGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.TransferNoIDGasFees != nil {
		for _, token := range g.TransferNoIDGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.OrderGasFees != nil {
		for _, token := range g.OrderGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.AddPairGasFees != nil {
		for _, token := range g.AddPairGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.MiningGasFees != nil {
		for _, token := range g.MiningGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.OnChainCancelOrderGasFees != nil {
		for _, token := range g.OnChainCancelOrderGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	if g.DepositFeeConfirmGasFees != nil {
		for _, token := range g.DepositFeeConfirmGasFees.Tokens {
			ids[token.TokenID] = token.TokenID
		}
	}
	return
}

type GasFee struct {
	Tokens           []*NewShowTokenData `json:"tokens"`
	Signature        string              `json:"signature"`
	Gas              string              `json:"gas"`
	OrderGasMultiple uint64              `json:"order_gas_multiple"`
}

type NewShowTokenData struct {
	TokenID uint64 `json:"token_id"`
	Volume  string `json:"volume"`
	Symbol  string `json:"symbol"`
	Price   string `json:"price"`
}
