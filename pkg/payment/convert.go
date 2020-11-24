package payment

import (
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

func ConvertModelOrderTypeToFuYou(orderType uint8) string {

	switch orderType {
	case model.OrderTypeWeChat:
		return constant.FuYouOrderTypeWeChat
	case model.OrderTypeAliPay:
		return constant.FuYouOrderTypeAliPay
	case model.OrderTypeUnionPay:
		return constant.FuYouOrderTypeUnionPay
	case model.OrderTypeBestPay:
		return constant.FuYouOrderTypeBestPay
	}

	return constant.FuYouOrderTypeWeChat
}
