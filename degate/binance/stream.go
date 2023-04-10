package binance

type ListenKeyResponse struct {
	Response
	Data *ListenKeyData `json:"data"`
}

type ListenKeyData struct {
	ListenKey string `json:"listenKey"`
}

type EmptyResponse struct {
	Response
}
