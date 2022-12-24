package main

import (
	"charites/middleware"
	pb "charites/proto"
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func SayHello(client pb.GreeterClient) error {
	// 像调用本地函数一样
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "linda"})
	if err != nil {
		log.Printf("client.SayHello err: %s", err)
		return err
	}
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}

func SayList(client pb.GreeterClient) error {
	stream, _ := client.SayList(context.Background(), &pb.HelloRequest{Name: "linda"})
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}
	return nil
}

func SayRecord(client pb.GreeterClient) error {
	stream, _ := client.SayRecord(context.Background())
	for n := 0; n < 6; n++ {
		_ = stream.Send(&pb.HelloRequest{Name: "linda"})
	}
	resp, _ := stream.CloseAndRecv()
	log.Printf("resp err: %v", resp)
	return nil
}

func SayRoute(client pb.GreeterClient) error {
	stream, _ := client.SayRoute(context.Background())
	for n := 0; n <= 6; n++ {
		_ = stream.Send(&pb.HelloRequest{Name: "linda"})
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}
	_ = stream.CloseSend()

	return nil
}

func HelloServer() {
	// 创建与服务端的连接句柄
	// conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())

	// consul 负载聚恒的方式连接
	consulStr := "consul://127.0.0.1:8500/shopping?healthy=true"
	// 客户端注册拦截器
	conn, _ := grpc.Dial(consulStr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.ClientUnaryInterceptor),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 指定round_robin策略

	)
	defer conn.Close()

	// 客户端对象，联系 errors 内部逻辑
	client := pb.NewGreeterClient(conn)
	_ = SayHello(client)
	// _ = SayList(client)
	// _ = SayRecord(client)
	// _ = SayRoute(client)
}

func main() {
	// TagServer(port)
	// var port string
	// flag.StringVar(&port, "p", "8081", "启动端口号")
	// flag.Parse()

	HelloServer()
}
