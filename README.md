# 两仪

`两仪` 是一套自研微服务中间件体系，其正式英文总品牌为 `Stellar Axis`。这套体系以“星轴”作为工程表达，以“八卦”作为秩序骨架，以“灵觉层”作为 AI 扩展方向，面向高并发、高可用、可治理、可演进的分布式系统场景。

本仓库当前聚焦四类内容：

- 总品牌与体系定位
- 核心组件命名与职责映射
- 仓库布局与模块边界
- 依赖约束与架构决策

## 核心定位

两仪不是若干开源组件的简单拼装，而是一套统一命名、统一分层、统一依赖约束的中间件体系。整体设计遵循以下原则：

- 以 `Stellar Axis（两仪）` 作为唯一总品牌
- 以 `英文正式名 + 中文八卦名` 作为核心组件展示方式
- 以英文工程名作为仓库、模块、依赖坐标的唯一命名口径
- 以协议、SDK、事件和控制面作为跨产品协作边界

双栈定位如下：

- Java 框架：`Stellar Steel（两仪·刚）`
- Go 框架：`Stellar Titan（两仪·柔）`

## 核心组件矩阵

| 领域 | 英文正式名 | 中文正式名 | 中文简称 | 职责 |
| :--- | :--- | :--- | :--- | :--- |
| 服务注册与发现 | `StarMap` | `乾仪` | `乾` | 服务元数据、实例发现、注册协调 |
| 配置中心 | `Nebula` | `坤仪` | `坤` | 动态配置、环境承载、配置下发 |
| 链路追踪 | `LightBeam` | `离鉴` | `离` | 观测、链路跟踪、请求照明 |
| 服务治理 | `Orbit` | `巽策` | `巽` | 路由、灰度、流量治理 |
| 限流熔断 | `Pulsar` | `艮闸` | `艮` | 限流、熔断、过载保护 |
| 分布式调度 | `Chronos` | `震策` | `震` | 定时任务、分布式调度 |
| 分布式锁 | `Singularity` | `坎锁` | `坎` | 并发竞争控制、唯一性保护 |
| 网关入口 | `EventHorizon` | `兑门` | `兑` | 网关、边界入口、协议汇聚 |
| 消息与事件流 | `TaiJi Flow` | `太极流` | `中枢` | MQ、事件流、异步解耦 |

推荐展示方式：

- `StarMap · 乾`
- `Nebula · 坤`
- `Orbit · 巽`
- `TaiJi Flow · 太极流`

## AI 灵觉层

在八卦治理矩阵之外，体系还定义了独立的 AI 灵觉层：

| 英文正式名 | 中文正式名 | 职责 |
| :--- | :--- | :--- |
| `Spirit Axis` | `灵觉层` | AI 能力总体品牌 |
| `ZhongFu Engine` | `中孚` | 模型推理与理解 |
| `DaChu Memory` | `大畜` | 记忆与知识蓄积 |
| `Sui Agent` | `随位` | Agent 编排与自治执行 |
| `Xian MCP` | `咸` | 上下文感知协议 |
| `Heng MCP` | `恒` | 工具执行协议 |
| `Meng MCP` | `蒙` | Prompt 引导协议 |

灵觉层不强行纳入八卦八位，而是作为治理层之上的认知扩展层存在。

## 工程命名口径

以下口径已经固定：

- 总品牌：`Stellar Axis（两仪）`
- Java 聚合仓库：`stellar-steel`
- Go 聚合仓库：`stellar-titan`
- 控制平面：`stellar-control-plane`
- CLI：`stellarctl`
- 公开依赖坐标根命名空间：`io.stellar.axis`

不再使用以下过渡命名：

- `Qian-Directory`
- `Kun-Profile`
- `Li-Trace`
- `Xun-Pilot`
- `Gen-Brake`
- `Zhen-Ticker`
- `Kan-Vault`
- `Dui-Portal`

## 仓库布局

整个体系采用“品牌聚合仓库 + 核心产品仓库 + 配套能力仓库”的三层结构。

顶层聚合仓库：

- `stellar-axis`
- `stellar-steel`
- `stellar-titan`
- `spirit-axis`
- `stellar-control-plane`
- `stellarctl`
- `stellar-examples`
- `stellar-deploy`

核心产品仓库：

- `starmap`
- `nebula`
- `lightbeam`
- `orbit`
- `pulsar`
- `chronos`
- `singularity`
- `event-horizon`
- `taiji-flow`

单产品统一目录骨架：

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

## 依赖规则

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

## 文档导航

- [命名提案](./Proposal.md)
- [命名规范](./Naming-Convention.md)
- [仓库布局](./Repo-Layout.md)
- [模块依赖](./Module-Dependency.md)
- [架构决策记录](./Architecture-Decision-Record.md)

## 架构口号

> 乾坤定基，巽离观行；艮坎守正，震兑通灵。

## License

Apache License 2.0
