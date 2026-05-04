---
title: Stellpulse Design Overview
outline: deep
---

# Stellpulse · Design Overview

> A flow-control and circuit-breaking platform for hotspot protection, capacity guarding, bulkhead isolation, and adaptive degradation.

## Control Plane

- The strategy center manages rule configuration, dynamic push, and version control.
- Rule publishing, canary rollout, and rollback are handled centrally.

## Data Plane

- The runtime engine makes in-process decisions quickly to keep limiting overhead low.
- A status-reporting module reports circuit windows, hotspot statistics, and system load.

## Continue Reading

- Start with the [Stellpulse product overview](/products/stellpulse/)
- Next: [System Architecture](/products/stellpulse/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellpulse/summary-design)
