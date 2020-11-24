package parameters

type PrePayRequest struct {
	Version                string `xml:"version" json:"version"`
	OrganizationCode       string `xml:"ins_cd" json:"ins_cd"`
	FuYouMerchantCode      string `xml:"mchnt_cd" json:"mchnt_cd"`
	TerminalId             string `xml:"term_id" json:"term_id"`
	RandomString           string `xml:"random_str" json:"random_str"`
	GoodsDescription       string `xml:"goods_des" json:"goods_des"`
	GoodsDetail            string `xml:"goods_detail" json:"goods_detail"`
	GoodsTag               string `xml:"goods_tag" json:"goods_tag"`
	ProductId              string `xml:"product_id" json:"product_id"`
	Additions              string `xml:"addn_inf" json:"addn_inf"`
	ChannelMerchantOrderId string `xml:"mchnt_order_no" json:"mchnt_order_no"`
	CurrencyType           string `xml:"curr_type" json:"curr_type"`
	OrderAmount            string `xml:"order_amt" json:"order_amt"`
	UserIP                 string `xml:"term_ip" json:"term_ip"`
	TransactionBeginTime   string `xml:"txn_begin_ts" json:"txn_begin_ts"`
	NotifyURL              string `xml:"notify_url" json:"notify_url"`
	LimitPay               string `xml:"limit_pay" json:"limit_pay"`
	PaymentType            string `xml:"trade_type" json:"trade_type"`
	OpenId                 string `xml:"openid" json:"openid"` // Deprecated
	SubOpenId              string `xml:"sub_openid" json:"sub_openid"`
	SubAppId               string `xml:"sub_appid" json:"sub_appid"`
	ReservedExpireMinute   int    `xml:"reserved_expire_minute" json:"reserved_expire_minute"`
}
