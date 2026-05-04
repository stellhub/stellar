---
title: Stellspec Deployment Model
outline: deep
---

# Stellspec · Deployment Model

> A log platform for unified collection, structured processing, search analysis, and retention governance.

## Deployment Topology

- Store hot and cold logs in separate tiers to balance query performance and cost.
- Audit logs should use independent storage and independent access control.

## Availability Strategy

- Query nodes and index nodes should be deployed separately to avoid contention.
- For high-throughput ingestion, use buffering queues and batched writes.

## Continue Reading

- Start with the [Stellspec product overview](/products/stellspec/)
- Previous: [System Architecture](/products/stellspec/architecture)
- Next: [Getting Started](/products/stellspec/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellspec/deployment)
