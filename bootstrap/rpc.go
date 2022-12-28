package bootstrap

import (
	"charites/global"
	"charites/middleware"
	pb "charites/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setupRPClient() error {
	global.Logger.Info("setupRPClient Goods Serivce...")
	goodsConn, err := grpc.Dial("127.0.0.1:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.ClientUnaryInterceptor),
	)
	if err != nil {
		return err
	}
	global.GoodsCli = pb.NewGoodsClient(goodsConn)

	global.Logger.Info("setupRPClient Stock Serivce...")
	stockConn, err := grpc.Dial("127.0.0.1:8092",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.ClientUnaryInterceptor),
	)
	if err != nil {
		return err
	}
	global.StockCli = pb.NewStockClient(stockConn)
	return nil
}
