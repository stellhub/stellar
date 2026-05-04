---
title: Stellspec Design Overview
outline: deep
---

# Stellspec · Design Overview

> A log platform for unified collection, structured processing, search analysis, and retention governance.

## Control Plane

- Collection rules, index templates, and lifecycle policy are orchestrated centrally.
- Audit logs and business logs can be managed in separate layers from the control plane.

## Data Plane

- The collection layer ingests logs from agents, sidecars, and gateways.
- The processing layer handles parsing, cleanup, desensitization, and index generation.
- The query layer provides full-text search, field filtering, and context correlation.

## Continue Reading

- Start with the [Stellspec product overview](/products/stellspec/)
- Next: [System Architecture](/products/stellspec/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellspec/summary-design)
