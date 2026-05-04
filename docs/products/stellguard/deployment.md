---
title: Stellguard Deployment Model
outline: deep
---

# Stellguard · Deployment Model

> A zero-trust security platform for identity verification, access control, policy evaluation, and secure service-to-service communication.

## Deployment Topology

- Deploy core identity and policy services in isolated, highly available clusters.
- Access control can be enforced collaboratively across gateways, sidecars, and application code.

## Availability Strategy

- Identity and policy services should support partitioned deployment and local cache capability.
- Certificate issuance and validation paths must account for disaster recovery and rotation windows.

## Continue Reading

- Start with the [Stellguard product overview](/products/stellguard/)
- Previous: [System Architecture](/products/stellguard/architecture)
- Next: [Getting Started](/products/stellguard/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellguard/deployment)
