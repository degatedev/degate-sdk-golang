package spot

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
)

func (c *Client) NewListenKey() (response *binance.ListenKeyResponse, err error) {
	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	res := &model.ListenKeyResponse{}
	err = c.PostByAbsPath(conf.WsPath+"userDataStream", header, nil, res)
	if err != nil {
		return
	}
	response = &binance.ListenKeyResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertListenerKey(res.Data)
	}
	return
}

func (c *Client) ReNewListenKey(param *model.ListenKeyParam) (response *binance.EmptyResponse, err error) {
	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	response = &binance.EmptyResponse{}
	err = c.PutByAbsPath(conf.WsPath+"userDataStream", header, param, response)
	return
}

func (c *Client) DeleteListenKey(param *model.ListenKeyParam) (response *binance.EmptyResponse, err error) {
	header, err := c.GetHeaderSign()
	if err != nil {
		return
	}
	response = &binance.EmptyResponse{}
	err = c.DeleteByAbsPath(conf.WsPath+"userDataStream", header, param, response)
	return
}
