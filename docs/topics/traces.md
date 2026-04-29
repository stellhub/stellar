---
title: 大型企业跨语言微服务链路追踪技术调研方案
category: 可观测性
summary: 对 OpenTelemetry、Tempo、Jaeger、SkyWalking 与 Zipkin 进行平台化对比，给出大型企业链路追踪底座的推荐架构与落地路径。
tags:
  - 链路追踪
  - OpenTelemetry
  - Tempo
  - Jaeger
  - SkyWalking
  - Zipkin
readingDirection: 适合在选型企业级 Trace 平台、设计统一采集链路或规划可观测底座升级时优先阅读。
outline: deep
---

# 大型企业跨语言微服务链路追踪技术调研方案

## 1. 结论

面向大型企业、跨语言、跨团队、跨环境的微服务链路追踪平台，推荐采用：

```text
OpenTelemetry Collector + Grafana Tempo Distributed + Grafana
```

该方案的核心判断是：

> **OpenTelemetry 负责标准化采集，Grafana Tempo Distributed 负责低成本、高扩展 Trace 存储与查询，Grafana 负责统一可观测入口。**

在五种方案中，推荐优先级如下：

| 排名 | 方案                                          | 推荐程度 | 结论                                                   |
| -: | ------------------------------------------- | ---: | ---------------------------------------------------- |
|  1 | **OpenTelemetry + Grafana Tempo + Grafana** |   最高 | 最适合大型企业平台化、跨语言、低成本、高扩展链路追踪                           |
|  2 | OpenTelemetry + Jaeger                      |   较高 | 适合成熟 Trace 平台、已有 Jaeger 体系、PoC 或中大型场景                |
|  3 | SkyWalking                                  |   中高 | 适合 Java APM、服务拓扑、开箱即用可观测平台，但标准化和底座灵活性弱于 OTel + Tempo |
|  4 | Jaeger 原生 SDK + Jaeger                      |    低 | 不建议新平台采用，Jaeger 官方已建议迁移到 OpenTelemetry               |
|  5 | Zipkin                                      |    低 | 适合轻量 tracing、历史兼容，不建议作为大型企业统一链路追踪底座                  |

最终建议：

```text
采集标准：OpenTelemetry
采集网关：OpenTelemetry Collector
Trace 后端：Grafana Tempo Distributed
展示入口：Grafana
存储底座：对象存储，例如 S3 / OSS / COS / MinIO / Ceph
兼容协议：OTLP / Jaeger / Zipkin
```

---

## 2. 建设目标

大型企业级链路追踪平台不应只满足“能看到调用链”，而应满足以下目标：

| 目标     | 说明                                      |
| ------ | --------------------------------------- |
| 跨语言接入  | Java、Go、Node.js、Python、C++、Rust、前端等统一接入 |
| 低侵入接入  | Java Agent、SDK、Sidecar、Gateway 多种方式并存   |
| 采集标准统一 | 避免被单一厂商、单一后端、单一 SDK 绑定                  |
| 高吞吐写入  | 支撑亿级、十亿级甚至更高 Span 写入                    |
| 低成本存储  | Trace 数据体量巨大，必须控制长期存储成本                 |
| 多租户隔离  | 按业务线、团队、环境、地域、应用进行隔离和治理                 |
| 采样治理   | 支持错误请求、慢请求、核心链路、VIP 租户差异化采样             |
| 可观测闭环  | Metrics、Logs、Traces 能够互相跳转和关联           |
| 平台治理能力 | 支持接入规范、标签规范、脱敏、限流、路由、权限、审计              |
| 云原生部署  | 适配 Kubernetes、多集群、多地域、多可用区部署            |

因此，企业级链路追踪平台的重点不只是 Trace UI，而是：

```text
标准化采集 + 可治理管道 + 高扩展存储 + 统一分析入口
```

---

## 3. 官方依据

### 3.1 OpenTelemetry 的官方定位

OpenTelemetry 官方文档将其定义为一个 vendor-neutral 的开源可观测性框架，用于生成、采集和导出 traces、metrics、logs 等遥测数据；官方还强调它是行业标准，并被大量厂商、库、服务和终端用户采用。([OpenTelemetry][1])

OpenTelemetry 官网还明确强调“vendor-neutral instrumentation”，即一次埋点后可以将遥测数据导出到 Jaeger、Prometheus、商业厂商或自建后端，并且可以在不修改业务代码的情况下切换后端。([OpenTelemetry][2])

这对大型企业非常关键。因为企业内部通常存在多语言、多框架、多团队、多供应商、多云环境。如果链路追踪直接绑定 Jaeger SDK、Zipkin SDK 或某个 APM Agent，后续迁移、治理和统一标准都会变得困难。

---

### 3.2 OpenTelemetry Collector 的官方定位

OpenTelemetry Collector 官方文档说明，它提供 vendor-agnostic 的接收、处理和导出遥测数据能力，可以减少运行和维护多个 agent/collector 的需要，并支持将开源可观测数据格式发送到一个或多个开源或商业后端。([OpenTelemetry][3])

这意味着 Collector 可以作为大型企业的统一遥测网关：

```text
业务服务
  -> OpenTelemetry Collector
  -> Tempo / Jaeger / Prometheus / Loki / 商业平台
```

它不是单纯的转发器，而是可观测数据治理层，适合实现：

```text
采样
脱敏
标签补全
租户识别
协议转换
流量限速
后端路由
多集群汇聚
```

---

### 3.3 Grafana Tempo 的官方定位

Grafana 官方文档将 Tempo 定义为 open-source、easy-to-use、high-scale distributed tracing backend，并说明它可以搜索 traces、从 spans 生成 metrics，并把 tracing 数据与 logs 和 metrics 关联起来。([Grafana Labs][4])

Grafana Tempo 官方 OSS 页面还强调，Tempo 是高扩展的分布式 tracing 后端，成本较低，只需要对象存储即可运行，并且深度集成 Grafana、Prometheus、Loki，同时可以接收 Jaeger、Zipkin、OpenTelemetry 等常见 tracing 协议。([Grafana Labs][5])

这正好匹配大型企业平台化链路追踪的几个关键诉求：

```text
高扩展
低成本
对象存储
协议兼容
Grafana 生态集成
Metrics / Logs / Traces 联动
```

---

### 3.4 Tempo 多租户能力

Tempo 官方文档说明，Tempo 是 multi-tenant distributed tracing backend，并通过 `X-Scope-OrgID` header 实现多租户隔离；该 header 用于限定写入和读取，使查询只返回对应租户的数据。([Grafana Labs][6])

大型企业内部通常存在多个事业部、团队、环境和地域，Trace 平台必须具备租户隔离能力。Tempo 的多租户模型适合在平台层设计：

```text
tenant.id
team
env
region
cluster
namespace
service.name
app.id
```

---

### 3.5 Tail Sampling 的官方依据

OpenTelemetry 官方文档说明，Tail Sampling 是在看到一个 Trace 的全部或大部分 Span 之后再决定是否采样，因此可以根据 Trace 中不同部分的信息做更精确的采样决策。([OpenTelemetry][7])

OpenTelemetry Collector 的 tail sampling processor 文档进一步说明，在使用 tail sampling 时，同一个 Trace 的所有 Span 必须进入同一个 Collector 实例，才能做出正确采样决策；扩展 Collector 时通常需要两层 Collector，一层做 load balancing，一层做 tail sampling。([GitHub][8])

这对大型企业非常重要。因为全量采集 Trace 往往会带来极高的网络、存储和查询成本，必须通过平台级采样策略控制成本。

---

## 4. 推荐架构

### 4.1 总体架构

```text
┌───────────────────────────────────────────────┐
│ Business Services                              │
│ Java / Go / Node.js / Python / C++ / Rust      │
└───────────────────────────────────────────────┘
                       │
                       │ OTLP / HTTP / gRPC
                       ▼
┌───────────────────────────────────────────────┐
│ Node-level OpenTelemetry Collector             │
│ - batch                                        │
│ - memory_limiter                               │
│ - resource processor                           │
│ - attributes processor                         │
└───────────────────────────────────────────────┘
                       │
                       ▼
┌───────────────────────────────────────────────┐
│ Gateway OpenTelemetry Collector Cluster        │
│ - tenant identification                        │
│ - tail sampling                                │
│ - PII filtering                                │
│ - routing                                      │
│ - rate limiting                                │
│ - load shedding                                │
└───────────────────────────────────────────────┘
                       │
                       ▼
┌───────────────────────────────────────────────┐
│ Grafana Tempo Distributed                      │
│ - distributor                                  │
│ - ingester                                     │
│ - querier                                      │
│ - query-frontend                               │
│ - compactor                                    │
│ - metrics-generator                            │
└───────────────────────────────────────────────┘
                       │
                       ▼
┌───────────────────────────────────────────────┐
│ Object Storage                                 │
│ S3 / OSS / COS / MinIO / Ceph                  │
└───────────────────────────────────────────────┘
                       │
                       ▼
┌───────────────────────────────────────────────┐
│ Grafana                                        │
│ Metrics + Logs + Traces                        │
└───────────────────────────────────────────────┘
```

---

### 4.2 推荐数据流

```text
应用服务产生 Span
  -> OpenTelemetry SDK / Java Agent 自动采集
  -> OTLP 发送到本地 Collector
  -> Collector 补充资源标签、批量发送
  -> Gateway Collector 做租户识别、脱敏、采样、路由
  -> Tempo Distributed 写入对象存储
  -> Grafana 查询、展示、联动 Metrics 和 Logs
```

---

### 4.3 为什么不建议应用直连 Tempo

大型企业场景不建议：

```text
应用服务 -> Tempo
```

原因是直连会失去平台治理能力：

| 问题      | 影响                           |
| ------- | ---------------------------- |
| 难统一采样   | 每个服务各自配置，策略不可控               |
| 难统一脱敏   | 敏感字段可能直接进入存储                 |
| 难做租户识别  | 缺少统一租户注入和鉴权                  |
| 难做流控    | 后端异常时应用侧容易受影响                |
| 难做多后端路由 | 不能灵活切换 Tempo / Jaeger / 商业平台 |
| 难做接入审计  | 无法集中统计服务接入质量                 |

正确做法是：

```text
应用服务 -> OpenTelemetry Collector -> Tempo
```

Collector 是企业级 Trace 平台的控制面入口。

---

## 5. 五种方案详细对比

## 5.1 方案一：OpenTelemetry + Grafana Tempo + Grafana

### 方案描述

```text
OpenTelemetry SDK / Agent
  -> OpenTelemetry Collector
  -> Grafana Tempo Distributed
  -> Grafana
```

该方案使用 OpenTelemetry 作为统一采集标准，使用 Tempo Distributed 作为分布式 Trace 后端，使用 Grafana 作为统一可视化入口。

### 官方依据

OpenTelemetry 官方强调 vendor-neutral，可生成、采集、导出 traces、metrics、logs，并被大量厂商和用户采用。([OpenTelemetry][1])

Tempo 官方强调其为 high-scale distributed tracing backend，可与 logs 和 metrics 关联；Tempo OSS 页面还说明其成本较低，只需要对象存储运行，并深度集成 Grafana、Prometheus、Loki，且支持 Jaeger、Zipkin、OpenTelemetry 等协议。([Grafana Labs][4])

### 优势

| 维度   | 优势                                       |
| ---- | ---------------------------------------- |
| 标准化  | OpenTelemetry 是当前云原生可观测事实标准              |
| 跨语言  | Java、Go、Node.js、Python、C++、Rust 等均可接入    |
| 后端解耦 | 业务代码不绑定 Tempo，后续可切换后端                    |
| 存储成本 | Tempo 以对象存储为主，适合海量 Trace 数据              |
| 扩展能力 | Tempo Distributed 支持分布式扩展                |
| 多租户  | Tempo 支持基于 `X-Scope-OrgID` 的多租户隔离        |
| 生态集成 | Grafana + Prometheus + Loki + Tempo 联动自然 |
| 协议兼容 | 可接收 OTLP、Jaeger、Zipkin 等协议               |
| 平台治理 | Collector 适合做采样、脱敏、路由、限流、租户识别            |

### 劣势

| 问题                | 说明                                    |
| ----------------- | ------------------------------------- |
| 建设复杂度高于单体 APM     | 需要设计 Collector、Tempo、对象存储、Grafana 等组件 |
| 采样策略需要平台化治理       | Tail Sampling 配置和扩展需要工程设计             |
| 查询体验依赖 Grafana 体系 | 对非 Grafana 体系团队有学习成本                  |
| 需要规范标签模型          | 没有统一标签规范时，Trace 查询价值会下降               |

### 适用场景

| 场景                          | 适配度 |
| --------------------------- | --: |
| 大型企业统一链路追踪平台                |  很高 |
| 跨语言微服务系统                    |  很高 |
| 云原生 / Kubernetes            |  很高 |
| 多团队、多租户治理                   |  很高 |
| 长期低成本存储                     |  很高 |
| Metrics / Logs / Traces 一体化 |  很高 |
| 单体 Java APM 快速落地            |  中等 |

### 结论

这是最适合作为大型企业统一链路追踪底座的方案。

---

## 5.2 方案二：OpenTelemetry + Jaeger

### 方案描述

```text
OpenTelemetry SDK / Agent
  -> OpenTelemetry Collector
  -> Jaeger
```

该方案使用 OpenTelemetry 进行采集，Jaeger 作为 Trace 存储、查询和 UI。

### 官方依据

Jaeger 官方文档说明其特性包括 OpenTelemetry compatible，并支持 Elasticsearch、OpenSearch、Cassandra、Badger、Kafka、Memory 等多种存储后端，还支持系统拓扑、服务依赖图和 adaptive sampling。([Jaeger][9])

OpenTelemetry 官方迁移文档也说明，Jaeger backend 自 v1.35 起可以通过 OTLP 接收 trace 数据。([OpenTelemetry][10])

### 优势

| 维度               | 优势                                  |
| ---------------- | ----------------------------------- |
| 成熟度              | Jaeger 历史悠久，社区成熟                    |
| OpenTelemetry 兼容 | 可以使用 OTel SDK / Collector 接入        |
| UI 简洁            | Trace 查询和链路展示直观                     |
| 存储选择多            | 支持 Elasticsearch、Cassandra、Badger 等 |
| 学习成本低            | 本地 all-in-one 部署简单                  |
| 适合兼容             | 适合已有 Jaeger 历史系统                    |

### 劣势

| 问题                     | 说明                                    |
| ---------------------- | ------------------------------------- |
| 存储运维压力较大               | 大规模常依赖 Elasticsearch / Cassandra      |
| 长期成本不如 Tempo           | ES 索引型存储面对海量 Trace 成本较高               |
| Grafana 生态联动弱于 Tempo   | Metrics / Logs / Traces 一体化体验不如 Tempo |
| 多租户平台能力需要额外建设          | Jaeger 本身不是最强多租户底座                    |
| 更适合 Trace 系统，不是统一可观测底座 | 企业平台化需要额外补很多治理能力                      |

### 适用场景

| 场景                   | 适配度 |
| -------------------- | --: |
| 已有 Jaeger 体系升级到 OTel |  很高 |
| 中大型微服务链路追踪           |   高 |
| 本地调试 / PoC           |  很高 |
| 超大型企业低成本长期 Trace 存储  |   中 |
| Grafana 全栈可观测平台      |   中 |

### 结论

`OpenTelemetry + Jaeger` 是可接受的成熟方案，但不是大型企业新建统一 Trace 平台的最佳底座。若企业已有 Jaeger，可以保留 Jaeger 作为兼容后端，同时逐步把采集侧迁移到 OpenTelemetry。

---

## 5.3 方案三：SkyWalking

### 方案描述

```text
SkyWalking Agent
  -> SkyWalking OAP
  -> SkyWalking Storage
  -> SkyWalking UI
```

SkyWalking 是 All-in-one APM 方案，包含 tracing、metrics、logging、拓扑、告警和 UI。

### 官方依据

SkyWalking 官方文档说明，它是一个开源可观测平台，用于从服务和云原生基础设施中收集、分析、聚合和可视化数据，并提供 distributed tracing、service mesh telemetry analysis、metrics aggregation、alerting、visualization 等能力。([GitHub][11])

SkyWalking 官网还强调其是 All-in-one APM solution，支持 end-to-end distributed tracing、service topology analysis、service-centric observability、API dashboards，并提供 Java、.NET Core、PHP、Node.js、Golang、LUA、Rust、C++、Client JavaScript、Python 等 agent。([Apache SkyWalking][12])

### 优势

| 维度          | 优势                                     |
| ----------- | -------------------------------------- |
| APM 能力完整    | Trace、Metrics、Logs、拓扑、告警、Dashboard 较完整 |
| Java 生态强    | 对 Java 微服务非常友好                         |
| 开箱即用        | 比 OTel + Tempo 自建组合更像完整产品              |
| 服务拓扑好       | 服务依赖关系展示能力强                            |
| Agent 覆盖多语言 | 官方支持多种语言 Agent                         |
| 适合业务团队直接使用  | UI 和 APM 能力更产品化                        |

### 劣势

| 问题                     | 说明                                                      |
| ---------------------- | ------------------------------------------------------- |
| 标准化程度弱于 OpenTelemetry  | 虽然支持多协议，但核心体系仍偏 SkyWalking 自身生态                         |
| 后端解耦弱                  | 业务接入、分析模型、平台能力更绑定 SkyWalking                            |
| 和 Grafana 体系整合不如 Tempo | 如果企业已选择 Grafana + Prometheus + Loki，SkyWalking 会形成另一套入口 |
| 大规模 Trace 存储成本需要单独评估   | 不像 Tempo 那样天然强调对象存储低成本模型                                |
| 平台可编排能力弱于 Collector 体系 | 采样、脱敏、路由、多后端分发不如 OTel Collector 灵活                      |

### 适用场景

| 场景                    | 适配度 |
| --------------------- | --: |
| Java 微服务 APM          |  很高 |
| 希望快速获得服务拓扑和 Dashboard |  很高 |
| 中大型内部 APM 平台          |   高 |
| 跨语言统一标准采集平台           |   中 |
| 超大型企业低成本 Trace 存储底座   |   中 |
| Grafana 统一可观测入口       |   中 |

### 结论

SkyWalking 是优秀的 APM 产品型方案，尤其适合 Java 技术栈占主导的企业。但如果目标是建设跨语言、开放标准、可插拔、低成本、高扩展的统一链路追踪底座，`OpenTelemetry + Tempo` 更适合作为主路线。

---

## 5.4 方案四：Jaeger 原生 SDK + Jaeger

### 方案描述

```text
Jaeger Client SDK
  -> Jaeger Agent / Collector
  -> Jaeger Storage
  -> Jaeger UI
```

这是早期 Jaeger 的典型接入方式，业务代码直接使用 Jaeger Client 或 OpenTracing API。

### 官方依据

OpenTelemetry 官方迁移文档明确说明，Jaeger 社区已经 deprecated Jaeger client libraries，并推荐迁移到 OpenTelemetry APIs、SDKs 和 instrumentations；同一文档还说明 Jaeger backend 自 v1.35 起可以接收 OTLP Trace 数据。([OpenTelemetry][10])

Jaeger Go Client 仓库也明确标注该库已 deprecated，不再接受除安全修复外的新 PR，并敦促用户迁移到 OpenTelemetry。([GitHub][13])

Jaeger Java Client 仓库同样标注 deprecated，说明 v1.8.1 是最终版本，不再接受新 PR，并敦促用户迁移到 OpenTelemetry。([GitHub][14])

### 优势

| 维度            | 优势                                     |
| ------------- | -------------------------------------- |
| 历史成熟          | 曾经是 Jaeger 的主流接入方式                     |
| 与 Jaeger 模型一致 | 直接面向 Jaeger 使用                         |
| 老系统兼容         | 已有 OpenTracing / Jaeger SDK 项目可以继续短期维护 |

### 劣势

| 问题        | 说明                         |
| --------- | -------------------------- |
| 官方已不推荐    | Jaeger Client 已 deprecated |
| 迁移风险高     | 新平台继续使用会扩大历史包袱             |
| 绑定 Jaeger | 业务代码和 Jaeger SDK 绑定，后端切换困难 |
| 跨语言标准化弱   | 不适合大型企业统一接入规范              |
| 生态趋势落后    | OpenTelemetry 已成为主流标准      |
| 维护风险      | 安全修复之外基本不再演进               |

### 适用场景

| 场景            | 适配度 |
| ------------- | --: |
| 老系统维持运行       |   中 |
| 新建企业 Trace 平台 |  很低 |
| 跨语言统一接入       |   低 |
| 长期技术路线        |  很低 |

### 结论

不建议新建平台采用该方案。已有 Jaeger 原生 SDK 的系统，应制定迁移计划，逐步迁移到 OpenTelemetry SDK / Agent / Collector。

---

## 5.5 方案五：Zipkin

### 方案描述

```text
Zipkin Instrumentation
  -> Zipkin Collector
  -> Zipkin Storage
  -> Zipkin Query
  -> Zipkin UI
```

Zipkin 是较早的分布式追踪系统，主要用于收集和查询服务架构中的延迟数据。

### 官方依据

Zipkin 官方首页说明，Zipkin 是 distributed tracing system，用于收集服务架构中排查延迟问题所需的 timing data，并提供数据收集和查询能力。([zipkin.io][15])

Zipkin 官方架构文档说明，Zipkin 最初使用 Cassandra 存储数据，后来将存储做成可插拔，并原生支持 Elasticsearch 和 MySQL。([zipkin.io][16])

### 优势

| 维度     | 优势                                 |
| ------ | ---------------------------------- |
| 简单     | 架构和概念清晰                            |
| 历史悠久   | 很多框架曾内置 Zipkin 支持                  |
| 轻量     | 适合小规模 tracing 或本地调试                |
| 存储可插拔  | 支持 Cassandra、Elasticsearch、MySQL 等 |
| 协议兼容价值 | 很多平台仍支持 Zipkin 协议接入                |

### 劣势

| 问题                     | 说明                                              |
| ---------------------- | ----------------------------------------------- |
| 平台化能力弱                 | 不适合作为大型企业统一治理底座                                 |
| 生态趋势较弱                 | 新建平台通常不会优先选择 Zipkin                             |
| 查询和分析能力有限              | 相比 Tempo TraceQL、SkyWalking APM、Jaeger UI，能力偏基础 |
| 大规模治理能力不足              | 多租户、采样、脱敏、路由等需要大量外围建设                           |
| 与 Grafana 全栈联动不如 Tempo | Metrics / Logs / Traces 一体化体验较弱                 |

### 适用场景

| 场景                | 适配度 |
| ----------------- | --: |
| 小型系统 tracing      |   中 |
| 历史系统兼容            |   中 |
| 本地调试              |   中 |
| 新建大型企业统一 Trace 平台 |   低 |
| 多语言标准化采集平台        |   低 |

### 结论

Zipkin 可以作为兼容协议保留，但不建议作为大型企业新建链路追踪平台的核心方案。

---

## 6. 横向评分对比

评分范围：1 到 5，5 为最好。

| 维度                     | OTel + Tempo + Grafana | OTel + Jaeger | SkyWalking | Jaeger SDK + Jaeger | Zipkin |
| ---------------------- | ---------------------: | ------------: | ---------: | ------------------: | -----: |
| 跨语言标准化                 |                      5 |             5 |          4 |                   2 |      3 |
| 后端解耦                   |                      5 |             4 |          3 |                   1 |      2 |
| 大规模写入能力                |                      5 |             4 |          4 |                   3 |      3 |
| 长期存储成本                 |                      5 |             3 |          3 |                   3 |      3 |
| 多租户能力                  |                      5 |             3 |          3 |                   2 |      2 |
| Grafana 生态集成           |                      5 |             3 |          2 |                   2 |      2 |
| Metrics/Logs/Traces 联动 |                      5 |             3 |          4 |                   2 |      2 |
| Java APM 开箱即用          |                      3 |             3 |          5 |                   3 |      2 |
| 平台治理能力                 |                      5 |             4 |          3 |                   2 |      2 |
| 迁移友好度                  |                      5 |             4 |          3 |                   1 |      2 |
| 长期技术趋势                 |                      5 |             4 |          4 |                   1 |      2 |

综合判断：

```text
OpenTelemetry + Grafana Tempo + Grafana
  > OpenTelemetry + Jaeger
  > SkyWalking
  > Zipkin
  > Jaeger 原生 SDK + Jaeger
```

其中 `Jaeger 原生 SDK + Jaeger` 的长期评分低于 Zipkin，主要原因是 Jaeger 官方已经明确推荐迁移到 OpenTelemetry，而不是继续使用 Jaeger Client。

---

## 7. 为什么推荐 OpenTelemetry + Tempo + Grafana

## 7.1 原因一：采集标准必须与后端解耦

大型企业不应让业务代码直接依赖某个 Trace 后端 SDK。正确方式是：

```text
业务代码接入 OpenTelemetry
后端可以是 Tempo / Jaeger / 商业平台 / 自研平台
```

OpenTelemetry 官方的 vendor-neutral 定位决定了它更适合作为企业统一接入标准。([OpenTelemetry][2])

这可以避免以下问题：

```text
Jaeger SDK 绑定
Zipkin SDK 绑定
APM Agent 绑定
后端迁移困难
跨语言规范不一致
```

---

## 7.2 原因二：Collector 是企业级治理层

大型企业 Trace 平台必须有集中治理能力：

```text
采样策略
字段脱敏
租户识别
标签标准化
数据路由
限流降级
多后端分发
接入质量监控
```

OpenTelemetry Collector 官方定位就是 vendor-agnostic 的接收、处理和导出遥测数据组件，并用于减少多个 agent/collector 的维护成本。([OpenTelemetry][3])

因此 Collector 不只是“中转站”，而是链路追踪平台的核心控制点。

---

## 7.3 原因三：Tempo 更适合海量 Trace 的低成本存储

Trace 数据天然体量巨大。一个 HTTP 请求可能产生多个 Span，一个核心链路可能产生几十个甚至上百个 Span。大型企业如果全量采集，很容易达到每天数十亿甚至更高 Span 规模。

Tempo 官方强调其高扩展、低成本、只需要对象存储运行，并深度集成 Grafana、Prometheus 和 Loki。([Grafana Labs][5])

这意味着 Tempo 在大规模场景下有明显优势：

```text
对象存储成本低
容量扩展简单
适合长期留存
减少 Elasticsearch / Cassandra 运维压力
```

---

## 7.4 原因四：Grafana 统一入口更适合可观测闭环

大型企业排障流程通常不是单看 Trace，而是：

```text
Metrics 告警
  -> Dashboard 定位异常服务
  -> Trace 查看慢调用链
  -> Logs 查看具体错误
  -> Profiling 定位代码热点
```

Tempo 官方文档说明它可以将 tracing 数据与 logs 和 metrics 关联起来。([Grafana Labs][4])

如果企业已经使用 Prometheus、Grafana、Loki，那么 Tempo 是天然匹配的 Trace 后端。

---

## 7.5 原因五：Tail Sampling 是大规模平台必需能力

在大型企业里，不建议简单全量采集所有 Trace。更合理的是：

```text
错误请求：100% 保留
慢请求：100% 保留
核心链路：高比例采样
普通成功请求：低比例采样
健康检查：丢弃
压测流量：低比例或独立隔离
VIP 租户：单独策略
灰度发布窗口：临时提高采样
```

OpenTelemetry 官方说明 Tail Sampling 可以根据 Trace 的多个部分做采样决策，比 Head Sampling 更精确。([OpenTelemetry][7])

因此，大型企业应将 Tail Sampling 作为平台能力，而不是让每个业务服务自己决定采样。

---

## 8. 推荐平台能力设计

### 8.1 接入层能力

| 能力                  | 说明                                         |
| ------------------- | ------------------------------------------ |
| Java Agent 接入       | Spring Boot、Dubbo、gRPC、HTTP Client 自动埋点    |
| Go SDK 接入           | Go 微服务标准化封装                                |
| Node.js / Python 接入 | 提供统一接入模板                                   |
| Gateway 接入          | API Gateway、Ingress、Service Mesh 注入 Trace  |
| 前端接入                | Web RUM 可选接入                               |
| MQ 链路接入             | Kafka、RocketMQ、Pulsar、RabbitMQ Trace 上下文传播 |

---

### 8.2 Collector 层能力

| 能力                   | 说明                     |
| -------------------- | ---------------------- |
| OTLP Receiver        | 统一接收 OTLP gRPC / HTTP  |
| Jaeger Receiver      | 兼容历史 Jaeger 协议         |
| Zipkin Receiver      | 兼容历史 Zipkin 协议         |
| Attributes Processor | 补充、修改、删除标签             |
| Resource Processor   | 统一资源属性                 |
| Batch Processor      | 批量发送，提高吞吐              |
| Memory Limiter       | 防止 Collector OOM       |
| Tail Sampling        | 按错误、慢请求、租户、服务采样        |
| Routing Processor    | 按租户、地域、环境路由            |
| Filter Processor     | 丢弃健康检查、低价值 Trace       |
| Transform Processor  | 标准化字段和语义               |
| Exporter             | 输出到 Tempo、Jaeger 或其他后端 |

---

### 8.3 标签规范

建议统一以下资源标签：

```text
tenant.id
org.id
team
env
region
az
cluster
namespace
service.name
service.version
service.instance.id
app.id
runtime.language
deployment.environment
```

建议统一以下 Span 标签：

```text
http.method
http.route
http.status_code
rpc.system
rpc.service
rpc.method
db.system
db.name
messaging.system
messaging.destination
error.type
error.message
```

没有标签规范，Trace 平台会快速失控。标签是后续采样、检索、聚合、权限、成本分摊的基础。

---

### 8.4 多租户设计

Tempo 官方支持通过 `X-Scope-OrgID` 实现多租户隔离。([Grafana Labs][6])

建议平台层设计如下：

```text
tenant.id -> X-Scope-OrgID
team      -> 成本分摊
env       -> prod / staging / test / dev
region    -> cn-shanghai / ap-singapore / us-west
cluster   -> k8s cluster
namespace -> Kubernetes namespace
```

Collector Gateway 负责将业务侧标签转换为 Tempo 多租户 header。

---

### 8.5 安全与合规

大型企业链路追踪必须前置治理敏感信息。

建议默认禁止进入 Trace 的字段：

```text
password
token
authorization
cookie
phone
email
id_card
bank_card
address
access_key
secret_key
```

Collector 层应执行：

```text
敏感字段删除
敏感字段 hash
敏感字段 mask
标签白名单
标签长度限制
高基数字段限制
```

尤其需要限制以下高基数字段：

```text
user.id
order.id
request.id
session.id
sql.full
http.url.full
```

高基数字段会显著放大存储、索引和查询成本。

---

## 9. 推荐采样策略

### 9.1 基础策略

```text
错误 Trace：100%
慢 Trace：100%
核心服务 Trace：5% - 20%
普通成功 Trace：0.1% - 1%
健康检查：0%
压测流量：单独租户或低比例
灰度发布：发布窗口临时提高
VIP 租户：独立策略
```

### 9.2 Collector 示例配置

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 2048

  batch:
    timeout: 5s
    send_batch_size: 8192

  attributes/sanitize:
    actions:
      - key: authorization
        action: delete
      - key: cookie
        action: delete
      - key: password
        action: delete
      - key: token
        action: delete

  tail_sampling:
    decision_wait: 10s
    num_traces: 50000
    policies:
      - name: error-traces
        type: status_code
        status_code:
          status_codes:
            - ERROR

      - name: slow-traces
        type: latency
        latency:
          threshold_ms: 1000

      - name: important-services
        type: string_attribute
        string_attribute:
          key: service.name
          values:
            - order-service
            - payment-service
            - gateway-service

      - name: default-probabilistic
        type: probabilistic
        probabilistic:
          sampling_percentage: 1

exporters:
  otlp/tempo:
    endpoint: tempo-distributor:4317
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers:
        - otlp
      processors:
        - memory_limiter
        - attributes/sanitize
        - tail_sampling
        - batch
      exporters:
        - otlp/tempo
```

代码注释如需添加，建议保持英文注释，便于跨团队维护。

---

## 10. 分阶段落地计划

### 第一阶段：PoC 验证

目标：验证技术链路可行。

范围：

```text
OpenTelemetry Java Agent
OpenTelemetry Go SDK
OpenTelemetry Collector
Tempo 单体或简化分布式部署
Grafana Trace 查询
```

验证项：

```text
HTTP 链路追踪
RPC 链路追踪
数据库 Span
Redis Span
MQ Trace 上下文传播
TraceId 与日志关联
Grafana 查询体验
```

---

### 第二阶段：平台化接入

目标：建立统一接入标准。

建设内容：

```text
统一接入文档
统一 SDK 封装
Java Agent 启动模板
Go SDK 初始化模板
Collector Helm Chart
标签规范
TraceId 日志规范
基础 Dashboard
```

---

### 第三阶段：生产级治理

目标：控制成本和风险。

建设内容：

```text
Tail Sampling
租户识别
敏感字段脱敏
高基数字段治理
Collector 限流
Collector 多集群部署
Tempo Distributed
对象存储生命周期策略
接入质量监控
```

---

### 第四阶段：深度诊断平台

目标：从“看链路”升级到“定位问题”。

建设内容：

```text
慢 Trace 自动分析
错误 Trace 聚合
服务拓扑
接口 SLA
发布前后 Trace 对比
服务依赖变更分析
熔断/限流事件关联 Trace
配置变更关联 Trace
告警跳转 Trace
Trace 跳转日志
```

---

## 11. 推荐最终技术栈

### 11.1 核心组件

| 模块       | 推荐选型                          |
| -------- | ----------------------------- |
| 采集标准     | OpenTelemetry                 |
| Java 接入  | OpenTelemetry Java Agent      |
| Go 接入    | OpenTelemetry Go SDK          |
| 采集网关     | OpenTelemetry Collector       |
| Trace 后端 | Grafana Tempo Distributed     |
| 展示入口     | Grafana                       |
| Metrics  | Prometheus / Mimir            |
| Logs     | Loki / Elasticsearch          |
| 存储       | S3 / OSS / COS / MinIO / Ceph |
| 部署       | Kubernetes + Helm             |
| 协议       | OTLP 为主，兼容 Jaeger / Zipkin    |

---

### 11.2 兼容策略

```text
新系统：
  OpenTelemetry SDK / Agent + OTLP

已有 Jaeger 系统：
  短期保留 Jaeger 协议
  Collector 接收后转发 Tempo
  中长期迁移到 OpenTelemetry SDK

已有 Zipkin 系统：
  短期保留 Zipkin receiver
  Collector 接收后转发 Tempo
  中长期迁移到 OTLP

Java APM 强诉求系统：
  可评估 SkyWalking 作为补充
  但不作为统一 Trace 底座
```

---

## 12. 最终建议

大型企业跨语言微服务链路追踪平台，推荐采用：

```text
OpenTelemetry Collector + Grafana Tempo Distributed + Grafana
```

推荐理由总结如下：

1. **OpenTelemetry 是跨语言、跨后端、vendor-neutral 的标准化采集方案**，可以避免业务代码绑定 Jaeger、Zipkin 或某个商业平台。([OpenTelemetry][1])
2. **OpenTelemetry Collector 是企业级遥测数据治理层**，适合做采样、脱敏、路由、租户识别、协议转换和多后端分发。([OpenTelemetry][3])
3. **Grafana Tempo 是高扩展 Trace 后端，面向大规模 tracing，并以对象存储降低长期存储成本**。([Grafana Labs][4])
4. **Tempo 支持多租户隔离**，适合大型企业按组织、团队、环境、地域进行平台化治理。([Grafana Labs][6])
5. **Grafana 体系可以形成 Metrics、Logs、Traces 的排障闭环**，更适合统一可观测平台建设。([Grafana Labs][4])
6. **Jaeger 原生 SDK 已不适合作为新平台路线**，Jaeger 社区和 OpenTelemetry 迁移文档均推荐迁移到 OpenTelemetry。([OpenTelemetry][10])
7. **Zipkin 更适合历史兼容和轻量场景**，不适合作为大型企业新建统一链路追踪平台底座。([zipkin.io][15])
8. **SkyWalking 适合 Java APM 产品化场景**，但如果目标是跨语言、开放标准、低成本、高扩展的统一 Trace 底座，OpenTelemetry + Tempo 更合适。([Apache SkyWalking][12])

最终架构建议：

```text
业务服务
  -> OpenTelemetry SDK / Agent
  -> OpenTelemetry Collector
  -> Grafana Tempo Distributed
  -> Object Storage
  -> Grafana
```

同时保留兼容入口：

```text
OTLP：主协议
Jaeger：历史系统兼容
Zipkin：历史系统兼容
```

该方案既符合当前云原生可观测技术趋势，也更适合大型企业长期建设统一链路追踪平台。

[1]: https://opentelemetry.io/docs/?utm_source=chatgpt.com "Documentation"
[2]: https://opentelemetry.io/?utm_source=chatgpt.com "OpenTelemetry"
[3]: https://opentelemetry.io/docs/collector/?utm_source=chatgpt.com "Collector"
[4]: https://grafana.com/docs/tempo/latest/?utm_source=chatgpt.com "Grafana Tempo documentation"
[5]: https://grafana.com/oss/tempo/?utm_source=chatgpt.com "Grafana Tempo OSS | Distributed tracing backend"
[6]: https://grafana.com/docs/tempo/latest/operations/manage-advanced-systems/multitenancy/?utm_source=chatgpt.com "Enable multi-tenancy | Grafana Tempo documentation"
[7]: https://opentelemetry.io/docs/concepts/sampling/?utm_source=chatgpt.com "Sampling"
[8]: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/tailsamplingprocessor/README.md?utm_source=chatgpt.com "Tail Sampling Processor"
[9]: https://www.jaegertracing.io/docs/latest/?utm_source=chatgpt.com "Introduction"
[10]: https://opentelemetry.io/docs/compatibility/migration/?utm_source=chatgpt.com "Migration"
[11]: https://github.com/apache/skywalking/blob/master/docs/README.md?utm_source=chatgpt.com "skywalking/docs/README.md at master"
[12]: https://skywalking.apache.org/?utm_source=chatgpt.com "Apache SkyWalking - Apache Software Foundation"
[13]: https://github.com/jaegertracing/jaeger-client-go?utm_source=chatgpt.com "jaegertracing/jaeger-client-go: 🛑 This library is ..."
[14]: https://github.com/jaegertracing/jaeger-client-java?utm_source=chatgpt.com "jaegertracing/jaeger-client-java: 🛑 This library is ..."
[15]: https://zipkin.io/?utm_source=chatgpt.com "OpenZipkin · A distributed tracing system"
[16]: https://zipkin.io/pages/architecture.html?utm_source=chatgpt.com "Architecture · OpenZipkin"
