package model

import "github.com/degatedev/degate-sdk-golang/degate/binance"

type ListenKeyResponse struct {
	binance.Response
	Data ListenKeyData `json:"data"`
}

type ListenKeyData struct {
	Expire    int64  `json:"expire"`
	ListenKey string `json:"listen_key"`
}

type ListenKeyParam struct {
	ListenKey string `json:"listen_key"`
}
