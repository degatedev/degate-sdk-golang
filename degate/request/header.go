package request

type Header struct {
	Owner         string `json:"owner"`
	Time          int64  `json:"time"`
	AccountId     uint32 `json:"account-id"`
	Authorization string `json:"authorization"`
}
