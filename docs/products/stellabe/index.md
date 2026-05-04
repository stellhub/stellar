---
title: Stellabe Design
outline: deep
---

# Stellabe · Astrolabe

<div class="product-logo">
  <img src="/logo/stellabe.png" alt="Stellabe Logo">
</div>

> A job-scheduling platform for timed orchestration, dependency scheduling, sharded execution, and runtime governance.

## Product Scope

### Objective

Stellabe provides unified scheduling for platform-level and business-level jobs, including batch processing, periodic execution, workflow orchestration, and failure retry.

### Boundaries

- Designed for scheduled jobs, workflows, and batch-processing scenarios.
- Designed for replay, backfill, and visualization of execution progress.
- It does not replace stream-processing engines; it focuses on orchestration and scheduling control.

## Core Capabilities

### Capabilities

- Supports cron scheduling, workflow DAGs, task dependencies, and sharded execution.
- Supports failure retry, idempotency control, backfill, and replay execution.
- Supports tenant isolation, task auditing, and runtime visualization.

### Engineering Value

- Brings scattered scripts and manual tasks into a unified governance model.
- Improves the maintainability and execution reliability of task orchestration.

## Reference Sections

- [Design Overview](/products/stellabe/summary-design)
- [System Architecture](/products/stellabe/architecture)
- [Deployment Model](/products/stellabe/deployment)
- [Getting Started](/products/stellabe/quick-start)
- [Configuration Guide](/products/stellabe/configuration)
- [API Reference](/products/stellabe/api-and-sdk)
- [Observability Guide](/products/stellabe/observability)

## Typical Use Cases

### Business Use Cases

- Scheduled reporting
- Data backfill

### Platform Use Cases

- Post-release task orchestration
- Platform-level workflow governance

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellabe/)
