package order

import (
	"charites/global"
	"charites/model"
	"charites/pkg/errcode"
	"charites/pkg/utils"
	pb "charites/proto"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
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

// OrderEntity 自定义结构体，实现两个方法
// 发送事务消息的时候，RocketMQ会根据情况自动调用这两个方法
type OrderEntity struct {
	OrderId int64
	Param   *pb.OrderReq
	Err     error
}

// 当发送prepare(half) message 成功后，这个方法(执行本地事务)就会被执行
func (o *OrderEntity) ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState {
	if o.Param == nil {
		o.Err = errcode.ToRPCError(errcode.ErrorOrderEntityParam)
		return primitive.RollbackMessageState
	}
	param := o.Param
	ctx := context.Background()

	// 请求商品微服务
	goodsDetail, err := global.GoodsCli.GetGoodsDetail(context.Background(), &pb.GetGoodsDetailReq{GoodsId: param.GoodsId})
	if err != nil {
		o.Err = errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
		return primitive.RollbackMessageState
	}
	// 拿到商品价格作为支付价格
	price, _ := strconv.ParseInt(goodsDetail.Price, 10, 64)

	// 请求库存微服务，扣减库存
	_, err = global.StockCli.ReduceStock(context.Background(), &pb.GoodsStockInfo{GoodsId: param.GoodsId, Num: param.Num})
	if err != nil {
		o.Err = errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
		return primitive.RollbackMessageState
	}

	// 创建订单与订单详情
	orderData := model.Order{
		UserId:         param.UserId,
		OrderId:        o.OrderId, // 雪花算法生成
		TradeId:        fmt.Sprintf("%d", o.OrderId),
		Status:         int64(100), // 创建订单初始状态
		ReceiveAddress: param.Address,
		ReceiveName:    param.Name,
		ReceivePhone:   param.Phone,
		PayAmount:      price * param.Num, // 该订单总价
	}

	marketPrice, _ := strconv.ParseInt(goodsDetail.MarketPrice, 10, 64)
	orderDetail := model.OrderDetail{
		UserId:    param.UserId,
		OrderId:   o.OrderId, // 雪花算法生成
		GoodsId:   param.GoodsId,
		Num:       param.Num,
		PayAmount: price * param.Num, // 该商品总价

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
		// 本地事务执行失败，但上一步库存已经扣减成功
		// o.Err = err 不要这个操作
		return primitive.CommitMessageState
	}

	// 说明本地事务执行成功，需要把之前的 half-message丢弃
	return primitive.RollbackMessageState
}

// 当发送prepare(half) message 没有响应时，broker会回查本地事务状态，此时这个方法被执行
func (o *OrderEntity) CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

// CreateOrderWithRoctetMQ 创建订单 RoctetMQ 实现分布式事务
func CreateOrderWithRoctetMQ(ctx context.Context, param *pb.OrderReq) error {
	// 生成订单号
	orderId := utils.GenInt64()

	orderEntity := &OrderEntity{
		OrderId: orderId,
		Param:   param,
	}
	p, err := rocketmq.NewTransactionProducer(
		orderEntity,
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.1.4:9876"})),
		producer.WithRetry(2),
		producer.WithGroupName("order_srv_1"), // 生产者组
	)
	if err != nil {
		global.Logger.Error("ErrorNewTransactionProducer", zap.Error(err))
		return errcode.ToRPCError(errcode.ErrorNewTransactionProducer)
	}

	// 消息中具体的载荷，定义为一个结构体，赞
	data := model.OrderGoodsStockInfo{
		OrderId: orderId,
		GoodsId: param.GoodsId,
		Num:     param.Num,
	}
	b, _ := json.Marshal(data)
	// 定义RocketMQ消息体
	msg := &primitive.Message{
		Topic: "stock_rollback", // 回滚库存，可定义为 Conf
		Body:  b,
	}

	// 发送事务消息
	res, err := p.SendMessageInTransaction(context.Background(), msg)
	fmt.Print(res)
	return nil
}
