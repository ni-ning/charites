package main

import (
	pb "charites/proto"
	"context"
	"flag"
	"io"
	"log"

	"google.golang.org/grpc"
)

var port1 string

func init() {
	flag.StringVar(&port1, "p", "8000", "启动端口号")
	flag.Parse()
}

func SayHello(client pb.GreeterClient) error {
	// 像调用本地函数一样
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "linda"})
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

func main() {
	// 创建与服务端的连接句柄
	conn, _ := grpc.Dial(":"+port1, grpc.WithInsecure())
	defer conn.Close()

	// 客户端对象，联系 errors 内部逻辑
	client := pb.NewGreeterClient(conn)
	_ = SayHello(client)
	_ = SayList(client)
	_ = SayRecord(client)
	_ = SayRoute(client)
}
