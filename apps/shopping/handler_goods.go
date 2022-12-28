package shopping

import (
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"
)

type GoodsServer struct {
	pb.UnimplementedGoodsServer
}

func NewGoodsServer() *GoodsServer {
	return &GoodsServer{}
}

func (g GoodsServer) GetGoodsListByRoomId(ctx context.Context, req *pb.GetGoodsListRoomReq) (*pb.GoodsListReply, error) {
	if req.RoomId <= 0 {
		return nil, errcode.ToRPCError(errcode.InvalidParams)
	}
	return GetGoodsListByRoomId(ctx, req.RoomId)
}

func (g GoodsServer) GetGoodsDetail(ctx context.Context, req *pb.GetGoodsDetailReq) (*pb.GoodsDetailReply, error) {
	if req.GoodsId <= 0 {
		return nil, errcode.ToRPCError(errcode.InvalidParams)
	}
	detail := &pb.GoodsDetailReply{}
	detail.GoodsId = req.GoodsId
	detail.Price = "9.99"
	return detail, nil
}
