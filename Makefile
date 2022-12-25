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


