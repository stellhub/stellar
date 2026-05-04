---
title: Stellpulse 详细设计
outline: deep
---

# Stellpulse · 脉冲

<div class="product-logo">
  <img src="/logo/stellpulse.png" alt="Stellpulse Logo">
</div>

> 流控熔断平台，负责热点保护、容量守卫、隔离舱和自适应降级。

## 产品定位

### 目标定位

Stellpulse 聚焦高并发场景下的可用性保护，通过限流、并发隔离、熔断降级和系统负载守护保证关键业务服务稳定运行。

### 适用边界

- 面向服务入口、关键接口和依赖调用的稳定性保护
- 面向突发流量、热点资源和依赖抖动场景
- 不替代服务治理，而是聚焦准入控制和退化保护

## 核心能力

### 能力清单

- 支持 QPS 限流、并发数限制和热点参数保护
- 支持熔断、舱壁隔离、预热与排队等待
- 支持基于错误率、慢调用和系统负载的自动降级

### 设计价值

- 为核心链路提供统一的可用性保护底座
- 将突发流量和依赖异常对系统的冲击控制在可预期范围内

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellpulse/summary-design)
- [架构组成](/zh/products/stellpulse/architecture)
- [部署形态](/zh/products/stellpulse/deployment)
- [快速入门](/zh/products/stellpulse/quick-start)
- [配置建议](/zh/products/stellpulse/configuration)
- [API 与 SDK](/zh/products/stellpulse/api-and-sdk)
- [可观测性](/zh/products/stellpulse/observability)

## 典型场景

### 业务场景

- 秒杀流量保护
- 核心链路降级兜底

### 平台场景

- 依赖故障隔离
- 平台统一流量守卫
