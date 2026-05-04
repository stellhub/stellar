---
title: Stellmap Deployment Model
outline: deep
---

# Stellmap · Deployment Model

> A registry center that manages service instance registration, health checks, discovery subscriptions, and topology change propagation.

## Deployment Topology

- Start with a three-node deployment in a single region and spread nodes across availability zones to reduce single-point risk.
- In multi-region scenarios, local discovery traffic can be served through read-only follower clusters.

## Availability Strategy

- Registry recovery should rely on both leases and replicated copies of the registry state.
- Clients should enable local caching and tolerate short-lived control-plane disconnects.

## Crash Recovery Flow

1. Read the latest local snapshot metadata.
2. If a valid snapshot exists, restore it into the state-machine working directory first.
3. Open `Pebble` and load registry data plus local metadata.
4. Open the `WAL` and recover `HardState`, `Entry`, and snapshot markers.
5. Discard log segments already covered by the snapshot index.
6. Apply the remaining committed entries to the state machine in order.
7. Rebuild the in-memory Raft node, apply watermarks, and membership view.
8. Only then transition the node into a serving state.

## Recovery Rules

- If the WAL is durable but the state machine has not applied the latest committed entries, replay the log at startup.
- If a snapshot file exists but metadata was not atomically switched, keep using the old snapshot.
- If Pebble reports an applied index behind the WAL commit index, continue applying entries until it catches up.
- If snapshot, WAL, and state-machine indexes disagree, prioritize the principle of never rolling back committed logs.

## Post-Recovery Checks

- Verify `HardState.Commit >= appliedIndex`.
- Verify the state machine's `appliedIndex` does not exceed the largest committed index in the WAL.
- Verify the snapshot `ConfState` matches the in-memory membership view after recovery.
- Verify whether the node is still part of the latest membership and whether it needs to fetch additional logs or install a newer snapshot from the leader.

## Snapshot and Log Compaction

- Trigger snapshot generation when `appliedIndex - snapshotIndex` crosses a configured threshold.
- Write snapshot metadata atomically after snapshot generation finishes.
- Retain only the WAL segments that are still necessary after the snapshot point.
- Let lagging nodes catch up by log replay first and switch to snapshot installation only when the gap becomes too large.

## Continue Reading

- Start with the [Stellmap product overview](/products/stellmap/)
- Previous: [System Architecture](/products/stellmap/architecture)
- Next: [Getting Started](/products/stellmap/quick-start)

## Chinese Source

- [Read the original Chinese page](/zh/products/stellmap/deployment)
