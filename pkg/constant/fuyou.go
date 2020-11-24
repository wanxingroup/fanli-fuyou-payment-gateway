package constant

const (
	FuYouAPIVersion   = "1.0"
	FuYouTerminalId   = "88888888"
	FuYouSuccessCode  = "000000"
	FuYouSuccessMsg   = "SUCCESS"
	FuYouChannel      = "fuiou"
	FuYouCurrencyType = "CNY"

	// JSAPI-公众号支付、FWC--支付宝服务窗、LETPAY-小程序
	FuYouPayTypeJSAPI  = "JSAPI"
	FuYouPayTypeFWC    = "FWC"
	FuYouPayTypeLETPAY = "LETPAY"

	FuYouXMLHeader = `<?xml version="1.0" encoding="GBK" standalone="yes"?>`

	// 富友订单类型:ALIPAY (统一下单、条码支付、服务窗支付), WECHAT(统一下单、条码支付，公众号支付,小程序),UNIONPAY,BESTPAY(翼支付)
	FuYouOrderTypeAliPay   = "ALIPAY"
	FuYouOrderTypeWeChat   = "WECHAT"
	FuYouOrderTypeUnionPay = "UNIONPAY"
	FuYouOrderTypeBestPay  = "BESTPAY"

	FuYouErrorCodeTransactionNotFound = 1010
)
