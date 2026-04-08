# Roadmap

## 1. Problem Analysis

当前体系已经完成了命名、仓库布局、依赖边界和架构决策层面的收敛，但仍然缺少一份按阶段推进的建设路线图。

如果没有路线图，后续会出现三个问题：

- 仓库创建顺序失控，容易先建“外围”再补“底座”
- 团队会同时启动过多产品，导致每条线都不完整
- 规范虽然齐全，但缺少明确的阶段性里程碑，无法判断体系建设是否进入下一阶段

因此，这份文档的目标是定义 `Stellar Axis（两仪）` 的分阶段建设顺序，明确每一阶段的目标、交付物和进入下一阶段的前置条件。

## 2. Design

### 2.1 路线图设计原则

路线图遵循以下原则：

- **先底座，后产品**：先完成命名、布局、依赖、聚合仓库，再进入核心产品建设
- **先主链路，后扩展**：优先打通注册、配置、治理、消息四条主链路
- **先接入，后生态**：先完成 Java/Go 接入路径，再建设控制平面、Operator、AI 扩展
- **先最小可运行，后全矩阵**：每个阶段都要求形成可验证产物，而不是停留在抽象设计

### 2.2 阶段划分

建议将整体建设划分为五个阶段：

- Phase 0：规则冻结
- Phase 1：底座仓库初始化
- Phase 2：主链路核心产品落地
- Phase 3：控制面与运行形态扩展
- Phase 4：AI 灵觉层与生态完善

## 3. Implementation

### 3.1 Phase 0：规则冻结

#### 目标

冻结体系级规则，避免后续边做边改口径。

#### 当前状态

该阶段已基本完成。

#### 核心交付物

- [Proposal.md](./Proposal.md)
- [Naming-Convention.md](./Naming-Convention.md)
- [Repo-Layout.md](./Repo-Layout.md)
- [Module-Dependency.md](./Module-Dependency.md)
- [Architecture-Decision-Record.md](./Architecture-Decision-Record.md)
- [README.md](./README.md)
- [README-EN.md](./README-EN.md)

#### 完成标准

- 总品牌固定
- 核心产品命名固定
- 仓库布局固定
- 依赖规则固定
- 首页文档与规范文档一致

### 3.2 Phase 1：底座仓库初始化

#### 目标

创建整个体系的仓库骨架和聚合入口，形成最小组织结构。

#### 建议优先创建

- `stellar-axis`
- `stellar-steel`
- `stellar-titan`
- `stellarctl`
- `stellar-control-plane`
- `stellar-examples`
- `stellar-deploy`

#### 核心交付物

- 每个仓库具备基础目录结构
- 每个仓库具备 README、License、占位模块
- Java 聚合仓库具备 BOM、core、starter、sdk 骨架
- Go 聚合仓库具备 runtime、sdk、sidecar 骨架

#### 里程碑

- Java 和 Go 双栈聚合仓库可初始化构建
- CLI 仓库可成功编译一个空命令
- 控制平面仓库具备前后端占位结构

#### 进入下一阶段的条件

- 所有聚合仓库完成初始化
- 所有目录结构与 [Repo-Layout.md](./Repo-Layout.md) 对齐
- 依赖方向符合 [Module-Dependency.md](./Module-Dependency.md)

### 3.3 Phase 2：主链路核心产品落地

#### 目标

优先打通最关键的基础设施主链路，形成“服务发现 + 配置 + 治理 + 消息”的最小可运行体系。

#### 第一优先级产品

- `starmap`
- `nebula`
- `orbit`
- `taiji-flow`

#### 第二优先级产品

- `lightbeam`
- `pulsar`

#### 第三优先级产品

- `chronos`
- `singularity`
- `event-horizon`

#### 核心交付物

- 每个优先级产品至少完成：
  - `api/`
  - `server/`
  - `clients/`
  - `deploy/`
  - `test/`
- Java SDK 和 Go SDK 可接入至少一条主链路
- `starmap + nebula + orbit + taiji-flow` 能形成一套最小演示链路

#### 里程碑

- 服务注册与发现跑通
- 动态配置下发跑通
- 基础路由与治理规则跑通
- 消息或事件链路跑通

#### 进入下一阶段的条件

- 第一优先级四个产品完成最小可运行版本
- `stellar-examples` 中具备一套 Java demo 和一套 Go demo
- 控制平面可接入最少两个核心产品的管理接口

### 3.4 Phase 3：控制面与运行形态扩展

#### 目标

在核心主链路稳定后，扩展统一控制能力和云原生运行形态。

#### 建设重点

- `stellar-control-plane`
- `stellarctl`
- `starters/`
- `sidecars/`
- `operators/`

#### 重点产物

- Java Starter 首批上线：
  - `starmap-spring-boot-starter`
  - `nebula-spring-boot-starter`
  - `orbit-spring-boot-starter`
  - `taiji-flow-spring-boot-starter`
- Sidecar 首批上线：
  - `orbit-sidecar`
  - `pulsar-sidecar`
  - `lightbeam-sidecar`
- Operator 首批上线：
  - `starmap-operator`
  - `nebula-operator`
  - `taiji-flow-operator`

#### 里程碑

- 控制平面可以查看和操作核心产品
- CLI 可以完成基础查询、发布、诊断
- 至少一个产品完成 K8s Operator 管理闭环
- 至少一个治理能力以 Sidecar 方式运行

#### 进入下一阶段的条件

- 控制平面完成第一批产品接入
- Java 接入路径稳定
- Go Sidecar 路径稳定
- 部署仓库有可复用安装模板

### 3.5 Phase 4：AI 灵觉层与生态完善

#### 目标

在治理底座成熟后，建设 AI 灵觉层，并补齐性能、示例、文档和发布生态。

#### 建设重点

- `spirit-axis`
- `zhongfu-engine`
- `dachu-memory`
- `sui-agent`
- `xian-mcp`
- `heng-mcp`
- `meng-mcp`

#### 生态补齐项

- `stellar-benchmarks`
- 完整部署模板
- 完整示例工程
- 发布流水线
- 文档站点

#### 里程碑

- AI 灵觉层可通过 MCP 与基础设施交互
- 形成至少一个 AI 驱动的治理或运维场景
- 性能压测和基准对比形成公开结果

### 3.6 各阶段优先级总览

| 阶段 | 重点 | 核心目标 |
| :--- | :--- | :--- |
| Phase 0 | 规则冻结 | 确立统一规则与文档口径 |
| Phase 1 | 聚合仓库初始化 | 搭好组织与代码骨架 |
| Phase 2 | 主链路核心产品 | 打通最小可运行中间件体系 |
| Phase 3 | 控制面与运行形态 | 补齐接入、边车、Operator、控制台 |
| Phase 4 | AI 与生态 | 建设灵觉层与完整生态能力 |

### 3.7 推荐建设顺序

推荐实际执行顺序：

1. `stellar-axis`
2. `stellar-steel`
3. `stellar-titan`
4. `starmap`
5. `nebula`
6. `orbit`
7. `taiji-flow`
8. `stellar-examples`
9. `stellar-control-plane`
10. `stellarctl`
11. `lightbeam`
12. `pulsar`
13. `chronos`
14. `singularity`
15. `event-horizon`
16. `spirit-axis`

#### 顺序说明

- `starmap`、`nebula`、`orbit`、`taiji-flow` 先行，是因为它们决定服务发现、配置、治理和异步通信四条主链路。
- `lightbeam` 和 `pulsar` 紧随其后，是因为观测和保护会直接影响体系可用性。
- `event-horizon` 放后一些，是因为入口层更适合建立在内部治理稳定之后。
- AI 灵觉层最后进入，是因为它应该建立在稳固的治理底座之上。

## 4. Complete Code

以下内容可作为路线图摘要直接复用：

```md
# Roadmap

## Phase 0

- 冻结命名、布局、依赖和架构决策

## Phase 1

- 初始化 `stellar-axis`
- 初始化 `stellar-steel`
- 初始化 `stellar-titan`
- 初始化 `stellar-control-plane`
- 初始化 `stellarctl`

## Phase 2

- 优先落地 `starmap`
- 优先落地 `nebula`
- 优先落地 `orbit`
- 优先落地 `taiji-flow`

## Phase 3

- 补齐 Starter、Sidecar、Operator
- 建设控制平面和 CLI

## Phase 4

- 建设 `spirit-axis`
- 建设 AI 灵觉层
- 补齐 benchmarks、deploy、examples
```
