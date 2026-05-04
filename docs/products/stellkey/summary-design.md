---
title: Stellkey Design Overview
outline: deep
---

# Stellkey · Design Overview

> A secret-management center for key lifecycle management, secret distribution, rotation auditing, and secure access control.

## Control Plane

- The control layer manages approval flow, rotation plans, and secret policy.
- Access boundaries, rotation cycles, and distribution policies are governed centrally.

## Data Plane

- The storage layer handles secret encryption, version management, and access isolation.
- The distribution layer handles short-lived credential issuance, retrieval cache, and expiration invalidation.

## Continue Reading

- Start with the [Stellkey product overview](/products/stellkey/)
- Next: [System Architecture](/products/stellkey/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellkey/summary-design)
