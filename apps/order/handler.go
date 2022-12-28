package order

import (
	pb "charites/proto"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	pb.UnimplementedOrderServer
}

func NewOrderServer() *OrderServer {
	return &OrderServer{}
}

// CreateOrder 创建订单
func (OrderServer) CreateOrder(ctx context.Context, req *pb.OrderReq) (*emptypb.Empty, error) {
	return CreateOrder(ctx, req)
}
