package order

import (
	pb "charites/proto"
	"context"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
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
	// return CreateOrder(ctx, req)

	err := CreateOrderWithRoctetMQ(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func OrderTimeoutHandle(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// TODO 处理订单超时
	return consumer.ConsumeSuccess, nil
}
