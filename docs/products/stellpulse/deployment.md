---
title: Stellpulse Deployment Model
outline: deep
---

# Stellpulse · Deployment Model

> A flow-control and circuit-breaking platform for hotspot protection, capacity guarding, bulkhead isolation, and adaptive degradation.

## Deployment Topology

- The preferred mode is SDK-based embedding inside business applications, optionally combined with gateway-level rules.
- For major promotion events, coordinate layered rate limiting with the gateway and service-governance platform.

## Availability Strategy

- Local rule cache should remain effective even when the control plane is unavailable.
- Rate-limiting and circuit-breaking rules should be deployed in multiple layers instead of concentrating on a single entry point.

## Continue Reading

- Start with the [Stellpulse product overview](/products/stellpulse/)
- Previous: [System Architecture](/products/stellpulse/architecture)
- Next: [Getting Started](/products/stellpulse/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellpulse/deployment)
