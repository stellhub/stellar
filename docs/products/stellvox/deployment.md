---
title: Stellvox Deployment Model
outline: deep
---

# Stellvox · Deployment Model

> An alerting platform for alert rules, deduplication, notification orchestration, and incident coordination.

## Deployment Topology

- Deploy event computation separately from notification delivery to improve resilience during alert peaks.
- For multi-region operation, support both local notification paths and a global aggregation channel.

## Availability Strategy

- Notification delivery should support retry, rate limiting, and multi-channel fallback.
- On-call configuration and rule configuration both need versioning and rapid rollback.

## Continue Reading

- Start with the [Stellvox product overview](/products/stellvox/)
- Previous: [System Architecture](/products/stellvox/architecture)
- Next: [Getting Started](/products/stellvox/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellvox/deployment)
