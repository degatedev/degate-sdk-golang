package test

import (
	"fmt"
	"github.com/BurntSushi/toml"
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
	fmt.Println("hahahhaha")
	_, err := toml.DecodeFile("./credential.toml", appConfig)
	if err != nil {
		panic(fmt.Sprintf("init appConfig failed, err=%+v", err))
	}
	fmt.Println("init appConfig success")
	fmt.Println(appConfig.BaseUrl)
	log.Init(appConfig.Debug)
}
