---
title: Stellkey Deployment Model
outline: deep
---

# Stellkey · Deployment Model

> A secret-management center for key lifecycle management, secret distribution, rotation auditing, and secure access control.

## Deployment Topology

- Deploy core storage separately from the access gateway to reduce exposure.
- For high-sensitivity workloads, use tiered custody for root keys and business keys.

## Availability Strategy

- Core secret storage should support replication and disaster-recovery restoration.
- The access gateway should support audit caching and fast synchronization when permissions are revoked.

## Continue Reading

- Start with the [Stellkey product overview](/products/stellkey/)
- Previous: [System Architecture](/products/stellkey/architecture)
- Next: [Getting Started](/products/stellkey/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellkey/deployment)
