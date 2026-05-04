---
title: Stellgate · 视界 · 架构组成
outline: deep
---

# 架构组成

> 网关入口平台，负责统一接入、认证鉴权、协议转换、流量治理与 API 暴露。

[返回产品首页](/zh/products/stellgate/)

## 核心组件

- Gateway Admin：管理路由、证书和插件配置
- Gateway Runtime：处理转发、负载和流量管控
- Plugin SDK：承载鉴权、审计和协议扩展逻辑

## 关键交互

- 控制面将路由和插件策略发布到运行时网关节点
- 请求经过插件链处理后再转发至目标上游服务
