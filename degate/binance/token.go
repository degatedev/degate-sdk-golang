package binance

type ShowTokenData struct {
	TokenID         uint64 `json:"token_id"`
	Chain           string `json:"chain"`
	Code            string `json:"code"`
	Symbol          string `json:"symbol"`
	Decimals        int32  `json:"decimals"`
	Volume          string `json:"volume"`
	ShowDecimals    int32  `json:"show_decimals"`
	IsQuotableToken bool   `json:"is_quotable_token"`
	IsGasToken      bool   `json:"is_gas_token"`
	IsListToken     bool   `json:"is_list_token"`
	Active          bool   `json:"active"`
	IsTrustedToken  bool   `json:"is_trusted_token"`
	Priority        uint64 `json:"priority"`
	MinStepSize     string `json:"min_step_size"`
	MaxSize         string `json:"max_size"`
}
