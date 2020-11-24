package notify

import (
	"fmt"
	"strconv"

	orderProtos "dev-gitlab.wanxingrowth.com/fanli/order/v2/pkg/rpc/protos"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/client"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

type OrderService struct {
	logger *logrus.Entry
}

const PaymentModeOfficialAccount = "officialAccount"
const PaymentModeMiniProgram = "miniProgram"
const PaymentModeAlipay = "alipay"
const PaymentModeOther = "other"

const PaymentProductWechat = "wechat"
const PaymentProductAlipay = "alipay"
const PaymentProductOther = "other"

func (service *OrderService) Notify(transaction *model.Transaction) error {

	service.logger.Debug("starting to call finish order")

	api := client.GetOrderService()
	if api == nil {

		service.logger.Error("order service is nil")
		return fmt.Errorf("order service not initialize, please configurate configuration file and restart service")
	}

	for _, source := range transaction.TransactionSources {

		reply, err := api.PaidOrder(context.Background(), &orderProtos.PaidOrderRequest{
			OrderId: source.SourceId,
			Payment: &orderProtos.Payment{
				TransactionId:  strconv.FormatUint(transaction.TransactionId, 10),
				PaidPrice:      uint64(transaction.Amount),
				PaymentChannel: transaction.Channel,
				PaymentMode:    service.convertPayTypeToPaymentMode(transaction.PayType),
				PaymentProduct: service.convertPayTypeToPaymentProduct(transaction.PayType),
				PaidTime:       uint64(transaction.TradeTime.Unix()),
			},
		})
		if err != nil {
			service.logger.WithError(err).Error("call finish order error")
			return err
		}

		if reply.GetError() != nil {
			service.logger.WithField("replyError", reply.GetError()).Error("finish order reply error")
			return fmt.Errorf("it was notified order service, but result error")
		}

		service.logger.WithField("reply", reply).Debug("finish order reply true")
	}

	service.logger.Debug("call finish order completed")

	return nil
}

func (service *OrderService) convertPayTypeToPaymentMode(payType uint8) string {

	switch payType {
	case model.PayTypeOfficialAccount:
		return PaymentModeOfficialAccount
	case model.PayTypeMiniProgram:
		return PaymentModeMiniProgram
	case model.PayTypeAliPay:
		return PaymentModeAlipay
	}

	return PaymentModeOther
}

func (service *OrderService) convertPayTypeToPaymentProduct(payType uint8) string {

	switch payType {
	case model.PayTypeOfficialAccount:
		return PaymentProductWechat
	case model.PayTypeMiniProgram:
		return PaymentProductWechat
	case model.PayTypeAliPay:
		return PaymentProductAlipay
	}

	return PaymentProductOther
}
