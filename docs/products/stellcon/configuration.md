---
title: Stellcon Configuration Guide
outline: deep
---

# Stellcon · Configuration Guide

> A metrics platform for collection, aggregation, query analysis, and capacity dashboards.

## Baseline Configuration

- Keep metric labels within a governable scope to avoid unbounded growth.
- Design histogram bucket boundaries according to actual latency distributions.

## Production Guidance

- Control high-cardinality labels through allowlists or dimension reduction.
- Layer retention periods and aggregation granularity according to cost and diagnostic need.

## Continue Reading

- Start with the [Stellcon product overview](/products/stellcon/)
- Previous: [Getting Started](/products/stellcon/quick-start)
- Next: [API Reference](/products/stellcon/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellcon/configuration)
