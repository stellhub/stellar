---
title: Stellorbit Configuration Guide
outline: deep
---

# Stellorbit · Configuration Guide

> A service-governance hub for routing, load balancing, traffic control orchestration, and service lifecycle governance.

## Baseline Configuration

- Route rules should be layered into global, tenant-level, and service-level strategies.
- Retry and timeout settings should be designed together to avoid amplifying failure traffic.

## Production Guidance

- Critical routing rules should include approval and rollback plans.
- Limit retry depth and timeout stacking to avoid request avalanches.

## Continue Reading

- Start with the [Stellorbit product overview](/products/stellorbit/)
- Previous: [Getting Started](/products/stellorbit/quick-start)
- Next: [API Reference](/products/stellorbit/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellorbit/configuration)
