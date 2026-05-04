---
title: Stellpulse System Architecture
outline: deep
---

# Stellpulse · System Architecture

> A flow-control and circuit-breaking platform for hotspot protection, capacity guarding, bulkhead isolation, and adaptive degradation.

## Component Model

- Rule Center: centrally manages rate-limiting and circuit-breaking rules.
- Runtime Engine: decides admission and degradation at request time.
- Metrics Reporter: reports hit, rejection, and recovery events.

## Interaction Flow

- Rules are pushed to runtime components and enforced locally at the request entry point.
- Rejection, circuit-break, and recovery events are returned to the monitoring and alerting systems.

## Continue Reading

- Start with the [Stellpulse product overview](/products/stellpulse/)
- Previous: [Design Overview](/products/stellpulse/summary-design)
- Next: [Deployment Model](/products/stellpulse/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellpulse/architecture)
