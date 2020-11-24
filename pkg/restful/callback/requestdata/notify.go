package requestdata

// 回调通知记录
type Notify struct {
	ResultCode                     string `xml:"result_code" json:"result_code"`                             // 错误代码，000000 成功,其他详细参见错误列表
	ResultMessage                  string `xml:"result_msg" json:"result_msg"`                               // 返回信息，返回错误原因
	OrganizationCode               string `xml:"ins_cd" json:"ins_cd"`                                       // 机构号，接入机构在富友的唯一代码
	FuYouMerchantCode              string `xml:"mchnt_cd" json:"mchnt_cd"`                                   // 商户号，富友分配的商户号
	TerminalId                     string `xml:"term_id" json:"term_id"`                                     // 终端号，富友分配的终端设备号
	RandomString                   string `xml:"random_str" json:"random_str"`                               // 随机字符串
	Signature                      string `xml:"sign" json:"sign"`                                           // 签名，详见签名生成算法
	UserId                         string `xml:"user_id" json:"user_id"`                                     // 用户在商户的ID
	OrderAmount                    int32  `xml:"order_amt" json:"order_amt"`                                 // 订单金额, 单位：分
	SettlementOrderAmount          int32  `xml:"settle_order_amt" json:"settle_order_amt"`                   // 应结订单金额，单位：分
	CurrencyType                   string `xml:"curr_type" json:"curr_type"`                                 // 货币种类
	FuYouTransactionId             string `xml:"transaction_id" json:"transaction_id"`                       // 富友渠道交易流水号
	FuYouMerchantOrderId           string `xml:"mchnt_order_no" json:"mchnt_order_no"`                       // 商户订单号, 商户系统内部的订单号（5 到30个字符、只能包含字母数字,区分大小写)）
	OrderType                      string `xml:"order_type" json:"order_type"`                               // 订单类型
	TransactionFinishTime          string `xml:"txn_fin_ts" json:"txn_fin_ts"`                               // 支付完成时间, 订单支付时间， 格式为yyyyMMddHHmmss
	ReservedSettlementDate         string `xml:"reserved_fy_settle_dt" json:"reserved_fy_settle_dt"`         // 富友清算日
	ReservedCouponFee              string `xml:"reserved_coupon_fee" json:"reserved_coupon_fee"`             // 结算优惠金额，单位：分
	ReservedBuyerLogonId           string `xml:"reserved_buyer_logon_id" json:"reserved_buyer_logon_id"`     // 买家在渠道登录账号
	ReservedPaymentChannelBillList string `xml:"reserved_fund_bill_list" json:"reserved_fund_bill_list"`     // 不定长支付宝交易资金渠道,详细渠道
	ReservedFuYouTracingId         string `xml:"reserved_fy_trace_no" json:"reserved_fy_trace_no"`           // 富友系统内部追踪号
	ReservedChannelOrderId         string `xml:"reserved_channel_order_id" json:"reserved_channel_order_id"` // 条码流水号，用户账单二维码对应的流水
	ReservedIsCreditCard           string `xml:"reserved_is_credit" json:"reserved_is_credit"`               // 表示信用卡或者花呗，0 表示其他(非信用方式) 不填，表示未知
	ReservedAdditions              string `xml:"reserved_addn_inf" json:"reserved_addn_inf"`                 // 附加数据
	ReservedSettlementAmount       string `xml:"reserved_settlement_amt" json:"reserved_settlement_amt"`     // 应结算订单金额，分为单位的整数只有成功交易才会返回如果使用了商户免充值优惠券，该值为订单金额-商户免充值如果没有使用商户免充值，该值等于订单金额
	ReservedBankType               string `xml:"reserved_bank_type" json:"reserved_bank_type"`               // 付款方式
	ReservedPromotionDetail        string `xml:"reserved_promotion_detail" json:"reserved_promotion_detail"` // 微信营销详情
}
