package main

import (
	"charites/middleware"
	pb "charites/proto"
	"context"
	"fmt"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func main() {
	// 建立连接 with grpc.DialOption
	conn, err := grpc.Dial("consul://127.0.0.1:8500/shopping?healthy=true",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.ClientUnaryInterceptor),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	// 判断连接 err 与 defer 关闭连接
	if err != nil {
		log.Fatalln("grpc.Dial err:", err)
	}
	defer conn.Close()

	// 获取操作gRPC服务端服务的client
	client := pb.NewStockClient(conn)

	// 客户端业务逻辑处理，如并发20次操作服务端服务
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.ReduceStock(context.Background(), &pb.GoodsStockInfo{GoodsId: 1, Num: 1})
			if err != nil {
				log.Printf("client.ReduceStock Error: %v\n", err)
				return
			}
			fmt.Printf("resp GoodsId:%d, Num:%d\n", resp.GoodsId, resp.Num)
		}()
	}
	wg.Wait()
}
