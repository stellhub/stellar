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

> `Stell-style name · original imagery`

### Core Middleware Matrix

| Domain | Original Imagery | **Stell-style Name** | Naming Rationale |
| :--- | :--- | :--- | :--- |
| **Service Registry** | StarMap (`星图`) | **Stellmap** | Stell + Map, a coordinate map for locating services. |
| **Configuration Center** | Nebula (`星云`) | **Stellnula** | A shortened blend of Stell and Nebula, suggesting an omnipresent cloud of config. |
| **Distributed Tracing** | StarTrace (`星迹`) | **Stelltrace** | Stell + Trace, following requests across the system like interstellar routes. |
| **Service Governance** | Orbit (`星轨`) | **Stellorbit** | Stell + Orbit, keeping services running on predictable paths. |
| **Rate Limiting and Circuit Breaking** | Pulsar (`脉冲`) | **Stellpulse** | Stell + Pulse, expressing precise control over traffic rhythm. |
| **Scheduling** | Astrolabe (`星盘`) | **Stellabe** | A shortened blend of Stell and Astrolabe, inspired by the classic celestial instrument. |
| **Distributed Locking** | Singularity (`奇点`) | **Stellpoint** | A singularity becomes a single point; Stell + Point emphasizes the unique execution point. |
| **Gateway** | EventHorizon (`视界`) | **Stellgate** | Gate is a natural gateway metaphor, forming a clear "stellar gate". |
| **Messaging** | CometFlow (`彗流`) | **Stellflow** | Stell + Flow, evoking an ordered data stream like a comet tail. |
| **Logging Platform** | Spectrum (`星谱`) | **Stellspec** | A shortened blend of Stell and Spectrum, focused on analyzing system "spectra". |
| **Metrics Platform** | Constellation (`星座`) | **Stellcon** | A shortened blend of Stell and Constellation, connecting metric points into constellations. |
| **Alerting Platform** | NovaSignal (`星讯`) | **Stellvox** | Vox means voice or message in Latin, giving the alerting plane a clearer signal metaphor. |
| **Zero Trust Platform** | StarShield (`星盾`) | **Stellguard** | Stell + Guard, a more active and dynamic defensive identity than Shield. |
| **Key Management** | StarKey (`星钥`) | **Stellkey** | Stell + Key, direct and strong for key management positioning. |

Recommended display examples:

- `Stellmap · 星图`
- `Stellnula · 星云`
- `Stelltrace · 星迹`
- `Stellspec · 星谱`
- `Stellorbit · 星轨`
- `Stellflow · 彗流`

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

- `stellmap-client-spring-boot-starter`
- `stellnula-client-spring-boot-starter`
- `stellorbit-governance-starter`
- `stellflow-client`

### Final Logging SDK Naming

`Stellspec` is the official product name for the logging platform. Its external SDKs use simplified names and do not include `otel` or signal-type markers in the primary SDK name.

The finalized names are:

- `stellspec-java-sdk`
- `stellspec-go-sdk`

Coordinate examples:

- `io.stellar.axis:stellspec-java-sdk`
- `io.stellar.axis:stellspec-go-sdk`

Naming semantics:

- `stellspec` is the product ownership for the logging platform
- `java` and `go` represent the language dimension
- `sdk` represents the engineering role

`OpenTelemetry` may still exist as an internal implementation or integration capability, but it is not part of the current external SDK naming.

### Java Package Naming Examples

- `io.stellar.axis.stellmap.client`
- `io.stellar.axis.stellnula.config`
- `io.stellar.axis.stellorbit.router`
- `io.stellar.axis.stellpulse.limiter`
- `io.stellar.axis.stellflow.producer`

## Gateway Family Naming

The gateway family uses `Stellgate · 视界` as the unified family name, with secondary qualifiers for specific responsibilities.

| Scenario | English Name | Chinese Name |
| :--- | :--- | :--- |
| Internal Gateway | `Stellgate Internal` | `视界·内域` |
| External Gateway | `Stellgate External` | `视界·外域` |
| LLM Gateway | `Stellgate LLM` | `视界·智域` |
| Edge Gateway | `Stellgate Edge` | `视界·边域` |

Recommended repository names:

- `stellgate-internal`
- `stellgate-external`
- `stellgate-llm`
- `stellgate-edge`

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

- `stellmap`
- `stellnula`
- `stelltrace`
- `stellorbit`
- `stellpulse`
- `stellabe`
- `stellpoint`
- `stellgate`
- `stellflow`
- `stellspec`
- `stellcon`
- `stellvox`
- `stellguard`
- `stellkey`

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
6. The finalized logging platform name is `Stellspec · 星谱`, and its SDK names are `stellspec-java-sdk` and `stellspec-go-sdk`

## License

Apache License 2.0
