---
title: Stellgate Deployment Model
outline: deep
---

# Stellgate · Deployment Model

> A gateway ingress platform for unified access, authentication, protocol translation, traffic governance, and API exposure.

## Deployment Topology

- Deploy external and internal gateways in separate layers.
- Edge nodes should connect locally and receive routes from the configuration center.

## Availability Strategy

- Gateway nodes should be deployed statelessly and combined with elastic scaling.
- Certificates, routes, and plugins should all support canary release and fast rollback.

## Continue Reading

- Start with the [Stellgate product overview](/products/stellgate/)
- Previous: [System Architecture](/products/stellgate/architecture)
- Next: [Getting Started](/products/stellgate/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellgate/deployment)
