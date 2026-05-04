---
title: Stellvox Design Overview
outline: deep
---

# Stellvox · Design Overview

> An alerting platform for alert rules, deduplication, notification orchestration, and incident coordination.

## Control Plane

- The rule layer manages strategy configuration, versioning, and progressive activation.
- Notification templates, on-call groups, and escalation paths are managed centrally.

## Data Plane

- The engine layer performs real-time computation, deduplication, convergence, and notification routing.
- The collaboration layer handles ticket linkage, on-call rotation, and disposal tracking.

## Continue Reading

- Start with the [Stellvox product overview](/products/stellvox/)
- Next: [System Architecture](/products/stellvox/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellvox/summary-design)
