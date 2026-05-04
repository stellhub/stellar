---
title: Stellnula System Architecture
outline: deep
---

# Stellnula · System Architecture

> A configuration center responsible for centralized storage, version management, progressive release, and dynamic distribution.

## Component Model

- Configuration Repository: stores configuration versions and metadata.
- Release Engine: handles canary release, phased rollout, and scheduled publishing.
- Client Agent: receives changes and persists local cache snapshots.

## Interaction Flow

- The release engine pushes incremental changes to clients according to the rollout plan.
- When remote fetch fails, the client falls back to local snapshots to keep startup stable.

## Continue Reading

- Start with the [Stellnula product overview](/products/stellnula/)
- Previous: [Design Overview](/products/stellnula/summary-design)
- Next: [Deployment Model](/products/stellnula/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellnula/architecture)
