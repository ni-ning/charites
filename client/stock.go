package main

import (
	"charites/middleware"
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func main() {
	conn, _ := grpc.Dial("consul://127.0.0.1:8500/shopping?healthy=true",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.ClientUnaryInterceptor),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 指定round_robin策略

	)
	defer conn.Close()

	client := pb.NewStockClient(conn)

	resp, err := client.ReduceStock(context.Background(), &pb.GoodsStockInfo{GoodsId: 1, Num: 1})
	if err != nil {
		sts := errcode.FromError(err)
		details := sts.Details()
		log.Printf("client.ReduceStock Error details: %v\n", details)
		return
	}
	fmt.Printf("【resp】GoodsId:%d, Num:%d\n", resp.GoodsId, resp.Num)
}
