package spot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/degatedev/degate-sdk-golang/util"

	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"github.com/shopspring/decimal"
)

func (c *Client) CreateAccount(param *model.AccountCreateParam) (response *model.AccountCreateResponse, err error) {
	if len(param.PrivateKey) == 0 {
		err = errors.New("privateKey is empty")
		return
	}
	if !model.IsETHAddress(param.Address) {
		err = errors.New("illegal address")
		return
	}
	if err = c.CheckExchangeAddress(); err != nil {
		return
	}
	if err = c.CheckChainId(); err != nil {
		return
	}

	accountResponse, err := c.GetAccountInfo(&model.AccountParam{
		Address: param.Address,
	})
	if err != nil {
		return
	}
	if !accountResponse.Success() && accountResponse.Code != 101002 {
		err = fmt.Errorf("fail get account code:%v message:%v", accountResponse.Code, accountResponse.Message)
		return
	}
	if accountResponse.Data != nil {
		if len(accountResponse.Data.PublicKeyX) > 0 || len(accountResponse.Data.PublicKeyY) > 0 {
			err = fmt.Errorf("account has created")
			return
		}
		if accountResponse.Data.ID > 0 {
			var updateResponse *model.AccountUpdateResponse
			updateResponse, err = c.UpdateAccount(&model.AccountUpdateParam{
				PrivateKey: param.PrivateKey,
			})
			if err != nil {
				return
			}
			if updateResponse != nil {
				err = model.Copy(response, &updateResponse)
				return
			}
			return
		}
	}

	r := &request.AccountCreateRequest{
		Owner:               param.Address,
		Nonce:               0,
		KeyNonce:            1,
		SignatureValidUntil: time.Now().Unix() + 10*365*24*60*60,
		ReferrerId:          param.ReferrerId,
	}
	var assetPrivateKey string
	assetPrivateKey, r.PublicKeyX, r.PublicKeyY, err = lib.CreateAppKey(param.PrivateKey, c.AppConfig.BaseUrl, uint(r.KeyNonce))
	if err != nil {
		return
	}

	r.ECDSASignature, err = lib.SignUpdateAccount(param.PrivateKey, param.Address, "0", 0, "0", strconv.FormatInt(r.SignatureValidUntil, 10), strconv.FormatInt(int64(r.Nonce), 10),
		r.PublicKeyX, r.PublicKeyY, c.AppConfig.ExchangeAddress, r.ReferrerId, c.AppConfig.ChainId)
	if err != nil {
		return
	}
	response = &model.AccountCreateResponse{}
	err = c.Post("account", nil, r, response)
	if err == nil && response.Success() {
		response.Data.AssetPrivateKey = assetPrivateKey
	}
	return
}

func (c *Client) UpdateAccount(param *model.AccountUpdateParam) (response *model.AccountUpdateResponse, err error) {
	var (
		gasFeeSymbol    string
		gasFee          *binance.GasFee
		feeTokenId      uint32
		feeVolume       string
		feeVolumeDec    decimal.Decimal
		assetPrivateKey string
		accountId       uint32
		nonce           int64
		keyNonce        int64
	)

	if err = c.CheckExchangeAddress(); err != nil {
		return
	}
	if err = c.CheckChainId(); err != nil {
		return
	}
	if len(param.PrivateKey) == 0 {
		err = errors.New("privateKey is empty")
		return
	}
	if !model.IsETHAddress(c.AppConfig.AccountAddress) {
		err = errors.New("illegal accountAddress")
		return
	}

	accountResponse, err := c.GetAccountInfo(&model.AccountParam{
		Address: c.AppConfig.AccountAddress,
	})
	if err != nil {
		return
	}
	if !accountResponse.Success() {
		err = fmt.Errorf("fail get account code:%v message:%v", accountResponse.Code, accountResponse.Message)
		return
	}
	if accountResponse.Data == nil {
		err = fmt.Errorf("fail get account")
		return
	}
	if accountResponse.Data.ID <= 0 {
		err = fmt.Errorf("error account not create")
		return
	}

	if len(param.Fee) == 0 {
		gasFeeSymbol = conf.GasFeeSymbol
	} else {
		gasFeeSymbol = param.Fee
	}

	if len(accountResponse.Data.PublicKeyX) == 0 || len(accountResponse.Data.PublicKeyY) == 0 {
		accountId = 0
		nonce = 0
		keyNonce = 1
		gasFeeSymbol = conf.GasFeeSymbol
		var tokenRes *model.TokensResponse
		tokenRes, err = c.TokenList(&model.TokenListParam{
			Symbols: gasFeeSymbol,
		})
		if err != nil {
			return
		}
		if len(tokenRes.Data) != 1 {
			err = fmt.Errorf("not find ETH token")
			return
		}
		feeTokenId = uint32(tokenRes.Data[0].Id)
		feeVolume = "0"
	} else {
		accountId = accountResponse.Data.ID
		nonce = accountResponse.Data.Nonce + 1
		keyNonce = accountResponse.Data.KeyNonce + 1
		//获取 gas fee
		var gasResponse *binance.GasFeeResponse
		exchangeClient := new(Client)
		exchangeClient.SetAppConfig(c.AppConfig)
		gasResponse, err = exchangeClient.GetGasFee()
		if err != nil {
			return
		}
		if gasResponse == nil || !gasResponse.Success() || gasResponse.Data == nil {
			err = errors.New("error get gasFee")
			return
		}

		if len(gasResponse.Data.UpdateAccountGasFees) > 0 {
			for _, gas := range gasResponse.Data.UpdateAccountGasFees {
				if strings.EqualFold(gas.Symbol, gasFeeSymbol) {
					gasFee = gas
					break
				}
			}
		}

		if gasFee == nil {
			err = errors.New("error not find gas fee")
			return
		}
		feeTokenId = uint32(gasFee.TokenId)
		feeVolume = gasFee.Volume
		feeVolumeDec, err = decimal.NewFromString(feeVolume)
		if err != nil {
			return
		}
		if feeVolumeDec.LessThanOrEqual(decimal.Zero) {
			err = fmt.Errorf("error gas fee volume must greater than 0")
			return
		}
	}

	r := &request.AccountUpdateRequest{
		AccountID:           accountId,
		Owner:               c.AppConfig.AccountAddress,
		Nonce:               nonce,
		KeyNonce:            keyNonce,
		SignatureValidUntil: time.Now().Unix() + 10*365*24*60*60,
		ReferrerId:          0,
		MaxFeeTokenId:       feeTokenId,
		MaxFeeVolume:        feeVolume,
	}
	assetPrivateKey, r.PublicKeyX, r.PublicKeyY, err = lib.CreateAppKey(param.PrivateKey, c.AppConfig.BaseUrl, uint(r.KeyNonce))
	if err != nil {
		return
	}
	r.ECDSASignature, err = lib.SignUpdateAccount(param.PrivateKey, c.AppConfig.AccountAddress, strconv.FormatUint(uint64(r.AccountID), 10), uint64(r.MaxFeeTokenId), r.MaxFeeVolume, strconv.FormatInt(r.SignatureValidUntil, 10), strconv.FormatInt(r.Nonce, 10),
		r.PublicKeyX, r.PublicKeyY, c.AppConfig.ExchangeAddress, r.ReferrerId, c.AppConfig.ChainId)
	if err != nil {
		return
	}
	response = &model.AccountUpdateResponse{}
	err = c.Put("account", nil, r, response)
	if err == nil && response.Success() {
		response.Data.AssetPrivateKey = assetPrivateKey
	}
	return
}

func (c *Client) GetAccountInfo(param *model.AccountParam) (response *model.AccountResponse, err error) {
	if !model.IsETHAddress(param.Address) {
		err = errors.New("illegal address")
		return
	}
	response = &model.AccountResponse{}
	err = c.Get("account", nil, param, response)
	return
}

func (c *Client) Account() (response *binance.AccountResponse, err error) {
	account, err := c.GetAccountInfo(&model.AccountParam{
		Address: c.AppConfig.AccountAddress,
	})
	response = &binance.AccountResponse{}
	if err != nil {
		return
	}
	if err = model.Copy(response, &account.Response); err != nil {
		return
	}
	if !response.Success() {
		return
	}
	if account.Data == nil {
		return
	}
	c.AppConfig.AccountId = account.Data.ID
	var balance *binance.BalanceResponse
	if account.Data.ID > 0 {
		balance, err = c.GetBalance(&model.AccountBalanceParam{
			Asset: "",
		})
		if err != nil {
			return
		}
		if err = model.Copy(response, &balance.Response); err != nil {
			return
		}
		if !response.Success() {
			return
		}
	}
	response.Data, err = lib.ConvertAccount(account.Data)
	if err != nil {
		return
	}
	if balance != nil {
		response.Data.Balances = balance.Data
	}
	return
}

func (c *Client) GetBalance(param *model.AccountBalanceParam) (response *binance.BalanceResponse, err error) {
	if !model.IsETHAddress(c.AppConfig.AccountAddress) {
		err = errors.New("illegal accountAddress")
		return
	}
	tokenIds, err := conf.Conf.GetTokenIds(param.Asset)
	if err != nil {
		return
	}
	r := &request.AccountBalancesRequest{
		AccountId: c.AppConfig.AccountId,
		Tokens:    tokenIds,
	}
	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	res := &model.AccountBalanceResponse{}
	err = c.Get("user/balances", header, r, res)
	if err != nil {
		return
	}
	response = &binance.BalanceResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() && len(res.Data) > 0 {
		response.Data, err = lib.ConvertBalances(res.Data)
	}
	return
}

func (c *Client) GetBalanceByTokenIds(tokenIds []uint64) (response *model.AccountBalanceResponse, err error) {
	if !model.IsETHAddress(c.AppConfig.AccountAddress) {
		err = errors.New("illegal accountAddress")
		return
	}
	if len(tokenIds) == 0 {
		err = errors.New("no tokenId")
		return
	}

	tokens := ""
	for _, id := range tokenIds {
		tokens += fmt.Sprint(id) + ","
	}
	r := &request.AccountBalancesRequest{
		AccountId: c.AppConfig.AccountId,
		Tokens:    tokens[0 : len(tokens)-1],
	}
	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	res := &model.AccountBalanceResponse{}
	err = c.Get("user/balances", header, r, res)
	if err != nil {
		return
	}

	return res, nil
}

func (c *Client) Transfer(param *model.TransferParam) (response *binance.TransferResponse, err error) {
	var (
		gasFeeSymbol string
		gasFee       *binance.GasFee
		gasFees      []*binance.GasFee
		tokenData    *binance.ShowTokenData
		quantity     decimal.Decimal
		volume       string
		feeTokenId   uint32
		feeVolume    = "0"
		toAccountId  uint32
	)

	if err = c.CheckEddsaSign(); err != nil {
		return
	}
	if err = c.CheckChainId(); err != nil {
		return
	}

	if len(param.PrivateKey) == 0 {
		err = errors.New("privateKey is empty")
		return
	}

	if !model.IsETHAddress(param.Address) {
		err = errors.New("illegal address")
		return
	}

	if param.ValidUntil < 0 {
		err = errors.New("illegal ValidUntil")
		return
	}
	if param.ValidUntil == 0 {
		param.ValidUntil = time.Now().Unix() + conf.ValidUntil
	}

	if len(param.Fee) == 0 {
		gasFeeSymbol = conf.GasFeeSymbol
	} else {
		gasFeeSymbol = param.Fee
	}

	token := conf.Conf.GetTokenInfo(param.Asset)
	if token == nil {
		err = errors.New("no config asset token")
		return
	}
	tokenData, err = c.GetTokenData(uint64(token.Id))
	if err != nil {
		return
	}
	if tokenData == nil {
		err = errors.New("asset token not find")
		return
	}
	if quantity = decimal.NewFromFloat(param.Amount); !quantity.IsPositive() {
		err = errors.New("illegal amount")
		return
	}
	if !util.IsEffectiveDigits(quantity.String(), conf.TransferEffectiveDigits) {
		err = fmt.Errorf("amount max digits %v", conf.TransferEffectiveDigits)
		return
	}
	quantity = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(tokenData.Decimals)))
	if !quantity.IsInteger() {
		err = fmt.Errorf("error amount max decimal %d", tokenData.Decimals)
		return
	}
	volume = quantity.String()

	gasResponse, err := c.GetGasFee()
	if err != nil {
		return
	}
	if gasResponse == nil || !gasResponse.Success() || gasResponse.Data == nil {
		err = errors.New("error get gasFee")
		return
	}

	accountResponse, err := c.GetAccountInfo(&model.AccountParam{
		Address: param.Address,
	})
	if err != nil {
		return
	}

	if accountResponse != nil && accountResponse.Data != nil && accountResponse.Data.ID > 0 {
		gasFees = gasResponse.Data.TransferGasFees
		toAccountId = accountResponse.Data.ID
	} else {
		gasFees = gasResponse.Data.TransferNoIDGasFees
	}

	if len(gasFees) > 0 {
		for _, gas := range gasFees {
			if strings.EqualFold(gas.Symbol, gasFeeSymbol) {
				gasFee = gas
				break
			}
		}
		if gasFee == nil {
			err = errors.New("error fee")
			return
		}
		feeTokenId = uint32(gasFee.TokenId)
		feeVolume = gasFee.Volume
	}

	feeVolumeDec, err := decimal.NewFromString(feeVolume)
	if err != nil || feeVolumeDec.LessThanOrEqual(decimal.Zero) {
		err = errors.New("illegal gas fee amount")
		return
	}

	storageIdResponse, err := c.GetStorageID(&request.StorageIdRequest{
		Owner:     c.AppConfig.AccountAddress,
		AccountId: c.AppConfig.AccountId,
		TokenId:   uint32(token.Id),
		Window:    1,
	})
	if err != nil {
		return
	}
	if storageIdResponse != nil && !storageIdResponse.Success() {
		err = errors.New("error get storageId")
		return
	}

	r := &request.TransferRequest{
		TransferID:  storageIdResponse.Data.ID,
		StorageId:   storageIdResponse.Data.StorageId,
		AccountId:   c.AppConfig.AccountId,
		ToAccountId: toAccountId,
		To:          strings.ToLower(param.Address),
		Token: request.Token{
			TokenId: uint32(token.Id),
			Volume:  volume,
		},
		MaxFee: request.Token{
			TokenId: feeTokenId,
			Volume:  feeVolume,
		},
		ValidUntil: param.ValidUntil,
	}
	if r.ECDSASignature, err = lib.SignTransferEcdsa(
		param.PrivateKey,
		c.AppConfig.AccountAddress,
		strconv.FormatUint(uint64(c.AppConfig.AccountId), 10),
		strconv.FormatUint(uint64(r.Token.TokenId), 10),
		r.Token.Volume,
		strconv.FormatUint(uint64(r.MaxFee.TokenId), 10),
		r.MaxFee.Volume,
		strconv.FormatInt(r.ValidUntil, 10),
		strconv.FormatUint(r.StorageId, 10),
		c.AppConfig.ExchangeAddress,
		r.To,
		c.AppConfig.ChainId); err != nil {
		return
	}
	if r.EDDSASignature, err = lib.SignTransfer(c.AppConfig.AssetPrivateKey, c.AppConfig.ExchangeAddress, r); err != nil {
		return
	}
	res := &model.TransferResponse{}
	err = c.Post("user/transfers", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.TransferResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Data != nil {
		response.Data = &binance.TransferID{
			TranId: res.Data.TransferID,
		}
	}
	return
}

func (c *Client) Withdraw(param *model.WithdrawParam) (response *binance.WithdrawResponse, err error) {
	var (
		gasFeeSymbol string
		gasFee       *binance.GasFee
		gasFees      []*binance.GasFee
		tokenData    *binance.ShowTokenData
		quantity     decimal.Decimal
		volume       string
		feeTokenId   uint32
		feeVolume    string
		feeVolumeDec decimal.Decimal
	)

	if err = c.CheckEddsaSign(); err != nil {
		return
	}
	if err = c.CheckChainId(); err != nil {
		return
	}

	if len(param.PrivateKey) == 0 {
		err = errors.New("privateKey is empty")
		return
	}

	if !model.IsETHAddress(param.Address) {
		err = errors.New("illegal address")
		return
	}

	if param.ValidUntil < 0 {
		err = errors.New("illegal ValidUntil")
		return
	}
	if param.ValidUntil == 0 {
		param.ValidUntil = time.Now().Unix() + conf.ValidUntil
	}

	if len(param.Fee) == 0 {
		gasFeeSymbol = conf.GasFeeSymbol
	} else {
		gasFeeSymbol = param.Fee
	}

	token := conf.Conf.GetTokenInfo(param.Coin)
	if token == nil {
		err = errors.New("no config coin token")
		return
	}
	tokenData, err = c.GetTokenData(uint64(token.Id))
	if err != nil {
		return
	}
	if tokenData == nil {
		err = errors.New("coin token not find")
		return
	}
	if quantity = decimal.NewFromFloat(param.Amount); !quantity.IsPositive() {
		err = errors.New("illegal amount")
		return
	}
	quantity = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(tokenData.Decimals)))
	if !quantity.IsInteger() {
		err = fmt.Errorf("error amount max decimal %d", tokenData.Decimals)
		return
	}
	volume = quantity.String()

	//获取 gas fee
	gasResponse, err := c.GetEstimatedWithdrawalGasFee(param.Address, tokenData.TokenID)
	if err != nil {
		return
	}
	if gasResponse == nil || !gasResponse.Success() || gasResponse.Data == nil {
		err = errors.New("error get gasFee")
		return
	}

	gasFees = gasResponse.Data.EstimatedWithdrawalGasFees

	if len(gasFees) == 0 {
		err = fmt.Errorf("not find gas fee token")
		return
	}
	for _, gas := range gasFees {
		if strings.EqualFold(gas.Symbol, gasFeeSymbol) {
			gasFee = gas
			break
		}
	}
	if gasFee == nil {
		err = errors.New("not find gas fee token")
		return
	}

	feeTokenId = uint32(gasFee.TokenId)
	feeVolume = gasFee.Volume
	if len(feeVolume) == 0 {
		err = errors.New("error gas fee is empty")
		return
	}
	feeVolumeDec, err = decimal.NewFromString(feeVolume)
	if err != nil {
		return
	}
	if feeVolumeDec.LessThanOrEqual(decimal.Zero) {
		err = errors.New("illegal gas fee volume")
		return
	}

	r := &request.WithdrawRequest{
		AccountId: c.AppConfig.AccountId,
		Token: request.Token{
			TokenId: uint32(token.Id),
			Volume:  volume,
		},
		MaxFee: request.Token{
			TokenId: feeTokenId,
			Volume:  feeVolume,
		},
		ValidUntil: param.ValidUntil,
		MinGas:     conf.MinGas,
		To:         param.Address,
	}

	storageIdResponse, err := c.GetStorageID(&request.StorageIdRequest{
		Owner:     c.AppConfig.AccountAddress,
		AccountId: c.AppConfig.AccountId,
		TokenId:   uint32(token.Id),
		Window:    1,
	})
	if err != nil {
		return
	}
	if storageIdResponse != nil && !storageIdResponse.Success() {
		err = errors.New("error get storageId")
		return
	}
	r.StorageId = storageIdResponse.Data.StorageId
	r.WithdrawID = storageIdResponse.Data.ID

	if r.EDDSASignature, err = lib.SignWithdrawEddsa(c.AppConfig.AssetPrivateKey, c.AppConfig.ExchangeAddress, r); err != nil {
		return
	}
	//ecdsa签名
	if r.ECDSASignature, err = lib.SignWithdrawEcdsa(param.PrivateKey,
		c.AppConfig.AccountAddress,
		strconv.FormatUint(uint64(r.AccountId), 10),
		strconv.FormatUint(uint64(r.Token.TokenId), 10),
		r.Token.Volume,
		strconv.FormatUint(uint64(r.MaxFee.TokenId), 10),
		r.MaxFee.Volume,
		r.MinGas,
		strconv.FormatUint(uint64(r.ValidUntil), 10),
		strconv.FormatUint(r.StorageId, 10),
		c.AppConfig.ExchangeAddress,
		r.To,
		c.AppConfig.ChainId); err != nil {
		return
	}
	res := &model.WithdrawResponse{}
	err = c.Post("user/withdrawals", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.WithdrawResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = &binance.WithdrawData{
			Id: r.WithdrawID,
		}
	}
	return
}

func (c *Client) WithdrawHistory(param *model.WithdrawsParam) (response *binance.WithdrawHistoryResponse, err error) {
	tokenIds, err := conf.Conf.GetTokenIds(param.Coin)
	if err != nil {
		return
	}

	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}

	if param.Limit <= 0 || param.Limit > 1000 {
		param.Limit = 1000
	}

	r := &model.WithdrawalParam{
		AccountId: int64(c.AppConfig.AccountId),
		Tokens:    tokenIds,
		Start:     param.StartTime,
		End:       param.EndTime,
		Limit:     param.Limit,
		Offset:    param.Offset,
	}

	if param.Status == 6 {
		r.Status = "COMPLETED"
	} else if param.Status == 4 {
		r.Status = "PROCESSING"
	} else if param.Status == 5 {
		r.Status = "FAILED"
	}
	res := &model.WithdrawsResponse{}
	err = c.Get("user/withdrawals", header, r, res)
	if err != nil {
		return
	}
	response = &binance.WithdrawHistoryResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() && len(res.Data.Data) > 0 {
		response.Data, err = lib.ConvertWithdraws(res.Data.Data)
	}
	return
}

func (c *Client) DepositHistory(param *model.DepositsParam) (response *binance.DepositHistoryResponse, err error) {
	tokenIds, err := conf.Conf.GetTokenIds(param.Coin)
	if err != nil {
		return
	}

	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}

	if param.Limit <= 0 || param.Limit > 1000 {
		param.Limit = 1000
	}

	r := &model.DGDepositsParam{
		AccountId: int64(c.AppConfig.AccountId),
		Tokens:    tokenIds,
		Start:     param.StartTime,
		End:       param.EndTime,
		Limit:     param.Limit,
		Offset:    param.Offset,
		TxId:      param.TxId,
	}
	if param.Status == 6 {
		r.Status = "PROCESSING"
	} else if param.Status == 1 {
		r.Status = "COMPLETED"
	}
	res := &model.DepositsResponse{}
	err = c.Get("user/deposits", header, r, res)
	if err != nil {
		return
	}
	response = &binance.DepositHistoryResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if len(res.Data.Data) > 0 {
		response.Data, err = lib.ConvertDeposits(res.Data.Data)
	}
	return
}

func (c *Client) TransferHistory(param *model.TransfersParam) (response *binance.TransferHistoryResponse, err error) {
	tokenIds, err := conf.Conf.GetTokenIds(param.Coin)
	if err != nil {
		return
	}

	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}

	if param.Limit <= 0 || param.Limit > 1000 {
		param.Limit = 1000
	}

	r := &model.WithdrawalParam{
		AccountId: int64(c.AppConfig.AccountId),
		Tokens:    tokenIds,
		Start:     param.StartTime,
		End:       param.EndTime,
		Limit:     param.Limit,
		Offset:    param.Offset,
	}
	res := &model.TransfersResponse{}
	err = c.Get("user/transfers", header, r, res)
	if err != nil {
		return
	}
	response = &binance.TransferHistoryResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Data != nil {
		response.Data = &binance.TransferHistoryData{
			Rows: []*binance.TransfersData{},
		}
		if len(res.Data.Data) > 0 {
			response.Data.Rows, err = lib.ConvertTransfers(res.Data.Data)
		}
	}
	return
}

func (c *Client) GetMyTradesWithTokenId(param *model.AccountTradesParam, base, quote uint32) (response *binance.TradeResponse, err error) {
	r := &model.TradesUserParam{
		AccountId: int64(c.AppConfig.AccountId),
		Start:     param.StartTime,
		OrderID:   param.OrderId,
		End:       param.EndTime,
		Limit:     param.Limit,
		Offset:    param.Offset,
	}
	r.Token1 = int64(base)
	r.Token2 = int64(quote)

	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}

	res := &model.TradesResponse{}
	err = c.Get("user/trades", header, r, res)
	if err != nil {
		return
	}
	response = &binance.TradeResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() && len(res.Data.Data) > 0 {
		response.Data, err = lib.ConvertTrades(res.Data.Data)
	}
	return
}

func (c *Client) MyTrades(param *model.AccountTradesParam) (response *binance.TradeResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("not find symbol")
		return
	}
	r := &model.TradesUserParam{
		AccountId: int64(c.AppConfig.AccountId),
		OrderID:   param.OrderId,
		FromId:    param.FromId,
		Start:     param.StartTime,
		End:       param.EndTime,
		Limit:     param.Limit,
	}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil {
		err = errors.New("not config base token")
		return
	}
	if quoteToken == nil {
		err = errors.New("not config quote token")
		return
	}
	r.Token1 = int64(baseToken.Id)
	r.Token2 = int64(quoteToken.Id)

	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}

	res := &model.TradesResponse{}
	err = c.Get("sdk/user/trades", header, r, res)
	if err != nil {
		return
	}
	response = &binance.TradeResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() && len(res.Data.Data) > 0 {
		response.Data, err = lib.ConvertTrades(res.Data.Data)
	}
	return
}

func (c *Client) GetOrder(param *model.OrderDetailParam) (res *model.OrderDetailResponse, response *binance.OrderResponse, err error) {
	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	res = &model.OrderDetailResponse{}
	err = c.Get("order", header, &model.OrderDetailsParam{
		OrderID:       param.OrderId,
		ClientOrderID: param.OrigClientOrderId,
	}, res)
	if err != nil {
		return
	}
	response = &binance.OrderResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertOrder(res.Data)
	}
	return
}

func (c *Client) GetHistoryOrders(param *model.OrdersParam) (response *binance.OrdersResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("not find symbol")
		return
	}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil {
		err = errors.New("not config baseToken")
		return
	}
	if quoteToken == nil {
		err = errors.New("not config quoteToken")
		return
	}
	r := &request.OrdersRequest{
		AccountId: c.AppConfig.AccountId,
		Token1:    int32(baseToken.Id),
		Token2:    int32(quoteToken.Id),
		Start:     param.StartTime,
		End:       param.EndTime,
		Limit:     param.Limit,
		OrderId:   param.OrderId,
	}
	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	res := &model.OrdersResponse{}
	err = c.Get("sdk/historyOrders", header, r, res)
	if err != nil {
		return
	}
	response = &binance.OrdersResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertOrders(res.Data.Data)
	}
	return
}

func (c *Client) GetOpenOrders(param *model.OrdersParam) (response *binance.OrdersResponse, err error) {
	r := &request.OrdersRequest{
		AccountId: c.AppConfig.AccountId,
	}
	if len(param.Symbol) > 0 {
		baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
		if baseToken == nil {
			err = errors.New("not config baseToken")
			return
		}
		if quoteToken == nil {
			err = errors.New("not config quoteToken")
			return
		}
		r.Token1 = int32(baseToken.Id)
		r.Token2 = int32(quoteToken.Id)
	} else {
		r.Token1 = -1
		r.Token2 = -1
	}

	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	res := &model.OrdersResponse{}
	err = c.Get("sdk/openOrders", header, r, res)
	if err != nil {
		return
	}
	response = &binance.OrdersResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertOrders(res.Data.Data)
	}
	return
}

func (c *Client) GetOpenOrdersWithTokenId(param *model.OrdersParamWithTokenId) (response *model.OrdersResponse, err error) {
	r := &request.OrdersRequest{
		AccountId: c.AppConfig.AccountId,
		Status:    model.OrderStatusOpen,
		Start:     param.StartTime,
		End:       param.EndTime,
		Limit:     param.Limit,
		Token1:    param.TokenId1,
		Token2:    param.TokenId2,
	}
	return c.getOrderLists(r, "orders")
}

func (c *Client) GetMakerOpenOrdersWithTokenId(param *model.OrdersParamWithTokenId) (response *model.OrdersResponse, err error) {
	r := &request.OrdersRequest{
		AccountId: c.AppConfig.AccountId,
		Status:    model.OrderStatusOpen,
		Start:     param.StartTime,
		End:       param.EndTime,
		Limit:     param.Limit,
		Token1:    param.TokenId1,
		Token2:    param.TokenId2,
	}
	return c.getOrderLists(r, "makerOrders")
}

func (c *Client) getOrderLists(r *request.OrdersRequest, uri string) (response *model.OrdersResponse, err error) {
	if r.Limit < 0 || r.Limit > 1000 {
		err = errors.New("illegal limit 0-1000")
		return
	}
	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	res := &model.OrdersResponse{}
	err = c.Get(uri, header, r, res)
	if err != nil {
		return
	}

	return res, nil
}
