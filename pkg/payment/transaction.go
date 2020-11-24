package payment

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/notify"
)

func UpdateOrderStatus(logger *logrus.Entry, transaction *model.Transaction, payOrderStatus uint8) (err error) {

	switch payOrderStatus {
	case model.PayOrderStatusPaying:
		// 支付中就不处理
		break
	default:
		// 不认识的状态也就不处理
		break
	case model.PayOrderStatusClose, model.PayOrderStatusFailed:
		// 更新数据库
		err = database.GetDB(constant.DatabaseConfigKey).
			Model(model.Transaction{}).
			Where(model.Transaction{TransactionId: transaction.TransactionId, PayOrderStatus: model.PayOrderStatusPaying}).
			Update(model.Transaction{PayOrderStatus: payOrderStatus}).
			Error
		break
	case model.PayOrderStatusSuccess:

		result := database.GetDB(constant.DatabaseConfigKey).
			Model(model.Transaction{}).
			Where(model.Transaction{TransactionId: transaction.TransactionId, PayOrderStatus: model.PayOrderStatusPaying}).
			Update(model.Transaction{PayOrderStatus: payOrderStatus})
		err = result.Error
		if err != nil {
			logger.WithError(err).Error("update transaction payOrderStatus error")
			break
		}

		if result.RowsAffected == 0 {
			return
		}

		err = notify.NewNotify(logger).Notify(transaction)
		if err != nil {
			logger.WithError(err).Error("notify service payStatus failed")
		}
	}

	return
}
