---
title: Stellpoint System Architecture
outline: deep
---

# Stellpoint · System Architecture

> A distributed locking and coordination center for mutual exclusion, leader election, and serialized access to critical resources.

## Component Model

- Lock API: receives lock, unlock, renewal, and query requests.
- Lease Manager: manages sessions, leases, and expiration cleanup.
- Coordination Log: stores lock ordering and fencing tokens.

## Interaction Flow

- Clients keep leases alive through sessions and periodic renewal.
- The coordination log emits fencing tokens to preserve holder-order semantics.

## Continue Reading

- Start with the [Stellpoint product overview](/products/stellpoint/)
- Previous: [Design Overview](/products/stellpoint/summary-design)
- Next: [Deployment Model](/products/stellpoint/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellpoint/architecture)
