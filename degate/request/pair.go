package request

type PairInfoRequest struct {
	Token1     uint64 `json:"token_1"`
	Token2     uint64 `json:"token_2"`
	Token1Code string `json:"token_1_code"`
	Token2Code string `json:"token_2_code"`
}
