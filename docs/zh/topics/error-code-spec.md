---
title: 错误码规范
category: 体系规范
summary: 以 HTTP、gRPC 与 OpenTelemetry 为主标准，统一定义基础错误码、业务异常扩展方式、错误响应载体与观测映射规则。
tags:
  - 错误码
  - OpenTelemetry
  - HTTP
  - gRPC
  - API 设计
readingDirection: 适合在设计 API 错误模型、跨协议错误响应、业务异常扩展与可观测错误聚合时优先阅读。
outline: deep
---

# 错误码规范

## 摘要

本文提出一套面向企业平台的统一错误码规范。该规范以 HTTP 原生状态码语义、gRPC 原生状态码语义、`google.rpc.Status` 错误模型与 OpenTelemetry 错误记录规范为基础，目标是在不破坏既有协议标准的前提下，建立一套可跨 HTTP 与 gRPC 复用、可被业务扩展、可被客户端稳定处理、可被观测系统统一聚合的错误语义体系。本文将错误空间划分为传输层状态、基础规范码、业务扩展原因、结构化错误详情与观测映射五层，主张优先复用小而稳定的标准错误集合，将业务差异放在结构化扩展层，而不是重新发明协议级错误码。

## 关键词

HTTP；gRPC；Problem Details；google.rpc.Status；OpenTelemetry；Error Code；API Design

## 1. 引言

错误处理是网络 API 设计中最容易被低估、却最直接影响客户端体验与平台治理能力的部分。很多系统虽然定义了大量错误码，但这些错误码往往混杂了三种不同问题：一是把 HTTP 或 gRPC 的传输状态和业务错误混为一谈，二是为每个业务场景发明新的大规模私有整数码，三是错误响应只服务于人类阅读，却缺乏稳定的机器可处理结构。

标准生态已经给出了相对成熟的答案。HTTP 使用 RFC 9110 定义状态码语义；HTTP 错误响应体有 RFC 9457 Problem Details；gRPC 使用一组定义明确的状态码，并在实践中广泛采用 `google.rpc.Status` 作为结构化错误模型；OpenTelemetry 则要求错误在 Trace、Metrics、Logs 中以低基数、可聚合的方式记录。因此，更合理的路线不是重新创造一套与标准并行的错误体系，而是在保留原生协议语义的前提下，建立一层跨协议一致的基础错误码与业务扩展模型。

## 2. 标准基线

本规范优先对齐以下标准：

- HTTP 状态码语义遵循 RFC 9110
- HTTP 错误体格式遵循 RFC 9457 Problem Details
- gRPC 状态码遵循 gRPC 官方 `Status Codes`
- 跨协议结构化错误模型优先参考 `google.rpc.Status`
- 业务扩展错误详情优先参考 `google.rpc.ErrorInfo`、`BadRequest`、`PreconditionFailure`
- 观测侧错误记录优先对齐 OpenTelemetry 的 `error.type`、`http.response.status_code`、`rpc.response.status_code`

推荐阅读：

1. [RFC 9110: HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
2. [RFC 9457: Problem Details for HTTP APIs](https://datatracker.ietf.org/doc/html/rfc9457)
3. [gRPC Status Codes](https://grpc.io/docs/guides/status-codes/)
4. [gRPC Error Handling](https://grpc.io/docs/guides/error/)
5. [Cloud API Design Guide](https://docs.cloud.google.com/apis/design)
6. [Cloud API Design Guide - Errors](https://cloud.google.com/apis/design/errors?hl=ja)
7. [Package google.rpc](https://cloud.google.com/functions/docs/reference/rpc/google.rpc)
8. [OpenTelemetry HTTP Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/http/)
9. [OpenTelemetry RPC Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/rpc/rpc-spans/)
10. [OpenTelemetry Error Attribute Registry](https://opentelemetry.io/docs/specs/semconv/attributes-registry/error/)
11. [OpenTelemetry Recording Errors](https://opentelemetry.io/docs/specs/semconv/general/recording-errors/)

## 3. 错误空间分层模型

### 3.1 分层原则

错误语义应分为以下五层：

1. 传输层状态：HTTP status code 或 gRPC status code
2. 基础规范码：跨协议一致的基础错误类别
3. 业务扩展原因：具体业务失败原因
4. 结构化详情：字段错误、前置条件失败、重试建议、资源标识等
5. 观测映射：用于 Trace、Metrics、Logs 聚合和筛选的低基数错误标识

这五层中，只有第一层和第二层应保持高度稳定和低基数；第三层允许有限扩展；第四层承载上下文；第五层服务观测聚合。

### 3.2 传输层状态不是业务错误码

传输层状态用于表达协议语义，而不是替代业务异常模型：

- HTTP 侧必须首先选择符合 RFC 9110 语义的状态码
- gRPC 侧必须首先选择符合 gRPC 语义的状态码
- 不允许为了业务方便重新定义 HTTP 状态码
- 不允许为了业务方便重新定义 gRPC 状态码

因此，本规范不鼓励“200 + 业务错误对象”这种做法，也不鼓励“所有失败都返回 500 / UNKNOWN”这种粗粒度处理。

### 3.3 基础规范码采用小错误空间

业界最成熟的跨协议错误码实践，是使用一小组稳定、标准、可泛化的基础错误类别，再通过结构化详情承载业务差异。Google API Design Guide 明确建议使用一小组标准错误，而不是在每个 API 中额外定义大量新错误码。这一思想与 gRPC 的 canonical status code 体系一致。

因此，本规范建议基础规范码直接对齐 `google.rpc.Code` / gRPC canonical status code 的命名空间，并将其作为跨协议的基础错误类别。

### 3.4 业务异常不扩展协议码，而扩展业务原因

业务异常的扩展点不应放在 HTTP status code 或 gRPC status code 上，而应放在业务原因层：

- 基础规范码回答“这是哪一类错误”
- 业务原因回答“具体为什么发生了这类错误”
- 结构化详情回答“客户端要修什么、重试什么、展示什么”

这也是 RFC 9457 与 `google.rpc.Status` 共同体现出的设计思想：状态码用于粗分类，响应体用于结构化细节。

## 4. 基础错误码规范

### 4.1 基础规范码集合

平台建议采用以下基础规范码作为统一错误空间：

| 基础规范码 | gRPC 状态 | HTTP 推荐状态 | 含义 |
| :--- | :--- | :--- | :--- |
| `INVALID_ARGUMENT` | `INVALID_ARGUMENT` | `400` | 参数本身无效，和系统当前状态无关 |
| `UNAUTHENTICATED` | `UNAUTHENTICATED` | `401` | 未认证或认证凭据无效 |
| `PERMISSION_DENIED` | `PERMISSION_DENIED` | `403` | 已认证但无权限 |
| `NOT_FOUND` | `NOT_FOUND` | `404` | 资源不存在 |
| `ALREADY_EXISTS` | `ALREADY_EXISTS` | `409` | 创建目标已存在 |
| `FAILED_PRECONDITION` | `FAILED_PRECONDITION` | `409` / `412` / `428` | 当前资源状态不满足执行前提 |
| `OUT_OF_RANGE` | `OUT_OF_RANGE` | `400` / `416` / `422` | 范围越界或语义越界 |
| `RESOURCE_EXHAUSTED` | `RESOURCE_EXHAUSTED` | `429` | 配额、限流、资源耗尽 |
| `CANCELLED` | `CANCELLED` | 无统一公共状态 | 调用被取消，通常由调用方触发 |
| `DEADLINE_EXCEEDED` | `DEADLINE_EXCEEDED` | `504` | 超时 |
| `ABORTED` | `ABORTED` | `409` | 并发冲突、中止、重试事务 |
| `UNIMPLEMENTED` | `UNIMPLEMENTED` | `501` | 未实现 |
| `UNAVAILABLE` | `UNAVAILABLE` | `503` | 服务临时不可用 |
| `INTERNAL` | `INTERNAL` | `500` | 系统内部错误 |
| `DATA_LOSS` | `DATA_LOSS` | `500` | 数据损坏或不可恢复错误 |
| `UNKNOWN` | `UNKNOWN` | `500` | 无法归类的未知错误 |

### 4.2 选择原则

基础规范码的选择应遵循以下原则：

1. 优先返回最具体的错误，而不是最泛化的错误。
2. 如果 `NOT_FOUND` 与 `FAILED_PRECONDITION` 都能解释，优先使用更具体者。
3. 如果 `ALREADY_EXISTS` 与 `FAILED_PRECONDITION` 都能解释，优先使用 `ALREADY_EXISTS`。
4. 如果参数本身错误，使用 `INVALID_ARGUMENT`；如果参数在当前资源状态下不成立，使用 `FAILED_PRECONDITION`。
5. 如果失败可以通过等待配额恢复或系统恢复解决，优先考虑 `RESOURCE_EXHAUSTED` 或 `UNAVAILABLE`，而不是 `INTERNAL`。

### 4.3 关于 HTTP 映射的说明

虽然 Google 文档中给出了 `google.rpc.Code` 到 HTTP 状态码的映射，但本规范不要求机械地一一照抄该映射作为 HTTP 公共契约。原因如下：

1. HTTP API 本身应首先遵循 RFC 9110 的原生语义。
2. 某些 gRPC canonical code 在 HTTP 中并没有完全等价的一一映射。
3. 某些内部网关状态，例如 `499 Client Closed Request`，并不是 RFC 9110 的公共标准状态码。

因此，HTTP 侧的实践应当是：

1. 先选最符合 RFC 9110 的 HTTP 状态码。
2. 再在响应体中携带基础规范码。
3. 不通过新造 HTTP 状态码来表达业务异常。

## 5. HTTP 错误模型

### 5.1 响应行规范

HTTP API 在发生错误时：

1. 必须返回符合 RFC 9110 的错误状态码
2. 不得使用 `200 OK` 包裹失败语义
3. 应尽量使用最具体的 4xx / 5xx 状态码

### 5.2 响应体规范

HTTP 错误响应体建议遵循 RFC 9457 `application/problem+json`。标准字段包括：

- `type`
- `title`
- `status`
- `detail`
- `instance`

本规范建议增加以下扩展成员：

| 字段 | 说明 |
| :--- | :--- |
| `code` | 基础规范码，值对齐 `google.rpc.Code` / gRPC canonical code 名称 |
| `domain` | 业务域，例如 `trade.order` |
| `reason` | 业务原因码，稳定、机器可读 |
| `retryable` | 是否建议客户端重试 |
| `request_id` | 请求标识 |
| `violations` | 字段错误或前置条件错误明细 |

### 5.3 `type` 与 `reason` 的关系

二者职责不同：

- `type` 是 RFC 9457 的问题类型标识，必须是 URI
- `reason` 是业务稳定原因码，是机器处理的主扩展点

推荐写法如下：

- `type`: `https://errors.stellhub.top/trade/order-not-cancellable`
- `reason`: `ORDER_NOT_CANCELLABLE`

### 5.4 HTTP 示例

```json
{
  "type": "https://errors.stellhub.top/trade/order-not-cancellable",
  "title": "Order cannot be cancelled",
  "status": 409,
  "detail": "The order has already entered settlement and can no longer be cancelled.",
  "instance": "/orders/123/cancel",
  "code": "FAILED_PRECONDITION",
  "domain": "trade.order",
  "reason": "ORDER_NOT_CANCELLABLE",
  "retryable": false,
  "request_id": "req-20260426-000001"
}
```

## 6. gRPC 错误模型

### 6.1 状态码规范

gRPC API 在发生错误时：

1. 必须返回标准 gRPC status code
2. 不得定义新的 gRPC status code
3. 应尽量返回最具体的 canonical code

### 6.2 结构化载体规范

gRPC 错误响应建议使用 `google.rpc.Status`：

- `code`: 基础规范码，对应 `google.rpc.Code`
- `message`: 开发者可读的错误消息
- `details`: 结构化错误详情

### 6.3 标准详情类型

对于 gRPC，本规范推荐优先复用以下标准详情类型：

| 类型 | 用途 |
| :--- | :--- |
| `google.rpc.ErrorInfo` | 业务扩展原因、错误域、元数据 |
| `google.rpc.BadRequest` | 字段校验失败 |
| `google.rpc.PreconditionFailure` | 资源状态或业务前置条件不满足 |
| `google.rpc.LocalizedMessage` | 用户可见本地化消息 |
| `google.rpc.RetryInfo` | 重试建议 |

其中：

- `ErrorInfo.reason` 应作为业务异常扩展的主字段
- `ErrorInfo.domain` 应作为业务域命名空间
- `ErrorInfo.metadata` 承载少量结构化上下文，不应塞入任意大对象

### 6.4 gRPC 示例

```json
{
  "code": 9,
  "message": "Order cannot be cancelled",
  "details": [
    {
      "@type": "type.googleapis.com/google.rpc.ErrorInfo",
      "domain": "trade.order",
      "reason": "ORDER_NOT_CANCELLABLE",
      "metadata": {
        "orderId": "123"
      }
    },
    {
      "@type": "type.googleapis.com/google.rpc.PreconditionFailure",
      "violations": [
        {
          "type": "ORDER_STATE",
          "subject": "orders/123",
          "description": "Order has already entered settlement."
        }
      ]
    }
  ]
}
```

## 7. 业务异常扩展规范

### 7.1 扩展原则

业务异常扩展应遵循以下原则：

1. 不扩展 HTTP 状态码
2. 不扩展 gRPC 状态码
3. 不扩展基础规范码集合
4. 只扩展 `domain`、`reason` 和结构化 `details`

### 7.2 业务原因码的推荐形式

业务扩展推荐采用字符串型原因码，而不是私有整数码作为主契约。原因如下：

1. RFC 9457 本身没有定义整数业务码标准。
2. `google.rpc.ErrorInfo.reason` 天然适合承载稳定字符串原因码。
3. 字符串原因码比私有整数码更易读、更易跨语言、更易跨协议。
4. OpenTelemetry 对 `error.type` 的要求是低基数、可预测，而不是整数优先。

推荐格式如下：

- 使用 `UPPER_SNAKE_CASE`
- 不包含资源 ID、用户名、订单号等高基数值
- 在同一 `domain` 下唯一

示例：

- `ORDER_NOT_CANCELLABLE`
- `INSUFFICIENT_BALANCE`
- `RATE_PLAN_EXPIRED`
- `USER_KYC_REQUIRED`

### 7.3 `domain` 的推荐形式

`domain` 用于界定原因码的命名空间。推荐形式如下：

- `trade.order`
- `identity.auth`
- `billing.invoice`
- `platform.quota`

不建议直接使用：

- 随意缩写
- 团队昵称
- 版本号
- 动态租户名

### 7.4 何时新增一个业务原因码

仅当以下条件至少满足其一时，才应新增原因码：

1. 客户端需要据此采取不同处理动作
2. 运维或支持团队需要据此快速区分不同故障类别
3. 同一基础规范码下，已有原因码无法稳定表达该场景

如果只是文案不同、字段值不同、资源 ID 不同，而客户端动作不变，则不应新增原因码。

### 7.5 字段错误与业务状态错误的扩展方式

不同类型的业务错误应使用不同扩展载体：

1. 字段校验失败：优先使用 `BadRequest` 或 `violations`
2. 资源状态不满足：优先使用 `PreconditionFailure`
3. 通用业务失败原因：优先使用 `ErrorInfo.reason`
4. 用户可见消息：优先使用 `LocalizedMessage` 或 Problem Details 中的本地化扩展成员

不要为每个字段校验失败单独发明一个新的全局业务错误码。

## 8. OpenTelemetry 观测映射

### 8.1 总体原则

OpenTelemetry 的要求是：错误记录应可预测、低基数、可聚合。因此，观测系统不应直接把无限扩展的业务原因码默认放入 Metrics 维度。

### 8.2 HTTP 侧映射

HTTP 侧至少应记录：

- `http.response.status_code`
- `error.type`

根据 OpenTelemetry HTTP 语义约定：

- 如果响应状态码表明错误，`error.type` 可以是状态码字符串、异常类型或低基数组件错误标识

本规范建议：

1. `http.response.status_code` 必须保留原生 HTTP 状态
2. `error.type` 默认优先使用 HTTP 状态码字符串或稳定低基数基础规范码
3. 业务扩展原因 `reason` 不应默认作为 Metrics 标签

### 8.3 gRPC 侧映射

gRPC 侧至少应记录：

- `rpc.response.status_code`
- `error.type`

根据 OpenTelemetry RPC 语义约定：

- 当有返回状态码且其表示错误时，`error.type` 应设置为该状态码

本规范建议：

1. `rpc.response.status_code` 必须保留原生 gRPC 状态
2. `error.type` 默认优先使用 gRPC canonical code
3. `ErrorInfo.reason` 可进入日志与 Trace 属性，但不应默认进入 Metrics 维度

### 8.4 平台补差属性

由于 OpenTelemetry 没有定义统一的跨业务错误原因字段，本规范允许增加最小补差属性：

- `stellar.error.code`: 基础规范码
- `stellar.error.domain`: 业务域
- `stellar.error.reason`: 业务原因码

统一约束如下：

1. 这些属性主要用于日志与 Trace 检索
2. 不应默认作为全平台 Metrics 标签
3. `stellar.error.reason` 必须保持稳定和低到中等基数，禁止拼接动态参数

### 8.5 错误码进入 Logs / Trace / Metrics 的统一字段表

为避免错误码在响应体、日志、Span 和指标之间出现多套并行表达，平台统一按下表落地：
可观测视角下的客户端 / 服务端日志与指标边界，参见 [可观测规范](/zh/topics/observability-spec)。

| 错误语义 | HTTP 来源 | gRPC 来源 | Logs | Trace | Metrics | 规则 |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 传输层状态 | Problem Details `status` | `rpc.response.status_code` / `google.rpc.Status.code` 对应的 canonical code | `http.response.status_code` 或 `rpc.response.status_code` | `http.response.status_code` 或 `rpc.response.status_code` | `http.response.status_code` 或 `rpc.response.status_code` | 必须保留原生协议状态，不得改写为平台私有码 |
| 基础规范码 | Problem Details 扩展字段 `code` | `google.rpc.Status.code` 对应的 canonical code 名称 | `stellar.error.code` | `stellar.error.code` | 默认不进入指标属性 | 统一使用 canonical code 名称，如 `FAILED_PRECONDITION` |
| 业务域 | Problem Details 扩展字段 `domain` | `google.rpc.ErrorInfo.domain` | `stellar.error.domain` | `stellar.error.domain` | 默认不进入指标属性 | 用于检索和诊断，不用于默认聚合 |
| 业务原因码 | Problem Details 扩展字段 `reason` | `google.rpc.ErrorInfo.reason` | `stellar.error.reason` | `stellar.error.reason` | 默认不进入指标属性 | 必须稳定、可预测、低到中等基数 |
| OTel 通用错误类型 | 由 HTTP 错误状态推导 | 由 gRPC 错误状态推导 | `error.type` | `error.type` | `error.type` | 成功时不设置；失败时必须与对应 Span / Metric 保持一致 |
| 人类可读详情 | `title`、`detail` | `message`、`LocalizedMessage` | `Body` 或领域日志属性 | Span Status Description 或异常事件 | 不进入指标 | 不得把详情文本复制成指标标签 |
| 字段校验与前置条件详情 | `violations` | `BadRequest`、`PreconditionFailure` | 结构化日志属性 | Span 事件或结构化属性 | 不进入指标 | 只承载诊断细节，不作为默认聚合维度 |
| 重试建议 | `retryable` | `RetryInfo` | 可选 `stellar.error.retryable` | 可选 `stellar.error.retryable` | 默认不进入指标属性 | 仅在客户端需要自动重试决策时暴露 |

统一约束如下：

1. `error.type` 是 OpenTelemetry 的通用错误聚合字段，`stellar.error.code` 是跨 HTTP / gRPC 的基础规范码字段，二者职责不同，不能混用。
2. 当 HTTP / gRPC 已经返回标准错误状态时，`error.type` 应优先对齐协议错误状态；`stellar.error.code` 再补充跨协议统一类别。
3. `stellar.error.domain` 与 `stellar.error.reason` 主要进入日志与 Trace，用于排障、检索与根因分析。
4. Metrics 默认只保留协议状态和 `error.type`，不默认引入 `stellar.error.code`、`stellar.error.domain`、`stellar.error.reason`，以控制基数。
5. 若某个业务需要按原因码看错误率，应通过显式白名单或单独指标实现，而不是扩大平台默认指标维度。

## 9. 落地建议

### 9.1 HTTP API

HTTP API 应遵循以下流程：

1. 先选择最符合 RFC 9110 的 HTTP 状态码
2. 再返回 RFC 9457 Problem Details 响应体
3. 在扩展字段中加入 `code`、`domain`、`reason`
4. 需要字段级错误时加入 `violations`

### 9.2 gRPC API

gRPC API 应遵循以下流程：

1. 先选择最符合语义的 gRPC status code
2. 使用 `google.rpc.Status` 作为错误载体
3. 用 `ErrorInfo` 承载 `domain` 与 `reason`
4. 用 `BadRequest`、`PreconditionFailure`、`RetryInfo` 承载结构化详情

### 9.3 业务 SDK

业务 SDK 应：

1. 优先基于基础规范码决定客户端分支逻辑
2. 必要时再基于 `domain + reason` 做更细粒度处理
3. 不应依赖错误文案做程序分支

## 10. 结论

错误码设计最重要的不是“是否有一张完整的枚举表”，而是“是否把不同层次的语义放在了正确的位置”。传输层状态必须交给 HTTP 和 gRPC 原生标准；跨协议的一致性应依靠一小组稳定的基础规范码；业务差异应通过 `domain + reason + details` 扩展，而不是通过新增协议码或泛滥的私有整数码来表达；观测系统则应优先聚合基础规范码和协议状态，而不是高基数业务原因。只有在这一分层模型下，错误处理才能同时兼顾规范性、可扩展性、可观测性与长期可维护性。
