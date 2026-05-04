---
title: Stellmap Design
outline: deep
---

# Stellmap · StarMap

<div class="product-logo">
  <img src="/logo/stellmap.png" alt="Stellmap Logo">
</div>

> A registry center that manages service instance registration, health checks, discovery subscriptions, and topology change propagation.

## Product Scope

### Objective

Stellmap provides unified service registration and discovery for microservices and infrastructure components, helping teams handle endpoint drift, environment switching, elastic scaling, and multi-cluster routing.

### Boundaries

- Designed for service registration, instance discovery, and subscription-driven change propagation.
- Designed for service directory governance across multiple environments, clusters, and canary routing scenarios.
- It does not carry business configuration or application-layer traffic governance logic.

## Core Capabilities

### Capabilities

- Supports service registration, instance deregistration, heartbeat renewal, and health checking.
- Supports namespace, cluster, label, and version-based instance isolation.
- Supports subscription push and local client caching to reduce lookup latency.

### Engineering Value

- Standardizes how service entry points are discovered, reducing the operational cost of maintaining endpoints manually.
- Uses labels and versions to support canary rollout and cross-region routing.

## Reference Sections

- [Design Overview](/products/stellmap/summary-design)
- [System Architecture](/products/stellmap/architecture)
- [Deployment Model](/products/stellmap/deployment)
- [Getting Started](/products/stellmap/quick-start)
- [Configuration Guide](/products/stellmap/configuration)
- [API Reference](/products/stellmap/api-and-sdk)
- [Observability Guide](/products/stellmap/observability)

## Typical Use Cases

### Business Use Cases

- Microservice instance discovery
- Multi-cluster canary routing

### Platform Use Cases

- Infrastructure node directory management
- Unified service catalog governance

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellmap/)
