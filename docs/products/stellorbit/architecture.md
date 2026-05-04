---
title: Stellorbit System Architecture
outline: deep
---

# Stellorbit · System Architecture

> A service-governance hub for routing, load balancing, traffic control orchestration, and service lifecycle governance.

## Component Model

- Rule Repository: stores routing and governance rules.
- Execution Engine: enforces strategy at traffic ingress and client-call points.
- State Synchronizer: reports instance health, capacity, and abnormal events.

## Interaction Flow

- Clients pull or receive governance rules and execute them immediately in the request path.
- Execution results and instance status are sent back for policy optimization.

## Continue Reading

- Start with the [Stellorbit product overview](/products/stellorbit/)
- Previous: [Design Overview](/products/stellorbit/summary-design)
- Next: [Deployment Model](/products/stellorbit/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellorbit/architecture)
