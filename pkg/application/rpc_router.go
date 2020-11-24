package application

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher"
	"google.golang.org/grpc"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/rpc/transaction"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

func registerRPCRouter(app *launcher.Application) {

	rpcService := app.GetRPCService()
	if rpcService == nil {

		log.GetLogger().WithField("stage", "onInit").Error("get rpc service is nil")
		return
	}

	connection := rpcService.GetRPCConnection()
	if connection == nil {
		log.GetLogger().WithField("stage", "onInit").Error("get rpc connection is nil")
		return
	}

	registerTransactionRPCRouter(connection)
}

func registerTransactionRPCRouter(connection *grpc.Server) {

	transactionService := &transaction.Service{}
	protos.RegisterPayControllerServer(connection, transactionService)
}
