package request

type Token struct {
	TokenId uint32 `json:"token_id" form:"token_id"`
	Volume  string `json:"volume"  form:"volume"`
}
