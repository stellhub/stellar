# Naming Convention

## 1. Problem Analysis

基于 [Proposal.md](./Proposal.md) 的决策，本文件用于把命名提案进一步收敛为可执行规范。

需要解决的核心问题有三个：

- **决策 A**：总品牌必须固定，不能在 `Stellar Axis`、`LiangYi`、`两仪` 之间来回切换。
- **决策 B**：核心组件必须有统一展示方式，不能有的用纯英文，有的用拼音混写，有的只写八卦。
- **决策 C**：工程命名和品牌命名必须分层，不能把宣传名称直接拿去当仓库名、包名、模块名。

如果没有这份规范，后续会很快出现以下问题：

- README、官网、仓库、SDK、代码包名使用不同口径。
- 核心中间件与配套工具的命名边界不清晰。
- 新增组件时没有统一规则，只能靠临时拍脑袋命名。

因此，这份文档的作用是给出最终命名结论，并解释为什么这样命名。

## 2. Design

### 2.1 总体命名原则

本体系采用三层命名结构：

- **品牌层**：定义整套体系的统一品牌。
- **产品层**：定义核心中间件与 AI 子体系的正式名称。
- **工程层**：定义仓库、包名、依赖、CLI、模块的命名格式。

三层之间的职责划分如下：

- 品牌层负责统一认知。
- 产品层负责统一展示。
- 工程层负责统一实现。

### 2.2 决策 A：总品牌固定

#### 最终命名

- 英文总品牌：`Stellar Axis`
- 中文总品牌：`两仪`
- 标准联合写法：`Stellar Axis（两仪）`

#### 命名解释

- `Stellar Axis` 强调“星轴”，适合作为国际化品牌、仓库组织名、域名和工程前缀。
- `两仪` 是中文世界观总名，承接“刚柔双栈”“阴阳相生”的哲学母题。
- `Stellar Axis（两仪）` 是唯一推荐的首次完整展示方式。

#### 使用规则

- 官网、白皮书、README 首次出现时，必须写作 `Stellar Axis（两仪）`。
- 在英文语境中，可单独使用 `Stellar Axis`。
- 在中文语境中，可单独使用 `两仪`。

#### 禁止写法

- 不使用 `LiangYi` 作为总品牌主名。
- 不使用 `LiangYi Framework` 作为正式对外总称。
- 不使用 `StellarAxis` 这种去空格拼接形式作为品牌展示名。

#### 原因说明

- `LiangYi` 更适合拼音注解，不适合作为全球工程品牌主名。
- `Stellar Axis` 在技术产品语境下更自然，也更适合作为仓库组织、域名和依赖坐标的上层品牌。

### 2.3 决策 B：核心组件统一采用双名制

#### 最终规则

所有核心中间件统一采用：

> **英文主名 + 中文八卦名**

标准写法：

> `EnglishName · 中文名`

例如：

- `StarMap · 乾`
- `Nebula · 坤`
- `LightBeam · 离`
- `Orbit · 巽`

#### 命名解释

- 英文主名负责工程辨识度与国际化传播。
- 中文八卦名负责秩序映射与哲学识别。
- 双名制使文档、官网、产品页、架构图可以同时表达“工程含义”和“卦位含义”。

#### 核心组件最终命名表

| 领域 | 英文正式名 | 中文正式名 | 中文简称 | 解释 |
| :--- | :--- | :--- | :--- | :--- |
| 服务注册中心 | `StarMap` | `乾仪` | `乾` | 天体坐标图，适合表达服务发现与节点定位 |
| 配置中心 | `Nebula` | `坤仪` | `坤` | 星云有承载、孕育之意，对应配置承载能力 |
| 链路追踪 | `LightBeam` | `离鉴` | `离` | 光束照见链路，对应离卦之明 |
| 服务治理 | `Orbit` | `巽策` | `巽` | 轨道天然对应流量调度、路由治理 |
| 限流熔断 | `Pulsar` | `艮闸` | `艮` | 脉冲星强调频率控制，艮强调阻断与止损 |
| 分布式调度 | `Chronos` | `震策` | `震` | Chronos 代表时间秩序，震代表发动与触发 |
| 分布式锁 | `Singularity` | `坎锁` | `坎` | 奇点表达唯一性与收束性，对应锁语义 |
| 智能网关 | `EventHorizon` | `兑门` | `兑` | 事件视界表达边界入口，兑对应开口与门户 |
| 消息流转中枢 | `TaiJi Flow` | `太极流` | `中枢` | 阴阳流转的中枢，适合作为消息与事件总线 |

#### 使用规则

- 首次出现使用 `英文正式名 · 中文简称` 或 `英文正式名 · 中文正式名`。
- 中文文档正文中，后续可以使用中文简称，如“乾”“坤”“兑门”“太极流”。
- 英文文档正文中，后续可以只使用英文正式名。

#### 推荐写法

- `StarMap · 乾 | Service Registry`
- `Nebula · 坤 | Config Center`
- `EventHorizon · 兑 | Gateway`

#### 不推荐写法

- `Qian-Directory`
- `Kun-Profile`
- `Dui-Portal`
- `Li-Trace`
- `Xun-Pilot`

#### 原因说明

- 这些中英混拼写法属于概念过渡阶段命名，不适合作为长期产品名。
- 它们在英文世界不自然，在中文语境中也不够纯粹。
- 统一切换到“双名制”后，品牌表达和工程表达会更稳定。

### 2.4 决策 C：英文负责工程，八卦负责秩序

#### 最终规则

- **核心基础设施组件**：必须同时具备英文工程名和中文八卦名。
- **配套工程工具**：只保留英文工程名，不强制挂八卦名。
- **工程实体**：仓库、模块、包名、依赖坐标统一只使用英文。

#### 命名解释

- 八卦名适用于“体系级角色”，不适用于所有工程对象。
- 工程对象强调可读性、可输入性、依赖管理与跨语言一致性，因此统一使用英文。
- 这样可以避免出现“中文品牌强、工程命名乱”的问题。

#### 适用范围

适合挂八卦名的对象：

- 核心中间件产品
- 架构图节点
- 官网产品卡片
- 中文方案文档

不适合挂八卦名的对象：

- Git 仓库名
- Maven ArtifactId
- Go Module
- CLI 工具名
- Sidecar、Operator、Starter、SDK 子模块

## 3. Implementation

### 3.1 最终命名总表

#### 总品牌

| 类型 | 最终名称 | 解释 |
| :--- | :--- | :--- |
| 英文总品牌 | `Stellar Axis` | 整套体系唯一英文总品牌 |
| 中文总品牌 | `两仪` | 整套体系唯一中文总品牌 |
| 首次标准展示 | `Stellar Axis（两仪）` | 官网、README、方案文档首次出现统一写法 |

#### 双栈框架

| 类型 | 最终名称 | 解释 |
| :--- | :--- | :--- |
| Java 框架 | `Stellar Steel` | 对应“两仪·刚”，强调厚重、完整、扩展性 |
| Go 框架 | `Stellar Titan` | 对应“两仪·柔”，强调高并发、轻量、边缘性能 |
| 中文 Java 名 | `两仪·刚` | 仅用于中文品牌展示 |
| 中文 Go 名 | `两仪·柔` | 仅用于中文品牌展示 |

#### 核心中间件

| 英文正式名 | 中文正式名 | 中文简称 | 最终解释 |
| :--- | :--- | :--- | :--- |
| `StarMap` | `乾仪` | `乾` | 服务元数据与发现中枢，代表天之统领 |
| `Nebula` | `坤仪` | `坤` | 系统配置承载底座，代表地之承载 |
| `LightBeam` | `离鉴` | `离` | 链路照明与观测系统，代表火之明察 |
| `Orbit` | `巽策` | `巽` | 路由、流量、治理策略中枢，代表风之渗透 |
| `Pulsar` | `艮闸` | `艮` | 限流、熔断、过载保护系统，代表山之止断 |
| `Chronos` | `震策` | `震` | 分布式时间调度系统，代表雷之发动 |
| `Singularity` | `坎锁` | `坎` | 分布式锁与唯一竞争控制，代表水之险陷 |
| `EventHorizon` | `兑门` | `兑` | 网关与入口边界系统，代表泽之开口 |
| `TaiJi Flow` | `太极流` | `中枢` | 消息与事件流转中心，代表阵法核心能量流 |

#### AI 灵觉层

| 英文正式名 | 中文正式名 | 最终解释 |
| :--- | :--- | :--- |
| `Spirit Axis` | `灵觉层` | AI Native 总体能力层 |
| `ZhongFu Engine` | `中孚` | 模型推理与理解中枢 |
| `DaChu Memory` | `大畜` | 记忆与知识蓄积层 |
| `Sui Agent` | `随位` | Agent 编排与自治执行层 |
| `Xian MCP` | `咸` | Context 感知协议族 |
| `Heng MCP` | `恒` | Tools 执行协议族 |
| `Meng MCP` | `蒙` | Prompt 引导协议族 |

### 3.2 工程命名规范

#### 仓库命名

统一规则：

- 全部使用小写英文
- 单词之间使用中划线 `-`
- 不使用拼音八卦名
- 不使用中文

推荐仓库名：

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
- `spirit-axis`

#### Maven 坐标命名

统一规则：

- GroupId 固定为 `io.stellar.axis`
- ArtifactId 使用英文工程名，不使用中文和拼音
- Starter、SDK、Client、BOM 等后缀表达工程角色

推荐示例：

- `io.stellar.axis:stellar-steel-core`
- `io.stellar.axis:stellar-steel-bom`
- `io.stellar.axis:starmap-client-spring-boot-starter`
- `io.stellar.axis:nebula-client-spring-boot-starter`
- `io.stellar.axis:orbit-governance-starter`
- `io.stellar.axis:taiji-flow-client`
- `io.stellar.axis:zhongfu-engine-client`

#### Go Module 命名

统一规则：

- 使用 `github.com/stellar-axis/` 作为上层命名空间
- 模块名使用英文工程名
- 不使用中文和拼音卦名

推荐示例：

- `github.com/stellar-axis/stellar-titan`
- `github.com/stellar-axis/starmap`
- `github.com/stellar-axis/nebula`
- `github.com/stellar-axis/orbit`
- `github.com/stellar-axis/taiji-flow`
- `github.com/stellar-axis/zhongfu-engine`

#### Java 包名命名

统一规则：

- 根包名固定为 `io.stellar.axis`
- 子包采用英文工程域
- 包结构体现职责，不体现中文卦位

推荐示例：

- `io.stellar.axis.starmap.client`
- `io.stellar.axis.nebula.config`
- `io.stellar.axis.orbit.router`
- `io.stellar.axis.pulsar.limiter`
- `io.stellar.axis.taijiflow.producer`

#### CLI / 平台工具命名

统一规则：

- 工具统一走英文工程名
- 总控工具绑定总品牌

最终建议：

- 总控 CLI：`stellarctl`
- 运维平台：`stellar-console`
- 管理平面：`stellar-control-plane`

#### Sidecar / Operator / Agent / Starter 命名

统一规则：

- 基础格式：`{product}-{role}`
- 角色后缀固定，不混用中文

推荐示例：

- `orbit-sidecar`
- `pulsar-sidecar`
- `starmap-operator`
- `nebula-operator`
- `stellar-steel-starter`
- `sui-agent-runtime`

### 3.3 文档与官网展示规范

#### 首次出现格式

首次出现时统一使用：

> `Stellar Axis（两仪）`

核心组件首次出现时统一使用：

> `EnglishName · 中文名`

例如：

- `StarMap · 乾`
- `Nebula · 坤`
- `TaiJi Flow · 太极流`

#### 后续简写规则

- 英文文档后续只写英文正式名。
- 中文文档后续可写中文简称。
- 架构图节点建议使用“双名制”。

#### 推荐标题格式

- `StarMap · 乾 | Service Registry`
- `Nebula · 坤 | Config Center`
- `Orbit · 巽 | Governance`
- `TaiJi Flow · 太极流 | Messaging`

### 3.4 禁止事项

以下写法统一不再使用：

- `Qian-Directory`
- `Kun-Profile`
- `Li-Trace`
- `Xun-Pilot`
- `Gen-Brake`
- `Zhen-Ticker`
- `Kan-Vault`
- `Dui-Portal`

禁止原因：

- 中英混拼，不利于国际化传播。
- 不满足统一工程命名规则。
- 与最终双名制命名体系冲突。

## 4. Complete Code

以下内容可作为最终命名规则摘要直接复用：

```md
# Naming Convention

## 最终品牌命名

- 英文总品牌：`Stellar Axis`
- 中文总品牌：`两仪`
- 标准首次展示：`Stellar Axis（两仪）`

## 最终双栈命名

- Java：`Stellar Steel（两仪·刚）`
- Go：`Stellar Titan（两仪·柔）`

## 最终核心组件命名

| 英文正式名 | 中文正式名 | 中文简称 |
| :--- | :--- | :--- |
| `StarMap` | `乾仪` | `乾` |
| `Nebula` | `坤仪` | `坤` |
| `LightBeam` | `离鉴` | `离` |
| `Orbit` | `巽策` | `巽` |
| `Pulsar` | `艮闸` | `艮` |
| `Chronos` | `震策` | `震` |
| `Singularity` | `坎锁` | `坎` |
| `EventHorizon` | `兑门` | `兑` |
| `TaiJi Flow` | `太极流` | `中枢` |

## 最终 AI 子体系命名

| 英文正式名 | 中文正式名 |
| :--- | :--- |
| `Spirit Axis` | `灵觉层` |
| `ZhongFu Engine` | `中孚` |
| `DaChu Memory` | `大畜` |
| `Sui Agent` | `随位` |
| `Xian MCP` | `咸` |
| `Heng MCP` | `恒` |
| `Meng MCP` | `蒙` |

## 命名规则

1. 总品牌固定为 `Stellar Axis（两仪）`。
2. 核心组件统一使用“英文正式名 + 中文八卦名”双名制。
3. 仓库、包名、依赖坐标、CLI、模块统一只使用英文工程名。
4. 不再使用 `Qian-Directory`、`Kun-Profile`、`Dui-Portal` 等中英混拼命名。

## 工程命名示例

- 仓库：`starmap`
- 仓库：`event-horizon`
- Maven：`io.stellar.axis:starmap-client-spring-boot-starter`
- Maven：`io.stellar.axis:stellar-steel-core`
- Go Module：`github.com/stellar-axis/taiji-flow`
- CLI：`stellarctl`
```
