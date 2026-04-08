# Architecture Decision Record

## 1. Problem Analysis

The current system already has clear naming, layout, and dependency decisions, but those decisions are spread across multiple documents. Without a single decision record, the team may still drift back to transitional naming, inconsistent repository structures, or invalid dependencies.

This document consolidates the architecture decisions that define the current form of the system.

It serves as the authoritative summary of:

- brand decisions
- naming decisions
- repository structure decisions
- dependency boundary decisions

Related documents:

- [Proposal.md](./Proposal.md)
- [Naming-Convention.md](./Naming-Convention.md)
- [Repo-Layout.md](./Repo-Layout.md)
- [Module-Dependency.md](./Module-Dependency.md)

## 2. Design

This ADR follows four decision groups:

- ADR-001: Global brand
- ADR-002: Product naming
- ADR-003: Repository topology
- ADR-004: Dependency boundaries

Each decision includes:

- context
- decision
- consequence

## 3. Implementation

### ADR-001: Global Brand

#### Context

The project previously mixed multiple top-level expressions such as `LiangYi`, `两仪`, and `Stellar Axis`. That creates ambiguity across documentation, repository names, and future ecosystem growth.

#### Decision

The official global brand is fixed as:

- English: `Stellar Axis`
- Chinese: `两仪`
- first combined presentation: `Stellar Axis（两仪）`

`LiangYi` is retained only as a phonetic reference, not as the primary English-facing brand.

#### Consequence

- English-facing documents, repository namespaces, and coordinates use `Stellar Axis`
- Chinese-facing materials use `两仪`
- mixed transitions such as `LiangYi Framework` are deprecated

### ADR-002: Core Product Naming

#### Context

The system needs names that work both as technical products and as part of a coherent philosophical system.

#### Decision

Every core infrastructure product uses:

> `English product name + Chinese Bagua name`

Final component set:

| Domain | English Name | Chinese Name | Short Name |
| :--- | :--- | :--- | :--- |
| Registry | `StarMap` | `乾仪` | `乾` |
| Config | `Nebula` | `坤仪` | `坤` |
| Tracing | `LightBeam` | `离鉴` | `离` |
| Governance | `Orbit` | `巽策` | `巽` |
| Limiter | `Pulsar` | `艮闸` | `艮` |
| Scheduler | `Chronos` | `震策` | `震` |
| Locking | `Singularity` | `坎锁` | `坎` |
| Gateway | `EventHorizon` | `兑门` | `兑` |
| Messaging | `TaiJi Flow` | `太极流` | `中枢` |

The AI extension layer is defined separately:

| English Name | Chinese Name |
| :--- | :--- |
| `Spirit Axis` | `灵觉层` |
| `ZhongFu Engine` | `中孚` |
| `DaChu Memory` | `大畜` |
| `Sui Agent` | `随位` |
| `Xian MCP` | `咸` |
| `Heng MCP` | `恒` |
| `Meng MCP` | `蒙` |

#### Consequence

- transitional hybrid names such as `Qian-Directory` and `Kun-Profile` are deprecated
- documentation and product cards use the dual-name convention
- repositories and engineering modules use English engineering names only

### ADR-003: Repository Topology

#### Context

Without a stable repository topology, different products will evolve incompatible structures and ownership boundaries.

#### Decision

The system uses a three-level repository model.

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

Each product repository follows the same skeleton:

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

Conditional modules are created only when needed:

- `starters/`
- `sidecars/`
- `operators/`
- `examples/`

#### Consequence

- each core product remains independently evolvable
- Java and Go have stable aggregate entry points
- platform and runtime concerns remain separated

### ADR-004: Dependency Boundaries

#### Context

Repository separation alone is not enough. Clear dependency rules are required to prevent architectural erosion.

#### Decision

The system uses strict one-way dependency rules.

Allowed dependency principles:

- depend on public APIs, SDKs, and stable contracts
- use events, public clients, and control-plane APIs for cross-product collaboration
- keep internal implementation private

Forbidden dependency principles:

- `starter -> server`
- `client -> server`
- `operator -> server`
- `control-plane -> product/internal`
- `productA/server -> productB/server`
- `examples -> product/internal`

Repository-level rules:

- `stellar-axis` has no runtime code dependencies
- `stellar-steel` and `stellar-titan` depend on public contracts only
- product repositories may depend on their own `api/` and shared runtime abstractions
- the control plane depends on public management surfaces, not internal implementations

#### Consequence

- internal implementation details remain encapsulated
- products can evolve independently
- cross-product integration remains contract-driven
- future splitting, replacement, or versioning becomes manageable

## 4. Complete Code

The current architecture is defined by the following fixed decisions:

```md
# Architecture Decision Record

## Fixed decisions

1. The only official global brand is `Stellar Axis（两仪）`.
2. Core products use `English product name + Chinese Bagua name`.
3. Repositories, modules, and dependency coordinates use English engineering names only.
4. The repository model is split into aggregate repositories, core product repositories, and supporting repositories.
5. Cross-product collaboration must happen through public APIs, SDKs, events, or control-plane interfaces.
6. Direct dependencies on internal implementation are forbidden across repository boundaries.
```
