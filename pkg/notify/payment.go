package notify

import (
	"fmt"

	paymentProtos "dev-gitlab.wanxingrowth.com/fanli/payment/pkg/rpc/protos"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/client"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

type PaymentService struct {
	logger *logrus.Entry
}

func (service *PaymentService) Notify(transaction *model.Transaction) error {

	service.logger.Debug("starting to call finish order")

	api := client.GetPaymentService()
	if api == nil {

		service.logger.Error("payment service is nil")
		return fmt.Errorf("payment service not initialize, please configurate configuration file and restart service")
	}

	for _, source := range transaction.TransactionSources {

		reply, err := api.FinishPayment(context.Background(), &paymentProtos.FinishPaymentRequest{PaymentId: source.SourceId})
		if err != nil {
			service.logger.WithError(err).Error("call finish payment error")
			return err
		}

		if reply.GetErr() != nil {
			service.logger.WithField("replyError", reply.GetErr()).Error("finish order reply false")
			return fmt.Errorf("it was notified payment service, but result false")
		}

		service.logger.WithField("reply", reply).Debug("finish order reply true")
	}

	service.logger.Debug("call finish order completed")

	return nil
}
