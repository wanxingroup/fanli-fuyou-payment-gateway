package transaction

import (
	"testing"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/request/requestid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/protos"
)

func TestService_PrePay(t *testing.T) {

	tests := []struct {
		ctx   context.Context
		input *protos.PrePayRequest
		want  *protos.PrePayReply
		err   error
	}{
		{
			ctx: context.WithValue(context.Background(), requestid.Key, "test-request-id"),
			input: &protos.PrePayRequest{
				SourceId:        10003,
				SourceType:      protos.SourceType_CardService,
				OrderAmount:     100,
				MerchantId:      100,
				FuyouMerchantId: "fuyou100",
				Description:     "i have a dream",
				PayType:         protos.PaymentType_LETPAY,
				OpenId:          "wechat-openid",
				UserIP:          "127.0.0.1",
				AppId:           "fakeAppId",
			},
			want: &protos.PrePayReply{
				MobilePaymentString: `{"appId":"wxfa089da95020ba1a","timeStamp":"1557649678","signType":"RSA","package":"prepay_id=wx121627586241613c7d311e003087612412","nonceStr":"683cedb4a5894ca9a681d4172ee6dab0","paySign":"hvMJreiPBOCrG45XppgPxmZnxL/mojcgv4aKKU0OGJNB5Cf/0e9CwhF8Uy4OOmjuwv4l/65xvCbom0UcNwJpQ9DQUc8lMRwYUwgKWsXgVcpDW0rw3yZLX7otQJQi8JGBagq3si55yIyca+1Pjg1sOXuRnh+/2yRHPj3wR2A+tqSM1Rl5Q0PrKmWgDocIU/5wIC4B/8uK7Kf+K/HtImaag/TTNsEKg/TEzzEshfhEVl5dYvMcoq1D4FkH/938eVmpUMqGI3DQMnBa1Oj6T12sjF9p/aPdV7S9OyVxPHUGNMqDoouXnACTd4HmkU3+G7/7oyYyozi2V+JDPX1lrO+Y1w=="}`,
			},
		},
		{
			ctx: context.WithValue(context.Background(), requestid.Key, "test-request-id"),
			input: &protos.PrePayRequest{
				SourceId:        10003,
				SourceType:      protos.SourceType_CardService,
				OrderAmount:     100,
				MerchantId:      100,
				FuyouMerchantId: "fuyou100",
				Description:     "i have a dream",
				PayType:         protos.PaymentType_LETPAY,
				OpenId:          "wechat-openid",
				UserIP:          "127.0.0.1",
				AppId:           "fakeAppId",
			},
			want: &protos.PrePayReply{
				MobilePaymentString: `{"appId":"wxfa089da95020ba1a","timeStamp":"1557649678","signType":"RSA","package":"prepay_id=wx121627586241613c7d311e003087612412","nonceStr":"683cedb4a5894ca9a681d4172ee6dab0","paySign":"hvMJreiPBOCrG45XppgPxmZnxL/mojcgv4aKKU0OGJNB5Cf/0e9CwhF8Uy4OOmjuwv4l/65xvCbom0UcNwJpQ9DQUc8lMRwYUwgKWsXgVcpDW0rw3yZLX7otQJQi8JGBagq3si55yIyca+1Pjg1sOXuRnh+/2yRHPj3wR2A+tqSM1Rl5Q0PrKmWgDocIU/5wIC4B/8uK7Kf+K/HtImaag/TTNsEKg/TEzzEshfhEVl5dYvMcoq1D4FkH/938eVmpUMqGI3DQMnBa1Oj6T12sjF9p/aPdV7S9OyVxPHUGNMqDoouXnACTd4HmkU3+G7/7oyYyozi2V+JDPX1lrO+Y1w=="}`,
			},
		},
	}

	for _, test := range tests {
		config.Config.FuYou.Host = "dev-api.wanxingrowth.com/mock/131"
		svc := &Service{}
		reply, err := svc.PrePay(test.ctx, test.input)
		if reply != nil && test.want != nil {
			test.want.TransactionId = reply.TransactionId
		}
		assert.Equal(t, test.want, reply, test)
		assert.Equal(t, test.err, err, test)
	}
}
