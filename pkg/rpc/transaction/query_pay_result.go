package transaction

import (
	"fmt"

	rpcLog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/payment"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

func (svc *Service) QueryPayResult(ctx context.Context, req *protos.QueryPayResultRequest) (*protos.QueryPayResultReply, error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger()).WithField("requestData", req)
	if req == nil {

		logger.Warn("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	transactions, err := payment.FindBySourceId(req.GetSourceId(), uint8(req.GetSourceType()))
	if err != nil {
		logger.WithError(err).Info("get transaction error")
		return &protos.QueryPayResultReply{
			Err: svc.convertError(err),
		}, nil
	}

	if len(transactions) == 0 {
		logger.WithError(err).Warn("transaction not found")
		return &protos.QueryPayResultReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeTransactionNotExist,
				Message: constant.ErrorMessageTransactionNotExist,
			},
		}, nil
	}

	logger = logger.WithField("transactions", transactions)

	for _, transaction := range transactions {

		if transaction.Forward == model.ForwardOutflow {
			continue
		}

		orderStatus, err := payment.QueryPayResultByTransaction(logger, transaction)
		if err != nil {
			logger.WithError(err).Error("query pay result error")
		}

		if transaction.PayOrderStatus != orderStatus {

			err = payment.UpdateOrderStatus(logger, transaction, orderStatus)
			if err != nil {
				logger.WithError(err).Error("update order orderStatus failed")
			}
		}

		return &protos.QueryPayResultReply{
			PayStatus: svc.transformPayStatus(orderStatus),
		}, nil
	}

	// 只返回正向流水的状态（将来如果需要展示反向流水的话，可以再扩容对应的处理方式）
	return &protos.QueryPayResultReply{
		Err: &protos.Error{
			Code:    constant.ErrorCodeTransactionNotExist,
			Message: constant.ErrorMessageTransactionNotExist,
		},
	}, nil
}

func (svc *Service) getQueryPayResultURL() string {
	return fmt.Sprintf("https://%s/commonQuery", config.Config.FuYou.GetHost())
}

func (svc *Service) transformPayStatus(payStatus uint8) protos.PayStatus {
	switch payStatus {
	case model.PayOrderStatusSuccess:
		return protos.PayStatus_Succeed
	case model.PayOrderStatusFailed:
		return protos.PayStatus_Failed
	case model.PayOrderStatusClose:
		return protos.PayStatus_Closed
	}
	return protos.PayStatus_Paying
}
