package util

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

var (
	OrderStabilityEffectiveDigitsGreaterThan1 = 5
	OrderStabilityEffectiveDigitsLessThan1    = 4
	OrderEffectiveDigitsLessThan1000          = 4

	OrderEffectiveDigitsGreaterThan10000 = 6
	OrderEffectiveDigitsLessThan10000    = 5
)

// 交易对类型为普通时：价格默认4位，如果数字>=10000时，可以支持6位;如果数字>=1000并且<10000时，可以支持5位;
// 交易对类型为稳定类：价格默认4位，如果数字>=1，可以支持5位
func TestGetEffectivePriceRound(t *testing.T) {
	type args struct {
		prices                            []decimal.Decimal
		isStable                          bool
		effectiveDigitsGreaterBigDigits   int
		effectiveDigitsGreaterSmallDigits int
		effectiveLess1000Digits           int
		isBuy                             bool
	}
	tests := []struct {
		name                string
		args                args
		wantConvertedPrices []decimal.Decimal
	}{
		{
			name: "isStable = true, isBuy = true, roundFloor",
			args: args{
				prices: []decimal.Decimal{
					decimal.NewFromFloat(0.1234567),
					decimal.NewFromFloat(1),
					decimal.NewFromFloat(1.1234567),
					decimal.NewFromFloat(10.191),
					decimal.NewFromFloat(1600.33),
				},
				isStable: true,
				isBuy:    true,
			},
			wantConvertedPrices: []decimal.Decimal{
				decimal.NewFromFloat(0.1234),
				decimal.NewFromFloat(1),
				decimal.NewFromFloat(1.1234),
				decimal.NewFromFloat(10.191),
				decimal.NewFromFloat(1600.3),
			},
		},
		{
			name: "isStable = true, isBuy = false, roundCeil",
			args: args{
				prices: []decimal.Decimal{
					decimal.NewFromFloat(0.1234567),
					decimal.NewFromFloat(1),
					decimal.NewFromFloat(1.1234567),
					decimal.NewFromFloat(10.191),
					decimal.NewFromFloat(1600.33),
				},
				isStable: true,
				isBuy:    false,
			},
			wantConvertedPrices: []decimal.Decimal{
				decimal.NewFromFloat(0.1235),
				decimal.NewFromFloat(1),
				decimal.NewFromFloat(1.1235),
				decimal.NewFromFloat(10.191),
				decimal.NewFromFloat(1600.4),
			},
		},
		{
			name: "isStable = false, isBuy = true, roundFloor",
			args: args{
				prices: []decimal.Decimal{
					decimal.NewFromFloat(30000.12),
					decimal.NewFromFloat(10000.1234567),
					decimal.NewFromFloat(10000),
					decimal.NewFromFloat(9999.1234567),
					decimal.NewFromFloat(9999.1234567),
					decimal.NewFromFloat(1600.33),
					decimal.NewFromFloat(1000),
					decimal.NewFromFloat(999.123456),
					decimal.NewFromFloat(10.191),
					decimal.NewFromFloat(10.195),
					decimal.NewFromFloat(1.234567),
					decimal.NewFromFloat(0.1234567),
				},
				isStable: false,
				isBuy:    true,
			},
			wantConvertedPrices: []decimal.Decimal{
				decimal.NewFromFloat(30000.1),
				decimal.NewFromFloat(10000.1),
				decimal.NewFromFloat(10000),
				decimal.NewFromFloat(9999.1),
				decimal.NewFromFloat(9999.1),
				decimal.NewFromFloat(1600.3),
				decimal.NewFromFloat(1000),
				decimal.NewFromFloat(999.1),
				decimal.NewFromFloat(10.19),
				decimal.NewFromFloat(10.19),
				decimal.NewFromFloat(1.234),
				decimal.NewFromFloat(0.1234),
			},
		},
		{
			name: "isStable = false, isBuy = false, roundCeil",
			args: args{
				prices: []decimal.Decimal{
					decimal.NewFromFloat(30000.12),
					decimal.NewFromFloat(10000.1234567),
					decimal.NewFromFloat(10000),
					decimal.NewFromFloat(9999.1234567),
					decimal.NewFromFloat(9999.1234567),
					decimal.NewFromFloat(1600.33),
					decimal.NewFromFloat(1000),
					decimal.NewFromFloat(999.123456),
					decimal.NewFromFloat(10.191),
					decimal.NewFromFloat(10.195),
					decimal.NewFromFloat(1.234567),
					decimal.NewFromFloat(0.1234567),
				},
				isStable: false,
				isBuy:    false,
			},
			wantConvertedPrices: []decimal.Decimal{
				decimal.NewFromFloat(30000.2),
				decimal.NewFromFloat(10000.2),
				decimal.NewFromFloat(10000),
				decimal.NewFromFloat(9999.2),
				decimal.NewFromFloat(9999.2),
				decimal.NewFromFloat(1600.4),
				decimal.NewFromFloat(1000),
				decimal.NewFromFloat(999.2),
				decimal.NewFromFloat(10.2),
				decimal.NewFromFloat(10.2),
				decimal.NewFromFloat(1.235),
				decimal.NewFromFloat(0.1235),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 交易对类型为普通时：价格默认4位，如果数字>=10000时，可以支持6位;如果数字>=1000并且<10000时，可以支持5位;
			// 交易对类型为稳定类：价格默认4位，如果数字>=1，可以支持5位
			if tt.args.isStable {
				tt.args.effectiveDigitsGreaterBigDigits = OrderStabilityEffectiveDigitsGreaterThan1
				tt.args.effectiveDigitsGreaterSmallDigits = OrderStabilityEffectiveDigitsLessThan1
				tt.args.effectiveLess1000Digits = OrderEffectiveDigitsLessThan1000
			} else {
				tt.args.effectiveDigitsGreaterBigDigits = OrderEffectiveDigitsGreaterThan10000
				tt.args.effectiveDigitsGreaterSmallDigits = OrderEffectiveDigitsLessThan10000
				tt.args.effectiveLess1000Digits = OrderEffectiveDigitsLessThan1000
			}
			if len(tt.args.prices) != len(tt.wantConvertedPrices) {
				t.Errorf("len(tt.args.prices) != len(tt.wantConvertedPrices)")
			}
			for idx, price := range tt.args.prices {
				roundTo := "Ceil"
				if tt.args.isBuy {
					roundTo = "Floor"
				}
				tt.name = fmt.Sprintf("isStable=%t, isBuy=%t, shouldRoundTo:%s", tt.args.isStable, tt.args.isBuy, roundTo)

				if gotR := GetEffectivePriceRound(price, tt.args.isStable, tt.args.effectiveDigitsGreaterBigDigits, tt.args.effectiveDigitsGreaterSmallDigits, tt.args.effectiveLess1000Digits, tt.args.isBuy); !gotR.Equal(tt.wantConvertedPrices[idx]) {
					t.Errorf("GetEffectivePriceRound() = %v, want %v", gotR, tt.wantConvertedPrices[idx])
				}
			}
		})
	}
}
