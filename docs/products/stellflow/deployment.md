---
title: Stellflow Deployment Model
outline: deep
---

# Stellflow · Deployment Model

> A message queue and event-stream platform for asynchronous decoupling, streaming distribution, burst smoothing, and event-driven integration.

## Deployment Topology

- Plan topics and partitions by business domain to avoid hotspot concentration.
- For cross-region traffic, deploy mirror replication or event-bridge channels.

## Availability Strategy

- Replica strategy must tolerate both single-node failures and single-availability-zone failures.
- Hot topics should use dedicated clusters or isolated resource pools.

## Continue Reading

- Start with the [Stellflow product overview](/products/stellflow/)
- Previous: [System Architecture](/products/stellflow/architecture)
- Next: [Getting Started](/products/stellflow/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellflow/deployment)
