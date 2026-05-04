---
title: Stellnula Deployment Model
outline: deep
---

# Stellnula · Deployment Model

> A configuration center responsible for centralized storage, version management, progressive release, and dynamic distribution.

## Deployment Topology

- Deploy the configuration center separately from the registry center so configuration traffic does not interfere with service discovery.
- At larger scale, split the read-only query layer from the release execution layer.

## Availability Strategy

- Separate the release path from the query path to reduce mutual impact during heavy change windows.
- Configuration snapshots should support cross-instance sharing and cold-start fallback.

## Continue Reading

- Start with the [Stellnula product overview](/products/stellnula/)
- Previous: [System Architecture](/products/stellnula/architecture)
- Next: [Getting Started](/products/stellnula/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellnula/deployment)
