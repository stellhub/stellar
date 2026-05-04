---
title: Stellpoint Deployment Model
outline: deep
---

# Stellpoint · Deployment Model

> A distributed locking and coordination center for mutual exclusion, leader election, and serialized access to critical resources.

## Deployment Topology

- Deploy an odd number of replicas, usually three to five nodes, to preserve quorum behavior.
- For high-risk resources, use dedicated namespaces or separate lock clusters.

## Availability Strategy

- Session timeout and lease expiration must be tuned carefully to avoid accidental release.
- For hotspot resources, split lock granularity and limit preemption frequency.

## Continue Reading

- Start with the [Stellpoint product overview](/products/stellpoint/)
- Previous: [System Architecture](/products/stellpoint/architecture)
- Next: [Getting Started](/products/stellpoint/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellpoint/deployment)
