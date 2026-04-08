# Proposal: 星轴 × 八卦微服务中间件命名方案

## 1. Problem Analysis

当前文档已经形成了两套命名语义：

- 一套是偏东方哲学表达的 **两仪 / 八卦 / 太极 / 灵觉层**。
- 一套是偏国际化工程表达的 **Stellar Axis / Stellar-* / StarMap / Nebula**。

这两套体系各自成立，但如果直接并行使用，后续会出现几个问题：

- 对外品牌、对内模块、代码仓库、SDK 包名之间容易出现风格漂移。
- 中文命名具备辨识度，但英文世界传播成本较高。
- 英文命名足够工程化，但如果脱离八卦体系，会削弱整套自研中间件的统一叙事。
- 后续继续扩展消息、AI、平台层时，如果没有统一规则，新组件容易失去命名秩序。

因此，这份提案的目标不是再发明第三套名字，而是把现有 README 中已经出现的两套世界观合并为一套“可持续扩展”的命名规则。

## 2. Design

### 2.1 命名总原则

建议采用“三层命名模型”：

- **品牌层**：统一使用 `Stellar Axis / 两仪`。
- **产品层**：统一使用“星轴英文名 + 八卦中文名”双轨表达。
- **工程层**：代码仓库、包名、制品名统一优先英文，中文只保留在文档、官网、注释和品牌叙事中。

这样做的理由：

- `Stellar Axis` 负责对外传播、开源生态、仓库命名、依赖坐标。
- `两仪` 负责中文世界观总名，承接“刚柔双栈”的哲学解释。
- 八卦负责每个核心中间件的“角色定位”。
- 星轴命名负责每个中间件的“工程识别度”。

### 2.2 建议的品牌架构

推荐采用以下结构：

| 层级 | 英文命名 | 中文命名 | 用途 |
| :--- | :--- | :--- | :--- |
| 总品牌 | **Stellar Axis** | **两仪** | 整套微服务中间件体系总称 |
| Java 技术栈 | **Stellar Steel** | **两仪·刚** | Java 侧框架与 SDK 发行体系 |
| Go 技术栈 | **Stellar Titan** | **两仪·柔** | Go 侧框架与 SDK 发行体系 |
| 治理矩阵 | **Axis Matrix** | **八卦阵** | 对八大核心中间件的总称 |
| 消息中枢 | **TaiJi Flow** / **CometFlow** | **太极流** | 消息、事件、异步流转层 |
| AI 灵觉层 | **Spirit Axis** | **灵觉层** | LLM、Memory、Agent、MCP 总称 |

其中：

- `Stellar Axis / 两仪` 作为总品牌最稳定，不建议再改。
- `Stellar Steel` 与 `Stellar Titan` 已经与“刚柔双栈”高度一致，建议保留。
- `Axis Matrix / 八卦阵` 可以作为官网、架构图、产品页中的总矩阵名。
- 消息中枢建议从 README 现有表述收敛到一个主名，避免 `TaiJi-Flow` 和 `Comet` 长期并行。

### 2.3 核心命名决策

#### 决策 A：总品牌不再切换

总品牌固定为：

- 英文：`Stellar Axis`
- 中文：`两仪`

不建议把总品牌命名成 `LiangYi` 为主、`Stellar Axis` 为辅。原因是：

- 英文品牌在仓库、域名、包名、Maven GroupId、Go Module 中更自然。
- `LiangYi` 更适合做中文世界观的拼音注解，而非唯一国际品牌。

因此推荐表达方式：

> **Stellar Axis（两仪）**

#### 决策 B：组件统一采用“双名制”

每个核心中间件统一采用：

> **English Product Name + Chinese Bagua Name**

例如：

- `StarMap · 乾`
- `Nebula · 坤`
- `LightBeam · 离`
- `Orbit · 巽`

这样兼顾三件事：

- 官网展示时有品牌辨识度。
- 代码仓库时有工程可读性。
- 中文语境下仍然保留东方哲学记忆点。

#### 决策 C：英文名负责工程，八卦名负责秩序

后续所有扩展组件遵循一个规则：

- 如果是核心基础设施组件，必须挂接到“星轴语义 + 八卦或周易语义”。
- 如果是非核心配套工具，则只保留 `Stellar-*` 工程前缀即可。

例如：

- 核心组件：`StarMap · 乾`
- 配套 CLI：`stellarctl`
- Java BOM：`stellar-steel-bom`
- Go Sidecar：`stellar-titan-sidecar`

### 2.4 核心组件推荐命名表

以下是基于当前 README 内容整理后的推荐主命名：

| 领域 | 推荐英文名 | 推荐中文名 | 卦位/意象 | 命名结论 |
| :--- | :--- | :--- | :--- | :--- |
| 服务注册中心 | **StarMap** | **乾仪** / **乾** | 乾 | 保留 `StarMap` 为主名，中文可简称“乾” |
| 配置中心 | **Nebula** | **坤仪** / **坤** | 坤 | 保留 `Nebula`，强调承载与孕育 |
| 链路追踪 | **LightBeam** | **离鉴** / **离** | 离 | 保留 `LightBeam`，中文建议强调“照见” |
| 服务治理 | **Orbit** | **巽策** / **巽** | 巽 | 保留 `Orbit`，天然适合路由与治理 |
| 限流熔断 | **Pulsar** | **艮闸** / **艮** | 艮 | 保留 `Pulsar`，中文建议体现“闸门” |
| 分布式任务 | **Chronos** | **震策** / **震** | 震 | `Chronos` 易理解，建议保留 |
| 分布式锁 | **Singularity** | **坎锁** / **坎** | 坎 | 保留 `Singularity`，突出唯一性 |
| 智能网关 | **EventHorizon** | **兑门** / **兑** | 兑 | 保留 `EventHorizon`，极具边界感 |
| 消息队列 / 事件流 | **TaiJi Flow** | **太极流** | 中枢 | 建议作为最终主名，不再与 `Comet` 并列 |

说明：

- 英文名建议继续沿用 README-EN 中现有版本，因为整体质量已经较高，且语义风格统一。
- 中文名不建议继续使用 `Qian-Directory`、`Kun-Profile` 这种“中英混拼”形态作为长期产品名。
- 中文侧更适合采用“卦位简称”或“卦位 + 职能字”的命名，例如“兑门”“艮闸”“离鉴”。

### 2.5 中文命名风格建议

中文产品名建议统一采用以下优先级：

1. **正式品牌名**：`StarMap · 乾`
2. **中文宣传名**：`乾仪`、`坤仪`、`兑门`
3. **内部简称**：`乾`、`坤`、`兑`

不推荐继续扩展的风格：

- `Qian-Directory`
- `Kun-Profile`
- `Dui-Portal`

原因：

- 这是过渡命名，不够纯粹。
- 既不利于中文传播，也不利于英文传播。
- 长期会增加文档和仓库名的不一致。

### 2.6 AI 灵觉层命名建议

README 中 AI 命名已具备较强辨识度，建议保留并做层次收束：

| 层级 | 英文建议 | 中文建议 | 说明 |
| :--- | :--- | :--- | :--- |
| AI 总层 | **Spirit Axis** | **灵觉层** | AI Native 扩展总称 |
| LLM 推理 | **ZhongFu** / **ZhongFu Engine** | **中孚** | 推理与理解中枢 |
| Memory | **DaChu** / **DaChu Memory** | **大畜** | 长短期知识蓄积 |
| Agent | **Sui** / **Sui Agent** | **随位** | 执行编排与自治代理 |
| MCP Context | **Xian MCP** | **咸** | 感知资源与上下文 |
| MCP Tools | **Heng MCP** | **恒** | 工具调用与执行 |
| MCP Prompt | **Meng MCP** | **蒙** | 智能引导与提示编排 |

建议：

- AI 子体系尽量不要再强行纳入八卦八位，而是作为“八卦治理阵”之外的“灵觉层”。
- 这能保持架构的分层清晰：底层治理负责秩序，灵觉层负责认知。

### 2.7 工程命名规范

建议统一如下：

#### 仓库命名

- `stellar-axis`：总仓库或总组织
- `starmap`
- `nebula`
- `lightbeam`
- `orbit`
- `pulsar`
- `chronos`
- `singularity`
- `event-horizon`
- `taiji-flow`

#### Java Maven 坐标

- GroupId：`io.stellar.axis`
- ArtifactId 示例：
  - `stellar-steel-core`
  - `starmap-client-spring-boot-starter`
  - `nebula-client-spring-boot-starter`
  - `orbit-governance-starter`
  - `taiji-flow-client`

#### Go Module

- `github.com/stellar-axis/stellar-titan`
- `github.com/stellar-axis/starmap`
- `github.com/stellar-axis/nebula`
- `github.com/stellar-axis/orbit`

#### CLI / 平台工具

- `stellarctl`
- `axisctl`
- `taijictl`

推荐优先：

- 总控工具使用 `stellarctl`

因为它与总品牌绑定最强，最利于记忆。

## 3. Implementation

基于以上分析，建议本项目正式采用以下命名方案。

### 3.1 最终推荐方案

#### 总品牌

- 英文：`Stellar Axis`
- 中文：`两仪`
- 标准展示：`Stellar Axis（两仪）`

#### 双栈品牌

- Java：`Stellar Steel（两仪·刚）`
- Go：`Stellar Titan（两仪·柔）`

#### 核心基础设施矩阵

| 英文主名 | 中文主名 | 中文简称 | 职责 |
| :--- | :--- | :--- | :--- |
| `StarMap` | `乾仪` | `乾` | 服务注册与发现 |
| `Nebula` | `坤仪` | `坤` | 配置中心 |
| `LightBeam` | `离鉴` | `离` | 链路追踪 |
| `Orbit` | `巽策` | `巽` | 服务治理与路由 |
| `Pulsar` | `艮闸` | `艮` | 限流、熔断、过载保护 |
| `Chronos` | `震策` | `震` | 分布式调度 |
| `Singularity` | `坎锁` | `坎` | 分布式锁 |
| `EventHorizon` | `兑门` | `兑` | 网关与流量入口 |
| `TaiJi Flow` | `太极流` | `中枢` | MQ / Event Streaming |

#### AI 灵觉层

| 英文主名 | 中文主名 | 职责 |
| :--- | :--- | :--- |
| `Spirit Axis` | `灵觉层` | AI 能力总体品牌 |
| `ZhongFu Engine` | `中孚` | LLM 推理 |
| `DaChu Memory` | `大畜` | Memory / RAG |
| `Sui Agent` | `随位` | Agent 编排 |
| `Xian MCP` | `咸` | 上下文感知 |
| `Heng MCP` | `恒` | 工具执行 |
| `Meng MCP` | `蒙` | Prompt 引导 |

### 3.2 命名使用规范

对外介绍时：

> Stellar Axis（两仪）是一套自研微服务中间件体系，其核心治理矩阵由 StarMap·乾、Nebula·坤、Orbit·巽 等组件构成。

在文档标题中：

- `StarMap · 乾 | Service Registry`
- `Nebula · 坤 | Config Center`
- `EventHorizon · 兑 | Intelligent Gateway`

在仓库和依赖命名中：

- 只使用英文名，不混用拼音卦名。

在中文文章或演讲中：

- 首次出现使用“双名制”。
- 后续可直接使用中文简称，如“乾”“坤”“兑门”“太极流”。

### 3.3 为什么这是当前最优解

这套方案同时满足四个目标：

- **统一性**：总品牌、组件名、SDK 名、AI 名称位于同一叙事体系下。
- **传播性**：英文名具备现代中间件品牌气质，适合开源与国际传播。
- **辨识度**：八卦中文名使体系区别于常见的英文直译中间件。
- **可扩展性**：后续新增组件时，可以继续遵循“星轴意象 + 东方秩序”的规则扩展。

## 4. Complete Code

以下内容可直接作为命名提案落地执行：

```md
# Proposal: 星轴 × 八卦微服务中间件命名方案

## 最终结论

本体系建议采用“**Stellar Axis（两仪）**”作为唯一总品牌。

- `Stellar Axis` 负责国际化传播、仓库命名、依赖坐标与工程生态。
- `两仪` 负责中文世界观表达，承接“刚柔双栈”的哲学母题。
- 核心中间件统一使用“**英文主名 + 八卦中文名**”的双名制。

## 双栈命名

- Java: `Stellar Steel（两仪·刚）`
- Go: `Stellar Titan（两仪·柔）`

## 核心组件命名

| 英文主名 | 中文主名 | 中文简称 | 职责 |
| :--- | :--- | :--- | :--- |
| `StarMap` | `乾仪` | `乾` | 服务注册与发现 |
| `Nebula` | `坤仪` | `坤` | 配置中心 |
| `LightBeam` | `离鉴` | `离` | 链路追踪 |
| `Orbit` | `巽策` | `巽` | 服务治理与路由 |
| `Pulsar` | `艮闸` | `艮` | 限流、熔断、过载保护 |
| `Chronos` | `震策` | `震` | 分布式调度 |
| `Singularity` | `坎锁` | `坎` | 分布式锁 |
| `EventHorizon` | `兑门` | `兑` | 智能网关 |
| `TaiJi Flow` | `太极流` | `中枢` | 消息队列 / 事件流 |

## AI 灵觉层命名

| 英文主名 | 中文主名 | 职责 |
| :--- | :--- | :--- |
| `Spirit Axis` | `灵觉层` | AI 总体能力层 |
| `ZhongFu Engine` | `中孚` | LLM 推理 |
| `DaChu Memory` | `大畜` | Memory / RAG |
| `Sui Agent` | `随位` | Agent 编排 |
| `Xian MCP` | `咸` | Context |
| `Heng MCP` | `恒` | Tools |
| `Meng MCP` | `蒙` | Prompt |

## 工程落地规则

1. 仓库、包名、依赖坐标统一使用英文名。
2. 官网、文档、架构图统一使用“双名制”。
3. 中文传播中优先使用“两仪、乾、坤、兑门、太极流、灵觉层”等名称。
4. 不再推荐继续使用 `Qian-Directory`、`Kun-Profile` 这类中英混拼命名。

## 推荐仓库示例

- `stellar-axis`
- `stellar-steel`
- `stellar-titan`
- `starmap`
- `nebula`
- `lightbeam`
- `orbit`
- `pulsar`
- `chronos`
- `singularity`
- `event-horizon`
- `taiji-flow`

## 推荐依赖坐标示例

- Maven GroupId: `io.stellar.axis`
- Java ArtifactId: `stellar-steel-core`
- Java ArtifactId: `starmap-client-spring-boot-starter`
- Go Module: `github.com/stellar-axis/starmap`

## 一句话定位

> Stellar Axis（两仪）是一套以星轴为工程表达、以八卦为秩序骨架、以灵觉层为智能扩展的自研微服务中间件体系。
```
