package main

import (
	_ "charites/bootstrap"
	"charites/global"
	pb "charites/proto"
	"context"
	"fmt"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func main() {
	// _ "charites/bootstrap" 需要初始化

	// 模拟 订单微服务8081 调用商品微服务和库存微服务8090
	resp, err := global.GoodsCli.GetGoodsDetail(context.Background(), &pb.GetGoodsDetailReq{GoodsId: 1})
	if err != nil {
		log.Printf("global.GoodsCli.GetGoodsDetail err: %v\n", err)
		return
	}
	fmt.Printf("global.GoodsCli.GetGoodsDetail: %#v\n", resp)

	resp2, err := global.StockCli.ReduceStock(context.Background(), &pb.GoodsStockInfo{GoodsId: 1, Num: 1})
	if err != nil {
		log.Printf("client.ReduceStock Error: %v\n", err)
		return
	}
	fmt.Printf("resp GoodsId:%d, Num:%d\n", resp2.GoodsId, resp2.Num)
}
