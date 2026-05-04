---
title: Stellflow 详细设计
outline: deep
---

# Stellflow · 彗流

<div class="product-logo">
  <img src="/logo/stellflow.png" alt="Stellflow Logo">
</div>

> 消息队列与事件流平台，负责异步解耦、流式分发、削峰填谷和事件驱动集成。

## 产品定位

### 目标定位

Stellflow 提供统一的消息投递与事件分发能力，支持业务异步化、系统集成和流式处理场景。

### 适用边界

- 面向异步解耦、事件驱动和数据分发场景
- 面向延迟消息、顺序消息和流式消费模型
- 不替代任务调度平台，而是聚焦消息持久化与分发

## 核心能力

### 能力清单

- 支持普通队列、主题订阅、延迟消息和顺序消息
- 支持消费组、重试队列、死信队列和回溯消费
- 支持事件桥接、流式处理和跨 Region 复制

### 设计价值

- 通过异步化降低系统耦合和峰值冲击
- 通过可靠投递与位点管理支撑事件驱动架构

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellflow/summary-design)
- [架构组成](/zh/products/stellflow/architecture)
- [部署形态](/zh/products/stellflow/deployment)
- [快速入门](/zh/products/stellflow/quick-start)
- [配置建议](/zh/products/stellflow/configuration)
- [API 与 SDK](/zh/products/stellflow/api-and-sdk)
- [可观测性](/zh/products/stellflow/observability)

## 典型场景

### 业务场景

- 订单事件驱动
- 异步削峰

### 平台场景

- 数据同步总线
- 跨系统事件桥接
