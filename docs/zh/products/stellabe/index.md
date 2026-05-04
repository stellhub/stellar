---
title: Stellabe 详细设计
outline: deep
---

# Stellabe · 星盘

<div class="product-logo">
  <img src="/logo/stellabe.png" alt="Stellabe Logo">
</div>

> 任务调度平台，负责定时任务编排、依赖调度、分片执行与运行治理。

## 产品定位

### 目标定位

Stellabe 提供面向平台级和业务级任务的统一调度能力，支持批处理、定时执行、工作流编排和失败重试。

### 适用边界

- 面向定时任务、工作流和批处理调度场景
- 面向任务回溯、补数和执行过程可视化
- 不替代流式计算引擎，而是聚焦任务编排和调度控制

## 核心能力

### 能力清单

- 支持 Cron 调度、工作流 DAG、任务依赖和分片执行
- 支持失败重试、幂等控制、补数和回溯执行
- 支持租户隔离、任务审计和运行态可视化

### 设计价值

- 将分散脚本和手工任务纳入统一治理体系
- 提高任务编排的可维护性和执行可靠性

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellabe/summary-design)
- [架构组成](/zh/products/stellabe/architecture)
- [部署形态](/zh/products/stellabe/deployment)
- [快速入门](/zh/products/stellabe/quick-start)
- [配置建议](/zh/products/stellabe/configuration)
- [API 与 SDK](/zh/products/stellabe/api-and-sdk)
- [可观测性](/zh/products/stellabe/observability)

## 典型场景

### 业务场景

- 定时报表
- 数据补数

### 平台场景

- 发布后置任务编排
- 平台级工作流治理
