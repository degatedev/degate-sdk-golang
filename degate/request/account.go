package request

type AccountCreateRequest struct {
	Owner               string `json:"owner"`
	PublicKeyX          string `json:"public_key_x"`
	PublicKeyY          string `json:"public_key_y"`
	ECDSASignature      string `json:"ecdsa_signature"`
	SignatureValidUntil int64  `json:"signature_valid_until"`
	ReferrerId          uint64 `json:"referrer_id"`
	Nonce               uint64 `json:"nonce"`
	KeyNonce            uint64 `json:"key_nonce"`
}

type AccountUpdateRequest struct {
	AccountID           uint32 `json:"account_id"`
	Owner               string `json:"owner"`
	Nonce               int64  `json:"nonce"`
	KeyNonce            int64  `json:"key_nonce"`
	PublicKeyX          string `json:"public_key_x"`
	PublicKeyY          string `json:"public_key_y"`
	MaxFeeTokenId       uint32 `json:"max_fee_token_id"`
	MaxFeeVolume        string `json:"max_fee_volume"`
	ReferrerId          uint64 `json:"referrer_id"`
	ECDSASignature      string `json:"ecdsa_signature"`
	SignatureValidUntil int64  `json:"signature_valid_until"`
}

type AccountBalancesRequest struct {
	AccountId uint32 `json:"account_id"`
	Tokens    string `json:"tokens"` // "0,1"
}

type TransferRequest struct {
	TransferID     string `json:"transfer_id"`
	AccountId      uint32 `json:"account_id"`
	ToAccountId    uint32 `json:"to_account_id"`
	To             string `json:"to"`
	Token          Token  `json:"token"`
	MaxFee         Token  `json:"max_fee"`
	StorageId      uint64 `json:"storage_id"`
	ValidUntil     int64  `json:"valid_until"`
	EDDSASignature string `json:"eddsa_signature"`
	ECDSASignature string `json:"ecdsa_signature"`
}

type WithdrawRequest struct {
	WithdrawID     string `json:"withdraw_id"`
	StorageId      uint64 `json:"storage_id"`
	AccountId      uint32 `json:"account_id"`
	Token          Token  `json:"token"`
	MaxFee         Token  `json:"max_fee"`
	ValidUntil     int64  `json:"valid_until"`
	MinGas         string `json:"min_gas"`
	To             string `json:"to"`
	EDDSASignature string `json:"eddsa_signature"`
	ECDSASignature string `json:"ecdsa_signature"`
}

type TradeFeeRequest struct {
	TokenID int `json:"token_id" form:"token_id"`
}