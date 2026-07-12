#!/usr/bin/env bash
# 停止 run-all.sh 启动的所有进程
# 用法：bash scripts/stop-all.sh
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RUN="$ROOT/.run"

for name in gateway-api user-api user-rpc notification-rpc product-rpc order-rpc payment-rpc etcd; do
  pidfile="$RUN/$name.pid"
  if [ -f "$pidfile" ]; then
    pid="$(cat "$pidfile")"
    if kill "$pid" 2>/dev/null; then
      echo ">> 已停止 $name (pid=$pid)"
    else
      echo ">> $name (pid=$pid) 未在运行"
    fi
    rm -f "$pidfile"
  fi
done

# 兜底：按端口清理可能残留的进程（Windows）
for port in 8000 8001 8005 8006 8007 8008 8081 2379; do
  pid="$(netstat -ano 2>/dev/null | grep ":$port " | grep LISTENING | awk '{print $NF}' | head -1)"
  [ -n "$pid" ] && taskkill //PID "$pid" //F >/dev/null 2>&1 && echo ">> 端口 $port 残留进程 $pid 已清理"
done
echo ">> 完成。"
