# Stellar Axis

`Stellar Axis` is the official global brand of the system. This repository preserves the final decisions for the `Stellar Axis` ecosystem and serves as the single authoritative source for branding, product naming, engineering naming, repository layout, and dependency boundaries.

Only two final documents are retained in this repository:

- Chinese version: `README.md`
- English version: `README-EN.md`

All proposal documents, roadmap drafts, bootstrap plans, and other process-oriented materials have been consolidated and removed.

## Final Brand Decisions

- Official English brand: `Stellar Axis`
- Official Chinese brand: `星轴`
- Standard first presentation: `Stellar Axis（星轴）`

Dual-stack positioning:

- Java framework: `Stellar Core（星核）`
- Go framework: `Stellar Pulse（星脉）`

Usage rules:

- Use `Stellar Axis（星轴）` the first time it appears in README files, whitepapers, or solution documents
- `星轴` may be used directly in Chinese contexts
- `Stellar Axis` may be used directly in English contexts
- Legacy philosophical aliases, pinyin branding, and mixed naming are no longer retained

## Core Product Naming

All core products follow this unified display rule:

> `English official name · Chinese cosmic name`

### Core Middleware Matrix

| Domain | English Name | Chinese Name |
| :--- | :--- | :--- |
| Service Registry | `StarMap` | `星图` |
| Configuration Center | `Nebula` | `星云` |
| Distributed Tracing | `StarTrace` | `星迹` |
| Service Governance | `Orbit` | `星轨` |
| Rate Limiting and Circuit Breaking | `Pulsar` | `脉冲` |
| Scheduling | `Astrolabe` | `星盘` |
| Distributed Locking | `Singularity` | `奇点` |
| Gateway | `EventHorizon` | `视界` |
| Messaging | `CometFlow` | `彗流` |
| Metrics Platform | `Constellation` | `星座` |
| Alerting Platform | `NovaSignal` | `星讯` |
| Zero Trust Platform | `StarShield` | `星盾` |
| Key Management | `StarKey` | `星钥` |

Recommended display examples:

- `StarMap · 星图`
- `Nebula · 星云`
- `StarTrace · 星迹`
- `Orbit · 星轨`
- `CometFlow · 彗流`

### AI Astral Layer Matrix

| English Name | Chinese Name | Responsibility |
| :--- | :--- | :--- |
| `Astral Layer` | `星穹层` | Overall AI capability brand |
| `Quasar Engine` | `类星引擎` | Model inference and understanding |
| `StarVault Memory` | `星库` | Memory and knowledge accumulation |
| `Orbit Agent` | `轨使` | Agent orchestration and autonomous execution |
| `Sensor MCP` | `星感` | Context sensing protocol |
| `Vector MCP` | `星行` | Tool execution protocol |
| `GuideStar MCP` | `导星` | Prompt guidance protocol |

## Engineering Naming Rules

Engineering entities use English names only and do not carry Chinese display names.

This applies to:

- Git repository names
- Maven `artifactId`
- Go modules
- Java package names
- CLI tool names
- Submodules such as `starter`, `client`, `sdk`, `sidecar`, and `operator`

Unified rules:

- Use lowercase English words only
- Use hyphens `-` between words
- Do not use Chinese
- Do not use pinyin
- Do not use mixed Chinese-English naming

### Fixed Engineering Names

- Java aggregate repository: `stellar-core`
- Go aggregate repository: `stellar-pulse`
- AI aggregate repository: `astral-layer`
- Control plane: `stellar-control-plane`
- CLI: `stellarctl`
- Maven root namespace: `io.stellar.axis`
- Go namespace root: `github.com/stellar-axis/`

### Engineering Naming Pattern

Recommended pattern:

> `{product}-{capability}-{role}`

Common examples:

- `starmap-client-spring-boot-starter`
- `nebula-client-spring-boot-starter`
- `orbit-governance-starter`
- `comet-flow-client`

### Final OTel SDK Naming

OpenTelemetry-related SDKs follow this final pattern:

> `{product}-otel-{signal}-{lang}-sdk`

Where:

- `{product}` is the product ownership, such as `startrace`
- `{signal}` is the signal type, such as `logs`, `traces`, or `metrics`
- `{lang}` is the language dimension, such as `java` or `go`

The finalized names are:

- `startrace-otel-logs-java-sdk`
- `startrace-otel-logs-go-sdk`

Coordinate examples:

- `io.stellar.axis:startrace-otel-logs-java-sdk`
- `io.stellar.axis:startrace-otel-logs-go-sdk`

Future extension examples:

- `startrace-otel-traces-java-sdk`
- `startrace-otel-metrics-java-sdk`
- `startrace-otel-traces-go-sdk`
- `startrace-otel-metrics-go-sdk`

If a language-level aggregation module is needed in the future, the following names may be used:

- `startrace-otel-java-sdk`
- `startrace-otel-go-sdk`

These aggregate names are reserved for umbrella modules and should not be used for the concrete `logs` SDK.

### Java Package Naming Examples

- `io.stellar.axis.starmap.client`
- `io.stellar.axis.nebula.config`
- `io.stellar.axis.orbit.router`
- `io.stellar.axis.pulsar.limiter`
- `io.stellar.axis.cometflow.producer`

## Gateway Family Naming

The gateway family uses `EventHorizon · 视界` as the unified family name, with secondary qualifiers for specific responsibilities.

| Scenario | English Name | Chinese Name |
| :--- | :--- | :--- |
| Internal Gateway | `EventHorizon Internal` | `视界·内域` |
| External Gateway | `EventHorizon External` | `视界·外域` |
| LLM Gateway | `EventHorizon LLM` | `视界·智域` |
| Edge Gateway | `EventHorizon Edge` | `视界·边域` |

Recommended repository names:

- `event-horizon-internal`
- `event-horizon-external`
- `event-horizon-llm`
- `event-horizon-edge`

## Final Repository Layout

The ecosystem uses a three-layer structure: aggregate repositories, core product repositories, and supporting repositories.

### Top-Level Aggregate Repositories

- `stellar-axis`
- `stellar-core`
- `stellar-pulse`
- `astral-layer`
- `stellar-control-plane`
- `stellarctl`
- `stellar-examples`
- `stellar-deploy`

### Core Product Repositories

- `starmap`
- `nebula`
- `startrace`
- `orbit`
- `pulsar`
- `astrolabe`
- `singularity`
- `event-horizon`
- `comet-flow`

### Unified Product Repository Skeleton

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

Module constraints:

- `docs/`, `api/`, `server/`, `clients/`, `deploy/`, and `test/` form the core skeleton
- `starters/` is created only when Java integration is provided
- `sidecars/` is created only when a sidecar runtime exists
- `operators/` is created only when Kubernetes lifecycle management is needed
- `examples/` is created only when independent minimal demos are needed

## Final Dependency Boundaries

The entire ecosystem follows one-way dependency constraints.

Allowed principles:

- Depend on public APIs, SDKs, and stable contracts
- Use control-plane APIs, event streams, and official clients for cross-product collaboration
- `starter` depends on `clients/api`
- `sidecar` depends on `clients/api/runtime`
- `operator` depends on `api` and Kubernetes Runtime

Forbidden principles:

- `starter -> server`
- `client -> server`
- `operator -> server`
- `control-plane -> product/internal`
- `productA/server -> productB/server`
- `examples -> product/internal`

Additional constraints:

- `stellar-axis` is for standards and entry only, not runtime code
- `stellar-core` and `stellar-pulse` may depend only on public contracts, not concrete server implementations
- Core products may collaborate only through APIs, SDKs, events, configuration, or the control plane
- The control plane may depend only on public management surfaces, not internal product implementations

## Final Summary

This repository has frozen the following final decisions:

1. The only official global brand is `Stellar Axis（星轴）`
2. Core products use the unified pattern `English official name · Chinese cosmic name`
3. Repositories, modules, package names, and dependency coordinates use English engineering names only
4. Repository topology is fixed as aggregate repositories, core product repositories, and supporting repositories
5. Cross-product collaboration must happen through public APIs, SDKs, events, and control-plane interfaces
6. The finalized OpenTelemetry Logs SDK names are `startrace-otel-logs-java-sdk` and `startrace-otel-logs-go-sdk`

## License

Apache License 2.0
