package binance

type AccountResponse struct {
	Response
	Data *Account `json:"data"`
}

type Account struct {
	ID               uint32     `json:"id"`
	Owner            string     `json:"owner"`
	PublicKeyX       string     `json:"public_key_x"`
	PublicKeyY       string     `json:"public_key_y"`
	ReferrerId       uint64     `json:"referrer_id"`
	Nonce            int64      `json:"nonce"`
	MakerCommission  int        `json:"makerCommission"`
	TakerCommission  int        `json:"takerCommission"`
	BuyerCommission  int        `json:"buyerCommission"`
	SellerCommission int        `json:"sellerCommission"`
	CanTrade         bool       `json:"canTrade"`
	CanWithdraw      bool       `json:"canWithdraw"`
	CanDeposit       bool       `json:"canDeposit"`
	UpdateTime       int        `json:"updateTime"`
	AccountType      string     `json:"accountType"`
	Balances         []*Balance `json:"balances"`
	Permissions      []string   `json:"permissions"`
}

type BalanceResponse struct {
	Response
	Data []*Balance `json:"data"`
}

type Balance struct {
	TokenId     uint64 `json:"tokenId"`
	Asset       string `json:"asset"`
	Free        string `json:"free"`
	Locked      string `json:"locked"`
	Freeze      string `json:"freeze"`
	Withdrawing string `json:"withdrawing"`
}

type WithdrawResponse struct {
	Response
	Data *WithdrawData `json:"data"`
}

type WithdrawData struct {
	Id string `json:"id"`
}

type WithdrawHistoryResponse struct {
	Response
	Data []*WithdrawHistory `json:"data"`
}

type WithdrawHistory struct {
	Address            string `json:"address"`
	AddressTo          string `json:"addressTo"`
	Amount             string `json:"amount"`
	ApplyTime          string `json:"applyTime"`
	Coin               string `json:"coin"`
	TransactionFeeCoin string `json:"transactionFeeCoin"`
	Id                 string `json:"id"`
	Network            string `json:"network"`
	TransferType       int    `json:"transferType"`
	Status             int    `json:"status"`
	TransactionFee     string `json:"transactionFee"`
	Info               string `json:"info"`
	TxId               string `json:"txId"`
}

type TransferHistoryResponse struct {
	Response
	Data *TransferHistoryData `json:"data"`
}

type TransferHistoryData struct {
	Rows []*TransfersData `json:"rows"`
}

type TransfersData struct {
	Asset       string `json:"asset"`
	Amount      string `json:"amount"`
	Status      string `json:"status"`
	TranId      string `json:"tranId"`
	FromAccount string `json:"fromAccount"`
	ToAccount   string `json:"toAccount"`
	Timestamp   int    `json:"timestamp"`
}

type DepositHistoryResponse struct {
	Response
	Data []*DepositHistory `json:"data"`
}

type DepositHistory struct {
	Amount       string `json:"amount"`
	Coin         string `json:"coin"`
	Network      string `json:"network"`
	Status       int    `json:"status"`
	Address      string `json:"address"`
	AddressTag   string `json:"addressTag"`
	TxId         string `json:"txId"`
	InsertTime   int    `json:"insertTime"`
	TransferType int    `json:"transferType"`
	ConfirmTimes string `json:"confirmTimes"`
}

type TradeFeeResponse struct {
	Response
	Data []*TradeFee `json:"data"`
}

type TradeFee struct {
	Symbol          string `json:"symbol"`
	MakerCommission string `json:"makerCommission"`
	TakerCommission string `json:"takerCommission"`
}

type TransferResponse struct {
	Response
	Data *TransferID `json:"data"`
}

type TransferID struct {
	TranId string `json:"tranId"`
}

type GasFeeTokenResponse struct {
	Response
	Data *GasFeeToken `json:"data"`
}

type GasFeeToken struct {
	WithdrawalGasFees         []*GasFee `json:"withdrawal_gas_fees"`
	TransferGasFees           []*GasFee `json:"transfer_gas_fees"`
	TransferNoIDGasFees       []*GasFee `json:"transfer_no_id_gas_fees"`
	OnChainCancelOrderGasFees []*GasFee `json:"on_chain_cancel_order_gas_fees"`
}
