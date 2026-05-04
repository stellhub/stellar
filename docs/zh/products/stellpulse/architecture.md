---
title: Stellpulse · 脉冲 · 架构组成
outline: deep
---

# 架构组成

> 流控熔断平台，负责热点保护、容量守卫、隔离舱和自适应降级。

[返回产品首页](/zh/products/stellpulse/)

## 核心组件

- Rule Center：统一管理限流与熔断规则
- Runtime Engine：进行请求准入和降级判定
- Metrics Reporter：上报命中、拒绝和恢复事件

## 关键交互

- 规则推送到运行时后在请求入口进行本地判断
- 拒绝、熔断和恢复事件回传到监控和告警体系
