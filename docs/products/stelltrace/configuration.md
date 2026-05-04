---
title: Stelltrace Configuration Guide
outline: deep
---

# Stelltrace · Configuration Guide

> A tracing platform for end-to-end trace collection, span analysis, cross-signal correlation, and issue localization.

## Baseline Configuration

- Use layered sampling in production and prioritize keeping errors and slow requests.
- Enable high-cardinality label allowlists for critical business paths.

## Production Guidance

- Use adaptive sampling on ingress traffic to balance cost and issue coverage.
- For asynchronous chains, enforce explicit consistency in context propagation.

## Continue Reading

- Start with the [Stelltrace product overview](/products/stelltrace/)
- Previous: [Getting Started](/products/stelltrace/quick-start)
- Next: [API Reference](/products/stelltrace/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stelltrace/configuration)
