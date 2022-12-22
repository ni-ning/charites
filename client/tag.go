package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"charites/pkg/errcode"
	pb "charites/proto"
)

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	// opts = append(opts,  grpc.WithInsecure())
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.DialContext(ctx, target, opts...)
}

func TagServer(port string) {
	ctx := context.Background()
	clientConn, _ := GetClientConn(ctx, "127.0.0.1:"+port, nil)
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{})

	if err != nil {
		sts := errcode.FromError(err)
		details := sts.Details()
		log.Printf("err details: %v\n", details)
		return
	}

	log.Printf("resp: %v\n", resp)
}

// func main() {
// 	var port string
// 	flag.StringVar(&port, "p", "8000", "启动端口号")
// 	flag.Parse()

// 	TagServer(port)
// }
