---
title: Stellguard · 星盾 · 架构组成
outline: deep
---

# 架构组成

> 零信任安全平台，负责身份验证、访问控制、策略评估和服务间安全通信。

[返回产品首页](/zh/products/stellguard/)

## 核心组件

- Identity Provider：管理身份、令牌和证书
- Policy Engine：计算授权策略和风险决策
- Enforcement Point：在网关和服务侧执行访问控制

## 关键交互

- 身份提供方签发主体凭证并与策略引擎共享上下文
- 执行点在请求入口实时拦截并依据策略给出决策
