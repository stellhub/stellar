---
title: Stellorbit Design Overview
outline: deep
---

# Stellorbit · Design Overview

> A service-governance hub for routing, load balancing, traffic control orchestration, and service lifecycle governance.

## Control Plane

- A strategy center manages routing rules and governance templates centrally.
- Rule publication, canary rollout, and rollback are all controlled from the control plane.

## Data Plane

- The execution plane enforces policies in real time through sidecars or SDKs.
- The state plane reports circuit state, service profiles, and capacity data back to the governance center.

## Continue Reading

- Start with the [Stellorbit product overview](/products/stellorbit/)
- Next: [System Architecture](/products/stellorbit/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellorbit/summary-design)
