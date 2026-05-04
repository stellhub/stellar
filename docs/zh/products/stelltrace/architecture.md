---
title: Stelltrace · 星迹 · 架构组成
outline: deep
---

# 架构组成

> 链路追踪平台，负责全链路 Trace、Span 采集、关联分析与问题定位。

[返回产品首页](/zh/products/stelltrace/)

## 核心组件

- Agent / SDK：埋点采集与上下文传播
- Collector：接收 Span 并进行预聚合
- Query API：提供检索、分析和可视化接口

## 关键交互

- Agent 将 Span 上报到 Collector 进行聚合与清洗
- 查询层按 Trace ID 还原链路并关联日志与指标
