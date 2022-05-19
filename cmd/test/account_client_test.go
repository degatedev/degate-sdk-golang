package test

import (
	"github.com/degatedev/degatesdk/conf"
	"github.com/degatedev/degatesdk/degate/lib"
	"github.com/degatedev/degatesdk/degate/model"
	"github.com/degatedev/degatesdk/degate/spot"
	"github.com/degatedev/degatesdk/log"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(&conf.AppConfig{})
	response, err := client.CreateAccount(&model.AccountCreateParam{
		Address:    "",
		PrivateKey: "",
	})
	if err != nil {
		log.Print("error CreateAccount: %v", err)
		return
	}
	if response.Success() {
		log.Print("CreateAccount success: %v", lib.String(response.Data))
	} else {
		log.Print("CreateAccount fail: %v", lib.String(response))
	}
}

func TestUpdateAccount(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.UpdateAccount(&model.AccountUpdateParam{
		PrivateKey: "",
	})
	if err != nil {
		log.Print("error UpdateAccount: %v", err)
		return
	}
	if response.Success() {
		log.Print("UpdateAccount success: %v", lib.String(response.Data))
	} else {
		log.Print("UpdateAccount fail: %v", lib.String(response))
	}
}

func TestGetAccount(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetAccount()
	if err != nil {
		log.Print("error Account: %v", err)
		return
	}
	if response.Success() {
		log.Print("Account success \n %v", lib.String(response.Data))
	} else {
		log.Print("Account fail: %v", lib.String(response))
	}
}

func TestGetBalance(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	response, err := client.GetBalance(&model.AccountBalanceParam{
		Asset: "ETH",
	})
	if err != nil {
		log.Print("error AccountBalance: %v", err)
		return
	}
	if response.Success() {
		log.Print("AccountBalance success: %v", lib.String(response.Data))
	} else {
		log.Print("AccountBalance fail: %v", lib.String(response))
	}
}

func TestTransfer(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Transfer(&model.TransferParam{
		PrivateKey: "",
		Address:    "",
		Asset:      "ETH",
		Amount:     0.001,
	})
	if err != nil {
		log.Print("error Transfer: %v", err)
		return
	}
	if response.Success() {
		log.Print("Transfer success \n %v", lib.String(response.Data))
	} else {
		log.Print("Transfer fail: %v", lib.String(response))
	}
}

func TestWithdraw(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Withdraw(&model.WithdrawParam{
		Address:    "",
		Coin:       "ETH",
		Amount:     0.001,
		PrivateKey: "",
	})
	if err != nil {
		log.Print("error Withdraw: %v", err)
		return
	}
	if response.Success() {
		log.Print("Withdraw success: %v", lib.String(response.Data))
	} else {
		log.Print("Withdraw fail: %v", lib.String(response))
	}
}

func TestGetDeposits(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetDeposits(&model.DepositsParam{
		Coin:  "ETH",
		Limit: 10,
	})

	if err != nil {
		log.Error("error DepositHistory: %v", err)
		return
	}

	if response.Success() {
		log.Print("DepositHistory success\n %v", lib.String(response.Data))
	} else {
		log.Print("DepositHistory fail: %v", lib.String(response))
	}
}

func TestGetWithdraws(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetWithdraws(&model.WithdrawsParam{
		Coin: "ETH",
	})
	if err != nil {
		log.Print("error Withdraws: %v", err)
		return
	}
	if response.Success() {
		log.Print("Withdraws success\n %v", lib.String(response.Data))
	} else {
		log.Print("Withdraws fail: %v", lib.String(response))
	}
}

func TestGetTransfers(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetTransfers(&model.TransfersParam{
		Coin: "ETH",
	})
	if err != nil {
		log.Print("error Transfers: %v", err)
		return
	}
	if response.Success() {
		log.Print("Transfers success: %v", lib.String(response.Data))
	} else {
		log.Print("Transfers fail: %v", lib.String(response))
	}
}

func TestGetMyTrades(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
		Decimals: 18,
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetMyTrades(&model.AccountTradesParam{
		Symbol: "ETHUSDC",
		Limit:  20,
		Offset: 0,
	})
	if err != nil {
		log.Print("error MyTrades: %v", err)
		return
	}
	if response.Success() {
		log.Print("MyTrades success: %v", lib.String(response.Data))
	} else {
		log.Print("MyTrades fail: %v", lib.String(response))
	}
}

func TestGetOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	_, r, err := c.GetOrder(&model.OrderDetailParam{
		OrderId: "1297850191237108107827054903296",
	})
	if err != nil {
		log.Print("Order error: %v", err)
		return
	}
	if r.Success() {
		log.Print("Order success\n %v", lib.String(r.Data))
	} else {
		log.Print("Order fail: %v", lib.String(r))
	}
}

func TestGetAllOrders(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
		Decimals: 18,
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetAllOrders(&model.OrdersParam{
		Symbol: "ETHUSDC",
		Limit:  20,
	})
	if err != nil {
		log.Print("AllOrders error: %v", err)
		return
	}
	if r.Success() {
		log.Print("AllOrders success \n%v", lib.String(r.Data))
	} else {
		log.Print("AllOrders fail: %v", lib.String(r))
	}
}

func TestGetOpenOrders(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
		Decimals: 18,
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
		Decimals: 18,
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetOpenOrders(&model.OrdersParam{
		Symbol: "ETHUSDC",
		Limit:  20,
	})
	if err != nil {
		log.Info("OpenOrders error: %v", err)
		return
	}
	if r.Success() {
		log.Print("OpenOrders success \n%v", lib.String(r.Data))
	} else {
		log.Print("OpenOrders fail: %v", lib.String(r))
	}
}
