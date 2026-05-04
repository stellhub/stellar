---
title: Stellnula 详细设计
outline: deep
---

# Stellnula · 星云

<div class="product-logo">
  <img src="/logo/stellnula.png" alt="Stellnula Logo">
</div>

> 配置中心，负责配置集中存储、版本管理、灰度发布与动态下发。

## 产品定位

### 目标定位

Stellnula 为服务、网关和平台组件提供统一配置治理能力，帮助团队消除分散配置、环境漂移和人工变更不可审计的问题。

### 适用边界

- 面向应用参数、平台开关和运行时配置分发
- 面向配置版本管理、灰度发布和审计回溯场景
- 敏感配置推荐与密钥中心联动，而非直接明文托管

## 核心能力

### 能力清单

- 支持命名空间、应用、环境和标签维度的配置隔离
- 支持版本历史、灰度发布、回滚与变更审计
- 支持推送监听、长轮询和配置快照回退

### 设计价值

- 降低应用发布和参数调整之间的耦合度
- 将配置治理纳入统一变更流程，减少误操作风险

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellnula/summary-design)
- [架构组成](/zh/products/stellnula/architecture)
- [部署形态](/zh/products/stellnula/deployment)
- [快速入门](/zh/products/stellnula/quick-start)
- [配置建议](/zh/products/stellnula/configuration)
- [API 与 SDK](/zh/products/stellnula/api-and-sdk)
- [可观测性](/zh/products/stellnula/observability)

## 典型场景

### 业务场景

- 微服务动态配置下发
- 灰度参数实验

### 平台场景

- 多环境配置治理
- 平台统一开关管理
