---
title: Stellflow Configuration Guide
outline: deep
---

# Stellflow · Configuration Guide

> A message queue and event-stream platform for asynchronous decoupling, streaming distribution, burst smoothing, and event-driven integration.

## Baseline Configuration

- Design separate topics for ordered-message workloads and high-throughput workloads.
- Set retry count and dead-letter thresholds according to business idempotency capability.

## Production Guidance

- Plan partition count based on throughput, concurrency, and hotspot-key distribution.
- Design offset-commit strategy together with consumer idempotency and transaction requirements.

## Continue Reading

- Start with the [Stellflow product overview](/products/stellflow/)
- Previous: [Getting Started](/products/stellflow/quick-start)
- Next: [API Reference](/products/stellflow/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellflow/configuration)
