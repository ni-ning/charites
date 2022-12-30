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

// OrderListener 自定义结构体，实现两个方法
// 发送事务消息的时候，RocketMQ会根据情况自动调用这两个方法
type OrderListener struct {
	OrderId int64
	Param   *pb.OrderReq
	Err     error
}

// 当发送prepare(half) message 成功后，这个方法(执行本地事务)就会被执行
func (o *OrderListener) ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState {
	if o.Param == nil {
		global.Logger.Error("ExecuteLocalTransaction param is nil")
		o.Err = errcode.ToRPCError(errcode.ErrorOrderEntityParam)
		// 库存未扣减
		return primitive.RollbackMessageState
	}
	param := o.Param
	ctx := context.Background()

	// 请求商品微服务，查询商品金额(营销相关)
	goodsDetail, err := global.GoodsCli.GetGoodsDetail(ctx, &pb.GetGoodsDetailReq{GoodsId: param.GoodsId})
	if err != nil {
		global.Logger.Error("GoodsCli.GetGoodsDetail failed", zap.Error(err))
		o.Err = errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
		// 库存未扣减
		return primitive.RollbackMessageState
	}
	// 拿到商品价格作为支付价格
	price, _ := strconv.ParseInt(goodsDetail.Price, 10, 64)

	// 请求库存微服务，扣减库存
	_, err = global.StockCli.ReduceStock(ctx, &pb.GoodsStockInfo{GoodsId: param.GoodsId, Num: param.Num, OrderId: o.OrderId})
	if err != nil {
		global.Logger.Error("StockCli.ReduceStock failed", zap.Error(err))
		o.Err = errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
		// 库存未扣减
		return primitive.RollbackMessageState
	}

	// 本地事务创建订单与订单详情
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
		UserId:      param.UserId,
		OrderId:     o.OrderId, // 雪花算法生成
		GoodsId:     param.GoodsId,
		Num:         param.Num,
		PayAmount:   price * param.Num, // 该商品总价
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
		return primitive.CommitMessageState
	}

	// 发送延迟消息
	// 1s 5s 10s...
	// 消息中具体的载荷，定义为一个结构体，赞
	data := model.OrderGoodsStockInfo{
		OrderId: o.OrderId,
		GoodsId: param.GoodsId,
		Num:     param.Num,
	}
	b, _ := json.Marshal(data)
	// 定义RocketMQ消息体
	msg := primitive.NewMessage(global.RocketMQSetting.TopicOrderPayTimeout, b)
	msg.WithDelayTimeLevel(3)
	_, err = global.Producer.SendSync(context.Background(), msg)
	if err != nil {
		// 延时消息发送失败
		return primitive.CommitMessageState
	}
	// 说明本地事务执行成功，不需要发送回滚库存的消息
	return primitive.RollbackMessageState
}

// 当发送prepare(half) message 没有响应时，broker会回查本地事务状态，此时这个方法被执行
func (o *OrderListener) CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState {
	// 检查本地是否订单创建成功即可
	var count int64
	global.DBEngine.
		WithContext(context.Background()).
		Model(&model.Order{}).Where("order_id = ?", o.OrderId).
		Count(&count)
	if count <= 0 {
		// 说明订单创建失败，需要回滚库存
		return primitive.CommitMessageState
	}
	// 不存回滚库存
	return primitive.RollbackMessageState
}

// CreateOrderWithRoctetMQ 创建订单 RoctetMQ 实现分布式事务
func CreateOrderWithRoctetMQ(ctx context.Context, param *pb.OrderReq) error {
	// 生成订单号
	orderId := utils.GenInt64()
	// 订单号+请求的参数传入到Listener中，Listener实现本地事务
	orderListener := &OrderListener{
		OrderId: orderId,
		Param:   param,
	}
	// orderListener 每个订单不同都需要创建
	p, err := rocketmq.NewTransactionProducer(
		orderListener,
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{global.RocketMQSetting.NameServer})),
		producer.WithRetry(2),
		producer.WithGroupName(global.RocketMQSetting.GroupOrderService), // 生产者组
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
		Topic: global.RocketMQSetting.TopicStockRollback, // 回滚库存，可定义为 Conf
		Body:  b,
	}
	// 发送事务消息
	res, _ := p.SendMessageInTransaction(context.Background(), msg)
	if res.State == primitive.CommitMessageState {
		// 回滚库存消息被正常投递，说明创建订单错误
		return errcode.ToRPCError(errcode.ErrorOrderCreate)
	}
	if orderListener.Err != nil {
		// 内部其他错误
		return orderListener.Err
	}
	return nil
}
