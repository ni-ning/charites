package tag

import (
	"charites/pkg/bapi"
	"charites/pkg/errcode"
	pb "charites/proto"
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type TagServer struct {
	pb.UnimplementedTagServiceServer
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, req *pb.GetTagListRequest) (*pb.GetTagListReplay, error) {
	// ctx 包含请求头，如metadata相关信息
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Printf("md: %+v\n", md)

	// 测试PRC异常的异常值
	// return nil, errcode.ToRPCError(errcode.ErrorGetTagListFail)

	api := bapi.NewAPI("http://127.0.0.1:8090")
	// 获取具体的业务数据
	body, err := api.GetTagList(ctx)
	if err != nil {
		return nil, err
	}
	// 业务数据映射到返回对象 GetTagListReply
	tagList := pb.GetTagListReplay{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, errcode.ToRPCError(errcode.Fail)
	}
	return &tagList, nil
}

// TestError 测试内部客户端解析details
func TestError() {
	err := errcode.ToRPCError(errcode.ErrorGetTagListFail)
	sts := errcode.FromError(err)
	details := sts.Details()
	fmt.Println(details...)
}
