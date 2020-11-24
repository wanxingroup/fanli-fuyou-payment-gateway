package notify

import (
	"fmt"

	cardProtos "dev-gitlab.wanxingrowth.com/fanli/card/pkg/rpc/protos"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/client"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

type CardService struct {
	logger *logrus.Entry
}

func (service *CardService) Notify(transaction *model.Transaction) error {

	service.logger.Debug("starting to call finish order")

	api := client.GetCardOrderService()
	if api == nil {

		service.logger.Error("card service is nil")
		return fmt.Errorf("card service not initialize, please configurate configuration file and restart service")
	}

	for _, source := range transaction.TransactionSources {

		reply, err := api.FinishOrder(context.Background(), &cardProtos.FinishOrderRequest{OrderId: source.SourceId})
		if err != nil {
			service.logger.WithError(err).Error("call finish order error")
			return err
		}

		if reply.GetErr() != nil {
			service.logger.WithField("replyError", reply.GetErr()).Error("finish order reply false")
			return fmt.Errorf("it was notified card service, but result false")
		}

		service.logger.WithField("reply", reply).Debug("finish order reply true")
	}

	service.logger.Debug("call finish order completed")

	return nil
}
