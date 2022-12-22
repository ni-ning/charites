package main

import (
	"charites/apps/helloword"
	"charites/apps/tag"
	_ "charites/bootstrap"
	"charites/global"
	pb "charites/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, helloword.NewGreeterServer())
	pb.RegisterTagServiceServer(s, tag.NewTagServer())

	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+global.ServerSetting.HttpPort)
	if err != nil {
		log.Fatalf("net.Listen err:%v", err)
	}

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("s.Serve err:%v", err)
	}
}
