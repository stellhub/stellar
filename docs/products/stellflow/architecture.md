---
title: Stellflow System Architecture
outline: deep
---

# Stellflow · System Architecture

> A message queue and event-stream platform for asynchronous decoupling, streaming distribution, burst smoothing, and event-driven integration.

## Component Model

- Broker: handles message write, storage, and distribution.
- Controller: manages topics, partitions, replicas, and failover.
- Client SDK: provides producer, consumer, and transactional-message capability.

## Interaction Flow

- Producers write messages to brokers and receive acknowledgement results.
- Consumer groups obtain partition ownership through offset management and rebalancing.

## Continue Reading

- Start with the [Stellflow product overview](/products/stellflow/)
- Previous: [Design Overview](/products/stellflow/summary-design)
- Next: [Deployment Model](/products/stellflow/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellflow/architecture)
