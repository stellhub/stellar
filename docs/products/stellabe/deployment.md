---
title: Stellabe Deployment Model
outline: deep
---

# Stellabe · Deployment Model

> A job-scheduling platform for timed orchestration, dependency scheduling, sharded execution, and runtime governance.

## Deployment Topology

- Deploy the control plane separately from the execution plane so job peaks do not impact the management console.
- For heavy workloads, split worker clusters by queue type or workload type.

## Availability Strategy

- Schedulers must guarantee single-trigger semantics and support preemption recovery.
- Executors should isolate resource pools by workload type to avoid mutual interference.

## Continue Reading

- Start with the [Stellabe product overview](/products/stellabe/)
- Previous: [System Architecture](/products/stellabe/architecture)
- Next: [Getting Started](/products/stellabe/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellabe/deployment)
