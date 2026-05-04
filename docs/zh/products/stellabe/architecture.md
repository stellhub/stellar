---
title: Stellabe · 星盘 · 架构组成
outline: deep
---

# 架构组成

> 任务调度平台，负责定时任务编排、依赖调度、分片执行与运行治理。

[返回产品首页](/zh/products/stellabe/)

## 核心组件

- Scheduler：触发任务并计算执行计划
- Executor：拉取任务、执行并回报结果
- Workflow Store：保存任务定义、运行记录和审计信息

## 关键交互

- Scheduler 依据触发条件生成执行计划并分派到 Executor
- Executor 上报执行状态和日志，形成完整运行记录
