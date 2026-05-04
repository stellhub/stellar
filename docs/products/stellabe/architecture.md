---
title: Stellabe System Architecture
outline: deep
---

# Stellabe · System Architecture

> A job-scheduling platform for timed orchestration, dependency scheduling, sharded execution, and runtime governance.

## Component Model

- Scheduler: triggers jobs and computes execution plans.
- Executor: pulls jobs, executes them, and reports results.
- Workflow Store: stores job definitions, run records, and audit information.

## Interaction Flow

- The scheduler generates execution plans according to trigger conditions and dispatches them to executors.
- Executors report execution state and logs to form a full runtime record.

## Continue Reading

- Start with the [Stellabe product overview](/products/stellabe/)
- Previous: [Design Overview](/products/stellabe/summary-design)
- Next: [Deployment Model](/products/stellabe/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellabe/architecture)
