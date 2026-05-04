---
title: Stellspec · 星谱 · 架构组成
outline: deep
---

# 架构组成

> 日志平台，负责统一采集、结构化处理、检索分析与日志留存治理。

[返回产品首页](/zh/products/stellspec/)

## 核心组件

- Log Agent：收集文件、标准输出和远程日志流
- Pipeline Engine：完成解析、脱敏和路由
- Query Console：提供检索、聚合和追踪关联视图

## 关键交互

- 采集 Agent 将日志送入处理管道进行结构化与脱敏
- 查询端按服务、Trace ID 或字段索引还原上下文
