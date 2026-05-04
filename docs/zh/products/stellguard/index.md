---
title: Stellguard 详细设计
outline: deep
---

# Stellguard · 星盾

<div class="product-logo">
  <img src="/logo/stellguard.png" alt="Stellguard Logo">
</div>

> 零信任安全平台，负责身份验证、访问控制、策略评估和服务间安全通信。

## 产品定位

### 目标定位

Stellguard 以零信任理念统一管理用户、服务和终端访问，确保“默认不信任、持续验证、最小权限”在平台内落地。

### 适用边界

- 面向用户、服务和终端的统一身份与访问控制
- 面向服务间 mTLS、上下文授权和风险拦截场景
- 不替代业务权限模型，而是提供统一安全底座

## 核心能力

### 能力清单

- 支持用户、服务和设备身份认证
- 支持细粒度策略评估、动态授权和风险拦截
- 支持服务间 mTLS、证书轮换和访问审计

### 设计价值

- 将身份、授权和通信安全统一到同一控制面
- 让最小权限和持续验证可以平台化实施

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellguard/summary-design)
- [架构组成](/zh/products/stellguard/architecture)
- [部署形态](/zh/products/stellguard/deployment)
- [快速入门](/zh/products/stellguard/quick-start)
- [配置建议](/zh/products/stellguard/configuration)
- [API 与 SDK](/zh/products/stellguard/api-and-sdk)
- [可观测性](/zh/products/stellguard/observability)

## 典型场景

### 业务场景

- 服务间零信任访问
- 高敏操作审计

### 平台场景

- 控制台统一认证
- 平台级持续验证治理
