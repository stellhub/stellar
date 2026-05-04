---
title: Stellspec System Architecture
outline: deep
---

# Stellspec · System Architecture

> A log platform for unified collection, structured processing, search analysis, and retention governance.

## Component Model

- Log Agent: collects files, stdout, and remote log streams.
- Pipeline Engine: performs parsing, desensitization, and routing.
- Query Console: provides search, aggregation, and trace-correlation views.

## Interaction Flow

- Collection agents send logs into processing pipelines for structuring and desensitization.
- Query interfaces reconstruct context by service, Trace ID, or indexed fields.

## Continue Reading

- Start with the [Stellspec product overview](/products/stellspec/)
- Previous: [Design Overview](/products/stellspec/summary-design)
- Next: [Deployment Model](/products/stellspec/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellspec/architecture)
