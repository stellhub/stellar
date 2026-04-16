# 星轴

`星轴` 的正式英文总品牌为 `Stellar Axis`。本仓库用于沉淀 `Stellar Axis（星轴）` 体系的最终方案，并作为品牌、产品、工程命名、仓库布局和依赖边界的唯一权威文档。

当前仓库只保留两份最终文档：

- 中文版：`README.md`
- 英文版：`README-EN.md`

所有过程性提案、路线图、初始化计划和阶段性描述均已收敛，不再单独保留。

## 最终品牌结论

- 英文总品牌：`Stellar Axis`
- 中文总品牌：`星轴`
- 首次标准展示：`Stellar Axis（星轴）`

双栈定位如下：

- Java 框架：`Stellar Core（星核）`
- Go 框架：`Stellar Pulse（星脉）`

使用规则如下：

- 官网、README、方案文档首次出现时，统一使用 `Stellar Axis（星轴）`
- 中文语境可直接使用 `星轴`
- 英文语境可直接使用 `Stellar Axis`
- 不再保留任何传统哲学别名、拼音品牌或混合命名

## 核心组件命名

核心组件统一采用：

> `英文正式名 · 中文宇宙名`

### 核心中间件矩阵

| 领域 | 英文正式名 | 中文正式名 |
| :--- | :--- | :--- |
| 注册中心 | `StarMap` | `星图` |
| 配置中心 | `Nebula` | `星云` |
| 链路追踪 | `StarTrace` | `星迹` |
| 服务治理 | `Orbit` | `星轨` |
| 流控熔断 | `Pulsar` | `脉冲` |
| 任务调度 | `Astrolabe` | `星盘` |
| 分布式锁 | `Singularity` | `奇点` |
| 网关入口 | `EventHorizon` | `视界` |
| 消息队列 | `CometFlow` | `彗流` |
| 指标平台 | `Constellation` | `星座` |
| 告警平台 | `NovaSignal` | `星讯` |
| 零信任平台 | `StarShield` | `星盾` |
| 密钥中心 | `StarKey` | `星钥` |

推荐展示方式：

- `StarMap · 星图`
- `Nebula · 星云`
- `StarTrace · 星迹`
- `Orbit · 星轨`
- `CometFlow · 彗流`

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

- `starmap-client-spring-boot-starter`
- `nebula-client-spring-boot-starter`
- `orbit-governance-starter`
- `comet-flow-client`

### OTel SDK 命名最终结论

OpenTelemetry 相关 SDK 统一采用：

> `{product}-otel-{signal}-{lang}-sdk`

其中：

- `{product}` 表示产品归属，例如 `startrace`
- `{signal}` 表示信号类型，例如 `logs`、`traces`、`metrics`
- `{lang}` 表示语言维度，例如 `java`、`go`

当前已经冻结的最终命名如下：

- `startrace-otel-logs-java-sdk`
- `startrace-otel-logs-go-sdk`

对应坐标示例：

- `io.stellar.axis:startrace-otel-logs-java-sdk`
- `io.stellar.axis:startrace-otel-logs-go-sdk`

扩展示例如下：

- `startrace-otel-traces-java-sdk`
- `startrace-otel-metrics-java-sdk`
- `startrace-otel-traces-go-sdk`
- `startrace-otel-metrics-go-sdk`

如果未来需要语言级聚合模块，可单独使用：

- `startrace-otel-java-sdk`
- `startrace-otel-go-sdk`

但这两个名称只建议用于聚合层，不建议用于具体的 `logs` SDK。

### Java 包名建议

- `io.stellar.axis.starmap.client`
- `io.stellar.axis.nebula.config`
- `io.stellar.axis.orbit.router`
- `io.stellar.axis.pulsar.limiter`
- `io.stellar.axis.cometflow.producer`

## 网关家族命名

网关家族主名统一使用 `EventHorizon · 视界`，通过二级限定词区分不同职责。

| 场景 | 英文建议 | 中文建议 |
| :--- | :--- | :--- |
| 内部网关 | `EventHorizon Internal` | `视界·内域` |
| 外部网关 | `EventHorizon External` | `视界·外域` |
| LLM 网关 | `EventHorizon LLM` | `视界·智域` |
| 边缘网关 | `EventHorizon Edge` | `视界·边域` |

对应仓库名建议：

- `event-horizon-internal`
- `event-horizon-external`
- `event-horizon-llm`
- `event-horizon-edge`

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

- `starmap`
- `nebula`
- `startrace`
- `orbit`
- `pulsar`
- `astrolabe`
- `singularity`
- `event-horizon`
- `comet-flow`

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

## 最终摘要

本仓库已经固定以下最终结论：

1. 唯一总品牌为 `Stellar Axis（星轴）`
2. 核心组件统一使用 `英文正式名 · 中文宇宙名`
3. 仓库、模块、包名、依赖坐标统一只使用英文工程名
4. 仓库结构固定为聚合仓库、核心产品仓库和配套能力仓库三层
5. 跨产品协作统一通过公开 API、SDK、事件和控制面完成
6. OpenTelemetry Logs SDK 最终命名固定为 `startrace-otel-logs-java-sdk` 与 `startrace-otel-logs-go-sdk`

## License

Apache License 2.0
