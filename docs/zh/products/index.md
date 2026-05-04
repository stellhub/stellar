---
title: 核心产品总览
outline: deep
---

# 核心产品总览

本页汇总 `Stell Hub（星际枢纽）` 核心中间件矩阵中的产品入口。每个中间件在 `products` 目录下都拥有独立目录，便于后续持续扩展子文档、部署说明和实践案例。

## 产品矩阵

<div class="product-grid">
  <a class="product-card" href="/zh/products/stellmap/">
    <div class="product-card-logo"><img src="/logo/stellmap.png" alt="Stellmap Logo"></div>
    <span class="product-card-domain">注册中心</span>
    <h3 class="product-card-title">Stellmap · 星图</h3>
    <p class="product-card-desc">统一管理服务注册、健康探测、实例发现与拓扑变更推送。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellnula/">
    <div class="product-card-logo"><img src="/logo/stellnula.png" alt="Stellnula Logo"></div>
    <span class="product-card-domain">配置中心</span>
    <h3 class="product-card-title">Stellnula · 星云</h3>
    <p class="product-card-desc">承载配置集中存储、灰度发布、动态下发与变更审计能力。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stelltrace/">
    <div class="product-card-logo"><img src="/logo/stelltrace.png" alt="Stelltrace Logo"></div>
    <span class="product-card-domain">链路追踪</span>
    <h3 class="product-card-title">Stelltrace · 星迹</h3>
    <p class="product-card-desc">面向分布式请求链路提供 Trace 采集、检索、分析和定位能力。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellorbit/">
    <div class="product-card-logo"><img src="/logo/stellorbit.png" alt="Stellorbit Logo"></div>
    <span class="product-card-domain">服务治理</span>
    <h3 class="product-card-title">Stellorbit · 星轨</h3>
    <p class="product-card-desc">统一编排路由、负载、重试、切流与服务生命周期治理策略。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellpulse/">
    <div class="product-card-logo"><img src="/logo/stellpulse.png" alt="Stellpulse Logo"></div>
    <span class="product-card-domain">流控熔断</span>
    <h3 class="product-card-title">Stellpulse · 脉冲</h3>
    <p class="product-card-desc">提供限流、熔断、降级、隔离舱和热点保护等高可用能力。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellabe/">
    <div class="product-card-logo"><img src="/logo/stellabe.png" alt="Stellabe Logo"></div>
    <span class="product-card-domain">任务调度</span>
    <h3 class="product-card-title">Stellabe · 星盘</h3>
    <p class="product-card-desc">支持定时任务、工作流编排、分片执行和失败恢复治理。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellpoint/">
    <div class="product-card-logo"><img src="/logo/stellpoint.png" alt="Stellpoint Logo"></div>
    <span class="product-card-domain">分布式锁</span>
    <h3 class="product-card-title">Stellpoint · 奇点</h3>
    <p class="product-card-desc">提供租约锁、选主锁、协调原语与关键资源互斥访问能力。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellgate/">
    <div class="product-card-logo"><img src="/logo/stellgate.png" alt="Stellgate Logo"></div>
    <span class="product-card-domain">网关入口</span>
    <h3 class="product-card-title">Stellgate · 视界</h3>
    <p class="product-card-desc">统一承接入口流量，提供认证鉴权、协议转换和 API 暴露能力。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellflow/">
    <div class="product-card-logo"><img src="/logo/stellflow.png" alt="Stellflow Logo"></div>
    <span class="product-card-domain">消息队列</span>
    <h3 class="product-card-title">Stellflow · 彗流</h3>
    <p class="product-card-desc">支撑异步解耦、事件驱动、流式分发和削峰填谷等消息场景。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellspec/">
    <div class="product-card-logo"><img src="/logo/stellspec.png" alt="Stellspec Logo"></div>
    <span class="product-card-domain">日志平台</span>
    <h3 class="product-card-title">Stellspec · 星谱</h3>
    <p class="product-card-desc">统一采集、解析、检索和留存应用日志、审计日志与平台日志。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellcon/">
    <div class="product-card-logo"><img src="/logo/stellcon.png" alt="Stellcon Logo"></div>
    <span class="product-card-domain">指标平台</span>
    <h3 class="product-card-title">Stellcon · 星座</h3>
    <p class="product-card-desc">提供统一指标采集、聚合分析、容量看板和 SLO 观测能力。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellvox/">
    <div class="product-card-logo"><img src="/logo/stellvox.png" alt="Stellvox Logo"></div>
    <span class="product-card-domain">告警平台</span>
    <h3 class="product-card-title">Stellvox · 星讯</h3>
    <p class="product-card-desc">统一处理规则触发、收敛降噪、值班通知和事件协同处置。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellguard/">
    <div class="product-card-logo"><img src="/logo/stellguard.png" alt="Stellguard Logo"></div>
    <span class="product-card-domain">零信平台</span>
    <h3 class="product-card-title">Stellguard · 星盾</h3>
    <p class="product-card-desc">面向身份、访问控制、mTLS 和持续验证的零信任安全中枢。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
  <a class="product-card" href="/zh/products/stellkey/">
    <div class="product-card-logo"><img src="/logo/stellkey.png" alt="Stellkey Logo"></div>
    <span class="product-card-domain">密钥中心</span>
    <h3 class="product-card-title">Stellkey · 星钥</h3>
    <p class="product-card-desc">统一托管机密、证书和访问策略，支撑轮换、分发和安全审计。</p>
    <span class="product-card-link">查看详细设计</span>
  </a>
</div>

## 文档结构约定

当前每个中间件目录已拆分为“首页概览 + 独立专题子页”的结构：

- 首页概览：产品定位、核心能力、典型场景
- 独立子页：概要设计
- 独立子页：架构组成
- 独立子页：部署形态
- 独立子页：快速入门
- 独立子页：配置建议
- 独立子页：API 与 SDK
- 独立子页：可观测性

这套结构兼顾了产品概览阅读和后续按专题持续扩展，也更适合继续沉淀部署说明、实践案例和版本化设计文档。
