package shopping

import (
	"charites/global"
	"charites/model"
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm/clause"
)

/*
实际架构或包含
handler -> biz -> dao

handler
  - 参数校验
  - 结构返回

biz 封装业务逻辑
  - 入参 ctx 贯穿始终
  - 出餐 error 是内置自定义error 还是 gRPC error

dao 操作数据库
 - ES
 - MySQL
 - Redis


- global.DBEngine.WithContext(ctx) 上下文
- 先定义变量，再Find(&roomGoodsList)取值
- Order("weight desc") 放到 Find(&roomGoodsList) 之前才能起作用
- 排序高级用法 FIELD(id,?)
*/

// GetGoodsListByRoomId
func GetGoodsListByRoomId(ctx context.Context, roomId int64) (*pb.GoodsListReply, error) {
	// 根据roomId获取对应关联对象列表
	var roomGoodsList []*model.RoomGoods
	db := global.DBEngine.WithContext(ctx)
	err := db.Model(&model.RoomGoods{}).Where("room_id = ?", roomId).Order("weight desc").Find(&roomGoodsList).Error
	if err != nil {
		return nil, errcode.ToRPCError(errcode.Fail)
	}
	if len(roomGoodsList) == 0 {
		return nil, errcode.ToRPCError(errcode.NotFound)
	}
	// 获取关联对象中的商品idList
	var (
		currentGoodsId int64
		goodsIdList    = make([]int64, 0, len(roomGoodsList))
	)
	for _, roomGoods := range roomGoodsList {
		goodsIdList = append(goodsIdList, roomGoods.GoodsId)
		if roomGoods.IsCurrent == 1 {
			currentGoodsId = roomGoods.GoodsId
		}
	}

	// 获取所有商品
	var goodsList []*model.Goods
	err = db.Model(&model.Goods{}).
		Where("id in ?", goodsIdList).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{goodsIdList}, WithoutParentheses: true},
		}).Find(&goodsList).Error
	if err != nil {
		return nil, errcode.ToRPCError(errcode.Fail)
	}
	if len(goodsList) == 0 {
		return nil, errcode.ToRPCError(errcode.NotFound)
	}

	// 封装返回数据
	goodsInfoList := make([]*pb.GoodsInfo, 0, len(goodsList))
	for _, goods := range goodsList {
		var headImgs []string
		json.Unmarshal([]byte(goods.HeadImgs), &headImgs)
		goodsInfoList = append(goodsInfoList, &pb.GoodsInfo{
			GoodsId:     goods.GoodsId,
			CategoryId:  goods.CategoryId,
			Status:      int32(goods.Status),
			Title:       goods.Title,
			MarketPrice: fmt.Sprintf("%.2f", float64(goods.MarketPrice/100)),
			Price:       fmt.Sprintf("%.2f", float64(goods.MarketPrice/100)),
			Brief:       goods.Brief,
			HeadImgs:    headImgs,
		})
	}
	reply := &pb.GoodsListReply{CurrentGoodsId: currentGoodsId, Data: goodsInfoList}
	return reply, nil
}
