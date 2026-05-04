---
title: Stellmap System Architecture
outline: deep
---

# Stellmap · System Architecture

> A registry center that manages service instance registration, health checks, discovery subscriptions, and topology change propagation.

## Naming Context

`StellMap` means a star map or coordinate chart within the Stell Hub product family.
That name reflects the role of a registry center: every service instance is like a point on a map, and the registry keeps track of its current location and state so callers can discover, locate, and navigate across the system.

## Implementation Scope

- This repository carries the Go implementation of StellMap.
- It is intended to evolve around service registration, service discovery, instance heartbeat and health state management, service metadata governance, namespace and grouping isolation, and integration contracts with governance, configuration, and control-plane modules.

## Engineering Goals

- Consistency first: a CP-oriented design that favors linearizable correctness over apparent availability.
- Lightweight operation: a single Raft group can carry the registry data plane without unnecessary coordination layers.
- High availability: service remains available to clients as long as a quorum survives.
- High concurrency: read and write paths stay short to avoid accidental bottlenecks.
- Fast recovery: crash recovery remains explicit, with clear ownership across WAL, storage, and snapshots.
- Future evolution: the design can grow into watch, lease, compression, and layered cache capabilities without reworking the core.

## Runtime Topology

The current `stellmapd` process exposes three listening surfaces: public HTTP, dedicated admin HTTP, and internal gRPC.
Together they form a replicated registry service where external clients use HTTP for registration and discovery, while nodes use gRPC for replication, snapshot transfer, and leader-oriented internal traffic.

## Consensus

- The registry cluster runs as a single Raft group built on `etcd-io/raft`.
- The core object is the service instance record, and its metadata scale is typically much smaller than a general-purpose database.
- A single consensus group greatly reduces implementation and operational complexity.
- For a registry center, consistency and maintainability are usually more important than horizontal sharding.

## Domain Model

- The external model is an instance registry rather than a general-purpose key-value product.
- Its logical primary key is `namespace / service / instanceId`.
- Instance records store endpoints, labels, metadata, lease TTL, and recent heartbeat information.
- Normalized service names and structured fields are retained together to support prefix subscription, permission governance, and aggregated monitoring.

## Consistency Guarantees

- StellMap explicitly adopts a CP architecture.
- During a network partition, minority nodes stop serving linearizable writes.
- The cluster prefers no divergence, no rollback, and no stale values over temporary write availability.
- All reads are linearizable by default, implemented through `ReadIndex` and local state-machine reads once `appliedIndex >= readIndex`.

## Membership Management

- Node membership changes support `Learner` and joint consensus.
- New nodes first join as learners, receive logs without voting, and are promoted only after they catch up and pass health checks.
- Large topology updates rely on `ConfChangeV2` and joint consensus to avoid quorum instability during transitions.

## Storage Layout

- `WAL` stores the Raft log and consensus metadata such as `Entry` and `HardState`.
- `Pebble` stores applied registry data plus a small amount of local metadata such as apply progress, snapshot markers, and control watermarks.
- `Snapshot` files are managed separately for recovery, validation, and atomic replacement.

## Storage Engine Choice

- Pebble is a Go-native LSM key-value engine that avoids `cgo` and fits frequent registration and heartbeat updates well.
- It supports range delete, snapshots, and batched writes, which are useful for instance cleanup, snapshotting, and batch apply.
- Its operational maturity and low maintenance burden make it a pragmatic fit for registry-state storage.

## Verification Strategy

- Core validation must include linearizable read testing, write-after-read guarantees, and candidate-set consistency checks.
- Membership testing must cover learner join, promotion, and joint consensus transitions.
- Failure testing must cover leader crashes, follower restart, partition tolerance, WAL corruption, snapshot corruption, and interrupted snapshot transfer.
- Load testing must cover high-frequency registration, heartbeat renewal, discovery queries, watch streams, and long-running soak scenarios.

## Module Map

- `raftnode` drives the replicated state machine and coordinates `Ready`, `ReadIndex`, and membership changes.
- `wal` persists the Raft log and manages segment rotation, sync, and recovery.
- `snapshot` exports, installs, validates, and cleans up snapshots.
- `storage` implements the state machine on top of Pebble and guarantees ordered, idempotent apply.
- `transport` separates external HTTP access from internal gRPC replication.

## Internal Protocol

- Internal replication is defined by Raft and snapshot protobuf contracts and is never exposed to business clients directly.
- `RaftTransport` batches ordinary Raft messages, while `SnapshotService` handles streamed snapshot upload and download.
- The gRPC layer converts protobuf payloads into local runtime models, and the runtime transport service executes the real message flow.

## Admin Control Surface

- Cluster membership changes, leader transfer, and cluster status inspection are all triggered through `stellmapctl`, even though execution still lands on the dedicated admin HTTP listener.
- Public business clients only reach `HTTPAddr`, while `stellmapctl` targets the local `AdminAddr` by default.
- Admin requests currently require both localhost source access and `Authorization: Bearer <token>` authentication.

## Continue Reading

- Start with the [Stellmap product overview](/products/stellmap/)
- Previous: [Design Overview](/products/stellmap/summary-design)
- Next: [Deployment Model](/products/stellmap/deployment)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellmap/architecture)
