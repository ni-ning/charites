# 微服务gRPC项目集-直播电商

`gRPC`框架从0到1搭建微服务平台，包括商品、库存与订单微服务，归纳总结核心技术点：

- <a href="/golangv2/shopping.html#_1-项目背景说明">1. 项目背景说明</a>
- <a href="/golangv2/shopping.html#_2-项目结构搭建">2. 项目结构搭建</a>
- <a href="/golangv2/shopping.html#_3-创建库表模型">3. 创建库表模型</a>
- <a href="/golangv2/shopping.html#_4-商品微服务">4. 商品微服务</a>
    * 创建库表模型
    * 编写proto文件
    * 生成proto代码
    * 实现服务端代码
    * gPRC-Gateway
    * Makefile 快速实现
- <a href="/golangv2/shopping.html#_5-库存微服务">5. 库存微服务</a>
    * 通用业务开发流程
    * 测试方法汇总
    * 并发资源竞争示例
    * 事务处理示例
    * 悲观锁实现并发
    * 乐观锁实现并发
    * 分布式锁实现并发
- <a href="/golangv2/shopping.html#_6-订单微服务">6. 订单微服务</a>
	* 微服务相互调用
    * 雪花算法订单号
    * 创建订单直接版
- <a href="/golangv2/shopping.html#_7-分布式事务">7. 分布式事务</a>
	* 分布式事务介绍
- <a href="/golangv2/shopping.html#_8-rocketmq入门">8. Rocketmq入门</a>
	* 本地安装RocketMQ
	* Go语言客户端
- <a href="golangv2/shopping.html#_9-分布式订单">9. 分布式订单</a>
	* 本地事务订单逻辑
	* 库存微服务回滚
	* 订单未支付延时消息
- <a href="/golangv2/shopping.html#_10-本地部署启动">10. 本地部署启动</a>
- 源码地址 [https://github.com/ni-ning/charites](https://github.com/ni-ning/charites)

## 1. 项目背景说明

- 常规直播电商业务与架构
```js
1. 直播
直播技术架构涉及很多方面
    1. 推拉流
        1. 推流端把实时的音视频数据推送到服务端，服务端（合流、转码、录制、转推、鉴黄）
        2. 拉流 看播
        3. 连麦PK --> 实时性要求很高
        4. 多人语音房，9连麦
        5. obs推流 --> 大型活动，专业音视频设备
    2. 技术栈：
        1. 前端：h5、ios、android  sdk
        2. 后端：C++、播放器、ffmpeg、webrtc、cdn
    3. 业务类
        1. webrtc 信令模块，任务模块
        2. im ：直播间聊天、私信等--> goim
        3. 点赞、送礼、排行榜、粉丝标签等等

2. 电商
电商业务涉及很多方面，包括无实物、有实物、O2O、B2C等，大型的电商架构
    1. 商品中心: SPU --> 品 (iphone13)、SKU --> item(iphone13 金色 128G)、类目中心
    2. 库存: 单一仓库、分区仓库
    3. 商户中心: 大品牌、经销商、核销、广告
    4. 订单中心: 订单、购物车
    5. 支付中心: 支付方式、定期支付、定金支付、混合支付、货到付款
    6. 物流中心: 寄快递、查快递
    7. 履约中心: 退货、换货、只退不换
    8. 用户中心: 地址服务、收藏服务、推荐
    9. 营销中心: 优惠券、满减券、折扣券、专属券、平台会员、店铺会员
    10. 广告推荐:
    11. 发票
```
完整的直播需要专门的音视频团队，或者采用三方的集成方案

本项目只实现部分微服务，以`打通后端架构，实践新技术`为目标，具体包括
1. `商品微服务` 侧重gRPC实现周边
2. `库存微服务` 侧重并发锁实现
3. `订单微服务` 侧重分布式事务实现


## 2. 项目结构搭建
```js
charites
    |- apps    // 实际项目中会拆分为不同的微服务项目
        |- order        // 订单微服务
        |- shoppig      // 商品微服务
        |- stock        // 库存微服务
    |- bootstrap    // 初始化各类配置
        |- init.go
        |- logger.go
        |- mysql.go
        |- redis.go
        |- rpc.go
        |- setting.go
        |- snowflake.go
    |- client   // 可作为测试客户端
    |- config   // 配置文件，如Server、App、Database
        |- config.yaml
    |- global       // 全局变量，如配置、数据库连接、日志
    |- middleware   // 拦截器，包括客户端和服务端
    |- model        // 模型数据
    |- pkg          // 项目公共模块
        |- errcode  // 错误码，定义NewError，自定义业务错误码
        |- logger   // 日志，定义NewLogger，结合初始化函数，得到全局变量
        |- registry // 注册服务中心封装工具
        |- setting  // 配置，定义NewSetting，结合初始化函数，得到全局变量
		|- utils    // 如获取出口IP
    |- proto        // gRPC 定义的传输模型 Protobuf
    |- storage
		|- sql	    // SQL文件
		|- logs		// 日志文件
    |- main.go
    |- Makefile     // 编译快捷命令
```


## 3. 创建库表模型

- 创建数据库表 `storage/sql/*.sql`
- 创建模型 `model/*.go`
- 库表直接关联关系

![](./assets/model.png)


## 4. 商品微服务

`gRPC`实现直播间商品列表时，需要的各类技术点


### 编写proto文件
```js
syntax = "proto3";

option go_package = ".;proto";

service Goods {
    // 获取直播间商品列表
    rpc GetGoodsListByRoomId(GetGoodsListRoomReq) returns (GoodsListReply){};

    // 获取商品详情
    rpc GetGoodsDetail(GetGoodsDetailReq) returns (GoodsDetailReply){};
}

message GetGoodsListRoomReq{
    int64 RoomId = 2;
}

message GoodsInfo {
    int64 GoodsId = 1;
    int64 CategoryId = 2;
    int32 Status = 3;
    string Title = 4;
    string MarketPrice = 5;
    string Price = 6;
    string Brief = 7;
    repeated string HeadImgs = 8;
}

message GoodsListReply {
    int64 CurrentGoodsId = 1;
    repeated GoodsInfo data = 2 ;
}

message GetGoodsDetailReq {
    int64 GoodsId = 1;
}

message GoodsDetailReply {
    int64 GoodsId = 1;
    int64 CategoryId = 2;
    int32 Status = 3;
    int64 Code = 4;
    string BrandName = 5;
    string Title = 6;
    string MarketPrice = 7;
    string Price = 8;
    string Brief = 9;
    repeated string HeadImgs = 10;
    repeated string Videos = 11;
    repeated string Detail = 12;
}
```

### 生成proto代码
```js
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/*.proto
```

### 实现服务端代码
```go
type GoodsServer struct {
	pb.UnimplementedGoodsServer
}

func NewGoodsServer() *GoodsServer {
	return &GoodsServer{}
}

func (g GoodsServer) GetGoodsListByRoomId(context.Context, *pb.GetGoodsListRoomReq) (*pb.GoodsListReply, error) {
	return nil, errcode.ToRPCError(errcode.Success)
}
func (g GoodsServer) GetGoodsDetail(context.Context, *pb.GetGoodsDetailReq) (*pb.GoodsDetailReply, error) {
	return nil, errcode.ToRPCError(errcode.Success)
}
```

### gPRC-Gateway

[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) 是一个`protoc`插件，生成一个反向代理服务器，实现通过`RESTful API`访问`gRPC`服务

![](./assets/gateway.png)

- 安装依赖库与插件
```go
// +build tools
package tools

import (
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

 go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

- 安装三方proto
```js 
// 下载 google 预定义proto文件 
https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto
https://github.com/googleapis/googleapis/blob/master/google/api/http.proto

// 拷贝到 protoc 编译器目录/include下
/Users/nining/go/install/protoc-3.20.1-osx-aarch_64/include
mkdir google/api
cp annotations.proto http.proto google/api/
```

- 重新定义proto
```js
service Goods {
    // 获取直播间商品列表
    rpc GetGoodsListByRoomId(GetGoodsListRoomReq) returns (GoodsListReply){
        option (google.api.http) = {
            post: "/v1/goods",
            body: "*"
        };
    };
    ...
}
```

- 生成代码
```sh
	protoc  \
    --go_out=.  \
    --go_opt=paths=source_relative  \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=.    \
    --grpc-gateway_opt paths=source_relative    \
    ./proto/*.proto
```
- 服务端启动HTTP代理
```go
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
		// 注册RegisterGoodsHandler
		err = pb.RegisterGoodsHandler(context.Background(), gwmux, conn)
		if err != nil {
			log.Fatalln("Failed to register gateway:", err)
		}
		gwServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", global.ServerSetting.HttpPort),
			Handler: gwmux,
		}
		// 提供gRPC-Gateway服务
		log.Printf("Serving gRPC-Gateway on http://%s:%d\n", ip.String(), global.ServerSetting.HttpPort)
		log.Fatalln(gwServer.ListenAndServe())
	}()
```

### Makefile 快速实现

- [Makefile 快速入门](https://zhuanlan.zhihu.com/p/350297509)
```js
.PHONY: all build run gotool clean

BINARY="charites_server"
PROTO_DIR=proto

all: gotool build

build:
	CGO_ENABLE=1 GOOS=darwin GOARCH=arm64 go build -o ${BINARY}

run:
	go run main.go

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY}; fi

gen:
	protoc  \
    --go_out=.  \
    --go_opt=paths=source_relative  \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=.    \
    --grpc-gateway_opt paths=source_relative    \
    $(shell find $(PROTO_DIR) -iname "*.proto")

hello:
	go run  client/helloworld.go

help:
	@echo "make build - 编译指定文件"
	@echo "make run - 直接运行项目"
	@echo "make clean - 删除编译文件"
	@echo "make gen - 生成pb及grpc代码"
```


## 5. 库存微服务

实际项目中会把商品服务、库存服务、订单服务拆分为不同的微服务，我们仅作为测试项目，代码写到同一个项目中

### 通用业务开发流程

参考 商品功能开发 模块
1. `storage/sql/stock.sql` 定义SQL语句
2. `model/stock.go` 定义数据模型
3. `proto/stock.proto` 生成gRPC代码结构
4. `make gen` 生成代码
5. `apps/stock/handler.go` 实现 `StockServer` 服务
6. `main.go` 注册服务 `pb.RegisterStockServer(s, stock.NewStockServer())`
7. `make run` 运行服务


### 测试方法汇总
```go
// 1. 命令行工具 grpcurl
grpcurl -plaintext -rpc-header 'authorization:"token"' 192.168.1.4:8081 list
grpcurl -plaintext -rpc-header 'authorization:"token"'  -d '{"GoodsId": 100, "Num":1000}' 192.168.1.4:8081  Stock.GetStock
// 2. 实现client端
go run client/stock.go
// 3. gPRC-Gateway HTTP形式
// 4. Swagger 文档
```

### 并发资源竞争示例
- 服务端未加锁示例
```go
// apps/stock/dao.go
func ReduceStock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	var stock model.Stock
	// 1. 查询现有库存
	db := global.DBEngine.WithContext(ctx).
		Model(&model.Stock{}).
		Where("id = ?", goodsId).
		First(&stock)
	// 不存在也会抛异常
	if db.Error != nil {
		global.Logger.Error("ErrorDBOperateStock", zap.String("error", db.Error.Error()))
		return nil, errcode.ToRPCError(errcode.ErrorDBOperateStock)
	}
	if db.RowsAffected == 0 {
		return nil, errcode.ToRPCError(errcode.ErrorNotFoundStock)
	}
	// 2. 校验库存
	if stock.Num-num < 0 {
		return nil, errcode.ErrorNotEnoughStock
	}
	// 3. 扣减库存并保存
	stock.Num -= num
	// global.DBEngine.WithContext(ctx).Save(&stock) // 更新所有字段
	err := global.DBEngine.WithContext(ctx).
		Model(&model.Stock{}).
		Where("id = ?", goodsId).
		Updates(map[string]interface{}{
			"num": stock.Num,
		}).Error
	if err != nil {
		global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
		return nil, errcode.ToRPCError(errcode.ErrorDBOperateStock)
	}
	return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
}
```
- 客户端并发20请求
```go
// go run client/stock.go
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
```
- 最终出产生资源竞争问题

### 事务处理示例
- 将多次数据库操作包装在一个事务中，实现要么成功要么失败，不会出现一个成功一个失败的情况
- 批量扣减库存可以使用事务，事务解决不了资源竞争问题

```go
func ReduceStockWithTransaction(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	var stock model.Stock
	err := global.DBEngine.Transaction(func(tx *gorm.DB) error {
		// 1. 查询现有库存
		db := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			First(&stock)
		// 不存在也会抛异常
		if db.Error != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", db.Error.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if db.RowsAffected == 0 {
			return errcode.ToRPCError(errcode.ErrorNotFoundStock)
		}
		// 2. 校验库存
		if stock.Num-num < 0 {
			return errcode.ErrorNotEnoughStock
		}
		// 3. 扣减库存并保存
		stock.Num -= num
		// global.DBEngine.WithContext(ctx).Save(&stock) // 更新所有字段
		err := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			Updates(map[string]interface{}{
				"num": stock.Num,
			}).Error
		if err != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		// return nil 提交事务，任何类型err回滚事务
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
}
```

### 悲观锁实现并发

- [通俗易懂 锁入门](https://zhuanlan.zhihu.com/p/71156910)
- [分布式锁看这篇就够了](https://zhuanlan.zhihu.com/p/42056183)
- [MySQL 行级锁](https://www.cnblogs.com/zping/p/10955750.html)

悲观锁，对一切事情比较悲观，我更新数据，就觉得所有人都要来跟我抢

从查询数据的时候就给这条数据加锁，保证只有我能更新

- 原生SQL
```sql
start();
select * from t1 where goods_id = 1 for update;
update t1 set num = 1 where goods_id = 1;
commit();
```
- 代码实现
```go
// 事务中添加 Clauses(clause.Locking{Strength: "UPDATE"}) 即可
func ReduceStockWithPessimisticLock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	var stock model.Stock
	err := global.DBEngine.Transaction(func(tx *gorm.DB) error {
		// 1. 查询现有库存
		db := tx.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			First(&stock)
		// 不存在也会抛异常
		if db.Error != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", db.Error.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if db.RowsAffected == 0 {
			return errcode.ToRPCError(errcode.ErrorNotFoundStock)
		}
		// 2. 校验库存
		if stock.Num-num < 0 {
			return errcode.ErrorNotEnoughStock
		}
		// 3. 扣减库存并保存
		stock.Num -= num
		// global.DBEngine.WithContext(ctx).Save(&stock) // 更新所有字段
		err := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			Updates(map[string]interface{}{
				"num": stock.Num,
			}).Error
		if err != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		// return nil 提交事务，任何类型err回滚事务
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
}
```
- 注意事项
1. 一定是基于索引来查询
2. 放到事务中处理


### 乐观锁实现并发
和悲观锁一样都是宏观的一个概念，本质上不算锁

乐观锁认为一般不会有人跟我竞争资源，通过version版本号在更新的时候做check

- 原生SQL
```sql
select goods_id,num,version from shopping_stock where goods_id = 1;
update shopping_stock set num=1,version=version+1 where goods_id = 1 and verison=verison;
```
- 代码实现
```go
func ReduceStockWithOptimisticLock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	for retry := 0; retry < 20; retry++ {
		var stock model.Stock
		// 1. 查询现有库存
		db := global.DBEngine.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			First(&stock)
		// 不存在也会抛异常
		if db.Error != nil {
			global.Logger.Error("model.Stock.First", zap.String("error", db.Error.Error()))
			return nil, errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if db.RowsAffected == 0 {
			return nil, errcode.ToRPCError(errcode.ErrorNotFoundStock)
		}
		// 2. 校验库存
		if stock.Num-num < 0 {
			return nil, errcode.ErrorNotEnoughStock
		}
		// 3. 扣减库存并保存
		ret := global.DBEngine.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ? and version = ?", goodsId, stock.Version).
			Updates(map[string]interface{}{
				"num":     stock.Num - 1,
				"version": stock.Version + 1,
			})
		if ret.Error != nil {
			global.Logger.Error("model.Stock.Updates", zap.String("error", ret.Error.Error()))
			return nil, errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if ret.RowsAffected == 0 {
			// 说明 version 被更新，重试即可
			continue
		}
		return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
	}
	return nil, errcode.ToRPCError(errcode.ErrorNeedRetryStock)
}
```
- 注意事项
1. `var stock model.Stock` 定义在 for 循环里面
2. `continue` 重试逻辑判断点


### 分布式锁实现并发

借助其他的组件：redis、zookeeper、etcd

基于redis实现：[https://github.com/go-redsync/redsync](https://github.com/go-redsync/redsync)

原生redis实现：setnx [https://www.redis.net.cn/order/3552.html](https://www.redis.net.cn/order/3552.html)

完善的基于redis的分布式锁：redlock [https://zhuanlan.zhihu.com/p/62769627](https://zhuanlan.zhihu.com/p/62769627)

- 代码实现
```go
func ReduceStockWithDistributedLock(ctx context.Context, goodsId, num int64) (*pb.GoodsStockInfo, error) {
	mutexname := fmt.Sprintf("reduce:stock:mutex:%d", goodsId)
	mutex := global.Redsync.NewMutex(mutexname)
	if err := mutex.Lock(); err != nil {
		return nil, errcode.ToRPCError(errcode.ErrorRedisLockStock)
	}
	defer mutex.Unlock()

	var stock model.Stock
	err := global.DBEngine.Transaction(func(tx *gorm.DB) error {
		// 1. 查询现有库存
		db := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			First(&stock)
		// 不存在也会抛异常
		if db.Error != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", db.Error.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		if db.RowsAffected == 0 {
			return errcode.ToRPCError(errcode.ErrorNotFoundStock)
		}
		// 2. 校验库存
		if stock.Num-num < 0 {
			return errcode.ErrorNotEnoughStock
		}
		// 3. 扣减库存并保存
		stock.Num -= num
		// global.DBEngine.WithContext(ctx).Save(&stock) // 更新所有字段
		err := tx.WithContext(ctx).
			Model(&model.Stock{}).
			Where("id = ?", goodsId).
			Updates(map[string]interface{}{
				"num": stock.Num,
			}).Error
		if err != nil {
			global.Logger.Error("ErrorDBOperateStock", zap.String("error", err.Error()))
			return errcode.ToRPCError(errcode.ErrorDBOperateStock)
		}
		// return nil 提交事务，任何类型err回滚事务
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.GoodsStockInfo{GoodsId: goodsId, Num: stock.Num}, nil
}
```

小结：
- 悲观锁阻塞事务 乐观锁回滚重试
- 乐观锁，本质上不加锁，适用于写操作少的场景


## 6. 订单微服务


### 微服务相互调用

- 修改`Makefile`，启用不同端口号实例微服务
```js
// 商品微服务
run_goods:
	go run main.go -p 8090

// 库存微服务
run_stock:
	go run main.go -p 8092
```

- 初始化微服务客户端
```go
// bootstrap/rpc.go
func setupRPClient() error {
	// 商品微服务客户端
	goodsConn, err := grpc.Dial("127.0.0.1:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.ClientUnaryInterceptor),
	)
	if err != nil {
		return err
	}
	global.GoodsCli = pb.NewGoodsClient(goodsConn)

	// 库存微服务客户端
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
```

实现订单微服务，直接通过全局客户端调用其他微服务接口
```go
resp, err := global.GoodsCli.GetGoodsDetail(context.Background(), &pb.GetGoodsDetailReq{GoodsId: 1})
resp, err := global.StockCli.ReduceStock(context.Background(), &pb.GoodsStockInfo{GoodsId: 1, Num: 1})
```

### 雪花算法订单号

- [分布式ID神器之雪花算法简介](https://zhuanlan.zhihu.com/p/85837641)
- [雪花算法go实现](https://github.com/bwmarrin/snowflake)

分布式服务，需要把雪花算法当成一个独立的服务部署

```go
import (
	"charites/global"
	"errors"
	"time"

	sf "github.com/bwmarrin/snowflake"
)
// global/snowflake.go
var SnowNode *sf.Node

// bootstrap/snowflake.go
const (
	_defaultStartTime = "2021-12-31"
)

func setupSnowflake(startTime string, machineId int64) error {
	if machineId < 0 {
		return errors.New("snowflake need machineId")
	}
	if len(startTime) == 0 {
		startTime = _defaultStartTime
	}
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	sf.Epoch = st.UnixNano() / 100_0000          // 时间戳，开始时间 69年
	global.SnowNode, err = sf.NewNode(machineId) // 机器编号，1024
	if err != nil {
		return err
	}
	return nil
}

err = setupSnowflake("", 1)
if err != nil {
    log.Fatalf("init.setupSnowflake err: %v", err)
}

// pkg/utils.go
func GenId() int64 {
	// 坑：前端展示不了 int64，需要String()
	return global.SnowNode.Generate().Int64()
}
```

### 创建订单直接版
```go
// CreateOrder 创建订单
func CreateOrder(ctx context.Context, req *pb.OrderReq) (*emptypb.Empty, error) {
	// 生成订单号
	orderId := utils.GenInt64()

	// 请求商品微服务
	goodsDetail, err := global.GoodsCli.GetGoodsDetail(context.Background(), &pb.GetGoodsDetailReq{GoodsId: req.GoodsId})
	if err != nil {
		return nil, errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
	}
	// 拿到商品价格作为支付价格
	price, _ := strconv.ParseInt(goodsDetail.Price, 10, 64)

	// 请求库存微服务，扣减库存
	_, err = global.StockCli.ReduceStock(context.Background(), &pb.GoodsStockInfo{GoodsId: req.GoodsId, Num: req.Num})
	if err != nil {
		return nil, errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
	}

	// 创建订单与订单详情
	orderData := model.Order{
		UserId:         req.UserId,
		OrderId:        orderId, // 雪花算法生成
		TradeId:        fmt.Sprintf("%d", orderId),
		Status:         int64(100), // 创建订单初始状态
		ReceiveAddress: req.Address,
		ReceiveName:    req.Name,
		ReceivePhone:   req.Phone,
		PayAmount:      price * req.Num, // 该订单总价
	}

	marketPrice, _ := strconv.ParseInt(goodsDetail.MarketPrice, 10, 64)
	orderDetail := model.OrderDetail{
		UserId:    req.UserId,
		OrderId:   orderId, // 雪花算法生成
		GoodsId:   req.GoodsId,
		Num:       req.Num,
		PayAmount: price * req.Num, // 该商品总价

		Title:       goodsDetail.Title,
		MarketPrice: marketPrice,
		Price:       price,
		Brief:       goodsDetail.Brief,
		HeadImgs:    strings.Join(goodsDetail.HeadImgs, ","),
		Videos:      strings.Join(goodsDetail.Videos, ","),
		Detail:      strings.Join(goodsDetail.Detail, ","),
	}
	err = global.DBEngine.Transaction(func(tx *gorm.DB) error {
		orderResult := tx.WithContext(ctx).Create(&orderData)
		if orderResult.Error != nil {
			return errcode.ToRPCError(errcode.ErrorCreateOrder)
		}
		orderDetailResult := tx.WithContext(ctx).Create(&orderDetail)
		if orderDetailResult.Error != nil {
			return errcode.ToRPCError(errcode.ErrorCreateOrderDetal)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
```

- 测试创建订单
```go
grpcurl \
-plaintext  \
-rpc-header 'authorization:"token"'  \
-d '{"GoodsId": 1, "Num": 2, "UserId": 1, "Address":"BJ", "Name":"linda", "Phone":"18210980038"}' \
192.168.1.4:8081  \
Order.CreateOrder
```

- 存在的问题
当扣减库存成功，但本地创建订单失败时，会导致数据不一致


## 7. 分布式事务

### 分布式事务介绍

微服务架构下带来的挑战：怎么解决分布式场景下数据一致性问题，分布式事务

- 讨论的前提：理论依据
```js
本地事务、分布式事务
强一致性、弱一致性、最终一致性
CAP理论：C一致性 A可用性 P分区容错性
BASE理论：面向的是大型高可用可扩展的分布式系统，和传统的事物ACID特性是相反的，它完全不同于ACID的强一致性模型，而是通过牺牲强一致性来获得可用性，并允许数据在一段时间内是不一致的，但最终达到一致状态
柔性事务
可见性(对外可查询)，全局唯一的标识用于查询
幂等操作，方便重试
```


- 常见分布式事务实现方式
```js


最大努力通知
本质：通过定期校对，实现数据一致性
- 支付宝/微信支付 通过回调的方式通知业务方支付状态
- callback --> 1 3 5 10 15 30 60
- 提供一个查询接口，业务方主动去查询
场景：适用于对业务最终一致性的时间敏感度低的系统
```


https://github.com/dtm-labs/dtm/blob/main/helper/README-cn.md


- [看一遍就理解：分布式事务详解](https://zhuanlan.zhihu.com/p/516554367)







## 8. RocketMQ入门

- 官文文档 [https://github.com/apache/rocketmq/tree/master/docs/cn](https://github.com/apache/rocketmq/tree/master/docs/cn)

### 本地安装RocketMQ

- 推荐使用docker-compose 快速搭建本地开发环境
- [https://github.com/foxiswho/docker-rocketmq](https://github.com/foxiswho/docker-rocketmq)

```go
git clone  https://github.com/foxiswho/docker-rocketmq.git
cd docker-rocketmq
cd rmq
```
- 修改一下`docker-compose.yml`文件，暂时使用 阿里云 镜像库里的4.7.0版本
```go
version: '3.5'

services:
  rmqnamesrv:
#    image: foxiswho/rocketmq:4.9.2
    image: registry.cn-hangzhou.aliyuncs.com/foxiswho/rocketmq:4.7.0
    container_name: rmqnamesrv
    ports:
      - 9876:9876
    volumes:
      - ./rmqs/logs:/home/rocketmq/logs
      - ./rmqs/store:/home/rocketmq/store
    environment:
      JAVA_OPT_EXT: "-Duser.home=/home/rocketmq -Xms512M -Xmx512M -Xmn128m"
    command: ["sh","mqnamesrv"]
    networks:
        rmq:
          aliases:
            - rmqnamesrv
  rmqbroker:
#    image: foxiswho/rocketmq:4.9.2
    image: registry.cn-hangzhou.aliyuncs.com/foxiswho/rocketmq:4.7.0
    container_name: rmqbroker
    ports:
      - 10909:10909
      - 10911:10911
    volumes:
      - ./rmq/logs:/home/rocketmq/logs
      - ./rmq/store:/home/rocketmq/store
      - ./rmq/brokerconf/broker.conf:/etc/rocketmq/broker.conf
    environment:
        JAVA_OPT_EXT: "-Duser.home=/home/rocketmq -Xms512M -Xmx512M -Xmn128m"
    command: ["sh","mqbroker","-c","/etc/rocketmq/broker.conf","-n","rmqnamesrv:9876","autoCreateTopicEnable=true"]
    depends_on:
      - rmqnamesrv
    networks:
      rmq:
        aliases:
          - rmqbroker

  rmqconsole:
    image: styletang/rocketmq-console-ng
    container_name: rmqconsole
    ports:
      - 8180:8080
    environment:
        JAVA_OPTS: "-Drocketmq.namesrv.addr=rmqnamesrv:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false"
    depends_on:
      - rmqnamesrv
    networks:
      rmq:
        aliases:
          - rmqconsole

networks:
  rmq:
    name: rmq
    driver: bridge
```
- 修改配置文件 `vim rmq/brokerconf/broker.conf`

```go
// 将33行取消注释，并将 `brokerIP1` 为你本机的IP地址
brokerIP1=192.168.1.4
```
- 执行本地安装
```go
chmod +x  start.sh
./start.sh

// 本地访问
http://localhost:8180
```
![](./assets/rmq.png)


### Go语言客户端

- [https://github.com/apache/rocketmq-client-go](https://github.com/apache/rocketmq-client-go)

- [https://github.com/apache/rocketmq-client-go/tree/master/examples](https://github.com/apache/rocketmq-client-go/tree/master/examples)


## 9. 分布式订单

- 基于RocketMQ事务消息实现订单微服务的分布式事务
- 逆向思路：先尝试返送回滚库存消息
	* 本地事务成功，撤销 滚库存消息
	* 本地事务失败，确认 滚库存消息

### 本地事务订单逻辑
- 按照 RocktMQ 事务消息实现两个方法 `ExecuteLocalTransaction`和 `CheckLocalTransaction`
```go
// OrderListener 自定义结构体，实现两个方法
// 发送事务消息的时候，RocketMQ会根据情况自动调用这两个方法
type OrderListener struct {
	OrderId int64
	Param   *pb.OrderReq
	Err     error
}

// 当发送prepare(half) message 成功后，这个方法(执行本地事务)就会被执行
func (o *OrderListener) ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState {
	if o.Param == nil {
		global.Logger.Error("ExecuteLocalTransaction param is nil")
		o.Err = errcode.ToRPCError(errcode.ErrorOrderEntityParam)
		// 库存未扣减
		return primitive.RollbackMessageState
	}
	param := o.Param
	ctx := context.Background()

	// 请求商品微服务，查询商品金额(营销相关)
	goodsDetail, err := global.GoodsCli.GetGoodsDetail(ctx, &pb.GetGoodsDetailReq{GoodsId: param.GoodsId})
	if err != nil {
		global.Logger.Error("GoodsCli.GetGoodsDetail failed", zap.Error(err))
		o.Err = errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
		// 库存未扣减
		return primitive.RollbackMessageState
	}
	// 拿到商品价格作为支付价格
	price, _ := strconv.ParseInt(goodsDetail.Price, 10, 64)

	// 请求库存微服务，扣减库存
	_, err = global.StockCli.ReduceStock(ctx, &pb.GoodsStockInfo{GoodsId: param.GoodsId, Num: param.Num, OrderId: o.OrderId})
	if err != nil {
		global.Logger.Error("StockCli.ReduceStock failed", zap.Error(err))
		o.Err = errcode.ToRPCError(errcode.ErrorRPCOrderToGoods)
		// 库存未扣减
		return primitive.RollbackMessageState
	}

	// 本地事务创建订单与订单详情
	orderData := model.Order{
		UserId:         param.UserId,
		OrderId:        o.OrderId, // 雪花算法生成
		TradeId:        fmt.Sprintf("%d", o.OrderId),
		Status:         int64(100), // 创建订单初始状态
		ReceiveAddress: param.Address,
		ReceiveName:    param.Name,
		ReceivePhone:   param.Phone,
		PayAmount:      price * param.Num, // 该订单总价
	}
	marketPrice, _ := strconv.ParseInt(goodsDetail.MarketPrice, 10, 64)
	orderDetail := model.OrderDetail{
		UserId:      param.UserId,
		OrderId:     o.OrderId, // 雪花算法生成
		GoodsId:     param.GoodsId,
		Num:         param.Num,
		PayAmount:   price * param.Num, // 该商品总价
		Title:       goodsDetail.Title,
		MarketPrice: marketPrice,
		Price:       price,
		Brief:       goodsDetail.Brief,
		HeadImgs:    strings.Join(goodsDetail.HeadImgs, ","),
		Videos:      strings.Join(goodsDetail.Videos, ","),
		Detail:      strings.Join(goodsDetail.Detail, ","),
	}
	err = global.DBEngine.Transaction(func(tx *gorm.DB) error {
		orderResult := tx.WithContext(ctx).Create(&orderData)
		if orderResult.Error != nil {
			return errcode.ToRPCError(errcode.ErrorCreateOrder)
		}
		orderDetailResult := tx.WithContext(ctx).Create(&orderDetail)
		if orderDetailResult.Error != nil {
			return errcode.ToRPCError(errcode.ErrorCreateOrderDetal)
		}
		return nil
	})
	if err != nil {
		// 本地事务执行失败，但上一步库存已经扣减成功
		return primitive.CommitMessageState
	}

	// 发送延迟消息
	// 不同等级：1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
	// 消息中具体的载荷，定义为一个结构体，赞
	data := model.OrderGoodsStockInfo{
		OrderId: o.OrderId,
		GoodsId: param.GoodsId,
		Num:     param.Num,
	}
	b, _ := json.Marshal(data)
	// 定义RocketMQ消息体
	msg := primitive.NewMessage(global.RocketMQSetting.TopicOrderPayTimeout, b)
	msg.WithDelayTimeLevel(3)
	_, err = global.Producer.SendSync(context.Background(), msg)
	if err != nil {
		// 延时消息发送失败
		global.Logger.Error("send delay msg failed", zap.Error(err))
		return primitive.CommitMessageState
	}
	// 说明本地事务执行成功，不需要发送回滚库存的消息
	return primitive.RollbackMessageState
}

// 当发送prepare(half) message 没有响应时，broker会回查本地事务状态，此时这个方法被执行
func (o *OrderListener) CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState {
	// 检查本地是否订单创建成功即可
	var count int64
	global.DBEngine.
		WithContext(context.Background()).
		Model(&model.Order{}).Where("order_id = ?", o.OrderId).
		Count(&count)
	if count <= 0 {
		// 说明订单创建失败，需要回滚库存
		return primitive.CommitMessageState
	}
	// 不存回滚库存
	return primitive.RollbackMessageState
}
```

### 库存微服务回滚

不能简单的收到回滚库存消息就回滚库存，因为有可能消息重复了，导致多次回滚，数据不一致的问题

- 库存微服务启动消息监听
```go
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

// RollbackMsgHandle 监听RocketMQ消息进行库存回滚的处理函数
// 需要考虑重复归还的问题(幂等性) --> 添加库存扣减记录
func RollbackMsgHandle(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var data model.OrderGoodsStockInfo
		err := json.Unmarshal(msgs[i].Body, &data)
		if err != nil {
			global.Logger.Error("json.Unmarshal RollbackMsg error", zap.Error(err))
			continue
		}
		err = RollbackStockByMsg(ctx, &data)
		if err != nil {
			return consumer.ConsumeRetryLater, nil
		}
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
```

### 订单未支付延时消息

1. 什么时机发送延迟消息？
   1. 创建的订单时候 --> 发延迟消息 -->30分钟
2. 发送方是谁？接收方又是谁？
   1. 订单服务发送
   2. 库存作为接收方的问题 --> 收到这个延迟消息就要回滚库存吗？
      1.  并不是，我们需要根据订单的状态去判断是否执行库存回滚
   3. 我们仍然选择在订单服务接收延时消息
      1. 收到消息就可以直接判断订单状态，
      2. 如果是**未支付状态**就发送一条回滚库存的消息给库存服务，复用上一步的`shopping_stock_rollback`这个topic

- 订单微服务监听超时消息
```go
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
```

## 10. 本地部署启动

```go
// 启动商品微服务
make run_goods

// 启动库存微服务
make run_stock

// 启动订单微服务
make run
```