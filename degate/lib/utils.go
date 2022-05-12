package lib

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/degatedev/degatesdk/degate/model"
	"github.com/degatedev/degatesdk/log"
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

func GetAmount(volume string, decimals int32, decimalDigits int) (amount string, err error) {
	q, err := decimal.NewFromString(volume)
	if err != nil {
		return
	}
	a := q.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals))), 32)
	if decimalDigits >= 0 {
		amount = a.Truncate(int32(decimalDigits)).StringFixed(int32(decimalDigits))
	} else {
		amount = a.String()
	}
	return
}

func FormatTime(t int64) string {
	if t == 0 {
		return ""
	}
	return time.Unix(t, 0).Format(timeFormat)
}

func GetSymbol(buyToken *model.ShowTokenData, sellToken *model.ShowTokenData, isBuy bool) (symbol string) {
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
