package stock

import (
	"charites/global"
	"charites/model"
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
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

	// return ReduceStockWithTransaction(ctx, req.GoodsId, req.Num)
	// return ReduceStockWithOptimisticLock(ctx, req.GoodsId, req.Num)
	// return ReduceStockWithPessimisticLock(ctx, req.GoodsId, req.Num)
	return ReduceStockWithDistributedLock(ctx, req)

}

// RollbackMsgHandle 监听RocketMQ消息进行库存回滚的处理函数
// 需要考虑重复归还的问题(幂等性) --> 添加库存扣减记录
func RollbackMsgHandle(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var data model.OrderGoodsStockInfo
		err := json.Unmarshal(msgs[i].Body, &data)
		if err != nil {
			global.Logger.Error("json.Unmarshal RollbackMsg error", zap.Error(err))
			continue
		}
		err = RollbackStockByMsg(ctx, &data)
		if err != nil {
			return consumer.ConsumeRetryLater, nil
		}
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
