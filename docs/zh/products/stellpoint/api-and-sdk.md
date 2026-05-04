---
title: Stellpoint · 奇点 · API 与 SDK
outline: deep
---

# API 与 SDK

> 分布式锁与协调中心，负责互斥控制、领导者选举和关键资源串行化访问。

[返回产品首页](/zh/products/stellpoint/)

## 接口分层

- 提供锁管理、选主监听和租约查询接口

## SDK 说明

- SDK 提供同步、异步和注解式加锁能力
- 推荐结合 fencing token 和业务版本号做双重保护
