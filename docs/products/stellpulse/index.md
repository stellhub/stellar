---
title: Stellpulse Design
outline: deep
---

# Stellpulse · Pulsar

<div class="product-logo">
  <img src="/logo/stellpulse.png" alt="Stellpulse Logo">
</div>

> A flow-control and circuit-breaking platform for hotspot protection, capacity guarding, bulkhead isolation, and adaptive degradation.

## Product Scope

### Objective

Stellpulse focuses on availability protection under high concurrency, using rate limiting, concurrency isolation, circuit breaking, degradation, and system load defense to keep critical services stable.

### Boundaries

- Designed for stability protection at service ingress, critical interfaces, and dependency calls.
- Designed for sudden traffic spikes, hotspot resources, and unstable dependencies.
- It does not replace service governance; it focuses on admission control and degradation protection.

## Core Capabilities

### Capabilities

- Supports QPS limiting, concurrency limiting, and hotspot parameter protection.
- Supports circuit breaking, bulkhead isolation, warmup, and request queuing.
- Supports automatic degradation based on error rate, slow calls, and system load.

### Engineering Value

- Provides a unified availability-protection foundation for critical call paths.
- Keeps the impact of traffic spikes and dependency failures within predictable bounds.

## Reference Sections

- [Design Overview](/products/stellpulse/summary-design)
- [System Architecture](/products/stellpulse/architecture)
- [Deployment Model](/products/stellpulse/deployment)
- [Getting Started](/products/stellpulse/quick-start)
- [Configuration Guide](/products/stellpulse/configuration)
- [API Reference](/products/stellpulse/api-and-sdk)
- [Observability Guide](/products/stellpulse/observability)

## Typical Use Cases

### Business Use Cases

- Traffic protection for flash-sale workloads
- Fallback for critical call paths

### Platform Use Cases

- Dependency failure isolation
- Unified platform-level traffic guard

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellpulse/)
