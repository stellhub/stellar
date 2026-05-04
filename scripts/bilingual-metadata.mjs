export const productDocPages = [
  {
    slug: "summary-design",
    titleEn: "Design Overview",
    titleZh: "概要设计",
    focusEn: "design goals, boundaries, and the main capability map",
    focusZh: "设计目标、边界与核心能力地图"
  },
  {
    slug: "architecture",
    titleEn: "System Architecture",
    titleZh: "架构组成",
    focusEn: "core modules, collaboration flows, and runtime responsibilities",
    focusZh: "核心模块、协作链路与运行时职责"
  },
  {
    slug: "deployment",
    titleEn: "Deployment Model",
    titleZh: "部署形态",
    focusEn: "deployment topology, dependencies, and environment planning",
    focusZh: "部署拓扑、依赖关系与环境规划"
  },
  {
    slug: "quick-start",
    titleEn: "Getting Started",
    titleZh: "快速入门",
    focusEn: "minimum setup, first-run steps, and validation points",
    focusZh: "最小化接入步骤、首次运行路径与验证点"
  },
  {
    slug: "configuration",
    titleEn: "Configuration Guide",
    titleZh: "配置建议",
    focusEn: "key parameters, rollout strategy, and configuration governance",
    focusZh: "关键参数、灰度发布策略与配置治理"
  },
  {
    slug: "api-and-sdk",
    titleEn: "API Reference",
    titleZh: "API 与 SDK",
    focusEn: "public contracts, integration surfaces, and client-side expectations",
    focusZh: "公开契约、集成面与客户端约束"
  },
  {
    slug: "observability",
    titleEn: "Observability Guide",
    titleZh: "可观测性",
    focusEn: "metrics, logs, traces, and operational feedback loops",
    focusZh: "指标、日志、链路与运维反馈闭环"
  }
];

export const products = [
  {
    slug: "stellmap",
    name: "Stellmap",
    chineseName: "星图",
    domainEn: "Registry Center",
    domainZh: "注册中心",
    originEn: "StarMap",
    originZh: "星图",
    rationaleEn: "Stell + Map, a service coordinate map for discovery and routing.",
    summaryEn:
      "Provides unified service registration and discovery for instance health, subscriptions, and topology change propagation.",
    positioningEn:
      "Stellmap keeps service endpoints stable across scaling, environment switches, and multi-cluster routing.",
    scenariosEn: ["Microservice instance discovery", "Multi-cluster traffic steering"]
  },
  {
    slug: "stellnula",
    name: "Stellnula",
    chineseName: "星云",
    domainEn: "Configuration Center",
    domainZh: "配置中心",
    originEn: "Nebula",
    originZh: "星云",
    rationaleEn: "A condensed form of Nebula, emphasizing widely distributed configuration presence.",
    summaryEn:
      "Provides centralized configuration storage, controlled rollout, dynamic distribution, and change-audit capabilities.",
    positioningEn:
      "Stellnula standardizes configuration lifecycle management across environments and services.",
    scenariosEn: ["Environment-aware configuration delivery", "Governed configuration rollout"]
  },
  {
    slug: "stelltrace",
    name: "Stelltrace",
    chineseName: "星迹",
    domainEn: "Distributed Tracing",
    domainZh: "链路追踪",
    originEn: "StarTrace",
    originZh: "星迹",
    rationaleEn: "Stell + Trace, tracking request paths across a distributed landscape.",
    summaryEn:
      "Provides trace collection, retrieval, and diagnostic analysis for distributed request paths.",
    positioningEn:
      "Stelltrace helps engineers reason about latency, dependency edges, and fault localization.",
    scenariosEn: ["Cross-service latency analysis", "Dependency troubleshooting"]
  },
  {
    slug: "stellorbit",
    name: "Stellorbit",
    chineseName: "星轨",
    domainEn: "Service Governance",
    domainZh: "服务治理",
    originEn: "Orbit",
    originZh: "星轨",
    rationaleEn: "Services stay on predictable operational trajectories, like bodies on an orbit.",
    summaryEn:
      "Provides runtime routing, load balancing, retry, and traffic-governance controls across service lifecycles.",
    positioningEn:
      "Stellorbit coordinates runtime governance policies without leaking them into business code.",
    scenariosEn: ["Progressive service rollout", "Runtime governance policy management"]
  },
  {
    slug: "stellpulse",
    name: "Stellpulse",
    chineseName: "脉冲",
    domainEn: "Rate Limiting and Circuit Breaking",
    domainZh: "流控熔断",
    originEn: "Pulsar",
    originZh: "脉冲",
    rationaleEn: "Stell + Pulse, focused on sensing and controlling live traffic rhythm.",
    summaryEn:
      "Provides rate limiting, circuit breaking, degradation, and hotspot protection for service availability.",
    positioningEn:
      "Stellpulse reduces blast radius during overload and abnormal dependency behavior.",
    scenariosEn: ["Runtime overload protection", "Dependency failure containment"]
  },
  {
    slug: "stellabe",
    name: "Stellabe",
    chineseName: "星盘",
    domainEn: "Job Scheduling",
    domainZh: "任务调度",
    originEn: "Astrolabe",
    originZh: "星盘",
    rationaleEn: "A shortened form of Astrolabe, referencing coordinated timing and positioning.",
    summaryEn:
      "Provides timed scheduling, workflow orchestration, sharded execution, and coordinated failure recovery.",
    positioningEn:
      "Stellabe manages time-based and workflow-driven execution with platform-level controls.",
    scenariosEn: ["Batch task orchestration", "Recovery-aware scheduled execution"]
  },
  {
    slug: "stellpoint",
    name: "Stellpoint",
    chineseName: "奇点",
    domainEn: "Distributed Lock",
    domainZh: "分布式锁",
    originEn: "Singularity",
    originZh: "奇点",
    rationaleEn: "A single decisive point, fitting exclusive ownership and coordination semantics.",
    summaryEn:
      "Provides lease locks, leader election, and coordination primitives for mutually exclusive access to critical resources.",
    positioningEn:
      "Stellpoint provides deterministic ownership where concurrent execution must be serialized.",
    scenariosEn: ["Leader election", "Critical section coordination"]
  },
  {
    slug: "stellgate",
    name: "Stellgate",
    chineseName: "视界",
    domainEn: "Gateway",
    domainZh: "网关入口",
    originEn: "Event Horizon",
    originZh: "视界",
    rationaleEn: "A gateway is the entry boundary, like a horizon separating external and internal traffic.",
    summaryEn:
      "Provides ingress traffic management, authentication, protocol translation, and controlled API exposure.",
    positioningEn:
      "Stellgate concentrates access control and entry-layer protocol concerns at the platform edge.",
    scenariosEn: ["North-south API exposure", "Centralized ingress governance"]
  },
  {
    slug: "stellflow",
    name: "Stellflow",
    chineseName: "彗流",
    domainEn: "Message Queue",
    domainZh: "消息队列",
    originEn: "CometFlow",
    originZh: "彗流",
    rationaleEn: "Ordered data streams that travel like a comet trail.",
    summaryEn:
      "Provides asynchronous decoupling, event delivery, stream distribution, and traffic smoothing for message workloads.",
    positioningEn:
      "Stellflow supports queue and stream patterns for loosely coupled system interactions.",
    scenariosEn: ["Event-driven integration", "Traffic peak smoothing"]
  },
  {
    slug: "stellspec",
    name: "Stellspec",
    chineseName: "星谱",
    domainEn: "Log Platform",
    domainZh: "日志平台",
    originEn: "Spectrum",
    originZh: "星谱",
    rationaleEn: "System behavior is inspected through its signal spectrum.",
    summaryEn:
      "Provides unified collection, parsing, retrieval, and retention for application, audit, and platform logs.",
    positioningEn:
      "Stellspec provides searchable operational evidence across applications and infrastructure.",
    scenariosEn: ["Centralized log analysis", "Audit trail retention"]
  },
  {
    slug: "stellcon",
    name: "Stellcon",
    chineseName: "星座",
    domainEn: "Metrics Platform",
    domainZh: "指标平台",
    originEn: "Constellation",
    originZh: "星座",
    rationaleEn: "Metrics points connect into recognizable constellations for operational reasoning.",
    summaryEn:
      "Provides metric collection, aggregation, capacity analysis, and SLO-oriented observability.",
    positioningEn:
      "Stellcon turns raw metrics into capacity, reliability, and service health signals.",
    scenariosEn: ["Service SLO tracking", "Capacity planning dashboards"]
  },
  {
    slug: "stellvox",
    name: "Stellvox",
    chineseName: "星讯",
    domainEn: "Alerting Platform",
    domainZh: "告警平台",
    originEn: "NovaSignal",
    originZh: "星讯",
    rationaleEn: "Vox means voice or message, matching an alerting system's notification role.",
    summaryEn:
      "Provides alert evaluation, deduplication, on-call notification, and incident-coordination workflows.",
    positioningEn:
      "Stellvox connects monitoring signals with real operator response workflows.",
    scenariosEn: ["Alert convergence and noise reduction", "On-call incident dispatch"]
  },
  {
    slug: "stellguard",
    name: "Stellguard",
    chineseName: "星盾",
    domainEn: "Zero Trust Platform",
    domainZh: "零信平台",
    originEn: "StarShield",
    originZh: "星盾",
    rationaleEn: "Dynamic guard semantics suit continuous verification better than a static shield metaphor alone.",
    summaryEn:
      "Provides identity, access control, mTLS, and continuous verification for zero-trust security programs.",
    positioningEn:
      "Stellguard centralizes trust enforcement across identities, workloads, and service boundaries.",
    scenariosEn: ["Service-to-service trust enforcement", "Identity-centric access control"]
  },
  {
    slug: "stellkey",
    name: "Stellkey",
    chineseName: "星钥",
    domainEn: "Key Management Center",
    domainZh: "密钥中心",
    originEn: "StarKey",
    originZh: "星钥",
    rationaleEn: "A direct and forceful name for secret and certificate management responsibilities.",
    summaryEn:
      "Provides centralized secret, certificate, and access-policy management with rotation and audit support.",
    positioningEn:
      "Stellkey secures cryptographic assets and standardizes how they are distributed and rotated.",
    scenariosEn: ["Secret lifecycle governance", "Certificate distribution and rotation"]
  }
];

export const topics = [
  {
    slug: "dsl",
    titleEn: "Why CUE Works Well as a Configuration DSL",
    titleZh: "最佳 DSL 语言：CUE",
    categoryEn: "Configuration Engineering",
    categoryZh: "配置工程",
    summaryEn:
      "A practical comparison of type constraints, reuse, validation, and multi-environment governance to explain why CUE is a strong fit for complex declarative configuration.",
    tagsEn: ["DSL", "CUE", "Configuration Language", "Config Governance"],
    readingDirectionEn:
      "Read this when evaluating configuration language choices, schema unification, or platform-level configuration engineering."
  },
  {
    slug: "service-naming",
    titleEn: "Service Naming for Very Large Enterprises",
    titleZh: "面向超大型企业的微服务命名体系研究",
    categoryEn: "Service Governance",
    categoryZh: "服务治理",
    summaryEn:
      "A naming-system discussion for large organizations that need stable, expressive, and governable service identities across many business domains.",
    tagsEn: ["Microservices", "Naming", "Governance", "Architecture"],
    readingDirectionEn:
      "Read this when defining service identity rules or cleaning up inconsistent naming across large service estates."
  },
  {
    slug: "observability-spec",
    titleEn: "Observability Specification",
    titleZh: "可观测规范",
    categoryEn: "Observability",
    categoryZh: "可观测性",
    summaryEn:
      "A baseline observability specification covering signals, naming, and operational expectations across infrastructure and application layers.",
    tagsEn: ["Observability", "Tracing", "Metrics", "Logging"],
    readingDirectionEn:
      "Read this when standardizing telemetry conventions or defining platform-wide observability contracts."
  },
  {
    slug: "traces",
    titleEn: "Tracing Research for Large-Scale Enterprises",
    titleZh: "大型企业跨语言微服务链路追踪技术调研方案",
    categoryEn: "Distributed Tracing",
    categoryZh: "链路追踪",
    summaryEn:
      "A research-oriented walkthrough of cross-language tracing design choices, interoperability concerns, and rollout considerations for large enterprises.",
    tagsEn: ["Tracing", "Microservices", "OpenTelemetry", "Research"],
    readingDirectionEn:
      "Read this when comparing tracing architectures or planning a platform-wide tracing rollout."
  },
  {
    slug: "error-code-spec",
    titleEn: "Error Code Specification",
    titleZh: "错误码规范",
    categoryEn: "Application Contract",
    categoryZh: "应用契约",
    summaryEn:
      "A specification for shaping error codes into a stable contract that is easier to govern, observe, and consume across teams.",
    tagsEn: ["Error Code", "Contract", "API", "Governance"],
    readingDirectionEn:
      "Read this when trying to make service errors more structured, machine-readable, and operationally useful."
  },
  {
    slug: "middleware-evolution",
    titleEn: "Why Enterprises Build Middleware Platforms",
    titleZh: "为什么企业要自研中间件",
    categoryEn: "Platform Strategy",
    categoryZh: "平台战略",
    summaryEn:
      "An engineering and organizational perspective on when self-built middleware becomes justified and what tradeoffs it introduces.",
    tagsEn: ["Middleware", "Platform", "Architecture", "Strategy"],
    readingDirectionEn:
      "Read this when evaluating build-vs-buy decisions or the long-term cost model of infrastructure platforms."
  },
  {
    slug: "distributed-consistency",
    titleEn: "Consistency Challenges in Distributed Systems",
    titleZh: "分布式系统中的一致性挑战及其解决路径",
    categoryEn: "Distributed Systems",
    categoryZh: "分布式系统",
    summaryEn:
      "A concise discussion of consistency challenges, failure modes, and the decision paths commonly used to address them in distributed systems.",
    tagsEn: ["Consistency", "Distributed Systems", "Reliability", "Transactions"],
    readingDirectionEn:
      "Read this when comparing consistency strategies or selecting a reliability model for cross-service workflows."
  },
  {
    slug: "distributed-system-registry-centers",
    titleEn: "Registry Centers in Distributed Systems",
    titleZh: "分布式系统注册中心意义、问题与主流实现",
    categoryEn: "Infrastructure Foundation",
    categoryZh: "基础设施",
    summaryEn:
      "A review of why registry centers exist, what problems they solve, and how mainstream implementations make different engineering tradeoffs.",
    tagsEn: ["Registry", "Discovery", "Distributed Systems", "Infrastructure"],
    readingDirectionEn:
      "Read this when evaluating service discovery patterns or studying registry-center implementation choices."
  }
];
