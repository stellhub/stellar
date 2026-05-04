---
title: Stellnula Design
outline: deep
---

# Stellnula · Nebula

<div class="product-logo">
  <img src="/logo/stellnula.png" alt="Stellnula Logo">
</div>

> A configuration center responsible for centralized storage, version management, progressive release, and dynamic distribution.

## Product Scope

### Objective

Stellnula provides unified configuration governance for services, gateways, and platform components, helping teams eliminate scattered configuration, environment drift, and unaudited manual changes.

### Boundaries

- Designed for application parameters, platform switches, and runtime configuration distribution.
- Designed for configuration version management, progressive rollout, and audit traceability.
- Sensitive configuration should be linked with the secret-management center instead of being stored in plain text.

## Core Capabilities

### Capabilities

- Supports configuration isolation by namespace, application, environment, and label.
- Supports version history, progressive release, rollback, and change auditing.
- Supports push-based listeners, long polling, and configuration snapshot fallback.

### Engineering Value

- Reduces coupling between application deployment and parameter adjustment.
- Brings configuration governance into a unified change-management workflow to reduce operator mistakes.

## Reference Sections

- [Design Overview](/products/stellnula/summary-design)
- [System Architecture](/products/stellnula/architecture)
- [Deployment Model](/products/stellnula/deployment)
- [Getting Started](/products/stellnula/quick-start)
- [Configuration Guide](/products/stellnula/configuration)
- [API Reference](/products/stellnula/api-and-sdk)
- [Observability Guide](/products/stellnula/observability)

## Typical Use Cases

### Business Use Cases

- Dynamic configuration delivery for microservices
- Canary parameter experiments

### Platform Use Cases

- Multi-environment configuration governance
- Unified platform switch management

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellnula/)
