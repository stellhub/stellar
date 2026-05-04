---
title: Stellmap 详细设计
outline: deep
---

# Stellmap · 星图

<div class="product-logo">
  <img src="/logo/stellmap.png" alt="Stellmap Logo">
</div>

> 注册中心，负责服务实例注册、健康检查、发现订阅与拓扑变更推送。

## 产品定位

### 目标定位

Stellmap 面向微服务与基础设施组件提供统一注册发现能力，解决服务地址漂移、环境切换、实例伸缩和多集群路由等问题。

### 适用边界

- 面向服务注册、实例发现和订阅推送场景
- 面向多环境、多集群和灰度路由的服务目录管理
- 不承载业务配置和应用层流量治理逻辑

## 核心能力

### 能力清单

- 支持服务注册、实例摘除、心跳续约与健康检查
- 支持命名空间、集群、标签和版本维度的实例隔离
- 支持订阅推送与本地缓存，降低查询延迟

### 设计价值

- 通过注册发现统一服务入口，降低人工维护地址成本
- 通过标签和版本隔离支撑灰度发布与跨地域路由

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellmap/summary-design)
- [架构组成](/zh/products/stellmap/architecture)
- [部署形态](/zh/products/stellmap/deployment)
- [快速入门](/zh/products/stellmap/quick-start)
- [配置建议](/zh/products/stellmap/configuration)
- [API 与 SDK](/zh/products/stellmap/api-and-sdk)
- [可观测性](/zh/products/stellmap/observability)

## 典型场景

### 业务场景

- 微服务实例发现
- 多集群灰度路由

### 平台场景

- 基础设施节点目录管理
- 统一服务目录治理
