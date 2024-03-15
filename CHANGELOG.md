# Changelog

## 1.7.0 - 2024-03-15
### Updated 
- Fix precision of quantity in newOrder response


## 1.6.9 - 2024-03-04
### Updated 
- Add minStepSize/maxSize in ShowTokenData struct 
- Return baseToken/quoteToken for ticker and pair endpoints

## 1.6.8 - 2024-02-19
### Updated 
- Parse deposit status 

## 1.6.7 - 2023-11-17
### Updated 
- Improve stardard withdrawal

## 1.6.6 - 2023-10-25
### Updated 
- RoundCeil for sell price, roundFloor for buy price 

## 1.6.5 - 2023-09-22
### Updated 
- Replace signature with access-token 

## 1.6.4 - 2023-07-12
### Updated 
- Replace minOrderPrice with minLimitOrderUSDValue in exchange info

## 1.6.3 - 2023-07-05
### Fixed
- Fix place order gas fee


## 1.6.2 - 2023-05-17
### Updated
- Update sign of cancel order request


## 1.6.1 - 2023-05-15
### Updated
- The price of place order is less than 1000, the effective digits are less than or equal to 4


## 1.6.0 - 2023-05-04
### Updated
Update sdk docs link
- https://api-docs.degate.com/spot


## 1.5.9 - 2023-05-04
### Updated
Update mainnet base endpoint
- https://mainnet-backend.degate.com
- wss://mainnet-ws.degate.com


## 1.5.8 - 2023-04-27
### Fixed
- Fix bugs


## 1.5.7 - 2023-04-26
### Fixed
- Fix bugs


## 1.5.6 - 2023-04-18
### Updated
Update testnet ws base endpoint
- wss://testnet-ws.degate.com


## 1.5.5 - 2023-04-04
### Fixed
- Fix bugs


## 1.5.4 - 2023-03-16
### Updated
- Update gas fee


## 1.5.3 - 2023-02-10
### Updated
- Update balanceUpdate Streams


## 1.5.2 - 2023-02-10
### Updated
- Update User Data Streams


## 1.5.1 - 2023-02-08
### Updated
- Update order sign


## 1.5.0 - 2023-01-18
### Updated
- Ticker24,TickerPrice response add pairId


## 1.4.9 - 2023-01-17
### Fixed
- Fixed new market order failed


## 1.4.6 - 2023-01-12
### Updated
- Ws accountTrade response add gasFee,tradeFee


## 1.4.5 - 2023-01-05
### Updated
- GetBalance response add tokenId


## 1.4.4 - 2022-12-20
### Updated
- DepositHistory,GetOpenOrders,GetHistoryOrders,MyTrades


## 1.4.3 - 2022-10-26
### Updated
- TradingKey changed to AssetPrivateKey


## 1.4.2 - 2022-10-18
### Fixed
- Fix Partial Book Depth


## 1.4.1 - 2022-09-02
### Updated
- Remove orderListId


## 1.4.0 - 2022-08-30
### Updated
- Update function name


## 1.3.2 - 2022-08-10
### Updated
- Withdraw,updateAccount max fee must greater than 0


## 1.3.1 - 2022-08-08
### Fixed
- Fix new buy order failure


## 1.3.0 - 2022-07-22
### Updated
- Update create storage_id


## 1.2.2 - 2022-07-12
### Updated
- AppPrivateKey changed to TradingKey


## 1.2.1 - 2022-07-01
### Updated
- Update README.md


## 1.2.0 - 2022-07-01
### Updated
- Update README.md


## 1.1.0 - 2022-06-29
### Added
- ExchangeInfo add rateLimits


## 1.0.9 - 2022-06-29
### Added
- GetTrades add fromId


## 1.0.8 - 2022-06-15
### Added
- Remove GetAllOrders, add GetHistoryOrders


## 1.0.7 - 2022-06-13
### Added
- Response add header


## 1.0.6 - 2022-06-09
### Added
- GasFee


## 1.0.5 - 2022-06-02
### Fixed
- limit default value


## 1.0.4 - 2022-06-01
### Fixed
- SubscribeTicker
- SubscribeUserData executionReport


## 1.0.3 - 2022-05-31
### Fixed
- GetTicker highPrice not return


## 1.0.2 - 2022-05-27
### Fixed
- GetDeposits status remove 0
- GetWithdraws status remove 2


## 1.0.1 - 2022-05-25
### Fixed
- GetAccount function return error when AccountId is 0


## 1.0.0 - 2022-05-12
### Added
- First release, please find details from README.md