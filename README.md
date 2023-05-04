# DeGate Public SDK Connector Go

This is a lightweight library that works as a connector to [DeGate public SDK](https://api-docs.degate.com/spot)

## Quick Start

The SDK is compiled by Go 1.18, you can import this SDK in your Golang project:

```go
import (
    "encoding/json"
	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
	"github.com/degatedev/degate-sdk-golang/log"
)
client := &spot.Client{}
var appConfig = &conf.AppConfig{}
client.SetAppConfig(appConfig)
response,err := client.Time()
if err != nil {
    log.Error("error: %v",err)
} else if !response.Success() {
    log.Error("fail: %v",response.Code)
} else {
    log.Info("success: %v",response.Data.ServerTime)
}


client = &spot.Client{}
var appConfig = &conf.AppConfig{
    AccountId: '<accountId>',
    AccountAddress: '<accountAddress>',
	AssetPrivateKey:'<DeGate AssetPrivateKey>',
}
client.SetAppConfig(appConfig)
response,err := client.Account()
if err != nil {
    log.Error("error: %v",err)
} else if !response.Success() {
    log.Error("fail: %v",response.Code)
} else {
    account,_ := json.Marshal(response.Data)
    log.Info("success: %v",string(account))
}
```

Please find `examples` folder for more endpoints

### Testnet

The [spot testnet](https://testnet.degate.com/) is available. In order to test on testnet:

```php
var appConfig = &conf.AppConfig{
    BaseUrl:"https://testnet.degate.com/",
}
```

### Display meta info

DeGate API server returns weight usage in the header of each response. This is very useful to indentify the current usage. To reveal this value, simpily intial the client with ShowHeader=True as:

```go
import (
    "encoding/json"
    "github.com/degatedev/degate-sdk-golang/conf"
    "github.com/degatedev/degate-sdk-golang/degate/spot"
    "github.com/degatedev/degate-sdk-golang/log"
)
client := &spot.Client{}
var appConfig = &conf.AppConfig{
    ShowHeader:true,
}
client.SetAppConfig(appConfig)
response,err := client.Time()
data,_ := json.Marshal(response)
log.Info("%v",string(data))
```

the returns will be like:

```go
{"data":{"serverTime":1590579942001},"header":{"Content-Type":["application/json;charset=utf-8"],"Transfer-Encoding":["chunked"],...}}
```

## Websocket

```go
import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"time"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
	"github.com/degatedev/degate-sdk-golang/degate/websocket"
	"github.com/degatedev/degate-sdk-golang/log"
)

conf.Conf.AddToken(&model.TokenInfo{
    Id:       0,
    Symbol:   "ETH",
})
conf.Conf.AddToken(&model.TokenInfo{
    Id:       8,
    Symbol:   "USDC",
})
client := &websocket.WebSocketClient{}
client.Init(&conf.AppConfig{})
err := client.SubscribeKline(&model.SubscribeKlineParam{
    Symbol:   "ETHUSDC",
    Interval: "1m",
}, func(message string) {
    log.Print(message)
})
if err != nil {
	log.Error("error: %v",err)
    return
}
time.Sleep(time.Minute * 30)
client.Stop()
```