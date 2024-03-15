package lib

import (
	"errors"
	"strings"

	"github.com/degatedev/degate-sdk-golang/conf"
	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"github.com/degatedev/degate-sdk-golang/degate/model"
	"github.com/shopspring/decimal"
)

func ConvertOrders(orders []*model.OrderList) []*binance.Order {
	var bos []*binance.Order
	for _, o := range orders {
		if o != nil {
			bos = append(bos, ConvertOrder(o))
		}
	}
	return bos
}

func ConvertOrder(o *model.OrderList) (order *binance.Order) {
	if o == nil {
		return
	}
	var (
		pow10 = decimal.NewFromInt(10)
	)
	order = &binance.Order{}
	order.Symbol = o.GetSymbol()
	order.OrderId = o.OrderID
	order.ClientOrderId = o.ClientOrderId
	if o.OrderType == 0 {
		order.Type = "LIMIT"
		order.TimeInForce = "GTC"
	} else if o.OrderType == 1 {
		order.Type = "MARKET"
		order.TimeInForce = "IOC"
	}
	if o.IsBuy {
		order.Side = "BUY"
	} else {
		order.Side = "SELL"
	}
	order.Time = o.CreateTime
	order.UpdateTime = o.UpdateTime
	if strings.EqualFold(o.Status, "canceled") {
		order.Status = "CANCELED"
	} else if strings.EqualFold(o.Status, "completed") {
		order.Status = "FILLED"
	} else if strings.EqualFold(o.Status, "completedPart") {
		order.Status = "PARTIALLY_FILLED"
	} else if strings.EqualFold(o.Status, "open") {
		order.Status = "NEW"
		order.IsWorking = true
	} else if strings.EqualFold(o.Status, "expired") {
		order.Status = "EXPIRED"
	}
	order.Price = o.Price
	filledBuyTokenVolume := o.FilledBuyTokenVolume
	if len(filledBuyTokenVolume) == 0 {
		filledBuyTokenVolume = "0"
	}
	filledBuyTokenDec, _ := decimal.NewFromString(filledBuyTokenVolume)
	filledSellTokenVolume := o.FilledSellTokenVolume
	if len(filledSellTokenVolume) == 0 {
		filledSellTokenVolume = "0"
	}
	filledSellTokenDec, _ := decimal.NewFromString(filledSellTokenVolume)
	if o.IsBuy {
		v, _ := decimal.NewFromString(o.BuyToken.Volume)
		order.OrigQty = v.DivRound(pow10.Pow(decimal.NewFromInt32(o.BuyToken.Decimals)), 32).String()
		order.ExecutedQty = filledBuyTokenDec.DivRound(pow10.Pow(decimal.NewFromInt32(o.BuyToken.Decimals)), 32).String()
		v, _ = decimal.NewFromString(o.SellToken.Volume)
		order.OrigQuoteOrderQty = v.DivRound(pow10.Pow(decimal.NewFromInt32(o.SellToken.Decimals)), 32).String()
		order.CummulativeQuoteQty = filledSellTokenDec.DivRound(pow10.Pow(decimal.NewFromInt32(o.SellToken.Decimals)), 32).String()
	} else {
		v, _ := decimal.NewFromString(o.SellToken.Volume)
		order.OrigQty = v.DivRound(pow10.Pow(decimal.NewFromInt32(o.SellToken.Decimals)), 32).String()
		order.ExecutedQty = filledSellTokenDec.DivRound(pow10.Pow(decimal.NewFromInt32(o.SellToken.Decimals)), 32).String()
		v, _ = decimal.NewFromString(o.BuyToken.Volume)
		order.OrigQuoteOrderQty = v.DivRound(pow10.Pow(decimal.NewFromInt32(o.BuyToken.Decimals)), 32).String()
		order.CummulativeQuoteQty = filledBuyTokenDec.DivRound(pow10.Pow(decimal.NewFromInt32(o.BuyToken.Decimals)), 32).String()
	}
	return order
}

func ConvertDepth(d *model.Depth) (depth *binance.Depth) {
	if d == nil {
		return
	}
	depth = &binance.Depth{
		LastUpdateId: d.LastUpdateID,
		Bids:         d.Bids,
		Asks:         d.Asks,
	}
	return
}

func ConvertWithdraws(ws []*model.WithdrawalData) (withdraws []*binance.WithdrawHistory, err error) {
	var wh *binance.WithdrawHistory
	for _, w := range ws {
		if w != nil {
			wh, err = ConvertWithdraw(w)
			if err != nil {
				return
			}
			withdraws = append(withdraws, wh)
		}
	}
	return
}

func ConvertWithdraw(w *model.WithdrawalData) (withdraw *binance.WithdrawHistory, err error) {
	if w == nil {
		return
	}
	withdraw = &binance.WithdrawHistory{
		Address:      w.ToAddress,
		ApplyTime:    FormatTime(w.UpdateTime),
		Id:           w.WithdrawID,
		TransferType: 0,
		TxId:         w.TxHash,
		Network:      "ETH",
	}

	if strings.EqualFold(w.Status, "PROCESSING") {
		withdraw.Status = 4
	} else if strings.EqualFold(w.Status, "FAILED") {
		withdraw.Status = 5
	} else if strings.EqualFold(w.Status, "COMPLETED") {
		withdraw.Status = 6
	}
	if w.Token != nil {
		withdraw.Coin = w.Token.Symbol
		withdraw.Amount = GetAmountNew(w.Token.Volume, w.Token.Decimals)
	}
	if w.FeeToken != nil {
		withdraw.TransactionFeeCoin = w.FeeToken.Symbol
		withdraw.TransactionFee = GetAmountNew(w.FeeToken.Volume, w.FeeToken.Decimals)
	}
	return
}

func ConvertDeposits(ws []*model.DepositData) (deposits []*binance.DepositHistory, err error) {
	var wh *binance.DepositHistory
	for _, w := range ws {
		if w != nil {
			wh, err = ConvertDeposit(w)
			if err != nil {
				return
			}
			deposits = append(deposits, wh)
		}
	}
	return
}

func ConvertDeposit(w *model.DepositData) (deposit *binance.DepositHistory, err error) {
	if w == nil {
		return
	}
	deposit = &binance.DepositHistory{
		Address:      w.Owner,
		InsertTime:   int(w.CreateTime * 1000),
		Network:      "ETH",
		TransferType: 0,
		TxId:         w.L2TrxID,
		ConfirmTimes: "12/12",
	}

	if strings.EqualFold(w.Status, "PROCESSING") {
		deposit.Status = 6
	} else if strings.EqualFold(w.Status, "COMPLETED") {
		deposit.Status = 1
	} else if strings.EqualFold(w.Status, "SUCCESS") {
		deposit.Status = 2
	} else if strings.EqualFold(w.Status, "FAILED") {
		deposit.Status = 3
	} else if strings.EqualFold(w.Status, "CANCELED") {
		deposit.Status = 4
	}
	if w.Token != nil {
		deposit.Coin = w.Token.Symbol
		deposit.Amount = GetAmountNew(w.Token.Volume, w.Token.Decimals)
	}
	return
}

func ConvertTransfers(ws []*model.TransfersDataDetail) (transfers []*binance.TransfersData, err error) {
	var wh *binance.TransfersData
	for _, w := range ws {
		if w != nil {
			wh, err = ConvertTransfer(w)
			if err != nil {
				return
			}
			transfers = append(transfers, wh)
		}
	}
	return
}

func ConvertTransfer(t *model.TransfersDataDetail) (transfer *binance.TransfersData, err error) {
	if t == nil {
		return
	}

	transfer = &binance.TransfersData{
		Status:      "CONFIRMED",
		TranId:      t.TransferID,
		Timestamp:   int(t.CreateTime * 1000),
		FromAccount: t.Address,
		ToAccount:   t.ToAddress,
	}
	if t.Token != nil {
		var volume decimal.Decimal
		transfer.Asset = t.Token.Symbol
		pow10 := decimal.NewFromInt32(10)
		volume, err = decimal.NewFromString(t.Token.Volume)
		if err != nil {
			return
		}
		transfer.Amount = volume.DivRound(pow10.Pow(decimal.NewFromInt32(t.Token.Decimals)), 32).String()
	}
	return
}

func ConvertAccount(a *model.Account) (account *binance.Account, err error) {
	account = &binance.Account{
		CanTrade:    a.CanTrade,
		CanWithdraw: a.CanWithdraw,
		CanDeposit:  a.CanDeposit,
	}
	if a != nil {
		account.ID = a.ID
		account.Owner = a.Owner
		account.PublicKeyX = a.PublicKeyX
		account.PublicKeyY = a.PublicKeyY
		account.ReferrerId = a.ReferrerId
		account.Nonce = a.Nonce
	}
	account.AccountType = "SPOT"
	account.Permissions = []string{"SPOT"}
	return
}

func ConvertBalances(bs []*model.Balances) (balances []*binance.Balance, err error) {
	var (
		dL decimal.Decimal
		wL decimal.Decimal
		oL decimal.Decimal
	)
	if bs != nil && len(bs) > 0 {
		for _, b := range bs {
			bb := &binance.Balance{}
			if b.Token != nil {
				bb.Asset = b.Token.Symbol
				bb.TokenId = b.Token.TokenID
			}
			bb.Free = GetAmountNew(b.Balance, b.Token.Decimals)
			if dL, err = decimal.NewFromString(b.FrozenDepositBalance); err != nil {
				return
			}
			if wL, err = decimal.NewFromString(b.FrozenWithdrawBalance); err != nil {
				return
			}
			if oL, err = decimal.NewFromString(b.FrozenOrderBalance); err != nil {
				return
			}
			bb.Freeze = GetAmountNew(dL.Add(wL).Add(oL).String(), b.Token.Decimals)
			bb.Withdrawing = GetAmountNew(wL.String(), b.Token.Decimals)
			balances = append(balances, bb)
		}
	}
	return
}

func ConvertTrades(trades []*model.TradeData) (bTrades []*binance.UserTrade, err error) {
	var bt *binance.UserTrade
	for _, t := range trades {
		if t != nil {
			if bt, err = ConvertTrade(t); err == nil && bt != nil {
				bTrades = append(bTrades, bt)
			}
		}
	}
	return
}

func ConvertTrade(t *model.TradeData) (trade *binance.UserTrade, err error) {
	if t == nil {
		return
	}
	trade = &binance.UserTrade{
		BuyOrderId:  t.BuyOrderId,
		SellOrderId: t.SellOrderId,
	}
	trade.Symbol = GetSymbol(t.FilledBuyToken, t.FilledSellToken, t.IsBuy)
	trade.Id = t.TradeId
	trade.PairId = t.PairId
	trade.OrderId = t.OrderId
	trade.Price = t.Price
	if t.IsBuy {
		trade.BaseTokenId = uint32(t.FilledBuyToken.TokenID)
		trade.QuoteTokenId = uint32(t.FilledSellToken.TokenID)
		trade.Qty = GetAmountNew(t.FilledBuyToken.Volume, t.FilledBuyToken.Decimals)
		trade.QuoteQty = GetAmountNew(t.FilledSellToken.Volume, t.FilledSellToken.Decimals)
	} else {
		trade.BaseTokenId = uint32(t.FilledSellToken.TokenID)
		trade.QuoteTokenId = uint32(t.FilledBuyToken.TokenID)
		trade.Qty = GetAmountNew(t.FilledSellToken.Volume, t.FilledSellToken.Decimals)
		trade.QuoteQty = GetAmountNew(t.FilledBuyToken.Volume, t.FilledBuyToken.Decimals)
	}
	trade.Commission = GetAmountNew(t.FilledGasFeeToken.Volume, t.FilledGasFeeToken.Decimals)
	trade.CommissionAsset = t.FilledGasFeeToken.Symbol
	trade.IsMaker = t.IsMaker
	trade.IsBuyer = t.IsBuy
	trade.Time = t.CreateTime
	trade.GasFee = GetAmountNew(t.FilledGasFeeToken.Volume, t.FilledGasFeeToken.Decimals)
	trade.TradeFee = GetAmountNew(t.FilledFeeToken.Volume, t.FilledFeeToken.Decimals)
	trade.AccountId = uint64(t.AccountID)
	trade.R = t.R
	return
}

func ConvertTradesHistory(trades []*model.TradeData) (bTrades []*binance.TradeHistory, err error) {
	var bt *binance.TradeHistory
	for _, t := range trades {
		if t != nil {
			if bt, err = ConvertTradeHistory(t); err == nil && bt != nil {
				bTrades = append(bTrades, bt)
			}
		}
	}
	return
}

func ConvertTradeHistory(t *model.TradeData) (trade *binance.TradeHistory, err error) {
	if t == nil {
		return
	}
	trade = &binance.TradeHistory{}
	trade.Id = t.TradeId
	trade.Price = t.Price
	if t.IsBuy {
		trade.Qty = GetAmountNew(t.FilledBuyToken.Volume, t.FilledBuyToken.Decimals)
		trade.QuoteQty = GetAmountNew(t.FilledSellToken.Volume, t.FilledSellToken.Decimals)
	} else {
		trade.Qty = GetAmountNew(t.FilledSellToken.Volume, t.FilledSellToken.Decimals)
		trade.QuoteQty = GetAmountNew(t.FilledBuyToken.Volume, t.FilledBuyToken.Decimals)
	}
	if t.IsBuy && t.IsMaker {
		trade.IsBuyerMaker = true
	}
	trade.Time = t.CreateTime
	return
}

func ConvertTicker(t *model.Ticker) (ticker *binance.Ticker) {
	if t == nil {
		return
	}
	ticker = &binance.Ticker{}
	ticker.PriceChange = t.PriceChange
	ticker.PriceChangePercent = t.PriceChangePercent
	ticker.WeightedAvgPrice = t.WeightedAvgPrice
	ticker.PrevClosePrice = t.PrevClosePrice
	ticker.LastPrice = t.LastPrice
	ticker.LastQty = t.LastQty
	ticker.BidPrice = t.BidPrice
	ticker.BidQty = t.BidQty
	ticker.AskPrice = t.AskPrice
	ticker.AskQty = t.AskQty
	ticker.OpenPrice = t.OpenPrice
	ticker.HighPrice = t.HighPrice
	ticker.LowPrice = t.LowPrice
	ticker.Volume = t.Volume
	ticker.QuoteVolume = t.QuoteVolume
	ticker.OpenTime = t.OpenTime
	ticker.CloseTime = t.CloseTime
	ticker.FirstId = t.FirstId
	ticker.LastId = t.LastId
	ticker.Count = int(t.Count)
	ticker.MakerFee = t.MakerFee
	ticker.TakerFee = t.TakerFee
	ticker.PairId = t.PairId
	ticker.BaseToken = t.BaseToken
	ticker.QuoteToken = t.QuoteToken
	return
}

func ConvertTime(t *model.TimeData) (time *binance.Time) {
	if t == nil {
		return
	}
	time = &binance.Time{
		ServerTime: int(t.Timestamp * 1000),
	}
	return
}

func ConvertExchangeInfo(t *model.ExchangeInfo) (time *binance.ExchangeInfo) {
	if t == nil {
		return
	}
	time = &binance.ExchangeInfo{
		ChainID:               t.ChainID,
		MinLimitOrderUSDValue: t.MinLimitOrderUSDValue,
		Timezone:              t.Timezone,
		ServerTime:            t.ServerTime,
		RateLimits:            t.RateLimits,
	}
	return
}

func ConvertOrderAck(symbol string, orderId string, clientOrderId string) (orderAck *binance.OrderAck) {
	orderAck = &binance.OrderAck{
		Symbol:        symbol,
		OrderId:       orderId,
		ClientOrderId: clientOrderId,
	}
	return
}

func ConvertOrderUpdate(symbol string, orderType string, side string, quoteToken *model.TokenInfo, baseToken *model.TokenInfo, feeToken *model.TokenInfo, price string, o *model.OrderUpdate) (orderResult *binance.OrderResult, err error) {
	var (
		sellVolume     decimal.Decimal
		sellFillVolume = decimal.NewFromInt(0)
		buyVolume      decimal.Decimal
		buyFillVolume  = decimal.NewFromInt(0)
		tokens         = map[uint32]*model.TokenInfo{
			uint32(quoteToken.Id): quoteToken,
			uint32(baseToken.Id):  baseToken,
			uint32(feeToken.Id):   quoteToken,
		}
	)
	orderResult = &binance.OrderResult{
		Type:  orderType,
		Side:  side,
		Price: price,
	}
	if orderType == model.OrderTypeLimit {
		orderResult.TimeInForce = "GTC"
	} else if orderType == model.OrderTypeMarket {
		orderResult.TimeInForce = "IOC"
	}
	orderResult.Symbol = symbol
	orderResult.OrderId = o.OrderId
	orderResult.TransactTime = o.TransactionTime
	orderResult.ClientOrderId = o.ClientOrderId
	if sellVolume, err = decimal.NewFromString(o.SellToken.Volume); err != nil {
		return
	}
	if len(o.SellToken.FilledVolume) > 0 {
		if sellFillVolume, err = decimal.NewFromString(o.SellToken.FilledVolume); err != nil {
			return
		}
	}
	if buyVolume, err = decimal.NewFromString(o.BuyToken.Volume); err != nil {
		return
	}
	if len(o.BuyToken.FilledVolume) > 0 {
		if buyFillVolume, err = decimal.NewFromString(o.BuyToken.FilledVolume); err != nil {
			return
		}
	}

	if o.IsBuy {
		orderResult.OrigQty = buyVolume.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(tokens[o.BuyToken.TokenId].Decimals)), 32).String()
		orderResult.ExecutedQty = buyFillVolume.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(tokens[o.BuyToken.TokenId].Decimals)), 32).String()
		orderResult.CummulativeQuoteQty = sellFillVolume.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(tokens[o.SellToken.TokenId].Decimals)), 32).String()
	} else {
		orderResult.OrigQty = sellVolume.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(tokens[o.SellToken.TokenId].Decimals)), 32).String()
		orderResult.ExecutedQty = sellFillVolume.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(tokens[o.SellToken.TokenId].Decimals)), 32).String()
		orderResult.CummulativeQuoteQty = buyFillVolume.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(tokens[o.BuyToken.TokenId].Decimals)), 32).String()
	}

	if strings.EqualFold(o.Status, "canceled") {
		orderResult.Status = "CANCELED"
	} else if strings.EqualFold(o.Status, "completed") {
		orderResult.Status = "FILLED"
	} else if strings.EqualFold(o.Status, "completedPart") {
		orderResult.Status = "PARTIALLY_FILLED"
	} else if strings.EqualFold(o.Status, "open") || strings.EqualFold(o.Status, "processing") {
		orderResult.Status = "NEW"
	} else if strings.EqualFold(o.Status, "expired") {
		orderResult.Status = "EXPIRED"
	}
	return
}

func ConvertOrderTrade(quoteToken *model.TokenInfo, baseToken *model.TokenInfo, feeToken *model.TokenInfo, trades []*model.TradeResult) []*binance.OrderFill {
	var (
		pow10 = decimal.NewFromInt32(10)
	)
	if len(trades) > 0 {
		var fills []*binance.OrderFill
		for _, trade := range trades {
			fill := &binance.OrderFill{
				Price: trade.Price,
			}
			fillSellA, _ := decimal.NewFromString(trade.FillSA)
			fillSellB, _ := decimal.NewFromString(trade.FillSB)
			if trade.OrderA.TokenS == uint32(quoteToken.Id) {
				// buy
				fill.Qty = fillSellB.DivRound(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)), 32).String()
				fill.Commission = decimal.NewFromInt(int64(trade.OrderA.FeeBips)).Mul(fillSellB).DivRound(decimal.NewFromInt(10000), 32).DivRound(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)), 32).String()
				fill.CommissionAsset = baseToken.Symbol
			} else {
				// sell
				fill.Qty = fillSellA.DivRound(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)), 32).String()
				fill.Commission = decimal.NewFromInt(int64(trade.OrderA.FeeBips)).Mul(fillSellB).DivRound(decimal.NewFromInt(10000), 32).DivRound(pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)), 32).String()
				fill.CommissionAsset = quoteToken.Symbol
			}
			fill.GasFeeCommissionAsset = feeToken.Symbol
			if feeVolume, err := decimal.NewFromString(trade.OrderA.Fee); err == nil {
				fill.GasFeeCommission = feeVolume.DivRound(pow10.Pow(decimal.NewFromInt32(feeToken.Decimals)), 32).String()
			}
			fills = append(fills, fill)
		}
		return fills
	}
	return nil
}

func ConvertOrderFull(symbol string, orderType string, side string, quoteToken *model.TokenInfo, baseToken *model.TokenInfo, feeToken *model.TokenInfo, price string, o *model.OrderUpdate) (orderResult *binance.OrderResult, err error) {
	orderResult, err = ConvertOrderUpdate(symbol, orderType, side, quoteToken, baseToken, feeToken, price, o)
	if err != nil {
		return
	}
	orderResult.Fills = ConvertOrderTrade(quoteToken, baseToken, feeToken, o.Trades)
	return
}

func ConvertListenerKey(data model.ListenKeyData) *binance.ListenKeyData {
	return &binance.ListenKeyData{
		ListenKey: data.ListenKey,
	}
}

func ConvertKlinePayload(payload *model.KlinePayload) (p *binance.KlinePayload, err error) {
	baseToken := conf.Conf.GetTokenInfoById(payload.B)
	if baseToken == nil {
		err = errors.New("not config token")
		return
	}
	quoteToken := conf.Conf.GetTokenInfoById(payload.U)
	if quoteToken == nil {
		err = errors.New("not config token")
		return
	}
	p = &binance.KlinePayload{
		E1: payload.E1,
		S:  baseToken.Symbol + quoteToken.Symbol,
	}
	p.E = payload.E
	p.K = binance.KlineData{
		T:  payload.K.T,
		T1: payload.K.T1,
		S:  baseToken.Symbol + quoteToken.Symbol,
		I:  conf.GetInterval(payload.K.I),
		F:  payload.K.F,
		L:  payload.K.L,
		O:  payload.K.O,
		C:  payload.K.C,
		H:  payload.K.H,
		L1: payload.K.L1,
		V:  payload.K.V,
		N:  payload.K.N,
		X:  payload.K.X,
		Q:  payload.K.Q,
		V1: payload.K.V1,
		Q1: payload.K.Q1,
	}
	return
}

func ConvertTradePayload(payload *model.TradePayload) (p *binance.TradePayload, err error) {
	baseToken := conf.Conf.GetTokenInfoById(payload.B)
	if baseToken == nil {
		err = errors.New("not config token")
		return
	}
	quoteToken := conf.Conf.GetTokenInfoById(payload.U)
	if quoteToken == nil {
		err = errors.New("not config token")
		return
	}
	p = &binance.TradePayload{
		E1: payload.E1,
		S:  baseToken.Symbol + quoteToken.Symbol,
		T:  payload.T,
		P:  payload.P,
		Q:  payload.Q,
		B:  payload.B1,
		A:  payload.A,
		T1: payload.T1,
		M:  payload.M,
	}
	p.E = payload.E
	return
}

func ConvertTickerPayload(payload *model.TickerPayload) (p *binance.TickerPayload, err error) {
	baseToken := conf.Conf.GetTokenInfoById(payload.B)
	if baseToken == nil {
		err = errors.New("not config token")
		return
	}
	quoteToken := conf.Conf.GetTokenInfoById(payload.U)
	if quoteToken == nil {
		err = errors.New("not config token")
		return
	}
	p = &binance.TickerPayload{
		E1: payload.E1,
		S:  baseToken.Symbol + quoteToken.Symbol,
		P:  payload.P,
		P1: payload.P1,
		W:  payload.W,
		X:  payload.X,
		C:  payload.C,
		Q:  payload.Q,
		B:  payload.B1,
		B1: payload.I,
		A:  payload.A,
		A1: payload.A1,
		O:  payload.O,
		H:  payload.H,
		L:  payload.L,
		V:  payload.V,
		Q1: payload.Q1,
		O1: payload.O1,
		C1: payload.C1,
		F:  payload.F,
		L1: payload.L1,
		N:  payload.N,
	}
	p.E = payload.E
	return
}

func ConvertBookTickerPayload(payload *model.BookTickerPayload) (p *binance.BookTickerPayload, err error) {
	baseToken := conf.Conf.GetTokenInfoById(payload.B)
	if baseToken == nil {
		err = errors.New("not config token")
		return
	}
	quoteToken := conf.Conf.GetTokenInfoById(payload.U)
	if quoteToken == nil {
		err = errors.New("not config token")
		return
	}
	p = &binance.BookTickerPayload{
		S:  baseToken.Symbol + quoteToken.Symbol,
		B:  payload.B1,
		B1: payload.I,
		A:  payload.A,
		A1: payload.A1,
	}
	p.E = payload.E
	return
}

func ConvertDepthPayload(payload *model.DepthPayload) (p *binance.DepthPayload, err error) {
	p = &binance.DepthPayload{
		LastUpdateId: payload.LastUpdateId,
		Bids:         payload.Bids,
		Asks:         payload.Asks,
	}
	return
}

func ConvertDepthUploadPayload(payload *model.DepthUpdatePayload) (p *binance.DepthUploadPayload, err error) {
	baseToken := conf.Conf.GetTokenInfoById(payload.B)
	if baseToken == nil {
		err = errors.New("not config token")
		return
	}
	quoteToken := conf.Conf.GetTokenInfoById(payload.Q)
	if quoteToken == nil {
		err = errors.New("not config token")
		return
	}
	p = &binance.DepthUploadPayload{
		E1: payload.E1,
		S:  baseToken.Symbol + quoteToken.Symbol,
		U:  payload.U1,
		U1: payload.U,
		B:  payload.B1,
		A:  payload.A,
	}
	p.E = payload.E
	return
}

func ConvertOutboundAccountPositionPayload(payload *model.OutboundAccountPositionPayload) (p *binance.OutboundAccountPositionPayload, err error) {
	p = &binance.OutboundAccountPositionPayload{
		Time: payload.Time,
		U:    payload.U,
	}
	if payload.B != nil {
		p.B = []*binance.OutboundAccountBalance{}
		bOut := &binance.OutboundAccountBalance{
			Symbol:  payload.B.Symbol,
			F:       payload.B.F,
			L:       payload.B.L,
			TokenId: payload.B.TokenId,
		}
		p.B = append(p.B, bOut)
	}
	p.E = payload.E
	return
}

func ConvertBalanceUpdatePayload(payload *model.BalanceUpdatePayload) (p *binance.BalanceUpdatePayload, err error) {
	p = &binance.BalanceUpdatePayload{
		Time:    payload.Time,
		Balance: payload.Balance,
	}
	p.E = payload.E
	p.Symbol = payload.Symbol
	p.T = payload.Time
	p.TokenId = payload.TokenId
	return
}

func ConvertExecutionReportPayload(payload *model.ExecutionReportPayload) (p *binance.ExecutionReportPayload, err error) {
	p = &binance.ExecutionReportPayload{
		Time:             payload.Time,
		OrderID:          payload.OrderID,
		Ts:               payload.Ts,
		CreateTime:       payload.CreateTime,
		TradingFee:       payload.TradingFee,
		Price:            payload.Price,
		IsOrderBook:      payload.IsOrderBook,
		TradingFeeSymbol: payload.TradingFeeSymbol,
		PairName:         payload.PairName,
		OriginalVolume:   payload.OriginalVolume,
		QuoteOrderQty:    payload.QuoteOrderQty,
		TotalDealVolume:  payload.TotalDealVolume,
		TotalDealAmount:  payload.TotalDealAmount,
		ClientOrderId:    payload.ClientOrderId,
		LastDealVolume:   payload.LastDealVolume,
		DealPrice:        payload.DealPrice,
		TradeID:          payload.TradeID,
		IsMaker:          payload.IsMaker,
		LastDealAmount:   payload.LastDealAmount,
		WorkingTime:      payload.WorkingTime,
		PairID:           payload.PairID,
		BaseTokenID:      payload.BaseTokenID,
		QuoteTokenID:     payload.QuoteTokenID,
	}
	p.E = payload.E
	if payload.OrderType == 0 {
		p.OrderType = "LIMIT"
	} else if payload.OrderType == 1 {
		p.OrderType = "MARKET"
	}
	if payload.IsBuy {
		p.Side = "BUY"
	} else {
		p.Side = "SELL"
	}
	if payload.CancelReason == 1 {
		p.CancelReason = "cancellation of spot orders by user"
	} else if payload.CancelReason == 2 {
		p.CancelReason = "cancellation of grid orders by user"
	} else if payload.CancelReason == 3 {
		p.CancelReason = "insufficient gas auto-cancel"
	} else if payload.CancelReason == 4 {
		p.CancelReason = "orders with very small remaining quantities are automatically cancelled"
	} else if payload.CancelReason == 5 {
		p.CancelReason = "insufficient counterparties in the order book"
	} else if payload.CancelReason == 6 {
		p.CancelReason = "expired spot limit order"
	} else if payload.CancelReason == 7 {
		p.CancelReason = "expired grid policy"
	} else if payload.CancelReason == 8 {
		p.CancelReason = "forced withdrawal"
	} else if payload.CancelReason == 9 {
		p.CancelReason = "account updating"
	} else if payload.CancelReason == 10 {
		p.CancelReason = "grid strategies are cancelled due to quantity errors"
	} else if payload.CancelReason == 11 {
		p.CancelReason = "cancel order on chain"
	} else if payload.CancelReason == 12 {
		p.CancelReason = "price protection triggered"
	}

	if strings.EqualFold(payload.Status, "OPEN") {
		if len(payload.LastDealVolume) > 0 && payload.LastDealVolume != "0" {
			p.X = "trade"
			p.X1 = "PARTIALLY_FILLED"
		} else {
			p.X1 = "NEW"
		}
	} else if payload.CancelReason == 3 || payload.CancelReason == 4 || payload.CancelReason == 5 || payload.CancelReason == 6 || payload.CancelReason == 7 || payload.CancelReason == 8 || payload.CancelReason == 9 || payload.CancelReason == 10 || payload.CancelReason == 12 {
		p.X1 = "cancel"
		p.X1 = "EXPIRED"
	} else if strings.EqualFold(payload.Status, "CANCELED") || payload.CancelReason == 1 || payload.CancelReason == 2 || payload.CancelReason == 11 {
		p.X = "cancel"
		p.X1 = "CANCELED"
	} else if strings.EqualFold(payload.Status, "COMPLETED") {
		p.X = "trade"
		if payload.Completed {
			p.X1 = "FILLED"
		} else {
			p.X1 = "PARTIALLY_FILLED"
		}
	} else if strings.EqualFold(payload.Status, "completedPart") {
		p.X = "trade"
		p.X1 = "PARTIALLY_FILLED"
	} else {
		p.X1 = strings.ToUpper(payload.Status)
	}
	return
}

//func ConvertOffChainFee(fees *model.OffChainFee) (offChainFee *binance.OffChainFee, err error) {
//	if fees == nil {
//		return
//	}
//	offChainFee = &binance.OffChainFee{}
//
//	offChainFee.UpdateAccountGasFees, err = ConvertGasFee(fees.UpdateAccountGasFees)
//
//	offChainFee.WithdrawalGasFees, err = ConvertGasFee(fees.WithdrawalGasFees)
//
//	offChainFee.WithdrawalOtherGasFees, err = ConvertGasFee(fees.WithdrawalOtherGasFees)
//
//	offChainFee.TransferGasFees, err = ConvertGasFee(fees.TransferGasFees)
//
//	offChainFee.TransferNoIDGasFees, err = ConvertGasFee(fees.TransferNoIDGasFees)
//
//	offChainFee.OrderGasFees, err = ConvertGasFee(fees.OrderGasFees)
//	offChainFee.OrderGasMultiple = fees.OrderGasMultiple
//
//	offChainFee.AddPairGasFees, err = ConvertGasFee(fees.AddPairGasFees)
//
//	offChainFee.MiningGasFees, err = ConvertGasFee(fees.MiningGasFees)
//
//	offChainFee.OnChainCancelOrderGasFees, err = ConvertGasFee(fees.OnChainCancelOrderGasFees)
//
//	return
//}

//func ConvertGasFee(fees []*model.ShowTokenData) (gasFees []*binance.GasFee, err error) {
//	if len(fees) == 0 {
//		return
//	}
//	var (
//		v decimal.Decimal
//	)
//	for _, f := range fees {
//		if f == nil {
//			continue
//		}
//		gasFee := &binance.GasFee{
//			TokenId:  f.TokenID,
//			Symbol:   f.Symbol,
//			Decimals: f.Decimals,
//			Volume:   f.Volume,
//		}
//		v, err = decimal.NewFromString(f.Volume)
//		if err != nil {
//			v = conf.Zero
//		}
//		gasFee.Quantity = v.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(f.Decimals)), 32).String()
//		gasFees = append(gasFees, gasFee)
//	}
//	return
//}

func ConvertGasFees(fees *model.GasFees, tokens []*model.TokenInfo) (offChainFee *binance.OffChainFee, err error) {
	if fees == nil {
		return
	}
	offChainFee = &binance.OffChainFee{}

	offChainFee.UpdateAccountGasFees, err = ConvertGasFee(fees.UpdateAccountGasFees, tokens)

	offChainFee.WithdrawalGasFees, err = ConvertGasFee(fees.WithdrawalGasFees, tokens)

	offChainFee.EstimatedWithdrawalGasFees, err = ConvertGasFee(fees.EstimatedWithdrawalGasFees, tokens)

	offChainFee.TransferGasFees, err = ConvertGasFee(fees.TransferGasFees, tokens)

	offChainFee.TransferNoIDGasFees, err = ConvertGasFee(fees.TransferNoIDGasFees, tokens)

	offChainFee.OrderGasFees, err = ConvertGasFee(fees.OrderGasFees, tokens)
	if fees.OrderGasFees != nil {
		offChainFee.OrderGasMultiple = uint32(fees.OrderGasFees.OrderGasMultiple)
	}

	offChainFee.AddPairGasFees, err = ConvertGasFee(fees.AddPairGasFees, tokens)

	offChainFee.MiningGasFees, err = ConvertGasFee(fees.MiningGasFees, tokens)

	offChainFee.OnChainCancelOrderGasFees, err = ConvertGasFee(fees.OnChainCancelOrderGasFees, tokens)

	return
}

func ConvertGasFee(fees *model.GasFee, tokens []*model.TokenInfo) (gasFees []*binance.GasFee, err error) {
	if fees == nil || len(fees.Tokens) == 0 {
		return
	}
	var (
		v decimal.Decimal
	)
	for _, f := range fees.Tokens {
		var decimals int32
		if f == nil {
			continue
		}
		for _, token := range tokens {
			if f.TokenID == uint64(token.Id) {
				decimals = token.Decimals
				break
			}
		}
		if decimals <= 0 {
			continue
		}
		gasFee := &binance.GasFee{
			TokenId:  f.TokenID,
			Symbol:   f.Symbol,
			Volume:   f.Volume,
			Decimals: decimals,
		}
		v, err = decimal.NewFromString(f.Volume)
		if err != nil {
			v = conf.Zero
		}
		gasFee.Quantity = v.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(decimals)), 32).String()
		gasFees = append(gasFees, gasFee)
	}
	return
}

func ConvertTradeFee(fees []*model.TradeFee) (tradeFees []*binance.TradeFee) {
	if len(fees) == 0 {
		return
	}
	for _, f := range fees {
		if f != nil {
			tradeFees = append(tradeFees, &binance.TradeFee{
				MakerCommission: f.MakerCommission,
				TakerCommission: f.TakerCommission,
				Symbol:          f.BaseToken.Symbol + f.QuoteToken.Symbol,
			})
		}
	}
	return
}

func ConvertTokenInfoToTokenData(info *model.TokenInfo) *binance.ShowTokenData {
	data := &binance.ShowTokenData{
		TokenID:         uint64(info.Id),
		Chain:           info.Chain,
		Code:            info.Code,
		Symbol:          info.Symbol,
		Decimals:        info.Decimals,
		Volume:          "0",
		ShowDecimals:    info.ShowDecimal,
		IsQuotableToken: info.IsQuotableToken,
		IsGasToken:      info.IsGasToken,
		IsListToken:     info.IsListToken,
		Active:          info.Active,
		IsTrustedToken:  info.IsTrustedToken,
		Priority:        info.Priority,
	}
	return data
}

func ConvertPairPrice(symbol map[uint64]string, t []*model.PairsPricesRes) (tickers []*binance.PairPrice) {
	if len(t) == 0 {
		return
	}
	for _, p := range t {
		ticker := &binance.PairPrice{
			Price:  p.Price,
			Symbol: symbol[p.PairID],
			PairId: p.PairID,
		}
		tickers = append(tickers, ticker)
	}

	return
}

func ConvertBookTicker(symbol string, t *model.BookTickerData) (ticker *binance.BookTicker) {
	if t == nil {
		return
	}
	ticker = &binance.BookTicker{
		Symbol: symbol,
	}
	if len(t.Asks) > 0 {
		ticker.AskPrice = t.Asks[0][0]
		ticker.AskQty = t.Asks[0][1]
	}
	if len(t.Bids) > 0 {
		ticker.BidPrice = t.Bids[0][0]
		ticker.BidQty = t.Bids[0][1]
	}
	return
}

func ConvertGasFeeToken(fees *binance.OffChainFee) (gasFee *binance.GasFeeToken) {
	if fees == nil {
		return
	}
	gasFee = &binance.GasFeeToken{
		WithdrawalGasFees:         fees.WithdrawalGasFees,
		TransferGasFees:           fees.TransferGasFees,
		TransferNoIDGasFees:       fees.TransferNoIDGasFees,
		OnChainCancelOrderGasFees: fees.OnChainCancelOrderGasFees,
	}
	return
}
