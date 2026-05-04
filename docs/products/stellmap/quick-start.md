---
title: Stellmap Getting Started
outline: deep
---

# Stellmap · Getting Started

> A registry center that manages service instance registration, health checks, discovery subscriptions, and topology change propagation.

## Setup

1. Create a namespace and a normalized service identity.
2. Integrate `stellmap-client` into the application and enable heartbeat renewal.
3. Enable service subscription and local cache support on the caller side.

## Validation

- Check through the console or API that the instance has been registered successfully.
- Simulate instance removal and confirm that subscribers receive the change event.

## Continue Reading

- Start with the [Stellmap product overview](/products/stellmap/)
- Previous: [Deployment Model](/products/stellmap/deployment)
- Next: [Configuration Guide](/products/stellmap/configuration)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellmap/quick-start)
