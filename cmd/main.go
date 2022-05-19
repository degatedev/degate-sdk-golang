package main

/*
#include "clib.h"
extern void CallBack(char *c, callback cb);
*/
import "C"

import (
	"encoding/json"
	"strings"
	"unsafe"

	"github.com/degatedev/degatesdk/conf"
	"github.com/degatedev/degatesdk/degate/binance"
	"github.com/degatedev/degatesdk/degate/lib"
	"github.com/degatedev/degatesdk/degate/model"
	"github.com/degatedev/degatesdk/degate/spot"
	"github.com/degatedev/degatesdk/degate/websocket"
	"github.com/degatedev/degatesdk/log"
)

//export send_request
func send_request(c *C.char, method *C.char, params *C.char) (res *C.char) {
	var (
		response      interface{}
		resposneBytes []byte
		cs            string
		p             string
		m             string
	)

	cs = C.GoString(c)
	m = C.GoString(method)
	if params != nil {
		p = C.GoString(params)
	}

	appConfig := &conf.AppConfig{}
	err := json.Unmarshal([]byte(cs), appConfig)
	if err != nil {
		panic(err)
	}
	conf.Debug = appConfig.Debug
	log.Info("config: %v\n", cs)
	log.Info("method: %v\n", m)
	log.Info("param : %v\n", p)

	if strings.EqualFold(m, "time") {
		response, err = GetTime(appConfig)
	} else if strings.EqualFold(m, "exchangeInfo") {
		response, err = GetExchangeInfo(appConfig)
	} else if strings.EqualFold(m, "gasFee") {
		response, err = GetGasFee(appConfig)
	} else if strings.EqualFold(m, "tradeFee") {
		response, err = GetTradeFee(appConfig, p)
	} else if strings.EqualFold(m, "createAccount") {
		response, err = CreateAccount(appConfig, p)
	} else if strings.EqualFold(m, "updateAccount") {
		response, err = UpdateAccount(appConfig, p)
	} else if strings.EqualFold(m, "account") {
		response, err = GetAccount(appConfig)
	} else if strings.EqualFold(m, "balance") {
		response, err = GetBalance(appConfig, p)
	} else if strings.EqualFold(m, "transfer") {
		response, err = Transfer(appConfig, p)
	} else if strings.EqualFold(m, "withdraw") {
		response, err = Withdraw(appConfig, p)
	} else if strings.EqualFold(m, "withdraws") {
		response, err = GetWithdraws(appConfig, p)
	} else if strings.EqualFold(m, "deposits") {
		response, err = GetDeposits(appConfig, p)
	} else if strings.EqualFold(m, "transfers") {
		response, err = GetTransfers(appConfig, p)
	} else if strings.EqualFold(m, "order") {
		response, err = GetOrder(appConfig, p)
	} else if strings.EqualFold(m, "orders") {
		response, err = GetAllOrders(appConfig, p)
	} else if strings.EqualFold(m, "openOrders") {
		response, err = GetOpenOrders(appConfig, p)
	} else if strings.EqualFold(m, "myTrades") {
		response, err = GetMyTrades(appConfig, p)
	} else if strings.EqualFold(m, "newOrder") {
		response, err = NewOrder(appConfig, p)
	} else if strings.EqualFold(m, "cancelOrder") {
		response, err = CancelOrder(appConfig, p)
	} else if strings.EqualFold(m, "cancelOrderOnChain") {
		response, err = CancelOrderOnChain(appConfig, p)
	} else if strings.EqualFold(m, "lastedTrades") {
		response, err = GetTrades(appConfig, p)
	} else if strings.EqualFold(m, "historyTrades") {
		response, err = GetHistoryTrades(appConfig, p)
	} else if strings.EqualFold(m, "depth") {
		response, err = GetDepth(appConfig, p)
	} else if strings.EqualFold(m, "klines") {
		response, err = GetKlines(appConfig, p)
	} else if strings.EqualFold(m, "ticker") {
		response, err = GetTicker(appConfig, p)
	} else if strings.EqualFold(m, "tickerPrice") {
		response, err = GetTickerPrice(appConfig, p)
	} else if strings.EqualFold(m, "bookTicker") {
		response, err = GetBookTicker(appConfig, p)
	} else if strings.EqualFold(m, "newListenKey") {
		response, err = NewListenKey(appConfig)
	} else if strings.EqualFold(m, "reNewListenKey") {
		response, err = ReNewListenKey(appConfig, p)
	} else if strings.EqualFold(m, "closeListenKey") {
		response, err = CloseListenKey(appConfig, p)
	} else if strings.EqualFold(m,"tokens") {
		response, err = GetTokens(appConfig, p)
	} else if strings.EqualFold(m,"cancelOpenOrders") {
		response, err = CancelAllOrders(appConfig, p)
	} else if strings.EqualFold(m, "ping") {
		response, err = Ping(appConfig)
	}

	if err != nil {
		res = C.CString(err.Error())
		return
	}
	resposneBytes, err = json.Marshal(response)
	if err != nil {
		res = C.CString(err.Error())
		return
	}
	res = C.CString(string(resposneBytes))
	return
}

//export send_subscribe
func send_subscribe(c *C.char, method *C.char, params *C.char, cb unsafe.Pointer) (res *C.char) {
	var (
		cs            string
		p             string
		m             string
		err           error
		client        *websocket.WebSocketClient
		responseBytes []byte
	)

	cs = C.GoString(c)
	m = C.GoString(method)
	if params != nil {
		p = C.GoString(params)
	}

	appConfig := &conf.AppConfig{}
	err = json.Unmarshal([]byte(cs), appConfig)
	if err != nil {
		panic(err)
	}
	conf.Debug = appConfig.Debug
	log.Info("config: %v\n", cs)
	log.Info("method: %v\n", m)
	log.Info("param : %v\n", p)

	if m == "stop" {
		err = Stop(p)
	} else if m == "kline" {
		client, err = SubscribeKline(appConfig, p, cb)
	} else if m == "ticker" {
		client, err = SubscribeTicker(appConfig, p, cb)
	} else if m == "bookTicker" {
		client, err = SubscribeBookTicker(appConfig, p, cb)
	} else if m == "trade" {
		client, err = SubscribeTrade(appConfig, p, cb)
	} else if m == "depth" {
		client, err = SubscribeDepth(appConfig, p, cb)
	} else if m == "depthUpdate" {
		client, err = SubscribeDepthUpdate(appConfig, p, cb)
	} else if m == "userData" {
		client, err = SubscribeUserData(appConfig, p, cb)
	}
	if err != nil {
		res = C.CString(err.Error())
		return
	}
	if client != nil {
		responseBytes, _ = json.Marshal(model.Client{
			Client: websocket.GetClientKey(client),
		})
		res = C.CString(string(responseBytes))
		return
	}
	return
}

func Stop(p string) (err error) {
	param := &model.StopParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	websocket.StopWebsocketClient(param.Client)
	return
}

func SubscribeKline(config *conf.AppConfig, p string, cb unsafe.Pointer) (c *websocket.WebSocketClient, err error) {
	param := &model.SubscribeKlineParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	c = &websocket.WebSocketClient{}
	c.Init(config)
	err = c.SubscribeKline(param, func(message string) {
		C.CallBack(C.CString(message), (*[0]byte)(cb))
	})
	if err != nil {
		return
	}
	return
}

func SubscribeTicker(config *conf.AppConfig, p string, cb unsafe.Pointer) (c *websocket.WebSocketClient, err error) {
	param := &model.SubscribeTickerParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	c = &websocket.WebSocketClient{}
	c.Init(config)
	err = c.SubscribeTicker(param, func(message string) {
		C.CallBack(C.CString(message), (*[0]byte)(cb))
	})
	if err != nil {
		return
	}
	return
}

func SubscribeBookTicker(config *conf.AppConfig, p string, cb unsafe.Pointer) (c *websocket.WebSocketClient, err error) {
	param := &model.SubscribeBookTickerParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	c = &websocket.WebSocketClient{}
	c.Init(config)
	err = c.SubscribeBookTicker(param, func(message string) {
		C.CallBack(C.CString(message), (*[0]byte)(cb))
	})
	if err != nil {
		return
	}
	return
}

func SubscribeTrade(config *conf.AppConfig, p string, cb unsafe.Pointer) (c *websocket.WebSocketClient, err error) {
	param := &model.SubscribeTradeParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	c = &websocket.WebSocketClient{}
	c.Init(config)
	err = c.SubscribeTrade(param, func(message string) {
		C.CallBack(C.CString(message), (*[0]byte)(cb))
	})
	if err != nil {
		return
	}
	return
}

func SubscribeDepth(config *conf.AppConfig, p string, cb unsafe.Pointer) (c *websocket.WebSocketClient, err error) {
	param := &model.SubscribeDepthParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	c = &websocket.WebSocketClient{}
	c.Init(config)
	err = c.SubscribeDepth(param, func(message string) {
		C.CallBack(C.CString(message), (*[0]byte)(cb))
	})
	if err != nil {
		return
	}
	return
}

func SubscribeDepthUpdate(config *conf.AppConfig, p string, cb unsafe.Pointer) (c *websocket.WebSocketClient, err error) {
	param := &model.SubscribeDepthUpdateParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	c = &websocket.WebSocketClient{}
	c.Init(config)
	err = c.SubscribeDepthUpdate(param, func(message string) {
		C.CallBack(C.CString(message), (*[0]byte)(cb))
	})
	if err != nil {
		return
	}
	return
}

func SubscribeUserData(config *conf.AppConfig, p string, cb unsafe.Pointer) (c *websocket.WebSocketClient, err error) {
	param := &model.SubscribeUserDataParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	c = &websocket.WebSocketClient{}
	c.Init(config)
	c.SubscribeUserData(param, func(message string) {
		C.CallBack(C.CString(message), (*[0]byte)(cb))
	})
	return
}

func GetTime(config *conf.AppConfig) (response *binance.TimeResponse, err error) {
	c := &spot.Client{}
	c.SetAppConfig(config)
	response, err = c.Time()
	return
}

func GetExchangeInfo(config *conf.AppConfig) (response *binance.ExchangeInfoResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	response, err = c.GetExchangeInfo()
	return
}

func GetGasFee(config *conf.AppConfig) (response *binance.GasFeeResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	response, err = c.GetGasFee()
	return
}

func GetTradeFee(config *conf.AppConfig, p string) (response *binance.TradeFeeResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.TradeFeeParam{}
	err = lib.ParseParam(p, param)
	if err != nil {
		return
	}
	response, err = c.GetTradeFee(param)
	return
}

func CreateAccount(config *conf.AppConfig, p string) (response *model.AccountCreateResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.AccountCreateParam{}
	err = lib.ParseParam(p, param)
	if err != nil {
		return
	}
	response, err = c.CreateAccount(param)
	return
}

func UpdateAccount(config *conf.AppConfig, p string) (response *model.AccountUpdateResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.AccountUpdateParam{}
	err = lib.ParseParam(p, param)
	if err != nil {
		return
	}
	response, err = c.UpdateAccount(param)
	return
}

func GetAccount(config *conf.AppConfig) (response *binance.AccountResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	response, err = c.GetAccount()
	return
}

func GetBalance(config *conf.AppConfig, p string) (response *binance.BalanceResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.AccountBalanceParam{}
	err = lib.ParseParam(p, param)
	if err != nil {
		return
	}
	response, err = c.GetBalance(param)
	return
}

func Transfer(config *conf.AppConfig, p string) (response *binance.TransferResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.TransferParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.Transfer(param)
	return
}

func Withdraw(config *conf.AppConfig, p string) (response *binance.WithdrawResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.WithdrawParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.Withdraw(param)
	return
}

func GetTransfers(config *conf.AppConfig, p string) (response *binance.TransferHistoryResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.TransfersParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetTransfers(param)
	return
}

func GetWithdraws(config *conf.AppConfig, p string) (response *binance.WithdrawHistoryResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.WithdrawsParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetWithdraws(param)
	return
}

func GetDeposits(config *conf.AppConfig, p string) (response *binance.DepositHistoryResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.DepositsParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetDeposits(param)
	return
}

func GetOrder(config *conf.AppConfig, p string) (response *binance.OrderResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.OrderDetailParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	_, response, err = c.GetOrder(param)
	return
}

func GetAllOrders(config *conf.AppConfig, p string) (response *binance.OrdersResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.OrdersParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetAllOrders(param)
	return
}

func GetOpenOrders(config *conf.AppConfig, p string) (response *binance.OrdersResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.OrdersParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetOpenOrders(param)
	return
}

func GetMyTrades(config *conf.AppConfig, p string) (response *binance.TradeResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.AccountTradesParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetMyTrades(param)
	return
}

func NewOrder(config *conf.AppConfig, p string) (response *binance.NewOrderResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.OrderParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.NewOrder(param)
	return
}

func CancelOrder(config *conf.AppConfig, p string) (response *binance.OrderCancelResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.CancelOrderParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.CancelOrder(param)
	return
}

func CancelOrderOnChain(config *conf.AppConfig, p string) (response *binance.OrderCancelResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.CancelOrderParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.CancelOrderOnChain(param)
	return
}

func GetTrades(config *conf.AppConfig, p string) (response *binance.TradeHistoryResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.TradeLastedParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetTrades(param)
	return
}

func GetHistoryTrades(config *conf.AppConfig, p string) (response *binance.TradeHistoryResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.TradeHistoryParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetHistoryTrades(param)
	return
}

func GetDepth(config *conf.AppConfig, p string) (response *binance.DepthResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.DepthParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetDepth(param)
	return
}

func GetKlines(config *conf.AppConfig, p string) (response *binance.KlineResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.KlineParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetKlines(param)
	return
}

func GetTicker(config *conf.AppConfig, p string) (response *binance.TickerResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.TickerParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetTicker(param)
	return
}

func GetTickerPrice(config *conf.AppConfig, p string) (response *binance.PairPriceResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.PairPriceParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetTickerPrice(param)
	return
}

func GetBookTicker(config *conf.AppConfig, p string) (response *binance.BookTickerResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.BookTickerParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetBookTicker(param)
	return
}

func NewListenKey(config *conf.AppConfig) (response *binance.ListenKeyResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	response, err = c.NewListenKey()
	return
}

func ReNewListenKey(config *conf.AppConfig, p string) (response *binance.EmptyResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.ListenKeyParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.ReNewListenKey(param)
	return
}

func CloseListenKey(config *conf.AppConfig, p string) (response *binance.EmptyResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.ListenKeyParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.DeleteListenKey(param)
	return
}

func GetTokens(config *conf.AppConfig, p string) (response *model.TokensResponse, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.TokenListParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.GetTokenList(param)
	return
}

func CancelAllOrders(config *conf.AppConfig, p string) (response *binance.Response, err error) {
	c := new(spot.Client)
	c.SetAppConfig(config)
	param := &model.CancelAllOrdersParam{}
	if err = lib.ParseParam(p, param); err != nil {
		return
	}
	response, err = c.CancelAllOrders(param.IncludeGrid)
	return
}

func Ping(config *conf.AppConfig) (response *binance.TimeResponse, err error) {
	c := &spot.Client{}
	c.SetAppConfig(config)
	response, err = c.Ping()
	return
}

func main() {}
