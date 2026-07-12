#!/usr/bin/env bash
# 一键启动所有服务：etcd + gateway-api + user-api + user.rpc + notification.rpc + product.rpc + order.rpc + payment.rpc
# 用法：bash scripts/run-all.sh
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RUN="$ROOT/.run"
BIN="$RUN/bin"
mkdir -p "$BIN"

# 本机工具目录（goctl/protoc/go），按需修改
export PATH="/d/bin:/d/protoc/bin:/c/Users/yihe/go/bin:$PATH"

# etcd 可执行文件路径（按需修改）
ETCD_BIN="${ETCD_BIN:-/d/etcd/etcd-v3.7.0-rc.0-windows-amd64/etcd}"
ETCDCTL_BIN="${ETCDCTL_BIN:-/d/etcd/etcd-v3.7.0-rc.0-windows-amd64/etcdctl}"

echo ">> 检查 etcd 是否已在运行..."
if "$ETCDCTL_BIN" --endpoints=127.0.0.1:2379 endpoint health >/dev/null 2>&1; then
  echo "   etcd 已在运行"
else
  echo ">> 启动 etcd..."
  rm -rf "$RUN/etcd-data"
  "$ETCD_BIN" --data-dir "$RUN/etcd-data" \
    --listen-client-urls http://127.0.0.1:2379 \
    --advertise-client-urls http://127.0.0.1:2379 \
    > "$RUN/etcd.log" 2>&1 &
  echo $! > "$RUN/etcd.pid"
  sleep 3
fi

echo ">> 编译服务..."
cd "$ROOT"
go build -o "$BIN/user-rpc.exe"         ./app/user/rpc
go build -o "$BIN/notification-rpc.exe" ./app/notification/rpc
go build -o "$BIN/product-rpc.exe"      ./app/product/rpc
go build -o "$BIN/order-rpc.exe"        ./app/order/rpc
go build -o "$BIN/payment-rpc.exe"      ./app/payment/rpc
go build -o "$BIN/user-api.exe"         ./app/user/api
go build -o "$BIN/gateway-api.exe"       ./app/gateway/api

echo ">> 启动 user.rpc (:8081)..."
( cd "$ROOT/app/user/rpc" && "$BIN/user-rpc.exe" -f etc/user.yaml > "$RUN/user-rpc.log" 2>&1 & echo $! > "$RUN/user-rpc.pid" )

echo ">> 启动 notification.rpc (:8005)..."
( cd "$ROOT/app/notification/rpc" && "$BIN/notification-rpc.exe" -f etc/notification.yaml > "$RUN/notification-rpc.log" 2>&1 & echo $! > "$RUN/notification-rpc.pid" )

echo ">> 启动 product.rpc (:8006)..."
( cd "$ROOT/app/product/rpc" && "$BIN/product-rpc.exe" -f etc/product.yaml > "$RUN/product-rpc.log" 2>&1 & echo $! > "$RUN/product-rpc.pid" )

echo ">> 启动 order.rpc (:8007)..."
( cd "$ROOT/app/order/rpc" && "$BIN/order-rpc.exe" -f etc/order.yaml > "$RUN/order-rpc.log" 2>&1 & echo $! > "$RUN/order-rpc.pid" )

echo ">> 启动 payment.rpc (:8008)..."
( cd "$ROOT/app/payment/rpc" && "$BIN/payment-rpc.exe" -f etc/payment.yaml > "$RUN/payment-rpc.log" 2>&1 & echo $! > "$RUN/payment-rpc.pid" )

echo ">> 等待 RPC 服务注册..."
sleep 3

echo ">> 启动 user-api (:8001)..."
( cd "$ROOT/app/user/api" && "$BIN/user-api.exe" -f etc/user-api.yaml > "$RUN/user-api.log" 2>&1 & echo $! > "$RUN/user-api.pid" )

echo ">> 启动 gateway-api (:8000)..."
( cd "$ROOT/app/gateway/api" && "$BIN/gateway-api.exe" -f etc/gateway-api.yaml > "$RUN/gateway-api.log" 2>&1 & echo $! > "$RUN/gateway-api.pid" )

sleep 4
echo ">> 全部启动完成。日志在 $RUN/*.log"
echo "   gateway-api : http://localhost:8000"
echo "   user-api    : http://localhost:8001"
echo "   停止：bash scripts/stop-all.sh"
