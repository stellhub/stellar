---
title: Stelltrace 详细设计
outline: deep
---

# Stelltrace · 星迹

<div class="product-logo">
  <img src="/logo/stelltrace.png" alt="Stelltrace Logo">
</div>

> 链路追踪平台，负责全链路 Trace、Span 采集、关联分析与问题定位。

## 产品定位

### 目标定位

Stelltrace 面向服务调用链、异步任务链和跨网关请求链提供统一追踪能力，帮助研发和运维快速识别慢请求、异常传播和跨系统依赖问题。

### 适用边界

- 面向请求链路、异步消息链路和定时任务链路追踪
- 面向问题定位、拓扑分析和跨系统依赖排查
- 不替代日志平台和指标平台，而是作为关联主线

## 核心能力

### 能力清单

- 支持 HTTP、gRPC、MQ、定时任务等链路追踪
- 支持采样策略、异常聚类和拓扑依赖视图
- 支持日志、指标和告警的 Trace 关联查询

### 设计价值

- 以 Trace ID 为中心串联多类观测数据
- 提升慢请求和异常传播的定位效率

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stelltrace/summary-design)
- [架构组成](/zh/products/stelltrace/architecture)
- [部署形态](/zh/products/stelltrace/deployment)
- [快速入门](/zh/products/stelltrace/quick-start)
- [配置建议](/zh/products/stelltrace/configuration)
- [API 与 SDK](/zh/products/stelltrace/api-and-sdk)
- [可观测性](/zh/products/stelltrace/observability)

## 典型场景

### 业务场景

- 慢请求排查
- 异常根因定位

### 平台场景

- 服务依赖拓扑治理
- 全链路观测关联分析
