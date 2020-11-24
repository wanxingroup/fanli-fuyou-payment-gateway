package transaction

import (
	"fmt"
	"strconv"
	"strings"

	rpcLog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/payment"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/transaction/parameters"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/fuyou"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/structtomap"
)

func (svc Service) CloseTransaction(ctx context.Context, req *protos.CloseTransactionRequest) (*protos.CloseTransactionReply, error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger()).WithField("requestData", req)

	sourceType := svc.convertProtobufSourceTypeToModel(req.GetSourceType())

	transactions, err := payment.FindBySourceId(req.GetSourceId(), sourceType)
	if err != nil {

		logger.Info("find transaction error")
		return &protos.CloseTransactionReply{
			Err: svc.convertError(err),
		}, nil
	}

	var transaction *model.Transaction
	for _, transactionData := range transactions {
		if transactionData.Forward != model.ForwardInflow {
			continue
		}

		transaction = transactionData
		break
	}

	if transaction == nil {

		logger.Info("transaction not exist")
		return &protos.CloseTransactionReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeTransactionNotExist,
				Message: constant.ErrorMessageTransactionNotExist,
			},
		}, nil
	}

	request, err := structtomap.StructToMap(svc.createCloseTransactionRequest(transaction))

	logger = logger.WithField("callAPIData", request)

	responseData, err := fuyou.NewRequest(logger).SendRequest(request, svc.getCloseTransactionRequestURL())
	if err != nil {

		logger.WithError(err).Error("call fuyou failed")
		return &protos.CloseTransactionReply{
			Err: svc.convertError(err),
		}, nil
	}

	logger = logger.WithField("response", responseData)

	reply := svc.parseCloseTransactionResponseToReply(logger, responseData)

	if reply.Err == nil {
		err = payment.Update(&model.Transaction{TransactionId: transaction.TransactionId, PayOrderStatus: model.PayOrderStatusClose})

		if err != nil {
			logger.WithError(err).Error("update transaction order status error")
		}
	}

	return reply, nil
}

func (svc Service) createCloseTransactionRequest(transaction *model.Transaction) *parameters.CloseTransactionRequest {

	return &parameters.CloseTransactionRequest{
		Version:           constant.FuYouAPIVersion,
		OrganizationCode:  config.Config.FuYou.GetOrganizationCode(),
		FuYouMerchantCode: transaction.ChannelMerchantId,
		TerminalId:        constant.FuYouTerminalId,
		RandomString:      strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
		TransactionId:     strconv.FormatUint(transaction.TransactionId, 10),
		OrderType:         payment.ConvertModelOrderTypeToFuYou(transaction.OrderType),
		AppId:             transaction.AppId,
	}
}

func (svc Service) getCloseTransactionRequestURL() string {
	return fmt.Sprintf("https://%s/closeorder", config.Config.FuYou.GetHost())
}

func (svc Service) parseCloseTransactionResponseToReply(logger *logrus.Entry, response map[string]string) *protos.CloseTransactionReply {

	if response["result_code"] == constant.FuYouSuccessCode {
		return &protos.CloseTransactionReply{}
	}

	errorCode, err := strconv.ParseInt(response["result_code"], 10, 64)
	errorMessage := response["result_msg"]
	if err != nil {
		logger.WithError(err).Error("convert error code error")
		errorCode = constant.ErrorCodeConvertErrorCodeFailed
		errorMessage = fmt.Sprintf("code: %s, message: %s", response["result_code"], response["result_msg"])
	}

	// 如果远端已经关闭支付，会返回找不到流水的情况
	// 因此直接返回没有异常，关单成功即可
	if errorCode == constant.FuYouErrorCodeTransactionNotFound {
		return &protos.CloseTransactionReply{}
	}

	return &protos.CloseTransactionReply{
		Err: &protos.Error{
			Code:    constant.ErrorCodeCloseTransactionResponseError,
			Message: constant.ErrorMessageCloseTransactionResponseError,
			Stack: &protos.Error{
				Code:    errorCode,
				Message: errorMessage,
			},
		},
	}
}

func (svc Service) convertProtobufSourceTypeToModel(sourceType protos.SourceType) uint8 {

	switch sourceType {
	case protos.SourceType_OrderService:
		return model.SourceTypeOrderService
	case protos.SourceType_CardService:
		return model.SourceTypeCardService
	case protos.SourceType_PaymentService:
		return model.SourceTypePaymentService
	}

	return model.SourceTypeOrderService
}
