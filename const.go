package indodax

type Method string

const (
	MethodGetInfo       Method = "getInfo"
	MethodTransHistory  Method = "transHistory"
	MethodTrade         Method = "trade"
	MethodTradeHistory  Method = "tradeHistory"
	MethodOpenOrders    Method = "openOrders"
	MethodOrderHistory  Method = "orderHistory"
	MethodGetOrder      Method = "getOrder"
	MethodCancelOrder   Method = "cancelOrder"
	MethodWithdrawCoin  Method = "withdrawCoin"
	MethodListDownline  Method = "listDownline"
	MethodCheckDownline Method = "checkDownline"
)
