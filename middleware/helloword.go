package middleware

import (
	"charites/pkg/errcode"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ClientUnaryInterceptor 客户端一元拦截器
func ClientUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	md := metadata.Pairs("authorization", "token")
	ctx = metadata.NewOutgoingContext(ctx, md)

	err := invoker(ctx, method, req, reply, cc, opts...)

	end := time.Now()
	fmt.Printf("ClientUnaryInterceptor: %s, start time: %s, end time: %s, err: %v\n", method, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return err
}

// ServerUnaryInterceptor 服务端一元拦截器
func ServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 健康检查不需要校验token
	if info.FullMethod != "/grpc.health.v1.Health/Check" {
		// authentication (token verification)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		if _, ok := md["authorization"]; !ok {
			return nil, errcode.ToRPCError(errcode.Unauthorized)
		}
	}

	m, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("RPC failed with error %v\n", err)
	}

	return m, err
}
