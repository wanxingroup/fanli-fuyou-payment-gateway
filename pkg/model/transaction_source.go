package model

const TableNameTransactionSource = "transaction_source"

// 业务流水与来源关联
type TransactionSource struct {
	SourceId      uint64 `gorm:"column:sourceId;type:bigint unsigned;primary_key;not null;comment:'我方使用渠道订单 ID'"`
	SourceType    uint8  `gorm:"column:sourceType;type:tinyint unsigned;primary_key;not null;comment:'我方使用渠道类型'"`
	TransactionId uint64 `gorm:"column:transactionId;type:bigint unsigned;primary_key;comment:'主键ID'"`
}

func (_ TransactionSource) TableName() string {

	return TableNameTransactionSource
}

const (

	// 使用渠道类型
	SourceTypeCardService    uint8 = 1 // 权益卡
	SourceTypeOrderService   uint8 = 2 // 订单中心
	SourceTypePaymentService uint8 = 3 //收款服务

)
