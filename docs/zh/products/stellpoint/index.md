---
title: Stellpoint 详细设计
outline: deep
---

# Stellpoint · 奇点

<div class="product-logo">
  <img src="/logo/stellpoint.png" alt="Stellpoint Logo">
</div>

> 分布式锁与协调中心，负责互斥控制、领导者选举和关键资源串行化访问。

## 产品定位

### 目标定位

Stellpoint 面向分布式系统中的共享资源竞争问题，提供统一锁服务、租约模型和协调原语，降低重复造轮子成本。

### 适用边界

- 面向互斥执行、领导者选举和顺序控制场景
- 面向共享资源写入保护和关键任务串行化
- 不替代业务事务，而是提供分布式协调基础能力

## 核心能力

### 能力清单

- 支持可重入锁、读写锁、租约锁和选主锁
- 支持锁续约、抢占、超时释放和死锁预防
- 支持基于资源路径的命名空间隔离和审计

### 设计价值

- 统一分布式锁语义，减少各业务重复实现
- 用租约和 fencing token 降低误锁和脏写风险

## 设计文档

以下设计章节已拆分为独立文档：

- [概要设计](/zh/products/stellpoint/summary-design)
- [架构组成](/zh/products/stellpoint/architecture)
- [部署形态](/zh/products/stellpoint/deployment)
- [快速入门](/zh/products/stellpoint/quick-start)
- [配置建议](/zh/products/stellpoint/configuration)
- [API 与 SDK](/zh/products/stellpoint/api-and-sdk)
- [可观测性](/zh/products/stellpoint/observability)

## 典型场景

### 业务场景

- 库存扣减串行化
- 定时任务互斥

### 平台场景

- 主节点选举
- 分布式协调中心
