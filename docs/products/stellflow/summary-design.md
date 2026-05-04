---
title: Stellflow Design Overview
outline: deep
---

# Stellflow · Design Overview

> A message queue and event-stream platform for asynchronous decoupling, streaming distribution, burst smoothing, and event-driven integration.

## Control Plane

- The governance layer manages topic configuration, quotas, audit, and event tracking.
- Topics, consumer groups, and replication strategy are administered centrally.

## Data Plane

- The write layer handles message durability, partition routing, and write acknowledgement.
- The distribution layer handles consumer-group coordination, offset management, and rebalancing.

## Continue Reading

- Start with the [Stellflow product overview](/products/stellflow/)
- Next: [System Architecture](/products/stellflow/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellflow/summary-design)
