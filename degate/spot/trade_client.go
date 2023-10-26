package spot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/degatedev/degate-sdk-golang/log"

	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"github.com/degatedev/degate-sdk-golang/degate/lib"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/degatedev/degate-sdk-golang/degate/request"
	"github.com/degatedev/degate-sdk-golang/util"
	"github.com/shopspring/decimal"
)

func (c *Client) NewOrder(param *model.OrderParam) (response *binance.NewOrderResponse, err error) {
	if param.Type == model.OrderTypeLimit {
		return c.PlaceOrder(param)
	} else if param.Type == model.OrderTypeMarket {
		return c.MarketOrder(param)
	}
	err = errors.New("illegal order type")
	return
}

func (c *Client) PlaceOrder(param *model.OrderParam) (response *binance.NewOrderResponse, err error) {
	var (
		isBuy            bool
		price            decimal.Decimal
		quantity         decimal.Decimal
		quoteQuantity    decimal.Decimal
		baseVolume       decimal.Decimal
		quoteVolume      decimal.Decimal
		sellVolume       decimal.Decimal
		buyVolume        decimal.Decimal
		feeVolume        decimal.Decimal
		sellToken        = model.Token{}
		buyToken         = model.Token{}
		balanceTokens    []uint64
		sellTokenBalance decimal.Decimal
		feeTokenBalance  decimal.Decimal
		pairInfo         *model.PairInfo
		ok               bool
	)

	if err = c.CheckEddsaSign(); err != nil {
		return
	}

	if param.ValidUntil < 0 {
		err = errors.New("validUntil illegal")
		return
	}
	if param.ValidUntil == 0 {
		param.ValidUntil = time.Now().Unix() + conf.ValidUntil
	}
	if param.Quantity <= 0 && param.QuoteOrderQty <= 0 {
		err = errors.New("quantity illegal")
		return
	}
	if param.Price <= 0 {
		err = errors.New("price illegal")
		return
	}
	if len(param.NewOrderRespType) == 0 {
		param.NewOrderRespType = "FULL"
	}

	if isBuy, err = IsBuy(param.Side); err != nil {
		return
	}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = fmt.Errorf("token not config")
		return
	}
	feeToken := quoteToken
	pairInfo, err = c.CheckOrderToken(baseToken, quoteToken)
	if err != nil {
		return
	}
	if !strings.EqualFold(pairInfo.BaseToken.Symbol+pairInfo.QuoteToken.Symbol, param.Symbol) {
		err = errors.New("error symbol")
		return
	}
	price = decimal.NewFromFloat(param.Price)
	if pairInfo.IsStable {
		ok = util.IsOrderEffectiveDigits(price, pairInfo.IsStable, conf.OrderStabilityEffectiveDigitsGreaterThan1, conf.OrderStabilityEffectiveDigitsLessThan1, 0)
	} else {
		ok = util.IsOrderEffectiveDigits(price, pairInfo.IsStable, conf.OrderEffectiveDigitsGreaterThan10000, conf.OrderEffectiveDigitsLessThan10000, conf.OrderEffectiveDigitsLessThan1000)
	}
	if !ok {
		err = errors.New("price is illegal")
		return
	}

	baseDecimalDigits := int(baseToken.Decimals)
	quoteDecimalDigits := int(quoteToken.Decimals)
	if param.Quantity > 0 {
		quantity = decimal.NewFromFloat(param.Quantity)
		quantity = util.GetEffectiveVolume(quantity, conf.EffectiveDigits, baseDecimalDigits, false)
		baseVolume = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)))
		if CheckOrderMaxVolume(baseVolume) {
			err = errors.New("quantity exceed max limit")
			return
		}

		quoteQuantity = price.Mul(quantity)
		quoteQuantity = util.GetEffectiveVolume(quoteQuantity, conf.EffectiveDigits, quoteDecimalDigits, isBuy)
		quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
		if CheckOrderMaxVolume(quoteVolume) {
			err = errors.New("quote quantity exceed max limit")
			return
		}
	} else {
		inputQuoteQuantity := decimal.NewFromFloat(param.QuoteOrderQty)
		inputQuoteQuantity = util.GetEffectiveVolume(inputQuoteQuantity, conf.EffectiveDigits, quoteDecimalDigits, isBuy)
		quantity = inputQuoteQuantity.DivRound(price, 32)
		quantity = util.GetEffectiveVolume(quantity, conf.EffectiveDigits, baseDecimalDigits, false)
		baseVolume = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)))
		if CheckOrderMaxVolume(baseVolume) {
			err = errors.New("quantity exceed max limit")
			return
		}

		quoteQuantity = price.Mul(quantity)
		quoteQuantity = util.GetEffectiveVolume(quoteQuantity, conf.EffectiveDigits, quoteDecimalDigits, isBuy)
		quoteQuantity = decimal.Max(quoteQuantity, inputQuoteQuantity)
		quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
		if CheckOrderMaxVolume(quoteVolume) {
			err = errors.New("quote quantity exceed max limit")
			return
		}
	}

	r := &model.DGOrderParam{
		AccountId:        uint64(c.AppConfig.AccountId),
		ValidUntil:       param.ValidUntil,
		NewOrderRespType: param.NewOrderRespType,
		ClientOrderId:    param.NewClientOrderId,
		Price:            price.String(),
	}

	if isBuy {
		buyToken.TokenId = uint64(uint32(baseToken.Id))
		buyToken.Volume = baseVolume.String()
		sellToken.TokenId = uint64(uint32(quoteToken.Id))
		sellToken.Volume = quoteVolume.String()
		r.FillAmountBOrs = true
	} else {
		sellToken.TokenId = uint64(uint32(baseToken.Id))
		sellToken.Volume = baseVolume.String()
		buyToken.TokenId = uint64(uint32(quoteToken.Id))
		buyToken.Volume = quoteVolume.String()
		r.FillAmountBOrs = false
	}
	r.SellToken = sellToken
	r.BuyToken = buyToken
	sellVolume, err = decimal.NewFromString(r.SellToken.Volume)
	if err != nil {
		return
	}
	buyVolume, err = decimal.NewFromString(r.BuyToken.Volume)
	if err != nil {
		return
	}

	balanceTokens = append(balanceTokens, r.SellToken.TokenId)
	balanceTokens = append(balanceTokens, uint64(feeToken.Id))
	balanceResponse, err := c.GetBalanceByTokenIds(balanceTokens)
	if err != nil {
		log.Error("GetBalanceByTokenIds error:%v", err)
		return
	}
	if !balanceResponse.Success() || len(balanceResponse.Data) == 0 {
		log.Error("insufficient balance")
		return
	}
	for _, balance := range balanceResponse.Data {
		if balance.Token.TokenID == sellToken.TokenId {
			if len(balance.Balance) == 0 {
				balance.Balance = "0"
			}
			sellTokenBalance, err = decimal.NewFromString(balance.Balance)
			if err != nil {
				return
			}
		}
		if balance.Token.TokenID == uint64(feeToken.Id) {
			if len(balance.Balance) == 0 {
				balance.Balance = "0"
			}
			feeTokenBalance, err = decimal.NewFromString(balance.Balance)
			if err != nil {
				return
			}
		}
	}
	if sellTokenBalance.LessThan(sellVolume) {
		err = fmt.Errorf("insufficient balance")
		return
	}
	if !isBuy {
		feeTokenBalance = feeTokenBalance.Add(buyVolume)
	}
	feeVolume, err = c.GetOrderGasFee(quoteQuantity, feeToken, feeTokenBalance.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(feeToken.Decimals)), 32))
	if err != nil {
		return
	}

	tickerRes, err := c.Ticker24(&model.TickerParam{
		Symbol: param.Symbol,
	})
	if err != nil {
		return
	}
	if !tickerRes.Success() || tickerRes.Data == nil {
		err = errors.New("failed get taker fee")
		return
	}
	r.FeeBips = tickerRes.Data.TakerFee

	r.FeeToken = model.Token{
		TokenId: uint64(uint32(feeToken.Id)),
		Volume:  feeVolume.String(),
	}

	storageIdResponse, err := c.GetStorageID(&request.StorageIdRequest{
		Owner:     c.AppConfig.AccountAddress,
		AccountId: c.AppConfig.AccountId,
		TokenId:   uint32(sellToken.TokenId),
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
	r.OrderID = storageIdResponse.Data.ID
	r.EDDSASignature, err = lib.SignOrderRequest(c.AppConfig.AssetPrivateKey, c.AppConfig.ExchangeAddress, r, c.AppConfig.UseTradeKey)
	if err != nil {
		return
	}

	header := &model.OrderHeader{}
	if len(param.Source) == 0 {
		header.Source = conf.OrderSource
	} else {
		header.Source = param.Source
	}
	response = &binance.NewOrderResponse{}
	if strings.EqualFold(r.NewOrderRespType, "RESULT") || strings.EqualFold(r.NewOrderRespType, "FULL") {
		res := &model.NewOrderResponse{}
		err = c.Post("order", header, r, res)
		if err != nil {
			return
		}
		if err = model.Copy(response, &res.Response); err != nil {
			return
		}
		if res.Success() {
			response.Data, err = lib.ConvertOrderFull(param.Symbol, model.OrderTypeLimit, param.Side, quoteToken, baseToken, feeToken, price.String(), res.Data)
		}
	} else {
		res := &model.NewOrderAckResponse{}
		err = c.Post("order", header, r, res)
		if err != nil {
			return
		}
		if err = model.Copy(response, &res.Response); err != nil {
			return
		}
		if res.Success() {
			response.Data = lib.ConvertOrderAck(param.Symbol, r.OrderID, r.ClientOrderId)
		}
	}
	return
}

func (c *Client) MarketOrder(param *model.OrderParam) (response *binance.NewOrderResponse, err error) {
	var (
		isBuy            bool
		quantity         = decimal.NewFromInt(0)
		quoteQuantity    = decimal.NewFromInt(0)
		price            decimal.Decimal
		baseVolume       decimal.Decimal
		quoteVolume      decimal.Decimal
		sellVolume       decimal.Decimal
		buyVolume        decimal.Decimal
		feeVolume        decimal.Decimal
		sellToken        = model.Token{}
		buyToken         = model.Token{}
		balanceTokens    []uint64
		sellTokenBalance decimal.Decimal
		feeTokenBalance  decimal.Decimal
		pairInfo         *model.PairInfo
	)

	if err = c.CheckEddsaSign(); err != nil {
		return
	}

	if param.ValidUntil < 0 {
		err = errors.New("validUntil illegal")
		return
	}
	if param.ValidUntil == 0 {
		param.ValidUntil = time.Now().Unix() + conf.ValidUntil
	}
	if len(param.NewOrderRespType) == 0 {
		param.NewOrderRespType = "FULL"
	}
	if isBuy, err = IsBuy(param.Side); err != nil {
		return
	}
	if param.Quantity <= 0 && param.QuoteOrderQty <= 0 {
		err = errors.New("quantity illegal")
		return
	}
	if param.Quantity > 0 && param.QuoteOrderQty > 0 {
		err = errors.New("quantity and quoteOrderQty need one")
		return
	}
	baseToken, quoteToken := conf.Conf.GetTokens(param.Symbol)
	if baseToken == nil || quoteToken == nil {
		err = fmt.Errorf("token not config")
		return
	}
	feeToken := quoteToken
	pairInfo, err = c.CheckOrderToken(baseToken, quoteToken)
	if err != nil {
		return
	}
	if !strings.EqualFold(pairInfo.BaseToken.Symbol+pairInfo.QuoteToken.Symbol, param.Symbol) {
		err = errors.New("error symbol")
		return
	}

	baseDecimalDigits := int(baseToken.Decimals)
	quoteDecimalDigits := int(quoteToken.Decimals)

	if param.Quantity > 0 {
		quantity = decimal.NewFromFloat(param.Quantity)
		baseVolume = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)))
		if conf.OrderMaxVolume.GreaterThan(decimal.Zero) && baseVolume.GreaterThan(conf.OrderMaxVolume) {
			baseVolume = conf.OrderMaxVolume
			quantity = conf.OrderMaxVolume.DivRound(conf.Pow10.Pow(decimal.NewFromInt(int64(baseDecimalDigits))), 32)
		}
	} else {
		quoteQuantity = decimal.NewFromFloat(param.QuoteOrderQty)
		quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
		if conf.OrderMaxVolume.GreaterThan(decimal.Zero) && quoteVolume.GreaterThan(conf.OrderMaxVolume) {
			quoteVolume = conf.OrderMaxVolume
			quoteQuantity = conf.OrderMaxVolume.DivRound(conf.Pow10.Pow(decimal.NewFromInt(int64(quoteDecimalDigits))), 32)
		}
	}

	depthResponse, err := c.Depth(&model.DepthParam{
		Symbol: param.Symbol,
		Limit:  2,
	})
	if err != nil {
		return
	}
	if !depthResponse.Success() {
		err = fmt.Errorf("fail get depth code:%v message:%v", depthResponse.Code, depthResponse.Message)
		return
	}
	if depthResponse.Data == nil {
		err = errors.New("no depth")
		return
	}

	if isBuy {
		if len(depthResponse.Data.Asks) == 0 || len(depthResponse.Data.Asks[0]) == 0 {
			err = errors.New("no sell depth")
			return
		}
		price, err = decimal.NewFromString(depthResponse.Data.Asks[0][0])
		if err != nil {
			return
		}
		if pairInfo.IsStable {
			price = util.GetEffectivePriceRound(price.Mul(conf.MarketOrderBuyAdjustment), pairInfo.IsStable, conf.OrderStabilityEffectiveDigitsGreaterThan1, conf.OrderStabilityEffectiveDigitsLessThan1, 0, isBuy)
		} else {
			price = util.GetEffectivePriceRound(price.Mul(conf.MarketOrderBuyAdjustment), pairInfo.IsStable, conf.OrderEffectiveDigitsGreaterThan10000, conf.OrderEffectiveDigitsLessThan10000, conf.OrderEffectiveDigitsLessThan1000, isBuy)
		}
		if param.Quantity > 0 {
			quantity = util.GetEffectiveVolume(quantity, conf.EffectiveDigits, baseDecimalDigits, false)
			baseVolume = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)))
			quoteQuantity = util.GetEffectiveVolume(quantity.Mul(price), conf.EffectiveDigits, quoteDecimalDigits, true)
			quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
		} else {
			quoteQuantity = util.GetEffectiveVolume(quoteQuantity, conf.EffectiveDigits, quoteDecimalDigits, true)
			quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
			quantity = util.GetEffectiveVolume(quoteQuantity.DivRound(price, 32), conf.EffectiveDigits, baseDecimalDigits, false)
			baseVolume = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)))
		}
		fixPrice := quoteQuantity.DivRound(quantity, 32)
		if pairInfo.IsStable {
			price = util.GetEffectivePriceRound(fixPrice, pairInfo.IsStable, conf.OrderStabilityEffectiveDigitsGreaterThan1, conf.OrderStabilityEffectiveDigitsLessThan1, 0, isBuy)
		} else {
			price = util.GetEffectivePriceRound(fixPrice, pairInfo.IsStable, conf.OrderEffectiveDigitsGreaterThan10000, conf.OrderEffectiveDigitsLessThan10000, conf.OrderEffectiveDigitsLessThan1000, isBuy)
		}
		quoteQuantity = util.GetEffectiveVolume(quantity.Mul(price), conf.EffectiveDigits, quoteDecimalDigits, true)
		quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
	} else {
		if len(depthResponse.Data.Bids) == 0 || len(depthResponse.Data.Bids[0]) == 0 {
			err = errors.New("no  buy depth")
			return
		}
		price, err = decimal.NewFromString(depthResponse.Data.Bids[len(depthResponse.Data.Bids)-1][0])
		if err != nil {
			return
		}
		if pairInfo.IsStable {
			price = util.GetEffectivePriceRound(price.Mul(conf.MarketOrderSellAdjustment), pairInfo.IsStable, conf.OrderStabilityEffectiveDigitsGreaterThan1, conf.OrderStabilityEffectiveDigitsLessThan1, 0, isBuy)
		} else {
			price = util.GetEffectivePriceRound(price.Mul(conf.MarketOrderSellAdjustment), pairInfo.IsStable, conf.OrderEffectiveDigitsGreaterThan10000, conf.OrderEffectiveDigitsLessThan10000, conf.OrderEffectiveDigitsLessThan1000, isBuy)
		}
		if param.Quantity > 0 {
			quantity = util.GetEffectiveVolume(quantity, conf.EffectiveDigits, baseDecimalDigits, true)
			baseVolume = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)))
			quoteQuantity = util.GetEffectiveVolume(quantity.Mul(price), conf.EffectiveDigits, quoteDecimalDigits, false)
			quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
		} else {
			quoteQuantity = util.GetEffectiveVolume(quoteQuantity, conf.EffectiveDigits, quoteDecimalDigits, false)
			quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
			quantity = util.GetEffectiveVolume(quoteQuantity.DivRound(price, 32), conf.EffectiveDigits, baseDecimalDigits, true)
			baseVolume = quantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)))
		}
		fixPrice := quoteQuantity.DivRound(quantity, 32)
		if pairInfo.IsStable {
			price = util.GetEffectivePriceRound(fixPrice, pairInfo.IsStable, conf.OrderStabilityEffectiveDigitsGreaterThan1, conf.OrderStabilityEffectiveDigitsLessThan1, 0, isBuy)
		} else {
			price = util.GetEffectivePriceRound(fixPrice, pairInfo.IsStable, conf.OrderEffectiveDigitsGreaterThan10000, conf.OrderEffectiveDigitsLessThan10000, conf.OrderEffectiveDigitsLessThan1000, isBuy)
		}
		quoteQuantity = util.GetEffectiveVolume(quantity.Mul(price), conf.EffectiveDigits, quoteDecimalDigits, false)
		quoteVolume = quoteQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)))
	}
	if CheckOrderMaxVolume(baseVolume) || CheckOrderMaxVolume(quoteVolume) {
		err = errors.New("quantity exceed max limit")
		return
	}

	r := &model.DGOrderParam{
		AccountId:        uint64(c.AppConfig.AccountId),
		ValidUntil:       param.ValidUntil,
		NewOrderRespType: param.NewOrderRespType,
		ClientOrderId:    param.NewClientOrderId,
		Price:            price.String(),
	}
	if isBuy {
		buyToken.TokenId = uint64(uint32(baseToken.Id))
		buyToken.Volume = baseVolume.String()
		sellToken.TokenId = uint64(uint32(quoteToken.Id))
		sellToken.Volume = quoteVolume.String()
		r.FillAmountBOrs = param.Quantity > 0
	} else {
		sellToken.TokenId = uint64(uint32(baseToken.Id))
		sellToken.Volume = baseVolume.String()
		buyToken.TokenId = uint64(uint32(quoteToken.Id))
		buyToken.Volume = quoteVolume.String()
		r.FillAmountBOrs = param.QuoteOrderQty > 0
	}
	r.SellToken = sellToken
	r.BuyToken = buyToken
	sellVolume, err = decimal.NewFromString(r.SellToken.Volume)
	if err != nil {
		return
	}
	buyVolume, err = decimal.NewFromString(r.BuyToken.Volume)
	if err != nil {
		return
	}

	balanceTokens = append(balanceTokens, r.SellToken.TokenId)
	balanceTokens = append(balanceTokens, uint64(feeToken.Id))
	balanceResponse, err := c.GetBalanceByTokenIds(balanceTokens)
	if err != nil {
		log.Error("GetBalanceByTokenIds error:%v", err)
		return
	}
	if !balanceResponse.Success() || len(balanceResponse.Data) == 0 {
		log.Error("insufficient balance")
		return
	}
	for _, balance := range balanceResponse.Data {
		if balance.Token.TokenID == sellToken.TokenId {
			if len(balance.Balance) == 0 {
				balance.Balance = "0"
			}
			sellTokenBalance, err = decimal.NewFromString(balance.Balance)
			if err != nil {
				return
			}
		}
		if balance.Token.TokenID == uint64(feeToken.Id) {
			if len(balance.Balance) == 0 {
				balance.Balance = "0"
			}
			feeTokenBalance, err = decimal.NewFromString(balance.Balance)
			if err != nil {
				return
			}
		}
	}
	if sellTokenBalance.LessThan(sellVolume) {
		err = fmt.Errorf("insufficient balance")
		return
	}
	if !isBuy {
		feeTokenBalance = feeTokenBalance.Add(buyVolume)
	}
	feeVolume, err = c.GetOrderGasFee(quoteQuantity, feeToken, feeTokenBalance.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(feeToken.Decimals)), 32))
	if err != nil {
		return
	}

	tickerRes, err := c.Ticker24(&model.TickerParam{
		Symbol: param.Symbol,
	})
	if err != nil {
		return
	}
	r.FeeBips = tickerRes.Data.TakerFee

	r.FeeToken = model.Token{
		TokenId: uint64(uint32(feeToken.Id)),
		Volume:  feeVolume.String(),
	}

	storageIdResponse, err := c.GetStorageID(&request.StorageIdRequest{
		Owner:     c.AppConfig.AccountAddress,
		AccountId: c.AppConfig.AccountId,
		TokenId:   uint32(sellToken.TokenId),
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
	r.OrderID = storageIdResponse.Data.ID
	r.EDDSASignature, err = lib.SignOrderRequest(c.AppConfig.AssetPrivateKey, c.AppConfig.ExchangeAddress, r, c.AppConfig.UseTradeKey)
	if err != nil {
		return
	}

	header := &model.OrderHeader{}
	if len(param.Source) == 0 {
		header.Source = conf.OrderSource
	} else {
		header.Source = param.Source
	}

	response = &binance.NewOrderResponse{}
	if strings.EqualFold(r.NewOrderRespType, "RESULT") || strings.EqualFold(r.NewOrderRespType, "FULL") {
		res := &model.NewOrderResponse{}
		err = c.Post("marketOrder", header, r, res)
		if err != nil {
			return
		}
		if err = model.Copy(response, &res.Response); err != nil {
			return
		}
		if res.Success() {
			response.Data, err = lib.ConvertOrderFull(param.Symbol, model.OrderTypeMarket, param.Side, quoteToken, baseToken, feeToken, price.String(), res.Data)
		}
	} else {
		res := &model.NewOrderAckResponse{}
		err = c.Post("marketOrder", header, r, res)
		if err != nil {
			return
		}
		if err = model.Copy(response, &res.Response); err != nil {
			return
		}
		if res.Success() {
			response.Data = lib.ConvertOrderAck(param.Symbol, r.OrderID, r.ClientOrderId)
		}
	}
	return
}

func (c *Client) CancelOpenOrders(includeGrid bool) (response *binance.Response, err error) {
	if err = c.CheckEddsaSign(); err != nil {
		return
	}

	r := &model.DGCancelAllParam{
		AccountId:   uint64(c.AppConfig.AccountId),
		IncludeGrid: includeGrid,
		Timestamp:   time.Now().Unix(),
	}

	ig := uint64(1)
	if !includeGrid {
		ig = 0
	}
	r.EDDSASignature, err = lib.SignCancelOrderNew(c.AppConfig.AssetPrivateKey, c.AppConfig.ExchangeAddress,
		uint64(c.AppConfig.AccountId), ig, "0", uint64(r.Timestamp), c.AppConfig.UseTradeKey)
	if err != nil {
		return
	}
	res := &model.CancelAllOrderResponse{}
	err = c.Delete("cancelAllOrders", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.Response{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if !res.Success() {
		err = fmt.Errorf("cancelAllOrder failed: code=%d status=%d msg=%v", res.Code, res.HttpStatusCode, res.Message)
	}
	return
}

func (c *Client) CancelOrder(param *model.CancelOrderParam) (response *binance.OrderCancelResponse, err error) {
	if len(param.OrderId) == 0 && len(param.OrigClientOrderId) == 0 {
		err = errors.New("both orderId and OrigClientOrderId empty")
		return
	}
	if err = c.CheckEddsaSign(); err != nil {
		return
	}
	detail, _, err := c.GetOrder(&model.OrderDetailParam{
		OrderId:           param.OrderId,
		OrigClientOrderId: param.OrigClientOrderId,
	})
	if err != nil {
		return
	}
	if detail == nil {
		err = errors.New("fail get order detail nil")
		return
	}
	if !detail.Success() {
		err = fmt.Errorf("fail get order detail code:%v message:%v", detail.Code, detail.Message)
		return
	}
	if detail.Data == nil {
		err = errors.New("fail get order detail nil")
		return
	}

	r := &model.DGCancelOrderParam{
		AccountId:     uint64(c.AppConfig.AccountId),
		OrderId:       param.OrderId,
		ClientOrderId: param.OrigClientOrderId,
		FeeToken: model.Token{
			TokenId: 0,
			Volume:  "0",
		},
	}

	r.EDDSASignature, err = lib.SignCancelOrderNew(c.AppConfig.AssetPrivateKey, c.AppConfig.ExchangeAddress, uint64(c.AppConfig.AccountId), uint64(detail.Data.StorageID), r.FeeToken.Volume, r.FeeToken.TokenId, c.AppConfig.UseTradeKey)
	if err != nil {
		return
	}
	res := &model.CancelOrderResponse{}
	err = c.Delete("order", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.OrderCancelResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		_, o, e := c.GetOrder(&model.OrderDetailParam{
			OrderId:           param.OrderId,
			OrigClientOrderId: param.OrigClientOrderId,
		})
		if e != nil {
			return
		}
		if o != nil && o.Success() && o.Data != nil {
			response.Data = o.Data
		}
	}
	return
}

func (c *Client) CancelOrderOnChain(param *model.CancelOrderParam) (response *binance.OrderCancelResponse, err error) {
	var (
		gasFeeSymbol string
		gasFee       *binance.GasFee
		feeTokenId   uint32
		feeVolume    = "0"
	)

	if err = c.CheckEddsaSign(); err != nil {
		return
	}

	if len(param.OrderId) == 0 && len(param.OrigClientOrderId) == 0 {
		err = errors.New("both orderId and OrigClientOrderId empty")
		return
	}

	if len(param.Fee) == 0 {
		gasFeeSymbol = conf.GasFeeSymbol
	} else {
		gasFeeSymbol = param.Fee
	}

	detail, _, err := c.GetOrder(&model.OrderDetailParam{
		OrderId:           param.OrderId,
		OrigClientOrderId: param.OrigClientOrderId,
	})
	if err != nil {
		return
	}
	if detail == nil {
		err = errors.New("fail get order detail")
		return
	}
	if !detail.Success() {
		err = fmt.Errorf("fail get order detail code:%v message:%v", detail.Code, detail.Message)
		return
	}
	if detail.Data == nil {
		err = errors.New("fail get order detail")
		return
	}

	var gasResponse *binance.GasFeeResponse
	gasResponse, err = c.GetGasFee()
	if err != nil {
		return
	}
	if gasResponse == nil || !gasResponse.Success() || gasResponse.Data == nil {
		err = errors.New("error get gasFee")
		return
	}

	if len(gasResponse.Data.OnChainCancelOrderGasFees) > 0 {
		for _, gas := range gasResponse.Data.OnChainCancelOrderGasFees {
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

	r := &model.DGCancelOrderParam{
		AccountId:     uint64(c.AppConfig.AccountId),
		OrderId:       param.OrderId,
		ClientOrderId: param.OrigClientOrderId,
		FeeToken: model.Token{
			TokenId: uint64(feeTokenId),
			Volume:  feeVolume,
		},
	}

	r.EDDSASignature, err = lib.SignCancelOrderNew(c.AppConfig.AssetPrivateKey, c.AppConfig.ExchangeAddress, uint64(c.AppConfig.AccountId), uint64(detail.Data.StorageID), r.FeeToken.Volume, r.FeeToken.TokenId, c.AppConfig.UseTradeKey)
	if err != nil {
		return
	}

	res := &model.CancelOrderResponse{}
	err = c.Delete("orderCancelOnChain", nil, r, res)
	if err != nil {
		return
	}
	response = &binance.OrderCancelResponse{}
	if err = model.Copy(response, &res.Response); err != nil {
		return
	}
	if res.Success() {
		_, o, e := c.GetOrder(&model.OrderDetailParam{
			OrderId:           param.OrderId,
			OrigClientOrderId: param.OrigClientOrderId,
		})
		if e != nil {
			return
		}
		if o != nil && o.Success() && o.Data != nil {
			response.Data = o.Data
		}
	}
	return
}

func (c *Client) CheckOrderToken(baseToken *model.TokenInfo, quoteToken *model.TokenInfo) (pairInfo *model.PairInfo, err error) {
	pairResponse, err := c.GetPair(&request.PairInfoRequest{
		Token1: uint64(baseToken.Id),
		Token2: uint64(quoteToken.Id),
	})
	if err != nil {
		return
	}
	if pairResponse == nil || !pairResponse.Success() || pairResponse.Data == nil || pairResponse.Data.BaseToken == nil || pairResponse.Data.QuoteToken == nil {
		err = fmt.Errorf("GetPairInfo fail")
		return
	}
	if uint64(baseToken.Id) != pairResponse.Data.BaseToken.TokenID || uint64(quoteToken.Id) != pairResponse.Data.QuoteToken.TokenID {
		err = fmt.Errorf("error symbol")
		return
	}
	baseToken.Decimals = pairResponse.Data.BaseToken.Decimals
	quoteToken.Decimals = pairResponse.Data.QuoteToken.Decimals
	pairInfo = pairResponse.Data
	return
}

func (c *Client) GetStorageID(param *request.StorageIdRequest) (response *model.StorageIdResponse, err error) {
	if !model.IsETHAddress(param.Owner) {
		err = errors.New("illegal address")
		return
	}
	header := &request.Header{
		Owner:     param.Owner,
		Time:      time.Now().Unix(),
		AccountId: param.AccountId,
	}
	if header.Authorization, _, err = c.GetAccessToken(); err != nil {
		return
	}
	response = &model.StorageIdResponse{}
	err = c.Get("storageId", header, param, response)
	if err != nil {
		return
	}
	if response.Success() {
		response.Data.ID = lib.GenerateOrderId(param.AccountId, uint64(param.TokenId), response.Data.StorageId)
	}
	return
}

func (c *Client) GetBatchStorageID(param *request.StorageIdRequest) (response *model.BatchStorageIdResponse, err error) {
	if !model.IsETHAddress(param.Owner) {
		err = errors.New("illegal address")
		return
	}
	header := &request.Header{
		Owner:     param.Owner,
		Time:      time.Now().Unix(),
		AccountId: param.AccountId,
	}
	if header.Authorization, _, err = c.GetAccessToken(); err != nil {
		return
	}
	response = &model.BatchStorageIdResponse{}
	err = c.Get("batchStorageId", header, param, response)
	if response.Success() && len(response.Data) > 0 {
		for _, s := range response.Data {
			s.ID = lib.GenerateOrderId(param.AccountId, uint64(param.TokenId), s.StorageId)
		}
	}
	return
}

func IsBuy(side string) (is bool, err error) {
	if strings.EqualFold(side, model.OrderSideBuy) {
		is = true
	} else if strings.EqualFold(side, model.OrderSideSell) {
		is = false
	} else {
		err = errors.New("illegal side")
	}
	return
}

func (c *Client) GetOrderGasFee(quoteQuantity decimal.Decimal, feeToken *model.TokenInfo, feeBalanceQuantity decimal.Decimal) (feeVolume decimal.Decimal, err error) {
	r, err := c.GetGasFee()
	if err != nil {
		log.Error("GetOrderGasFee exchangeClient.GetGasFee() error: %v", err)
		return
	}
	if r == nil || !r.Success() || r.Data == nil || len(r.Data.OrderGasFees) == 0 {
		err = fmt.Errorf("GetOrderGasFee exchangeClient.GetGasFee() fail")
		log.Error(err.Error())
		return
	}

	var (
		feeQuantity    decimal.Decimal
		singleQuantity decimal.Decimal
		multiple       decimal.Decimal
	)
	multiple = decimal.NewFromInt32(int32(r.Data.OrderGasMultiple))
	if !multiple.IsPositive() {
		multiple = decimal.NewFromInt(1)
	}
	for _, gasFee := range r.Data.OrderGasFees {
		if gasFee.TokenId == uint64(feeToken.Id) {
			singleQuantity, err = decimal.NewFromString(gasFee.Quantity)
			if err != nil {
				return
			}
			break
		}
	}
	if singleQuantity.LessThanOrEqual(decimal.Zero) {
		err = errors.New("illegal gas fee amount")
		return
	}

	feeQuantity = decimal.Max(singleQuantity, util.CalculateQuoteFeeVolume(quoteQuantity)).Mul(multiple)
	feeDecimalDigits := int(feeToken.Decimals)
	feeQuantity = util.GetEffectiveVolume(feeQuantity, conf.GasFeeEffectiveDigits, feeDecimalDigits, true)
	feeBalanceQuantity = util.GetEffectiveVolume(feeBalanceQuantity, conf.GasFeeEffectiveDigits, feeDecimalDigits, false)
	feeQuantity = decimal.Min(feeBalanceQuantity, feeQuantity)
	feeVolume = feeQuantity.Mul(conf.Pow10.Pow(decimal.NewFromInt32(feeToken.Decimals)))
	return
}

func CheckOrderMaxVolume(volume decimal.Decimal) (isExceed bool) {
	return conf.OrderMaxVolume.GreaterThan(conf.Zero) && volume.GreaterThan(conf.OrderMaxVolume)
}
