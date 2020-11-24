package transaction

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	rpcLog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	fuyouUtils "dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/fuyou"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/structtomap"
)

func (svc *Service) PrePay(ctx context.Context, req *protos.PrePayRequest) (reply *protos.PrePayReply, err error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger())
	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}
	logger = logger.WithField("requestData", req)

	logger.Info("requested prepay")

	err = svc.validatePrePayParameters(req)
	if err != nil {

		logger.WithError(err).Info("validate parameter error")

		return &protos.PrePayReply{
			Err: svc.convertError(err),
		}, nil
	}

	transaction, err := svc.tryToFindTransaction(logger, req)
	if err != nil {

		logger.WithError(err).Error("try to find transaction error")
		return &protos.PrePayReply{
			Err: svc.convertError(err),
		}, nil
	}

	if transaction != nil {

		return &protos.PrePayReply{
			TransactionId:       transaction.TransactionId,
			MobilePaymentString: transaction.PayInfo,
		}, nil
	}

	requestData := svc.buildPrePayRequestData(req)

	logger = logger.WithField("callData", requestData)

	logger.Debugf("callData: %#v", requestData)

	responseData, err := svc.sendPrePayRequest(logger, requestData)
	if err != nil {

		logger.WithError(err).
			Error("call prepay request error")
		return &protos.PrePayReply{
			Err: svc.convertError(err),
		}, nil
	}

	logger = logger.WithField("responseData", responseData)

	reply = svc.parsePrePayResponse2Reply(logger, responseData)
	if reply.GetErr() != nil {
		logger.WithField("reply", reply).Info("parse prepay response to reply error")
		return reply, nil
	}

	transactionData := svc.buildTransactionData(req)
	transactionData.ChannelMerchantOrderId = requestData.ChannelMerchantOrderId

	reply.TransactionId = transactionData.TransactionId

	err = svc.parsePrePayResponseDataToTransactionData(transactionData, responseData)
	if err != nil {

		logger.WithError(err).
			Error("parse prepay response to transaction data error")
		return &protos.PrePayReply{
			Err: svc.convertError(err),
		}, nil
	}

	logger = logger.WithField("transactionData", transactionData)

	err = svc.savePrePayTransaction(transactionData)
	if err != nil {

		logger.WithError(err).
			Error("save transaction data error")
		return &protos.PrePayReply{
			Err: svc.convertError(err),
		}, nil
	}

	return reply, nil
}

func (svc *Service) validatePrePayParameters(req *protos.PrePayRequest) error {

	const descriptionMaxLength = 128 // 富友支持最长的商品描述（字节）
	return validation.ValidateStruct(req,
		validation.Field(&req.MerchantId, validation.Required.ErrorObject(
			validation.NewError(strconv.Itoa(constant.ErrorCodeMerchantIdRequired), constant.ErrorMessageMerchantIdRequired),
		)),
		validation.Field(&req.FuyouMerchantId, validation.Required.ErrorObject(
			validation.NewError(strconv.Itoa(constant.ErrorCodeFuYouMerchantIdRequired), constant.ErrorMessageFuYouMerchantIdRequired),
		)),
		validation.Field(&req.AppId, validation.Required.ErrorObject(
			validation.NewError(strconv.Itoa(constant.ErrorCodeAppIdRequired), constant.ErrorMessageAppIdRequired),
		)),
		validation.Field(&req.OpenId, validation.Required.ErrorObject(
			validation.NewError(strconv.Itoa(constant.ErrorCodeOpenIdRequired), constant.ErrorMessageOpenIdRequired),
		)),
		validation.Field(&req.PayType,
			validation.Required.ErrorObject(
				validation.NewError(strconv.Itoa(constant.ErrorCodePayTypeRequired), constant.ErrorMessagePayTypeRequired),
			),
			validation.In(protos.PaymentType_JSAPI, protos.PaymentType_FWC, protos.PaymentType_LETPAY).ErrorObject(
				validation.NewError(strconv.Itoa(constant.ErrorCodePayTypeNotSupportOption), constant.ErrorMessagePayTypeNotSupportOption),
			),
		),
		validation.Field(&req.Description,
			validation.Required.ErrorObject(
				validation.NewError(strconv.Itoa(constant.ErrorCodeDescriptionRequired), constant.ErrorMessageDescriptionRequired),
			),
			validation.Length(0, descriptionMaxLength).ErrorObject(
				validation.NewError(strconv.Itoa(constant.ErrorCodeDescriptionTooLarge), constant.ErrorMessageDescriptionTooLarge),
			),
		),
	)
}

func (svc *Service) buildPrePayRequestData(req *protos.PrePayRequest) (requestData *parameters.PrePayRequest) {

	return &parameters.PrePayRequest{
		Version:                constant.FuYouAPIVersion,
		OrganizationCode:       config.Config.FuYou.GetOrganizationCode(),
		FuYouMerchantCode:      req.GetFuyouMerchantId(),
		TerminalId:             constant.FuYouTerminalId,
		RandomString:           strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
		GoodsDescription:       req.GetDescription(),
		ChannelMerchantOrderId: fuyouUtils.GeneralFuYouOrderId(),
		CurrencyType:           constant.FuYouCurrencyType,
		OrderAmount:            strconv.FormatUint(req.GetOrderAmount(), 10),
		UserIP:                 req.GetUserIP(),
		TransactionBeginTime:   time.Now().Format("20060102150405"),
		NotifyURL:              svc.getCallbackURL(),
		PaymentType:            svc.convertPayTypeString(req.GetPayType()),
		SubOpenId:              req.GetOpenId(),
		SubAppId:               req.GetAppId(),
		ReservedExpireMinute:   svc.getExpireMinutes(req),
	}
}

func (svc *Service) getExpireMinutes(req *protos.PrePayRequest) int {

	if req.GetSourceType() == protos.SourceType_CardService {
		return constant.TransactionExpireMinuteCardService
	}

	return constant.TransactionExpireMinuteOrderService
}

func (svc *Service) getCallbackURL() string {

	return fmt.Sprintf("https://%s/api/paymentgateway/fuyou/callback", config.Config.GetCallbackHost())
}

func (svc *Service) sendPrePayRequest(logger *logrus.Entry, requestData *parameters.PrePayRequest) (map[string]string, error) {

	data, err := structtomap.StructToMap(requestData)
	if err != nil {
		return nil, err
	}

	logger.WithField("data", data).Debug("struct to map")

	responseData, err := fuyou.NewRequest(logger).SendRequest(data, svc.getPrePayRequestURL())
	if err != nil {
		return nil, err
	}

	return responseData, err
}

func (svc *Service) getPrePayRequestURL() string {
	return fmt.Sprintf("https://%s/wxPreCreate", config.Config.FuYou.GetHost())
}

func (svc *Service) parsePrePayResponse2Reply(logger *logrus.Entry, res map[string]string) (reply *protos.PrePayReply) {

	resultCode := res["result_code"]
	if resultCode == constant.FuYouSuccessCode {
		reply = &protos.PrePayReply{
			MobilePaymentString: res["reserved_pay_info"],
		}
	} else {
		errorCode, convertError := strconv.ParseInt(resultCode, 10, 64)
		errorMessage := res["result_msg"]
		if convertError != nil {
			logger.WithError(convertError).WithField("result_code", resultCode).Error("convert result code error")
			errorCode = constant.ErrorCodeConvertErrorCodeFailed
			errorMessage = fmt.Sprintf("code: %s, message: %s", res["result_code"], res["result_msg"])
		}
		reply = &protos.PrePayReply{
			Err: &protos.Error{
				Code:    errorCode,
				Message: errorMessage,
			},
		}
	}

	return
}

func (svc *Service) buildTransactionData(req *protos.PrePayRequest) *model.Transaction {

	transactionId := idcreator.NextID()
	return &model.Transaction{
		TransactionId:         transactionId,
		Forward:               model.ForwardInflow,
		PayOrderStatus:        model.PayOrderStatusPaying,
		PayType:               svc.convertPayTypeModel(req.GetPayType()),
		OrderType:             model.OrderTypeWeChat,
		AppId:                 req.GetAppId(),
		ChannelMerchantId:     req.GetFuyouMerchantId(),
		MerchantId:            req.GetMerchantId(),
		OpenId:                req.GetOpenId(),
		Amount:                int64(req.GetOrderAmount()),
		Currency:              constant.FuYouCurrencyType,
		ClientIP:              req.GetUserIP(),
		RefundRemainingAmount: int64(req.GetOrderAmount()),
		ServiceCharge:         0,
		TransactionSources: []model.TransactionSource{
			{
				SourceId:      req.GetSourceId(),
				SourceType:    uint8(req.GetSourceType()),
				TransactionId: transactionId,
			},
		},
		Time: databases.Time{},
	}
}

func (svc *Service) parsePrePayResponseDataToTransactionData(transactionData *model.Transaction, responseData map[string]string) (err error) {

	transactionData.ErrorCode = responseData["result_code"]
	transactionData.ErrorMessage = responseData["result_msg"]
	transactionData.TracingId = responseData["reserved_fy_trace_no"]
	transactionData.PayInfo = responseData["reserved_pay_info"]
	return
}

func (svc *Service) savePrePayTransaction(transactionData *model.Transaction) (err error) {

	err = database.GetDB(constant.DatabaseConfigKey).Create(transactionData).Error
	return
}

func (svc *Service) tryToFindTransaction(logger *logrus.Entry, req *protos.PrePayRequest) (*model.Transaction, error) {

	transactions, err := payment.FindBySourceId(req.GetSourceId(), svc.convertProtobufSourceTypeToModel(req.GetSourceType()))
	if err != nil {

		logger.WithError(err).Error("find database error")

		return nil, err
	}

	if len(transactions) == 0 {
		return nil, nil
	}

	for _, transaction := range transactions {

		if transaction.Forward != model.ForwardInflow {
			continue
		}

		return transaction, nil
	}

	return nil, nil
}
