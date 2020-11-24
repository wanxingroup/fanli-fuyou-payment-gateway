package mock

import (
	orderProtos "dev-gitlab.wanxingrowth.com/fanli/order/v2/pkg/rpc/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type OrderPaidOrder struct {
}

func (_ OrderPaidOrder) OrderDetails(ctx context.Context, in *orderProtos.OrderDetailsRequest, opts ...grpc.CallOption) (*orderProtos.OrderDetailsReply, error) {

	return nil, nil
}

func (_ OrderPaidOrder) CancelOrder(ctx context.Context, in *orderProtos.CancelOrderRequest, opts ...grpc.CallOption) (*orderProtos.CancelOrderReply, error) {

	return nil, nil
}

func (_ OrderPaidOrder) PaidOrder(ctx context.Context, in *orderProtos.PaidOrderRequest, opts ...grpc.CallOption) (*orderProtos.PaidOrderReply, error) {

	return &orderProtos.PaidOrderReply{
		Success: true,
	}, nil
}

func (_ OrderPaidOrder) GetOrderList(ctx context.Context, in *orderProtos.GetOrderListRequest, opts ...grpc.CallOption) (*orderProtos.GetOrderListReply, error) {
	return nil, nil
}

func (_ OrderPaidOrder) CreateOrder(ctx context.Context, in *orderProtos.CreateOrderRequest, opts ...grpc.CallOption) (*orderProtos.CreateOrderReply, error) {
	return nil, nil
}

func (_ OrderPaidOrder) DeliverOrder(ctx context.Context, in *orderProtos.DeliverOrderRequest, opts ...grpc.CallOption) (*orderProtos.DeliverOrderReply, error) {
	return nil, nil
}

func (_ OrderPaidOrder) ReceivedOrder(ctx context.Context, in *orderProtos.ReceivedOrderRequest, opts ...grpc.CallOption) (*orderProtos.ReceivedOrderReply, error) {
	return nil, nil
}

func (_ OrderPaidOrder) CommentOrder(ctx context.Context, in *orderProtos.CommentOrderRequest, opts ...grpc.CallOption) (*orderProtos.CommentOrderReply, error) {
	return nil, nil
}
