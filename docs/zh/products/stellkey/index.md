---
title: Stellkey 详细设计
outline: deep
---

# Stellkey · 星钥

<div class="product-logo">
  <img src="/logo/stellkey.png" alt="Stellkey Logo">
</div>

> 密钥中心，负责密钥生命周期管理、机密分发、轮换审计与安全访问控制。

## 产品定位

### 目标定位

Stellkey 为应用和平台组件提供统一密钥托管与机密分发能力，降低明文密钥散落、人工轮换困难和审计缺失带来的安全风险。

### 适用边界

- 面向密钥、证书、令牌和机密配置的统一托管
- 面向访问审批、轮换审计和短期凭证分发
- 不直接替代配置中心，而是为敏感值治理提供安全底座

## 核心能力

### 能力清单

- 支持密钥、证书、令牌和机密配置统一托管
- 支持版本轮换、临时凭证、访问审批和审计追踪
- 支持与配置中心、零信平台和网关联动分发

### 设计价值

- 统一机密生命周期管理，降低明文泄露风险
- 让敏感配置分发和访问审计具备平台级能力

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellkey/summary-design)
- [架构组成](/zh/products/stellkey/architecture)
- [部署形态](/zh/products/stellkey/deployment)
- [快速入门](/zh/products/stellkey/quick-start)
- [配置建议](/zh/products/stellkey/configuration)
- [API 与 SDK](/zh/products/stellkey/api-and-sdk)
- [可观测性](/zh/products/stellkey/observability)

## 典型场景

### 业务场景

- 数据库凭证托管
- TLS 证书管理

### 平台场景

- 配置中心敏感值引用
- 平台级机密生命周期治理
