---
title: Stelltrace Deployment Model
outline: deep
---

# Stelltrace · Deployment Model

> A tracing platform for end-to-end trace collection, span analysis, cross-signal correlation, and issue localization.

## Deployment Topology

- Decouple the ingestion layer from the storage layer so collectors can scale horizontally.
- Store hot and cold data in different tiers to control long-term retention cost.

## Availability Strategy

- Collector clusters should support batch buffering and failover.
- Query nodes should be separated from index nodes to avoid resource contention.

## Continue Reading

- Start with the [Stelltrace product overview](/products/stelltrace/)
- Previous: [System Architecture](/products/stelltrace/architecture)
- Next: [Getting Started](/products/stelltrace/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stelltrace/deployment)
