---
title: Stellgate Design
outline: deep
---

# Stellgate · Event Horizon

<div class="product-logo">
  <img src="/logo/stellgate.png" alt="Stellgate Logo">
</div>

> A gateway ingress platform for unified access, authentication, protocol translation, traffic governance, and API exposure.

## Product Scope

### Objective

Stellgate acts as the external entry point and internal traffic-orchestration hub, receiving client requests, applying unified security policies, and forwarding traffic to backend services.

### Boundaries

- Designed for external APIs, internal service ingress, and edge-access scenarios.
- Designed for protocol access, authentication and authorization, forwarding governance, and plugin extensibility.
- It does not replace a service-governance platform; it focuses on ingress-side traffic handling.

## Core Capabilities

### Capabilities

- Supports HTTP, gRPC, WebSocket, and SSE access.
- Supports authentication, rate limiting, routing, rewriting, and plugin-based extension.
- Supports multi-tenant, multi-domain, and multi-environment gateway orchestration.

### Engineering Value

- Centralizes ingress governance, security, and protocol processing into a unified gateway.
- Improves consistency for both external API exposure and internal traffic orchestration.

## Reference Sections

- [Design Overview](/products/stellgate/summary-design)
- [System Architecture](/products/stellgate/architecture)
- [Deployment Model](/products/stellgate/deployment)
- [Getting Started](/products/stellgate/quick-start)
- [Configuration Guide](/products/stellgate/configuration)
- [API Reference](/products/stellgate/api-and-sdk)
- [Observability Guide](/products/stellgate/observability)

## Typical Use Cases

### Business Use Cases

- Unified Open API ingress
- LLM gateway and edge access

### Platform Use Cases

- Internal service gateway
- Edge traffic ingress

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellgate/)
