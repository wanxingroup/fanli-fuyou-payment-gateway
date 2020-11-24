package client

import (
	"context"

	orderProtos "dev-gitlab.wanxingrowth.com/fanli/order/v2/pkg/rpc/protos"
	"google.golang.org/grpc"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

var orderRPCService orderProtos.OrderControllerClient

func InitOrderService() {

	log.GetLogger().Info("starting init order rpc service")

	var ctx = context.Background()

	var rpcConfig, exist = config.Config.RPCServices[constant.RPCOrderServiceConfigKey]
	if !exist {
		log.GetLogger().Error("order rpc service configuration not exist")
		return
	}

	if rpcConfig.GetConnectTimeout() > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.TODO(), rpcConfig.GetConnectTimeout())
		defer cancel()
	}

	conn, err := grpc.DialContext(ctx, rpcConfig.GetAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.GetLogger().WithError(err).Error("order rpc service connect failed")
		return
	}

	orderRPCService = orderProtos.NewOrderControllerClient(conn)

	log.GetLogger().Info("order rpc service init succeed")
}

func GetOrderService() orderProtos.OrderControllerClient {

	return orderRPCService
}

func SetOrderMockService(mock orderProtos.OrderControllerClient) {

	orderRPCService = mock
}
