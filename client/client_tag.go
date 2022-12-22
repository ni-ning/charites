package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"charites/pkg/errcode"
	pb "charites/proto"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}

func main1() {
	// 1. 直接创建ctx
	ctx := context.Background()
	// 2. timeout + cancel()
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	// defer cancel()

	// 1. metadata.Pairs
	// md := metadata.Pairs("k1", "v1", "k2", "v2")
	// 2.metadata.New
	// md := metadata.New(map[string]string{
	// 	"k3": "v3",
	// 	"k4": "v4",
	// })

	md := metadata.New(map[string]string{
		"k3": "v3",
		"k4": "v4",
	})
	// 1. 新增md
	newCtx := metadata.NewOutgoingContext(ctx, md)

	// 2. 附加md信息
	// newCtx := metadata.AppendToOutgoingContext(ctx, "ctx_name", "linda")

	clientConn, _ := GetClientConn(newCtx, "127.0.0.1:"+port, nil)
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{})

	if err != nil {
		sts := errcode.FromError(err)
		details := sts.Details()
		log.Printf("err details: %v\n", details)
		return
	}

	log.Printf("resp: %v\n", resp)
}
