---
title: Stellflow · 彗流 · 架构组成
outline: deep
---

# 架构组成

> 消息队列与事件流平台，负责异步解耦、流式分发、削峰填谷和事件驱动集成。

[返回产品首页](/zh/products/stellflow/)

## 核心组件

- Broker：承载消息写入、存储与分发
- Controller：管理主题、分区、副本和故障转移
- Client SDK：提供生产、消费和事务消息能力

## 关键交互

- 生产端将消息写入 Broker 并获得确认结果
- 消费组通过位点管理和重平衡机制获取分区消费权
