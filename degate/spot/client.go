package spot

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"github.com/degatedev/degate-sdk-golang/internal"
)

type Client struct {
	AppConfig  *conf.AppConfig
	httpClient *internal.HttpClient
}

func (c *Client) SetAppConfig(config *conf.AppConfig) *Client {
	c.AppConfig = config
	if len(config.BaseUrl) == 0 {
		config.BaseUrl = conf.BaseUrl
	}
	if config.Timeout <= 0 {
		config.Timeout = conf.Timeout
	}

	if conf.Conf == nil {
		conf.Conf = new(conf.Config).Init()
	}
	conf.Conf.AddTokens(config.Tokens)
	return c
}

func (c *Client) GetHttpClient() *internal.HttpClient {
	if c.httpClient == nil {
		c.httpClient = internal.New(time.Duration(c.AppConfig.Timeout)*time.Second, 0.01, c.AppConfig.ShowHeader)
	}
	return c.httpClient
}

func (c *Client) CheckEddsaSign() (err error) {
	if !model.IsETHAddress(c.AppConfig.AccountAddress) {
		err = errors.New("illegal AccountAddress")
		return
	}
	if len(c.AppConfig.AssetPrivateKey) == 0 {
		err = errors.New("AssetPrivateKey is empty")
		return
	}
	err = c.CheckExchangeAddress()
	if err != nil {
		return
	}
	if !model.IsETHAddress(c.AppConfig.ExchangeAddress) {
		err = errors.New("illegal exchange address")
		return
	}
	return
}

func (c *Client) CheckChainId() (err error) {
	if c.AppConfig.ChainId == 0 {
		var response *model.ExchangeInfoResponse
		response, err = c.GetExchangeInfoInner()
		if err != nil {
			return
		}
		if response.Success() && response.Data != nil {
			c.AppConfig.ChainId = response.Data.ChainID
			if orderMaxVolume, e := decimal.NewFromString(response.Data.OrderMaxVolume); e == nil {
				conf.OrderMaxVolume = orderMaxVolume
			}
		}
	}
	return
}

func (c *Client) CheckExchangeAddress() (err error) {
	if len(c.AppConfig.ExchangeAddress) == 0 {
		var response *model.ExchangeInfoResponse
		response, err = c.GetExchangeInfoInner()
		if err != nil {
			return
		}
		if response.Success() && response.Data != nil {
			c.AppConfig.ExchangeAddress = response.Data.ExchangeAddress
			c.AppConfig.ChainId = response.Data.ChainID
			if orderMaxVolume, e := decimal.NewFromString(response.Data.OrderMaxVolume); e == nil {
				conf.OrderMaxVolume = orderMaxVolume
			}
		}
	}
	return
}

func (c *Client) GetHeaderSign() (header *request.Header, err error) {
	header = &request.Header{
		Owner:     c.AppConfig.AccountAddress,
		Time:      time.Now().Unix(),
		AccountId: c.AppConfig.AccountId,
	}
	if header.Authorization, _, err = c.GetAccessToken(); err != nil {
		return
	}
	return
}

func (c *Client) GetAccessToken() (token string, exp int64, err error) {
	if len(c.AppConfig.AccessToken) != 0 {
		exp = c.AppConfig.AccessTokenExpireTime
		if exp-int64(time.Hour)*24 > time.Now().Unix() {
			token = c.AppConfig.AccessToken
			return
		}
	}
	var signature string
	var t = time.Now()
	if signature, err = lib.SignHeader(c.AppConfig.AssetPrivateKey, c.AppConfig.AccountAddress, t.UnixMilli()); err != nil {
		return
	}
	req := &model.AccessTokenParam{
		Owner:          c.AppConfig.AccountAddress,
		EDDSASignature: signature,
		Time:           t.UnixMilli(),
		UseTradeKey:    c.AppConfig.UseTradeKey == 1,
	}
	res := &model.AccessTokenResponse{}
	err = c.Post("access/token", nil, req, res)
	if err != nil {
		return
	}

	if res.Success() && res.Data != nil {
		token = "Bearer " + res.Data.Token
		exp = res.Data.Expire
		c.AppConfig.AccessToken = token
		c.AppConfig.AccessTokenExpireTime = res.Data.Expire
	}
	return
}

func (c *Client) GetUrl(path string) string {
	url := c.AppConfig.BaseUrl
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url + conf.OrderBookPath + path
}

func (c *Client) GetUrlByAbsPath(path string) string {
	url := c.AppConfig.BaseUrl
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url + path
}

func (c *Client) GetByAbsPath(path string, header interface{}, params interface{}, response interface{}) (err error) {
	return c.GetByUrl(c.GetUrlByAbsPath(path), header, params, response)
}

func (c *Client) Get(path string, header interface{}, params interface{}, response interface{}) (err error) {
	return c.GetByUrl(c.GetUrl(path), header, params, response)
}

func (c *Client) Post(path string, header interface{}, params interface{}, response interface{}) (err error) {
	return c.PostUrl(c.GetUrl(path), header, params, response)
}

func (c *Client) PostByAbsPath(path string, header interface{}, params interface{}, response interface{}) (err error) {
	return c.PostUrl(c.GetUrlByAbsPath(path), header, params, response)
}

func (c *Client) Delete(path string, header interface{}, params interface{}, response interface{}) (err error) {
	return c.DeleteUrl(c.GetUrl(path), header, params, response)
}

func (c *Client) DeleteByAbsPath(path string, header interface{}, params interface{}, response interface{}) (err error) {
	return c.DeleteUrl(c.GetUrlByAbsPath(path), header, params, response)
}

func (c *Client) Put(path string, header interface{}, params interface{}, response interface{}) (err error) {
	return c.PutUrl(c.GetUrl(path), header, params, response)
}

func (c *Client) PutByAbsPath(path string, header interface{}, params interface{}, response interface{}) (err error) {
	return c.PutUrl(c.GetUrlByAbsPath(path), header, params, response)
}

func (c *Client) processHeader(header interface{}) (newHeader interface{}) {
	if header == nil {
		headerMap := map[string]interface{}{}
		headerMap["use-trade-key"] = c.AppConfig.UseTradeKey
		newHeader = headerMap
		return
	} else {
		b, e := json.Marshal(header)
		if e == nil {
			var headerMap map[string]interface{}
			e = json.Unmarshal(b, &headerMap)
			if e == nil {
				if _, ok := headerMap["use-trade-key"]; !ok {
					headerMap["use-trade-key"] = c.AppConfig.UseTradeKey
					newHeader = headerMap
					return
				}
			}
		}
	}
	newHeader = header
	return
}

func (c *Client) GetByUrl(url string, header interface{}, params interface{}, response interface{}) (err error) {
	header = c.processHeader(header)
	return c.GetHttpClient().GetJSON(url, header, params, response)
}

func (c *Client) PostUrl(url string, header interface{}, params interface{}, response interface{}) (err error) {
	header = c.processHeader(header)
	return c.GetHttpClient().PostJSON(url, header, params, response)
}

func (c *Client) DeleteUrl(url string, header interface{}, params interface{}, response interface{}) (err error) {
	header = c.processHeader(header)
	return c.GetHttpClient().DeleteJSON(url, header, params, response)
}

func (c *Client) PutUrl(url string, header interface{}, params interface{}, response interface{}) (err error) {
	header = c.processHeader(header)
	return c.GetHttpClient().PutJSON(url, header, params, response)
}
