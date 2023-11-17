package test

import (
	"errors"
	"testing"

	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"github.com/degatedev/degate-sdk-golang/degate/spot"
)

func TestCreateAccount(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(&conf.AppConfig{})
	response, err := client.CreateAccount(&model.AccountCreateParam{
		Address:    "",
		PrivateKey: "",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestUpdateAccount(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.UpdateAccount(&model.AccountUpdateParam{
		PrivateKey: "",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestAccount(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.Account()
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestAccessToken(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	token, exp, err := client.GetAccessToken()
	if err != nil {
		t.Errorf("%v", err)
	} else {
		t.Logf("%v-%v", token, exp)
	}
}

func TestGetStorageId(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	storageIdResponse, err := client.GetStorageID(&request.StorageIdRequest{
		Owner:     client.AppConfig.AccountAddress,
		AccountId: client.AppConfig.AccountId,
		TokenId:   0,
		Window:    1,
	})
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if storageIdResponse != nil && !storageIdResponse.Success() {
		err = errors.New("error get storageId")
		t.Errorf("%v", err)
		return
	}
}
func TestGetBalance(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetBalance(&model.AccountBalanceParam{})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
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
		Address:    "0xedb70A95AaEEEcB9074F609aABCa7e1754870797",
		Asset:      "ETH",
		Amount:     0.0001,
		PrivateKey: "",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
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
		Address:    "0xBA2b5FEae299808b119FD410370D388B2fBF744b",
		Coin:       "ETH",
		Amount:     0.0001,
		PrivateKey: "",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestDepositHistory(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.DepositHistory(&model.DepositsParam{
		//TxId:  "0xa8fab89af2afade696d17e5f51b52a6239fe962a4f737baafac138403281130a",
	})

	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestWithdrawHistory(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.WithdrawHistory(&model.WithdrawsParam{
		Coin:   "ETH",
		Status: 6,
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestTransferHistory(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.TransferHistory(&model.TransfersParam{
		Coin: "ETH",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestMyTrades(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     45,
		Symbol: "TESB",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     3,
		Symbol: "USDT",
	})
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.MyTrades(&model.AccountTradesParam{
		Symbol: "TESBUSDT",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}

func TestGetOrder(t *testing.T) {
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	_, r, err := c.GetOrder(&model.OrderDetailParam{
		OrderId: "113644221776254868476966985531396",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestGetHistoryOrders(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     2,
		Symbol: "USDC",
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetHistoryOrders(&model.OrdersParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestGetOpenOrders(t *testing.T) {
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     0,
		Symbol: "ETH",
	})
	conf.Conf.AddToken(&model.TokenInfo{
		Id:     2,
		Symbol: "USDC",
	})
	c := new(spot.Client)
	c.SetAppConfig(appConfig)
	r, err := c.GetOpenOrders(&model.OrdersParam{
		Symbol: "ETHUSDC",
	})
	if err != nil {
		t.Errorf("%v", err)
	} else if r.Success() {
		t.Logf("%v", lib.String(r.Data))
	} else {
		t.Logf("%v", lib.String(r))
	}
}

func TestGetEstimatedWithdrawalGasFees(t *testing.T) {
	client := new(spot.Client)
	client.SetAppConfig(appConfig)
	response, err := client.GetEstimatedWithdrawalGasFee(appConfig.AccountAddress, 0)
	if err != nil {
		t.Errorf("%v", err)
	} else if response.Success() {
		t.Logf("%v", lib.String(response.Data))
	} else {
		t.Logf("%v", lib.String(response))
	}
}
