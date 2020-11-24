package callback

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/client"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/client/mock"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

func TestController_Notify(t *testing.T) {

	controller := &Controller{}
	tests := []struct {
		init  func()
		input string
		want  string
		check func()
	}{
		{
			init: func() {
				database.GetDB(constant.DatabaseConfigKey).Create(&model.Transaction{
					TransactionId:          10002,
					Forward:                model.ForwardInflow,
					PayOrderStatus:         model.PayOrderStatusPaying,
					PayType:                model.PayTypeMiniProgram,
					OrderType:              model.OrderTypeWeChat,
					AppId:                  "fakeAppId",
					ChannelMerchantId:      "fuyou-merchantId",
					ChannelMerchantOrderId: "10002",
					MerchantId:             100,
					OpenId:                 "fakeOpenId",
					Amount:                 100,
					Currency:               constant.FuYouCurrencyType,
					ClientIP:               "127.0.0.1",
					ErrorCode:              "000000",
					ErrorMessage:           "",
					ExpireTime:             0,
					TracingId:              "",
					NotifyStatus:           model.NotifyStatusCreate,
					RefundRemainingAmount:  100,
					TradeTime:              time.Now(),
					Channel:                "",
					ServiceCharge:          0,
					PayInfo:                "",
					TransactionSources: []model.TransactionSource{
						{
							TransactionId: 10002,
							SourceId:      10000,
							SourceType:    model.SourceTypeCardService,
						},
					},
					Time: databases.Time{},
				})
				client.SetCardOrderMockService(&mock.CardOrderFinishOrder{})
			},
			input: `req=%3C%3Fxml+version%3D%221.0%22+encoding%3D%22GBK%22+standalone%3D%22yes%22%3F%3E%3Cxml%3E%3Ccurr_type%3ECNY%3C%2Fcurr_type%3E%3Cins_cd%3E08M0026894%3C%2Fins_cd%3E%3Cmchnt_cd%3E0002900F3057141%3C%2Fmchnt_cd%3E%3Cmchnt_order_no%3E10002%3C%2Fmchnt_order_no%3E%3Corder_amt%3E1%3C%2Forder_amt%3E%3Corder_type%3ELETPAY%3C%2Forder_type%3E%3Crandom_str%3E09VRJF5KXIZYY3XQI7D09LLJ0X3MUPH4%3C%2Frandom_str%3E%3Creserved_addn_inf%3E%3C%2Freserved_addn_inf%3E%3Creserved_bank_type%3EOTHERS%3C%2Freserved_bank_type%3E%3Creserved_buyer_logon_id%3E%3C%2Freserved_buyer_logon_id%3E%3Creserved_channel_order_id%3E14192020030216292272100049%3C%2Freserved_channel_order_id%3E%3Creserved_coupon_fee%3E0%3C%2Freserved_coupon_fee%3E%3Creserved_fund_bill_list%3E%3C%2Freserved_fund_bill_list%3E%3Creserved_fy_settle_dt%3E20200302%3C%2Freserved_fy_settle_dt%3E%3Creserved_fy_trace_no%3E030500595279%3C%2Freserved_fy_trace_no%3E%3Creserved_is_credit%3E0%3C%2Freserved_is_credit%3E%3Creserved_settlement_amt%3E1%3C%2Freserved_settlement_amt%3E%3Cresult_code%3E000000%3C%2Fresult_code%3E%3Cresult_msg%3ESUCCESS%3C%2Fresult_msg%3E%3Csettle_order_amt%3E1%3C%2Fsettle_order_amt%3E%3Csign%3EFiatOmh6FxzEab6Fs0ZemYyIl9heAM6fAfOAWG3J8KaZ1usVgdgQhnkkqWZxu%2Fk8danv7KmhpM3Du2woHsRsX495vhrec540zEbezPYNgcL6IYWJGnRvBIobDPampmNF2eRcyTak3akrabi8d3Cjlpsij2tHKtVvLucZsBkUSzo%3D%3C%2Fsign%3E%3Cterm_id%3E%3C%2Fterm_id%3E%3Ctransaction_id%3E4200000477202003028039906349%3C%2Ftransaction_id%3E%3Ctxn_fin_ts%3E20200302162954%3C%2Ftxn_fin_ts%3E%3Cuser_id%3Eo2elW48bhulySYwDPjI6QTNISl94%3C%2Fuser_id%3E%3C%2Fxml%3E`,
			want:  "1",
			check: func() {
				var transaction = &model.Transaction{}
				err := database.GetDB(constant.DatabaseConfigKey).Find(&model.Transaction{
					TransactionId: 10002,
				}).First(transaction).Error
				assert.Nil(t, err)
				assert.NotNil(t, transaction)
				assert.Equal(t, model.PayOrderStatusSuccess, transaction.PayOrderStatus, transaction)
			},
		},
		{
			input: `req=%3C%3Fxml+version%3D%221.0%22+encoding%3D%22GBK%22+standalone%3D%22yes%22%3F%3E%3Cxml%3E%3Ccurr_type%3ECNY%3C%2Fcurr_type%3E%3Cins_cd%3E08M0026894%3C%2Fins_cd%3E%3Cmchnt_cd%3E0002900F3057141%3C%2Fmchnt_cd%3E%3Cmchnt_order_no%3E10002%3C%2Fmchnt_order_no%3E%3Corder_amt%3E1%3C%2Forder_amt%3E%3Corder_type%3ELETPAY%3C%2Forder_type%3E%3Crandom_str%3E09VRJF5KXIZYY3XQI7D09LLJ0X3MUPH4%3C%2Frandom_str%3E%3Creserved_addn_inf%3E%3C%2Freserved_addn_inf%3E%3Creserved_bank_type%3EOTHERS%3C%2Freserved_bank_type%3E%3Creserved_buyer_logon_id%3E%3C%2Freserved_buyer_logon_id%3E%3Creserved_channel_order_id%3E14192020030216292272100049%3C%2Freserved_channel_order_id%3E%3Creserved_coupon_fee%3E0%3C%2Freserved_coupon_fee%3E%3Creserved_fund_bill_list%3E%3C%2Freserved_fund_bill_list%3E%3Creserved_fy_settle_dt%3E20200302%3C%2Freserved_fy_settle_dt%3E%3Creserved_fy_trace_no%3E030500595279%3C%2Freserved_fy_trace_no%3E%3Creserved_is_credit%3E0%3C%2Freserved_is_credit%3E%3Creserved_settlement_amt%3E1%3C%2Freserved_settlement_amt%3E%3Cresult_code%3E000000%3C%2Fresult_code%3E%3Cresult_msg%3ESUCCESS%3C%2Fresult_msg%3E%3Csettle_order_amt%3E1%3C%2Fsettle_order_amt%3E%3Csign%3EFiatOmh6FxzEab6Fs0ZemYyIl9heAM6fAfOAWG3J8KaZ1usVgdgQhnkkqWZxu%2Fk8danv7KmhpM3Du2woHsRsX495vhrec540zEbezPYNgcL6IYWJGnRvBIobDPampmNF2eRcyTak3akrabi8d3Cjlpsij2tHKtVvLucZsBkUSzo%3D%3C%2Fsign%3E%3Cterm_id%3E%3C%2Fterm_id%3E%3Ctransaction_id%3E4200000477202003028039906349%3C%2Ftransaction_id%3E%3Ctxn_fin_ts%3E20200302162954%3C%2Ftxn_fin_ts%3E%3Cuser_id%3Eo2elW48bhulySYwDPjI6QTNISl94%3C%2Fuser_id%3E%3C%2Fxml%3E`,
			want:  "1",
		},
		{
			input: `req=%3C%3Fxml+version%3D%221.0%22+encoding%3D%22GBK%22+standalone%3D%22yes%22%3F%3E%3Cxml%3E%3Ccurr_type%3ECNY%3C%2Fcurr_type%3E%3Cins_cd%3E08M0026894%3C%2Fins_cd%3E%3Cmchnt_cd%3E0002900F3057141%3C%2Fmchnt_cd%3E%3Cmchnt_order_no%3EerrorCode%3C%2Fmchnt_order_no%3E%3Corder_amt%3E1%3C%2Forder_amt%3E%3Corder_type%3ELETPAY%3C%2Forder_type%3E%3Crandom_str%3E09VRJF5KXIZYY3XQI7D09LLJ0X3MUPH4%3C%2Frandom_str%3E%3Creserved_addn_inf%3E%3C%2Freserved_addn_inf%3E%3Creserved_bank_type%3EOTHERS%3C%2Freserved_bank_type%3E%3Creserved_buyer_logon_id%3E%3C%2Freserved_buyer_logon_id%3E%3Creserved_channel_order_id%3E14192020030216292272100049%3C%2Freserved_channel_order_id%3E%3Creserved_coupon_fee%3E0%3C%2Freserved_coupon_fee%3E%3Creserved_fund_bill_list%3E%3C%2Freserved_fund_bill_list%3E%3Creserved_fy_settle_dt%3E20200302%3C%2Freserved_fy_settle_dt%3E%3Creserved_fy_trace_no%3E030500595279%3C%2Freserved_fy_trace_no%3E%3Creserved_is_credit%3E0%3C%2Freserved_is_credit%3E%3Creserved_settlement_amt%3E1%3C%2Freserved_settlement_amt%3E%3Cresult_code%3E000000%3C%2Fresult_code%3E%3Cresult_msg%3ESUCCESS%3C%2Fresult_msg%3E%3Csettle_order_amt%3E1%3C%2Fsettle_order_amt%3E%3Csign%3EFiatOmh6FxzEab6Fs0ZemYyIl9heAM6fAfOAWG3J8KaZ1usVgdgQhnkkqWZxu%2Fk8danv7KmhpM3Du2woHsRsX495vhrec540zEbezPYNgcL6IYWJGnRvBIobDPampmNF2eRcyTak3akrabi8d3Cjlpsij2tHKtVvLucZsBkUSzo%3D%3C%2Fsign%3E%3Cterm_id%3E%3C%2Fterm_id%3E%3Ctransaction_id%3E4200000477202003028039906349%3C%2Ftransaction_id%3E%3Ctxn_fin_ts%3E20200302162954%3C%2Ftxn_fin_ts%3E%3Cuser_id%3Eo2elW48bhulySYwDPjI6QTNISl94%3C%2Fuser_id%3E%3C%2Fxml%3E`,
			want:  "",
		},
		{
			input: `req=%3C%3Fxml+version%3D%221.0%22+encoding%3D%22GBK%22+standalone%3D%22yes%22%3F%3E%3Cxml%3E%3Ccurr_type%3ECNY%3C%2Fcurr_type%3E%3Cins_cd%3E08M0026894%3C%2Fins_cd%3E%3Cmchnt_cd%3E0002900F3057141%3C%2Fmchnt_cd%3E%3Cmchnt_order_no%3E100000004%3C%2Fmchnt_order_no%3E%3Corder_amt%3E1%3C%2Forder_amt%3E%3Corder_type%3ELETPAY%3C%2Forder_type%3E%3Crandom_str%3E09VRJF5KXIZYY3XQI7D09LLJ0X3MUPH4%3C%2Frandom_str%3E%3Creserved_addn_inf%3E%3C%2Freserved_addn_inf%3E%3Creserved_bank_type%3EOTHERS%3C%2Freserved_bank_type%3E%3Creserved_buyer_logon_id%3E%3C%2Freserved_buyer_logon_id%3E%3Creserved_channel_order_id%3E14192020030216292272100049%3C%2Freserved_channel_order_id%3E%3Creserved_coupon_fee%3E0%3C%2Freserved_coupon_fee%3E%3Creserved_fund_bill_list%3E%3C%2Freserved_fund_bill_list%3E%3Creserved_fy_settle_dt%3E20200302%3C%2Freserved_fy_settle_dt%3E%3Creserved_fy_trace_no%3E030500595279%3C%2Freserved_fy_trace_no%3E%3Creserved_is_credit%3E0%3C%2Freserved_is_credit%3E%3Creserved_settlement_amt%3E1%3C%2Freserved_settlement_amt%3E%3Cresult_code%3E000000%3C%2Fresult_code%3E%3Cresult_msg%3ESUCCESS%3C%2Fresult_msg%3E%3Csettle_order_amt%3E1%3C%2Fsettle_order_amt%3E%3Csign%3EFiatOmh6FxzEab6Fs0ZemYyIl9heAM6fAfOAWG3J8KaZ1usVgdgQhnkkqWZxu%2Fk8danv7KmhpM3Du2woHsRsX495vhrec540zEbezPYNgcL6IYWJGnRvBIobDPampmNF2eRcyTak3akrabi8d3Cjlpsij2tHKtVvLucZsBkUSzo%3D%3C%2Fsign%3E%3Cterm_id%3E%3C%2Fterm_id%3E%3Ctransaction_id%3E4200000477202003028039906349%3C%2Ftransaction_id%3E%3Ctxn_fin_ts%3E20200302162954%3C%2Ftxn_fin_ts%3E%3Cuser_id%3Eo2elW48bhulySYwDPjI6QTNISl94%3C%2Fuser_id%3E%3C%2Fxml%3E`,
			want:  "",
		},
		{
			init: func() {
				database.GetDB(constant.DatabaseConfigKey).Create(&model.Transaction{
					TransactionId:          10003,
					Forward:                model.ForwardInflow,
					PayOrderStatus:         model.PayOrderStatusPaying,
					PayType:                model.PayTypeMiniProgram,
					OrderType:              model.OrderTypeWeChat,
					AppId:                  "fakeAppId",
					ChannelMerchantId:      "fuyou-merchantId",
					ChannelMerchantOrderId: "10003",
					MerchantId:             100,
					OpenId:                 "fakeOpenId",
					Amount:                 100,
					Currency:               constant.FuYouCurrencyType,
					ClientIP:               "127.0.0.1",
					ErrorCode:              "000000",
					ErrorMessage:           "",
					ExpireTime:             0,
					TracingId:              "",
					NotifyStatus:           model.NotifyStatusCreate,
					RefundRemainingAmount:  100,
					TradeTime:              time.Now(),
					Channel:                "",
					ServiceCharge:          0,
					PayInfo:                "",
					TransactionSources: []model.TransactionSource{
						{
							TransactionId: 10003,
							SourceId:      10001,
							SourceType:    model.SourceTypeOrderService,
						},
					},
					Time: databases.Time{},
				})
				client.SetCardOrderMockService(&mock.CardOrderFinishOrder{})
			},
			input: `req=%3C%3Fxml+version%3D%221.0%22+encoding%3D%22GBK%22+standalone%3D%22yes%22%3F%3E%3Cxml%3E%3Ccurr_type%3ECNY%3C%2Fcurr_type%3E%3Cins_cd%3E08M0026894%3C%2Fins_cd%3E%3Cmchnt_cd%3E0002900F3057141%3C%2Fmchnt_cd%3E%3Cmchnt_order_no%3E10002%3C%2Fmchnt_order_no%3E%3Corder_amt%3E1%3C%2Forder_amt%3E%3Corder_type%3ELETPAY%3C%2Forder_type%3E%3Crandom_str%3E09VRJF5KXIZYY3XQI7D09LLJ0X3MUPH4%3C%2Frandom_str%3E%3Creserved_addn_inf%3E%3C%2Freserved_addn_inf%3E%3Creserved_bank_type%3EOTHERS%3C%2Freserved_bank_type%3E%3Creserved_buyer_logon_id%3E%3C%2Freserved_buyer_logon_id%3E%3Creserved_channel_order_id%3E14192020030216292272100049%3C%2Freserved_channel_order_id%3E%3Creserved_coupon_fee%3E0%3C%2Freserved_coupon_fee%3E%3Creserved_fund_bill_list%3E%3C%2Freserved_fund_bill_list%3E%3Creserved_fy_settle_dt%3E20200302%3C%2Freserved_fy_settle_dt%3E%3Creserved_fy_trace_no%3E030500595279%3C%2Freserved_fy_trace_no%3E%3Creserved_is_credit%3E0%3C%2Freserved_is_credit%3E%3Creserved_settlement_amt%3E1%3C%2Freserved_settlement_amt%3E%3Cresult_code%3E000000%3C%2Fresult_code%3E%3Cresult_msg%3ESUCCESS%3C%2Fresult_msg%3E%3Csettle_order_amt%3E1%3C%2Fsettle_order_amt%3E%3Csign%3EFiatOmh6FxzEab6Fs0ZemYyIl9heAM6fAfOAWG3J8KaZ1usVgdgQhnkkqWZxu%2Fk8danv7KmhpM3Du2woHsRsX495vhrec540zEbezPYNgcL6IYWJGnRvBIobDPampmNF2eRcyTak3akrabi8d3Cjlpsij2tHKtVvLucZsBkUSzo%3D%3C%2Fsign%3E%3Cterm_id%3E%3C%2Fterm_id%3E%3Ctransaction_id%3E4200000477202003028039906349%3C%2Ftransaction_id%3E%3Ctxn_fin_ts%3E20200302162954%3C%2Ftxn_fin_ts%3E%3Cuser_id%3Eo2elW48bhulySYwDPjI6QTNISl94%3C%2Fuser_id%3E%3C%2Fxml%3E`,
			want:  "1",
			check: func() {
				var transaction = &model.Transaction{}
				err := database.GetDB(constant.DatabaseConfigKey).Find(&model.Transaction{
					TransactionId: 10003,
				}).First(transaction).Error
				assert.Nil(t, err)
				assert.NotNil(t, transaction)
				assert.Equal(t, model.PayOrderStatusSuccess, transaction.PayOrderStatus, transaction)
			},
		},
	}

	for _, test := range tests {

		if test.init != nil {
			test.init()
		}

		var err error
		assert.Nil(t, err)
		bodyReader := strings.NewReader(test.input)

		resp := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(resp)
		ctx.Request = httptest.NewRequest("POST", "/api/paymentgateway/fuyou/callback", bodyReader)
		ctx.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		controller.Notify(ctx)
		resp.Flush()
		assert.Equal(t, test.want, resp.Body.String(), test)

		if test.check != nil {
			test.check()
		}
	}
}
