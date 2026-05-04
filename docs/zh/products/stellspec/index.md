---
title: Stellspec 详细设计
outline: deep
---

# Stellspec · 星谱

<div class="product-logo">
  <img src="/logo/stellspec.png" alt="Stellspec Logo">
</div>

> 日志平台，负责统一采集、结构化处理、检索分析与日志留存治理。

## 产品定位

### 目标定位

Stellspec 为业务日志、审计日志和平台日志提供统一采集与检索能力，帮助团队建立标准化日志规范和快速排障路径。

### 适用边界

- 面向应用日志、审计日志和平台运行日志统一治理
- 面向日志检索、追踪关联和留存策略管理
- 不替代指标与链路追踪，而是承接全文检索和上下文还原

## 核心能力

### 能力清单

- 支持文本日志、结构化日志和审计日志统一采集
- 支持字段解析、脱敏、索引路由和生命周期管理
- 支持 Trace 关联、日志聚类和异常模式检索

### 设计价值

- 建立统一日志规范，降低检索和定位成本
- 通过脱敏和生命周期管理兼顾安全与成本

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellspec/summary-design)
- [架构组成](/zh/products/stellspec/architecture)
- [部署形态](/zh/products/stellspec/deployment)
- [快速入门](/zh/products/stellspec/quick-start)
- [配置建议](/zh/products/stellspec/configuration)
- [API 与 SDK](/zh/products/stellspec/api-and-sdk)
- [可观测性](/zh/products/stellspec/observability)

## 典型场景

### 业务场景

- 应用日志检索
- 故障排查与根因定位

### 平台场景

- 审计追踪
- 平台级日志留存治理
