---
title: Stellcon · 星座 · 架构组成
outline: deep
---

# 架构组成

> 指标平台，负责指标采集、聚合计算、查询分析和容量看板建设。

[返回产品首页](/zh/products/stellcon/)

## 核心组件

- Metrics Gateway：承接指标采集流量
- Time Series Store：保存时序数据和聚合结果
- Dashboard API：提供查询、看板和规则配置接口

## 关键交互

- 指标从采集端汇入 Gateway 后写入时序存储
- 查询端按时间窗口和标签条件聚合出可视化结果
