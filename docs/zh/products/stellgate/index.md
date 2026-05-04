---
title: Stellgate 详细设计
outline: deep
---

# Stellgate · 视界

<div class="product-logo">
  <img src="/logo/stellgate.png" alt="Stellgate Logo">
</div>

> 网关入口平台，负责统一接入、认证鉴权、协议转换、流量治理与 API 暴露。

## 产品定位

### 目标定位

Stellgate 作为外部入口和内部流量编排枢纽，负责承接客户端请求、统一安全策略并将流量转发至后端服务。

### 适用边界

- 面向外部 API、内部服务入口和边缘接入场景
- 面向协议接入、认证鉴权、转发治理和插件扩展
- 不替代服务治理平台，而是承接入口侧流量处理

## 核心能力

### 能力清单

- 支持 HTTP、gRPC、WebSocket 和 SSE 接入
- 支持鉴权、限流、路由、改写和插件化扩展
- 支持多租户、多域名和多环境网关编排

### 设计价值

- 将入口治理、安全和协议处理收口到统一网关
- 提升对外 API 暴露和内部流量编排的一致性

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellgate/summary-design)
- [架构组成](/zh/products/stellgate/architecture)
- [部署形态](/zh/products/stellgate/deployment)
- [快速入门](/zh/products/stellgate/quick-start)
- [配置建议](/zh/products/stellgate/configuration)
- [API 与 SDK](/zh/products/stellgate/api-and-sdk)
- [可观测性](/zh/products/stellgate/observability)

## 典型场景

### 业务场景

- Open API 统一入口
- LLM 网关与边缘接入

### 平台场景

- 内部服务网关
- 边缘流量接入
