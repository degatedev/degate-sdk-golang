package lib

import (
	"encoding/json"
	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"math/big"
	"strings"
	"time"

	"github.com/degatedev/degate-sdk-golang/log"
	"github.com/shopspring/decimal"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func String(o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		log.Error("String err: %v", err)
		return ""
	}
	return string(b)
}

func GetAmountNew(volume string, decimals int32) (amount string) {
	q, err := decimal.NewFromString(volume)
	if err != nil {
		log.Error("GetAmountNew error:%v", err)
		return
	}
	amount = q.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals))), 32).String()
	return
}

func FormatTime(t int64) string {
	if t == 0 {
		return ""
	}
	return time.Unix(t, 0).Format(timeFormat)
}

func GetSymbol(buyToken *binance.ShowTokenData, sellToken *binance.ShowTokenData, isBuy bool) (symbol string) {
	if buyToken == nil || sellToken == nil {
		return
	}
	if isBuy {
		symbol = strings.ToUpper(buyToken.Symbol) + strings.ToUpper(sellToken.Symbol)
	} else {
		symbol = strings.ToUpper(sellToken.Symbol) + strings.ToUpper(buyToken.Symbol)
	}
	return
}

func ParseParam(p string, obj interface{}) (err error) {
	err = json.Unmarshal([]byte(p), obj)
	if err != nil {
		return
	}
	return
}

func GenerateOrderId(accountId uint32, tokenId uint64, storageId uint64) (orderId string) {
	var (
		x, y, z    big.Int
		i, e, f, g = big.NewInt(2), big.NewInt(96), big.NewInt(64), big.NewInt(32)
	)
	ts := time.Now().Unix()
	x.Mul(big.NewInt(int64(accountId)), x.Exp(i, e, nil))
	y.Mul(big.NewInt(ts), y.Exp(i, f, nil))
	z.Mul(big.NewInt(int64(tokenId)), z.Exp(i, g, nil))
	x.Add(&x, &y).Add(&x, &z).Add(&x, big.NewInt(int64(storageId)))
	return x.String()
}
