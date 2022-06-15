package test

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
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
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success: %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
	}
}

func TestUpdateAccount(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.UpdateAccount(&model.AccountUpdateParam{
		PrivateKey: "",
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success: %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
	}
}

func TestGetAccount(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetAccount()
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success: %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
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
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success: %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
	}
	if appConfig.ShowHeader {
		t.Logf("header: %v", lib.String(response.Header))
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
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
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
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success: %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
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
		Status: 1,
		Limit: 10,
	})

	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	if response.Success() {
		t.Logf("success %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
	}
}

func TestGetWithdraws(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetWithdraws(&model.WithdrawsParam{
		Coin: "ETH",
		Status:6,
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
	}
}

func TestGetTransfers(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetTransfers(&model.TransfersParam{
		Coin: "ETH",
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success: %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
	}
}

func TestGetMyTrades(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetMyTrades(&model.AccountTradesParam{
		Symbol: "ETHUSDC",
		Limit:  20,
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if response.Success() {
		t.Logf("success: %v", lib.String(response.Data))
	} else {
		t.Errorf("fail: %v", lib.String(response))
	}
}

func TestGetOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	_, r, err := c.GetOrder(&model.OrderDetailParam{
		OrderId: "1297850191237108107827054903296",
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("success %v", lib.String(r.Data))
	} else {
		t.Errorf("fail: %v", lib.String(r))
	}
}

func TestGetHistoryOrders(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetHistoryOrders(&model.OrdersParam{
		Symbol:    "ETHUSDC",
		Limit: 20,
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Errorf("%v", lib.String(r))
	}
}

func TestGetOpenOrders(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       0,
		Symbol:   "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:       8,
		Symbol:   "USDC",
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetOpenOrders(&model.OrdersParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	if r.Success() {
		t.Logf("success %v", lib.String(r.Data))
	} else {
		t.Errorf("fail: %v", lib.String(r))
	}
}