---
title: Stellcon Design Overview
outline: deep
---

# Stellcon · Design Overview

> A metrics platform for collection, aggregation, query analysis, and capacity dashboards.

## Control Plane

- Collection rules, label boundaries, and dashboard templates are configured centrally.
- SLO targets and aggregation standards are distributed from the control plane.

## Data Plane

- The collection layer handles scraping, pushing, and remote write ingestion.
- The storage layer handles compression, sharding, and long-window aggregation.
- The analysis layer handles alert queries, SLO computation, and trend forecasting.

## Continue Reading

- Start with the [Stellcon product overview](/products/stellcon/)
- Next: [System Architecture](/products/stellcon/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellcon/summary-design)
