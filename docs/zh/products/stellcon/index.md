---
title: Stellcon 详细设计
outline: deep
---

# Stellcon · 星座

<div class="product-logo">
  <img src="/logo/stellcon.png" alt="Stellcon Logo">
</div>

> 指标平台，负责指标采集、聚合计算、查询分析和容量看板建设。

## 产品定位

### 目标定位

Stellcon 为应用、网关、中间件和基础设施提供统一指标体系，支撑容量评估、SLO 观测和告警阈值计算。

### 适用边界

- 面向系统运行指标、业务指标和平台资源指标治理
- 面向趋势分析、容量规划和 SLO 目标管理
- 不替代日志检索，而是聚焦时序指标采集与聚合查询

## 核心能力

### 能力清单

- 支持 Counter、Gauge、Histogram 和 Summary 指标模型
- 支持统一标签规范、聚合规则和多维查询
- 支持看板模板、SLO 计算和容量趋势分析

### 设计价值

- 形成统一指标口径，减少各系统观测语义分裂
- 为告警、容量和稳定性治理提供统一数据底座

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellcon/summary-design)
- [架构组成](/zh/products/stellcon/architecture)
- [部署形态](/zh/products/stellcon/deployment)
- [快速入门](/zh/products/stellcon/quick-start)
- [配置建议](/zh/products/stellcon/configuration)
- [API 与 SDK](/zh/products/stellcon/api-and-sdk)
- [可观测性](/zh/products/stellcon/observability)

## 典型场景

### 业务场景

- SLO 观测
- 容量规划

### 平台场景

- 统一监控看板
- 多租户指标治理
