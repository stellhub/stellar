---
title: Stelltrace · 星迹 · 部署形态
outline: deep
---

# 部署形态

> 链路追踪平台，负责全链路 Trace、Span 采集、关联分析与问题定位。

[返回产品首页](/zh/products/stelltrace/)

## 推荐拓扑

- 建议采集层与存储层解耦，Collector 可水平扩容
- 冷热数据建议分层存储，控制长期留存成本

## 高可用策略

- Collector 集群需支持批量缓冲和故障转移
- 查询节点建议与索引节点分离，避免相互争抢资源
