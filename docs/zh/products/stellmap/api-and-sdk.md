---
title: Stellmap · 星图 · API 与 SDK
outline: deep
---

# API 与 SDK

> 注册中心，负责服务实例注册、健康检查、发现订阅与拓扑变更推送。

[返回产品首页](/zh/products/stellmap/)

## 读写流程

### 写请求

1. 客户端通过 `HTTP API` 发起注册、注销或心跳续约请求
2. 非 Leader 节点重定向或转发到 Leader
3. Leader 将实例变更命令编码为 proposal，提交给单 Raft 组
4. 日志追加到本地 `WAL`
5. 日志复制到多数派并提交
6. 状态机按顺序 apply 到 `Pebble`
7. 返回写成功结果

### 读请求

1. 客户端发起实例查询请求
2. 请求到达 Leader，或被 Follower 转发到 Leader
3. Leader 执行 `ReadIndex`
4. 等待本地 `appliedIndex` 追平到 `readIndex`
5. 从 `Pebble` 读取最新实例注册表视图
6. 返回线性一致结果

## 通信设计

### 外部 API

对三方业务接入只提供公共 `HTTP API`。

对外 `HTTP API` 的价值：

- 降低接入门槛
- 方便脚本、Sidecar 和多语言客户端调用
- 适合注册、注销、查询、健康上报等开放接口

### 管理面 API

成员变更、Leader 转移和集群状态查询不挂在公共 HTTP listener 上，而是挂在独立的 `admin HTTP` listener 上。

设计约束：

- 公共 `HTTP` 只承载业务数据面和健康检查
- `admin HTTP` 只承载控制面动作，默认绑定到本地回环地址，例如 `127.0.0.1:18080`
- `admin HTTP` 当前额外强制只接受来源为 `127.0.0.1` 的请求
- `admin HTTP` 需要固定 token 鉴权，请求头格式为 `Authorization: Bearer <token>`
- `stellmapctl` 是控制面的唯一入口
- 这意味着当前控制面默认只支持本机运维；如需跨机器运维，需要后续单独放宽来源限制并补更完整鉴权

### 内部 Transport

集群内部复制面统一采用 `gRPC`。

这样设计的原因：

- 内外通信面职责不同，不应共用同一套协议假设
- 单 Raft 集群内部需要承载 Raft message、快照、Leader 转发等协议
- `gRPC` 对内部机器间通信更适合，协议清晰、强类型、支持流式快照传输

内部 transport 主要承载：

- Raft message 传输
- Snapshot 文件/分片传输
- Leader 转发读写请求
- 节点管理与健康探测

策略：

- 外部开放面只提供 `HTTP API`
- 内部复制面统一使用 `gRPC`
- 外部 HTTP 与内部 `gRPC` 共享同一套核心状态机与权限校验逻辑，避免双份实现
