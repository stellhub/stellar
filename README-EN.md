# Stellar Axis

`Stellar Axis` is a self-developed microservice middleware system. Its official Chinese brand name is `两仪`. The system uses the stellar axis as its engineering metaphor, the Bagua as its structural order, and the Spirit Axis as its AI expansion layer for high-concurrency, high-availability, and governable distributed systems.

This repository currently focuses on four areas:

- overall brand and system positioning
- component naming and role mapping
- repository layout and module boundaries
- dependency constraints and architecture decisions

## Positioning

Stellar Axis is not a loose collection of unrelated middleware pieces. It is a unified system with consistent naming, repository layout, and dependency boundaries.

The system follows these rules:

- `Stellar Axis` is the only official global English brand
- core components are presented as `English product name + Chinese Bagua name`
- repositories, modules, packages, and coordinates use English engineering names only
- cross-product collaboration happens through APIs, SDKs, events, and the control plane

Dual-stack positioning:

- Java framework: `Stellar Steel`
- Go framework: `Stellar Titan`

## Core Component Matrix

| Domain | English Name | Chinese Name | Short Name | Responsibility |
| :--- | :--- | :--- | :--- | :--- |
| Service discovery | `StarMap` | `乾仪` | `乾` | service metadata, instance discovery, registration |
| Config center | `Nebula` | `坤仪` | `坤` | dynamic configuration and environment distribution |
| Tracing | `LightBeam` | `离鉴` | `离` | tracing, observability, request illumination |
| Governance | `Orbit` | `巽策` | `巽` | routing, traffic governance, rollout control |
| Limiting and circuit breaking | `Pulsar` | `艮闸` | `艮` | rate limiting, circuit breaking, overload protection |
| Distributed scheduling | `Chronos` | `震策` | `震` | scheduling and distributed job execution |
| Distributed locking | `Singularity` | `坎锁` | `坎` | exclusivity and concurrency control |
| Gateway | `EventHorizon` | `兑门` | `兑` | ingress, boundary traffic, protocol aggregation |
| Messaging and event flow | `TaiJi Flow` | `太极流` | `中枢` | messaging, event streaming, asynchronous decoupling |

Recommended presentation style:

- `StarMap · 乾`
- `Nebula · 坤`
- `Orbit · 巽`
- `TaiJi Flow · 太极流`

## AI Extension Layer

Beyond the Bagua governance matrix, the system defines an independent AI-native extension layer:

| English Name | Chinese Name | Responsibility |
| :--- | :--- | :--- |
| `Spirit Axis` | `灵觉层` | overall AI capability layer |
| `ZhongFu Engine` | `中孚` | reasoning and model understanding |
| `DaChu Memory` | `大畜` | memory and knowledge accumulation |
| `Sui Agent` | `随位` | agent orchestration and autonomous execution |
| `Xian MCP` | `咸` | context perception protocol |
| `Heng MCP` | `恒` | tool execution protocol |
| `Meng MCP` | `蒙` | prompt guidance protocol |

The Spirit Axis is intentionally separate from the eight Bagua positions. It is an intelligence layer built above the governance layer.

## Engineering Naming Rules

The following names are fixed:

- global brand: `Stellar Axis`
- Java aggregate repository: `stellar-steel`
- Go aggregate repository: `stellar-titan`
- control plane: `stellar-control-plane`
- CLI: `stellarctl`
- coordinate namespace: `io.stellar.axis`

The following transitional names are deprecated:

- `Qian-Directory`
- `Kun-Profile`
- `Li-Trace`
- `Xun-Pilot`
- `Gen-Brake`
- `Zhen-Ticker`
- `Kan-Vault`
- `Dui-Portal`

## Repository Layout

The system uses a three-level repository model:

Top-level aggregate repositories:

- `stellar-axis`
- `stellar-steel`
- `stellar-titan`
- `spirit-axis`
- `stellar-control-plane`
- `stellarctl`
- `stellar-examples`
- `stellar-deploy`

Core product repositories:

- `starmap`
- `nebula`
- `lightbeam`
- `orbit`
- `pulsar`
- `chronos`
- `singularity`
- `event-horizon`
- `taiji-flow`

Unified product repository skeleton:

```text
{product}/
├── docs/
├── api/
├── server/
├── clients/
├── starters/
├── sidecars/
├── operators/
├── deploy/
├── test/
└── examples/
```

## Dependency Rules

The system uses strict one-way dependency rules.

Allowed patterns:

- depend on public APIs, SDKs, and stable protocols
- use the control plane, event flow, or official clients for cross-product collaboration
- `starter` depends on `clients/api`
- `sidecar` depends on `clients/api/runtime`
- `operator` depends on `api` and Kubernetes runtime

Forbidden patterns:

- `starter -> server`
- `client -> server`
- `operator -> server`
- `control-plane -> product/internal`
- `productA/server -> productB/server`
- `examples -> product/internal`

## Documentation

- [Naming Proposal](./Proposal.md)
- [Naming Convention](./Naming-Convention.md)
- [Repository Layout](./Repo-Layout.md)
- [Module Dependency](./Module-Dependency.md)
- [Architecture Decision Record](./Architecture-Decision-Record.md)

## Motto

> Establish the foundation with Qian and Kun, observe the system through Xun and Li, guard order with Gen and Kan, activate flow through Zhen and Dui.

## License

Apache License 2.0
