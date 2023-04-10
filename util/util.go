package util

import (
	"encoding/binary"
	"github.com/degatedev/degate-sdk-golang/conf"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	one                    = decimal.NewFromInt(1)
	effectPointPrice       = decimal.NewFromInt(10000)
	effectPointStablePrice = decimal.NewFromInt(1)
)

func TruncateLastDecimal(d decimal.Decimal) decimal.Decimal {
	if d.Equal(decimal.NewFromBigInt(d.BigInt(), 0)) {
		return d
	}
	s := d.String()
	return d.Truncate(int32(len(s) - strings.Index(s, ".") - 2))
}

func ceil(d decimal.Decimal, places int32) decimal.Decimal {
	t := d.Truncate(places)
	if t.Equal(d) {
		return t
	}
	return t.Add(decimal.NewFromInt(1).DivRound(decimal.NewFromInt32(10).Pow(decimal.NewFromInt32(places)), places))
}

func IsEffectiveDigits(volume string, effectiveDigits int) (is bool) {
	volume = strings.Replace(volume, ".", "", -1)
	volume = strings.Replace(volume, "0", " ", -1)
	volume = strings.TrimSpace(volume)
	volume = strings.Replace(volume, " ", "0", -1)
	if len(volume) > effectiveDigits {
		return false
	}
	return true
}

// GetEffectiveVolume
// Example: usdc decimals=6
// v is 1 not 1*1e6
func GetEffectiveVolume(v decimal.Decimal, effectiveDigits int, effectiveDecimal int, isCeil bool) (volume decimal.Decimal) {
	vDeciaml := v.String()
	vBigInt := v.BigInt().String()
	if len(vBigInt) >= effectiveDigits {
		newInt := vBigInt[0:effectiveDigits]
		if isCeil {
			leftDecimal, _ := decimal.NewFromString(vDeciaml[effectiveDigits:])
			if leftDecimal.GreaterThan(decimal.NewFromInt(0)) {
				mNewInt, _ := decimal.NewFromString(newInt)
				mNewInt = mNewInt.Add(decimal.NewFromInt(1))
				newInt = mNewInt.String()
			}
		}
		for index := 0; index < len(vBigInt)-effectiveDigits; index++ {
			newInt = newInt + "0"
		}
		volume, _ = decimal.NewFromString(newInt)
		return
	}

	var places int32
	if base, _ := decimal.NewFromString("1"); v.GreaterThanOrEqual(base) {
		intSize := len(vBigInt)
		if effectiveDigits-intSize <= effectiveDecimal {
			places = int32(effectiveDigits - intSize)
			//volume = ceil(v, int32(effectiveDigits-intSize))
		} else {
			places = int32(effectiveDecimal)
			//volume = ceil(v, int32(effectiveDecimal))
		}
	} else {
		e := 0
		for i, s := range v.String() {
			if x, _ := strconv.Atoi(string(s)); i > 1 && x > 0 {
				e = i - 2
				break
			}
		}
		if e+effectiveDigits <= effectiveDecimal {
			places = int32(e + effectiveDigits)
			//volume = ceil(v, int32(e+effectiveDigits))
		} else {
			places = int32(effectiveDecimal)
			//volume = ceil(v, int32(effectiveDecimal))
		}
	}
	if isCeil {
		volume = ceil(v, places)
	} else {
		volume = v.Truncate(places)
	}
	return
}

// GetEffectiveDecimalVolume
// Example: usdc decimals=6
// v is 1 not 1*1e6
func GetEffectiveDecimalVolume(v decimal.Decimal, decimalsDigits int32, isCeil bool) (volume decimal.Decimal) {
	vBigInt, _ := decimal.NewFromString(v.BigInt().String())
	if v.Equal(vBigInt) {
		volume = v
		return
	}
	d := v.Sub(vBigInt)
	if isCeil {
		volume = ceil(d, decimalsDigits)
	} else {
		volume = d.Truncate(decimalsDigits)
	}
	volume = vBigInt.Add(volume)
	return
}

func GetEffectivePriceNew(p decimal.Decimal) (r decimal.Decimal) {
	var (
		ps     = p.String()
		price  string
		digits int
		ds     string
	)

	if p.GreaterThanOrEqual(effectPointPrice) {
		for _, s := range ps {
			ds = string(s)
			if ds == "." {
				break
			}
			if digits != conf.OrderEffectiveDigitsGreaterThan10000 {
				price = price + ds
				digits++
			} else {
				price = price + "0"
			}
		}
	} else {
		for index, s := range ps {
			ds = string(s)
			if index == 0 && ds == "0" {
				price = price + ds
				continue
			}
			if ds == "." {
				price = price + ds
				continue
			}
			price = price + ds
			if digits > 0 || ds != "0" {
				digits++
			}
			if digits == conf.OrderEffectiveDigitsLessThan10000 {
				break
			}
		}
	}
	r, _ = decimal.NewFromString(price)
	return
}

func GetEffectivePriceStable(p decimal.Decimal) (r decimal.Decimal) {
	var (
		ps       = p.String()
		price    string
		digits   int
		ds       string
		hasPoint bool
	)

	if p.GreaterThanOrEqual(effectPointStablePrice) {
		for _, s := range ps {
			ds = string(s)
			if ds == "." {
				price = price + "."
				hasPoint = true
			} else if digits != conf.OrderStabilityEffectiveDigitsGreaterThan1 {
				price = price + ds
				digits++
			} else if hasPoint {
				break
			} else {
				price = price + "0"
			}
		}
	} else {
		for index, s := range ps {
			ds = string(s)
			if index == 0 && ds == "0" {
				price = price + ds
				continue
			}
			if ds == "." {
				price = price + ds
				continue
			}
			price = price + ds
			if digits > 0 || ds != "0" {
				digits++
			}
			if digits == conf.OrderStabilityEffectiveDigitsLessThan1 {
				break
			}
		}
	}
	if strings.HasSuffix(price, ".") {
		price = price[0 : len(price)-1]
	}
	r, _ = decimal.NewFromString(price)
	return
}

func GetEffectivePriceRound(p decimal.Decimal, isStable bool, effectiveDigitsGreaterBigDigits int, effectiveDigitsGreaterSmallDigits int) (r decimal.Decimal) {
	var effectiveDigits int
	if isStable {
		if p.GreaterThanOrEqual(effectPointStablePrice) {
			effectiveDigits = effectiveDigitsGreaterBigDigits
		} else {
			effectiveDigits = effectiveDigitsGreaterSmallDigits
		}
	} else {
		if p.GreaterThanOrEqual(effectPointPrice) {
			effectiveDigits = effectiveDigitsGreaterBigDigits
		} else {
			effectiveDigits = effectiveDigitsGreaterSmallDigits
		}
	}

	if p.GreaterThanOrEqual(one) {
		intSize := len(p.BigInt().String())
		r = p.Round(int32(effectiveDigits - intSize))
	} else {
		e := 0
		for i, s := range p.String() {
			if x, _ := strconv.Atoi(string(s)); i > 1 && x > 0 {
				e = i - 2
				break
			}
		}
		r = p.Round(int32(e + effectiveDigits))
	}
	return
}

func CalculatePrice(quotaVolume1, baseVolume2 string, coefficient decimal.Decimal) (r decimal.Decimal) {
	v1, _ := decimal.NewFromString(quotaVolume1)
	v2, _ := decimal.NewFromString(baseVolume2)
	p := v1.DivRound(v2, 32).Mul(coefficient)
	if base, _ := decimal.NewFromString("1"); p.GreaterThanOrEqual(base) {
		intSize := len(p.BigInt().String())
		if intSize > 5 {
			r = p.Round(int32(6 - intSize))
		} else {
			r = p.Round(int32(5 - intSize))
		}
	} else {
		e := 0
		for i, s := range p.String() {
			if x, _ := strconv.Atoi(string(s)); i > 1 && x > 0 {
				e = i - 2
				break
			}
		}
		r = p.Round(int32(e + 5))
	}
	return
}

func CalculatePriceStable(quotaVolume1, baseVolume2 string, coefficient decimal.Decimal, isStable bool) (r decimal.Decimal) {
	v1, _ := decimal.NewFromString(quotaVolume1)
	v2, _ := decimal.NewFromString(baseVolume2)
	p := v1.DivRound(v2, 32).Mul(coefficient)
	if p.GreaterThanOrEqual(effectPointPrice) && !isStable {
		intSize := len(p.BigInt().String())
		r = p.Round(int32(6 - intSize))
	} else if p.GreaterThanOrEqual(effectPointStablePrice) {
		intSize := len(p.BigInt().String())
		r = p.Round(int32(5 - intSize))
	} else {
		e := 0
		for i, s := range p.String() {
			if x, _ := strconv.Atoi(string(s)); i > 1 && x > 0 {
				e = i - 2
				break
			}
		}
		if isStable {
			r = p.Round(int32(e + 4))
		} else {
			r = p.Round(int32(e + 5))
		}
	}
	return
}

func CalculateQuoteFeeVolume(v decimal.Decimal) (volume decimal.Decimal) {
	if !v.IsPositive() {
		return decimal.NewFromInt(0)
	}
	var result string
	vs := v.DivRound(conf.Pow10.Pow(decimal.NewFromInt(6)), 32).String()
	re := regexp.MustCompile("[1-9]")
	vs = re.ReplaceAllString(vs, "1")
	index := strings.Index(vs, "1")
	result = vs[0 : index+1]
	if index+1 < len(vs) {
		re = regexp.MustCompile("[1-9]")
		result += re.ReplaceAllString(vs[index+1:], "0")
	}
	volume, _ = decimal.NewFromString(result)
	volume = decimal.Max(volume, conf.MinFeeVolume)
	return
}

func GetStorageIdFormOrderId(id string) (storageId uint32) {
	i, b := new(big.Int).SetString(id, 10)
	if b {
		var a = make([]byte, 16)
		i.FillBytes(a)
		return binary.BigEndian.Uint32(a[12:16])
	}
	return
}
