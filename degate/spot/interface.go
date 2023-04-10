package spot

import (
	"github.com/degatedev/degate-sdk-golang/degate/binance"
	"github.com/degatedev/degate-sdk-golang/degate/model"
)

type SpotClient interface {
	Time() (response *binance.TimeResponse, err error)
	GasFee() (response *binance.GasFeeTokenResponse, err error)
	GetTradeFee(param *model.TradeFeeParam) (response *binance.TradeFeeResponse, err error)
	TokenList(param *model.TokenListParam) (response *model.TokensResponse, err error)
	ExchangeInfo() (response *binance.ExchangeInfoResponse, err error)
	Account() (response *binance.AccountResponse, err error)
	GetBalance(param *model.AccountBalanceParam) (response *binance.BalanceResponse, err error)
	Transfer(param *model.TransferParam) (response *binance.TransferResponse, err error)
	Withdraw(param *model.WithdrawParam) (response *binance.WithdrawResponse, err error)
	DepositHistory(param *model.DepositsParam) (response *binance.DepositHistoryResponse, err error)
	WithdrawHistory(param *model.WithdrawsParam) (response *binance.WithdrawHistoryResponse, err error)
	TransferHistory(param *model.TransfersParam) (response *binance.TransferHistoryResponse, err error)
	MyTrades(param *model.AccountTradesParam) (response *binance.TradeResponse, err error)
	GetOrder(param *model.OrderDetailParam) (res *model.OrderDetailResponse, response *binance.OrderResponse, err error)
	GetHistoryOrders(param *model.OrdersParam) (response *binance.OrdersResponse, err error)
	GetOpenOrders(param *model.OrdersParam) (response *binance.OrdersResponse, err error)
	NewOrder(param *model.OrderParam) (response *binance.NewOrderResponse, err error)
	CancelOrder(param *model.CancelOrderParam) (response *binance.OrderCancelResponse, err error)
	CancelOpenOrders(includeGrid bool) (response *binance.Response, err error)
	CancelOrderOnChain(param *model.CancelOrderParam) (response *binance.OrderCancelResponse, err error)
	Ping() (response *binance.PingResponse, err error)
	Depth(param *model.DepthParam) (response *binance.DepthResponse, err error)
	Trades(param *model.TradeLastedParam) (response *binance.TradeHistoryResponse, err error)
	TradesHistory(param *model.TradeHistoryParam) (response *binance.TradeHistoryResponse, err error)
	Klines(param *model.KlineParam) (response *binance.KlineResponse, err error)
	Ticker24(param *model.TickerParam) (response *binance.TickerResponse, err error)
	TickerPrice(param *model.PairPriceParam) (response *binance.PairPriceResponse, err error)
	BookTicker(param *model.BookTickerParam) (response *binance.BookTickerResponse, err error)
	NewListenKey() (response *binance.ListenKeyResponse, err error)
	ReNewListenKey(param *model.ListenKeyParam) (response *binance.EmptyResponse, err error)
	DeleteListenKey(param *model.ListenKeyParam) (response *binance.EmptyResponse, err error)
}
