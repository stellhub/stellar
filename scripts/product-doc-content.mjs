export const productDocs = {
  stellmap: {
    overview: {
      tagline:
        "A registry center that manages service instance registration, health checks, discovery subscriptions, and topology change propagation.",
      goal:
        "Stellmap provides unified service registration and discovery for microservices and infrastructure components, helping teams handle endpoint drift, environment switching, elastic scaling, and multi-cluster routing.",
      boundaries: [
        "Designed for service registration, instance discovery, and subscription-driven change propagation.",
        "Designed for service directory governance across multiple environments, clusters, and canary routing scenarios.",
        "It does not carry business configuration or application-layer traffic governance logic."
      ],
      capabilities: [
        "Supports service registration, instance deregistration, heartbeat renewal, and health checking.",
        "Supports namespace, cluster, label, and version-based instance isolation.",
        "Supports subscription push and local client caching to reduce lookup latency."
      ],
      values: [
        "Standardizes how service entry points are discovered, reducing the operational cost of maintaining endpoints manually.",
        "Uses labels and versions to support canary rollout and cross-region routing."
      ],
      businessScenarios: ["Microservice instance discovery", "Multi-cluster canary routing"],
      platformScenarios: [
        "Infrastructure node directory management",
        "Unified service catalog governance"
      ]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The control plane manages service metadata, adapts registration protocols, and converges instance state.",
              "Service definitions, grouping labels, and subscription relationships are all maintained centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The data plane is responsible for high-availability queries, incremental push delivery, and local client-side cache support.",
              "A replicated consistency layer and lease model ensure instance information can be recovered after failures."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Naming Context",
            paragraphs: [
              "`StellMap` means a star map or coordinate chart within the Stell Hub product family.",
              "That name reflects the role of a registry center: every service instance is like a point on a map, and the registry keeps track of its current location and state so callers can discover, locate, and navigate across the system."
            ]
          },
          {
            heading: "Implementation Scope",
            bullets: [
              "This repository carries the Go implementation of StellMap.",
              "It is intended to evolve around service registration, service discovery, instance heartbeat and health state management, service metadata governance, namespace and grouping isolation, and integration contracts with governance, configuration, and control-plane modules."
            ]
          },
          {
            heading: "Engineering Goals",
            bullets: [
              "Consistency first: a CP-oriented design that favors linearizable correctness over apparent availability.",
              "Lightweight operation: a single Raft group can carry the registry data plane without unnecessary coordination layers.",
              "High availability: service remains available to clients as long as a quorum survives.",
              "High concurrency: read and write paths stay short to avoid accidental bottlenecks.",
              "Fast recovery: crash recovery remains explicit, with clear ownership across WAL, storage, and snapshots.",
              "Future evolution: the design can grow into watch, lease, compression, and layered cache capabilities without reworking the core."
            ]
          },
          {
            heading: "Runtime Topology",
            paragraphs: [
              "The current `stellmapd` process exposes three listening surfaces: public HTTP, dedicated admin HTTP, and internal gRPC.",
              "Together they form a replicated registry service where external clients use HTTP for registration and discovery, while nodes use gRPC for replication, snapshot transfer, and leader-oriented internal traffic."
            ]
          },
          {
            heading: "Consensus",
            bullets: [
              "The registry cluster runs as a single Raft group built on `etcd-io/raft`.",
              "The core object is the service instance record, and its metadata scale is typically much smaller than a general-purpose database.",
              "A single consensus group greatly reduces implementation and operational complexity.",
              "For a registry center, consistency and maintainability are usually more important than horizontal sharding."
            ]
          },
          {
            heading: "Domain Model",
            bullets: [
              "The external model is an instance registry rather than a general-purpose key-value product.",
              "Its logical primary key is `namespace / service / instanceId`.",
              "Instance records store endpoints, labels, metadata, lease TTL, and recent heartbeat information.",
              "Normalized service names and structured fields are retained together to support prefix subscription, permission governance, and aggregated monitoring."
            ]
          },
          {
            heading: "Consistency Guarantees",
            bullets: [
              "StellMap explicitly adopts a CP architecture.",
              "During a network partition, minority nodes stop serving linearizable writes.",
              "The cluster prefers no divergence, no rollback, and no stale values over temporary write availability.",
              "All reads are linearizable by default, implemented through `ReadIndex` and local state-machine reads once `appliedIndex >= readIndex`."
            ]
          },
          {
            heading: "Membership Management",
            bullets: [
              "Node membership changes support `Learner` and joint consensus.",
              "New nodes first join as learners, receive logs without voting, and are promoted only after they catch up and pass health checks.",
              "Large topology updates rely on `ConfChangeV2` and joint consensus to avoid quorum instability during transitions."
            ]
          },
          {
            heading: "Storage Layout",
            bullets: [
              "`WAL` stores the Raft log and consensus metadata such as `Entry` and `HardState`.",
              "`Pebble` stores applied registry data plus a small amount of local metadata such as apply progress, snapshot markers, and control watermarks.",
              "`Snapshot` files are managed separately for recovery, validation, and atomic replacement."
            ]
          },
          {
            heading: "Storage Engine Choice",
            bullets: [
              "Pebble is a Go-native LSM key-value engine that avoids `cgo` and fits frequent registration and heartbeat updates well.",
              "It supports range delete, snapshots, and batched writes, which are useful for instance cleanup, snapshotting, and batch apply.",
              "Its operational maturity and low maintenance burden make it a pragmatic fit for registry-state storage."
            ]
          },
          {
            heading: "Verification Strategy",
            bullets: [
              "Core validation must include linearizable read testing, write-after-read guarantees, and candidate-set consistency checks.",
              "Membership testing must cover learner join, promotion, and joint consensus transitions.",
              "Failure testing must cover leader crashes, follower restart, partition tolerance, WAL corruption, snapshot corruption, and interrupted snapshot transfer.",
              "Load testing must cover high-frequency registration, heartbeat renewal, discovery queries, watch streams, and long-running soak scenarios."
            ]
          },
          {
            heading: "Module Map",
            bullets: [
              "`raftnode` drives the replicated state machine and coordinates `Ready`, `ReadIndex`, and membership changes.",
              "`wal` persists the Raft log and manages segment rotation, sync, and recovery.",
              "`snapshot` exports, installs, validates, and cleans up snapshots.",
              "`storage` implements the state machine on top of Pebble and guarantees ordered, idempotent apply.",
              "`transport` separates external HTTP access from internal gRPC replication."
            ]
          },
          {
            heading: "Internal Protocol",
            bullets: [
              "Internal replication is defined by Raft and snapshot protobuf contracts and is never exposed to business clients directly.",
              "`RaftTransport` batches ordinary Raft messages, while `SnapshotService` handles streamed snapshot upload and download.",
              "The gRPC layer converts protobuf payloads into local runtime models, and the runtime transport service executes the real message flow."
            ]
          },
          {
            heading: "Admin Control Surface",
            bullets: [
              "Cluster membership changes, leader transfer, and cluster status inspection are all triggered through `stellmapctl`, even though execution still lands on the dedicated admin HTTP listener.",
              "Public business clients only reach `HTTPAddr`, while `stellmapctl` targets the local `AdminAddr` by default.",
              "Admin requests currently require both localhost source access and `Authorization: Bearer <token>` authentication."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Start with a three-node deployment in a single region and spread nodes across availability zones to reduce single-point risk.",
              "In multi-region scenarios, local discovery traffic can be served through read-only follower clusters."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Registry recovery should rely on both leases and replicated copies of the registry state.",
              "Clients should enable local caching and tolerate short-lived control-plane disconnects."
            ]
          },
          {
            heading: "Crash Recovery Flow",
            ordered: [
              "Read the latest local snapshot metadata.",
              "If a valid snapshot exists, restore it into the state-machine working directory first.",
              "Open `Pebble` and load registry data plus local metadata.",
              "Open the `WAL` and recover `HardState`, `Entry`, and snapshot markers.",
              "Discard log segments already covered by the snapshot index.",
              "Apply the remaining committed entries to the state machine in order.",
              "Rebuild the in-memory Raft node, apply watermarks, and membership view.",
              "Only then transition the node into a serving state."
            ]
          },
          {
            heading: "Recovery Rules",
            bullets: [
              "If the WAL is durable but the state machine has not applied the latest committed entries, replay the log at startup.",
              "If a snapshot file exists but metadata was not atomically switched, keep using the old snapshot.",
              "If Pebble reports an applied index behind the WAL commit index, continue applying entries until it catches up.",
              "If snapshot, WAL, and state-machine indexes disagree, prioritize the principle of never rolling back committed logs."
            ]
          },
          {
            heading: "Post-Recovery Checks",
            bullets: [
              "Verify `HardState.Commit >= appliedIndex`.",
              "Verify the state machine's `appliedIndex` does not exceed the largest committed index in the WAL.",
              "Verify the snapshot `ConfState` matches the in-memory membership view after recovery.",
              "Verify whether the node is still part of the latest membership and whether it needs to fetch additional logs or install a newer snapshot from the leader."
            ]
          },
          {
            heading: "Snapshot and Log Compaction",
            bullets: [
              "Trigger snapshot generation when `appliedIndex - snapshotIndex` crosses a configured threshold.",
              "Write snapshot metadata atomically after snapshot generation finishes.",
              "Retain only the WAL segments that are still necessary after the snapshot point.",
              "Let lagging nodes catch up by log replay first and switch to snapshot installation only when the gap becomes too large."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Create a namespace and a normalized service identity.",
              "Integrate `stellmap-client` into the application and enable heartbeat renewal.",
              "Enable service subscription and local cache support on the caller side."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Check through the console or API that the instance has been registered successfully.",
              "Simulate instance removal and confirm that subscribers receive the change event."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "In production, separate registration timeout, removal threshold, and push batch size instead of relying on a single global timing model.",
              "For cross-datacenter traffic, add region labels and nearest-routing policies."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Tune registration timeout and heartbeat intervals according to service startup time and network quality.",
              "Balance push batch size and local cache refresh frequency against freshness and throughput."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "Write Contract",
            ordered: [
              "A client sends a registration, deregistration, or heartbeat renewal request through the public HTTP API.",
              "If the request lands on a follower, it is redirected or forwarded to the leader.",
              "The leader encodes the instance change into a proposal and submits it to the single Raft group.",
              "The log is appended to the local WAL.",
              "The log is replicated to a quorum and committed.",
              "The state machine applies the change to Pebble in order.",
              "The write success response is returned."
            ]
          },
          {
            heading: "Read Contract",
            ordered: [
              "A client sends an instance query request.",
              "The request reaches the leader, or a follower forwards it to the leader path.",
              "The leader executes `ReadIndex`.",
              "The node waits until local `appliedIndex` catches up to `readIndex`.",
              "The latest registry view is read from Pebble.",
              "A linearizable response is returned."
            ]
          },
          {
            heading: "Public API Contract",
            bullets: [
              "Third-party integration is exposed through the public HTTP API only.",
              "That keeps onboarding simple for scripts, sidecars, and multi-language clients while preserving a clear boundary for registration, discovery, query, and health-reporting interfaces."
            ]
          },
          {
            heading: "Admin API Contract",
            bullets: [
              "Membership changes, leader transfer, and cluster-status inspection run on a dedicated admin HTTP listener instead of the public business listener.",
              "The public HTTP surface carries only business data-plane and health endpoints.",
              "The admin listener is bound to loopback by default, currently accepts only localhost traffic, and requires `Authorization: Bearer <token>` authentication.",
              "`stellmapctl` is the only supported control-plane entry point today."
            ]
          },
          {
            heading: "Internal Transport Contract",
            bullets: [
              "Cluster-internal replication uniformly uses gRPC.",
              "The external API and internal replication surface should not share protocol assumptions because the internal channel carries Raft messages, snapshots, and leader-forwarded read or write traffic.",
              "External HTTP and internal gRPC still share the same core state-machine and permission-validation logic to avoid double implementation."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on registration success rate, renewal latency, push backlog, and cache hit ratio."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for instance flapping, abnormal lease expiration, and push-channel backlog.",
              "Watch for replication delay in registration data across clusters."
            ]
          }
        ]
      }
    }
  },
  stellnula: {
    overview: {
      tagline:
        "A configuration center responsible for centralized storage, version management, progressive release, and dynamic distribution.",
      goal:
        "Stellnula provides unified configuration governance for services, gateways, and platform components, helping teams eliminate scattered configuration, environment drift, and unaudited manual changes.",
      boundaries: [
        "Designed for application parameters, platform switches, and runtime configuration distribution.",
        "Designed for configuration version management, progressive rollout, and audit traceability.",
        "Sensitive configuration should be linked with the secret-management center instead of being stored in plain text."
      ],
      capabilities: [
        "Supports configuration isolation by namespace, application, environment, and label.",
        "Supports version history, progressive release, rollback, and change auditing.",
        "Supports push-based listeners, long polling, and configuration snapshot fallback."
      ],
      values: [
        "Reduces coupling between application deployment and parameter adjustment.",
        "Brings configuration governance into a unified change-management workflow to reduce operator mistakes."
      ],
      businessScenarios: ["Dynamic configuration delivery for microservices", "Canary parameter experiments"],
      platformScenarios: ["Multi-environment configuration governance", "Unified platform switch management"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The control plane handles configuration orchestration, approval flows, and release-plan management.",
              "Configuration layering and environment isolation are defined centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The storage layer keeps both text and structured configuration in a versioned model.",
              "Clients consume incremental updates and local disaster-recovery fallbacks through watch channels."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Configuration Repository: stores configuration versions and metadata.",
              "Release Engine: handles canary release, phased rollout, and scheduled publishing.",
              "Client Agent: receives changes and persists local cache snapshots."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "The release engine pushes incremental changes to clients according to the rollout plan.",
              "When remote fetch fails, the client falls back to local snapshots to keep startup stable."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Deploy the configuration center separately from the registry center so configuration traffic does not interfere with service discovery.",
              "At larger scale, split the read-only query layer from the release execution layer."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Separate the release path from the query path to reduce mutual impact during heavy change windows.",
              "Configuration snapshots should support cross-instance sharing and cold-start fallback."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Create the application and environment.",
              "Publish the first configuration version and define a rollback point.",
              "Integrate `stellnula-client` into the application to listen for changes."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Change a non-sensitive configuration item and confirm that the client hot reloads successfully.",
              "Roll back and verify that the client returns to the previous stable version."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Manage core configuration through layered categories such as `base configuration / secret references / environment differences`.",
              "Store only reference identifiers for sensitive values and let the secret-management center handle decryption."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Keep configuration names stable and traceable to avoid ambiguity across teams.",
              "During progressive rollout, validate with a small traffic segment first and expand gradually."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose configuration query, listener, publish, and rollback interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK supports preload at startup, runtime refresh, and change callbacks.",
              "Critical configuration changes should include local validation and callback protection."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on release success rate, client latency, rollback count, and configuration drift."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for failed releases, listener backlog, and clients that have not refreshed for a long time.",
              "Watch for drift and inconsistency for the same configuration across environments."
            ]
          }
        ]
      }
    }
  },
  stelltrace: {
    overview: {
      tagline:
        "A tracing platform for end-to-end trace collection, span analysis, cross-signal correlation, and issue localization.",
      goal:
        "Stelltrace provides unified tracing for service call graphs, asynchronous tasks, and cross-gateway requests so engineering and operations teams can quickly identify slow requests, error propagation, and cross-system dependencies.",
      boundaries: [
        "Designed for request tracing across synchronous calls, asynchronous messaging, and scheduled tasks.",
        "Designed for troubleshooting, topology analysis, and cross-system dependency investigation.",
        "It does not replace the log platform or metrics platform; instead, it acts as their correlation backbone."
      ],
      capabilities: [
        "Supports tracing across HTTP, gRPC, MQ, scheduled tasks, and other request paths.",
        "Supports sampling policies, anomaly clustering, and dependency-topology views.",
        "Supports trace-correlated queries across logs, metrics, and alerts."
      ],
      values: [
        "Uses Trace ID as the centerline that connects multiple observability signals.",
        "Improves troubleshooting efficiency for slow requests and abnormal propagation."
      ],
      businessScenarios: ["Slow request troubleshooting", "Root-cause analysis for anomalies"],
      platformScenarios: ["Service dependency topology governance", "Full-path observability correlation analysis"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The ingestion layer handles protocol access, context propagation, and sampling decisions.",
              "Sampling strategy, field conventions, and retention windows are managed centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The processing layer aggregates spans, builds indexes, and computes dependency relationships.",
              "The query layer serves single-trace details, topology search, and troubleshooting views."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Agent / SDK: handles instrumentation and context propagation.",
              "Collector: receives spans and performs pre-aggregation.",
              "Query API: provides search, analysis, and visualization interfaces."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "Agents report spans to the collector for aggregation and cleanup.",
              "The query layer reconstructs execution paths by Trace ID and correlates them with logs and metrics."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Decouple the ingestion layer from the storage layer so collectors can scale horizontally.",
              "Store hot and cold data in different tiers to control long-term retention cost."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Collector clusters should support batch buffering and failover.",
              "Query nodes should be separated from index nodes to avoid resource contention."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Enable the Trace SDK inside the application.",
              "Configure the sampling strategy and reporting endpoint.",
              "Use Trace ID in the platform to retrieve the full call chain."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Trigger a cross-service request and confirm the full trace is visible.",
              "Intentionally create a slow call and verify that slow-trace filtering works."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Use layered sampling in production and prioritize keeping errors and slow requests.",
              "Enable high-cardinality label allowlists for critical business paths."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Use adaptive sampling on ingress traffic to balance cost and issue coverage.",
              "For asynchronous chains, enforce explicit consistency in context propagation."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose trace query, dependency analysis, and sampling-policy interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK supports auto instrumentation, manual instrumentation, and custom context-propagation extensions.",
              "Gateways, task systems, and messaging components should adopt the same context standard."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on span throughput, sampling hit rate, index latency, and query latency."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for collector backpressure, index backlog, and uncontrolled high-cardinality labels.",
              "Watch for broken traces and failures in context propagation."
            ]
          }
        ]
      }
    }
  },
  stellorbit: {
    overview: {
      tagline:
        "A service-governance hub for routing, load balancing, traffic control orchestration, and service lifecycle governance.",
      goal:
        "Stellorbit provides unified governance policies for service-to-service calls, solving dynamic control problems such as canary release, geo-routing, version isolation, and failover.",
      boundaries: [
        "Designed for routing, retries, traffic shifting, and failure isolation in service call paths.",
        "Designed for governance rule orchestration across multiple versions, regions, and tenants.",
        "It does not store registry or configuration data directly; instead, it consumes those external contracts."
      ],
      capabilities: [
        "Supports weighted, label-based, version-based, and geo-based routing.",
        "Supports load balancing, retries, timeout policy, and unhealthy-instance ejection.",
        "Supports progressive activation and rollback of governance rules."
      ],
      values: [
        "Moves traffic governance from embedded code to platform-level orchestration.",
        "Improves control over canary release and active-active traffic switching."
      ],
      businessScenarios: ["Canary release", "Active-active traffic switching"],
      platformScenarios: ["Service failure isolation", "Cross-cluster governance orchestration"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "A strategy center manages routing rules and governance templates centrally.",
              "Rule publication, canary rollout, and rollback are all controlled from the control plane."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The execution plane enforces policies in real time through sidecars or SDKs.",
              "The state plane reports circuit state, service profiles, and capacity data back to the governance center."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Rule Repository: stores routing and governance rules.",
              "Execution Engine: enforces strategy at traffic ingress and client-call points.",
              "State Synchronizer: reports instance health, capacity, and abnormal events."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "Clients pull or receive governance rules and execute them immediately in the request path.",
              "Execution results and instance status are sent back for policy optimization."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Support three governance modes: pure SDK, sidecar, and gateway-integrated governance.",
              "Deploy the core control plane with multiple replicas to keep rule publication continuous."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Rule cache must support local persistence and disconnect tolerance.",
              "Gateway-side and service-side governance can be combined to reduce single-point policy failure."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Define version labels and traffic groups for the service.",
              "Configure baseline timeout, retry, and load-balancing policies.",
              "Publish a canary routing rule and observe governance metrics."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Roll out a canary rule to a small amount of traffic and verify hit behavior.",
              "Simulate an unhealthy instance and confirm ejection and traffic switching behave as expected."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Route rules should be layered into global, tenant-level, and service-level strategies.",
              "Retry and timeout settings should be designed together to avoid amplifying failure traffic."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Critical routing rules should include approval and rollback plans.",
              "Limit retry depth and timeout stacking to avoid request avalanches."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose governance rule query, canary publishing, and state-reporting interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK exposes extensions for load balancing, routing, circuit breaking, and retries.",
              "In sidecar deployments, public rules and business rules should be managed in separate layers."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on rule hit rate, retry amplification factor, circuit-break count, and traffic-switch latency."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for abnormal rule spread, switching delay, and rollback windows.",
              "Watch for dependency amplification and distorted capacity caused by excessive retries."
            ]
          }
        ]
      }
    }
  },
  stellpulse: {
    overview: {
      tagline:
        "A flow-control and circuit-breaking platform for hotspot protection, capacity guarding, bulkhead isolation, and adaptive degradation.",
      goal:
        "Stellpulse focuses on availability protection under high concurrency, using rate limiting, concurrency isolation, circuit breaking, degradation, and system load defense to keep critical services stable.",
      boundaries: [
        "Designed for stability protection at service ingress, critical interfaces, and dependency calls.",
        "Designed for sudden traffic spikes, hotspot resources, and unstable dependencies.",
        "It does not replace service governance; it focuses on admission control and degradation protection."
      ],
      capabilities: [
        "Supports QPS limiting, concurrency limiting, and hotspot parameter protection.",
        "Supports circuit breaking, bulkhead isolation, warmup, and request queuing.",
        "Supports automatic degradation based on error rate, slow calls, and system load."
      ],
      values: [
        "Provides a unified availability-protection foundation for critical call paths.",
        "Keeps the impact of traffic spikes and dependency failures within predictable bounds."
      ],
      businessScenarios: ["Traffic protection for flash-sale workloads", "Fallback for critical call paths"],
      platformScenarios: ["Dependency failure isolation", "Unified platform-level traffic guard"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The strategy center manages rule configuration, dynamic push, and version control.",
              "Rule publishing, canary rollout, and rollback are handled centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The runtime engine makes in-process decisions quickly to keep limiting overhead low.",
              "A status-reporting module reports circuit windows, hotspot statistics, and system load."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Rule Center: centrally manages rate-limiting and circuit-breaking rules.",
              "Runtime Engine: decides admission and degradation at request time.",
              "Metrics Reporter: reports hit, rejection, and recovery events."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "Rules are pushed to runtime components and enforced locally at the request entry point.",
              "Rejection, circuit-break, and recovery events are returned to the monitoring and alerting systems."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "The preferred mode is SDK-based embedding inside business applications, optionally combined with gateway-level rules.",
              "For major promotion events, coordinate layered rate limiting with the gateway and service-governance platform."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Local rule cache should remain effective even when the control plane is unavailable.",
              "Rate-limiting and circuit-breaking rules should be deployed in multiple layers instead of concentrating on a single entry point."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Configure baseline QPS thresholds and slow-call thresholds for each interface.",
              "Define circuit-break windows and recovery strategies for critical dependencies.",
              "Validate rule hits and degradation paths in a performance-testing environment."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Simulate over-threshold traffic and confirm that the rate-limiting rule is triggered.",
              "Simulate dependency failures and confirm the circuit-break and recovery window behave as expected."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Circuit recovery should use half-open probing to avoid an instant full-volume recovery.",
              "Hotspot-parameter rules should include allowlists and fallback logic."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Different interfaces should have independent thresholds based on capacity class and business criticality.",
              "Degradation logic should be designed together with fallback pages or default responses."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose rule query, dynamic push, and event-subscription interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK supports annotations, programmatic APIs, and gateway adapters.",
              "Use the same rule model on both the gateway and service side."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on rejection rate, circuit-break duration, hotspot distribution, and system load."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for false-positive rule hits, hotspot spikes, and sudden load surges.",
              "Watch whether circuit recovery pacing matches dependency recovery pacing."
            ]
          }
        ]
      }
    }
  },
  stellabe: {
    overview: {
      tagline:
        "A job-scheduling platform for timed orchestration, dependency scheduling, sharded execution, and runtime governance.",
      goal:
        "Stellabe provides unified scheduling for platform-level and business-level jobs, including batch processing, periodic execution, workflow orchestration, and failure retry.",
      boundaries: [
        "Designed for scheduled jobs, workflows, and batch-processing scenarios.",
        "Designed for replay, backfill, and visualization of execution progress.",
        "It does not replace stream-processing engines; it focuses on orchestration and scheduling control."
      ],
      capabilities: [
        "Supports cron scheduling, workflow DAGs, task dependencies, and sharded execution.",
        "Supports failure retry, idempotency control, backfill, and replay execution.",
        "Supports tenant isolation, task auditing, and runtime visualization."
      ],
      values: [
        "Brings scattered scripts and manual tasks into a unified governance model.",
        "Improves the maintainability and execution reliability of task orchestration."
      ],
      businessScenarios: ["Scheduled reporting", "Data backfill"],
      platformScenarios: ["Post-release task orchestration", "Platform-level workflow governance"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The orchestration layer defines workflows, analyzes dependencies, and publishes tasks.",
              "Task versions, approval flow, and release process are managed centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The scheduling layer computes triggers, allocates shards, and controls execution windows.",
              "The execution layer handles worker registration, status reporting, and failure recovery."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Scheduler: triggers jobs and computes execution plans.",
              "Executor: pulls jobs, executes them, and reports results.",
              "Workflow Store: stores job definitions, run records, and audit information."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "The scheduler generates execution plans according to trigger conditions and dispatches them to executors.",
              "Executors report execution state and logs to form a full runtime record."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Deploy the control plane separately from the execution plane so job peaks do not impact the management console.",
              "For heavy workloads, split worker clusters by queue type or workload type."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Schedulers must guarantee single-trigger semantics and support preemption recovery.",
              "Executors should isolate resource pools by workload type to avoid mutual interference."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Define tasks and executors.",
              "Configure cron schedules or DAG dependencies.",
              "Observe the latest execution log and retry history."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Manually trigger a task and confirm that an executor can pick it up normally.",
              "Simulate a failure and verify that retry and compensation chains work."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Task timeout, retry count, and concurrency should be configured independently by task type.",
              "For backfill tasks, define time windows and resource ceilings."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Critical tasks must define explicit idempotency and failure-compensation logic.",
              "Limit DAG depth and fan-out width to avoid scheduling storms."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose task management, execution control, log query, and workflow publication interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK provides Java and Go executor integration templates plus heartbeat protocols.",
              "Executors must explicitly report progress, heartbeat, and final state."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on trigger latency, success rate, retry depth, and executor utilization."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for task backlog, executor disconnects, and retry avalanches.",
              "Watch how failures at key DAG nodes propagate to downstream paths."
            ]
          }
        ]
      }
    }
  },
  stellpoint: {
    overview: {
      tagline:
        "A distributed locking and coordination center for mutual exclusion, leader election, and serialized access to critical resources.",
      goal:
        "Stellpoint addresses shared-resource contention in distributed systems by providing a unified lock service, lease model, and coordination primitives, reducing repeated one-off implementations.",
      boundaries: [
        "Designed for mutual exclusion, leader election, and ordered execution.",
        "Designed for shared-resource write protection and serialization of critical tasks.",
        "It does not replace business transactions; it provides foundational distributed coordination."
      ],
      capabilities: [
        "Supports reentrant locks, read-write locks, lease locks, and leader-election locks.",
        "Supports lock renewal, preemption, timeout release, and deadlock prevention.",
        "Supports namespace isolation and auditability based on resource paths."
      ],
      values: [
        "Unifies distributed-lock semantics and avoids repeated custom implementations in business systems.",
        "Uses leases and fencing tokens to reduce accidental lock misuse and dirty writes."
      ],
      businessScenarios: ["Serialized inventory deduction", "Mutual exclusion for scheduled tasks"],
      platformScenarios: ["Leader election", "Distributed coordination center"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The lock control plane manages the resource model, permissions, and lock-strategy publication.",
              "Resource paths, lease policies, and permission boundaries are maintained centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The coordination storage layer persists lock state and lease information through a consistent log.",
              "Clients preserve ordering guarantees through session keepalive and fencing tokens."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Lock API: receives lock, unlock, renewal, and query requests.",
              "Lease Manager: manages sessions, leases, and expiration cleanup.",
              "Coordination Log: stores lock ordering and fencing tokens."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "Clients keep leases alive through sessions and periodic renewal.",
              "The coordination log emits fencing tokens to preserve holder-order semantics."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Deploy an odd number of replicas, usually three to five nodes, to preserve quorum behavior.",
              "For high-risk resources, use dedicated namespaces or separate lock clusters."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Session timeout and lease expiration must be tuned carefully to avoid accidental release.",
              "For hotspot resources, split lock granularity and limit preemption frequency."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Define the resource path and lease duration.",
              "Acquire the lock inside business code and execute the critical logic.",
              "Use fencing tokens to prevent stale holders from writing."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Simulate two instances competing for the same resource and confirm only one gets the lock.",
              "Simulate an abnormal exit of the lock holder and confirm the lease is reclaimed after expiration."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Long-running transactions should coordinate automatic lease renewal with business timeout behavior.",
              "For hotspot resources, add backoff retries and concurrency quotas."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Lock granularity should be designed together with business idempotency strategy.",
              "Critical resources should emit audit logs and abnormal-alert signals."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose interfaces for lock management, leader-election watches, and lease queries."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK supports synchronous, asynchronous, and annotation-based locking.",
              "Use fencing tokens together with business version numbers for double protection."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on lock wait time, lock-acquisition failure rate, renewal success rate, and session drift."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for session drift, hotspot lock contention, and abnormal lease reclamation.",
              "Watch for coordination-log delay and quorum-node instability."
            ]
          }
        ]
      }
    }
  },
  stellgate: {
    overview: {
      tagline:
        "A gateway ingress platform for unified access, authentication, protocol translation, traffic governance, and API exposure.",
      goal:
        "Stellgate acts as the external entry point and internal traffic-orchestration hub, receiving client requests, applying unified security policies, and forwarding traffic to backend services.",
      boundaries: [
        "Designed for external APIs, internal service ingress, and edge-access scenarios.",
        "Designed for protocol access, authentication and authorization, forwarding governance, and plugin extensibility.",
        "It does not replace a service-governance platform; it focuses on ingress-side traffic handling."
      ],
      capabilities: [
        "Supports HTTP, gRPC, WebSocket, and SSE access.",
        "Supports authentication, rate limiting, routing, rewriting, and plugin-based extension.",
        "Supports multi-tenant, multi-domain, and multi-environment gateway orchestration."
      ],
      values: [
        "Centralizes ingress governance, security, and protocol processing into a unified gateway.",
        "Improves consistency for both external API exposure and internal traffic orchestration."
      ],
      businessScenarios: ["Unified Open API ingress", "LLM gateway and edge access"],
      platformScenarios: ["Internal service gateway", "Edge traffic ingress"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The control plane manages routes, certificates, plugins, and release strategy.",
              "Gateway instances, domains, upstream services, and plugin chains are orchestrated centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The data plane uses stateless forwarding nodes that scale horizontally.",
              "An extension layer handles authentication, observability, transformation, and traffic control through plugin chains."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Gateway Admin: manages route, certificate, and plugin configuration.",
              "Gateway Runtime: performs forwarding, load balancing, and traffic control.",
              "Plugin SDK: carries authentication, audit, and protocol-extension logic."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "The control plane publishes route and plugin strategy to runtime gateway nodes.",
              "Requests are processed through the plugin chain before being forwarded to the target upstream service."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Deploy external and internal gateways in separate layers.",
              "Edge nodes should connect locally and receive routes from the configuration center."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Gateway nodes should be deployed statelessly and combined with elastic scaling.",
              "Certificates, routes, and plugins should all support canary release and fast rollback."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Create a gateway instance and domain mapping.",
              "Configure upstream services, route rules, and authentication policies.",
              "Publish the route and validate traffic through a canary domain."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Send an authenticated request and confirm both routing and the plugin chain take effect.",
              "Validate a small amount of canary traffic and confirm the release result is observable."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Routes should be layered by domain, path, tenant, and version.",
              "Security plugins and observability plugins are good candidates for a global default chain."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Use separate certificates and plugin strategies for public-facing and internal gateways.",
              "Long-lived connection protocols need dedicated timeout, retry, and connection-limit tuning."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose route management, certificate management, plugin management, and release interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK provides a plugin-development framework and gateway extension context.",
              "Plugin design should distinguish pre-authentication, pre-forward, and post-response stages."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on gateway latency, 5xx ratio, rate-limit hits, and upstream retry count."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for upstream errors, plugin execution latency, and connection-pool congestion.",
              "Watch for certificate expiry, canary release failure, and domain-configuration drift."
            ]
          }
        ]
      }
    }
  },
  stellflow: {
    overview: {
      tagline:
        "A message queue and event-stream platform for asynchronous decoupling, streaming distribution, burst smoothing, and event-driven integration.",
      goal:
        "Stellflow provides a unified messaging and event-distribution capability for asynchronous business flows, system integration, and streaming workloads.",
      boundaries: [
        "Designed for asynchronous decoupling, event-driven interaction, and data-distribution scenarios.",
        "Designed for delayed messages, ordered messages, and stream-consumption models.",
        "It does not replace a job-scheduling platform; it focuses on durable message storage and distribution."
      ],
      capabilities: [
        "Supports ordinary queues, topic subscription, delayed messages, and ordered messages.",
        "Supports consumer groups, retry queues, dead-letter queues, and replay consumption.",
        "Supports event bridging, streaming processing, and cross-region replication."
      ],
      values: [
        "Reduces coupling and traffic shocks through asynchronous design.",
        "Supports event-driven architecture through reliable delivery and offset management."
      ],
      businessScenarios: ["Order event-driven integration", "Asynchronous burst smoothing"],
      platformScenarios: ["Data synchronization bus", "Cross-system event bridging"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The governance layer manages topic configuration, quotas, audit, and event tracking.",
              "Topics, consumer groups, and replication strategy are administered centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The write layer handles message durability, partition routing, and write acknowledgement.",
              "The distribution layer handles consumer-group coordination, offset management, and rebalancing."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Broker: handles message write, storage, and distribution.",
              "Controller: manages topics, partitions, replicas, and failover.",
              "Client SDK: provides producer, consumer, and transactional-message capability."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "Producers write messages to brokers and receive acknowledgement results.",
              "Consumer groups obtain partition ownership through offset management and rebalancing."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Plan topics and partitions by business domain to avoid hotspot concentration.",
              "For cross-region traffic, deploy mirror replication or event-bridge channels."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Replica strategy must tolerate both single-node failures and single-availability-zone failures.",
              "Hot topics should use dedicated clusters or isolated resource pools."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Create a topic and a consumer group.",
              "Send business events through the SDK.",
              "Configure consumer retry, dead-letter handling, and offset-commit policy."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Send a test message and confirm that consumers can process it successfully.",
              "Simulate consumption failure and verify retry queues and dead-letter queues behave as expected."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Design separate topics for ordered-message workloads and high-throughput workloads.",
              "Set retry count and dead-letter thresholds according to business idempotency capability."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Plan partition count based on throughput, concurrency, and hotspot-key distribution.",
              "Design offset-commit strategy together with consumer idempotency and transaction requirements."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose topic management, message delivery, consumer offset, and audit-query interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK supports synchronous, asynchronous, batch, and transactional message models.",
              "Consumers should support batch pull, manual commit, and exception rollback strategy."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on production latency, consumer backlog, retry rate, and partition hotspots."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for partition imbalance, consumer backlog, and replica synchronization delay.",
              "Watch for retry storms and dead-letter backlog growth."
            ]
          }
        ]
      }
    }
  },
  stellspec: {
    overview: {
      tagline:
        "A log platform for unified collection, structured processing, search analysis, and retention governance.",
      goal:
        "Stellspec provides unified ingestion and retrieval for business logs, audit logs, and platform logs so teams can establish consistent logging standards and faster troubleshooting paths.",
      boundaries: [
        "Designed for unified governance of application logs, audit logs, and platform runtime logs.",
        "Designed for log search, trace correlation, and retention-policy management.",
        "It does not replace metrics or tracing; it focuses on full-text search and contextual reconstruction."
      ],
      capabilities: [
        "Supports unified ingestion of plain-text logs, structured logs, and audit logs.",
        "Supports field parsing, desensitization, index routing, and lifecycle management.",
        "Supports trace correlation, log clustering, and anomaly-pattern search."
      ],
      values: [
        "Creates consistent logging conventions and reduces search and troubleshooting cost.",
        "Balances security and cost through desensitization and lifecycle management."
      ],
      businessScenarios: ["Application log search", "Troubleshooting and root-cause analysis"],
      platformScenarios: ["Audit tracing", "Platform-level retention governance"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "Collection rules, index templates, and lifecycle policy are orchestrated centrally.",
              "Audit logs and business logs can be managed in separate layers from the control plane."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The collection layer ingests logs from agents, sidecars, and gateways.",
              "The processing layer handles parsing, cleanup, desensitization, and index generation.",
              "The query layer provides full-text search, field filtering, and context correlation."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Log Agent: collects files, stdout, and remote log streams.",
              "Pipeline Engine: performs parsing, desensitization, and routing.",
              "Query Console: provides search, aggregation, and trace-correlation views."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "Collection agents send logs into processing pipelines for structuring and desensitization.",
              "Query interfaces reconstruct context by service, Trace ID, or indexed fields."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Store hot and cold logs in separate tiers to balance query performance and cost.",
              "Audit logs should use independent storage and independent access control."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Query nodes and index nodes should be deployed separately to avoid contention.",
              "For high-throughput ingestion, use buffering queues and batched writes."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Deploy the log collection agent.",
              "Configure standard fields and collection rules for the application.",
              "Search logs by Trace ID or service name in the query layer."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Emit a standard structured log and confirm that fields are parsed correctly.",
              "Use Trace ID to verify that trace-related logs can be correlated end to end."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Standardize log field naming to avoid index fragmentation.",
              "Sensitive fields should be desensitized in the collection pipeline rather than during query time."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Index high-cardinality fields carefully to avoid runaway storage and query cost.",
              "Configure retention windows in layers based on audit, troubleshooting, and compliance requirements."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose interfaces for log ingestion, search, index management, and lifecycle configuration."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK injects standard fields and bridges log context across application code.",
              "Use it together with the Trace SDK so both signals share the same context fields."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on collection success rate, index latency, query latency, and storage growth."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for collection loss, index backlog, and hotspot query fields.",
              "Watch for storage expansion, hot-cold migration issues, and failed desensitization rules."
            ]
          }
        ]
      }
    }
  },
  stellcon: {
    overview: {
      tagline:
        "A metrics platform for collection, aggregation, query analysis, and capacity dashboards.",
      goal:
        "Stellcon provides a unified metrics system for applications, gateways, middleware, and infrastructure, supporting capacity estimation, SLO tracking, and alert-threshold computation.",
      boundaries: [
        "Designed for governance of runtime metrics, business metrics, and platform resource metrics.",
        "Designed for trend analysis, capacity planning, and SLO target management.",
        "It does not replace log search; it focuses on time-series collection and aggregated query."
      ],
      capabilities: [
        "Supports Counter, Gauge, Histogram, and Summary metric models.",
        "Supports unified label conventions, aggregation rules, and multidimensional queries.",
        "Supports dashboard templates, SLO calculation, and capacity-trend analysis."
      ],
      values: [
        "Creates consistent metric semantics and reduces fragmentation across systems.",
        "Provides a unified data foundation for alerting, capacity planning, and reliability governance."
      ],
      businessScenarios: ["SLO observation", "Capacity planning"],
      platformScenarios: ["Unified monitoring dashboards", "Multi-tenant metrics governance"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "Collection rules, label boundaries, and dashboard templates are configured centrally.",
              "SLO targets and aggregation standards are distributed from the control plane."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The collection layer handles scraping, pushing, and remote write ingestion.",
              "The storage layer handles compression, sharding, and long-window aggregation.",
              "The analysis layer handles alert queries, SLO computation, and trend forecasting."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Metrics Gateway: receives metrics ingestion traffic.",
              "Time Series Store: stores time-series data and aggregation results.",
              "Dashboard API: provides query, dashboard, and rule-configuration interfaces."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "Metrics enter through the gateway and are written into time-series storage.",
              "Query interfaces aggregate results into visual views by time window and label conditions."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "In production, separate hot and cold storage and split storage clusters by tenant or business domain.",
              "For high-cardinality scenarios, use label allowlists and sampling strategy."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Scale collection and query paths independently to avoid peak interference.",
              "Store long-window aggregation and raw detail in different tiers to control cost."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Integrate the metrics SDK or expose a standard scrape endpoint.",
              "Import the system's default dashboard templates.",
              "Create alerts and SLO targets for key metrics."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Check whether the application is reporting standard runtime metrics successfully.",
              "Verify that a key SLO dashboard can render valid numbers."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Keep metric labels within a governable scope to avoid unbounded growth.",
              "Design histogram bucket boundaries according to actual latency distributions."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Control high-cardinality labels through allowlists or dimension reduction.",
              "Layer retention periods and aggregation granularity according to cost and diagnostic need."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose query, aggregation, rule-management, and dashboard-template interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK injects standard business metrics and platform metadata.",
              "Inject base labels such as service name, environment, and version consistently."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on collection latency, write-loss rate, query latency, and high-cardinality risk."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for label explosion, hotspot queries, and remote-write backlog.",
              "Watch for SLO semantic drift and incorrectly configured aggregation rules."
            ]
          }
        ]
      }
    }
  },
  stellvox: {
    overview: {
      tagline:
        "An alerting platform for alert rules, deduplication, notification orchestration, and incident coordination.",
      goal:
        "Stellvox ingests metrics, logs, traces, and platform events into a unified alerting center that supports consistent triggering, convergence, and response coordination.",
      boundaries: [
        "Designed for unified alert governance across metrics, logs, traces, and platform events.",
        "Designed for notification orchestration, on-call collaboration, and closed-loop incident handling.",
        "It does not replace monitoring data sources; it acts as the decision and notification center."
      ],
      capabilities: [
        "Supports threshold alerts, anomaly detection, event alerts, and composite rules.",
        "Supports alert grouping, suppression, silence windows, and escalation strategy.",
        "Supports multi-channel notification, on-call scheduling, and event backtracking."
      ],
      values: [
        "Converges multi-source anomalies into governable incidents.",
        "Reduces the impact of alert storms on on-call efficiency and troubleshooting speed."
      ],
      businessScenarios: ["Coordinated incident handling", "On-call rotation"],
      platformScenarios: ["Unified platform alert center", "Multi-source event convergence governance"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The rule layer manages strategy configuration, versioning, and progressive activation.",
              "Notification templates, on-call groups, and escalation paths are managed centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The engine layer performs real-time computation, deduplication, convergence, and notification routing.",
              "The collaboration layer handles ticket linkage, on-call rotation, and disposal tracking."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Alert Rule Engine: executes alert rules and convergence logic.",
              "Notification Hub: integrates email, IM, SMS, and webhook channels.",
              "On-call Center: manages scheduling, escalation, acknowledgement, and closure."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "The rule engine aggregates multi-source events and deduplicates, suppresses, and groups them.",
              "The notification center sends incidents to the proper owners according to on-call strategy."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Deploy event computation separately from notification delivery to improve resilience during alert peaks.",
              "For multi-region operation, support both local notification paths and a global aggregation channel."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Notification delivery should support retry, rate limiting, and multi-channel fallback.",
              "On-call configuration and rule configuration both need versioning and rapid rollback."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Connect metrics, log, and trace event sources.",
              "Create key alert rules and notification policies.",
              "Configure on-call groups and escalation paths."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Trigger a test alert and confirm grouping and the notification chain work.",
              "Verify that acknowledgement, silence, and escalation flows behave as expected."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Alert-convergence windows should be tuned by failure type to avoid storm-style notification.",
              "Critical alerts should be linked to a runbook and responsible team."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Rule granularity and notification frequency must balance response speed against noise control.",
              "High-priority alerts should use multi-channel delivery and escalation."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose rule management, notification template, on-call management, and event-query interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK supports direct reporting of custom business events and callback handling.",
              "Business-side integrations should declare event severity, source, and ownership labels explicitly."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on alert hit rate, false-positive rate, MTTA, and notification success rate."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for notification delay, on-call rule drift, and broken escalation paths.",
              "Watch for alert storms and false-positive ratios caused by the same fault."
            ]
          }
        ]
      }
    }
  },
  stellguard: {
    overview: {
      tagline:
        "A zero-trust security platform for identity verification, access control, policy evaluation, and secure service-to-service communication.",
      goal:
        "Stellguard applies zero-trust principles to user, service, and endpoint access so that default deny, continuous verification, and least privilege can be implemented consistently across the platform.",
      boundaries: [
        "Designed for unified identity and access control across users, services, and endpoints.",
        "Designed for mTLS, context-based authorization, and risk interception across service boundaries.",
        "It does not replace business permission models; it provides a consistent security foundation."
      ],
      capabilities: [
        "Supports identity authentication for users, services, and devices.",
        "Supports fine-grained policy evaluation, dynamic authorization, and risk interception.",
        "Supports service-to-service mTLS, certificate rotation, and access auditing."
      ],
      values: [
        "Unifies identity, authorization, and communication security under the same control plane.",
        "Makes least privilege and continuous verification enforceable at the platform level."
      ],
      businessScenarios: ["Zero-trust access between services", "Audit for highly sensitive operations"],
      platformScenarios: ["Unified console authentication", "Platform-level continuous-verification governance"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The identity layer handles principal registration, credential issuance, and identity mapping.",
              "The policy layer makes context-based authorization decisions and risk assessments."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The enforcement layer applies interception uniformly across gateways, sidecars, and SDKs.",
              "Service-to-service communication is validated through mTLS, tokens, and context attributes."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Identity Provider: manages identities, tokens, and certificates.",
              "Policy Engine: computes authorization strategy and risk decisions.",
              "Enforcement Point: enforces access control at gateways and services."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "The identity provider issues credentials and shares identity context with the policy engine.",
              "The enforcement point intercepts requests in real time and returns a policy decision."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Deploy core identity and policy services in isolated, highly available clusters.",
              "Access control can be enforced collaboratively across gateways, sidecars, and application code."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Identity and policy services should support partitioned deployment and local cache capability.",
              "Certificate issuance and validation paths must account for disaster recovery and rotation windows."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Register service identities and user identity sources.",
              "Configure baseline access policies and certificate issuance chains.",
              "Enable enforcement points in a gateway or sidecar."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Send a controlled access request and confirm the authentication and authorization path works.",
              "Simulate an abnormal source and verify risk interception and audit records."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Default policy should deny unknown sources and declare allowed ranges explicitly.",
              "Certificate rotation cycles should be decoupled from business release windows."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Manage service identities and user identities in separate domains to reduce permission crossover.",
              "Highly sensitive interfaces should add stricter dynamic risk-validation policy."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose identity management, policy management, token issuance, and audit-query interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK supplies auth interceptors, mTLS configuration, and context propagation.",
              "Use the same context-attribute model on both the gateway and service side."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on authentication latency, rejection rate, certificate validity, and risk-hit rate."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for certificate expiry, false-positive policy blocks, and identity-source synchronization failures.",
              "Watch for authorization-decision latency and the distribution of risk-policy hits."
            ]
          }
        ]
      }
    }
  },
  stellkey: {
    overview: {
      tagline:
        "A secret-management center for key lifecycle management, secret distribution, rotation auditing, and secure access control.",
      goal:
        "Stellkey provides unified secret custody and secret distribution for applications and platform components, reducing risk from scattered plain-text keys, difficult manual rotation, and missing audit trails.",
      boundaries: [
        "Designed for unified management of keys, certificates, tokens, and sensitive configuration.",
        "Designed for access approval, rotation audit, and short-lived credential distribution.",
        "It does not replace the configuration center directly; it provides a secure foundation for sensitive-value governance."
      ],
      capabilities: [
        "Supports unified custody of keys, certificates, tokens, and secret configuration.",
        "Supports version rotation, temporary credentials, access approval, and audit tracking.",
        "Supports coordinated distribution with the configuration center, zero-trust platform, and gateways."
      ],
      values: [
        "Unifies secret lifecycle management and reduces the risk of plain-text exposure.",
        "Makes sensitive-configuration delivery and access auditing available as platform-level capabilities."
      ],
      businessScenarios: ["Database credential custody", "TLS certificate management"],
      platformScenarios: ["Sensitive-value references from the configuration center", "Platform-level secret lifecycle governance"]
    },
    pages: {
      "summary-design": {
        sections: [
          {
            heading: "Control Plane",
            bullets: [
              "The control layer manages approval flow, rotation plans, and secret policy.",
              "Access boundaries, rotation cycles, and distribution policies are governed centrally."
            ]
          },
          {
            heading: "Data Plane",
            bullets: [
              "The storage layer handles secret encryption, version management, and access isolation.",
              "The distribution layer handles short-lived credential issuance, retrieval cache, and expiration invalidation."
            ]
          }
        ]
      },
      architecture: {
        sections: [
          {
            heading: "Component Model",
            bullets: [
              "Secret Store: stores secrets, certificates, and version records.",
              "Rotation Engine: performs rotation, revocation, and expiry reminder.",
              "Access Gateway: centralizes authentication, audit, and distribution control."
            ]
          },
          {
            heading: "Interaction Flow",
            bullets: [
              "Applications obtain short-lived credentials or secret references through the access gateway.",
              "The rotation engine triggers rotation by policy and synchronizes version state updates."
            ]
          }
        ]
      },
      deployment: {
        sections: [
          {
            heading: "Deployment Topology",
            bullets: [
              "Deploy core storage separately from the access gateway to reduce exposure.",
              "For high-sensitivity workloads, use tiered custody for root keys and business keys."
            ]
          },
          {
            heading: "Availability Strategy",
            bullets: [
              "Core secret storage should support replication and disaster-recovery restoration.",
              "The access gateway should support audit caching and fast synchronization when permissions are revoked."
            ]
          }
        ]
      },
      "quick-start": {
        sections: [
          {
            heading: "Setup",
            ordered: [
              "Create a secret namespace and access policy.",
              "Write the first secret version and bind it to a rotation plan.",
              "Use the SDK inside the application to pull short-lived credentials or secret references."
            ]
          },
          {
            heading: "Validation",
            bullets: [
              "Retrieve a test secret and confirm that an audit record is written.",
              "Trigger a rotation task and confirm the application can obtain the new version."
            ]
          }
        ]
      },
      configuration: {
        sections: [
          {
            heading: "Baseline Configuration",
            bullets: [
              "Use different rotation cycles for highly sensitive secrets and ordinary secrets.",
              "Applications should keep only references and short-lived cache, avoiding long-lived plain-text persistence."
            ]
          },
          {
            heading: "Production Guidance",
            bullets: [
              "Approval flows should be layered by secret class, environment, and access source.",
              "All sensitive values should be passed by reference to avoid copy-based spread."
            ]
          }
        ]
      },
      "api-and-sdk": {
        sections: [
          {
            heading: "API Surface",
            bullets: [
              "Expose secret management, version management, approval, rotation, and audit interfaces."
            ]
          },
          {
            heading: "SDK Guidance",
            bullets: [
              "The SDK supports secret retrieval, cache refresh, and invalidation listeners.",
              "Applications should have fallback behavior for refresh failures both at startup and runtime."
            ]
          }
        ]
      },
      observability: {
        sections: [
          {
            heading: "Metrics",
            bullets: [
              "Focus on secret access volume, rotation success rate, expiry risk, and audit anomalies."
            ]
          },
          {
            heading: "Operational Watchpoints",
            bullets: [
              "Watch for abnormal access, approval backlog, and rotation failure rate.",
              "Watch for certificate expiry, cache invalidation issues, and multi-environment reference drift."
            ]
          }
        ]
      }
    }
  }
};



