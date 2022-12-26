package stock

import (
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"
)

type StockServer struct {
	pb.UnimplementedStockServer
}

func NewStockServer() *StockServer {
	return &StockServer{}
}

func (StockServer) SetStock(ctx context.Context, req *pb.GoodsStockInfo) (*pb.GoodsStockInfo, error) {
	if req.GoodsId <= 0 || req.Num <= 0 {
		return nil, errcode.ToRPCError(errcode.InvalidParams)
	}

	return SetStock(ctx, req.GoodsId, req.Num)
}

func (StockServer) GetStock(ctx context.Context, req *pb.GoodsStockInfo) (*pb.GoodsStockInfo, error) {
	if req.GoodsId <= 0 {
		return nil, errcode.ToRPCError(errcode.InvalidParams)
	}

	return GetStock(ctx, req.GoodsId)
}

func (StockServer) ReduceStock(ctx context.Context, req *pb.GoodsStockInfo) (*pb.GoodsStockInfo, error) {
	if req.GoodsId <= 0 || req.Num <= 0 {
		return nil, errcode.ToRPCError(errcode.InvalidParams)
	}

	return ReduceStock(ctx, req.GoodsId, req.Num)
}
