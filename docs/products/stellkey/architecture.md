---
title: Stellkey System Architecture
outline: deep
---

# Stellkey · System Architecture

> A secret-management center for key lifecycle management, secret distribution, rotation auditing, and secure access control.

## Component Model

- Secret Store: stores secrets, certificates, and version records.
- Rotation Engine: performs rotation, revocation, and expiry reminder.
- Access Gateway: centralizes authentication, audit, and distribution control.

## Interaction Flow

- Applications obtain short-lived credentials or secret references through the access gateway.
- The rotation engine triggers rotation by policy and synchronizes version state updates.

## Continue Reading

- Start with the [Stellkey product overview](/products/stellkey/)
- Previous: [Design Overview](/products/stellkey/summary-design)
- Next: [Deployment Model](/products/stellkey/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellkey/architecture)
