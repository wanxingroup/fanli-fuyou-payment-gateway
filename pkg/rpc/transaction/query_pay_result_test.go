package transaction

import (
	"testing"
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/request/requestid"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/client"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/client/mock"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/protos"
)

func TestService_QueryPayResult(t *testing.T) {

	tests := []struct {
		init  func()
		ctx   context.Context
		input *protos.QueryPayResultRequest
		want  *protos.QueryPayResultReply
		err   error
	}{
		{
			init: func() {
				database.GetDB(constant.DatabaseConfigKey).Create(&model.Transaction{
					TransactionId:          10000,
					Forward:                model.ForwardInflow,
					PayOrderStatus:         model.PayOrderStatusSuccess,
					PayType:                model.PayTypeMiniProgram,
					OrderType:              model.OrderTypeWeChat,
					AppId:                  "fakeAppId",
					ChannelMerchantId:      "fuyou-merchantId",
					ChannelMerchantOrderId: "10000",
					MerchantId:             100,
					OpenId:                 "fakeOpenId",
					Amount:                 100,
					Currency:               constant.FuYouCurrencyType,
					ClientIP:               "127.0.0.1",
					ErrorCode:              "000000",
					ErrorMessage:           "",
					ExpireTime:             0,
					TracingId:              "",
					NotifyStatus:           model.NotifyStatusSuccess,
					RefundRemainingAmount:  100,
					TradeTime:              time.Now(),
					Channel:                "",
					ServiceCharge:          0,
					PayInfo:                "",
					TransactionSources: []model.TransactionSource{
						{
							TransactionId: 10000,
							SourceId:      10000,
							SourceType:    model.SourceTypeCardService,
						},
					},
					Time: databases.Time{},
				})

				client.SetCardOrderMockService(mock.CardOrderFinishOrder{})
			},
			ctx: context.WithValue(context.Background(), requestid.Key, "test-request-id"),
			input: &protos.QueryPayResultRequest{
				SourceId:   10000,
				SourceType: protos.SourceType_CardService,
			},
			want: &protos.QueryPayResultReply{
				PayStatus: protos.PayStatus_Succeed,
			},
		},
	}

	for _, test := range tests {
		config.Config.FuYou.Host = "dev-api.wanxingrowth.com/mock/131"
		svc := &Service{}

		if test.init != nil {
			test.init()
		}

		reply, err := svc.QueryPayResult(test.ctx, test.input)
		assert.Equal(t, test.want, reply, test)
		assert.Equal(t, test.err, err, test)
	}
}
