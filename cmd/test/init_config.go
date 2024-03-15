package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/log"
)

var appConfig = &conf.AppConfig{
	AccountId:       0,
	AccountAddress:  "",
	AssetPrivateKey: "",
	BaseUrl:         "",
}

func init() {
	log.Init(appConfig.Debug)
}
