package shopping

import (
	"charites/global"
	"charites/model"
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"
	"encoding/json"
	"fmt"
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
	var goods model.Goods
	db := global.DBEngine.WithContext(ctx)
	err := db.Model(&model.Goods{}).Debug().Where("id = ?", req.GoodsId).First(&goods).Error
	if err != nil {
		return nil, errcode.ToRPCError(errcode.Fail)
	}

	var headImgs []string
	var Videos []string
	var Detail []string
	json.Unmarshal([]byte(goods.HeadImgs), &headImgs)
	json.Unmarshal([]byte(goods.Videos), &Videos)
	json.Unmarshal([]byte(goods.Detail), &Detail)

	detail := &pb.GoodsDetailReply{}
	detail.GoodsId = req.GoodsId
	detail.Title = goods.Title
	detail.MarketPrice = fmt.Sprintf("%d", goods.MarketPrice)
	detail.Price = fmt.Sprintf("%d", goods.Price)
	detail.Brief = goods.Brief
	detail.HeadImgs = headImgs
	detail.Videos = Videos
	detail.Detail = Detail

	fmt.Printf("detail===========> %#v\n", detail)

	return detail, nil
}
