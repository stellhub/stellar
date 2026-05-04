---
title: Stellpoint Design
outline: deep
---

# Stellpoint · Singularity

<div class="product-logo">
  <img src="/logo/stellpoint.png" alt="Stellpoint Logo">
</div>

> A distributed locking and coordination center for mutual exclusion, leader election, and serialized access to critical resources.

## Product Scope

### Objective

Stellpoint addresses shared-resource contention in distributed systems by providing a unified lock service, lease model, and coordination primitives, reducing repeated one-off implementations.

### Boundaries

- Designed for mutual exclusion, leader election, and ordered execution.
- Designed for shared-resource write protection and serialization of critical tasks.
- It does not replace business transactions; it provides foundational distributed coordination.

## Core Capabilities

### Capabilities

- Supports reentrant locks, read-write locks, lease locks, and leader-election locks.
- Supports lock renewal, preemption, timeout release, and deadlock prevention.
- Supports namespace isolation and auditability based on resource paths.

### Engineering Value

- Unifies distributed-lock semantics and avoids repeated custom implementations in business systems.
- Uses leases and fencing tokens to reduce accidental lock misuse and dirty writes.

## Reference Sections

- [Design Overview](/products/stellpoint/summary-design)
- [System Architecture](/products/stellpoint/architecture)
- [Deployment Model](/products/stellpoint/deployment)
- [Getting Started](/products/stellpoint/quick-start)
- [Configuration Guide](/products/stellpoint/configuration)
- [API Reference](/products/stellpoint/api-and-sdk)
- [Observability Guide](/products/stellpoint/observability)

## Typical Use Cases

### Business Use Cases

- Serialized inventory deduction
- Mutual exclusion for scheduled tasks

### Platform Use Cases

- Leader election
- Distributed coordination center

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellpoint/)
