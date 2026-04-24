---
title: 可观测规范
outline: deep
---

# 可观测规范

本文档统一定义 `Stell Hub（星际枢纽）` 体系下所有中间件、所有业务服务共同遵循的可观测规范，覆盖基础环境变量、全局请求上下文、链路传递与指标语义。

这份规范用于解决以下问题：

- 平台统一注入应用身份、部署拓扑和观测上下文
- 所有中间件在日志、链路、指标之间共享一致的基础语义
- 网关、服务、SDK、Agent 与 Collector 在请求传播和指标聚合上保持同构
- 降低产品间字段映射成本，避免长期演进后的语义漂移

## 规范定位

本规范是全局运行规范中与观测体系直接相关的统一约束，适用于：

- 所有业务服务
- 所有中间件 SDK
- 所有 HTTP / gRPC 网关
- 所有 Agent、Sidecar、Collector
- 所有日志、链路、指标、告警平台
- 所有平台注入系统、部署脚本与镜像装机脚本

## 第一部分：基础环境变量规范

### 命名原则

- 全局统一前缀使用 `STELLAR_`
- 全部使用大写英文与下划线命名
- 变量语义必须长期稳定，不允许在单个产品内重定义
- 基础元数据与产品级配置分层，不混用命名空间

### 优先级原则

所有中间件统一建议采用以下优先级：

1. 代码显式配置
2. 产品级环境变量
3. 全局 `STELLAR_*` 环境变量
4. 中间件默认值

说明如下：

- `STELLAR_*` 是平台统一注入的全局基础元数据
- 产品级变量只允许覆盖本产品局部行为
- 不允许产品级变量改变全局基础元数据的定义

### 必选基础变量

以下变量建议作为企业级运行环境的最小必选集合：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_APP_NAME` | `user-service` | 应用或服务名称 |
| `STELLAR_APP_NAMESPACE` | `stellar.trade` | 逻辑业务域或应用命名空间 |
| `STELLAR_APP_VERSION` | `1.4.2` | 当前发布版本 |
| `STELLAR_APP_INSTANCE_ID` | `user-service-7f6d9d6d7b-2xk9p` | 当前运行实例唯一标识 |
| `STELLAR_ENV` | `dev` / `test` / `prod` | 部署环境 |

### 推荐拓扑变量

以下变量建议由平台统一注入，供日志、链路、指标、配置、网关等能力共同使用：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_CLUSTER` | `cluster-sh-prod-01` | 集群标识 |
| `STELLAR_REGION` | `cn-east-1` | 区域标识 |
| `STELLAR_ZONE` | `cn-east-1a` | 可用区标识 |
| `STELLAR_IDC` | `sh-a` | 机房或园区标识 |
| `STELLAR_HOST_NAME` | `node-01` | 主机名 |
| `STELLAR_HOST_IP` | `10.10.0.11` | 主机 IP |

### Kubernetes 推荐变量

如果运行在 Kubernetes 中，建议额外统一注入以下变量：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_NODE_NAME` | `worker-node-01` | 节点名称 |
| `STELLAR_K8S_NAMESPACE` | `trade` | Kubernetes Namespace |
| `STELLAR_POD_NAME` | `user-service-7f6d9d6d7b-2xk9p` | Pod 名称 |
| `STELLAR_POD_IP` | `172.20.10.23` | Pod IP |
| `STELLAR_CONTAINER_NAME` | `app` | 容器名称 |

### 统一语义约束

- `APP_NAME`
  表示业务应用或工作负载的稳定名称，不使用主机名、Pod 名或镜像名替代
- `APP_NAMESPACE`
  表示逻辑业务域，不等同于 Kubernetes Namespace
- `APP_VERSION`
  表示本次发布版本，不使用 Git 分支名替代
- `APP_INSTANCE_ID`
  表示单个运行实例唯一标识，允许使用 Pod 名、实例 ID 或平台生成值
- `ENV`
  只表示部署环境，不混入区域、租户、集群或机房信息
- `CLUSTER / REGION / ZONE / IDC`
  只表示基础设施拓扑，不与环境语义混用

### 产品级扩展约束

在 `STELLAR_*` 之外，各产品可以保留自己的产品级环境变量，例如：

- 日志平台：`STELLSPEC_*`
- 链路平台：`STELLTRACE_*`
- 配置中心：`STELLNULA_*`
- 服务治理：`STELLORBIT_*`

但产品级变量只允许描述：

- 本产品的接入地址
- 本产品的协议与格式
- 本产品的采样、级别、超时、批量参数
- 对基础元数据的局部覆盖

不允许产品级变量重新发明：

- 应用名语义
- 环境语义
- 区域语义
- 集群语义

### 平台注入建议

推荐统一由平台在如下阶段注入：

- Kubernetes Admission Webhook
- Helm / Kustomize 模板
- 宿主机装机脚本
- 容器运行时入口脚本
- 发布平台注入流程

推荐来源包括：

- `metadata.name`
- `metadata.namespace`
- `spec.nodeName`
- `status.podIP`
- `status.hostIP`
- 发布系统生成的版本号
- CMDB / 资源平台中的区域、机房、集群信息

## 第二部分：全局请求上下文规范

### Header 设计原则

- 优先复用标准协议头
- 平台扩展头统一使用 `X-Stellar-` 前缀
- Header 语义必须稳定，不因单个产品变化而变化
- Header 只承载请求上下文，不承载大体积业务数据

### 标准透传头

以下标准头建议原样透传并优先复用：

| Header | 说明 |
| :--- | :--- |
| `traceparent` | W3C Trace Context 主链路头 |
| `tracestate` | W3C Trace Context 扩展状态 |
| `baggage` | W3C Baggage 扩展上下文 |
| `X-Request-Id` | 请求唯一 ID，兼容传统生态 |
| `X-Forwarded-For` | 客户端来源 IP 链 |
| `X-Forwarded-Proto` | 原始请求协议 |
| `X-Forwarded-Host` | 原始请求主机 |
| `X-Real-IP` | 客户端真实 IP |

### Stellar 平台扩展头

以下平台扩展头建议所有业务和中间件统一识别：

| Header | 示例 | 说明 |
| :--- | :--- | :--- |
| `X-Stellar-Request-Id` | `req-20260416-000001` | 平台统一请求 ID |
| `X-Stellar-Session-Id` | `sess-8fa12c` | 会话 ID |
| `X-Stellar-User-Id` | `u-1001` | 用户 ID |
| `X-Stellar-Tenant-Id` | `t-01` | 租户 ID |
| `X-Stellar-Device-Id` | `d-88f1` | 设备 ID |
| `X-Stellar-Client-Ip` | `10.20.30.40` | 平台归一后的客户端 IP |
| `X-Stellar-Source-App` | `stellgate-external` | 来源应用 |
| `X-Stellar-Source-Service` | `gateway` | 来源服务 |
| `X-Stellar-Source-Region` | `cn-east-1` | 来源区域 |
| `X-Stellar-Env` | `prod` | 来源环境 |
| `X-Stellar-Gray-Tag` | `gray-a` | 灰度标识 |
| `X-Stellar-Canary-Tag` | `canary-v2` | 金丝雀发布标识 |

### Header 统一约束

- 如果 `traceparent` 已存在，必须优先沿用，不得随意重建主链路
- 如果 `X-Stellar-Request-Id` 不存在，最外层网关必须生成
- `X-Request-Id` 与 `X-Stellar-Request-Id` 建议保持一致或建立确定性映射
- `User-Id`、`Tenant-Id` 等身份头必须经过鉴权后注入，不允许由不可信外部请求直接覆盖
- 平台扩展头进入外部边界时，应由网关进行白名单透传与重建

### gRPC 映射建议

如果是 gRPC 调用，建议映射到 metadata 中保持同名小写键，例如：

- `x-stellar-request-id`
- `x-stellar-tenant-id`
- `x-stellar-source-app`

### 对日志与链路的约束

所有日志 SDK、链路 SDK 与网关应尽量将以下 Header 自动映射为上下文字段：

- `traceparent`
- `X-Stellar-Request-Id`
- `X-Stellar-Tenant-Id`
- `X-Stellar-User-Id`
- `X-Stellar-Gray-Tag`

## 第三部分：全局指标规范

### 指标设计原则

- 优先使用单调、可聚合、可下钻的指标模型
- 指标命名统一使用小写英文与下划线
- 指标名中明确业务对象与语义，不使用含糊缩写
- 高基数信息不得默认进入指标标签

### 指标命名模式

推荐统一采用：

> `{domain}_{object}_{metric}_{unit}`

常见示例：

- `http_server_requests_total`
- `http_server_request_duration_ms`
- `rpc_client_requests_total`
- `mq_consumer_messages_total`
- `db_client_connections_active`

### 指标类型建议

| 类型 | 适用场景 |
| :--- | :--- |
| Counter | 请求总量、错误总量、消息总量 |
| Gauge | 当前连接数、队列积压数、线程池活跃数 |
| Histogram | 延迟、大小、耗时分布 |

### 必选通用标签

以下标签建议作为全平台默认低基数标签：

| 标签 | 说明 |
| :--- | :--- |
| `app_name` | 应用名 |
| `app_namespace` | 应用命名空间 |
| `app_version` | 应用版本 |
| `env` | 环境 |
| `cluster` | 集群 |
| `region` | 区域 |
| `zone` | 可用区 |

### 场景标签建议

不同场景建议增加以下低基数标签：

| 场景 | 推荐标签 |
| :--- | :--- |
| HTTP 服务端 | `method`、`route`、`status_code` |
| gRPC 服务端 | `rpc_system`、`rpc_service`、`rpc_method`、`rpc_code` |
| 数据库 | `db_system`、`db_operation`、`db_name` |
| MQ | `mq_system`、`topic`、`consumer_group` |

### 禁止默认进入指标标签的字段

以下字段通常属于高基数信息，不应默认作为指标标签：

- `user_id`
- `tenant_id`
- `request_id`
- `session_id`
- `device_id`
- `pod_name`
- `host_ip`
- `trace_id`
- `span_id`
- 原始 URL 全路径

如确需观测，应优先进入：

- 日志字段
- Trace 属性
- 离线分析数据

### HTTP 指标最小集合

所有网关与 HTTP 服务建议至少提供：

| 指标名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `http_server_requests_total` | Counter | 请求总数 |
| `http_server_request_errors_total` | Counter | 错误总数 |
| `http_server_request_duration_ms` | Histogram | 请求耗时 |
| `http_server_inflight_requests` | Gauge | 当前处理中请求数 |

### RPC 指标最小集合

所有 RPC 服务建议至少提供：

| 指标名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `rpc_server_requests_total` | Counter | RPC 请求总数 |
| `rpc_server_request_errors_total` | Counter | RPC 错误总数 |
| `rpc_server_request_duration_ms` | Histogram | RPC 耗时 |

### 中间件客户端指标最小集合

数据库、缓存、MQ、配置中心、注册中心等客户端建议至少提供：

| 指标名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `middleware_client_requests_total` | Counter | 调用总数 |
| `middleware_client_errors_total` | Counter | 调用错误数 |
| `middleware_client_duration_ms` | Histogram | 调用耗时 |

### 指标聚合约束

- 所有通用仪表盘应优先依赖低基数标签
- 版本、环境、区域、集群必须可以稳定聚合
- 指标单位必须进入指标名或统一元数据中，不允许模糊表达
- 同一语义的指标在不同产品中必须统一命名

## 第四部分：平台与业务落地约束

### 对网关的约束

- 负责生成或校验请求 ID
- 负责清洗不可信 Header
- 负责把标准链路头和平台扩展头稳定传入下游

### 对 SDK 的约束

- 默认识别标准链路头与 `X-Stellar-*` 扩展头
- 默认读取 `STELLAR_*` 基础元数据，并映射到日志、链路、指标语义
- 日志、链路、指标 SDK 应尽量共享同一份上下文语义
- 允许产品级环境变量做本地覆盖，但不得破坏全局语义
- 不得自行定义与平台规范冲突的新头或新指标语义

### 对业务应用的约束

- 业务应用应优先复用平台标准 Header
- 业务自定义 Header 应使用业务域前缀，不与平台头冲突
- 业务自定义指标应遵循统一命名和标签边界

### 落地建议

- 平台统一注入 `STELLAR_*` 元数据，业务侧只做局部补充
- 网关统一负责链路头、请求 ID 和不可信 Header 清洗
- SDK 在日志、链路、指标三个面向上共享同一份上下文模型
- 仪表盘与告警规则优先依赖低基数标签构建，以保证稳定聚合
