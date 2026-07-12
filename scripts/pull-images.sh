#!/usr/bin/env bash
# 弹性拉取基础设施镜像（Docker Hub 在本网络被墙，走 daocloud 镜像源并不断重试）。
# 拉全后自动 docker compose up -d。可反复运行，已就绪的镜像会跳过。
# 用法：bash scripts/pull-images.sh [总时长秒，默认1800]
export PATH="/c/Program Files/Docker/Docker/resources/bin:$PATH"
MIRROR=docker.m.daocloud.io
BUDGET="${1:-1800}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# canonical|mirror-path
IMAGES="
mysql:8.0|library/mysql:8.0
redis:7|library/redis:7
rabbitmq:3.13-management|library/rabbitmq:3.13-management
mongo:7|library/mongo:7
minio/minio:latest|minio/minio:latest
jaegertracing/all-in-one:1.57|jaegertracing/all-in-one:1.57
prom/prometheus:v2.53.0|prom/prometheus:v2.53.0
grafana/grafana:11.1.0|grafana/grafana:11.1.0
"

has_image() { docker image inspect "$1" >/dev/null 2>&1; }

while [ "$SECONDS" -lt "$BUDGET" ]; do
  missing=0
  echo "$IMAGES" | while IFS='|' read -r canonical path; do
    [ -z "$canonical" ] && continue
    has_image "$canonical" && continue
    echo ">> [$SECONDS s] pulling $canonical ..."
    if docker pull "$MIRROR/$path" >/dev/null 2>&1; then
      docker tag "$MIRROR/$path" "$canonical" && echo "   OK $canonical"
    else
      echo "   miss $canonical (稍后重试)"
    fi
  done
  # 统计还缺几个
  missing=0
  for pair in $IMAGES; do
    canonical="${pair%%|*}"
    [ -z "$canonical" ] && continue
    has_image "$canonical" || missing=$((missing+1))
  done
  if [ "$missing" -eq 0 ]; then
    echo ">> 全部镜像就绪，启动 compose ..."
    docker compose -f "$ROOT/deploy/docker/docker-compose.yml" up -d
    echo ">> DONE"
    exit 0
  fi
  echo ">> 还缺 $missing 个镜像，10s 后重试 ..."
  sleep 10
done
echo ">> 时间预算 ${BUDGET}s 用尽，仍有镜像未拉全。可再次运行本脚本继续。"
exit 1
