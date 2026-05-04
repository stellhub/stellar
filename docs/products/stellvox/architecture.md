---
title: Stellvox System Architecture
outline: deep
---

# Stellvox · System Architecture

> An alerting platform for alert rules, deduplication, notification orchestration, and incident coordination.

## Component Model

- Alert Rule Engine: executes alert rules and convergence logic.
- Notification Hub: integrates email, IM, SMS, and webhook channels.
- On-call Center: manages scheduling, escalation, acknowledgement, and closure.

## Interaction Flow

- The rule engine aggregates multi-source events and deduplicates, suppresses, and groups them.
- The notification center sends incidents to the proper owners according to on-call strategy.

## Continue Reading

- Start with the [Stellvox product overview](/products/stellvox/)
- Previous: [Design Overview](/products/stellvox/summary-design)
- Next: [Deployment Model](/products/stellvox/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellvox/architecture)
