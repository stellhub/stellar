---
title: Stelltrace System Architecture
outline: deep
---

# Stelltrace · System Architecture

> A tracing platform for end-to-end trace collection, span analysis, cross-signal correlation, and issue localization.

## Component Model

- Agent / SDK: handles instrumentation and context propagation.
- Collector: receives spans and performs pre-aggregation.
- Query API: provides search, analysis, and visualization interfaces.

## Interaction Flow

- Agents report spans to the collector for aggregation and cleanup.
- The query layer reconstructs execution paths by Trace ID and correlates them with logs and metrics.

## Continue Reading

- Start with the [Stelltrace product overview](/products/stelltrace/)
- Previous: [Design Overview](/products/stelltrace/summary-design)
- Next: [Deployment Model](/products/stelltrace/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stelltrace/architecture)
