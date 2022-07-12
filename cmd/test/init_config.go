package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/log"
)

var appConfig = &conf.AppConfig{
	AccountId:      2475,
	AccountAddress: "0xBA2b5FEae299808b119FD410370D388B2fBF744b",
	TradingKey:     "",
}

func init() {
	log.Init(appConfig.Debug)
}
