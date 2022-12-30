package stock

import (
	"charites/global"
	"charites/model"
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SetStock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	return &pb.GoodsStockInfo{GoodsId: 100, Num: 100}, nil
}

func GetStock(ctx context.Context, goodsId int64) (*pb.GoodsStockInfo, error) {
	return &pb.GoodsStockInfo{GoodsId: 200, Num: 200}, nil
}

func ReduceStockWithTransaction(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	var stock model.Stock
	err := global.DBEngine.Transaction(func(tx *gorm.DB) error {
		// 1. 查询现有库存
		db := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			First(&stock)
		// 不存在也会抛异常
		if db.Error != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", db.Error.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if db.RowsAffected == 0 {
			return errcode.ToRPCError(errcode.ErrorNotFoundStock)
		}
		// 2. 校验库存
		if stock.Num-num < 0 {
			return errcode.ErrorNotEnoughStock
		}
		// 3. 扣减库存并保存
		stock.Num -= num
		// global.DBEngine.WithContext(ctx).Save(&stock) // 更新所有字段
		err := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			Updates(map[string]interface{}{
				"num": stock.Num,
			}).Error
		if err != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		// return nil 提交事务，任何类型err回滚事务
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
}

func ReduceStockWithPessimisticLock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	var stock model.Stock
	err := global.DBEngine.Transaction(func(tx *gorm.DB) error {
		// 1. 查询现有库存
		db := tx.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			First(&stock)
		// 不存在也会抛异常
		if db.Error != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", db.Error.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if db.RowsAffected == 0 {
			return errcode.ToRPCError(errcode.ErrorNotFoundStock)
		}
		// 2. 校验库存
		if stock.Num-num < 0 {
			return errcode.ErrorNotEnoughStock
		}
		// 3. 扣减库存并保存
		stock.Num -= num
		// global.DBEngine.WithContext(ctx).Save(&stock) // 更新所有字段
		err := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			Updates(map[string]interface{}{
				"num": stock.Num,
			}).Error
		if err != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		// return nil 提交事务，任何类型err回滚事务
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
}

func ReduceStockWithOptimisticLock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	for retry := 0; retry < 20; retry++ {
		var stock model.Stock
		// 1. 查询现有库存
		db := global.DBEngine.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			First(&stock)
		// 不存在也会抛异常
		if db.Error != nil {
			global.Logger.Error("model.Stock.First", zap.String("error", db.Error.Error()))
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
		ret := global.DBEngine.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ? and version = ?", goodsId, stock.Version).
			Updates(map[string]interface{}{
				"num":     stock.Num - 1,
				"version": stock.Version + 1,
			})
		if ret.Error != nil {
			global.Logger.Error("model.Stock.Updates", zap.String("error", ret.Error.Error()))
			return nil, errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if ret.RowsAffected == 0 {
			// 说明 version 被更新，重试即可
			continue
		}
		return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
	}
	return nil, errcode.ToRPCError(errcode.ErrorNeedRetryStock)
}

func ReduceStockWithDistributedLock(ctx context.Context, param *pb.GoodsStockInfo) (*pb.GoodsStockInfo, error) {
	mutexname := fmt.Sprintf("reduce:stock:mutex:%d", param.GoodsId)
	mutex := global.Redsync.NewMutex(mutexname)
	if err := mutex.Lock(); err != nil {
		return nil, errcode.ToRPCError(errcode.ErrorRedisLockStock)
	}
	defer mutex.Unlock()

	var stock model.Stock
	err := global.DBEngine.Transaction(func(tx *gorm.DB) error {
		// 1. 查询现有库存
		db := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", param.GoodsId).
			First(&stock)
		// 不存在也会抛异常
		if db.Error != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", db.Error.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if db.RowsAffected == 0 {
			return errcode.ToRPCError(errcode.ErrorNotFoundStock)
		}
		// 2. 校验库存
		if stock.Num-param.Num < 0 {
			return errcode.ErrorNotEnoughStock
		}
		// 3. 扣减库存并保存
		stock.Num -= param.Num
		stock.Lock += param.Num

		// global.DBEngine.WithContext(ctx).Save(&stock) // 更新所有字段
		err := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", param.GoodsId).
			Updates(map[string]interface{}{
				"num":  stock.Num,
				"lock": stock.Lock,
			}).Error
		if err != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		// 新增预扣减库存记录
		stockRecord := model.StockRecord{
			OrderId: param.OrderId,
			GoodsId: param.GoodsId,
			Num:     param.Num,
			Status:  1,
		}
		err = tx.WithContext(ctx).Save(&stockRecord).Error
		if err != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		// return nil 提交事务，任何类型err回滚事务
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.GoodsStockInfo{GoodsId: param.GoodsId, Num: stock.Num}, nil
}

// RollbackStockByMsg 回滚商品库存
func RollbackStockByMsg(ctx context.Context, data *model.OrderGoodsStockInfo) error {
	var sr model.StockRecord
	global.DBEngine.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.StockRecord{}).
			Where("order_id = ? and goods_id = ? and status = 1", data.OrderId, data.GoodsId).
			First(&sr).Error
		// 没找到记录：确实没记录或者已经回滚
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		if err != nil {
			return err
		}

		var s model.Stock
		err = tx.Model(&model.Stock{}).Where("goods_id = ?", data.GoodsId).First(&s).Error
		if err != nil {
			return err
		}
		s.Num += data.Num
		s.Lock -= data.Num
		err = tx.Model(&model.Stock{}).Save(&s).Error
		if err != nil {
			return err
		}
		// 2已回滚
		sr.Status = 2
		err = tx.Model(&model.StockRecord{}).Save(&sr).Error
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}
