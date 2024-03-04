package binance

type Depth struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type DepthResponse struct {
	Response
	Data *Depth `json:"data"`
}

type TradeHistoryResponse struct {
	Response
	Data []*TradeHistory `json:"data"`
}

type TradeHistory struct {
	Id           string `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

type KlineResponse struct {
	Response
	Data [][]interface{} `json:"data"`
}

type TickerResponse struct {
	Response
	Data *Ticker `json:"data"`
}

type Ticker struct {
	Symbol             string         `json:"symbol"`
	PriceChange        string         `json:"priceChange"`
	PriceChangePercent string         `json:"priceChangePercent"`
	WeightedAvgPrice   string         `json:"weightedAvgPrice"`
	PrevClosePrice     string         `json:"prevClosePrice"`
	LastPrice          string         `json:"lastPrice"`
	LastQty            string         `json:"lastQty"`
	BidPrice           string         `json:"bidPrice"`
	BidQty             string         `json:"bidQty"`
	AskPrice           string         `json:"askPrice"`
	AskQty             string         `json:"askQty"`
	OpenPrice          string         `json:"openPrice"`
	HighPrice          string         `json:"highPrice"`
	LowPrice           string         `json:"lowPrice"`
	Volume             string         `json:"volume"`
	QuoteVolume        string         `json:"quoteVolume"`
	OpenTime           int64          `json:"openTime"`
	CloseTime          int64          `json:"closeTime"`
	FirstId            string         `json:"firstId"`
	LastId             string         `json:"lastId"`
	Count              int            `json:"count"`
	MakerFee           string         `json:"makerFee"`
	TakerFee           string         `json:"takerFee"`
	PairId             uint64         `json:"pairId"`
	BaseToken          *ShowTokenData `json:"baseToken"`
	QuoteToken         *ShowTokenData `json:"quoteToken"`
}

type PairPriceResponse struct {
	Response
	Data []*PairPrice `json:"data"`
}

type PairPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	PairId uint64 `json:"pairId"`
}

type BookTickerResponse struct {
	Response
	Data *BookTicker `json:"data"`
}

type BookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}
