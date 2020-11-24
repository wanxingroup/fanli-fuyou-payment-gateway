package payment

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/jinzhu/gorm"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

func FindByTransactionId(transactionId uint64) (*model.Transaction, error) {
	transaction := &model.Transaction{
		TransactionId: transactionId,
	}
	err := database.GetDB(constant.DatabaseConfigKey).
		Preload("TransactionSources").
		Where(transaction).
		First(transaction).
		Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return transaction, nil
}

func FindByTransactionIds(transactionIds []uint64) ([]*model.Transaction, error) {

	list := make([]*model.Transaction, 0, len(transactionIds))

	err := database.GetDB(constant.DatabaseConfigKey).
		Preload("TransactionSources").
		Model(model.Transaction{}).
		Where("`transactionId` IN (?)", transactionIds).
		Find(&list).
		Error

	if err != nil {

		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return list, nil
}

func FindByChannelMerchantOrderId(channelMerchantOrderId string) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ChannelMerchantOrderId: channelMerchantOrderId,
	}
	err := database.GetDB(constant.DatabaseConfigKey).
		Preload("TransactionSources").
		Where(transaction).
		First(transaction).
		Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return transaction, nil
}

func FindBySourceId(sourceId uint64, sourceType uint8) ([]*model.Transaction, error) {

	records := make([]*model.TransactionSource, 0)

	err := database.GetDB(constant.DatabaseConfigKey).
		Where("`sourceId` = ? AND `sourceType` = ?", sourceId, sourceType).
		Find(&records).
		Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	transactionIds := make([]uint64, 0, len(records))

	for _, record := range records {

		transactionIds = append(transactionIds, record.TransactionId)
	}

	return FindByTransactionIds(transactionIds)
}

func Update(transaction *model.Transaction) (err error) {
	err = database.GetDB(constant.DatabaseConfigKey).Model(&transaction).UpdateColumns(transaction).Error
	return
}
