---
title: Stellguard Design
outline: deep
---

# Stellguard · StarShield

<div class="product-logo">
  <img src="/logo/stellguard.png" alt="Stellguard Logo">
</div>

> A zero-trust security platform for identity verification, access control, policy evaluation, and secure service-to-service communication.

## Product Scope

### Objective

Stellguard applies zero-trust principles to user, service, and endpoint access so that default deny, continuous verification, and least privilege can be implemented consistently across the platform.

### Boundaries

- Designed for unified identity and access control across users, services, and endpoints.
- Designed for mTLS, context-based authorization, and risk interception across service boundaries.
- It does not replace business permission models; it provides a consistent security foundation.

## Core Capabilities

### Capabilities

- Supports identity authentication for users, services, and devices.
- Supports fine-grained policy evaluation, dynamic authorization, and risk interception.
- Supports service-to-service mTLS, certificate rotation, and access auditing.

### Engineering Value

- Unifies identity, authorization, and communication security under the same control plane.
- Makes least privilege and continuous verification enforceable at the platform level.

## Reference Sections

- [Design Overview](/products/stellguard/summary-design)
- [System Architecture](/products/stellguard/architecture)
- [Deployment Model](/products/stellguard/deployment)
- [Getting Started](/products/stellguard/quick-start)
- [Configuration Guide](/products/stellguard/configuration)
- [API Reference](/products/stellguard/api-and-sdk)
- [Observability Guide](/products/stellguard/observability)

## Typical Use Cases

### Business Use Cases

- Zero-trust access between services
- Audit for highly sensitive operations

### Platform Use Cases

- Unified console authentication
- Platform-level continuous-verification governance

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellguard/)
