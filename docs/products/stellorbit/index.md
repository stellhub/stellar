---
title: Stellorbit Design
outline: deep
---

# Stellorbit · Orbit

<div class="product-logo">
  <img src="/logo/stellorbit.png" alt="Stellorbit Logo">
</div>

> A service-governance hub for routing, load balancing, traffic control orchestration, and service lifecycle governance.

## Product Scope

### Objective

Stellorbit provides unified governance policies for service-to-service calls, solving dynamic control problems such as canary release, geo-routing, version isolation, and failover.

### Boundaries

- Designed for routing, retries, traffic shifting, and failure isolation in service call paths.
- Designed for governance rule orchestration across multiple versions, regions, and tenants.
- It does not store registry or configuration data directly; instead, it consumes those external contracts.

## Core Capabilities

### Capabilities

- Supports weighted, label-based, version-based, and geo-based routing.
- Supports load balancing, retries, timeout policy, and unhealthy-instance ejection.
- Supports progressive activation and rollback of governance rules.

### Engineering Value

- Moves traffic governance from embedded code to platform-level orchestration.
- Improves control over canary release and active-active traffic switching.

## Reference Sections

- [Design Overview](/products/stellorbit/summary-design)
- [System Architecture](/products/stellorbit/architecture)
- [Deployment Model](/products/stellorbit/deployment)
- [Getting Started](/products/stellorbit/quick-start)
- [Configuration Guide](/products/stellorbit/configuration)
- [API Reference](/products/stellorbit/api-and-sdk)
- [Observability Guide](/products/stellorbit/observability)

## Typical Use Cases

### Business Use Cases

- Canary release
- Active-active traffic switching

### Platform Use Cases

- Service failure isolation
- Cross-cluster governance orchestration

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellorbit/)
