package spot

import (
	"fmt"

	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/request"

	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
)

func (c *Client) GasFee() (response *binance.GasFeeTokenResponse, err error) {
	res,err := c.GetGasFee()
	if err != nil {
		return
	}
	response = &binance.GasFeeTokenResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertGasFeeToken(res.Data)
	}
	return
}

func (c *Client) GetGasFee() (response *binance.GasFeeResponse, err error) {
	res := &model.GasFeeResponse{}
	err = c.Get("user/offChainFee", nil, nil, res)
	if err != nil {
		return
	}
	response = &binance.GasFeeResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data, _ = lib.ConvertOffChainFee(res.Data)
	}
	return
}

func (c *Client) GetTradeFee(param *model.TradeFeeParam) (response *binance.TradeFeeResponse, err error) {
	res := &model.TradeFeeResponse{}
	token := conf.Conf.GetTokenInfo(param.Symbol)
	if token == nil {
		err = fmt.Errorf("no config token")
		return
	}
	err = c.Get("user/tradeFee", nil, &request.TradeFeeRequest{
		TokenID: token.Id,
	}, res)
	if err != nil {
		return
	}
	response = &binance.TradeFeeResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertTradeFee(res.Data)
	}
	return
}

func (c *Client) Time() (response *binance.TimeResponse, err error) {
	res := &model.TimeResponse{}
	err = c.Get("server/status", nil, nil, res)
	if err != nil {
		return
	}
	response = &binance.TimeResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertTime(res.Data)
	}
	return
}

func (c *Client) GetTokens() (response *model.TokensResponse, err error) {
	response = &model.TokensResponse{}
	err = c.Get("exchange/tokens", nil, nil, response)
	return
}

func (c *Client) GetTokenList(param *model.TokenListParam) (response *model.TokensResponse, err error) {
	response = &model.TokensResponse{}
	err = c.Get("exchange/tokenList", nil, param, response)
	return
}

func (c *Client) GetTokenData(tokenId uint64) (tokenData *model.ShowTokenData, err error) {
	response := &model.TokensResponse{}
	param := model.TokenListParam{Ids: fmt.Sprint(tokenId)}
	err = c.Get("exchange/tokenList", nil, param, response)
	if err != nil {
		return
	}
	if len(response.Data) == 0 {
		err = fmt.Errorf("not found token by id: %d", tokenId)
		return
	}
	tokenData = lib.ConvertTokenInfoToTokenData(response.Data[0])
	return
}

func (c *Client) GetExchangeInfo() (response *binance.ExchangeInfoResponse, err error) {
	res := &model.ExchangeInfoResponse{}
	err = c.Get("exchange/info", nil, nil, res)
	if err != nil {
		return
	}
	response = &binance.ExchangeInfoResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertExchangeInfo(res.Data)
	}
	return
}

func (c *Client) GetExchangeInfoInner() (response *model.ExchangeInfoResponse, err error) {
	response = &model.ExchangeInfoResponse{}
	err = c.Get("exchange/info", nil, nil, response)
	if err != nil {
		return
	}
	return
}

func (c *Client) GetPair(request *request.PairInfoRequest) (response *model.PairInfoResponse, err error) {
	response = &model.PairInfoResponse{}
	err = c.Get("pair", nil, request, response)
	if err != nil {
		return
	}
	return
}
