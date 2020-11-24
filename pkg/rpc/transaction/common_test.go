package transaction

import (
	"fmt"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/protos"
)

func TestService_ConvertError(t *testing.T) {

	tests := []struct {
		input error
		want  *protos.Error
	}{
		{
			input: fmt.Errorf("test error"),
			want: &protos.Error{
				Code:    constant.ErrorCodeInternalError,
				Message: "test error",
			},
		},
		{
			input: validation.NewError("100", "test error"),
			want: &protos.Error{
				Code:    100,
				Message: "test error",
			},
		},
		{
			input: validation.NewError("errorNumber", "test error"),
			want: &protos.Error{
				Code:    constant.ErrorCodeInternalError,
				Message: "test error",
			},
		},
	}

	for _, test := range tests {

		svc := &Service{}
		assert.Equal(t, test.want, svc.convertError(test.input), test)
	}
}

func TestService_ConvertPayTypeString(t *testing.T) {

	tests := []struct {
		input protos.PaymentType
		want  string
	}{
		{
			input: protos.PaymentType_JSAPI,
			want:  constant.FuYouPayTypeJSAPI,
		},
		{
			input: protos.PaymentType_FWC,
			want:  constant.FuYouPayTypeFWC,
		},
		{
			input: protos.PaymentType_LETPAY,
			want:  constant.FuYouPayTypeLETPAY,
		},
		{
			input: protos.PaymentType(-1),
			want:  constant.FuYouPayTypeJSAPI,
		},
	}

	for _, test := range tests {

		svc := &Service{}
		assert.Equal(t, test.want, svc.convertPayTypeString(test.input), test)
	}
}

func TestService_ConvertPayTypeModel(t *testing.T) {

	tests := []struct {
		input protos.PaymentType
		want  uint8
	}{
		{
			input: protos.PaymentType_JSAPI,
			want:  model.PayTypeOfficialAccount,
		},
		{
			input: protos.PaymentType_FWC,
			want:  model.PayTypeAliPay,
		},
		{
			input: protos.PaymentType_LETPAY,
			want:  model.PayTypeMiniProgram,
		},
		{
			input: protos.PaymentType(-1),
			want:  model.PayTypeMiniProgram,
		},
	}

	for _, test := range tests {

		svc := &Service{}
		assert.Equal(t, test.want, svc.convertPayTypeModel(test.input), test)
	}
}

func TestService_ConvertModelOrderTypeToFuYou(t *testing.T) {

	tests := []struct {
		input uint8
		want  string
	}{
		{
			input: model.OrderTypeUnknown,
			want:  constant.FuYouOrderTypeWeChat,
		},
		{
			input: model.OrderTypeWeChat,
			want:  constant.FuYouOrderTypeWeChat,
		},
		{
			input: model.OrderTypeAliPay,
			want:  constant.FuYouOrderTypeAliPay,
		},
		{
			input: model.OrderTypeUnionPay,
			want:  constant.FuYouOrderTypeUnionPay,
		},
		{
			input: model.OrderTypeBestPay,
			want:  constant.FuYouOrderTypeBestPay,
		},
	}

	for _, test := range tests {

		svc := &Service{}
		assert.Equal(t, test.want, svc.convertModelOrderTypeToFuYou(test.input), test)
	}
}
