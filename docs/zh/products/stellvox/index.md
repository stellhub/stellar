---
title: Stellvox 详细设计
outline: deep
---

# Stellvox · 星讯

<div class="product-logo">
  <img src="/logo/stellvox.png" alt="Stellvox Logo">
</div>

> 告警平台，负责告警规则、收敛降噪、通知编排和事件协同处置。

## 产品定位

### 目标定位

Stellvox 将指标、日志、链路和平台事件统一接入告警中心，实现统一触发、统一收敛和统一协同。

### 适用边界

- 面向指标、日志、链路和平台事件的统一告警治理
- 面向通知编排、值班协同和事件闭环处置
- 不替代监控数据源，而是作为告警决策与通知中枢

## 核心能力

### 能力清单

- 支持阈值告警、异常检测、事件告警和组合规则
- 支持告警分组、抑制、静默和升级策略
- 支持多渠道通知、值班编排和事件回溯

### 设计价值

- 将多来源异常统一收敛为可治理事件
- 降低告警风暴对值班和排障效率的冲击

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellvox/summary-design)
- [架构组成](/zh/products/stellvox/architecture)
- [部署形态](/zh/products/stellvox/deployment)
- [快速入门](/zh/products/stellvox/quick-start)
- [配置建议](/zh/products/stellvox/configuration)
- [API 与 SDK](/zh/products/stellvox/api-and-sdk)
- [可观测性](/zh/products/stellvox/observability)

## 典型场景

### 业务场景

- 故障协同处置
- 值班轮值

### 平台场景

- 平台统一告警中心
- 多源事件收敛治理
