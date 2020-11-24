package client

import (
	"context"

	paymentProtos "dev-gitlab.wanxingrowth.com/fanli/payment/pkg/rpc/protos"
	"google.golang.org/grpc"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

var paymentRPCService paymentProtos.PaymentControllerClient

func InitPaymentService() {

	log.GetLogger().Info("starting init payment rpc service")

	var ctx = context.Background()

	var rpcConfig, exist = config.Config.RPCServices[constant.RPCPaymentServiceConfigKey]
	if !exist {
		log.GetLogger().Error("payment rpc service configuration not exist")
		return
	}

	if rpcConfig.GetConnectTimeout() > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.TODO(), rpcConfig.GetConnectTimeout())
		defer cancel()
	}

	conn, err := grpc.DialContext(ctx, rpcConfig.GetAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.GetLogger().WithError(err).Error("payment rpc service connect failed")
		return
	}

	paymentRPCService = paymentProtos.NewPaymentControllerClient(conn)
	log.GetLogger().Info("payment rpc service init succeed")
}

func GetPaymentService() paymentProtos.PaymentControllerClient {

	return paymentRPCService
}

func SetPaymentMockService(mock paymentProtos.PaymentControllerClient) {

	paymentRPCService = mock
}
