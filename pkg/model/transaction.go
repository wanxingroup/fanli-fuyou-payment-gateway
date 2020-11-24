package model

import (
	"time"

	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameTransaction = "transaction"

// 业务流水记录
type Transaction struct {
	TransactionId          uint64              `gorm:"column:transactionId;type:bigint unsigned;primary_key;comment:'主键ID'"`
	Forward                uint8               `gorm:"column:payOrderType;type:tinyint unsigned;not null;default:1;comment:'现金流向，1：流入 2：流出'"`
	PayOrderStatus         uint8               `gorm:"column:payOrderStatus;type:tinyint unsigned;not null;default:1;comment:'支付状态： 1 => 支付中；2 => 成功； => 失败；4 => 已关闭；8 => 部分退款；9 => 全部退款'"`
	PayType                uint8               `gorm:"column:payType;type:tinyint unsigned;not null;comment:'支付方式：1 => 公众号支付；2 => 支付宝服务窗；3 => 小程序'"`
	OrderType              uint8               `gorm:"column:orderType;type:tinyint unsigned;not null;comment:'订单类型：1 => ALIPAY (统一下单、条码支付、服务窗支付)；2 => WECHAT(统一下单、条码支付，公众号支付,小程序)；3 => UNIONPAY（银联）；4 => ,BESTPAY(翼支付)'"`
	AppId                  string              `gorm:"column:appId;type:varchar(32);not null;default:'';comment:'小程序或公众号的appId'"`
	ChannelMerchantId      string              `gorm:"column:channelMerchantId;type:varchar(32);not null;default:'';comment:'支付渠道商户ID'"`
	ChannelMerchantOrderId string              `gorm:"column:channelMerchantOrderId;type:char(26);index:channelMerchantOrderId;not null;default:'';comment:'富友的商户订单号'"`
	MerchantId             uint64              `gorm:"column:merchantId;type:bigint unsigned;index:idx_merchant_id;not null;default:'0';comment:'我方商家ID'"`
	OpenId                 string              `gorm:"column:openId;type:varchar(40);not null;default:'';comment:'用户第三方平台认证号'"`
	Amount                 int64               `gorm:"column:amount;type:bigint;not null;default:0;comment:'支付金额,单位：分'"`
	Currency               string              `gorm:"column:currency;type:varchar(10);not null;default:'CNY';comment:'三位货币代码,人民币:cny'"`
	ClientIP               string              `gorm:"column:clientIP;type:varchar(32);not_null;default:'';comment:'客户端ip'"`
	ErrorCode              string              `gorm:"column:errorCode;type:varchar(10);not null;default:'';comment:'第三方渠道支付错误码'"`
	ErrorMessage           string              `gorm:"column:errMsg;type:varchar(50);not null;default:'';comment:'第三方渠道支付错误描述'"`
	ExpireTime             uint16              `gorm:"column:expireTime;type:smallint unsigned;not null;default:'0';comment:'订单失效时间（分钟）'"`
	TracingId              string              `gorm:"column:tracingId;type:varchar(32);not null;default:'';comment:'渠道内部追踪号'"`
	NotifyStatus           uint8               `gorm:"column:notifyStatus;type:tinyint unsigned;not null;default:'1';comment:'通知状态：1 => 创建；2 => 成功；3 => 失败'"`
	RefundRemainingAmount  int64               `gorm:"column:refundRemainingAmount;type:bigint unsigned;not null;default:0;comment:'可退款金额，单位：分'"`
	TradeTime              time.Time           `gorm:"column:tradeTime;type:datetime;index:idx_trade_time;not null;default:CURRENT_TIMESTAMP;comment:'交易时间'"`
	Channel                string              `gorm:"column:channel;type:varchar(16);not null;default:'fuiou';comment:'渠道'"`
	ServiceCharge          uint32              `gorm:"column:serviceCharge;type:int unsigned;not null;default:0;comment:'商户手续费，单位：分'"`
	PayInfo                string              `gorm:"column:payInfo;type:varchar(4000);not null;default:'';comment:'创建支付时的支付信息'"`
	TransactionSources     []TransactionSource `gorm:"foreignkey:TransactionId"`
	databases.Time
}

func (_ Transaction) TableName() string {

	return TableNameTransaction
}

const (
	// 现金流向
	ForwardInflow  = 1
	ForwardOutflow = 2

	// 支付状态 1:支付中,2成功,3:失败,4:已关闭,8:部分退款,9:全部退款
	PayOrderStatusPaying  uint8 = 1
	PayOrderStatusSuccess uint8 = 2
	PayOrderStatusFailed  uint8 = 3
	PayOrderStatusClose   uint8 = 4

	// 支付方式 JSAPI-公众号支付、FWC--支付宝服务窗、LETPAY-小程序
	PayTypeOfficialAccount uint8 = 1
	PayTypeAliPay          uint8 = 2
	PayTypeMiniProgram     uint8 = 3

	// 订单类型:ALIPAY (统一下单、条码支付、服务窗支付), WECHAT(统一下单、条码支付，公众号支付,小程序),UNIONPAY,BESTPAY(翼支付)
	OrderTypeUnknown  uint8 = 0 // 未识别
	OrderTypeWeChat   uint8 = 1 // 微信支付
	OrderTypeAliPay   uint8 = 2 // 支付宝
	OrderTypeUnionPay uint8 = 3 // 银联
	OrderTypeBestPay  uint8 = 4 // 翼支付

	// 通知状态 1:创建 2成功 3失败
	NotifyStatusCreate  uint8 = 1
	NotifyStatusSuccess uint8 = 2
	NotifyStatusFailed  uint8 = 3
)
