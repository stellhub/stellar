---
title: Stellmap API Reference
outline: deep
---

# Stellmap · API Reference

> A registry center that manages service instance registration, health checks, discovery subscriptions, and topology change propagation.

## Write Contract

1. A client sends a registration, deregistration, or heartbeat renewal request through the public HTTP API.
2. If the request lands on a follower, it is redirected or forwarded to the leader.
3. The leader encodes the instance change into a proposal and submits it to the single Raft group.
4. The log is appended to the local WAL.
5. The log is replicated to a quorum and committed.
6. The state machine applies the change to Pebble in order.
7. The write success response is returned.

## Read Contract

1. A client sends an instance query request.
2. The request reaches the leader, or a follower forwards it to the leader path.
3. The leader executes `ReadIndex`.
4. The node waits until local `appliedIndex` catches up to `readIndex`.
5. The latest registry view is read from Pebble.
6. A linearizable response is returned.

## Public API Contract

- Third-party integration is exposed through the public HTTP API only.
- That keeps onboarding simple for scripts, sidecars, and multi-language clients while preserving a clear boundary for registration, discovery, query, and health-reporting interfaces.

## Admin API Contract

- Membership changes, leader transfer, and cluster-status inspection run on a dedicated admin HTTP listener instead of the public business listener.
- The public HTTP surface carries only business data-plane and health endpoints.
- The admin listener is bound to loopback by default, currently accepts only localhost traffic, and requires `Authorization: Bearer <token>` authentication.
- `stellmapctl` is the only supported control-plane entry point today.

## Internal Transport Contract

- Cluster-internal replication uniformly uses gRPC.
- The external API and internal replication surface should not share protocol assumptions because the internal channel carries Raft messages, snapshots, and leader-forwarded read or write traffic.
- External HTTP and internal gRPC still share the same core state-machine and permission-validation logic to avoid double implementation.

## Continue Reading

- Start with the [Stellmap product overview](/products/stellmap/)
- Previous: [Configuration Guide](/products/stellmap/configuration)
- Next: [Observability Guide](/products/stellmap/observability)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellmap/api-and-sdk)
