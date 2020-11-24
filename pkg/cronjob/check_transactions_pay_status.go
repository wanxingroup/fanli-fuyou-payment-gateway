package cronjob

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/payment"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

func CheckTransactionsPayStatus() {

	var lastTransactionId uint64
	for {
		transactions := make([]*model.Transaction, 0)
		query := database.GetDB(constant.DatabaseConfigKey).
			Model(model.Transaction{}).
			Preload("TransactionSources").
			Where("`transactionId` > ? AND `payOrderStatus` = ?", lastTransactionId, model.PayOrderStatusPaying).
			Order("`transactionId` ASC")
		err := query.
			Find(&transactions).
			Error

		if err != nil {

			log.GetLogger().WithError(err).Error("get transaction failed")
			break
		}

		var count uint64
		err = query.Count(&count).Error
		if err != nil {

			log.GetLogger().WithError(err).Error("get transaction count failed")
			break
		}

		if count == 0 {
			break
		}

		for _, transaction := range transactions {

			lastTransactionId = transaction.TransactionId

			logger := log.GetLogger().WithField("transaction", transaction)
			payOrderStatus, err := payment.QueryPayStatus(logger, transaction.TransactionId)
			if err != nil {

				if payment.IsFuYouDataNotCompletion(err) {

					continue
				}
				logger.WithError(err).Error("query pay status error")
				continue
			}

			err = payment.UpdateOrderStatus(logger, transaction, payOrderStatus)
			if err != nil {
				logger.WithError(err).Error("update order status failed")
			}
		}
	}
}
