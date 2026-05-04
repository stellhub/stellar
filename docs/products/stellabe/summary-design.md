---
title: Stellabe Design Overview
outline: deep
---

# Stellabe · Design Overview

> A job-scheduling platform for timed orchestration, dependency scheduling, sharded execution, and runtime governance.

## Control Plane

- The orchestration layer defines workflows, analyzes dependencies, and publishes tasks.
- Task versions, approval flow, and release process are managed centrally.

## Data Plane

- The scheduling layer computes triggers, allocates shards, and controls execution windows.
- The execution layer handles worker registration, status reporting, and failure recovery.

## Continue Reading

- Start with the [Stellabe product overview](/products/stellabe/)
- Next: [System Architecture](/products/stellabe/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellabe/summary-design)
