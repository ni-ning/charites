package stock

import (
	"charites/global"
	"charites/model"
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"

	"go.uber.org/zap"
)

func SetStock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	return &pb.GoodsStockInfo{GoodsId: 100, Num: 100}, nil
}

func GetStock(ctx context.Context, goodsId int64) (*pb.GoodsStockInfo, error) {
	return &pb.GoodsStockInfo{GoodsId: 200, Num: 200}, nil
}

func ReduceStock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	var stock model.Stock
	// 1. 查询现有库存
	db := global.DBEngine.WithContext(ctx).
		Model(&model.Stock{}).
		Where("id = ?", goodsId).
		First(&stock)
	// 不存在也会抛异常
	if db.Error != nil {
		global.Logger.Error("ErrorDBOperateStock", zap.String("error", db.Error.Error()))
		return nil, errcode.ToRPCError(errcode.ErrorDBOperateStock)
	}
	if db.RowsAffected == 0 {
		return nil, errcode.ToRPCError(errcode.ErrorNotFoundStock)
	}
	// 2. 校验库存
	if stock.Num-num < 0 {
		return nil, errcode.ErrorNotEnoughStock
	}
	// 3. 扣减库存并保存
	stock.Num -= num
	// global.DBEngine.WithContext(ctx).Save(&stock) // 更新所有字段
	err := global.DBEngine.WithContext(ctx).
		Model(&model.Stock{}).
		Where("id = ?", goodsId).
		Updates(map[string]interface{}{
			"num": stock.Num,
		}).Error
	if err != nil {
		global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
		return nil, errcode.ToRPCError(errcode.ErrorDBOperateStock)
	}
	return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
}
