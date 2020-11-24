package application

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

func autoMigration() {

	db := database.GetDB(constant.DatabaseConfigKey)
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.TransactionSource{})
	db.AutoMigrate(&model.Callback{})
}
