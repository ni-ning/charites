package order

import (
	"charites/global"
	"charites/model"
	"charites/pkg/errcode"
	"charites/pkg/utils"
	pb "charites/proto"
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// CreateOrder 创建订单
func CreateOrder(ctx context.Context, req *pb.OrderReq) (*emptypb.Empty, error) {
	// 生成订单号
	orderId := utils.GenInt64()

	// 请求商品微服务
	goodsDetail, err := global.GoodsCli.GetGoodsDetail(context.Background(), &pb.GetGoodsDetailReq{GoodsId: req.GoodsId})
	if err != nil {
		return nil, errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
	}
	// 拿到商品价格作为支付价格
	price, _ := strconv.ParseInt(goodsDetail.Price, 10, 64)

	// 请求库存微服务，扣减库存
	_, err = global.StockCli.ReduceStock(context.Background(), &pb.GoodsStockInfo{GoodsId: req.GoodsId, Num: req.Num})
	if err != nil {
		return nil, errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
	}

	// 创建订单与订单详情
	orderData := model.Order{
		UserId:         req.UserId,
		OrderId:        orderId, // 雪花算法生成
		TradeId:        fmt.Sprintf("%d", orderId),
		Status:         int64(100), // 创建订单初始状态
		ReceiveAddress: req.Address,
		ReceiveName:    req.Name,
		ReceivePhone:   req.Phone,
		PayAmount:      price * req.Num, // 该订单总价
	}

	marketPrice, _ := strconv.ParseInt(goodsDetail.MarketPrice, 10, 64)
	orderDetail := model.OrderDetail{
		UserId:    req.UserId,
		OrderId:   orderId, // 雪花算法生成
		GoodsId:   req.GoodsId,
		Num:       req.Num,
		PayAmount: price * req.Num, // 该商品总价

		Title:       goodsDetail.Title,
		MarketPrice: marketPrice,
		Price:       price,
		Brief:       goodsDetail.Brief,
		HeadImgs:    strings.Join(goodsDetail.HeadImgs, ","),
		Videos:      strings.Join(goodsDetail.Videos, ","),
		Detail:      strings.Join(goodsDetail.Detail, ","),
	}
	err = global.DBEngine.Transaction(func(tx *gorm.DB) error {
		orderResult := tx.WithContext(ctx).Create(&orderData)
		if orderResult.Error != nil {
			return errcode.ToRPCError(errcode.ErrorCreateOrder)
		}
		orderDetailResult := tx.WithContext(ctx).Create(&orderDetail)
		if orderDetailResult.Error != nil {
			return errcode.ToRPCError(errcode.ErrorCreateOrderDetal)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
