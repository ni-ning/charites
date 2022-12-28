package order

import (
	"charites/pkg/utils"
	pb "charites/proto"
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateOrder 创建订单
func CreateOrder(ctx context.Context, req *pb.OrderReq) (*emptypb.Empty, error) {
	fmt.Println("req.GoodsId", req.GoodsId)
	fmt.Println("req.Num", req.Num)
	fmt.Println("OrderId", utils.GenInt64())

	return &emptypb.Empty{}, nil
}
