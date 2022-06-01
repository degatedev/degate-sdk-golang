package spot

import (
	"errors"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"strconv"

	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/log"
)

func (c *Client) Ping() (response *binance.PingResponse, err error) {
	response = &binance.PingResponse{}
	err = c.GetByAbsPath("ping", nil, nil, response)
	if err != nil {
		return
	}
	return
}

func (c *Client) GetTrades(param *model.TradeLastedParam) (response *binance.TradeHistoryResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	p := &model.TradesParam{
		Symbol: param.Symbol,
		Limit:  param.Limit,
	}
	response, err = c.GetTradeRecords(p)
	return
}

func (c *Client) GetHistoryTrades(param *model.TradeHistoryParam) (response *binance.TradeHistoryResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	p := &model.TradesParam{
		Symbol: param.Symbol,
		Limit:  param.Limit,
		FromId: param.FromId,
	}
	response, err = c.GetTradeRecords(p)
	return
}

func (c *Client) GetTradeRecords(param *model.TradesParam) (response *binance.TradeHistoryResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	if param.Limit < 0 {
		err = errors.New("illegal limit")
		return
	}
	if param.Limit > 1000 {
		err = errors.New("illegal limit max 1000")
		return
	}
	if param.Limit == 0 {
		param.Limit = 500
	}
	r := &model.DGTradesParam{
		Start:  param.StartTime,
		End:    param.EndTime,
		Limit:  param.Limit,
		Offset: param.Offset,
		FromId: param.FromId,
	}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil {
		err = errors.New("not config base symbol")
		return
	}
	if quoteToken == nil {
		err = errors.New("not config quote symbol")
		return
	}
	r.Token1 = uint64(baseToken.Id)
	r.Token2 = uint64(quoteToken.Id)

	res := &model.TradesResponse{}
	err = c.Get("trades", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.TradeHistoryResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() && len(res.Data.Data) > 0 {
		response.Data, err = lib.ConvertTradesHistory(res.Data.Data)
	}
	return
}

func (c *Client) GetDepth(param *model.DepthParam) (response *binance.DepthResponse, err error) {
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = errors.New("error symbol")
		return
	}
	if param.Limit < 0 {
		err = errors.New("limit illegal")
		return
	}
	if param.Limit == 0 {
		param.Limit = conf.DepthLimit
	}
	r := &model.DGDepthParam{
		QuoteTokenID: uint64(quoteToken.Id),
		BaseTokenID:  uint64(baseToken.Id),
		Size:         param.Limit,
	}
	res := &model.DepthResponse{}
	err = c.GetByAbsPath(conf.WsPath+"depth", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.DepthResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertDepth(res.Data)
	}
	return
}

func (c *Client) GetDepthByDGParam(baseTokenId, quoteTokenId uint64, limit int64) (response *binance.DepthResponse, err error) {
	if limit < 0 {
		err = errors.New("illegal limit")
		return
	}
	param := &model.DGDepthParam{
		BaseTokenID:  baseTokenId,
		QuoteTokenID: quoteTokenId,
		Size:         limit,
	}
	if param.Size == 0 {
		param.Size = 100 //conf.DepthLimit
	}

	res := &model.DepthResponse{}
	err = c.GetByAbsPath("order-book-ws-api/depth", nil, param, res)
	if err != nil {
		return
	}
	log.Info("depth response: %v", res.Data)
	response = &binance.DepthResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		response.Data = lib.ConvertDepth(res.Data)
	}
	return
}

func (c *Client) GetKlines(param *model.KlineParam) (response *binance.KlineResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	if param.Limit < 0 {
		err = errors.New("illegal limit")
		return
	}
	if param.Limit > 1000 {
		err = errors.New("illegal limit max 1000")
		return
	}
	if param.Limit == 0 {
		param.Limit = 500
	}
	if len(param.Interval) == 0 {
		param.Interval = "1m"
	}
	granulary := conf.KlineInterval[param.Interval]
	if granulary == 0 {
		err = errors.New("illegal Interval")
		return
	}
	r := &model.KlinesParam{
		Start:       param.StartTime,
		End:         param.EndTime,
		Limit:       int(param.Limit),
		Granularity: granulary,
	}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil {
		err = errors.New("not config base symbol")
		return
	}
	if quoteToken == nil {
		err = errors.New("not config quote symbol")
		return
	}
	r.QuoteTokenID = uint64(quoteToken.Id)
	r.BaseTokenID = uint64(baseToken.Id)

	response = &binance.KlineResponse{}
	err = c.GetByAbsPath(conf.WsPath+"klines", nil, r, response)
	return
}

func (c *Client) GetTicker(param *model.TickerParam) (response *binance.TickerResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	r := &model.DGTickerParam{}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil {
		err = errors.New("not config base symbol")
		return
	}
	if quoteToken == nil {
		err = errors.New("not config quote symbol")
		return
	}
	r.QuoteTokenID = uint64(quoteToken.Id)
	r.BaseTokenID = uint64(baseToken.Id)

	res := &model.TickerResponse{}
	err = c.GetByAbsPath(conf.WsPath+"ticker", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.TickerResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	response.Data = lib.ConvertTicker(res.Data)
	if response.Data != nil {
		response.Data.Symbol = param.Symbol
	}
	return
}

func (c *Client) GetTickerPrice(param *model.PairPriceParam) (response *binance.PairPriceResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil {
		err = errors.New("not config base symbol")
		return
	}
	if quoteToken == nil {
		err = errors.New("not config quote symbol")
		return
	}
	pairRes, err := c.GetPair(&request.PairInfoRequest{
		Token1: uint64(baseToken.Id),
		Token2: uint64(quoteToken.Id),
	})
	if err != nil {
		return
	}
	if pairRes == nil || !pairRes.Success() || pairRes.Data == nil {
		err = errors.New("not find pair")
		return
	}

	r := &model.DGPairPriceParam{
		Pairs: strconv.Itoa(int(pairRes.Data.PairID)),
	}
	res := &model.PairPriceResponse{}
	err = c.GetByAbsPath(conf.WsPath+"pairs/prices", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.PairPriceResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if len(res.Data) > 0 {
		response.Data = lib.ConvertPairPrice(param.Symbol, res.Data[0])
	}
	return
}

func (c *Client) GetBookTicker(param *model.BookTickerParam) (response *binance.BookTickerResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	r := &model.DGTickerParam{}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil {
		err = errors.New("not config base symbol")
		return
	}
	if quoteToken == nil {
		err = errors.New("not config quote symbol")
		return
	}
	r.QuoteTokenID = uint64(quoteToken.Id)
	r.BaseTokenID = uint64(baseToken.Id)

	res := &model.BookTickerResponse{}
	err = c.GetByAbsPath(conf.WsPath+"ticker/bookTicker", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.BookTickerResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	response.Data = lib.ConvertBookTicker(param.Symbol, res.Data)
	return
}
