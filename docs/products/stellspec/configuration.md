---
title: Stellspec Configuration Guide
outline: deep
---

# Stellspec · Configuration Guide

> A log platform for unified collection, structured processing, search analysis, and retention governance.

## Baseline Configuration

- Standardize log field naming to avoid index fragmentation.
- Sensitive fields should be desensitized in the collection pipeline rather than during query time.

## Production Guidance

- Index high-cardinality fields carefully to avoid runaway storage and query cost.
- Configure retention windows in layers based on audit, troubleshooting, and compliance requirements.

## Continue Reading

- Start with the [Stellspec product overview](/products/stellspec/)
- Previous: [Getting Started](/products/stellspec/quick-start)
- Next: [API Reference](/products/stellspec/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellspec/configuration)
