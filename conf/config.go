package conf

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/degatedev/degate-sdk-golang/degate/model"
)

var (
	BaseUrl               = "https://v1-mainnet-backend.degate.com"
	WebsocketBaseUrl      = "wss://v1-mainnet-ws.degate.com"
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
	OrderMaxVolume                            = decimal.NewFromInt(10).Pow(decimal.NewFromInt(28))
	Pow10                                     = decimal.NewFromInt(10)
	Zero                                      = decimal.NewFromInt(0)
	MarketOrderBuyAdjustment                  = decimal.NewFromFloat(1.1)
	MarketOrderSellAdjustment                 = decimal.NewFromFloat(0.9)
	GasFeeSymbol                              = "ETH"
	MinGas                                    = "10"
	OrderEffectiveDigitsGreaterThan10000      = 6
	OrderEffectiveDigitsLessThan10000         = 5
	OrderEffectiveDigitsLessThan1000          = 4
	OrderStabilityEffectiveDigitsGreaterThan1 = 5
	OrderStabilityEffectiveDigitsLessThan1    = 4
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
	Debug                 bool
	BaseUrl               string
	WebsocketBaseUrl      string
	ExchangeAddress       string
	ChainId               int64
	Timeout               int
	AccountId             uint32
	AccountAddress        string
	AssetPrivateKey       string
	EffectDigits          uint
	EffectDecimals        uint
	Tokens                []*model.TokenInfo
	ShowHeader            bool
	UseTradeKey           int64
	AccessToken           string
	AccessTokenExpireTime int64
}

//export Config
type Config struct {
	Tokens      sync.Map //map[string]*model.TokenInfo
	QuoteTokens sync.Map //map[string]*model.TokenInfo
}

func (c *Config) Init() *Config {
	c.Tokens = sync.Map{}      //map[string]*model.TokenInfo{}
	c.QuoteTokens = sync.Map{} //map[string]*model.TokenInfo{}
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
	c.Tokens.Store(strings.ToUpper(token.Symbol), token)
	if token.IsQuotableToken {
		c.QuoteTokens.Store(strings.ToUpper(token.Symbol), token)
	}
}

func (c *Config) GetTokenInfo(symbol string) *model.TokenInfo {
	if len(symbol) == 0 {
		return nil
	}
	tokenObj, ok := c.Tokens.Load(strings.ToUpper(symbol))
	if !ok {
		return nil
	}
	return tokenObj.(*model.TokenInfo)
}

func (c *Config) GetTokenInfoById(id int) (t *model.TokenInfo) {
	if id < 0 {
		return
	}
	c.Tokens.Range(func(key, value any) bool {
		token := value.(*model.TokenInfo)
		if token != nil && token.Id == id {
			t = token
			return false
		}
		return true
	})
	return
}

func (c *Config) GetTokens(symbol string) (baseToken *model.TokenInfo, quoteToken *model.TokenInfo) {
	if len(symbol) == 0 {
		return
	}
	symbol = strings.ToUpper(symbol)
	c.QuoteTokens.Range(func(key, value any) bool {
		s := key.(string)
		t := value.(*model.TokenInfo)
		if t != nil && strings.HasSuffix(symbol, s) {
			quoteToken = t
			baseSymbol := symbol[0 : len(symbol)-len(s)]
			if len(baseSymbol) > 0 {
				baseToken = c.GetTokenInfo(baseSymbol)
			}
			return false
		}
		return true
	})
	if quoteToken == nil {
		c.Tokens.Range(func(key, value any) bool {
			s := key.(string)
			t := value.(*model.TokenInfo)
			if t != nil && strings.HasSuffix(symbol, s) {
				quoteToken = t
				baseSymbol := symbol[0 : len(symbol)-len(s)]
				if len(baseSymbol) > 0 {
					baseToken = c.GetTokenInfo(baseSymbol)
				}
				return false
			}
			return true
		})
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
			token = c.GetTokenInfo(strings.ToUpper(v))
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
