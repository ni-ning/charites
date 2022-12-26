package stock

import (
	pb "charites/proto"
	"context"
)

func SetStock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	return &pb.GoodsStockInfo{GoodsId: 100, Num: 100}, nil
}

func GetStock(ctx context.Context, goodsId int64) (*pb.GoodsStockInfo, error) {
	return &pb.GoodsStockInfo{GoodsId: 200, Num: 200}, nil
}

func ReduceStock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {

	return &pb.GoodsStockInfo{GoodsId: 300, Num: 300}, nil
	// global.DBEngine.Transaction(func(tx *gorm.DB) error {
	// 	var goods model.Goods
	// 	err := tx.WithContext(ctx).Model(&model.Goods{}).Where("id = ?", 1).
	// 		First(&goods).Error
	// 	fmt.Print(err)
	// 	return nil
	// })
}
