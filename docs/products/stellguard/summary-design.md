---
title: Stellguard Design Overview
outline: deep
---

# Stellguard · Design Overview

> A zero-trust security platform for identity verification, access control, policy evaluation, and secure service-to-service communication.

## Control Plane

- The identity layer handles principal registration, credential issuance, and identity mapping.
- The policy layer makes context-based authorization decisions and risk assessments.

## Data Plane

- The enforcement layer applies interception uniformly across gateways, sidecars, and SDKs.
- Service-to-service communication is validated through mTLS, tokens, and context attributes.

## Continue Reading

- Start with the [Stellguard product overview](/products/stellguard/)
- Next: [System Architecture](/products/stellguard/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellguard/summary-design)
