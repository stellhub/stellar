---
title: 体系总览
outline: deep
---

# 星际枢纽

`星际枢纽` 的正式英文总品牌为 `Stell Hub`。本仓库用于沉淀 `Stell Hub（星际枢纽）` 体系的最终方案，并作为品牌、产品、工程命名、仓库布局和依赖边界的唯一权威文档。

当前仓库以一份站点首页和若干规范文档为主：

- 首页：`docs/index.md`
- 体系总览：`docs/overview.md`
- 核心产品总览：`docs/products/index.md`
- 可观测规范：`docs/products/observability-spec.md`

## 最终品牌结论

- 英文总品牌：`Stell Hub`
- 中文总品牌：`星际枢纽`
- 首次标准展示：`Stell Hub（星际枢纽）`

双栈定位如下：

- Java 框架：`Stellar Core（星核）`
- Go 框架：`Stellar Pulse（星脉）`

使用规则如下：

- 官网、README、方案文档首次出现时，统一使用 `Stell Hub（星际枢纽）`
- 中文语境可直接使用 `星际枢纽`
- 英文语境可直接使用 `Stell Hub`

## 核心组件命名

核心组件统一采用：

> `Stell 风格命名 · 原始意象`

### 核心中间件矩阵

| 领域 | 原始意象 | **Stell 风格命名** | 命名逻辑解析 |
| :--- | :--- | :--- | :--- |
| **注册中心** | 星图 (StarMap) | **Stellmap** | Stell + Map，定位服务的地理坐标图。 |
| **配置中心** | 星云 (Nebula) | **Stellnula** | Stell + Nebula 缩写，像星云一样无处不在。 |
| **链路追踪** | 星迹 (StarTrace) | **Stelltrace** | Stell + Trace，追踪请求在星际间的流转。 |
| **服务治理** | 星轨 (Orbit) | **Stellorbit** | Stell + Orbit，让服务在预定的轨道上运行。 |
| **流控熔断** | 脉冲 (Pulsar) | **Stellpulse** | Stell + Pulse，精准掌控流量脉搏。 |
| **任务调度** | 星盘 (Astrolabe) | **Stellabe** | Stell + Astrolabe 缩写，古老的星际定位仪。 |
| **分布式锁** | 奇点 (Singularity) | **Stellpoint** | 奇点即一点。Stell + Point，锁定唯一的执行点。 |
| **网关入口** | 视界 (EventHorizon) | **Stellgate** | 网关通常叫 Gate。Stell + Gate，星际之门。 |
| **消息队列** | 彗流 (CometFlow) | **Stellflow** | Stell + Flow，像彗星尾迹一样的有序数据流。 |
| **日志平台** | 星谱 (Spectrum) | **Stellspec** | Stell + Spectrum 缩写，分析系统“光谱”异常。 |
| **指标平台** | 星座 (Constellation) | **Stellcon** | Stell + Constellation 缩写，由指标点连成的星座。 |
| **告警平台** | 星讯 (NovaSignal) | **Stellvox** | Vox 在拉丁语中意为声音/讯息。Stell + Vox。 |
| **零信平台** | 星盾 (StarShield) | **Stellguard** | Stell + Guard，星际守卫，比 Shield 更有动态感。 |
| **密钥中心** | 星钥 (StarKey) | **Stellkey** | Stell + Key，简单有力，符合密钥中心定位。 |

推荐展示方式：

- `Stellmap · 星图`
- `Stellnula · 星云`
- `Stelltrace · 星迹`
- `Stellspec · 星谱`
- `Stellorbit · 星轨`
- `Stellflow · 彗流`

对应详细设计文档入口：

- [Stellmap · 星图](/products/stellmap/)
- [Stellnula · 星云](/products/stellnula/)
- [Stelltrace · 星迹](/products/stelltrace/)
- [Stellorbit · 星轨](/products/stellorbit/)
- [Stellpulse · 脉冲](/products/stellpulse/)
- [Stellabe · 星盘](/products/stellabe/)
- [Stellpoint · 奇点](/products/stellpoint/)
- [Stellgate · 视界](/products/stellgate/)
- [Stellflow · 彗流](/products/stellflow/)
- [Stellspec · 星谱](/products/stellspec/)
- [Stellcon · 星座](/products/stellcon/)
- [Stellvox · 星讯](/products/stellvox/)
- [Stellguard · 星盾](/products/stellguard/)
- [Stellkey · 星钥](/products/stellkey/)

矩阵配图如下：

![核心中间件矩阵图](/core-middleware-matrix-dependency.svg)

### AI 星穹层矩阵

| 英文正式名 | 中文正式名 | 职责 |
| :--- | :--- | :--- |
| `Astral Layer` | `星穹层` | AI 能力总体品牌 |
| `Quasar Engine` | `类星引擎` | 模型推理与理解 |
| `StarVault Memory` | `星库` | 记忆与知识蓄积 |
| `Orbit Agent` | `轨使` | Agent 编排与自治执行 |
| `Sensor MCP` | `星感` | 上下文感知协议 |
| `Vector MCP` | `星行` | 工具执行协议 |
| `GuideStar MCP` | `导星` | Prompt 引导协议 |

## 工程命名规范

工程实体统一只使用英文，不挂中文宇宙名。

适用对象：

- Git 仓库名
- Maven `artifactId`
- Go module
- Java 包名
- CLI 工具名
- `starter`、`client`、`sdk`、`sidecar`、`operator` 等子模块

统一规则如下：

- 全部使用小写英文
- 单词之间使用中划线 `-`
- 不使用中文
- 不使用拼音
- 不使用中英混拼

### 固定工程命名

- Java 聚合仓库：`stellar-core`
- Go 聚合仓库：`stellar-pulse`
- AI 聚合仓库：`astral-layer`
- 控制平面：`stellar-control-plane`
- CLI：`stellarctl`
- 公开依赖坐标根命名空间：`io.stellar.axis`
- Go 上层命名空间：`github.com/stellar-axis/`

### 工程命名模式

推荐使用：

> `{product}-{capability}-{role}`

常见示例：

- `stellmap-client-spring-boot-starter`
- `stellnula-client-spring-boot-starter`
- `stellorbit-governance-starter`
- `stellflow-client`

### 日志平台 SDK 命名最终结论

`Stellspec` 作为日志平台主名，对外 SDK 统一采用简化命名，不再把 `otel` 或信号类型写入 SDK 主名称。

当前已经冻结的最终命名如下：

- `stellspec-java-sdk`
- `stellspec-go-sdk`

对应坐标示例：

- `io.stellar.axis:stellspec-java-sdk`
- `io.stellar.axis:stellspec-go-sdk`

命名含义如下：

- `stellspec` 表示日志平台产品归属
- `java`、`go` 表示语言维度
- `sdk` 表示工程角色

`OpenTelemetry` 可作为内部实现能力或适配能力存在，但不进入当前对外 SDK 主命名。

### Java 包名建议

- `io.stellar.axis.stellmap.client`
- `io.stellar.axis.stellnula.config`
- `io.stellar.axis.stellorbit.router`
- `io.stellar.axis.stellpulse.limiter`
- `io.stellar.axis.stellflow.producer`

## 网关家族命名

网关家族主名统一使用 `Stellgate · 视界`，通过二级限定词区分不同职责。

| 场景 | 英文建议 | 中文建议 |
| :--- | :--- | :--- |
| 内部网关 | `Stellgate Internal` | `视界·内域` |
| 外部网关 | `Stellgate External` | `视界·外域` |
| LLM 网关 | `Stellgate LLM` | `视界·智域` |
| 边缘网关 | `Stellgate Edge` | `视界·边域` |

对应仓库名建议：

- `stellgate-internal`
- `stellgate-external`
- `stellgate-llm`
- `stellgate-edge`

## 仓库布局最终结论

体系采用“品牌聚合仓库 + 核心产品仓库 + 配套能力仓库”的三层结构。

### 顶层聚合仓库

- `stellar-axis`
- `stellar-core`
- `stellar-pulse`
- `astral-layer`
- `stellar-control-plane`
- `stellarctl`
- `stellar-examples`
- `stellar-deploy`

### 核心产品仓库

- `stellmap`
- `stellnula`
- `stelltrace`
- `stellorbit`
- `stellpulse`
- `stellabe`
- `stellpoint`
- `stellgate`
- `stellflow`
- `stellspec`
- `stellcon`
- `stellvox`
- `stellguard`
- `stellkey`

### 单产品统一目录骨架

```text
{product}/
├── docs/
├── api/
├── server/
├── clients/
├── starters/
├── sidecars/
├── operators/
├── deploy/
├── test/
└── examples/
```

模块约束如下：

- `docs/`、`api/`、`server/`、`clients/`、`deploy/`、`test/` 为核心骨架
- `starters/` 仅在提供 Java 接入时创建
- `sidecars/` 仅在存在边车形态时创建
- `operators/` 仅在存在 Kubernetes 生命周期管理需求时创建
- `examples/` 仅在需要独立最小示例时创建

## 依赖边界最终结论

整套体系统一采用单向依赖约束。

允许的原则：

- 依赖公开 API、SDK 和稳定协议
- 通过控制面、事件流和官方客户端实现跨产品协作
- `starter` 依赖 `clients/api`
- `sidecar` 依赖 `clients/api/runtime`
- `operator` 依赖 `api` 和 Kubernetes Runtime

禁止的原则：

- `starter -> server`
- `client -> server`
- `operator -> server`
- `control-plane -> product/internal`
- `productA/server -> productB/server`
- `examples -> product/internal`

进一步约束如下：

- `stellar-axis` 只承载规范与入口，不承载运行时代码
- `stellar-core` 和 `stellar-pulse` 只能依赖公开契约，不能依赖具体服务端实现
- 核心产品之间只能通过 API、SDK、事件、配置或控制面交互
- 控制平面只能依赖公开管理面，不能依赖产品内部实现

## 全局运行规范

除品牌、命名、仓库布局与依赖边界之外，体系还固定以下全局运行规范：

- 可观测规范：`docs/products/observability-spec.md`

适用范围如下：

- 所有业务服务
- 所有中间件 SDK
- 所有网关、Agent、Sidecar、Collector
- 所有部署平台与运行平台

约束原则如下：

- 平台统一注入基础环境变量，所有中间件默认识别
- 标准链路头与平台扩展 Header 统一透传、统一清洗、统一映射
- 指标命名、标签边界与聚合口径统一
- 各产品可做局部扩展，但不得破坏全局语义

## 最终摘要

本仓库已经固定以下最终结论：

1. 唯一总品牌为 `Stell Hub（星际枢纽）`
2. 核心组件统一使用 `英文正式名 · 中文宇宙名`
3. 仓库、模块、包名、依赖坐标统一只使用英文工程名
4. 仓库结构固定为聚合仓库、核心产品仓库和配套能力仓库三层
5. 跨产品协作统一通过公开 API、SDK、事件和控制面完成
6. 日志平台最终命名固定为 `Stellspec · 星谱`，对应 SDK 命名固定为 `stellspec-java-sdk` 与 `stellspec-go-sdk`
7. 全体系基础环境变量、请求上下文与指标语义统一遵循 `docs/products/observability-spec.md`

## License

Apache License 2.0
