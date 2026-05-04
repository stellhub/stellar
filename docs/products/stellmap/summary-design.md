---
title: Stellmap Design Overview
outline: deep
---

# Stellmap · Design Overview

> A registry center that manages service instance registration, health checks, discovery subscriptions, and topology change propagation.

## Control Plane

- The control plane manages service metadata, adapts registration protocols, and converges instance state.
- Service definitions, grouping labels, and subscription relationships are all maintained centrally.

## Data Plane

- The data plane is responsible for high-availability queries, incremental push delivery, and local client-side cache support.
- A replicated consistency layer and lease model ensure instance information can be recovered after failures.

## Continue Reading

- Start with the [Stellmap product overview](/products/stellmap/)
- Next: [System Architecture](/products/stellmap/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellmap/summary-design)
