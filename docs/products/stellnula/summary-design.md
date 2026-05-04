---
title: Stellnula Design Overview
outline: deep
---

# Stellnula · Design Overview

> A configuration center responsible for centralized storage, version management, progressive release, and dynamic distribution.

## Control Plane

- The control plane handles configuration orchestration, approval flows, and release-plan management.
- Configuration layering and environment isolation are defined centrally.

## Data Plane

- The storage layer keeps both text and structured configuration in a versioned model.
- Clients consume incremental updates and local disaster-recovery fallbacks through watch channels.

## Continue Reading

- Start with the [Stellnula product overview](/products/stellnula/)
- Next: [System Architecture](/products/stellnula/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellnula/summary-design)
