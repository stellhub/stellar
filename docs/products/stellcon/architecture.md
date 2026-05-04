---
title: Stellcon System Architecture
outline: deep
---

# Stellcon · System Architecture

> A metrics platform for collection, aggregation, query analysis, and capacity dashboards.

## Component Model

- Metrics Gateway: receives metrics ingestion traffic.
- Time Series Store: stores time-series data and aggregation results.
- Dashboard API: provides query, dashboard, and rule-configuration interfaces.

## Interaction Flow

- Metrics enter through the gateway and are written into time-series storage.
- Query interfaces aggregate results into visual views by time window and label conditions.

## Continue Reading

- Start with the [Stellcon product overview](/products/stellcon/)
- Previous: [Design Overview](/products/stellcon/summary-design)
- Next: [Deployment Model](/products/stellcon/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellcon/architecture)
