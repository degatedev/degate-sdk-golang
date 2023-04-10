package spot

import (
	"errors"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"strconv"
	"strings"

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

func (c *Client) Trades(param *model.TradeLastedParam) (response *binance.TradeHistoryResponse, err error) {
	p := &model.TradesParam{
		Symbol: param.Symbol,
		Limit:  param.Limit,
	}
	response, err = c.GetTradeRecords(p)
	return
}

func (c *Client) TradesHistory(param *model.TradeHistoryParam) (response *binance.TradeHistoryResponse, err error) {
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
		err = errors.New("not find symbol")
		return
	}
	r := &model.DGTradesParam{
		Limit:  param.Limit,
		FromId: param.FromId,
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
	r.Token1 = uint64(baseToken.Id)
	r.Token2 = uint64(quoteToken.Id)

	res := &model.TradesResponse{}
	err = c.Get("sdk/trades", nil, r, res)
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

func (c *Client) Depth(param *model.DepthParam) (response *binance.DepthResponse, err error) {
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = errors.New("error symbol")
		return
	}
	if param.Limit <= 0 {
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

func (c *Client) Klines(param *model.KlineParam) (response *binance.KlineResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	if param.Limit <= 0 {
		param.Limit = 500
	}
	if param.Limit > 1000 {
		param.Limit = 1000
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

func (c *Client) Ticker24(param *model.TickerParam) (response *binance.TickerResponse, err error) {
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

func (c *Client) TickerPrice(param *model.PairPriceParam) (response *binance.PairPriceResponse, err error) {
	if len(param.Symbol) == 0 {
		err = errors.New("no symbol")
		return
	}
	var (
		pairIds   string
		symbols   = strings.Split(param.Symbol, ",")
		symbolMap = map[uint64]string{}
	)
	for _, symbol := range symbols {
		var pairRes *model.PairInfoResponse
		baseToken, quoteToken := conf.Conf.GetTokens(symbol)
		if baseToken == nil {
			err = errors.New("not config base symbol")
			return
		}
		if quoteToken == nil {
			err = errors.New("not config quote symbol")
			return
		}
		pairRes, err = c.GetPair(&request.PairInfoRequest{
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
		pairIds += strconv.Itoa(int(pairRes.Data.PairID)) + ","
		symbolMap[pairRes.Data.PairID] = symbol
	}
	if len(pairIds) > 0 {
		pairIds = pairIds[0 : len(pairIds)-1]
	}
	if len(pairIds) == 0 {
		err = errors.New("not find pair")
		return
	}

	r := &model.DGPairPriceParam{
		Pairs: pairIds,
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
		response.Data = lib.ConvertPairPrice(symbolMap, res.Data)
	}
	return
}

func (c *Client) BookTicker(param *model.BookTickerParam) (response *binance.BookTickerResponse, err error) {
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
