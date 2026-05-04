---
title: Stellcon Design
outline: deep
---

# Stellcon · Constellation

<div class="product-logo">
  <img src="/logo/stellcon.png" alt="Stellcon Logo">
</div>

> A metrics platform for collection, aggregation, query analysis, and capacity dashboards.

## Product Scope

### Objective

Stellcon provides a unified metrics system for applications, gateways, middleware, and infrastructure, supporting capacity estimation, SLO tracking, and alert-threshold computation.

### Boundaries

- Designed for governance of runtime metrics, business metrics, and platform resource metrics.
- Designed for trend analysis, capacity planning, and SLO target management.
- It does not replace log search; it focuses on time-series collection and aggregated query.

## Core Capabilities

### Capabilities

- Supports Counter, Gauge, Histogram, and Summary metric models.
- Supports unified label conventions, aggregation rules, and multidimensional queries.
- Supports dashboard templates, SLO calculation, and capacity-trend analysis.

### Engineering Value

- Creates consistent metric semantics and reduces fragmentation across systems.
- Provides a unified data foundation for alerting, capacity planning, and reliability governance.

## Reference Sections

- [Design Overview](/products/stellcon/summary-design)
- [System Architecture](/products/stellcon/architecture)
- [Deployment Model](/products/stellcon/deployment)
- [Getting Started](/products/stellcon/quick-start)
- [Configuration Guide](/products/stellcon/configuration)
- [API Reference](/products/stellcon/api-and-sdk)
- [Observability Guide](/products/stellcon/observability)

## Typical Use Cases

### Business Use Cases

- SLO observation
- Capacity planning

### Platform Use Cases

- Unified monitoring dashboards
- Multi-tenant metrics governance

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellcon/)
