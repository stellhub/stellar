# Stellar Axis 全局请求 Header 与指标规范

本文档定义 `Stellar Axis（星轴）` 体系下所有中间件、所有业务服务统一遵循的请求 Header 规范与指标规范。

设计目标如下：

- 统一跨服务请求上下文语义
- 统一网关、服务、SDK、Agent 对请求头的识别方式
- 统一指标命名、标签边界与聚合口径
- 降低中间件之间字段映射成本，避免长期演进后语义漂移

## 适用范围

本规范适用于：

- 所有 HTTP / gRPC 网关
- 所有业务服务
- 所有 SDK
- 所有 Agent 与 Collector
- 所有日志、链路、指标、告警平台

## 第一部分：全局请求 Header 规范

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

## 第二部分：全局指标规范

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

## 第三部分：对平台与业务的约束

### 对网关的约束

- 负责生成或校验请求 ID
- 负责清洗不可信 Header
- 负责把标准链路头和平台扩展头稳定传入下游

### 对 SDK 的约束

- 默认识别标准链路头与 `X-Stellar-*` 扩展头
- 日志、链路、指标 SDK 应尽量共享同一份上下文语义
- 不得自行定义与平台规范冲突的新头或新指标语义

### 对业务应用的约束

- 业务应用应优先复用平台标准 Header
- 业务自定义 Header 应使用业务域前缀，不与平台头冲突
- 业务自定义指标应遵循统一命名和标签边界
