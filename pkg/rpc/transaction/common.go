package transaction

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/protos"
)

func (svc *Service) convertError(err error) *protos.Error {

	if validationError, ok := err.(validation.Error); ok {
		errorCode, convertError := strconv.Atoi(validationError.Code())
		if convertError != nil {
			errorCode = constant.ErrorCodeInternalError
		}
		return &protos.Error{
			Code:    int64(errorCode),
			Message: validationError.Error(),
		}
	}

	return &protos.Error{
		Code:    constant.ErrorCodeInternalError,
		Message: err.Error(),
	}
}

func (svc *Service) convertPayTypeString(payType protos.PaymentType) string {

	switch payType {
	case protos.PaymentType_JSAPI:
		return constant.FuYouPayTypeJSAPI
	case protos.PaymentType_FWC:
		return constant.FuYouPayTypeFWC
	case protos.PaymentType_LETPAY:
		return constant.FuYouPayTypeLETPAY
	}

	return constant.FuYouPayTypeJSAPI
}

func (svc *Service) convertPayTypeModel(payType protos.PaymentType) uint8 {

	switch payType {
	case protos.PaymentType_JSAPI:
		return model.PayTypeOfficialAccount
	case protos.PaymentType_FWC:
		return model.PayTypeAliPay
	case protos.PaymentType_LETPAY:
		return model.PayTypeMiniProgram
	}

	return model.PayTypeMiniProgram
}

func (svc *Service) convertModelOrderTypeToFuYou(orderType uint8) string {

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
