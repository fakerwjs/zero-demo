# zero-demo Makefile
# 说明：goctl 依赖 protoc / protoc-gen-go / protoc-gen-go-grpc，需确保它们在 PATH 中。

GOCTL := goctl
STYLE := go_zero

.PHONY: help tidy build vet gen-user-rpc gen-user-api gen-notification-rpc run-user-rpc run-user-api run-notification-rpc

help:
	@echo "make tidy                 - go mod tidy"
	@echo "make build                - 编译所有服务"
	@echo "make vet                  - go vet 全部包"
	@echo "make gen-user-rpc         - 重新生成 user rpc（改 proto 后）"
	@echo "make gen-user-api         - 重新生成 user api（改 .api 后）"
	@echo "make gen-notification-rpc - 重新生成 notification rpc（改 proto 后）"
	@echo "make run-user-rpc         - 启动 user rpc"
	@echo "make run-user-api         - 启动 user api"
	@echo "make run-notification-rpc - 启动 notification rpc"

tidy:
	go mod tidy

build:
	go build ./...

vet:
	go vet ./...

gen-user-rpc:
	cd app/user/rpc && $(GOCTL) rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style=$(STYLE)

gen-user-api:
	cd app/user/api && $(GOCTL) api go -api user.api -dir . --style=$(STYLE)

gen-notification-rpc:
	cd app/notification/rpc && $(GOCTL) rpc protoc notification.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style=$(STYLE)

run-user-rpc:
	cd app/user/rpc && go run user.go -f etc/user.yaml

run-user-api:
	cd app/user/api && go run user.go -f etc/user-api.yaml

run-notification-rpc:
	cd app/notification/rpc && go run notification.go -f etc/notification.yaml
