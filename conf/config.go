package conf

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"

	"github.com/degatedev/degate-sdk-golang/degate/model"
)

var (
	Debug                 = false
	BaseUrl               = "https://backend-testnet.degate.com/"
	WebsocketBaseUrl      = "wss://ws-testnet.degate.com/"
	OrderBookPath         = "order-book-api/"
	WsPath                = "order-book-ws-api/"
	EffectiveDigits       = 6
	GasFeeEffectiveDigits = 3
	Timeout               = 60
	ValidUntil            = int64(60 * 60 * 24 * 30)
	DepthLimit            = int64(100)
	MaxRetryNum           = 11
	PongInterval          = 5
	ConnectInterval       = 2 //2 4 6 8 10 12 14 16 18 20 22
	Conf                  = new(Config).Init()
	KlineInterval         = map[string]int{
		"1m":  60,
		"3m":  180,
		"5m":  300,
		"15m": 900,
		"30m": 1800,
		"1h":  3600,
		"2h":  7200,
		"4h":  14400,
		"6h":  21600,
		"12h": 43200,
		"1d":  86400,
		"1w":  604800,
		"1M":  2678400,
		"3M":  8035200,
		"6M":  16070400,
		"1y":  31536000,
	}
	OrderSource                               = "goSdk"
	OrderMaxVolume                            = decimal.NewFromInt(0)
	Pow10                                     = decimal.NewFromInt(10)
	Zero                                      = decimal.NewFromInt(0)
	MinFeeVolume                              = decimal.NewFromFloat(0.3)
	MarketOrderBuyAdjustment                  = decimal.NewFromFloat(1.1)
	MarketOrderSellAdjustment                 = decimal.NewFromFloat(0.9)
	GasFeeSymbol                              = "ETH"
	MinGas                                    = "10"
	OrderEffectiveDigitsGreaterThan10000      = 6
	OrderEffectiveDigitsLessThan10000         = 5
	OrderStabilityEffectiveDigitsGreaterThan1 = 5
	OrderStabilityEffectiveDigitsLessThan1    = 4
	OrderPriceDiff                            = decimal.NewFromInt(1).DivRound(decimal.NewFromInt(10000), 32)
	TransferEffectiveDigits                   = 7
)

func GetInterval(interval int) string {
	for k, v := range KlineInterval {
		if v == interval {
			return k
		}
	}
	return ""
}

//export AppConfig
type AppConfig struct {
	Debug            bool
	BaseUrl          string
	WebsocketBaseUrl string
	ExchangeAddress  string
	ChainId          int64
	Timeout          int
	AccountId        uint32
	AccountAddress   string
	TradingKey       string
	EffectDigits     uint
	EffectDecimals   uint
	Tokens           []*model.TokenInfo
	ShowHeader       bool
}

//export Config
type Config struct {
	Tokens      map[string]*model.TokenInfo
	QuoteTokens map[string]*model.TokenInfo
}

func (c *Config) Init() *Config {
	c.Tokens = map[string]*model.TokenInfo{}
	c.QuoteTokens = map[string]*model.TokenInfo{}
	return c
}

func (c *Config) AddTokens(tokens []*model.TokenInfo) {
	if tokens != nil {
		for _, t := range tokens {
			c.AddToken(t)
		}
	}
}

func (c *Config) AddToken(token *model.TokenInfo) {
	if token == nil {
		return
	}
	c.Tokens[strings.ToUpper(token.Symbol)] = token
	if token.IsQuotableToken {
		c.QuoteTokens[strings.ToUpper(token.Symbol)] = token
	}
}

func (c *Config) GetTokenInfo(symbol string) *model.TokenInfo {
	if len(symbol) == 0 {
		return nil
	}
	return c.Tokens[strings.ToUpper(symbol)]
}

func (c *Config) GetTokenInfoById(id int) *model.TokenInfo {
	if id < 0 {
		return nil
	}
	for _, t := range c.Tokens {
		if t.Id == id {
			return t
		}
	}
	return nil
}

func (c *Config) GetTokens(symbol string) (baseToken *model.TokenInfo, quoteToken *model.TokenInfo) {
	if len(symbol) == 0 {
		return
	}
	symbol = strings.ToUpper(symbol)
	for s, t := range c.QuoteTokens {
		if strings.HasSuffix(symbol, s) {
			quoteToken = t
			baseSymbol := symbol[0 : len(symbol)-len(s)]
			if len(baseSymbol) > 0 {
				baseToken = c.Tokens[baseSymbol]
			}
			break
		}
	}
	if quoteToken == nil {
		for s, t := range c.Tokens {
			if strings.HasSuffix(symbol, s) {
				quoteToken = t
				baseSymbol := symbol[0 : len(symbol)-len(s)]
				if len(baseSymbol) > 0 {
					baseToken = c.Tokens[baseSymbol]
				}
				break
			}
		}
	}
	return
}

func (c *Config) GetTokenIds(tokens string) (ids string, err error) {
	var token *model.TokenInfo
	if len(tokens) == 0 {
		return
	}
	s := strings.Split(tokens, ",")
	for _, v := range s {
		if len(v) > 0 {
			token = c.Tokens[strings.ToUpper(v)]
			if token == nil {
				err = fmt.Errorf("not config token %v", v)
				return
			} else {
				ids += strconv.FormatUint(uint64(token.Id), 10) + ","
			}
		}
	}
	if len(ids) > 0 {
		ids = ids[0 : len(ids)-1]
	}
	return
}
