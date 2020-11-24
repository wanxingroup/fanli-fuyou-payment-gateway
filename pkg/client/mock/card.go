package mock

import (
	"dev-gitlab.wanxingrowth.com/fanli/card/pkg/rpc/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type CardOrderFinishOrder struct {
}

func (_ CardOrderFinishOrder) CreateOrder(ctx context.Context, in *protos.CreateOrderRequest, opts ...grpc.CallOption) (*protos.CreateOrderReply, error) {

	return nil, nil
}

func (_ CardOrderFinishOrder) FinishOrder(ctx context.Context, in *protos.FinishOrderRequest, opts ...grpc.CallOption) (*protos.FinishOrderReply, error) {

	return &protos.FinishOrderReply{}, nil
}
func (_ CardOrderFinishOrder) GetOrder(ctx context.Context, in *protos.GetOrderRequest, opts ...grpc.CallOption) (*protos.GetOrderReply, error) {
	return nil, nil
}

type CardFinishOrderResultFalse struct {
}

func (_ CardFinishOrderResultFalse) CreateOrder(ctx context.Context, in *protos.CreateOrderRequest, opts ...grpc.CallOption) (*protos.CreateOrderReply, error) {

	return nil, nil
}

func (_ CardFinishOrderResultFalse) FinishOrder(ctx context.Context, in *protos.FinishOrderRequest, opts ...grpc.CallOption) (*protos.FinishOrderReply, error) {

	return &protos.FinishOrderReply{Err: &protos.Error{
		Code:    20000,
		Message: "test failed",
	}}, nil
}

func (_ CardFinishOrderResultFalse) GetOrder(ctx context.Context, in *protos.GetOrderRequest, opts ...grpc.CallOption) (*protos.GetOrderReply, error) {
	return nil, nil
}

type CardReplyNil struct {
}

func (_ CardReplyNil) CreateOrder(ctx context.Context, in *protos.CreateOrderRequest, opts ...grpc.CallOption) (*protos.CreateOrderReply, error) {

	return nil, nil
}

func (_ CardReplyNil) FinishOrder(ctx context.Context, in *protos.FinishOrderRequest, opts ...grpc.CallOption) (*protos.FinishOrderReply, error) {

	return nil, nil
}

func (_ CardReplyNil) GetOrder(ctx context.Context, in *protos.GetOrderRequest, opts ...grpc.CallOption) (*protos.GetOrderReply, error) {
	return nil, nil
}
