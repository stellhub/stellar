---
title: Stellkey · 星钥 · 架构组成
outline: deep
---

# 架构组成

> 密钥中心，负责密钥生命周期管理、机密分发、轮换审计与安全访问控制。

[返回产品首页](/zh/products/stellkey/)

## 核心组件

- Secret Store：保存机密、证书和版本记录
- Rotation Engine：执行轮换、吊销和过期提醒
- Access Gateway：统一鉴权、审计和分发控制

## 关键交互

- 应用通过 Access Gateway 获取短期凭证或机密引用
- Rotation Engine 按策略触发轮换并同步更新版本状态
