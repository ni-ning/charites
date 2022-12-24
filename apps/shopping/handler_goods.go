package shopping

import (
	pb "charites/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GoodsServer struct {
	pb.UnimplementedGoodsServer
}

func NewGoodsServer() *GoodsServer {
	return &GoodsServer{}
}

func (g GoodsServer) GetGoodsListByRoomId(context.Context, *pb.GetGoodsListRoomReq) (*pb.GoodsListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGoodsListByRoomId not implemented")
}
func (g GoodsServer) GetGoodsDetail(context.Context, *pb.GetGoodsDetailReq) (*pb.GoodsDetailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGoodsDetail not implemented")
}
