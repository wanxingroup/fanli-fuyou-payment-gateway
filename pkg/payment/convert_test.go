package payment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

func TestConvertModelOrderTypeToFuYou(t *testing.T) {

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

		assert.Equal(t, test.want, ConvertModelOrderTypeToFuYou(test.input), test)
	}
}
