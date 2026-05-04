---
title: Stelltrace Design
outline: deep
---

# Stelltrace · StarTrace

<div class="product-logo">
  <img src="/logo/stelltrace.png" alt="Stelltrace Logo">
</div>

> A tracing platform for end-to-end trace collection, span analysis, cross-signal correlation, and issue localization.

## Product Scope

### Objective

Stelltrace provides unified tracing for service call graphs, asynchronous tasks, and cross-gateway requests so engineering and operations teams can quickly identify slow requests, error propagation, and cross-system dependencies.

### Boundaries

- Designed for request tracing across synchronous calls, asynchronous messaging, and scheduled tasks.
- Designed for troubleshooting, topology analysis, and cross-system dependency investigation.
- It does not replace the log platform or metrics platform; instead, it acts as their correlation backbone.

## Core Capabilities

### Capabilities

- Supports tracing across HTTP, gRPC, MQ, scheduled tasks, and other request paths.
- Supports sampling policies, anomaly clustering, and dependency-topology views.
- Supports trace-correlated queries across logs, metrics, and alerts.

### Engineering Value

- Uses Trace ID as the centerline that connects multiple observability signals.
- Improves troubleshooting efficiency for slow requests and abnormal propagation.

## Reference Sections

- [Design Overview](/products/stelltrace/summary-design)
- [System Architecture](/products/stelltrace/architecture)
- [Deployment Model](/products/stelltrace/deployment)
- [Getting Started](/products/stelltrace/quick-start)
- [Configuration Guide](/products/stelltrace/configuration)
- [API Reference](/products/stelltrace/api-and-sdk)
- [Observability Guide](/products/stelltrace/observability)

## Typical Use Cases

### Business Use Cases

- Slow request troubleshooting
- Root-cause analysis for anomalies

### Platform Use Cases

- Service dependency topology governance
- Full-path observability correlation analysis

## Chinese Source

- [Read the original Chinese product page](/zh/products/stelltrace/)
