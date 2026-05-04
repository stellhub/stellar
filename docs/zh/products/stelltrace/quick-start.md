---
title: Stelltrace · 星迹 · 快速入门
outline: deep
---

# 快速入门

> 链路追踪平台，负责全链路 Trace、Span 采集、关联分析与问题定位。

[返回产品首页](/zh/products/stelltrace/)

## 接入步骤

1. 在应用中启用 Trace SDK。
2. 配置采样策略和上报地址。
3. 通过 Trace ID 在平台中检索完整调用链。

## 首次验证

- 发起一次跨服务调用并确认链路已完整展示
- 故意制造慢调用，验证慢链路筛选是否生效
