package main

import (
	"charites/apps/helloword"
	"charites/apps/order"
	"charites/apps/shopping"
	"charites/apps/stock"
	"charites/apps/tag"
	_ "charites/bootstrap"
	"charites/global"
	"charites/middleware"
	"charites/pkg/registry"
	"charites/pkg/utils"
	pb "charites/proto"
	"context"
	"flag"
	"net/http"

	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/health"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var grpcPort int
var httpPort int

func init() {
	flag.IntVar(&grpcPort, "p", 8081, "gRPC端口号")
	flag.Parse()
	httpPort = grpcPort + 1
}

func StartStockConsume() {
	// 库存微服务启动消息监听
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{global.RocketMQSetting.NameServer})),
		consumer.WithGroupName(global.RocketMQSetting.GroupStockService),
	)
	// 监听Topck
	err := c.Subscribe(global.RocketMQSetting.TopicStockRollback, consumer.MessageSelector{}, stock.RollbackMsgHandle)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = c.Start()
	if err != nil {
		panic(err)
	}
}

func StartOrderConsume() {
	// 订单微服务监听超时消息
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{global.RocketMQSetting.NameServer})),
		consumer.WithGroupName(global.RocketMQSetting.GroupOrderService),
	)
	// 订阅topic
	err := c.Subscribe(global.RocketMQSetting.TopicOrderPayTimeout, consumer.MessageSelector{}, order.OrderTimeoutHandle)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = c.Start()
	if err != nil {
		panic(err)
	}
}

func main() {
	if grpcPort != 0 {
		global.ServerSetting.GrpcPort = grpcPort
	}
	if httpPort != 0 {
		global.ServerSetting.HttpPort = httpPort
	}

	ip, _ := utils.GetOutBoundIp()

	// 创建 gRPC 服务端启动对象，NewServer构造函数支持选项，如服务端拦截器
	// s := grpc.NewServer()
	s := grpc.NewServer(grpc.UnaryInterceptor(middleware.ServerUnaryInterceptor))

	// 注册 gRPC 接口服务1：业务接口服务
	pb.RegisterGreeterServer(s, helloword.NewGreeterServer())
	pb.RegisterTagServiceServer(s, tag.NewTagServer())
	pb.RegisterGoodsServer(s, shopping.NewGoodsServer())
	pb.RegisterStockServer(s, stock.NewStockServer())
	pb.RegisterOrderServer(s, order.NewOrderServer())

	// 注册 gRPC 接口服务2：三方插件接口服务
	reflection.Register(s)                               // grpcurl
	healthpb.RegisterHealthServer(s, health.NewServer()) // 健康检查

	// 监听 TCP 端口号，底层通用 net 库
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", global.ServerSetting.GrpcPort))
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	// 启动 gRPC 服务端轮询，为阻塞服务，结合goroutine实现 *赞*
	go func() {
		// 启动 RPC 服务
		log.Printf("Serving gRPC on %s:%d\n", ip.String(), global.ServerSetting.GrpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("s.Serve err: %v\n", err)
		}
	}()

	// 注册服务到注册中心
	client := registry.NewClient()
	client.RegisterService(global.ServerSetting.ServiceName, ip.String(), global.ServerSetting.GrpcPort)

	// gRPC-Gateway
	go func() {
		// 创建一个连接到我们刚刚启动的 gRPC 服务器的客户端连接
		// gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
		conn, err := grpc.DialContext(
			context.Background(),
			fmt.Sprintf("%s:%d", ip.String(), global.ServerSetting.GrpcPort),
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalln("grpc.DialContext err:", err)
		}
		gwmux := runtime.NewServeMux()
		// gRPC服务映射为HTTP服务
		err = pb.RegisterGoodsHandler(context.Background(), gwmux, conn)
		if err != nil {
			log.Fatalln("pb.RegisterGoodsHandler err:", err)
		}
		err = pb.RegisterStockHandler(context.Background(), gwmux, conn)
		if err != nil {
			log.Fatalln("pb.RegisterStockHandler err:", err)
		}
		err = pb.RegisterOrderHandler(context.Background(), gwmux, conn)
		if err != nil {
			log.Fatalln("pb.RegisterOrderHandler err:", err)
		}

		gwServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", global.ServerSetting.HttpPort),
			Handler: gwmux,
		}
		// 提供gRPC-Gateway服务
		log.Printf("Serving gRPC-Gateway on http://%s:%d\n", ip.String(), global.ServerSetting.HttpPort)
		log.Fatalln(gwServer.ListenAndServe())
	}()

	// 关闭服务流程
	quitChan := make(chan os.Signal, 1) // 在代码里接收操作系统发来的中断信号
	// syscall.SIGTERM(kill)、syscall.SIGINT(ctrl+c)、syscall.SIGKILL(kill -9)
	signal.Notify(quitChan, syscall.SIGTERM, syscall.SIGINT)
	<-quitChan // 一直卡住，直到收到中断信号

	global.Logger.Info("*服务关闭清理流程*")
	serviceId := fmt.Sprintf("%s-%s-%d", global.ServerSetting.ServiceName, ip, global.ServerSetting.GrpcPort)
	global.Logger.Info("注销服务: ", zap.String("serviceId", serviceId))
	client.DeregisterService(serviceId)
}
