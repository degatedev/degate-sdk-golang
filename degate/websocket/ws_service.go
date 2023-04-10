package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type PriceLevel [2]string

// WsDepthEvent define websocket depth book event
type WsDepthEvent struct {
	Event            string       `json:"e"`
	Time             int64        `json:"E"`
	TransactionTime  int64        `json:"T"`
	Symbol           string       `json:"s"`
	FirstUpdateID    int64        `json:"U"`
	LastUpdateID     int64        `json:"u"`
	SeqId            uint64       `json:"s"`
	PrevLastUpdateID int64        `json:"pu"`
	Bids             []PriceLevel `json:"b"`
	Asks             []PriceLevel `json:"a"`
}

//   "e": "trade",
//   "E": 123456789,
//   "B": 1,           // baseTokenID
//   "U": 0,           // quoteTokenID
//   "t": 12345,
//   "p": "0.001",
//   "q": "100",
//   "b": 88,
//   "a": 50,
//   "T": 123456785,
//   "m": true,

// WsAggTradeEvent define websocket aggTrade event
type WsAggTradeEvent struct {
	Event        string `json:"e"`
	Time         int64  `json:"E"`
	BaseTokenId  uint64 `json:"B"` // baseTokenID
	QuoteTokenId uint64 `json:"U"` // quoteTokenID
	TxId         string `json:"t"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	BuyOrderId   string `json:"b"`
	SellOrderId  string `json:"a"`
	TradeTime    uint64 `json:"T"`
	Maker        bool   `json:"m"`
}

// "e": "24hrTicker",
// "E": 123456789,
// "B": 1,             // BaseTokenID
// "U": 2,             // QuoteTokenID
// "p": "0.0015",
// "P": "250.00",
// "w": "0.0018",
// "x": "0.0009",
// "c": "0.0025",
// "Q": "10",
// "b": "0.0024",
// "I": "10",
// "a": "0.0026",
// "A": "100",
// "o": "0.0010",
// "h": "0.0025",
// "l": "0.0010",
// "v": "10000",
// "q": "18",
// "O": 0,
// "C": 86400000,
// "F": 0,
// "L": 18150,
// "n": 18151,
// "s": 180,
// "S": 1
type WsTickerEvent struct {
	Event           string `json:"e"`
	Time            int64  `json:"E"`
	BaseTokenId     uint64 `json:"B"`
	QuoteTokenId    uint64 `json:"U"` // quoteTokenID
	Price           string `json:"p"`
	Percent         string `json:"P"`
	AvgPrice        string `json:"w"`
	PrevPrice0      string `json:"x"`
	LastPrice       string `json:"c"`
	Quantity        string `json:"Q"`
	BestBidPrice    string `json:"b"`
	BestBidQuantity string `json:"I"`
	BestAskPrice    string `json:"a"`
	BestAskQuantity string `json:"A"`
	PrevPrice1      string `json:"o"`
	PriceHigh       string `json:"h"`
	PriceLow        string `json:"l"`
	Amount          string `json:"v"`
	Volume          string `json:"q"`
	StartTime       int64  `json:"O"`
	EndTime         int64  `json:"C"`
	FirstTradeId    string `json:"F"`
	LastTradeId     string `json:"L"`
	TradeCount      int64  `json:"n"`
	BasePrice       string `json:"s"`
	QuotePrice      string `json:"S"`
}

// "e": "kline",     // event
// "B": 1,           // baseTokenID
// "U": 0,           // quoteTokenID
// "k": {
//   "t": 123400000,
//   "T": 123460000,
//   "B": 1,         // baseTokenID
//   "U": 0,         // quoteTokenID
//   "i": 60,
//   "f": "100",
//   "L": "200",
//   "o": "0.0010",
//   "c": "0.0020",
//   "h": "0.0025",
//   "l": "0.0015",
//   "v": "1000",
//   "n": 100,
//   "x": false,
//   "q": "1.0000",
//   "V": "500",
//   "Q": "0.500",
//   "u": "111111"
// }

type WsKlineEvent struct {
	StartTime    int64  `json:"t"`
	EndTime      int64  `json:"T"`
	BaseTokenId  uint64 `json:"B"`
	QuoteTokenId uint64 `json:"U"` // quoteTokenID
	Interval     uint32 `json:"i"`
	FirstTradeId string `json:"f"`
	LastTradeId  string `json:"L"`
	OpenPrice    string `json:"o"`
	ClosePrice   string `json:"c"`
	HighPrice    string `json:"h"`
	LowPrice     string `json:"l"`
	Volume       string `json:"v"`
	Trades       uint32 `json:"n"`
	Complete     bool   `json:"x"`
	Amount       string `json:"q"`
	TakerVolume  string `json:"V"`
	TakerAmount  string `json:"Q"`
	TradeId      string `json:"u"`
}

const (
	WsEventTypeAccountUpdate = iota
	WsEventTypeBalanceUpdate
	WsEventTypeOrderUpdate
	WsEventTypeOrderTrade
)

var (
	// baseWsMainUrl       = "wss://dev-dg-ws.mykey007.com/ws"
	// baseCombinedMainURL = "wss://dev-dg-ws.mykey007.com/stream?streams="

	UserDataEventTypeOutboundAccountPosition UserDataEventType = "outboundAccountPosition"
	UserDataEventTypeBalanceUpdate           UserDataEventType = "balanceUpdate"
	UserDataEventTypeExecutionReport         UserDataEventType = "executionReport"
	UserDataEventTypeAccountTrade            UserDataEventType = "accountTrade"
	UserDataEventTypeListStatus              UserDataEventType = "ListStatus"
)

// WsDepthHandler handle websocket depth event
type WsDepthHandler func(event *WsDepthEvent)

// WsTradeHandler handle websocket trade event
type WsTradeHandler func(event *WsAggTradeEvent)

// WsTickerHandler handle websocket ticker event
type WsTickerHandler func(event *WsTickerEvent)

// WsKlineHandler handle websocket kline event
type WsKlineHandler func(event *WsKlineEvent)

// getWsEndpoint return the base endpoint of the WS according the UseTestnet flag
// func getWsEndpoint() string {
// 	// if UseTestnet {
// 	// 	return baseWsTestnetUrl
// 	// }
// 	return baseWsMainUrl
// }

// getCombinedEndpoint return the base endpoint of the combined stream according the UseTestnet flag
// func getCombinedEndpoint() string {
// 	return baseCombinedMainURL
// }

func getDepthEndpoint(baseTokenId, quoteTokenId uint64,
	baseURL string,
	levels string, rate int,
	inc, combined bool) (string, *SubscribeParam) {
	baseURL = strings.TrimSuffix(baseURL, "/")
	base := baseURL + "/ws"
	if combined {
		base = baseURL + "/stream?streams=" // getCombinedEndpoint()
	}
	ch := "depth" // "labelDepth0" //"depth"
	if inc {
		// 增量模式 没有levels
		levels = ""
	}

	s := ""
	if rate == 100 {
		s = fmt.Sprintf("%d.%d@%s%s_100ms", baseTokenId, quoteTokenId, ch, levels)
	} else {
		s = fmt.Sprintf("%d.%d@%s%s", baseTokenId, quoteTokenId, ch, levels)
	}
	param := SubscribeParam{
		Method: "SUBSCRIBE",
		Params: []string{s},
		ID:     uint64(time.Now().Unix()),
	}
	return base, &param
}

func getTradeEndpoint(endpoint string, baseId, quoteId uint64) (string, *SubscribeParam) {
	return endpoint, &SubscribeParam{
		Method: "SUBSCRIBE",
		Params: []string{fmt.Sprintf("%d.%d@trade", baseId, quoteId)},
		ID:     uint64(time.Now().Unix()),
	}
}

func getTickerEndpoint(endpoint string, baseId, quoteId uint64) (string, *SubscribeParam) {
	return endpoint, &SubscribeParam{
		Method: "SUBSCRIBE",
		Params: []string{fmt.Sprintf("%d.%d@ticker", baseId, quoteId)},
		ID:     uint64(time.Now().Unix()),
	}
}

func getKlineEndpoint(endpoint string, baseId, quoteId uint64, interval uint32) (string, *SubscribeParam) {
	return endpoint, &SubscribeParam{
		Method: "SUBSCRIBE",
		Params: []string{fmt.Sprintf("%d.%d@kline_%d", baseId, quoteId, interval)},
		ID:     uint64(time.Now().Unix()),
	}
}

func getLabelDepthEndpoint(baseTokenId, quoteTokenId uint64, base string, label string, combined bool) (string, *SubscribeParam) {
	// base := getWsEndpoint()
	// if combined {
	// 	base = getCombinedEndpoint()
	// }

	s := fmt.Sprintf("%d.%d@labelDepth%s", baseTokenId, quoteTokenId, label)
	param := SubscribeParam{
		Method: "SUBSCRIBE",
		Params: []string{s},
		ID:     uint64(time.Now().Unix()),
	}
	return base, &param
}

func wsPartialDepthServe(baseTokenId, quoteTokenId uint64, baseURL string, levels int, rate int, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	if levels != 5 && levels != 10 && levels != 20 {
		return nil, nil, errors.New("invalid levels")
	}
	levelsStr := fmt.Sprintf("%d", levels)
	endpoint, param := getDepthEndpoint(baseTokenId, quoteTokenId, baseURL, levelsStr, rate, false, false)
	return wsDepthServe(endpoint, param, handler, errHandler)
}

// WsPartialDepthServe serve websocket partial depth handler.
func WsPartialDepthServe(baseTokenId, quoteTokenId uint64, baseURL string, levels int, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return wsPartialDepthServe(baseTokenId, quoteTokenId, baseURL, levels, -1, handler, errHandler)
}

// WsPartialDepthServeWithRate serve websocket partial depth handler with rate.
func WsPartialDepthServeWithRate(baseTokenId, quoteTokenId uint64, baseURL string, levels int, rate int, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	return wsPartialDepthServe(baseTokenId, quoteTokenId, baseURL, levels, rate, handler, errHandler)
}

// WsDiffDepthServe serve websocket diff. depth handler.
func WsDiffDepthServe(baseTokenId, quoteTokenId uint64, baseURL string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, param := getDepthEndpoint(baseTokenId, quoteTokenId, baseURL, "", 100, true, false)
	return wsDepthServe(endpoint, param, handler, errHandler)
}

/*
// WsCombinedDepthServe is similar to WsPartialDepthServe, but it for multiple symbols
func WsCombinedDepthServe(symbolLevels map[string]string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for s, l := range symbolLevels {
		endpoint += fmt.Sprintf("%s@depth%s", strings.ToLower(s), l) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}
		event := new(WsDepthEvent)
		data := j.Get("data").MustMap()
		event.Event = data["e"].(string)
		event.Time, _ = data["E"].(json.Number).Int64()
		event.TransactionTime, _ = data["T"].(json.Number).Int64()
		event.Symbol = data["s"].(string)
		event.FirstUpdateID, _ = data["U"].(json.Number).Int64()
		event.LastUpdateID, _ = data["u"].(json.Number).Int64()
		event.PrevLastUpdateID, _ = data["pu"].(json.Number).Int64()
		bidsLen := len(data["b"].([]interface{}))
		event.Bids = make([]PriceLevel, bidsLen)
		for i := 0; i < bidsLen; i++ {
			item := data["b"].([]interface{})[i].([]interface{})
			event.Bids[i] = PriceLevel{item[0].(string), item[1].(string)}
			// 	Price:    item[0].(string),
			// 	Quantity: item[1].(string),
			// }
		}
		asksLen := len(data["a"].([]interface{}))
		event.Asks = make([]PriceLevel, asksLen)
		for i := 0; i < asksLen; i++ {
			item := data["a"].([]interface{})[i].([]interface{})
			event.Asks[i] = PriceLevel{item[0].(string), item[1].(string)}
			// 	Price:    item[0].(string),
			// 	Quantity: item[1].(string),
			// }
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedDiffDepthServe is similar to WsDiffDepthServe, but it for multiple symbols
func WsCombinedDiffDepthServe(symbols []string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@depth", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}
		event := new(WsDepthEvent)
		data := j.Get("data").MustMap()
		event.Event = data["e"].(string)
		event.Time, _ = data["E"].(json.Number).Int64()
		event.TransactionTime, _ = data["T"].(json.Number).Int64()
		event.Symbol = data["s"].(string)
		event.FirstUpdateID, _ = data["U"].(json.Number).Int64()
		event.LastUpdateID, _ = data["u"].(json.Number).Int64()
		event.PrevLastUpdateID, _ = data["pu"].(json.Number).Int64()
		bidsLen := len(data["b"].([]interface{}))
		event.Bids = make([]PriceLevel, bidsLen)
		for i := 0; i < bidsLen; i++ {
			item := data["b"].([]interface{})[i].([]interface{})
			event.Bids[i] = PriceLevel{item[0].(string), item[1].(string)}
			// 	Price:    item[0].(string),
			// 	Quantity: item[1].(string),
			// }
		}
		asksLen := len(data["a"].([]interface{}))
		event.Asks = make([]PriceLevel, asksLen)
		for i := 0; i < asksLen; i++ {
			item := data["a"].([]interface{})[i].([]interface{})
			event.Asks[i] = PriceLevel{item[0].(string), item[1].(string)}
			// 	Price:    item[0].(string),
			// 	Quantity: item[1].(string),
			// }
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
*/

// WsDiffDepthServeWithRate serve websocket diff. depth handler with rate.
func WsLabelDepthServe(baseTokenId, quoteTokenId uint64, baseURL string, label int, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, param := getLabelDepthEndpoint(baseTokenId, quoteTokenId, baseURL, fmt.Sprint(label), false)
	fmt.Println(endpoint, param)
	return wsDepthServe(endpoint, param, handler, errHandler)
}

// WsDiffDepthServeWithRate serve websocket diff. depth handler with rate.
func WsDiffDepthServeWithRate(baseTokenId, quoteTokenId uint64, baseURL string, rate int, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, param := getDepthEndpoint(baseTokenId, quoteTokenId, baseURL, "", rate, true, false)
	// fmt.Println(endpoint, param)
	return wsDepthServe(endpoint, param, handler, errHandler)
}

// WsAggTradeServe serve websocket diff. depth handler with rate.
func WsAggTradeServe(baseTokenId, quoteTokenId uint64, baseURL string, handler WsTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, param := getTradeEndpoint(baseURL, baseTokenId, quoteTokenId)
	// fmt.Println(endpoint, param)
	return wsAggTradeServe(endpoint, param, handler, errHandler)
}

func wsAggTradeServe(endpoint string, param *SubscribeParam, handler WsTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		// fmt.Println(string(message))
		// j, err := newJSON(message)
		// if err != nil {
		// 	errHandler(err)
		// 	return
		// }
		var event WsAggTradeEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		if event.Time == 0 {
			// 订阅成功消息
			return
		}

		handler(&event)
	}
	return wsServe(cfg, param, wsHandler, errHandler)
}

// WsAggTradeServe serve websocket diff. depth handler with rate.
func WsTickerServe(baseTokenId, quoteTokenId uint64, baseURL string, handler WsTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, param := getTickerEndpoint(baseURL, baseTokenId, quoteTokenId)
	// fmt.Println(endpoint, param)
	return wsTickerServe(endpoint, param, handler, errHandler)
}

func wsTickerServe(endpoint string, param *SubscribeParam, handler WsTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		// fmt.Println(string(message))
		// j, err := newJSON(message)
		// if err != nil {
		// 	errHandler(err)
		// 	return
		// }
		var event WsTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		if event.Time == 0 {
			// 订阅成功消息
			return
		}

		handler(&event)
	}
	return wsServe(cfg, param, wsHandler, errHandler)
}

// WsAggTradeServe serve websocket diff. depth handler with rate.
func WsKlineServe(baseTokenId, quoteTokenId uint64, baseURL string, interval uint32, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, param := getKlineEndpoint(baseURL, baseTokenId, quoteTokenId, interval)
	// fmt.Println(endpoint, param)
	return wsKlineServe(endpoint, param, handler, errHandler)
}

func wsKlineServe(endpoint string, param *SubscribeParam, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		// fmt.Println(string(message))
		// j, err := newJSON(message)
		// if err != nil {
		// 	errHandler(err)
		// 	return
		// }
		var event WsKlineEvent
		item := struct {
			Event string        `json:"e"`
			E     uint          `json:"E"`
			B     uint          `json:"B"`
			K     *WsKlineEvent `json:"k"`
		}{
			K: &event,
		}
		err := json.Unmarshal(message, &item)
		if err != nil {
			errHandler(err)
			return
		}
		if item.K.EndTime == 0 {
			// 订阅成功消息
			return
		}

		handler(&event)
	}
	return wsServe(cfg, param, wsHandler, errHandler)
}

func wsDepthServe(endpoint string, param *SubscribeParam, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	// var rateStr string
	// if rate != nil {
	// 	switch *rate {
	// 	case 100 * time.Millisecond:
	// 		rateStr = "_100ms"
	// 	default: // 1000
	// 		rateStr = ""
	// 	}
	// }
	// endpoint := fmt.Sprintf("%s/%d.%d@depth@%s%s", getWsEndpoint(), baseTokenId, quoteTokenId, levels, rateStr)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		// fmt.Println(string(message))
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}
		event := new(WsDepthEvent)
		event.Event = j.Get("e").MustString()
		event.Time = j.Get("E").MustInt64()
		if event.Event == "depth" {
			// depth 与 depthUpdate 数据格式不同
			event.LastUpdateID = j.Get("last_update_id").MustInt64()
			bidsLen := len(j.Get("bids").MustArray())
			event.Bids = make([]PriceLevel, bidsLen)
			for i := 0; i < bidsLen; i++ {
				item := j.Get("bids").GetIndex(i)
				event.Bids[i] = PriceLevel{item.GetIndex(0).MustString(), item.GetIndex(1).MustString()}
				// 	Price:    item.GetIndex(0).MustString(),
				// 	Quantity: item.GetIndex(1).MustString(),
				// }
			}
			asksLen := len(j.Get("asks").MustArray())
			event.Asks = make([]PriceLevel, asksLen)
			for i := 0; i < asksLen; i++ {
				item := j.Get("asks").GetIndex(i)
				event.Asks[i] = PriceLevel{item.GetIndex(0).MustString(), item.GetIndex(1).MustString()}
				// 	Price:    item.GetIndex(0).MustString(),
				// 	Quantity: item.GetIndex(1).MustString(),
				// }
			}
			handler(event)
			return
		}
		event.TransactionTime = j.Get("T").MustInt64()
		// event.Symbol = j.Get("s").MustString()
		event.FirstUpdateID = j.Get("U").MustInt64()
		event.LastUpdateID = j.Get("u").MustInt64()
		event.SeqId = j.Get("s").MustUint64()
		// event.PrevLastUpdateID = j.Get("pu").MustInt64()
		bidsLen := len(j.Get("b").MustArray())
		event.Bids = make([]PriceLevel, bidsLen)
		for i := 0; i < bidsLen; i++ {
			item := j.Get("b").GetIndex(i)
			event.Bids[i] = PriceLevel{item.GetIndex(0).MustString(), item.GetIndex(1).MustString()}
			// 	Price:    item.GetIndex(0).MustString(),
			// 	Quantity: item.GetIndex(1).MustString(),
			// }
		}
		asksLen := len(j.Get("a").MustArray())
		event.Asks = make([]PriceLevel, asksLen)
		for i := 0; i < asksLen; i++ {
			item := j.Get("a").GetIndex(i)
			event.Asks[i] = PriceLevel{item.GetIndex(0).MustString(), item.GetIndex(1).MustString()}
			// 	Price:    item.GetIndex(0).MustString(),
			// 	Quantity: item.GetIndex(1).MustString(),
			// }
		}
		handler(event)
	}
	return wsServe(cfg, param, wsHandler, errHandler)
}

type UserDataEventType string

// TimeInForceType define time in force type of order
type TimeInForceType string

// type WsUserDataEventHeader struct {
// 	Event string `json:"e"`
// 	Time  int64  `json:"E"`
// }

// WsUserDataEvent define user data event
type WsUserDataEvent struct {
	EventTyp      uint
	Event         UserDataEventType `json:"e"`
	Time          int64             `json:"E"`
	AccountUpdate WsAccountUpdate   `json:"B"`
	BalanceUpdate WsBalanceUpdate
	OrderUpdate   WsOrderUpdate
	OrderTrade    WsOrderTrade
}

// WsAccountUpdate define account update
// 	"e": "outboundAccountPosition",
// 	"E": 1564034571105,
// 	"u": 1564034571073,
// 	"B":
// 	  {
// 		"a": 1,                     // tokenID
// 		"f": "10000.000000",
// 		"l": "0.000000",
// 		"o": "1",
// 		"d": "1",
// 		"w": "1"
// 	  }
type WsAccountUpdate struct {
	TokenId        uint64 `json:"a"`
	Available      string `json:"f"`
	FrozenTotal    string `json:"l"`
	FrozenOrder    string `json:"o"`
	FrozenDeposit  string `json:"d"`
	FrozenWithdraw string `json:"w"`
}

// "e": "balanceUpdate",
// "E": 1638436619185,
// "a": 5,                 // TokenID
// "v": "138347.815785"
type WsBalanceUpdate struct {
	TokenId   uint64 `json:"a"`
	Available string `json:"v"`
}

// {
//     "e": "executionReport",
//     "E": 1637380246232,
//     "s": 339,                               // StorageID
//     "o": "980942284077759781128233812307",
//     "a": 12,
//     "p": 6,
//     "S": 4,
//     "B": 5,
//     "G": 0,                                 // gas fee token id
//     "F": 5,
//     "M": 0,
//     "f": "0",
//     "l": "0",
//     "c": "0",
//     "k": "0",
//     "u": "canceled",
//     "y": false,
//     "b": false,
//     "v": 1639972221,
//     "d": "1111",
//     "D": "1111",
//     "T": 0,
//     "t": 1637351446,
//     "r": 1637351422,
//     "U": 1637351446
// }
type WsOrderUpdate struct {
	StorageId     uint64 `json:"s"`
	OrderId       string `json:"o"`
	AccountId     uint64 `json:"a"`
	PairId        uint64 `json:"p"`
	SellTokenId   uint64 `json:"S"`
	BuyTokenId    uint64 `json:"B"`
	GasTokenId    uint64 `json:"G"`
	FeeTokenId    uint64 `json:"F"`
	MaxFeeRatio   uint64 `json:"M"`
	FilledSellVol string `json:"f"`
	GasUsed       string `json:"l"`
	FilledBuyVol  string `json:"c"`
	FeeUsed       string `json:"k"`
	Status        string `json:"u"`
	Y             bool   `json:"y"`
	BuyOrSell     bool   `json:"b"`
	EndTime       uint64 `json:"v"`
	SellTokenVol  string `json:"d"`
	BuyTokenVol   string `json:"D"`
	OrderType     uint   `json:"T"`
	CloseTime     uint64 `json:"t"`
	CreateTime    uint64 `json:"r"`
	UpdateTime    uint64 `json:"U"`
	CancelReason  uint32 `json:"C"`
	Price         string `json:"P"`
	MaxFee        string `json:"m"` // max fee
}

// "e": "accountTrade",
// "E": 123456789,
// "B": 1,
// "U": 0,
// "P": 10,
// "t": 12345,
// "p": "0.001",
// "q": "100",
// "b": 88,
// "y": 0,
// "c": "100",
// "o": 10,
// "k": true,
// "a": 50,
// "Y": 0,
// "C": "100",
// "O": 11,
// "K": false,
// "T": 123456785,
// "m": true,
// "M": true,        // is buy
// {"e":"accountTrade","E":1649729838364,"B":38,"U":2,"t":"53919893573506938846026056646645071265414976847805002301702219526183",
// "p":"2972.9","q":"0.567","Q":"1685.63","b":"268116561428386642823475130407","y":6,"c":"1685630000","o":2,"k":false,
// "a":"901941931658575567896478931478","Y":0,"C":"567000000000000000","O":38,"K":false,"T":1649729838212,"m":true,"M":false}
type WsOrderTrade struct {
	BaseTokenId  uint64 `json:"B"`
	QuoteTokenId uint64 `json:"U"`
	TradeId      string `json:"t"`
	Price        string `json:"p"`
	PairId       uint64 `json:"P"` // pair_id
	Quantity     string `json:"q"`
	QuoteQty     string `json:"Q"`
	BuyOrderId   string `json:"b"`
	SellOrderId  string `json:"a"`
	TradeTime    uint64 `json:"T"`
	Maker        bool   `json:"m"` // is this order maker
	IsBuy        bool   `json:"M"`
	GasFee       string `json:"f"`
	TradeFee     string `json:"r"`
}

type WsOCOOrder struct {
	Symbol        string `json:"s"`
	OrderId       int64  `json:"i"`
	ClientOrderId string `json:"c"`
}

// WsUserDataHandler handle WsUserDataEvent
type WsUserDataHandler func(event *WsUserDataEvent)

// WsUserDataServe serve user data handler with listen key
func WsUserDataServe(listenKey string, baseURL string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	baseURL = strings.TrimSuffix(baseURL, "/")
	endpoint := fmt.Sprintf("%s/ws/%s", baseURL, listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}

		// fmt.Printf("DeGate WsUserData msg: %v\n", string(message))

		// header := new(WsUserDataEventHeader)
		event := new(WsUserDataEvent)

		// err = json.Unmarshal(message, header)
		// if err != nil {
		// 	errHandler(err)
		// 	return
		// }

		switch UserDataEventType(j.Get("e").MustString()) {
		case UserDataEventTypeOutboundAccountPosition:
			err = json.Unmarshal(message, &event)
			if err != nil {
				errHandler(err)
				return
			}
			event.EventTyp = WsEventTypeAccountUpdate

		case UserDataEventTypeBalanceUpdate:
			err = json.Unmarshal(message, &event.BalanceUpdate)
			if err != nil {
				errHandler(err)
				return
			}
			event.EventTyp = WsEventTypeBalanceUpdate

		case UserDataEventTypeExecutionReport:
			err = json.Unmarshal(message, &event.OrderUpdate)
			if err != nil {
				fmt.Println(string(message))
				errHandler(err)
				return
			}
			// fmt.Println(string(message))
			event.EventTyp = WsEventTypeOrderUpdate

			// Unmarshal has case sensitive problem
			// event.TransactionTime = j.Get("T").MustInt64()
		// event.OrderUpdate.TransactionTime = j.Get("T").MustInt64()
		// event.OrderUpdate.Id = j.Get("i").MustInt64()
		// event.OrderUpdate.TradeId = j.Get("t").MustInt64()
		// event.OrderUpdate.FeeAsset = j.Get("N").MustString()
		// case UserDataEventTypeListStatus:
		// 	err = json.Unmarshal(message, &event.OCOUpdate)
		// 	if err != nil {
		// 		errHandler(err)
		// 		return
		// 	}
		case UserDataEventTypeAccountTrade:
			// fmt.Printf("trade msg: %v\n", string(message))
			err = json.Unmarshal(message, &event.OrderTrade)
			if err != nil {
				errHandler(err)
				return
			}
			event.EventTyp = WsEventTypeOrderTrade
		}

		handler(event)
	}
	return wsServe(cfg, nil, wsHandler, errHandler)
}
