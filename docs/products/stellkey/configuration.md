---
title: Stellkey Configuration Guide
outline: deep
---

# Stellkey · Configuration Guide

> A secret-management center for key lifecycle management, secret distribution, rotation auditing, and secure access control.

## Baseline Configuration

- Use different rotation cycles for highly sensitive secrets and ordinary secrets.
- Applications should keep only references and short-lived cache, avoiding long-lived plain-text persistence.

## Production Guidance

- Approval flows should be layered by secret class, environment, and access source.
- All sensitive values should be passed by reference to avoid copy-based spread.

## Continue Reading

- Start with the [Stellkey product overview](/products/stellkey/)
- Previous: [Getting Started](/products/stellkey/quick-start)
- Next: [API Reference](/products/stellkey/api-and-sdk)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellkey/configuration)
