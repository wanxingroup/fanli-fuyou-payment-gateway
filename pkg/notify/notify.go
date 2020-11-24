package notify

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

type Notify struct {
	logger *logrus.Entry
}

func NewNotify(logger *logrus.Entry) *Notify {
	return &Notify{
		logger: logger,
	}
}

func (n *Notify) Notify(transaction *model.Transaction) (err error) {

	n.logger = n.logger.WithField("transaction", transaction)
	n.logger.Debug("starting to notify")

	for _, source := range transaction.TransactionSources {

		var service Service
		switch source.SourceType {
		case model.SourceTypeOrderService:
			n.logger.Debug("init order service")
			service = &OrderService{logger: n.logger}

		case model.SourceTypeCardService:
			n.logger.Debug("init card service")
			service = &CardService{logger: n.logger}

		case model.SourceTypePaymentService:
			n.logger.Debug("init payment service")
			service = &PaymentService{logger: n.logger}

		default:
			n.logger.WithField("source", source).
				Warn("not match source type")
			continue
		}

		n.logger.WithField("service", service).
			Debug("notify service")

		notifyError := service.Notify(transaction)
		if notifyError != nil {

			n.logger.WithField("service", service).
				WithError(notifyError).
				Error("notify error")
			if err != nil {
				err = errors.Wrap(err, notifyError.Error())
			} else {
				err = notifyError
			}
		}
	}

	if err != nil {
		n.logger.WithError(err).Errorf("notify service error")
	} else {

		n.logger.Debug("notify completed")
	}

	return err
}
