package parameters

type QueryPayResultRequest struct {
	Version                string `xml:"version" json:"version"`
	OrganizationCode       string `xml:"ins_cd" json:"ins_cd"`
	FuYouMerchantCode      string `xml:"mchnt_cd" json:"mchnt_cd"`
	TerminalId             string `xml:"term_id" json:"term_id"`
	RandomString           string `xml:"random_str" json:"random_str"`
	ChannelMerchantOrderId string `xml:"mchnt_order_no" json:"mchnt_order_no"`
	OrderType              string `xml:"order_type" json:"order_type"`
}
