---
title: Stellvox · 星讯 · 架构组成
outline: deep
---

# 架构组成

> 告警平台，负责告警规则、收敛降噪、通知编排和事件协同处置。

[返回产品首页](/zh/products/stellvox/)

## 核心组件

- Alert Rule Engine：执行规则和收敛逻辑
- Notification Hub：对接邮件、IM、短信和 Webhook
- On-call Center：管理排班、升级和确认闭环

## 关键交互

- 规则引擎汇总多源事件并进行去重、抑制和分组
- 通知中心根据值班策略将事件投递到对应责任人
