# Module Dependency

## 1. Problem Analysis

在 [Proposal.md](./Proposal.md)、[Naming-Convention.md](./Naming-Convention.md) 和 [Repo-Layout.md](./Repo-Layout.md) 中，我们已经分别定义了：

- 体系命名
- 命名规范
- 仓库与目录布局

但如果没有进一步定义“谁可以依赖谁，谁不能依赖谁”，那么即使仓库拆分正确，代码结构仍然会很快失控。

常见失控方式包括：

- 聚合仓库反向依赖某个具体产品仓库，导致上层被下层污染。
- Starter 直接依赖服务端实现，导致客户端接入层和服务端耦合。
- 控制平面直接依赖某个产品的内部模块，而不是依赖公开 API。
- Operator 依赖业务 SDK 或控制台前端代码，造成部署层与业务层缠绕。
- 不同产品之间直接引用彼此内部实现，而不是通过协议或 SDK 通信。

因此，这份文档的目标是定义 `Stellar Axis（两仪）` 体系下的仓库级与模块级依赖规则。

## 2. Design

### 2.1 依赖设计原则

所有仓库和模块统一遵循以下原则：

- **单向依赖**：上层可以依赖下层抽象，下层不能反向依赖上层。
- **稳定依赖**：依赖必须尽量指向稳定接口、协议和 SDK，而不是指向内部实现。
- **同层隔离**：同一层的核心产品之间禁止直接依赖内部代码。
- **跨产品走协议**：跨产品交互统一通过 API、SDK、事件、配置或控制面协议完成。
- **实现不外泄**：`internal`、`domain`、`infrastructure` 等内部实现禁止被跨仓库直接依赖。

### 2.2 逻辑层级划分

建议将整个体系抽象为六层：

#### L0：规范层

- `stellar-axis`
- 各类规范文档

职责：

- 定义命名、布局、架构与治理规则

#### L1：协议与模型层

- `api/`
- `protobuf/`
- `openapi/`
- 通用模型定义

职责：

- 定义稳定契约

#### L2：公共运行时与 SDK 层

- `stellar-steel`
- `stellar-titan`
- 各产品 `clients/`
- 各产品 `sdk/`

职责：

- 提供接入能力、客户端封装和运行时抽象

#### L3：产品实现层

- 各产品 `server/`
- 各产品核心实现模块

职责：

- 实现具体中间件能力

#### L4：部署与控制层

- `stellar-control-plane`
- `stellarctl`
- `operators/`
- `deploy/`

职责：

- 运维、发布、控制、生命周期管理

#### L5：示例与测试层

- `stellar-examples`
- `stellar-benchmarks`
- `examples/`
- `test/`

职责：

- 示例、验证、压测、基准对比

### 2.3 总体依赖方向

统一依赖方向如下：

```text
L0 规范层
  ↓
L1 协议与模型层
  ↓
L2 公共运行时与 SDK 层
  ↓
L3 产品实现层
  ↓
L4 部署与控制层
  ↓
L5 示例与测试层
```

这里的含义不是“上层一定依赖下层全部内容”，而是“只允许朝这个方向建立依赖，不允许反向依赖”。

## 3. Implementation

### 3.1 仓库级依赖规则

#### `stellar-axis`

允许依赖：

- 无代码级依赖

禁止依赖：

- 所有具体产品实现仓库
- 所有 SDK、Starter、Sidecar、Operator 模块

解释：

- `stellar-axis` 是规范与入口仓库，不是运行时代码仓库。

#### `stellar-steel`

允许依赖：

- 各产品公开协议定义
- 各产品公开客户端接口
- 通用第三方基础库

禁止依赖：

- 任意产品仓库的 `server/`
- 任意产品仓库的 `internal/`
- `stellar-control-plane`
- `frontend/`

解释：

- `stellar-steel` 是 Java 聚合运行时，只能依赖稳定契约，不能反向依赖具体服务端实现。

#### `stellar-titan`

允许依赖：

- 各产品公开协议定义
- 各产品公开客户端接口
- 通用第三方基础库

禁止依赖：

- 任意产品仓库的 `server/`
- 任意产品仓库的 `internal/`
- `stellar-control-plane`

解释：

- `stellar-titan` 与 `stellar-steel` 同理，是 Go 运行时聚合，不允许被具体产品实现污染。

#### `spirit-axis`

允许依赖：

- LLM、Memory、Agent 所需的公开协议与 SDK
- MCP 协议定义
- 必要的通用基础库

禁止依赖：

- 核心产品仓库内部实现
- 控制平面前端
- 具体产品 Sidecar/Operator 内部逻辑

解释：

- AI 灵觉层可以消费平台能力，但应通过公开协议接入，而不是侵入具体产品实现。

#### 核心产品仓库

以 `starmap`、`nebula`、`orbit`、`taiji-flow` 等为代表。

允许依赖：

- 自身仓库内的 `api/`
- `stellar-steel` 或 `stellar-titan` 中的公共运行时抽象
- 必要的第三方基础组件

禁止依赖：

- 其他核心产品仓库的 `server/`
- 其他核心产品仓库的 `internal/`
- `stellar-control-plane` 的 `backend/`、`frontend/`
- `stellar-examples`

解释：

- 核心产品之间只能通过协议、SDK、事件或控制面 API 交互，不能直接吃彼此内部代码。

#### `stellar-control-plane`

允许依赖：

- 各产品公开 API
- 各产品公开 SDK
- 各产品 Operator/管理接口的稳定契约

禁止依赖：

- 各产品 `internal/`
- 各产品 `domain/`
- 各产品 `infrastructure/`
- `starters/`
- `sidecars/`

解释：

- 控制平面应该面向“产品公开管理面”编程，而不是直接侵入产品内部实现。

#### `stellarctl`

允许依赖：

- 各产品公开管理 API Client
- 通用认证、配置、输出模块

禁止依赖：

- 产品服务端内部实现
- 控制台前端代码
- Operator 控制器实现

解释：

- CLI 是控制入口，不是服务端内部工具集。

#### `stellar-examples`

允许依赖：

- 各产品公开 SDK
- `stellar-steel`
- `stellar-titan`

禁止依赖：

- 任意产品内部模块
- 任意控制平面内部模块

解释：

- 示例只能展示标准接入方式，不能通过内部依赖“作弊”。

### 3.2 单产品仓库内部依赖规则

以统一产品骨架为基准：

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

#### `api/`

允许依赖：

- 通用模型库
- 协议定义库

禁止依赖：

- `server/`
- `clients/`
- `starters/`
- `sidecars/`
- `operators/`

解释：

- `api/` 必须是最稳定、最底层的契约层。

#### `server/`

允许依赖：

- `api/`
- 必要的公共运行时
- 必要的第三方基础库

禁止依赖：

- `starters/`
- `examples/`
- 其他产品的 `server/`

解释：

- `server/` 是产品主实现，不能反向依赖接入层和示例层。

#### `clients/`

允许依赖：

- `api/`
- 公共运行时

禁止依赖：

- `server/`
- `sidecars/`
- `operators/`

解释：

- Client 必须独立于服务端主实现。

#### `starters/`

允许依赖：

- `clients/`
- `api/`
- `stellar-steel` 公共模块

禁止依赖：

- `server/`
- `sidecars/`
- `operators/`
- 控制平面模块

解释：

- Starter 只负责接入装配，不允许把服务端主实现塞进去。

#### `sidecars/`

允许依赖：

- `api/`
- `clients/`
- `stellar-titan` 或边车运行时公共模块

禁止依赖：

- `starters/`
- 控制平面前端
- 其他产品的 `server/`

解释：

- Sidecar 是独立进程形态，只依赖协议和运行时，不依赖服务端内部实现。

#### `operators/`

允许依赖：

- `api/`
- Kubernetes Operator 运行时
- 必要的管理 API Client

禁止依赖：

- `server/`
- `starters/`
- 前端模块
- 业务示例代码

解释：

- Operator 负责生命周期管理，不负责承载业务主逻辑。

#### `deploy/`

允许依赖：

- 镜像名
- 配置模板
- 发布产物清单

禁止依赖：

- 代码模块级直接引用

解释：

- `deploy/` 是部署资产，不是代码依赖层。

#### `test/`

允许依赖：

- `api/`
- `server/`
- `clients/`
- 测试专用工具模块

禁止依赖：

- 作为生产依赖被其他模块引用

解释：

- 测试代码可以依赖生产代码，反向不成立。

### 3.3 聚合仓库内部依赖规则

#### `stellar-steel`

推荐依赖方向：

```text
stellar-steel-bom
  ↓
stellar-steel-common
  ↓
stellar-steel-core
  ↓
sdks/*
  ↓
starters/*
  ↓
examples/*
```

规则：

- `bom` 不依赖业务模块。
- `common` 不依赖 `sdk` 和 `starter`。
- `core` 可以依赖 `common`，不能依赖具体 `starter`。
- `starter` 可以依赖 `sdk` 与 `core`，不能依赖任意产品 `server/`。

#### `stellar-titan`

推荐依赖方向：

```text
pkg/runtime
  ↓
pkg/*
  ↓
sdk/*
  ↓
sidecars/*
  ↓
examples/*
```

规则：

- `pkg/runtime` 和 `pkg/*` 是稳定公共层。
- `sdk/*` 可以依赖 `pkg/*`。
- `sidecars/*` 可以依赖 `sdk/*` 和 `pkg/*`。
- `sidecars/*` 不得反向要求 `sdk/*` 依赖它们。

### 3.4 允许的跨产品交互方式

核心产品之间允许以下交互方式：

- 通过 HTTP/gRPC/OpenAPI/Protobuf 调用公开接口
- 通过官方 SDK 调用公开能力
- 通过 `TaiJi Flow` 进行异步事件通信
- 通过 `Nebula` 下发配置
- 通过控制平面进行统一编排

核心产品之间禁止以下交互方式：

- 直接引用其他产品的 `internal` 代码
- 直接引用其他产品的 `server/domain/infrastructure`
- 直接共享数据库表结构作为隐式接口
- 通过复制代码方式复用能力

### 3.5 特殊禁止规则

以下依赖关系一律禁止：

- `starter -> server`
- `operator -> server`
- `operator -> frontend`
- `client -> server`
- `productA/server -> productB/server`
- `control-plane -> product/internal`
- `examples -> product/internal`
- `stellar-axis -> any-runtime-code`

原因：

- 这些依赖会打破层级、破坏封装，并导致后续升级和拆仓成本急剧上升。

### 3.6 最终依赖矩阵

#### 允许依赖

| From | Allow |
| :--- | :--- |
| `stellar-axis` | 文档级引用，无代码依赖 |
| `stellar-steel` | `api`、公开 SDK、公共模型 |
| `stellar-titan` | `api`、公开 SDK、公共模型 |
| `spirit-axis` | 公开协议、公开 SDK、MCP 协议 |
| `product/server` | 本产品 `api`、公共运行时 |
| `product/clients` | 本产品 `api`、公共运行时 |
| `product/starters` | 本产品 `clients`、本产品 `api`、`stellar-steel` |
| `product/sidecars` | 本产品 `api`、本产品 `clients`、`stellar-titan` |
| `product/operators` | 本产品 `api`、管理 API Client、K8s Runtime |
| `stellar-control-plane` | 各产品公开 API、公开 SDK |
| `stellarctl` | 各产品公开管理 API Client |
| `stellar-examples` | 公开 SDK、公共运行时 |

#### 禁止依赖

| From | Forbidden |
| :--- | :--- |
| `stellar-axis` | 任意运行时代码 |
| `stellar-steel` | 任意产品 `server/internal` |
| `stellar-titan` | 任意产品 `server/internal` |
| `product/starters` | `server` |
| `product/clients` | `server` |
| `product/operators` | `server`、`frontend` |
| `product/sidecars` | 其他产品 `server/internal` |
| `stellar-control-plane` | 任意产品 `internal/domain/infrastructure` |
| `stellar-examples` | 任意产品内部实现 |

## 4. Complete Code

以下内容可作为依赖规则摘要直接复用：

```md
# Module Dependency

## 依赖原则

1. 统一单向依赖，禁止反向依赖。
2. 统一依赖公开协议、SDK 和稳定接口，禁止依赖内部实现。
3. 核心产品之间通过 API、SDK、事件或控制面交互，禁止直接引用彼此内部代码。

## 仓库级规则

- `stellar-axis`：只做规范，不依赖任何运行时代码
- `stellar-steel`：可依赖公开 API 与 SDK，不可依赖产品服务端实现
- `stellar-titan`：可依赖公开 API 与 SDK，不可依赖产品服务端实现
- `stellar-control-plane`：可依赖公开管理接口，不可依赖产品内部模块
- `stellarctl`：可依赖公开管理 API Client，不可依赖服务端内部代码

## 模块级规则

- `api ->` 不能依赖 `server/clients/starters/sidecars/operators`
- `server ->` 可依赖 `api`，不可依赖 `starter/examples`
- `clients ->` 可依赖 `api`，不可依赖 `server`
- `starters ->` 可依赖 `clients/api`，不可依赖 `server`
- `sidecars ->` 可依赖 `clients/api/runtime`，不可依赖其他产品内部实现
- `operators ->` 可依赖 `api` 与 K8s Runtime，不可依赖 `server/frontend`

## 明确禁止

- `starter -> server`
- `client -> server`
- `operator -> server`
- `control-plane -> product/internal`
- `productA/server -> productB/server`
- `examples -> product/internal`
```
