---
title: Stellcon Deployment Model
outline: deep
---

# Stellcon · Deployment Model

> A metrics platform for collection, aggregation, query analysis, and capacity dashboards.

## Deployment Topology

- In production, separate hot and cold storage and split storage clusters by tenant or business domain.
- For high-cardinality scenarios, use label allowlists and sampling strategy.

## Availability Strategy

- Scale collection and query paths independently to avoid peak interference.
- Store long-window aggregation and raw detail in different tiers to control cost.

## Continue Reading

- Start with the [Stellcon product overview](/products/stellcon/)
- Previous: [System Architecture](/products/stellcon/architecture)
- Next: [Getting Started](/products/stellcon/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellcon/deployment)
