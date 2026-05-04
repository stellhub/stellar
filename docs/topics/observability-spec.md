---
title: Observability Specification
category: Observability
summary: A baseline observability specification covering signals, naming, and operational expectations across infrastructure and application layers.
tags:
  - Observability
  - Tracing
  - Metrics
  - Logging
readingDirection: Read this when standardizing telemetry conventions or defining platform-wide observability contracts.
outline: deep
---

# Observability Specification

## Overview

A baseline observability specification covering signals, naming, and operational expectations across infrastructure and application layers.

## Abstract

An enterprise observability specification should not start by inventing its own vocabulary. It should start from existing standards, especially OpenTelemetry and Kubernetes, and only introduce custom fields where those standards are insufficient. That is the central position of this article: standard-first, fact-reuse first, and minimal custom supplementation.

The goal is to make logs, traces, and metrics align across platforms without building a second parallel semantics layer that future teams must constantly translate.

## 1. Why a Standard-First Spec Matters

Many organizations begin observability standardization by defining custom environment variables, custom request headers, and custom metric names. This feels fast in the short term, but over time it creates:

- duplicated semantics
- inconsistent field meaning across systems
- long-lived mapping layers between collection and analysis systems
- weaker ecosystem compatibility

A better approach is:

- OpenTelemetry as the primary semantic standard
- Kubernetes metadata as the primary source of runtime facts
- `STELLAR_*` and `X-Stellar-*` only as a minimal compatibility layer

## 2. Goals and Baseline

This specification is meant to align:

- business services
- middleware SDKs
- HTTP and gRPC gateways
- agents, sidecars, and collectors
- logs, traces, metrics, and alerting platforms
- deployment and platform injection systems

Its goals are:

1. unify resource identity
2. unify log key-value perspectives
3. unify propagation semantics
4. unify metric naming and attribute boundaries
5. clarify responsibility by platform layer

The baseline should come from OpenTelemetry resource semantics, service semantics, HTTP and RPC conventions, logs data model, W3C Trace Context, W3C Baggage, and Kubernetes metadata.

## 3. Resource Semantics Model

### 3.1 Core Principles

Resource identity should follow these rules:

1. prefer standard OpenTelemetry resource attributes
2. reuse Kubernetes metadata as the source of truth
3. treat `STELLAR_*` only as an adaptation layer
4. avoid multiple primary definitions for the same fact

### 3.2 Recommended Identity Fields

The most important service identity fields are:

- `service.name`
- `service.namespace`
- `service.version`
- `service.instance.id`
- `deployment.environment.name`

These should remain semantically stable. For example:

- `service.name` is not a pod name
- `service.namespace` is not the same thing as `k8s.namespace.name`
- `service.version` is not a branch name
- `service.instance.id` is not just a host IP

### 3.3 Infrastructure and Kubernetes Fields

Infrastructure semantics should remain distinct:

- `host.name` means host identity
- `k8s.node.name` means Kubernetes node identity
- `k8s.pod.ip` is a network address, not an instance identity

This distinction matters because observability data becomes misleading when topology fields are overloaded.

## 4. Log KV Perspective Model

OpenTelemetry defines a logs data model, but teams still need platform-level guidance on what should appear in client and server logs.

### 4.1 Common Log Fields

Every structured log entry should ideally expose:

- timestamp
- observed timestamp
- severity
- body
- trace correlation fields such as `TraceId` and `SpanId`
- service identity
- environment identity
- stable platform metadata such as pod and node references

### 4.2 Client-Side Log Perspective

A client log describes what the caller observed during an outbound request. It should focus on:

- caller service identity
- method or operation
- target address
- protocol attributes
- response status
- low-cardinality error classification
- authenticated business context where appropriate

Client logs should not pretend to know the callee's internal pod or process identity unless that fact is truly observed and propagated.

### 4.3 Server-Side Log Perspective

A server log describes what the receiving service actually handled. It should focus on:

- the receiving service identity
- inbound request metadata
- peer information if available
- response classification
- local decision context
- server-side error semantics

The key rule is simple: the client records what it sees as a caller, and the server records what it sees as a handler. Blurring the two produces confusing telemetry.

## 5. Context Propagation Model

The recommended default is:

- W3C Trace Context for trace propagation
- W3C Baggage for business metadata propagation
- `X-Stellar-*` only for a very small compatibility surface

For HTTP and gRPC, propagation should remain consistent enough that logs, traces, and metrics can be correlated without per-team conventions.

The practical principle is that platform-specific headers should not compete with standard propagation formats when the standard already works.

## 6. Metrics Model

Metrics should also remain standard-first.

### 6.1 General Rules

- reuse OpenTelemetry semantic conventions whenever possible
- keep labels governable
- avoid high-cardinality explosions
- separate client-side and server-side perspectives clearly

### 6.2 Client and Server Boundaries

Client metrics should describe outbound request behavior. Server metrics should describe request handling behavior. Mixing the two leads to incorrect SLO reasoning.

### 6.3 Other Middleware Scenarios

Databases, queues, and infrastructure middleware need the same discipline:

- standard resource identity
- explicit operation boundaries
- low-cardinality error dimensions
- careful label control

## 7. Layered Platform Responsibilities

A usable observability standard is not only a field list. It must define which layer owns what.

- release platform: injects stable runtime metadata
- gateways: enrich request context and preserve propagation correctness
- SDKs: emit standards-aligned telemetry
- collectors: perform translation, filtering, routing, and policy enforcement
- analysis platforms: consume standardized signals instead of inventing local semantics

Without explicit layering, every component starts "helping" in a different way and the telemetry model fragments again.

## 8. Acceptance Principles

An observability standard should be accepted only if it:

- improves consistency without fighting the ecosystem
- reduces duplicated naming
- makes machine correlation easier
- keeps custom extensions minimal
- remains understandable to both platform teams and service teams

## 9. Conclusion

The strongest observability standards are not the most customized ones. They are the ones that align with widely used semantics, reuse runtime facts that already exist, and constrain customization to the few places where it is genuinely necessary.

For enterprise observability, the right direction is not to out-invent OpenTelemetry. It is to adopt it rigorously, add only the smallest required compatibility layer, and keep logs, traces, and metrics anchored to one coherent model.

## Chinese Reference

- [Read the original Chinese article](/zh/topics/observability-spec)
