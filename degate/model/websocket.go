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
	B  int    `json:"B"` //baseTokenID
	U  int    `json:"U"` //quoteTokenID
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
	F  int    `json:"F"`
	L1 int    `json:"L"`
	N  int    `json:"n"`
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
	P  [][]string `json:"p"`
	B1 [][]string `json:"b"`
	A  [][]string `json:"a"`
}

type OutboundAccountPositionPayload struct {
	binance.Payload
	E1 int                     `json:"E"`
	U  int                     `json:"u"`
	B  *OutboundAccountBalance `json:"B"`
}

type OutboundAccountBalance struct {
	A int    `json:"a"`
	F string `json:"f"`
	L string `json:"l"`
	O string `json:"o"`
	D string `json:"d"`
	W string `json:"w"`
}

type BalanceUpdatePayload struct {
	binance.Payload
	E1 int    `json:"E"`
	A  int    `json:"a"`
	V  string `json:"v"`
}

type ExecutionReportPayload struct {
	binance.Payload
	E1 int    `json:"E"`
	S  int    `json:"s"`
	O  string `json:"o"`
	A  int    `json:"a"`
	P  int    `json:"p"`
	S1 int    `json:"S"`
	B  int    `json:"B"`
	G  int    `json:"G"`
	F  int    `json:"F"`
	M  int    `json:"M"`
	F1 string `json:"f"`
	L  string `json:"l"`
	C  string `json:"c"`
	K  string `json:"k"`
	U  string `json:"u"`
	Y  bool   `json:"y"`
	B1 bool   `json:"b"`
	V  int    `json:"v"`
	D  string `json:"d"`
	D1 string `json:"D"`
	T  int    `json:"T"`
	T1 int    `json:"t"`
	R  int    `json:"r"`
	U1 int    `json:"U"`
}

type SubscribeUserDataParam struct {
	ListenKey string `json:"listen_key"`
}
