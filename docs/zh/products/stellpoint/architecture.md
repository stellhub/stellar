---
title: Stellpoint · 奇点 · 架构组成
outline: deep
---

# 架构组成

> 分布式锁与协调中心，负责互斥控制、领导者选举和关键资源串行化访问。

[返回产品首页](/zh/products/stellpoint/)

## 核心组件

- Lock API：受理加锁、解锁、续约和查询请求
- Lease Manager：管理会话、租约和过期回收
- Coordination Log：保存锁顺序与 fencing token

## 关键交互

- 客户端通过会话保持租约有效并周期续约
- 协调日志生成 fencing token 保障持有者顺序语义
