---
title: 可观测规范
category: 体系规范
summary: 以 OpenTelemetry 与 Kubernetes 为主标准，统一定义资源语义、日志 KV 视角、上下文传播、客户端与服务端指标模型及平台落地职责，Stellar 仅承担最小补差角色。
tags:
  - 可观测性
  - OpenTelemetry
  - Kubernetes
  - 规范设计
  - 平台工程
readingDirection: 适合在准备统一日志、链路、指标语义并推动平台化落地时优先阅读。
outline: deep
---

# 可观测规范

## 摘要

本文给出一份面向企业平台的可观测统一规范。该规范以 OpenTelemetry 资源语义、日志数据模型、上下文传播与指标语义约定为主线，以 Kubernetes 原生元数据为事实来源，以 `STELLAR_*` 与 `X-Stellar-*` 作为最小补差层，目标是在不破坏开源标准的前提下，建立统一、稳定、可治理的日志、链路与指标协同模型。本文从资源身份、日志 KV 视角、HTTP 与 gRPC 传播、客户端与服务端指标模型、平台分层职责与验收原则六个方面展开，给出一条“标准优先、事实复用、补差最小”的实现路径。

## 关键词

OpenTelemetry；Kubernetes；Resource Semantics；Logs Data Model；Trace Context；Baggage；Metrics；Platform Engineering

## 1. 引言

企业在建设可观测平台时，往往会先从统一环境变量、统一请求头和统一指标名入手。这种方式在项目早期能快速降低接入成本，但如果平台定义与 OpenTelemetry、Kubernetes 并行的命名体系，就会逐步引入新的复杂度：同一语义出现多个名字，同一字段在不同系统中含义不同，采集链路与分析平台之间需要长期维护额外映射关系，最终导致平台治理能力增强的同时，生态兼容性与长期可维护性下降。

因此，本规范不再以“平台自定义字段”为主，而是将 OpenTelemetry 视为主标准，将 Kubernetes 视为主事实来源，将平台自定义字段压缩到最小范围。对于资源身份，优先对齐 `service.*`、`k8s.*`、`host.*`、`cloud.*`；对于日志 KV，优先对齐 OpenTelemetry Logs Data Model 与 HTTP / RPC 语义属性；对于链路传播，优先对齐 W3C Trace Context 与 W3C Baggage；对于 HTTP、RPC、gRPC 等指标，优先复用 OpenTelemetry 语义约定中已有的标准指标与标准属性。只有在开源标准尚未覆盖、或平台治理确有刚性需求时，才允许使用 `STELLAR_*` 与 `X-Stellar-*` 补齐差异。

## 2. 规范目标与标准基线

### 2.1 规范目标

本规范面向以下对象：

- 所有业务服务
- 所有中间件 SDK
- 所有 HTTP / gRPC 网关
- 所有 Agent、Sidecar、Collector
- 所有日志、链路、指标、告警平台
- 所有平台注入系统、部署脚本与镜像入口

规范目标如下：

1. 统一资源身份语义，避免同一事实存在多个并行定义。
2. 统一日志 KV 视角，避免客户端与服务端观察边界混淆。
3. 统一上下文传播语义，避免链路头和平台头重复造轮子。
4. 统一指标命名与属性边界，避免平台私有指标体系与标准指标体系并行。
5. 明确平台各层职责，避免发布平台、网关、SDK、Collector 重复实现或相互覆盖。

### 2.2 标准基线

本规范优先对齐以下标准：

- OpenTelemetry Resources
- OpenTelemetry Service Semantic Conventions
- OpenTelemetry HTTP Semantic Conventions
- OpenTelemetry RPC / gRPC Semantic Conventions
- OpenTelemetry Logs Data Model
- OpenTelemetry Propagators API
- OpenTelemetry Baggage
- Kubernetes 原生对象元数据与状态字段

推荐阅读：

1. [OpenTelemetry Resources](https://opentelemetry.io/docs/concepts/resources/)
2. [OpenTelemetry Service Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/resource/service/)
3. [OpenTelemetry Kubernetes Resource Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/resource/k8s/)
4. [OpenTelemetry Host Resource Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/resource/host/)
5. [OpenTelemetry HTTP Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/http/)
6. [OpenTelemetry RPC Metrics](https://opentelemetry.io/docs/specs/semconv/rpc/rpc-metrics/)
7. [OpenTelemetry gRPC Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/rpc/grpc/)
8. [OpenTelemetry Logs Data Model](https://opentelemetry.io/docs/specs/otel/logs/data-model/)
9. [Trace Context in non-OTLP Log Formats](https://opentelemetry.io/docs/specs/otel/compatibility/logging_trace_context/)
10. [OpenTelemetry Propagators API](https://opentelemetry.io/docs/specs/otel/context/api-propagators/)
11. [OpenTelemetry Baggage](https://opentelemetry.io/docs/concepts/signals/baggage/)

## 3. 资源语义模型

### 3.1 总体原则

资源语义实现遵循以下原则：

1. 资源身份优先使用 OpenTelemetry 标准属性表达。
2. Kubernetes 已提供的元数据优先作为事实来源，不额外创造同义事实。
3. `STELLAR_*` 不是主定义，只是最小补差或桥接入口。
4. 同一语义只能有一个主定义，不允许在不同产品中重解释。

### 3.2 资源属性优先级

所有 SDK、Agent、Collector 与平台注入系统，统一建议采用以下优先级：

1. 代码或 SDK 显式设置的 OpenTelemetry 标准资源属性
2. `OTEL_SERVICE_NAME`、`OTEL_RESOURCE_ATTRIBUTES` 等标准 OTel 配置
3. Kubernetes / Host / Cloud Resource Detector 自动识别结果
4. 平台补差注入的 `STELLAR_*`
5. 中间件默认值

这意味着，如果某个运行时已经可以直接提供标准 OTel Resource，就不应先写入 `STELLAR_*` 再进行二次翻译。

### 3.3 应用身份语义

应用身份字段应优先对齐以下标准属性：

| 语义 | 标准属性 | 示例 | 说明 |
| :--- | :--- | :--- | :--- |
| 服务名 | `service.name` | `user-service` | 逻辑服务名称 |
| 服务命名空间 | `service.namespace` | `stellar.trade` | 逻辑业务域或服务组 |
| 服务版本 | `service.version` | `1.4.2` | 当前发布版本 |
| 实例 ID | `service.instance.id` | `627cc493-f310-47de-96bd-71410b7dec09` | 单个服务实例唯一标识 |
| 部署环境 | `deployment.environment.name` | `prod` | 部署环境 |

统一约束如下：

1. `service.name` 表示稳定服务名，不使用 Pod 名、主机名或镜像名替代。
2. `service.namespace` 表示逻辑业务域，不等同于 `k8s.namespace.name`。
3. `service.version` 表示发布版本，不使用 Git 分支名替代。
4. `deployment.environment.name` 只表示环境，不混入区域、集群、租户或机房信息。
5. `service.instance.id` 必须在同一 `service.namespace + service.name` 维度内唯一。

### 3.4 实例 ID 约束

`service.instance.id` 是实例身份，不是网络地址，也不是拓扑位置。因此：

1. 不建议直接使用 `Pod IP` 作为 `service.instance.id`。
2. 不建议直接使用 `host.ip` 作为 `service.instance.id`。
3. 推荐优先使用平台生成的实例唯一值，或应用启动时生成并稳定持有的 UUID。
4. 若一个 Pod 明确只承载一个服务实例，且需要简单落地，可以使用 `Pod UID` 或 `Pod Name` 作为折中。

### 3.5 基础设施与 Kubernetes 语义

基础设施与 Kubernetes 相关语义应优先对齐以下标准属性：

| 语义 | 标准属性 | 推荐来源 |
| :--- | :--- | :--- |
| 主机名 | `host.name` | Host Detector |
| 主机 IP | `host.ip` | Host Detector、`status.hostIP` |
| 区域 | `cloud.region` | 云平台元数据、资源平台 |
| 可用区 | `cloud.availability_zone` | 云平台元数据、资源平台 |
| 集群名 | `k8s.cluster.name` | 集群配置、资源平台 |
| Namespace | `k8s.namespace.name` | `metadata.namespace` |
| Pod 名 | `k8s.pod.name` | `metadata.name` |
| Pod UID | `k8s.pod.uid` | `metadata.uid` |
| Pod IP | `k8s.pod.ip` | `status.podIP` |
| 节点名 | `k8s.node.name` | `spec.nodeName` |
| 容器名 | `k8s.container.name` | Pod 模板容器名 |

统一约束如下：

1. `host.name` 表达主机语义，`k8s.node.name` 表达 Kubernetes 节点语义。
2. 两者在 Kubernetes 场景中可能同值，但语义上不应混为一个字段。
3. `k8s.pod.ip` 表达 Pod 网络地址，不能提升为实例身份字段。
4. `k8s.container.name` 表示 Pod 规范中的容器名，不等同于运行时全局容器 ID。

### 3.6 `STELLAR_*` 的边界

`STELLAR_*` 只用于以下场景：

1. 非 OTel 组件需要统一环境变量入口。
2. 当前运行时无法直接设置标准 OTel Resource，需要桥接输入。
3. 平台内部存在 OTel 尚未稳定覆盖的统一拓扑字段，例如 `STELLAR_IDC`。

推荐映射关系如下：

| `STELLAR_*` | 首选标准属性 |
| :--- | :--- |
| `STELLAR_APP_NAME` | `service.name` |
| `STELLAR_APP_NAMESPACE` | `service.namespace` |
| `STELLAR_APP_VERSION` | `service.version` |
| `STELLAR_APP_INSTANCE_ID` | `service.instance.id` |
| `STELLAR_ENV` | `deployment.environment.name` |
| `STELLAR_CLUSTER` | `k8s.cluster.name` |
| `STELLAR_REGION` | `cloud.region` |
| `STELLAR_ZONE` | `cloud.availability_zone` |
| `STELLAR_HOST_NAME` | `host.name` |
| `STELLAR_HOST_IP` | `host.ip` |
| `STELLAR_NODE_NAME` | `k8s.node.name` |
| `STELLAR_K8S_NAMESPACE` | `k8s.namespace.name` |
| `STELLAR_POD_NAME` | `k8s.pod.name` |
| `STELLAR_POD_IP` | `k8s.pod.ip` |
| `STELLAR_CONTAINER_NAME` | `k8s.container.name` |

新系统应优先读取标准 OTel 属性，而不是优先依赖 `STELLAR_*`。

## 4. 日志 KV 视角模型

### 4.1 说明

OpenTelemetry 日志数据模型明确了 `Timestamp`、`ObservedTimestamp`、`TraceId`、`SpanId`、`TraceFlags`、`SeverityText`、`SeverityNumber`、`Body`、`Resource`、`Attributes` 等基础字段，但没有直接提供一份必须统一的“客户端日志 KV 模板”或“服务端日志 KV 模板”。因此，本节是在 OpenTelemetry Logs Data Model 以及 HTTP / RPC 语义约定基础上的实现性推导。

换句话说，下面关于“客户端应该看到哪些字段”“服务端应该看到哪些字段”的建议，是基于标准语义做出的平台约束，而不是 OpenTelemetry 文档中逐字段给出的原样模板。

### 4.2 公共日志基础字段

无论客户端还是服务端，单条结构化日志都应优先具备以下基础字段：

| 类别 | 字段 | 说明 |
| :--- | :--- | :--- |
| 时间 | `Timestamp` | 事件发生时间 |
| 采集时间 | `ObservedTimestamp` | 被 SDK / Collector 观测到的时间 |
| 严重级别 | `SeverityText` / `SeverityNumber` | 日志等级 |
| 消息体 | `Body` | 日志正文 |
| 链路关联 | `TraceId` / `SpanId` / `TraceFlags` | 与 Trace 相关联的上下文 |
| 资源身份 | `service.name` | 当前记录日志的服务 |
| 资源身份 | `service.namespace` | 当前记录日志的服务命名空间 |
| 资源身份 | `service.version` | 当前记录日志的服务版本 |
| 资源身份 | `service.instance.id` | 当前记录日志的实例 ID |
| 资源身份 | `deployment.environment.name` | 当前记录日志的环境 |

若运行在 Kubernetes 中，还应通过 Resource 或日志属性稳定带出：

- `k8s.namespace.name`
- `k8s.pod.name`
- `k8s.pod.uid`
- `k8s.node.name`

应用原生日志字段可以继续保留 `level`，但进入平台统一日志模型时，必须规范化映射为 `SeverityText` 与 `SeverityNumber`。

### 4.3 客户端日志 KV 视角

客户端日志描述的是“当前服务向外发起一次调用”时，调用方自己看到的事实。推荐关注以下字段：

| 类别 | 字段 | 说明 |
| :--- | :--- | :--- |
| 当前服务身份 | `service.*` | 记录日志的一方，即调用方自身 |
| 请求方法 | `http.request.method` 或 `rpc.method` | 发起的调用方法 |
| 目标地址 | `server.address` / `server.port` | 客户端视角下的目标地址 |
| 协议细节 | `url.scheme`、`network.protocol.name`、`network.protocol.version` | 调用协议信息 |
| 返回结果 | `http.response.status_code` 或 `rpc.response.status_code` | 返回状态 |
| 错误信息 | `error.type` | 客户端侧观测到的低基数错误类型 |
| 平台补差 | `X-Stellar-Request-Id` | 平台请求 ID |
| 平台补差 | `X-Stellar-Tenant-Id`、`X-Stellar-User-Id` | 经鉴权后透传的业务上下文 |

客户端日志不应默认记录以下被调方内部字段：

- 被调方的 `service.instance.id`
- 被调方的 `k8s.pod.name`
- 被调方的 `k8s.node.name`
- 被调方的 `host.ip`

原因在于这些字段属于服务端内部运行事实，客户端既不稳定掌握，也不应把它们误当成“当前调用方观察到的事实”。

### 4.4 服务端日志 KV 视角

服务端日志描述的是“当前服务收到并处理一次请求”时，被调方自己看到的事实。推荐关注以下字段：

| 类别 | 字段 | 说明 |
| :--- | :--- | :--- |
| 当前服务身份 | `service.*` | 记录日志的一方，即被调方自身 |
| 入口方法 | `http.request.method` 或 `rpc.method` | 收到的请求方法 |
| 路由语义 | `http.route` | 服务端路由模板，必须保持低基数 |
| 本端地址 | `server.address` / `server.port` | 当前服务监听地址 |
| 返回结果 | `http.response.status_code` 或 `rpc.response.status_code` | 最终返回状态 |
| 错误信息 | `error.type` | 服务端侧观测到的低基数错误类型 |
| 来源补充 | `X-Stellar-Source-Service` | 若存在，表示调用方服务标识 |
| 平台补差 | `X-Stellar-Request-Id`、`X-Stellar-Tenant-Id`、`X-Stellar-User-Id` | 统一透传后的上下文 |

服务端日志不应默认去记录“后续下游被调方”的内部字段。如果当前服务在处理入口请求后继续访问别的服务，那些信息应记录在新的客户端日志中，而不是混入本次入口服务端日志。

### 4.5 客户端与服务端日志边界

为了避免视角混淆，统一执行以下规则：

1. 调用方日志不记录被调方实例级内部拓扑。
2. 服务端入口日志不记录当前请求后续下游调用的内部拓扑。
3. 若同一请求同时存在入口处理和下游调用，应分别形成“服务端日志”和“客户端日志”两组记录。
4. `X-Stellar-Source-Service` 是服务端用来识别调用方的补充字段，不是客户端用来记录自己身份的主字段；客户端自身身份应由 `service.name` 表达。

## 5. 上下文传播模型

### 5.1 总体原则

上下文传播以 OpenTelemetry Propagators API 为主标准。默认传播组合应为：

- `traceparent`
- `tracestate`
- `baggage`

因此：

1. HTTP 请求头的主传播协议不应自定义。
2. gRPC metadata 的主传播协议也不应自定义，只是载体由 HTTP Header 变为 metadata。
3. `X-Stellar-*` 只能承载标准传播头之外的最小平台补差上下文。

### 5.2 HTTP 传播语义

HTTP 侧平台实现应满足以下要求：

1. 网关、Sidecar、SDK 必须优先提取和注入 `traceparent`、`tracestate`、`baggage`。
2. 可保留 `X-Request-Id` 作为请求标识补充，但其地位低于 W3C Trace Context。
3. 若某些业务上下文适合跨服务透传，应优先评估是否进入 `baggage`。
4. 只有在不适合进入 `baggage` 时，才使用 `X-Stellar-*` 补差。

HTTP span 与指标属性应优先复用以下标准字段：

- `http.request.method`
- `http.response.status_code`
- `http.route`
- `url.scheme`
- `url.path`
- `url.query`
- `server.address`
- `server.port`
- `network.protocol.name`
- `network.protocol.version`

### 5.3 gRPC 传播语义

gRPC 侧平台实现应满足以下要求：

1. 标准传播上下文仍使用 W3C Trace Context 与 Baggage，只是通过 metadata 承载。
2. 自定义 metadata 键必须使用小写形式。
3. 若需要采集 request / response metadata，不应默认全量采集，而应使用白名单显式开启。

gRPC / RPC span 与指标属性应优先复用以下标准字段：

- `rpc.system.name`
- `rpc.method`
- `rpc.response.status_code`
- `server.address`
- `server.port`
- `error.type`

对 gRPC 场景，`rpc.system.name` 应固定为 `grpc`。

### 5.4 `X-Stellar-*` 最小集合

平台建议保留的最小补差头如下：

| Header | 用途 | 说明 |
| :--- | :--- | :--- |
| `X-Stellar-Request-Id` | 平台请求 ID | 可与 `X-Request-Id` 保持一致 |
| `X-Stellar-Session-Id` | 会话上下文 | 仅在业务确有需要时透传 |
| `X-Stellar-User-Id` | 用户上下文 | 必须经鉴权后注入 |
| `X-Stellar-Tenant-Id` | 租户上下文 | 必须经鉴权后注入 |
| `X-Stellar-Device-Id` | 设备上下文 | 仅在业务确有需要时透传 |
| `X-Stellar-Client-Ip` | 平台归一后的客户端 IP | 不替代标准网络属性 |
| `X-Stellar-Source-Service` | 调用方服务标识 | 应对齐调用方 `service.name` |
| `X-Stellar-Source-Region` | 来源区域 | 用于平台治理补差 |
| `X-Stellar-Env` | 来源环境 | 用于平台治理补差 |
| `X-Stellar-Gray-Tag` | 灰度标记 | 用于流量治理 |
| `X-Stellar-Canary-Tag` | 金丝雀标记 | 用于发布治理 |

统一约束如下：

1. 新规范不定义 `X-Stellar-Source-App`。
2. 对调用方身份，优先复用资源语义 `service.name`、`service.namespace`。
3. 外部边界进入内部链路前，所有 `X-Stellar-*` 都必须经过网关清洗、校验或重建。

## 6. 指标模型

### 6.1 总体原则

指标建模遵循以下原则：

1. 已有 OpenTelemetry 标准语义的场景，必须直接复用标准指标名。
2. 指标名、单位、属性集合应尽量直接复用 OTel 语义约定。
3. 不为 HTTP、RPC、gRPC、数据库、消息等已定义场景再发明平台私有指标名。
4. 高基数信息不得默认进入指标属性。
5. 客户端指标与服务端指标必须从不同视角建模，不能互相混入。

### 6.2 客户端 Metrics 视角

客户端 Metrics 描述的是“当前服务对外部依赖发起调用时的表现”。它回答的问题包括：

- 我调用下游是否成功
- 我调用下游用了多久
- 我依赖的目标地址或协议是否异常

客户端指标的主语始终是“调用方自己”，而不是“被调方内部状态”。

#### 6.2.1 HTTP 客户端 Metrics

HTTP 客户端至少应对齐以下标准指标：

| 指标名 | 类型 | 单位 | 地位 |
| :--- | :--- | :--- | :--- |
| `http.client.request.duration` | Histogram | `s` | 必选 |
| `http.client.active_requests` | UpDownCounter | `{request}` | 可选 |

推荐属性集合如下：

- `http.request.method`
- `http.response.status_code`
- `server.address`
- `server.port`
- `url.scheme`
- `network.protocol.name`
- `network.protocol.version`
- `error.type`

客户端 HTTP 指标不应默认带出：

- `http.route`
- 被调方 `service.instance.id`
- 被调方 `k8s.pod.name`
- 被调方 `host.ip`

其中，`http.route` 是服务端路由语义，实例级拓扑属于被调方内部事实。

#### 6.2.2 RPC / gRPC 客户端 Metrics

RPC / gRPC 客户端至少应对齐以下标准指标：

| 指标名 | 类型 | 单位 | 地位 |
| :--- | :--- | :--- | :--- |
| `rpc.client.call.duration` | Histogram | `s` | 必选 |

推荐属性集合如下：

- `rpc.system.name`
- `rpc.method`
- `rpc.response.status_code`
- `server.address`
- `server.port`
- `error.type`

统一约束如下：

1. gRPC 调用固定设置 `rpc.system.name = grpc`。
2. `rpc.method` 记录客户端视角的调用方法。
3. 客户端指标不记录服务端内部实例拓扑。

### 6.3 服务端 Metrics 视角

服务端 Metrics 描述的是“当前服务作为被调方时，对入口请求的处理表现”。它回答的问题包括：

- 我作为服务端收到了多少请求
- 我的入口请求耗时怎样
- 我的路由、方法、状态码分布如何

服务端指标的主语始终是“当前被调服务自己”，而不是它正在访问的下游依赖。

#### 6.3.1 HTTP 服务端 Metrics

HTTP 服务端至少应对齐以下标准指标：

| 指标名 | 类型 | 单位 | 地位 |
| :--- | :--- | :--- | :--- |
| `http.server.request.duration` | Histogram | `s` | 必选 |
| `http.server.active_requests` | UpDownCounter | `{request}` | 可选 |

推荐属性集合如下：

- `http.request.method`
- `http.response.status_code`
- `http.route`
- `server.address`
- `server.port`
- `url.scheme`
- `network.protocol.name`
- `network.protocol.version`
- `error.type`

统一约束如下：

1. `http.request.method` 应优先使用标准已知方法集合，未知值归并为 `_OTHER`。
2. `http.route` 必须保持低基数，不能用原始路径替代。
3. `server.address` 与 `server.port` 若来源于请求头，应注意高基数与攻击面风险。

#### 6.3.2 RPC / gRPC 服务端 Metrics

RPC / gRPC 服务端至少应对齐以下标准指标：

| 指标名 | 类型 | 单位 | 地位 |
| :--- | :--- | :--- | :--- |
| `rpc.server.call.duration` | Histogram | `s` | 必选 |

推荐属性集合如下：

- `rpc.system.name`
- `rpc.method`
- `rpc.response.status_code`
- `error.type`
- `server.address`
- `server.port`

统一约束如下：

1. gRPC 调用固定设置 `rpc.system.name = grpc`。
2. `rpc.method` 表达 RPC 接口视角的方法名，而不是具体实现函数名。
3. `rpc.response.status_code` 表达 RPC / gRPC 返回码。

### 6.4 数据库、消息与其它中间件场景

数据库、消息、缓存等场景，若 OpenTelemetry 已定义对应语义，应直接复用其标准指标名与标准属性集合。平台不再定义 `middleware_*` 一类兜底指标名作为长期主标准。平台侧只负责采集、聚合、导出与治理，而不改写标准指标语义。

### 6.5 默认资源维度与高基数边界

以下资源属性建议作为全平台默认低基数聚合维度：

- `service.name`
- `service.namespace`
- `service.version`
- `deployment.environment.name`
- `k8s.cluster.name`
- `cloud.region`
- `cloud.availability_zone`
- `k8s.namespace.name`
- `k8s.node.name`

以下字段不应默认进入指标属性：

- `user_id`
- `tenant_id`
- `request_id`
- `session_id`
- `trace_id`
- `span_id`
- `k8s.pod.name`
- `host.ip`
- 原始 URL 全路径

这类信息应优先进入日志字段、Trace 属性或离线分析链路。

### 6.6 错误码进入 Logs / Trace / Metrics 的统一规则

错误码规范与可观测规范必须使用同一套信号映射。统一原则是：协议状态保留原生语义，`error.type` 用于 OTel 通用错误聚合，`stellar.error.*` 只承担跨协议错误语义补差。
完整的错误分层与响应模型，参见 [错误码规范](/zh/topics/error-code-spec)。

| 错误语义 | 日志字段 | Trace 字段 | Metrics 字段 | 规则 |
| :--- | :--- | :--- | :--- | :--- |
| HTTP / gRPC 原生错误状态 | `http.response.status_code` 或 `rpc.response.status_code` | `http.response.status_code` 或 `rpc.response.status_code` | `http.response.status_code` 或 `rpc.response.status_code` | 必须保留原生协议状态 |
| OTel 通用错误类型 | `error.type` | `error.type` 与 Span Status = `Error` | `error.type` | 成功时不设置；失败时三种信号保持一致 |
| 基础规范码 | `stellar.error.code` | `stellar.error.code` | 默认不进入指标属性 | 对齐 canonical code 名称，如 `INVALID_ARGUMENT`、`UNAVAILABLE` |
| 业务域 | `stellar.error.domain` | `stellar.error.domain` | 默认不进入指标属性 | 用于检索和排障，不作默认聚合维度 |
| 业务原因码 | `stellar.error.reason` | `stellar.error.reason` | 默认不进入指标属性 | 必须稳定且低到中等基数 |
| 人类可读详情 | `Body` 或领域日志属性 | Span Status Description、异常事件或领域属性 | 不进入指标 | 禁止把动态详情文本写成指标标签 |
| 字段校验 / 前置条件细节 | 结构化日志属性 | Span 事件或结构化属性 | 不进入指标 | 仅用于诊断，不用于平台默认聚合 |

统一约束如下：

1. `error.type` 不等于业务错误码。它服务于 OpenTelemetry 错误聚合，优先使用协议错误状态或稳定低基数错误标识。
2. `stellar.error.code`、`stellar.error.domain`、`stellar.error.reason` 主要进入日志与 Trace，不默认进入 Metrics。
3. 发生异常时，Trace 应设置 Span Status 为 `Error`，并记录 `error.type`；若有异常对象，应按 OpenTelemetry 规范记录异常事件。
4. `error.message` 或任意自由文本错误描述，不应作为 Span 常驻聚合属性，也不应作为 Metrics 标签。
5. 需要按业务原因码做指标分析时，必须显式评审基数风险，并通过白名单单独开启。

## 7. 平台分层职责

### 7.1 发布平台

发布平台负责生成服务身份与部署环境的主配置。其核心职责包括：

1. 明确并下发 `service.name`、`service.namespace`、`service.version`、`deployment.environment.name`。
2. 设计 `service.instance.id` 的生成规则，并确保其不依赖 IP。
3. 优先输出 `OTEL_SERVICE_NAME` 与 `OTEL_RESOURCE_ATTRIBUTES`。
4. 仅在运行时无法直接设置标准 Resource 时，补充 `STELLAR_*`。

### 7.2 Helm / WebHook

Helm 模板或 Admission WebHook 负责把 Kubernetes 运行时元数据稳定暴露给工作负载，并映射到标准资源语义。其核心职责包括：

1. 提供 `metadata.name`、`metadata.uid`、`metadata.namespace`、`spec.nodeName`、`status.podIP`、`status.hostIP`。
2. 将这些事实直接映射到 `k8s.*` 与 `host.*`。
3. 若平台维护集群、区域与可用区信息，再补充 `k8s.cluster.name`、`cloud.region`、`cloud.availability_zone`。

### 7.3 网关

网关是标准传播链路的第一责任点。其核心职责包括：

1. 提取并注入 `traceparent`、`tracestate`、`baggage`。
2. 生成或规范化 `X-Stellar-Request-Id` 与 `X-Request-Id`。
3. 清洗外部不可信的 `X-Stellar-*`。
4. 对租户、用户、设备等头执行鉴权后再注入。
5. 只保留 `X-Stellar-Source-Service`，并确保其值对齐调用方 `service.name`。

### 7.4 SDK

SDK 的职责是把标准资源、标准传播头和标准指标暴露给业务代码，避免业务层重复实现。其核心职责包括：

1. 默认优先读取标准 OTel Resource。
2. 自动提取与注入 `traceparent`、`tracestate`、`baggage`。
3. 支持最小补差头，如 `X-Stellar-Request-Id`、`X-Stellar-Tenant-Id`、`X-Stellar-User-Id`、`X-Stellar-Source-Service`。
4. 分别暴露客户端视角和服务端视角的标准指标，如 `http.server.request.duration`、`http.client.request.duration`、`rpc.server.call.duration`、`rpc.client.call.duration`。
5. 保持日志、Trace、Metrics 共用同一份资源上下文模型。

### 7.5 Collector

Collector 负责统一接收、加工、导出与治理遥测数据，但不改写主语义。其核心职责包括：

1. 启用标准 OTLP 接收链路。
2. 优先保留上游写入的标准资源属性。
3. 只在来源明确且上游缺失时补充标准属性。
4. 不把 `k8s.pod.ip` 或 `host.ip` 改写成 `service.instance.id`。
5. 对高基数属性、敏感字段、请求头与 metadata 采集建立白名单和脱敏规则。
6. 保证客户端与服务端指标出口均保持 OTel 标准指标名。

## 8. 落地验收原则

一套可观测实施工作只有在以下条件同时满足时，才可视为完成：

1. 标准资源属性已成为主实现路径。
2. 标准传播头已成为主传播路径。
3. 标准 OTel 指标名已成为主暴露路径。
4. `STELLAR_*` 与 `X-Stellar-*` 仅保留最小补差集合。
5. 发布平台、模板、网关、SDK、Collector 的实现边界保持一致。
6. `service.name`、`service.namespace`、`deployment.environment.name` 在日志、Trace、Metrics 三个信号面向上保持一致。
7. 同一服务既能产出客户端日志 / 指标，也能产出服务端日志 / 指标，且两类视角不混淆。
8. 跨 HTTP 与 gRPC 的调用链路中，trace 与最小补差上下文不会丢失。

## 9. 结论

对于企业平台而言，最稳妥的可观测实现路线不是重新创造一套企业私有标准，而是以 OpenTelemetry 为主标准，以 Kubernetes 为事实来源，以最小补差层补齐平台治理需求。在这一前提下，客户端视角与服务端视角必须被严格拆开：客户端日志和指标只描述“我发出的调用”，服务端日志和指标只描述“我收到并处理的请求”。`STELLAR_*` 与 `X-Stellar-*` 可以存在，但它们只能处于补差和治理边界内，而不能成为资源语义、传播协议或指标语义的主定义。只有坚持这一边界，平台才能既保持统一治理能力，又持续复用开源生态的工具链与语义一致性。
