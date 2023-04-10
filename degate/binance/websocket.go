package binance

type Payload struct {
	E string `json:"e"`
}

type Payload2 struct {
	Type string `json:"type"`
}

type ErrorPayload struct {
	Payload
	M string `json:"m"`
}

type KlinePayload struct {
	Payload
	E1 int       `json:"E"`
	S  string    `json:"s"`
	K  KlineData `json:"k"`
}

type KlineData struct {
	T  int    `json:"t"`
	T1 int    `json:"T"`
	S  string `json:"s"`
	I  string `json:"i"`
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
}

type TradePayload struct {
	Payload
	E1 int    `json:"E"`
	S  string `json:"s"`
	T  string `json:"t"`
	P  string `json:"p"`
	Q  string `json:"q"`
	B  string `json:"b"`
	A  string `json:"a"`
	T1 int    `json:"T"`
	M  bool   `json:"m"`
}

type TickerPayload struct {
	Payload
	E1 int    `json:"E"`
	S  string `json:"s"`
	P  string `json:"p"`
	P1 string `json:"P"`
	W  string `json:"w"`
	X  string `json:"x"`
	C  string `json:"c"`
	Q  string `json:"Q"`
	B  string `json:"b"`
	B1 string `json:"B"`
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
}

type BookTickerPayload struct {
	Payload
	S  string `json:"s"`
	B  string `json:"b"`
	B1 string `json:"B"`
	A  string `json:"a"`
	A1 string `json:"A"`
}

type DepthPayload struct {
	LastUpdateId int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type DepthUploadPayload struct {
	Payload
	E1 int        `json:"E"`
	S  string     `json:"s"`
	U  int        `json:"U"`
	U1 int        `json:"u"`
	B  [][]string `json:"b"`
	A  [][]string `json:"a"`
}

type OutboundAccountPositionPayload struct {
	Payload
	Time int                       `json:"E"`
	U    int                       `json:"u"`
	B    []*OutboundAccountBalance `json:"B"`
}

type OutboundAccountBalance struct {
	Symbol  string `json:"a"`
	F       string `json:"f"`
	L       string `json:"l"`
	TokenId uint32 `json:"i"`
}

type BalanceUpdatePayload struct {
	Payload
	Time    int64  `json:"E"`
	Symbol  string `json:"a"`
	Balance string `json:"v"`
	TokenId uint32 `json:"i"`
	T       int64  `json:"T"`
}

type ExecutionReportPayload struct {
	Payload
	Time             int    `json:"E"`
	PairName         string `json:"s"`
	ClientOrderId    string `json:"c"`
	Side             string `json:"S"`
	OrderType        string `json:"o"`
	OriginalVolume   string `json:"q"`
	Price            string `json:"p"`
	X                string `json:"x"`
	X1               string `json:"X"`
	CancelReason     string `json:"r"`
	OrderID          string `json:"i"`
	TotalDealVolume  string `json:"z"`
	TradingFee       string `json:"n"`
	TradingFeeSymbol string `json:"N"`
	Ts               int64  `json:"T"`
	TradeID          string `json:"t"`
	IsOrderBook      bool   `json:"w"`
	CreateTime       int64  `json:"O"`
	TotalDealAmount  string `json:"Z"`
	QuoteOrderQty    string `json:"Q"`
	LastDealVolume   string `json:"l"`
	DealPrice        string `json:"L"`
	IsMaker          bool   `json:"m"`
	LastDealAmount   string `json:"Y"`
	WorkingTime      int64  `json:"W"`
	PairID           uint64 `json:"P"`
	BaseTokenID      uint32 `json:"B"`
	QuoteTokenID     uint32 `json:"U"`
}
