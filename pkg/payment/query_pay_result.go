package payment

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/transaction/parameters"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/fuyou"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/structtomap"
)

var TransactionNotExist = fmt.Errorf("transaction not exist")

func IsTransactionNotExist(err error) bool {
	return err == TransactionNotExist
}

var FuYouDataNotCompletion = fmt.Errorf("fu you data not completion")

func IsFuYouDataNotCompletion(err error) bool {
	return err == FuYouDataNotCompletion
}

func QueryPayStatus(logger *logrus.Entry, transactionId uint64) (orderStatus uint8, err error) {

	transaction, err := FindByTransactionId(transactionId)
	if err != nil {
		logger.WithError(err).Info("get transaction error")
		return
	}

	if transaction == nil {
		logger.WithError(err).Warn("transaction not found")
		err = TransactionNotExist
		return
	}

	return QueryPayResultByTransaction(logger, transaction)
}

// @return model.PayOrderStatusPaying, model.PayOrderStatusSuccess, model.PayOrderStatusFailed, model.PayOrderStatusClose
func QueryPayResultByTransaction(logger *logrus.Entry, transaction *model.Transaction) (orderStatus uint8, err error) {

	if len(transaction.ChannelMerchantOrderId) == 0 ||
		len(transaction.ChannelMerchantId) == 0 {

		err = FuYouDataNotCompletion
		return
	}

	request := createQueryPayResultRequestData(transaction)

	requestData, err := structtomap.StructToMap(request)
	if err != nil {
		logger.WithError(err).Error("convert struct to map error")
		return
	}

	responseData, err := fuyou.NewRequest(logger).SendRequest(requestData, getQueryPayResultURL())

	if err != nil {
		logger.WithError(err).Error("call api error")
		return
	}

	logger = logger.WithField("responseData", responseData)

	orderStatus = parseQueryPayResultResponseToReply(responseData)
	return
}

func createQueryPayResultRequestData(transaction *model.Transaction) *parameters.QueryPayResultRequest {

	return &parameters.QueryPayResultRequest{
		Version:                constant.FuYouAPIVersion,
		OrganizationCode:       config.Config.FuYou.GetOrganizationCode(),
		FuYouMerchantCode:      transaction.ChannelMerchantId,
		TerminalId:             constant.FuYouTerminalId,
		RandomString:           strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
		ChannelMerchantOrderId: transaction.ChannelMerchantOrderId,
		OrderType:              ConvertModelOrderTypeToFuYou(transaction.OrderType),
	}
}

func parseQueryPayResultResponseToReply(responseData map[string]string) uint8 {

	return transformPayStatus(responseData["trans_stat"])
}

func transformPayStatus(orderStatus string) uint8 {
	switch orderStatus {
	case "SUCCESS":
		return model.PayOrderStatusSuccess
	case "PAYERROR":
		return model.PayOrderStatusFailed
	case "CLOSED":
		return model.PayOrderStatusClose
	}
	return model.PayOrderStatusPaying
}

func getQueryPayResultURL() string {
	return fmt.Sprintf("https://%s/commonQuery", config.Config.FuYou.GetHost())
}
