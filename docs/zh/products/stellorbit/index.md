---
title: Stellorbit 详细设计
outline: deep
---

# Stellorbit · 星轨

<div class="product-logo">
  <img src="/logo/stellorbit.png" alt="Stellorbit Logo">
</div>

> 服务治理中枢，负责路由、负载、限流策略编排以及服务生命周期治理。

## 产品定位

### 目标定位

Stellorbit 为服务间调用提供统一治理策略，解决灰度发布、地域路由、版本隔离和故障切换等动态治理问题。

### 适用边界

- 面向服务调用链中的路由、重试、切流和故障隔离
- 面向多版本、多地域、多租户的治理规则编排
- 不直接承载注册和配置存储，而是消费外部契约信息

## 核心能力

### 能力清单

- 支持权重路由、标签路由、版本路由和地域路由
- 支持负载均衡、重试、超时和故障摘除编排
- 支持治理规则的灰度生效和回滚

### 设计价值

- 让流量治理从代码内嵌转为平台化编排
- 提升灰度发布和多活切流的可控性

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellorbit/summary-design)
- [架构组成](/zh/products/stellorbit/architecture)
- [部署形态](/zh/products/stellorbit/deployment)
- [快速入门](/zh/products/stellorbit/quick-start)
- [配置建议](/zh/products/stellorbit/configuration)
- [API 与 SDK](/zh/products/stellorbit/api-and-sdk)
- [可观测性](/zh/products/stellorbit/observability)

## 典型场景

### 业务场景

- 灰度发布
- 多活切流

### 平台场景

- 服务故障隔离
- 跨集群治理编排
