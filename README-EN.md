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

| Domain | English Full Name | Chinese Full Name | Short Name | Responsibility |
| :--- | :--- | :--- | :--- | :--- |
| <nobr>Service discovery</nobr> | `StarMap` | `乾仪` | `乾` | service metadata, instance discovery, registration |
| <nobr>Config center</nobr> | `Nebula` | `坤仪` | `坤` | dynamic configuration and environment distribution |
| <nobr>Tracing</nobr> | `LightBeam` | `离鉴` | `离` | tracing, observability, request illumination |
| <nobr>Governance</nobr> | `Orbit` | `巽策` | `巽` | routing, traffic governance, rollout control |
| <nobr>Limiting and circuit breaking</nobr> | `Pulsar` | `艮闸` | `艮` | rate limiting, circuit breaking, overload protection |
| <nobr>Distributed scheduling</nobr> | `Chronos` | `震策` | `震` | scheduling and distributed job execution |
| <nobr>Distributed locking</nobr> | `Singularity` | `坎锁` | `坎` | exclusivity and concurrency control |
| <nobr>Gateway</nobr> | `EventHorizon` | `兑门` | `兑` | ingress, boundary traffic, protocol aggregation |
| <nobr>Messaging and event flow</nobr> | `TaiJi Flow` | `太极流` | `中枢` | messaging, event streaming, asynchronous decoupling |

Recommended presentation style:

- `StarMap · 乾`
- `Nebula · 坤`
- `Orbit · 巽`
- `TaiJi Flow · 太极流`

## Naming Rationale for Core Components

These names are not direct functional translations. Each one binds an engineering role to a stellar metaphor and a Bagua-derived Chinese concept. The English name carries technical branding, while the Chinese name carries structural and philosophical meaning.

### StarMap · 乾仪

- `StarMap` refers to a celestial map or coordinate chart. A service registry is effectively the coordinate system of the platform, allowing services to locate and discover one another.
- `乾仪` draws from the idea that `Qian` represents Heaven in the *I Ching*. Qian stands for primacy, order, and overarching direction, which matches the role of global service metadata.
- The character `仪` suggests a formal model, pattern, or governing frame, indicating that this is the system-wide reference standard rather than a single node.

### Nebula · 坤仪

- `Nebula` evokes the stellar nursery: a space that carries matter and gives rise to formation. That aligns with a config center, which carries environment state and nurtures runtime behavior.
- `坤仪` draws from `Kun`, which represents Earth in the *I Ching*. Kun is the principle of bearing, containing, and sustaining, which fits a configuration substrate.
- `仪` again signals a formal system of order rather than a random set of config values.

### LightBeam · 离鉴

- `LightBeam` highlights the act of illumination. Distributed tracing exists to make invisible call paths visible across the system.
- `离鉴` draws from `Li`, which in the *I Ching* is associated with fire, brightness, and visibility.
- `鉴` means mirror, reflection, or discernment. The point is not only to illuminate the system, but also to inspect and understand it clearly.

### Orbit · 巽策

- `Orbit` expresses governed movement along an established path. That directly fits routing, policy-driven traffic shaping, and service governance.
- `巽策` draws from `Xun`, associated with wind in the *I Ching*. Wind penetrates gently but thoroughly, which matches the way governance influences traffic behavior without brute force.
- `策` means strategy, planning, and policy. The name emphasizes rule-based guidance rather than static interception.

### Pulsar · 艮闸

- `Pulsar` suggests stable cadence and precise periodic control. That maps well to quotas, thresholds, and controlled frequency in rate limiting and circuit breaking.
- `艮闸` draws from `Gen`, associated with stillness and stopping in the *I Ching*.
- `闸` means gate or sluice. The name makes the engineering meaning explicit: it is a protective gate that stops overload and isolates failure.

### Chronos · 震策

- `Chronos` represents time order and temporal control, making it a natural fit for scheduling.
- `震策` draws from `Zhen`, associated with thunder and activation in the *I Ching*.
- `策` emphasizes orchestration and dispatch strategy, not just clock-based triggering.

### Singularity · 坎锁

- `Singularity` captures the idea of collapsing many concurrent contenders into one exclusive control point.
- `坎锁` draws from `Kan`, associated with water, danger, and constrained passage in the *I Ching*.
- `锁` states the engineering role directly: exclusive locking that turns disorderly contention into controlled access.

### EventHorizon · 兑门

- `EventHorizon` represents a boundary surface between inside and outside. A gateway is exactly that boundary for ingress traffic.
- `兑门` draws from `Dui`, associated with openness, speech, and the mouth or opening in the *I Ching*.
- `门` makes the meaning explicit: it is the formal gate through which external requests enter the system.

### TaiJi Flow · 太极流

- `TaiJi Flow` emphasizes flow as the medium of events, messages, and asynchronous movement across the platform.
- `太极流` draws from the classical line that Taiji gives rise to LiangYi. In this system, messaging is treated as the central medium that keeps all components connected.
- The name is intentionally central rather than peripheral, because the messaging layer is the circulation core of the overall architecture.

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

## Naming Rationale for the AI Layer

### Spirit Axis · 灵觉层

- `Spirit Axis` describes the axis of cognition, understanding, and intelligent response.
- `灵觉层` means an awareness layer rather than a single model endpoint. It frames AI as a structured layer of reasoning, memory, and execution.

### ZhongFu Engine · 中孚

- `ZhongFu` comes from Hexagram 61 of the *I Ching*, carrying the idea of inner trust, sincerity, and meaningful resonance.
- As a reasoning engine, the name suggests semantic understanding with internal coherence rather than shallow pattern matching.

### DaChu Memory · 大畜

- `DaChu` comes from Hexagram 26, associated with accumulation, reserve, and cultivated storage.
- That makes it a strong fit for memory, retrieval, and knowledge retention.

### Sui Agent · 随位

- `Sui` carries the sense of following conditions and adapting to circumstance.
- `位` adds the idea of role and position, so the name conveys adaptive execution within a defined operational role.

### Xian MCP · 咸

- `Xian` comes from Hexagram 31, centered on resonance and mutual sensing.
- That makes it appropriate for context perception and resource awareness.

### Heng MCP · 恒

- `Heng` comes from Hexagram 32, centered on constancy and continuity.
- It fits the tool execution protocol because reliable action must be stable, repeatable, and durable.

### Meng MCP · 蒙

- `Meng` comes from Hexagram 4, associated with instruction, awakening, and guided understanding.
- It fits prompt guidance because the role of prompting is to orient the model toward correct interpretation and action.

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
