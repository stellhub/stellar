---
title: Stellpoint Design Overview
outline: deep
---

# Stellpoint · Design Overview

> A distributed locking and coordination center for mutual exclusion, leader election, and serialized access to critical resources.

## Control Plane

- The lock control plane manages the resource model, permissions, and lock-strategy publication.
- Resource paths, lease policies, and permission boundaries are maintained centrally.

## Data Plane

- The coordination storage layer persists lock state and lease information through a consistent log.
- Clients preserve ordering guarantees through session keepalive and fencing tokens.

## Continue Reading

- Start with the [Stellpoint product overview](/products/stellpoint/)
- Next: [System Architecture](/products/stellpoint/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellpoint/summary-design)
