---
title: Stellorbit Deployment Model
outline: deep
---

# Stellorbit · Deployment Model

> A service-governance hub for routing, load balancing, traffic control orchestration, and service lifecycle governance.

## Deployment Topology

- Support three governance modes: pure SDK, sidecar, and gateway-integrated governance.
- Deploy the core control plane with multiple replicas to keep rule publication continuous.

## Availability Strategy

- Rule cache must support local persistence and disconnect tolerance.
- Gateway-side and service-side governance can be combined to reduce single-point policy failure.

## Continue Reading

- Start with the [Stellorbit product overview](/products/stellorbit/)
- Previous: [System Architecture](/products/stellorbit/architecture)
- Next: [Getting Started](/products/stellorbit/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellorbit/deployment)
