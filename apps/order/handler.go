package order

import (
	"charites/global"
	"charites/model"
	pb "charites/proto"
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	pb.UnimplementedOrderServer
}

func NewOrderServer() *OrderServer {
	return &OrderServer{}
}

// CreateOrder 创建订单
func (OrderServer) CreateOrder(ctx context.Context, req *pb.OrderReq) (*emptypb.Empty, error) {
	// return CreateOrder(ctx, req)

	err := CreateOrderWithRoctetMQ(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func OrderTimeoutHandle(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var data model.OrderGoodsStockInfo
		err := json.Unmarshal(msgs[i].Body, &data)
		if err != nil {
			global.Logger.Error("json.Unmarshal RollbackMsg failed", zap.Error(err))
			continue
		}

		// 查订单表
		// 1. 如果订单为已支付状态则不处理
		// 2. 如果订单为未支付状态则发送一条回滚库存的消息
		var order model.Order
		err = global.DBEngine.WithContext(ctx).
			Model(&model.Order{}).
			Where("order_id = ?", data.OrderId).
			First(&order).Error
		if err != nil {
			global.Logger.Error("mysql.QueryOrder failed", zap.Error(err))
			return consumer.ConsumeRetryLater, nil // 稍后再试
		}

		if order.OrderId == data.OrderId && order.Status == 100 { // 待支付
			msg := &primitive.Message{
				Topic: global.RocketMQSetting.TopicStockRollback,
				Body:  msgs[i].Body,
			}
			_, err = global.Producer.SendSync(context.Background(), msg)
			if err != nil {
				global.Logger.Error("send rollback msg failed", zap.Error(err))
				return consumer.ConsumeRetryLater, nil // 稍后再试
			}
			// 发送回滚库存成功，将订单状态设置为关闭
			global.DBEngine.WithContext(ctx).
				Model(&model.Order{}).
				Where("order_id = ?", data.OrderId).
				Updates(map[string]interface{}{
					"status": 300,
				})
		}
	}

	return consumer.ConsumeSuccess, nil
}
