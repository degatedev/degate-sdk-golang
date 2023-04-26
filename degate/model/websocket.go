package model

import (
	"github.com/degatedev/degate-sdk-golang/degate/binance"
)

type SubscribeParam struct {
	Id int
}

type Client struct {
	Client string `json:"client"`
}

type StopParam struct {
	SubscribeParam
	Client string
}

type SubscribeKlineParam struct {
	SubscribeParam
	Symbol   string
	Interval string
}

type KlinePayload struct {
	binance.Payload
	E1 int       `json:"E"`
	B  int       `json:"B"` // baseTokenID
	U  int       `json:"U"` // quoteTokenID
	K  KlineData `json:"k"`
}

type KlineData struct {
	T  int    `json:"t"`
	T1 int    `json:"T"`
	B  int    `json:"B"`
	U  int    `json:"U"`
	I  int    `json:"i"`
	F  string `json:"f"`
	L  string `json:"L"`
	O  string `json:"o"`
	C  string `json:"c"`
	H  string `json:"h"`
	L1 string `json:"l"`
	V  string `json:"v"`
	N  int    `json:"n"`
	X  bool   `json:"x"`
	Q  string `json:"q"`
	V1 string `json:"V"`
	Q1 string `json:"Q"`
	U1 string `json:"u"`
}

type SubscribeTradeParam struct {
	SubscribeParam
	Symbol string
}

type TradePayload struct {
	binance.Payload
	E1 int    `json:"E"`
	B  int    `json:"B"`
	U  int    `json:"U"`
	T  string `json:"t"`
	P  string `json:"p"`
	Q  string `json:"q"`
	B1 string `json:"b"`
	A  string `json:"a"`
	T1 int    `json:"T"`
	M  bool   `json:"m"`
}

type SubscribeTickerParam struct {
	SubscribeParam
	Symbol string
}

type TickerPayload struct {
	binance.Payload
	E1 int    `json:"E"`
	B  int    `json:"B"`
	U  int    `json:"U"`
	P  string `json:"p"`
	P1 string `json:"P"`
	W  string `json:"w"`
	X  string `json:"x"`
	C  string `json:"c"`
	Q  string `json:"Q"`
	B1 string `json:"b"`
	I  string `json:"I"`
	A  string `json:"a"`
	A1 string `json:"A"`
	O  string `json:"o"`
	H  string `json:"h"`
	L  string `json:"l"`
	V  string `json:"v"`
	Q1 string `json:"q"`
	O1 int    `json:"O"`
	C1 int    `json:"C"`
	F  string `json:"F"`
	L1 string `json:"L"`
	N  int    `json:"n"`
	U1 string `json:"u"`
}

type SubscribeBookTickerParam struct {
	SubscribeParam
	Symbol string
}

type BookTickerPayload struct {
	binance.Payload
	B  int    `json:"B"`
	U  int    `json:"U"`
	B1 string `json:"b"`
	I  string `json:"I"`
	A  string `json:"a"`
	A1 string `json:"A"`
}

type SubscribeDepthParam struct {
	SubscribeParam
	Symbol string
	Level  int
	Speed  int
}

type DepthPayload struct {
	LastUpdateId int        `json:"last_update_id"`
	PairId       int        `json:"pair_id"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
	QuoteTokenId int        `json:"quote_token_id"`
	BaseTokenId  int        `json:"base_token_id"`
}

type SubscribeDepthUpdateParam struct {
	SubscribeParam
	Symbol string
	Speed  int
}

type DepthUpdatePayload struct {
	binance.Payload
	E1 int        `json:"E"`
	B  int        `json:"B"`
	Q  int        `json:"Q"`
	U  int        `json:"U"`
	U1 int        `json:"u"`
	P  int        `json:"p"`
	B1 [][]string `json:"b"`
	A  [][]string `json:"a"`
}

type OutboundAccountPositionPayload struct {
	binance.Payload
	Time int                     `json:"E"`
	U    int                     `json:"u"`
	B    *OutboundAccountBalance `json:"B"`
}

type OutboundAccountBalance struct {
	TokenId uint32 `json:"a"`
	F       string `json:"f"`
	L       string `json:"l"`
	O       string `json:"o"`
	D       string `json:"d"`
	W       string `json:"w"`
	Symbol  string `json:"s"`
}

type BalanceUpdatePayload struct {
	binance.Payload
	Time    int64  `json:"E"`
	TokenId uint32 `json:"a"`
	Balance string `json:"v"`
	Symbol  string `json:"s"`
}

type ExecutionReportPayload struct {
	binance.Payload
	Time             int    `json:"E"`
	OrderID          string `json:"i"`
	PairName         string `json:"s"`
	ClientOrderId    string `json:"c"`
	IsBuy            bool   `json:"S"`
	Ts               int64  `json:"T"`
	CreateTime       int64  `json:"O"`
	TradingFee       string `json:"n"`
	CancelReason     uint32 `json:"r"`
	Price            string `json:"p"`
	OrderType        uint16 `json:"o"`
	Status           string `json:"X"`
	LastDealVolume   string `json:"l"`
	IsOrderBook      bool   `json:"w"`
	TradingFeeSymbol string `json:"N"`
	OriginalVolume   string `json:"q"`
	QuoteOrderQty    string `json:"Q"`
	TotalDealVolume  string `json:"z"`
	TotalDealAmount  string `json:"Z"`
	DealPrice        string `json:"L"`
	TradeID          string `json:"t"`
	IsMaker          bool   `json:"m"`
	LastDealAmount   string `json:"Y"`
	WorkingTime      int64  `json:"W"`
	PairID           uint64 `json:"P"`
	BaseTokenID      uint32 `json:"B"`
	QuoteTokenID     uint32 `json:"U"`
	GseFee           string `json:"g"`
	MaxFee           string `json:"M"`
	GridId           uint32 `json:"G"`
	MaxFeeBips       int    `json:"I"`
	ValidUntil       int64  `json:"v"`
	CloseTime        int64  `json:"C"`
	UpdateTime       int64  `json:"d"`
	Completed        bool   `json:"h"`
	B1               bool   `json:"b"`
}

type SubscribeUserDataParam struct {
	ListenKey string `json:"listen_key"`
}
