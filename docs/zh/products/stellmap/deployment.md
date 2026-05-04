---
title: Stellmap · 星图 · 部署形态
outline: deep
---

# 部署形态

> 注册中心，负责服务实例注册、健康检查、发现订阅与拓扑变更推送。

[返回产品首页](/zh/products/stellmap/)

## 推荐拓扑

- 单 Region 三节点起步，跨可用区部署以降低单点故障
- 可在多 Region 模式下通过只读从集群承载本地发现流量

## 高可用策略

- 注册数据建议按租约和副本双重机制做故障恢复
- 客户端需启用本地缓存与短暂失联容忍策略

## 崩溃恢复流程

`StellMap` 的恢复顺序需要严格围绕 `Snapshot -> Pebble -> WAL Replay` 展开。

### 启动恢复步骤

1. 读取本地最新快照元信息
2. 如果存在有效快照，则先恢复快照文件到状态机工作目录
3. 打开 `Pebble`，加载实例注册表数据和本地元数据
4. 打开 `WAL`，读取 `HardState`、`Entry` 和快照点位
5. 根据快照 index 丢弃已被快照覆盖的旧日志
6. 将剩余未 apply 的 committed entries 依序 apply 到状态机
7. 重建内存中的 Raft 节点、apply watermarks、membership 视图
8. 对外进入可服务状态

### 崩溃点处理原则

- 如果 WAL 已持久化但状态机未来得及 apply，重启后按日志重放
- 如果快照文件已生成但元信息未原子切换，则仍以旧快照为准
- 如果 Pebble 中保存的 applied index 落后于 WAL committed index，则继续补 apply
- 如果发现快照、WAL、状态机三者 index 不一致，优先保证“不回退已提交日志”的原则

### 恢复后的校验项

- `HardState.Commit >= appliedIndex`
- 状态机记录的 `appliedIndex` 不超过 WAL 中的最大 committed index
- 快照中的 `ConfState` 与恢复后内存 membership 一致
- 当前节点是否仍在最新 membership 中
- 是否需要继续从 Leader 拉取缺失日志或安装新快照

## 快照与日志压缩

为避免 Raft Log 无限增长，需要周期性生成快照并截断旧日志。

基本策略：

- 当 `appliedIndex - snapshotIndex` 超过阈值时触发快照
- 快照生成后，原子写入快照元数据
- WAL 仅保留快照点之后的必要日志
- 新节点或落后节点优先走日志追平，差距过大时改走快照安装
