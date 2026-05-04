---
title: Stellkey Design
outline: deep
---

# Stellkey · StarKey

<div class="product-logo">
  <img src="/logo/stellkey.png" alt="Stellkey Logo">
</div>

> A secret-management center for key lifecycle management, secret distribution, rotation auditing, and secure access control.

## Product Scope

### Objective

Stellkey provides unified secret custody and secret distribution for applications and platform components, reducing risk from scattered plain-text keys, difficult manual rotation, and missing audit trails.

### Boundaries

- Designed for unified management of keys, certificates, tokens, and sensitive configuration.
- Designed for access approval, rotation audit, and short-lived credential distribution.
- It does not replace the configuration center directly; it provides a secure foundation for sensitive-value governance.

## Core Capabilities

### Capabilities

- Supports unified custody of keys, certificates, tokens, and secret configuration.
- Supports version rotation, temporary credentials, access approval, and audit tracking.
- Supports coordinated distribution with the configuration center, zero-trust platform, and gateways.

### Engineering Value

- Unifies secret lifecycle management and reduces the risk of plain-text exposure.
- Makes sensitive-configuration delivery and access auditing available as platform-level capabilities.

## Reference Sections

- [Design Overview](/products/stellkey/summary-design)
- [System Architecture](/products/stellkey/architecture)
- [Deployment Model](/products/stellkey/deployment)
- [Getting Started](/products/stellkey/quick-start)
- [Configuration Guide](/products/stellkey/configuration)
- [API Reference](/products/stellkey/api-and-sdk)
- [Observability Guide](/products/stellkey/observability)

## Typical Use Cases

### Business Use Cases

- Database credential custody
- TLS certificate management

### Platform Use Cases

- Sensitive-value references from the configuration center
- Platform-level secret lifecycle governance

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellkey/)
