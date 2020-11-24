package callback

import (
	"fmt"
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/notify"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/payment"
	fuyouUtils "dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/fuyou"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/structtomap"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/transform"
)

// @ID Notify
// @Summary fuyou callback pay notifications
// @Description fuyou payment service will send notification when user paid order
// @Tags callback
// @Accept application/x-www-form-urlencoded
// @Param req formData string true "notification message body"
// @Success 200 {string} int	"result"
// @Router /api/paymentgateway/fuyou/callback [post]
func (controller *Controller) Notify(c *gin.Context) {

	requestString := c.PostForm("req")
	logger := log.GetLogger().WithField("req", requestString)
	req, err := transform.Decode(requestString, "gbk")
	if err != nil {
		logger.WithError(err).Error("decode fuyou pay notify error")
		return
	}

	logger = logger.WithField("req", req)

	var callbackRecord = &model.Callback{}
	logger.Info("receive fy pay notify")

	err = transform.XmlGbk.Bind([]byte(req), callbackRecord)
	if err != nil {
		logger.WithError(err).Error("transform fuyou pay notify failed")
		return
	}

	err = controller.handler(logger, callbackRecord)
	if err != nil {
		logger.WithError(err).Error("handle fy pay notify failed")
		return
	}
	response.Response(c, 1)
}

func (controller *Controller) handler(logger *logrus.Entry, callback *model.Callback) (err error) {

	callbackData, err := structtomap.StructToMap(callback)
	if err != nil {

		return
	}

	logger = logger.WithField("channelMerchantOrderId", callback.ChannelMerchantOrderId)

	if err = fuyouUtils.FYVerify(callbackData, callback.Signature); err != nil {
		logger.WithError(err).Error("verify signature failed")
		return fmt.Errorf("verify signature failed")
	}

	// 校验是否已处理或不存在相应支付单
	transaction, err := payment.FindByChannelMerchantOrderId(callback.ChannelMerchantOrderId)
	if err != nil {
		logger.WithError(err).Errorf("find transaction error")
		return
	}

	if transaction == nil {
		logger.Errorf("transaction record not found")
		return fmt.Errorf("transaction record not found")
	}

	if transaction.NotifyStatus == model.NotifyStatusSuccess {
		logger.Errorf("transaction was processed")
		return
	}

	var payOrderStatus uint8
	if callback.ResultCode == constant.FuYouSuccessCode {
		t, err := time.ParseInLocation("20060102150405", callback.TransactionFinishTime, time.Local)
		if err != nil {
			logger.WithField("transactionFinishTime", callback.TransactionFinishTime).Error("parse paid time error")
			return err
		}
		transaction.TradeTime = t
		payOrderStatus = model.PayOrderStatusSuccess
		err = notify.NewNotify(logger).Notify(transaction)
		if err != nil {
			logger.WithError(err).Error("notify service payStatus failed")
		}
	} else {
		payOrderStatus = model.PayOrderStatusFailed
	}

	// 更新支付单
	if err = payment.Update(&model.Transaction{
		TransactionId:  transaction.TransactionId,
		TradeTime:      transaction.TradeTime,
		PayOrderStatus: payOrderStatus,
		TracingId:      callback.ReservedFuYouTracingId,
		NotifyStatus:   model.NotifyStatusSuccess,
		ServiceCharge:  uint32(transaction.Amount - callback.SettlementOrderAmount),
	}); err != nil {
		return
	}

	return
}
