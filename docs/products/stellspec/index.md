---
title: Stellspec Design
outline: deep
---

# Stellspec · Spectrum

<div class="product-logo">
  <img src="/logo/stellspec.png" alt="Stellspec Logo">
</div>

> A log platform for unified collection, structured processing, search analysis, and retention governance.

## Product Scope

### Objective

Stellspec provides unified ingestion and retrieval for business logs, audit logs, and platform logs so teams can establish consistent logging standards and faster troubleshooting paths.

### Boundaries

- Designed for unified governance of application logs, audit logs, and platform runtime logs.
- Designed for log search, trace correlation, and retention-policy management.
- It does not replace metrics or tracing; it focuses on full-text search and contextual reconstruction.

## Core Capabilities

### Capabilities

- Supports unified ingestion of plain-text logs, structured logs, and audit logs.
- Supports field parsing, desensitization, index routing, and lifecycle management.
- Supports trace correlation, log clustering, and anomaly-pattern search.

### Engineering Value

- Creates consistent logging conventions and reduces search and troubleshooting cost.
- Balances security and cost through desensitization and lifecycle management.

## Reference Sections

- [Design Overview](/products/stellspec/summary-design)
- [System Architecture](/products/stellspec/architecture)
- [Deployment Model](/products/stellspec/deployment)
- [Getting Started](/products/stellspec/quick-start)
- [Configuration Guide](/products/stellspec/configuration)
- [API Reference](/products/stellspec/api-and-sdk)
- [Observability Guide](/products/stellspec/observability)

## Typical Use Cases

### Business Use Cases

- Application log search
- Troubleshooting and root-cause analysis

### Platform Use Cases

- Audit tracing
- Platform-level retention governance

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellspec/)
