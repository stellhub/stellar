---
title: Stellguard System Architecture
outline: deep
---

# Stellguard · System Architecture

> A zero-trust security platform for identity verification, access control, policy evaluation, and secure service-to-service communication.

## Component Model

- Identity Provider: manages identities, tokens, and certificates.
- Policy Engine: computes authorization strategy and risk decisions.
- Enforcement Point: enforces access control at gateways and services.

## Interaction Flow

- The identity provider issues credentials and shares identity context with the policy engine.
- The enforcement point intercepts requests in real time and returns a policy decision.

## Continue Reading

- Start with the [Stellguard product overview](/products/stellguard/)
- Previous: [Design Overview](/products/stellguard/summary-design)
- Next: [Deployment Model](/products/stellguard/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellguard/architecture)
