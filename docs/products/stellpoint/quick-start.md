---
title: Stellpoint Getting Started
outline: deep
---

# Stellpoint · Getting Started

> A distributed locking and coordination center for mutual exclusion, leader election, and serialized access to critical resources.

## Setup

1. Define the resource path and lease duration.
2. Acquire the lock inside business code and execute the critical logic.
3. Use fencing tokens to prevent stale holders from writing.

## Validation

- Simulate two instances competing for the same resource and confirm only one gets the lock.
- Simulate an abnormal exit of the lock holder and confirm the lease is reclaimed after expiration.

## Continue Reading

- Start with the [Stellpoint product overview](/products/stellpoint/)
- Previous: [Deployment Model](/products/stellpoint/deployment)
- Next: [Configuration Guide](/products/stellpoint/configuration)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellpoint/quick-start)
