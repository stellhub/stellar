---
title: Stellflow Design
outline: deep
---

# Stellflow · CometFlow

<div class="product-logo">
  <img src="/logo/stellflow.png" alt="Stellflow Logo">
</div>

> A message queue and event-stream platform for asynchronous decoupling, streaming distribution, burst smoothing, and event-driven integration.

## Product Scope

### Objective

Stellflow provides a unified messaging and event-distribution capability for asynchronous business flows, system integration, and streaming workloads.

### Boundaries

- Designed for asynchronous decoupling, event-driven interaction, and data-distribution scenarios.
- Designed for delayed messages, ordered messages, and stream-consumption models.
- It does not replace a job-scheduling platform; it focuses on durable message storage and distribution.

## Core Capabilities

### Capabilities

- Supports ordinary queues, topic subscription, delayed messages, and ordered messages.
- Supports consumer groups, retry queues, dead-letter queues, and replay consumption.
- Supports event bridging, streaming processing, and cross-region replication.

### Engineering Value

- Reduces coupling and traffic shocks through asynchronous design.
- Supports event-driven architecture through reliable delivery and offset management.

## Reference Sections

- [Design Overview](/products/stellflow/summary-design)
- [System Architecture](/products/stellflow/architecture)
- [Deployment Model](/products/stellflow/deployment)
- [Getting Started](/products/stellflow/quick-start)
- [Configuration Guide](/products/stellflow/configuration)
- [API Reference](/products/stellflow/api-and-sdk)
- [Observability Guide](/products/stellflow/observability)

## Typical Use Cases

### Business Use Cases

- Order event-driven integration
- Asynchronous burst smoothing

### Platform Use Cases

- Data synchronization bus
- Cross-system event bridging

## Chinese Source

- [Read the original Chinese product page](/zh/products/stellflow/)
