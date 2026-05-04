---
title: Stellgate Configuration Guide
outline: deep
---

# Stellgate · Configuration Guide

> A gateway ingress platform for unified access, authentication, protocol translation, traffic governance, and API exposure.

## Baseline Configuration

- Routes should be layered by domain, path, tenant, and version.
- Security plugins and observability plugins are good candidates for a global default chain.

## Production Guidance

- Use separate certificates and plugin strategies for public-facing and internal gateways.
- Long-lived connection protocols need dedicated timeout, retry, and connection-limit tuning.

## Continue Reading

- Start with the [Stellgate product overview](/products/stellgate/)
- Previous: [Getting Started](/products/stellgate/quick-start)
- Next: [API Reference](/products/stellgate/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellgate/configuration)
