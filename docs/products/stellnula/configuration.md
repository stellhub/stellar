---
title: Stellnula Configuration Guide
outline: deep
---

# Stellnula · Configuration Guide

> A configuration center responsible for centralized storage, version management, progressive release, and dynamic distribution.

## Baseline Configuration

- Manage core configuration through layered categories such as `base configuration / secret references / environment differences`.
- Store only reference identifiers for sensitive values and let the secret-management center handle decryption.

## Production Guidance

- Keep configuration names stable and traceable to avoid ambiguity across teams.
- During progressive rollout, validate with a small traffic segment first and expand gradually.

## Continue Reading

- Start with the [Stellnula product overview](/products/stellnula/)
- Previous: [Getting Started](/products/stellnula/quick-start)
- Next: [API Reference](/products/stellnula/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellnula/configuration)
