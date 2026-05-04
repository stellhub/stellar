---
title: Stelltrace Design Overview
outline: deep
---

# Stelltrace · Design Overview

> A tracing platform for end-to-end trace collection, span analysis, cross-signal correlation, and issue localization.

## Control Plane

- The ingestion layer handles protocol access, context propagation, and sampling decisions.
- Sampling strategy, field conventions, and retention windows are managed centrally.

## Data Plane

- The processing layer aggregates spans, builds indexes, and computes dependency relationships.
- The query layer serves single-trace details, topology search, and troubleshooting views.

## Continue Reading

- Start with the [Stelltrace product overview](/products/stelltrace/)
- Next: [System Architecture](/products/stelltrace/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stelltrace/summary-design)
