package model

import (
	"github.com/degatedev/degate-sdk-golang/degate/binance"
)

type TradeFeeParam struct {
	Symbol string `json:"symbol"`
}

type TransferParam struct {
	Asset      string
	Amount     float64
	Address    string
	PrivateKey string
	Fee        string
	ValidUntil int64
}

type TransferResponse struct {
	binance.Response
	Data *TransfersDataDetail `json:"data"`
}

type AccountParam struct {
	Address string `json:"owner"`
}

type Account struct {
	ID          uint32 `json:"id"`
	Owner       string `json:"owner"`
	PublicKeyX  string `json:"public_key_x"`
	PublicKeyY  string `json:"public_key_y"`
	ReferrerId  uint64 `json:"referrer_id"`
	Nonce       int64  `json:"nonce"`
	KeyNonce    int64  `json:"key_nonce"`
	CanTrade    bool   `json:"can_trade"`
	CanWithdraw bool   `json:"can_withdraw"`
	CanDeposit  bool   `json:"can_deposit"`
}

type AccountResponse struct {
	binance.Response
	Data *Account `json:"data"`
}

type AccountCreateParam struct {
	Address    string
	PrivateKey string
	ReferrerId uint64
}

type AccountUpdateParam struct {
	PrivateKey string
	Fee        string
}

type AccountCreateData struct {
	Account
	AssetPrivateKey string
}

type AccountCreateResponse struct {
	binance.Response
	Data AccountCreateData `json:"data"`
}

type AccountUpdateData struct {
	Account
	AssetPrivateKey string
}

type AccountUpdateResponse struct {
	binance.Response
	Data AccountUpdateData `json:"data"`
}

type AccountBalanceParam struct {
	Asset string `json:"Asset"`
}

type AccountBalanceResponse struct {
	binance.Response
	Data []*Balances `json:"data"`
}

type WithdrawParam struct {
	Coin       string
	Address    string
	Amount     float64
	PrivateKey string
	ValidUntil int64
	Fee        string
}

type WithdrawResponse struct {
	binance.Response
	Data *InitiateWithdrawal `json:"data"`
}

type DepositsParam struct {
	Coin      string
	Status    int
	StartTime int64
	EndTime   int64
	Offset    int64
	Limit     int64
	TxId      string
}

type DepositsData struct {
	binance.ListData
	Data []*DepositData `json:"list"`
}

type DepositsResponse struct {
	binance.Response
	Data DepositsData `json:"data"`
}

type WithdrawsParam struct {
	Coin      string
	Status    int
	StartTime int64
	EndTime   int64
	Offset    int64
	Limit     int64
}

type WithdrawsData struct {
	binance.ListData
	Data []*WithdrawalData `json:"list"`
}

type WithdrawsResponse struct {
	binance.Response
	Data WithdrawsData `json:"data"`
}

type TransfersParam struct {
	Coin      string
	StartTime int64
	EndTime   int64
	Offset    int64
	Limit     int64
}

type TransfersData struct {
	binance.ListData
	Data []*TransfersDataDetail `json:"list"`
}

type TransfersResponse struct {
	binance.Response
	Data *TransfersData `json:"data"`
}

type AccountTradesParam struct {
	Symbol    string
	OrderId   string
	FromId    string //tradeId
	StartTime int64
	EndTime   int64
	Limit     int64
	Offset    int64
}

type TradesData struct {
	binance.ListData
	Data []*TradeData `json:"list"`
}

type TradesResponse struct {
	binance.Response
	Data TradesData `json:"data"`
}

type TradeFee struct {
	MakerCommission string                 `json:"maker_commission"`
	TakerCommission string                 `json:"taker_commission"`
	BaseToken       *binance.ShowTokenData `json:"base_token"`
	QuoteToken      *binance.ShowTokenData `json:"quote_token"`
}

type TradeFeeResponse struct {
	binance.Response
	Data []*TradeFee `json:"data"`
}
