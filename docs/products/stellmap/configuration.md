---
title: Stellmap Configuration Guide
outline: deep
---

# Stellmap · Configuration Guide

> A registry center that manages service instance registration, health checks, discovery subscriptions, and topology change propagation.

## Baseline Configuration

- In production, separate registration timeout, removal threshold, and push batch size instead of relying on a single global timing model.
- For cross-datacenter traffic, add region labels and nearest-routing policies.

## Production Guidance

- Tune registration timeout and heartbeat intervals according to service startup time and network quality.
- Balance push batch size and local cache refresh frequency against freshness and throughput.

## Continue Reading

- Start with the [Stellmap product overview](/products/stellmap/)
- Previous: [Getting Started](/products/stellmap/quick-start)
- Next: [API Reference](/products/stellmap/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellmap/configuration)
