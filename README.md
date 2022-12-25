
## 代码生成
```
protoc  \
    --go_out=.  \
    --go_opt=paths=source_relative  \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=.    \
    --grpc-gateway_opt paths=source_relative    \
    ./proto/*.proto

```

## 安装调试工具

```
go get github.com/fullstorydev/grpcurl
// 如果报错 可进入 ~/go/pkg 试试安装
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```