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
	F  int    `json:"F"`
	L1 int    `json:"L"`
	N  int    `json:"n"`
}

type BookTickerPayload struct {
	Payload
	//U  int    `json:"u"`
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
	E1 int                       `json:"E"`
	U  int                       `json:"u"`
	B  []*OutboundAccountBalance `json:"B"`
}

type OutboundAccountBalance struct {
	A string `json:"a"`
	F string `json:"f"`
	L string `json:"l"`
}

type BalanceUpdatePayload struct {
	Payload
	E1 int    `json:"E"`
	A  string `json:"a"`
	D  string `json:"d"`
	T  int    `json:"t"`
}

type ExecutionReportPayload struct {
	Payload
	E1 int    `json:"E"`
	S  string `json:"s"`
	C  string `json:"c"`
	S1 string `json:"S"`
	O  string `json:"o"`
	F  string `json:"f"`
	Q  string `json:"q"`
	P  string `json:"p"`
	P1 string `json:"P"`
	F1 string `json:"F"`
	G  int    `json:"g"`
	C1 string `json:"C"`
	X  string `json:"x"`
	X1 string `json:"X"`
	R  string `json:"r"`
	I  string `json:"i"`
	L  string `json:"l"`
	Z  string `json:"z"`
	L1 string `json:"L"`
	N  string `json:"n"`
	N1 string `json:"N"`
	T  int    `json:"T"`
	T1 int    `json:"t"`
	I1 int    `json:"I"`
	W  bool   `json:"w"`
	M  bool   `json:"m"`
	M1 bool   `json:"M"`
	O1 int    `json:"O"`
	Z1 string `json:"Z"`
	Y  string `json:"Y"`
	Q1 string `json:"Q"`
}
