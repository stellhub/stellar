---
title: Stellgate System Architecture
outline: deep
---

# Stellgate · System Architecture

> A gateway ingress platform for unified access, authentication, protocol translation, traffic governance, and API exposure.

## Component Model

- Gateway Admin: manages route, certificate, and plugin configuration.
- Gateway Runtime: performs forwarding, load balancing, and traffic control.
- Plugin SDK: carries authentication, audit, and protocol-extension logic.

## Interaction Flow

- The control plane publishes route and plugin strategy to runtime gateway nodes.
- Requests are processed through the plugin chain before being forwarded to the target upstream service.

## Continue Reading

- Start with the [Stellgate product overview](/products/stellgate/)
- Previous: [Design Overview](/products/stellgate/summary-design)
- Next: [Deployment Model](/products/stellgate/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellgate/architecture)
