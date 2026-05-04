---
title: Stellgate Design Overview
outline: deep
---

# Stellgate · Design Overview

> A gateway ingress platform for unified access, authentication, protocol translation, traffic governance, and API exposure.

## Control Plane

- The control plane manages routes, certificates, plugins, and release strategy.
- Gateway instances, domains, upstream services, and plugin chains are orchestrated centrally.

## Data Plane

- The data plane uses stateless forwarding nodes that scale horizontally.
- An extension layer handles authentication, observability, transformation, and traffic control through plugin chains.

## Continue Reading

- Start with the [Stellgate product overview](/products/stellgate/)
- Next: [System Architecture](/products/stellgate/architecture)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellgate/summary-design)
