package util

import (
	"github.com/degatedev/degate-sdk-golang/conf"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	one                    = decimal.NewFromInt(1)
	effectPointPrice       = decimal.NewFromInt(10000)
	effect1000PointPrice   = decimal.NewFromInt(1000)
	effectPointStablePrice = decimal.NewFromInt(1)
)

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

func GetEffectivePriceRound(p decimal.Decimal, isStable bool, effectiveDigitsGreaterBigDigits int, effectiveDigitsGreaterSmallDigits int, effectiveLess1000Digits int, isBuy bool) (r decimal.Decimal) {
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
		} else if p.LessThanOrEqual(effect1000PointPrice) {
			effectiveDigits = effectiveLess1000Digits
		} else {
			effectiveDigits = effectiveDigitsGreaterSmallDigits
		}
	}

	var roundPlaces int32
	if p.GreaterThanOrEqual(one) {
		intSize := len(p.BigInt().String())
		roundPlaces = int32(effectiveDigits - intSize)
	} else {
		e := 0
		for i, s := range p.String() {
			if x, _ := strconv.Atoi(string(s)); i > 1 && x > 0 {
				e = i - 2
				break
			}
		}
		roundPlaces = int32(e + effectiveDigits)
	}
	if isBuy {
		r = p.RoundFloor(roundPlaces)
	} else {
		r = p.RoundCeil(roundPlaces)
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
	return
}

func IsOrderEffectiveDigits(p decimal.Decimal, isStable bool, effectiveDigitsGreaterBigDigits int, effectiveDigitsGreaterSmallDigits int, effectiveLess1000Digits int) bool {
	if isStable {
		if p.GreaterThanOrEqual(effectPointStablePrice) {
			return IsEffectiveDigits(p.String(), effectiveDigitsGreaterBigDigits)
		} else {
			return IsEffectiveDigits(p.String(), effectiveDigitsGreaterSmallDigits)
		}
	} else {
		if p.GreaterThanOrEqual(effectPointPrice) {
			return IsEffectiveDigits(p.String(), effectiveDigitsGreaterBigDigits)
		} else if p.LessThanOrEqual(effect1000PointPrice) {
			return IsEffectiveDigits(p.String(), effectiveLess1000Digits)
		} else {
			return IsEffectiveDigits(p.String(), effectiveDigitsGreaterSmallDigits)
		}
	}
}
