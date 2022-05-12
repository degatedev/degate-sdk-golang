package lib

import (
	"errors"
	"strings"

	"github.com/degatedev/degatesdk/conf"
	"github.com/degatedev/degatesdk/degate/binance"
	"github.com/degatedev/degatesdk/degate/model"
	"github.com/degatedev/degatesdk/util"

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
		zero  = decimal.NewFromInt(0)
	)
	order = &binance.Order{}
	order.Symbol = o.GetSymbol()
	order.OrderId = o.OrderID
	order.OrderListId = -1
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
		if o.FillAmountBOrs {
			fv, _ := decimal.NewFromString(o.FilledBuyTokenVolume)
			v, _ := decimal.NewFromString(o.BuyToken.Volume)
			if fv.Equals(v) {
				order.Status = "FILLED"
			} else {
				order.Status = "PARTIALLY_FILLED"
			}
		} else {
			fv, _ := decimal.NewFromString(o.FilledSellTokenVolume)
			v, _ := decimal.NewFromString(o.SellToken.Volume)
			if fv.Equal(v) {
				order.Status = "FILLED"
			} else {
				order.Status = "PARTIALLY_FILLED"
			}
		}
	} else if strings.EqualFold(o.Status, "open") {
		order.Status = "NEW"
		order.IsWorking = true
	} else {
		order.Status = ""
	}

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
		if o.OrderType == 1 {
			if filledBuyTokenDec.GreaterThan(zero) {
				order.Price = util.CalculatePrice(filledSellTokenVolume, filledBuyTokenVolume, decimal.NewFromInt(10).Pow(decimal.NewFromInt32(o.BuyToken.Decimals-o.SellToken.Decimals))).String()
			}
		} else {
			order.Price = util.CalculatePrice(o.SellToken.Volume, o.BuyToken.Volume, decimal.NewFromInt(10).Pow(decimal.NewFromInt32(o.BuyToken.Decimals-o.SellToken.Decimals))).String()
		}
	} else {
		v, _ := decimal.NewFromString(o.SellToken.Volume)
		order.OrigQty = v.DivRound(pow10.Pow(decimal.NewFromInt32(o.SellToken.Decimals)), 32).String()
		order.ExecutedQty = filledSellTokenDec.DivRound(pow10.Pow(decimal.NewFromInt32(o.SellToken.Decimals)), 32).String()
		v, _ = decimal.NewFromString(o.BuyToken.Volume)
		order.OrigQuoteOrderQty = v.DivRound(pow10.Pow(decimal.NewFromInt32(o.BuyToken.Decimals)), 32).String()
		order.CummulativeQuoteQty = filledBuyTokenDec.DivRound(pow10.Pow(decimal.NewFromInt32(o.BuyToken.Decimals)), 32).String()
		if o.OrderType == 1 {
			if filledSellTokenDec.GreaterThan(zero) {
				order.Price = util.CalculatePrice(filledBuyTokenVolume, filledSellTokenVolume, decimal.NewFromInt(10).Pow(decimal.NewFromInt32(o.SellToken.Decimals-o.BuyToken.Decimals))).String()
			}
		} else {
			order.Price = util.CalculatePrice(o.BuyToken.Volume, o.SellToken.Volume, decimal.NewFromInt(10).Pow(decimal.NewFromInt32(o.SellToken.Decimals-o.BuyToken.Decimals))).String()
		}
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
	}

	if strings.EqualFold(w.Status, "processed") {
		withdraw.Status = 2
	} else if strings.EqualFold(w.Status, "PROCESSING") {
		withdraw.Status = 4
	} else if strings.EqualFold(w.Status, "FAILED") {
		withdraw.Status = 5
	} else if strings.EqualFold(w.Status, "COMPLETED") {
		withdraw.Status = 6
	}
	if w.Token != nil {
		withdraw.Coin = w.Token.Symbol
		withdraw.Amount, err = GetAmount(w.Token.Volume, w.Token.Decimals, conf.EffectiveDecimal)
		if err != nil {
			return
		}
	}
	if w.FeeToken != nil {
		withdraw.TransactionFeeCoin = w.FeeToken.Symbol
		withdraw.TransactionFee, err = GetAmount(w.FeeToken.Volume, w.FeeToken.Decimals, conf.EffectiveDecimal)
		if err != nil {
			return
		}
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

// ConvertDeposit (0:pending,6: credited but cannot withdraw, 1:success)
func ConvertDeposit(w *model.DepositData) (deposit *binance.DepositHistory, err error) {
	if w == nil {
		return
	}
	deposit = &binance.DepositHistory{
		Address:      w.Owner,
		InsertTime:   int(w.CreateTime * 1000),
		TransferType: 0,
		TxId:         w.L2TrxID,
	}

	if strings.EqualFold(w.Status, "PROCESSING") {
		deposit.Status = 0
	} else if strings.EqualFold(w.Status, "COMPLETED") {
		deposit.Status = 1
	} else if strings.EqualFold(w.Status, "PROCESSED") {
		deposit.Status = 6
	}
	if w.Token != nil {
		deposit.Coin = w.Token.Symbol
		deposit.Amount, err = GetAmount(w.Token.Volume, w.Token.Decimals, conf.EffectiveDecimal)
		if err != nil {
			return
		}
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
	account = &binance.Account{}
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
			}
			if bb.Free, err = GetAmount(b.Balance, b.Token.Decimals, conf.EffectiveDecimal); err != nil {
				return
			}
			if dL, err = decimal.NewFromString(b.FrozenDepositBalance); err != nil {
				return
			}
			if wL, err = decimal.NewFromString(b.FrozenWithdrawBalance); err != nil {
				return
			}
			if oL, err = decimal.NewFromString(b.FrozenOrderBalance); err != nil {
				return
			}
			if bb.Freeze, err = GetAmount(dL.Add(wL).Add(oL).String(), b.Token.Decimals, conf.EffectiveDecimal); err != nil {
				return
			}
			balances = append(balances, bb)
		}
	}
	return
}

func ConvertTrades(trades []*model.TradeData) (bTrades []*binance.Trade, err error) {
	var bt *binance.Trade
	for _, t := range trades {
		if t != nil {
			if bt, err = ConvertTrade(t); err == nil && bt != nil {
				bTrades = append(bTrades, bt)
			}
		}
	}
	return
}

func ConvertTrade(t *model.TradeData) (trade *binance.Trade, err error) {
	if t == nil {
		return
	}
	trade = &binance.Trade{}
	trade.Symbol = GetSymbol(t.FilledBuyToken, t.FilledSellToken, t.IsBuy)
	trade.Id = t.ID
	trade.PairId = t.PairId
	trade.TradeId = t.TradeId
	trade.OrderId = t.OrderId
	trade.OrderListId = -1
	if t.IsBuy {
		trade.BaseTokenId = uint32(t.FilledBuyToken.TokenID)
		trade.QuoteTokenId = uint32(t.FilledSellToken.TokenID)
		trade.Price = util.CalculatePrice(t.FilledSellToken.Volume, t.FilledBuyToken.Volume, decimal.NewFromInt(10).Pow(decimal.NewFromInt32(t.FilledBuyToken.Decimals-t.FilledSellToken.Decimals))).String()
		if trade.Qty, err = GetAmount(t.FilledBuyToken.Volume, t.FilledBuyToken.Decimals, conf.EffectiveDecimal); err != nil {
			return
		}
		if trade.QuoteQty, err = GetAmount(t.FilledSellToken.Volume, t.FilledSellToken.Decimals, conf.EffectiveDecimal); err != nil {
			return
		}
	} else {
		trade.BaseTokenId = uint32(t.FilledSellToken.TokenID)
		trade.QuoteTokenId = uint32(t.FilledBuyToken.TokenID)
		trade.Price = util.CalculatePrice(t.FilledBuyToken.Volume, t.FilledSellToken.Volume, decimal.NewFromInt(10).Pow(decimal.NewFromInt32(t.FilledSellToken.Decimals-t.FilledBuyToken.Decimals))).String()
		if trade.Qty, err = GetAmount(t.FilledSellToken.Volume, t.FilledSellToken.Decimals, conf.EffectiveDecimal); err != nil {
			return
		}
		if trade.QuoteQty, err = GetAmount(t.FilledBuyToken.Volume, t.FilledBuyToken.Decimals, conf.EffectiveDecimal); err != nil {
			return
		}
	}

	if trade.Commission, err = GetAmount(t.FilledGasFeeToken.Volume, t.FilledGasFeeToken.Decimals, conf.EffectiveDecimal); err != nil {
		return
	}
	trade.CommissionAsset = t.FilledGasFeeToken.Symbol
	trade.IsMaker = t.IsMaker
	trade.IsBuyer = t.IsBuy
	// create time is the order's tradetime
	trade.Time = t.CreateTime
	trade.GasFee, _ = GetAmount(t.FilledGasFeeToken.Volume, t.FilledGasFeeToken.Decimals, conf.EffectiveDecimal)
	trade.TradeFee, _ = GetAmount(t.FilledFeeToken.Volume, t.FilledFeeToken.Decimals, conf.EffectiveDecimal)
	trade.AccountId = uint64(t.AccountID)
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
	trade.Id = t.ID
	if t.IsBuy {
		trade.Price = util.CalculatePrice(t.FilledSellToken.Volume, t.FilledBuyToken.Volume, decimal.NewFromInt(10).Pow(decimal.NewFromInt32(t.FilledBuyToken.Decimals-t.FilledSellToken.Decimals))).String()
		if trade.Qty, err = GetAmount(t.FilledBuyToken.Volume, t.FilledBuyToken.Decimals, conf.EffectiveDecimal); err != nil {
			return
		}
		if trade.QuoteQty, err = GetAmount(t.FilledSellToken.Volume, t.FilledSellToken.Decimals, conf.EffectiveDecimal); err != nil {
			return
		}
	} else {
		trade.Price = util.CalculatePrice(t.FilledBuyToken.Volume, t.FilledSellToken.Volume, decimal.NewFromInt(10).Pow(decimal.NewFromInt32(t.FilledSellToken.Decimals-t.FilledBuyToken.Decimals))).String()
		if trade.Qty, err = GetAmount(t.FilledSellToken.Volume, t.FilledSellToken.Decimals, conf.EffectiveDecimal); err != nil {
			return
		}
		if trade.QuoteQty, err = GetAmount(t.FilledBuyToken.Volume, t.FilledBuyToken.Decimals, conf.EffectiveDecimal); err != nil {
			return
		}
	}
	if t.IsBuy && t.IsMaker {
		trade.IsBuyerMaker = true
	}
	trade.Time = t.CreateTime * 1000
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
	ticker.OpenPrice = t.HighPrice
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
		ChainID:              t.ChainID,
		ExchangeAddress:      t.ExchangeAddress,
		DepositAddress:       t.DepositAddress,
		WithdrawalsAddress:   t.WithdrawalsAddress,
		SpotTradeAddress:     t.SpotTradeAddress,
		OrderCancelAddress:   t.OrderCancelAddress,
		OrderEffectiveDigits: t.OrderEffectiveDigits,
		MinOrderPrice:        t.MinOrderPrice,
		MaxFeeBipsMax:        t.MaxFeeBipsMax,
		Timezone:             t.Timezone,
		ServerTime:           t.ServerTime,
		OrderMaxVolume:       t.OrderMaxVolume,
	}
	return
}

func ConvertOrderAck(symbol string, orderId string, clientOrderId string) (orderAck *binance.OrderAck) {
	orderAck = &binance.OrderAck{
		Symbol:        symbol,
		OrderId:       orderId,
		OrderListId:   -1,
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
			uint32(baseToken.Id):  quoteToken,
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
	orderResult.OrderListId = -1
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
		if o.FillAmountBors {
			if buyFillVolume.GreaterThanOrEqual(buyVolume) {
				orderResult.Status = "FILLED"
			} else {
				orderResult.Status = "PARTIALLY_FILLED"
			}
		} else {
			if sellFillVolume.GreaterThanOrEqual(sellVolume) {
				orderResult.Status = "FILLED"
			} else {
				orderResult.Status = "PARTIALLY_FILLED"
			}
		}
	} else if strings.EqualFold(o.Status, "open") || strings.EqualFold(o.Status, "processing") {
		orderResult.Status = "NEW"
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
			var price decimal.Decimal
			fill := &binance.OrderFill{}
			fillSellA, _ := decimal.NewFromString(trade.FillSA)
			fillSellB, _ := decimal.NewFromString(trade.FillSB)
			if trade.OrderA.TokenS == uint32(quoteToken.Id) {
				if trade.IsStatble {
					price = util.GetEffectivePriceRound(fillSellA.DivRound(fillSellB, 32).Mul(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals-quoteToken.Decimals))), trade.IsStatble, conf.OrderStabilityEffectiveDigitsGreaterThan1, conf.OrderStabilityEffectiveDigitsLessThan1)
				} else {
					price = util.GetEffectivePriceRound(fillSellA.DivRound(fillSellB, 32).Mul(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals-quoteToken.Decimals))), trade.IsStatble, conf.OrderEffectiveDigitsGreaterThan10000, conf.OrderEffectiveDigitsLessThan10000)
				}
				fill.Price = price.String()
				fill.Qty = fillSellB.DivRound(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)), 32).String()
				fill.Commission = decimal.NewFromInt(int64(trade.OrderA.FeeBips)).Mul(fillSellB).DivRound(decimal.NewFromInt(10000), 32).DivRound(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)), 32).String()
				fill.CommissionAsset = baseToken.Symbol
			} else {
				if trade.IsStatble {
					price = util.GetEffectivePriceRound(fillSellB.DivRound(fillSellA, 32).Mul(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals-quoteToken.Decimals))), trade.IsStatble, conf.OrderStabilityEffectiveDigitsGreaterThan1, conf.OrderStabilityEffectiveDigitsLessThan1)
				} else {
					price = util.GetEffectivePriceRound(fillSellB.DivRound(fillSellA, 32).Mul(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals-quoteToken.Decimals))), trade.IsStatble, conf.OrderEffectiveDigitsGreaterThan10000, conf.OrderEffectiveDigitsLessThan10000)
				}
				fill.Price = price.String()
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
		E1: payload.E1,
		U:  payload.U,
	}
	if payload.B != nil {
		p.B = []*binance.OutboundAccountBalance{}
		bOut := &binance.OutboundAccountBalance{
			F: payload.B.F,
			L: payload.B.L,
		}
		token := conf.Conf.GetTokenInfoById(payload.B.A)
		if token != nil {
			bOut.A = token.Symbol
		}
		p.B = append(p.B, bOut)
	}
	p.E = payload.E
	return
}

func ConvertBalanceUpdatePayload(payload *model.BalanceUpdatePayload) (p *binance.BalanceUpdatePayload, err error) {
	p = &binance.BalanceUpdatePayload{
		E1: payload.E1,
		D:  payload.V,
	}
	token := conf.Conf.GetTokenInfoById(payload.A)
	if token != nil {
		p.A = token.Symbol
	}
	p.E = payload.E
	return
}

func ConvertExecutionReportPayload(payload *model.ExecutionReportPayload) (p *binance.ExecutionReportPayload, err error) {
	p = &binance.ExecutionReportPayload{
		E1: payload.E1,
		F:  "GTC",
		G:  -1,
		I:  payload.O,
		T:  payload.U1,
		O1: payload.R,
		N:  payload.K,
	}
	p.E = payload.E
	if payload.T == 0 {
		p.O = "LIMIT"
	} else if payload.T == 1 {
		p.O = "MARKET"
	}
	if strings.EqualFold(payload.U, "OPEN") {
		p.X1 = "NEW"
		p.W = true
	} else if strings.EqualFold(payload.U, "CANCELED") {
		p.X1 = "CANCELED"
	} else if strings.EqualFold(payload.U, "COMPLETED") {
		if fDec, err := decimal.NewFromString(payload.F1); err == nil {
			if dDec, err := decimal.NewFromString(payload.D); err == nil {
				if fDec.GreaterThanOrEqual(dDec) {
					p.X1 = "COMPLETED"
				} else {
					p.X1 = "PARTIALLY_FILLED"
				}
			}
		}
	}
	var (
		pow10       = decimal.NewFromInt(10)
		baseVolume  string
		quoteVolume string
		baseDec     decimal.Decimal
		quoteDec    decimal.Decimal
		feeVolume   decimal.Decimal
		excVolume   string
		excDec      decimal.Decimal
		feeToken    *model.TokenInfo
		sellToken   *model.TokenInfo
		buyToken    *model.TokenInfo
		baseToken   *model.TokenInfo
		quoteToken  *model.TokenInfo
	)

	if feeToken = conf.Conf.GetTokenInfoById(payload.F); feeToken != nil {
		p.N1 = feeToken.Symbol
		if len(payload.L) > 0 {
			if feeVolume, err = decimal.NewFromString(payload.L); err != nil {
				return
			}
		}
		p.N = feeVolume.DivRound(pow10.Pow(decimal.NewFromInt32(feeToken.Decimals)), 32).String()
	}

	sellToken = conf.Conf.GetTokenInfoById(payload.S1)
	buyToken = conf.Conf.GetTokenInfoById(payload.B)
	if payload.B1 {
		p.S1 = "BUY"
		baseVolume = payload.D1
		quoteVolume = payload.D
		excVolume = payload.C
		baseToken = buyToken
		quoteToken = sellToken
	} else {
		p.S1 = "SELL"
		baseVolume = payload.D
		quoteVolume = payload.D1
		excVolume = payload.F1
		baseToken = sellToken
		quoteToken = buyToken
	}
	if baseToken != nil && quoteToken != nil {
		p.S = baseToken.Symbol + quoteToken.Symbol
		if len(baseVolume) > 0 {
			if baseDec, err = decimal.NewFromString(baseVolume); err != nil {
				return
			}
		}
		p.Q = baseDec.DivRound(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)), 32).StringFixedBank(int32(conf.EffectiveDecimal))

		if len(quoteVolume) > 0 {
			if quoteDec, err = decimal.NewFromString(quoteVolume); err != nil {
				return
			}
		}
		p.Q1 = quoteDec.DivRound(pow10.Pow(decimal.NewFromInt32(quoteToken.Decimals)), 32).StringFixedBank(int32(conf.EffectiveDecimal))

		price := util.GetEffectivePriceRound(quoteDec.DivRound(baseDec, 32).Mul(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals-quoteToken.Decimals))), false, conf.OrderEffectiveDigitsGreaterThan10000, conf.OrderEffectiveDigitsLessThan10000)
		p.P = price.String()

		if len(excVolume) > 0 {
			if excDec, err = decimal.NewFromString(excVolume); err != nil {
				return
			}
		}
		p.Z = excDec.DivRound(pow10.Pow(decimal.NewFromInt32(baseToken.Decimals)), 32).StringFixedBank(int32(conf.EffectiveDecimal))
	}
	return
}

func ConvertOffChainFee(fees *model.OffChainFee) (offChainFee *binance.OffChainFee, err error) {
	if fees == nil {
		return
	}
	offChainFee = &binance.OffChainFee{}

	offChainFee.UpdateAccountGasFees, err = ConvertGasFee(fees.UpdateAccountGasFees)

	offChainFee.WithdrawalGasFees, err = ConvertGasFee(fees.WithdrawalGasFees)

	offChainFee.WithdrawalOtherGasFees, err = ConvertGasFee(fees.WithdrawalOtherGasFees)

	offChainFee.TransferGasFees, err = ConvertGasFee(fees.TransferGasFees)

	offChainFee.TransferNoIDGasFees, err = ConvertGasFee(fees.TransferNoIDGasFees)

	offChainFee.OrderGasFees, err = ConvertGasFee(fees.OrderGasFees)
	offChainFee.OrderGasMultiple = fees.OrderGasMultiple

	offChainFee.AddPairGasFees, err = ConvertGasFee(fees.AddPairGasFees)

	offChainFee.MiningGasFees, err = ConvertGasFee(fees.MiningGasFees)

	offChainFee.OnChainCancelOrderGasFees, err = ConvertGasFee(fees.OnChainCancelOrderGasFees)
	return
}

func ConvertGasFee(fees []*model.ShowTokenData) (gasFees []*binance.GasFee, err error) {
	if len(fees) == 0 {
		return
	}
	var (
		v decimal.Decimal
	)
	for _, f := range fees {
		if f == nil {
			continue
		}
		gasFee := &binance.GasFee{
			TokenId:  f.TokenID,
			Symbol:   f.Symbol,
			Decimals: f.Decimals,
			Volume:   f.Volume,
		}
		v, err = decimal.NewFromString(f.Volume)
		if err != nil {
			v = conf.Zero
		}
		gasFee.Quantity = v.DivRound(conf.Pow10.Pow(decimal.NewFromInt32(f.Decimals)), 32).String()
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

func ConvertTokenInfoToTokenData(info *model.TokenInfo) *model.ShowTokenData {
	data := &model.ShowTokenData{
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

func ConvertPairPrice(symbol string, t *model.PairsPricesRes) (ticker *binance.PairPrice) {
	if t == nil {
		return
	}
	ticker = &binance.PairPrice{
		Price:  t.Price,
		Symbol: symbol,
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
