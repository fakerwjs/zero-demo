# zero-demo

基于 [go-zero](https://github.com/zeromicro/go-zero) 的微服务电商系统，采用 **APISIX + API Gateway + RPC 服务** 的分层架构。

## 架构概览

```
客户端 → APISIX Gateway (9080) → go-zero Gateway (8000) → RPC 服务
```

**流量走向**：
- **APISIX**：前置网关，负责统一入口、JWT 鉴权、IP 限流、CORS、路由转发
- **go-zero Gateway**：业务网关，负责业务路由、协议转换 (HTTP ↔ gRPC)、链路追踪

## 服务与端口

### 应用服务

| 服务 | 类型 | 端口 | 指标端口 | etcd Key | 说明 |
|------|------|------|----------|----------|------|
| **apisix** | Gateway | 9080 | 9180 | - | 前置网关，统一入口 |
| **apisix-dashboard** | Web | 9000 | - | - | APISIX 管理控制台 |
| gateway-api | HTTP | 8000 | 9100 | - | API 网关，业务路由 |
| user-api | HTTP | 8001 | 9101 | - | 用户独立 HTTP API |
| user.rpc | gRPC | 8081 | 9102 | user.rpc | 用户核心逻辑 |
| product.rpc | gRPC | 8006 | 9104 | product.rpc | 商品管理、库存扣减 |
| order.rpc | gRPC | 8007 | 9105 | order.rpc | 订单管理、超时取消 |
| payment.rpc | gRPC | 8008 | 9106 | payment.rpc | 支付管理 |
| notification.rpc | gRPC | 8005 | 9103 | notification.rpc | 通知服务 |

### 基础设施（docker-compose）

| 组件 | 端口 | 用途 | 访问地址 |
|------|------|------|----------|
| etcd | 2379 | 服务注册发现 | - |
| MySQL | 3306 | 主存储 | root/root123456 |
| Redis | 6379 | model 缓存 | - |
| RabbitMQ | 5672/15672 | 异步队列/管理台 | http://localhost:15672 (guest/guest) |
| MongoDB | 27017 | 文档存储 | - |
| MinIO | 9000/9001 | 对象存储/控制台 | http://localhost:9001 (minioadmin/minioadmin) |
| Jaeger | 16686/4317 | 链路追踪/OTLP | http://localhost:16686 |
| Prometheus | 9090 | 指标采集 | http://localhost:9090 |
| Grafana | 3000 | 监控大盘 | http://localhost:3000 (admin/admin) |
| Elasticsearch | 9200/9300 | 日志存储 | - |
| Kibana | 5601 | 日志可视化 | http://localhost:5601 |

## 目录结构

```
zero-demo/
├── app/
│   ├── gateway/api/           # API 网关（HTTP 入口）
│   ├── user/                  # 用户服务（API + RPC）
│   ├── product/rpc/           # 商品 RPC 服务
│   ├── order/rpc/             # 订单 RPC 服务
│   ├── payment/rpc/           # 支付 RPC 服务
│   └── notification/rpc/      # 通知 RPC 服务
├── pkg/                       # 公共库（JWT、MQ、MinIO、MongoDB）
├── deploy/                    # 部署配置
│   ├── docker/                # Docker Compose
│   ├── apisix/                # APISIX 配置
│   ├── mysql/init/            # 数据库初始化脚本
│   ├── prometheus/            # Prometheus 配置
│   └── grafana/               # Grafana 配置
├── scripts/                   # 启动脚本
├── docs/                      # 文档
├── go.mod
└── go.sum
```

## 快速开始

```bash
# 1. 启动全套基础设施（含 APISIX）
cd deploy/docker && docker compose up -d && cd ../..

# 2. 一键启动所有服务
bash scripts/run-all.sh

# 3. 通过 APISIX 访问 API
curl http://localhost:9080/api/v1/product/list

# 停止服务
bash scripts/stop-all.sh

# 停止基础设施
cd deploy/docker && docker compose down
```

## API 接口示例

**通过 APISIX 网关访问（推荐）**：

```bash
# 注册
curl -X POST http://localhost:9080/api/v1/user/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"123456","mobile":"13800000000"}'

# 登录（返回 accessToken）
curl -X POST http://localhost:9080/api/v1/user/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"123456"}'

# 查询用户信息（需带 JWT）
curl http://localhost:9080/api/v1/user/info \
  -H 'Authorization: Bearer <accessToken>'
```

## 管理界面

| 服务 | 地址 | 用户名/密码 |
|------|------|-------------|
| APISIX Dashboard | http://localhost:9000 | admin/admin |
| Grafana | http://localhost:3000 | admin/admin |
| Jaeger | http://localhost:16686 | - |
| RabbitMQ | http://localhost:15672 | guest/guest |
| MinIO | http://localhost:9001 | minioadmin/minioadmin |
| Kibana | http://localhost:5601 | - |

## 技术栈

- **框架**: go-zero v1.10.2
- **语言**: Go 1.25+
- **RPC**: gRPC + Protocol Buffers
- **网关**: APISIX + go-zero Gateway
- **服务发现**: etcd
- **数据库**: MySQL 8.0 + Redis 7 + MongoDB 7
- **消息队列**: RabbitMQ 3.13
- **对象存储**: MinIO
- **可观测性**: Jaeger (链路追踪) + Prometheus + Grafana (监控) + ELK (日志)

## 文档

| 文档 | 说明 |
|------|------|
| [docs/1-项目架构梳理.md](docs/1-项目架构梳理.md) | 完整架构说明 |
| [docs/2-快速启动指南.md](docs/2-快速启动指南.md) | 快速上手启动 |
| [docs/3-部署文档.md](docs/3-部署文档.md) | 详细部署说明 |
| [docs/4-执行文档.md](docs/4-执行文档.md) | 核心业务流程 |
| [docs/5-基础设施使用文档.md](docs/5-基础设施使用文档.md) | 基础设施配置参考 |
| [docs/6-代码审查报告.md](docs/6-代码审查报告.md) | 质量检查报告 |
| [docs/7-优化需求清单.md](docs/7-优化需求清单.md) | 改进计划 |
| [docs/8-优化实施报告.md](docs/8-优化实施报告.md) | 实施记录 |
| [docs/9-RabbitMQ使用教程.md](docs/9-RabbitMQ使用教程.md) | RabbitMQ 详细使用教程 |
