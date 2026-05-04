---
title: 分布式系统中的一致性挑战及其解决路径
category: 分布式理论
summary: 从 FLP、CAP、PACELC 到 Paxos、Raft、Zab、PBFT，再到 Spanner、etcd 与 KRaft，系统梳理一致性问题的理论边界与工程落地。
tags:
  - 分布式一致性
  - 共识协议
  - Paxos
  - Raft
readingDirection: 适合在补分布式系统理论基础，或想把共识协议和工业系统实现串起来时阅读。
---

# 分布式系统中的一致性挑战及其解决路径

## ——从 FLP、Paxos 到 Spanner、etcd 与 KRaft

### 摘要

分布式系统的发展，本质上是把“单机内的确定性”换成了“多机、多副本、网络不确定条件下的可扩展性与容错性”。GFS 和 Bigtable 这类早期工业系统已经明确表明，互联网级系统首先面对的是廉价机器、频繁故障、海量数据与跨机架通信，而不是教科书式的理想环境；随之而来的核心问题不是“怎么把数据复制几份”，而是“复制之后，谁是合法主节点、谁的写入有效、不同副本何时算一致、跨分片事务如何保持正确”。([research.google.com][1])

讨论分布式一致性时，最常见也最不严谨的错误，是把“副本共识”“事务提交”“客户端可见一致性”混为一谈。严格说，共识算法解决的是复制状态机/复制日志顺序的一致；事务系统还必须叠加锁、时间戳、成员变更、租约或 2PC 才能给出 serializability、strict serializability 或 external consistency 这类更强语义。Spanner/F1 明确是在 Paxos 之上再做 2PC 与事务控制；而 etcd 则把 KV API 直接做到 strict serializability。([research.google.com][2])

本文的核心结论是：分布式一致性并不是被某一个“神算法”一次性解决的，而是沿着“先承认不可能，再在可行假设下保证安全，最后把理论做成工程”的路径逐层解决的。理论边界由 FLP、CAP 与 PACELC 划定；安全复制由 Viewstamped Replication、Paxos、Raft、Zab、PBFT 等协议建立；真正落地则依赖 Chubby、ZooKeeper、Dynamo、Megastore、Percolator、Spanner、F1、etcd、TiDB/TiKV、Kafka KRaft 等工程系统的长期验证。([MIT CSAIL][3])

### 关键词

分布式系统；一致性；共识算法；Paxos；Raft；Spanner；脑裂；强一致性；弱一致性

---

## 一、背景：为什么分布式系统一定会撞上一致性问题

分布式系统不是“把单机程序部署到很多机器上”那么简单。GFS 论文明确指出，其目标是在廉价通用硬件上为大规模数据密集型应用提供容错与高吞吐；Bigtable 进一步展示了这种架构已经在 Google 内部形成大规模生产实践。也就是说，系统一旦进入多副本、多节点、多网络路径环境，局部故障就会成为常态，而不是例外。([research.google.com][1])

更麻烦的是，现实网络不是“断了”与“通了”两种状态，而可能表现为高延迟、部分丢包、单向阻塞、抖动甚至时钟漂移。Google SRE 明确把网络分区类问题描述为：看起来像完全分区的故障，实际上可能只是网络很慢、部分消息丢失，或某一个方向被限流。正因为网络故障不是二元的，系统几乎不可能靠简单心跳把“对方宕机”和“对方只是很慢”可靠地区分出来。([sre.google][4])

这就是一致性问题的起点：当多个副本都可能继续运行，而网络又不能及时给出真相时，系统必须回答三个残酷问题：第一，谁有资格当主；第二，谁的日志顺序算准；第三，客户端看到的数据语义到底是什么。真正成熟的分布式系统，都是围绕这三个问题设计出来的。([sre.google][4])

---

## 二、问题 1：分布式系统在演进过程中出现了哪些问题与挑战

### 1. 部分故障与网络不确定性

单机系统通常只有“正常/崩溃”两种视角，而分布式系统里最难受的是部分故障：某些节点活着、某些链路慢着、某些消息丢着、某些副本还在回放日志。FLP 证明告诉我们，在纯异步模型里，只要允许哪怕一个进程故障，就不存在一个既总能终止又总正确的确定性共识算法。CAP 则进一步说明，在异步网络模型下，不可能同时在所有公平执行中保证 Availability 与 Atomic Consistency。Abadi 的 PACELC 又把问题说得更接地气：即便没有分区，系统也仍然要在一致性与延迟之间权衡。([MIT CSAIL][3])

### 2. 领导者选举、日志复制与副本收敛

只要系统采用主从或 leader/follower 复制，就一定会面临领导者失效、日志乱序、旧主复活、重复提交、落后副本追赶等问题。Viewstamped Replication 很早就把“主副本崩溃后如何重组选主并继续复制”讲清楚；Paxos 解决的是如何在副本之间就某个值达成一致；Raft 则把 leader election、log replication、safety、membership changes 直接拆开讲清楚，这就是它后来大规模流行的根本原因。([普林斯顿大学计算机科学系][5])

### 3. 成员变更与扩缩容

真正的生产系统不会一直维持固定副本集合。节点要下线、扩容、迁移机房、替换硬件，副本集合就必须变。难点在于：成员变更过程中不能同时让旧配置和新配置都“各自合法”，否则会直接引出脑裂。Raft 的 joint consensus 给出的就是工业上非常正确的答案：过渡期同时要求旧配置多数与新配置多数，从而保证安全。这个设计非常重要，因为它回答的不是“怎么加节点”，而是“怎么在加节点时不破坏一致性”。([Raft][6])

### 4. 跨分片、跨复制组事务

共识只能保证单个复制日志的顺序一致，但现实数据库很快就会遇到跨分片事务。Percolator 明确指出 Bigtable 本身不提供多行事务，因此它要在其之上补上 multirow transactions；Spanner 则更进一步，把每个数据目录放到 Paxos group 中，再在跨组事务上叠加 2PL 和 2PC；F1 继续在 Spanner 之上构建分布式 SQL。很多工程文档把 2PC 吹成“一致性算法”，这是不严谨的：2PC 解决原子提交，不解决在异步故障环境下的副本共识，所以它必须站在 Paxos/Raft 这类复制一致性底座上。([USENIX][7])

### 5. 全局时间与可见顺序

当系统跨地域部署时，仅靠“消息先到先提交”不够，因为真实世界的因果顺序、提交顺序和客户端观察顺序可能脱钩。Spanner 的突破不只是“用了 Paxos”，而是把 TrueTime 引入事务时间戳分配，使其能够提供 external consistency，并且官方文档明确说明 external consistency 比 linearizability 更强。换句话说，全球强一致数据库不是只靠共识算法撑起来的，还要靠时间不确定性的显式建模。([research.google.com][2])

### 6. 一致性、可用性与延迟的工程权衡

不是所有业务都应该追求最强一致。Dynamo 的设计目标就是“always-on”，因此它在某些故障场景下牺牲一致性以换取可用性，并把冲突解决的一部分责任交给应用。Cassandra 官方文档明确承认其具有 eventually consistent semantics，并通过 consistency levels 让用户在每次操作上做 `R + W > N` 类权衡。这条路线没有错，但它是把复杂度从数据库内核转移给了业务侧。([All Things Distributed][8])

---

## 三、问题 2：这些挑战分别被哪些论文解决了

先说结论：有些问题不是“被解决”，而是“被证明不可能在理想模型下彻底解决”。FLP 说明纯异步 + 故障下无法保证确定性共识总能推进；Gilbert 与 Lynch 形式化了 CAP；Abadi 提出 PACELC，指出即使没有分区，系统也要在一致性与延迟间抉择。工业界后来不是推翻这些结果，而是接受它们的边界。([MIT CSAIL][3])

在“如何让系统最终活起来”这个问题上，Dwork、Lynch、Stockmeyer 的 partial synchrony 模型，以及 Chandra、Toueg 的 failure detectors，是理论转向工程的关键一步。它们的价值不在于让世界变同步，而在于告诉我们：只要承认现实系统通常“多数时候近似同步”，并给出足够可靠的失败怀疑机制，就能在不放弃安全性的前提下获得活性。Chubby 论文对这一点说得非常实在：Paxos 保安全不依赖时间假设，但活性需要时钟。([MIT CSAIL][9])

在崩溃故障模型下，Viewstamped Replication、Paxos 和 Raft 是三条最核心的主线。VR 是早期主副本复制的一套完整答案；Paxos 给出了通用共识框架；Raft 在等价于 multi-Paxos 的前提下，通过强 leader、随机选举、joint consensus 等设计，大幅降低了理解和实现门槛。工业上真正跑赢的，不是“最花哨”的论文，而是那些把安全性、可实现性和可运维性平衡得最好的协议。Raft 爆红，不是因为它比 Paxos 更“高深”，而是因为它更适合人类团队写对。([普林斯顿大学计算机科学系][5])

在性能优化方向，Fast Paxos 试图把学习值的消息延迟从三跳降到两跳；Mencius 面向 WAN 强一致复制，强调负载均衡与跨地域低延迟；EPaxos 追求 wide-area 下的一跳提交与无单点 leader 热点；Flexible Paxos 则重新审视 quorum 交叉条件，指出两阶段都用严格多数并非必要。这些论文都很重要，但要实话实说：它们对工业界的影响更多体现在“思想输入”和特定场景优化，而不是像 Raft 那样直接成为云原生基础设施的标准教材。([Microsoft][10])

在协调服务方向，Zab 和 ZooKeeper 解决的是“如何把一致性复制真正变成大家能用的协调内核”。ZooKeeper 论文把它定义为互联网规模系统的协调服务；Zab 则提供面向 primary-backup 的高性能原子广播，并在论文中给出生产经验与性能数据。这个方向的意义在于：它没有把每个业务都逼成“自己实现 Paxos”，而是抽象出锁、命名、配置、会话、watch 这些高频协调能力。([USENIX][11])

在弱一致性路线，Dynamo 是决定性论文。它明确提出，为了 always-on，可在某些故障场景下牺牲强一致，转而使用版本向量、异步复制与应用辅助冲突解决，并证明 eventually consistent storage system 可以在严苛生产环境中工作。这条路线后来被 Cassandra 工程化地继承。([All Things Distributed][8])

在全球强一致事务路线，Megastore、Percolator、Spanner、F1 构成了清晰演进链。Megastore 把强一致和高可用结合在细粒度分区内；Percolator 在 Bigtable 之上补足多行事务；Spanner 首次把 externally-consistent distributed transactions 做到全球规模；F1 则把分布式 SQL 建立在 Spanner 之上，并给出生产化业务验证。这个序列真正回答了“共识算法之后，数据库一致性还差什么”——答案是事务、锁、时间戳和分片布局。([cidrdb.org][12])

在更强的故障模型上，PBFT 则给出了拜占庭容错的实用方案。它证明了在异步环境中，面对恶意或任意行为节点，也可以把复制服务做到可用，并给出实现与性能数据。BFT 在传统企业数据库里不是主流，但在区块链、跨组织系统和高对抗环境中意义很大。([css.csail.mit.edu][13])

---

## 四、问题 3：这些问题被哪些工程实践验证过

最早的成熟工程答案之一是 Google 的 Chubby。Chubby 论文明确承认它不是新算法研究，而是把 Paxos 这类共识思想做成工程服务；更关键的是，GFS 和 Bigtable 都使用 Chubby 来进行 master 选举、服务发现和元数据定位。这说明工业界最先验证的，不是“让每个业务都实现共识”，而是“先把共识封装成协调基础设施”。([research.google.com][14])

Yahoo 的 ZooKeeper/Zab 路线则证明了另一件事：协调服务可以成为互联网级基础设施，并且可以用更高层 API 支撑锁、配置、命名、会话与 watch。Zab 论文还给出生产经验，表明其实现能达到每秒数万次广播，足以支撑 Web 规模应用。([USENIX][11])

Amazon 的 Dynamo 路线证明了弱一致并不是学术妥协，而是真正的生产策略。论文直接说明 Dynamo 是 some of Amazon’s core services 的 highly available key-value store，用于 always-on 场景；Cassandra 官方文档则明确表明自己继承了 Dynamo 与 Bigtable 的思路，采用 eventually consistent semantics，并支持按操作调节一致性级别。这条路线后来在社交、消息、购物车、日志和可容忍短暂陈旧读的业务中非常长寿。([All Things Distributed][8])

Google 的强一致事务路线验证得更彻底。Percolator 已被用于构建 Google 搜索索引，并把搜索结果中文档平均年龄降低了 50%；Spanner 论文说明其是全球同步复制、支持 externally-consistent distributed transactions 的数据库；F1 论文则明确写到它自 2012 年初开始管理全部 AdWords 广告活动数据。这说明“全球强一致数据库做不成生产”这种说法，早就被现实打脸了，只是代价是系统设计必须非常重。([USENIX][7])

云原生时代的代表验证是 etcd 与 Kubernetes。etcd 官方定义自己是 strongly consistent distributed key-value store，并说明其可以在网络分区中处理 leader 选举；Kubernetes 官方则明确指出 etcd 是所有集群数据的 backing store。换句话说，今天最主流的容器控制面，本质上就是建立在 Raft 风格一致性存储之上的。([etcd][15])

数据库工程里，TiDB/TiKV 是非常典型的 Raft 落地。官方文档说明 TiKV 以 Region 为单位做 Raft replication 与 membership management，每个 Region 的多个副本构成一个 Raft Group；TiDB 提供 ACID transaction，适合银行转账这类强一致场景，其事务模型又明确借鉴并优化了 Google Percolator。这个组合说明，现代分布式数据库的主流路线已经是“分片 + 多个共识组 + MVCC + 分布式事务”。([PingCAP Docs][16])

消息系统方面，Kafka 已经完成从 ZooKeeper 模式到 KRaft 模式的转型。Kafka 官方文档明确说明，在 KRaft 模式下，controller 使用 Raft 协议管理元数据；Kafka 4.0 公告则说明这是首个完全不依赖 ZooKeeper 的主版本发布，KRaft 成为默认模式。这说明共识算法已经不再只是“数据库的事”，而是整个分布式基础设施的通用内核。([kafka.apache.org][17])

---

## 五、问题 4：一致性共识算法有哪些，它们如何演进

从主线看，一致性共识算法可以分为两大类：崩溃故障容错（CFT）和拜占庭容错（BFT）。工业主流几乎都在 CFT 线上，因为大多数数据中心故障来自宕机、重启、网络抖动、磁盘损坏，而不是恶意节点；PBFT 这类 BFT 更多用于高对抗环境。([css.csail.mit.edu][13])

CFT 的历史脉络大致是这样的：1988 年的 Viewstamped Replication 先把 primary-backup 复制与主故障恢复讲清楚；1998/2001 年 Paxos 形成经典共识框架；2005/2006 年 Fast Paxos 追求更低消息延迟；2008 年 Mencius 面向 WAN 负载均衡；2011 年 Zab 把原子广播用于协调服务；2013 年 EPaxos 尝试降低 wide-area 提交延迟与 leader 热点；2014 年 Raft 通过强 leader 和分治式设计大幅改善可理解性；2016 年 Flexible Paxos 则证明传统多数派并非唯一 quorum 选择。这个演进的本质，不是安全性从无到有，而是从“能证明对”逐步走向“更快、更好实现、更适合运维”。([普林斯顿大学计算机科学系][5])

我对这条演进链的判断很明确：Paxos 在理论地位上无可替代，Raft 在工程传播上胜出，Zab 在协调服务场景里足够务实，而 Fast/EPaxos/Flexible Paxos 这类工作更多代表“优化方向”。真正大面积进入工业标准栈的，是 Paxos family、Raft 和 Zab，不是所有论文里的变体。这个判断并不保守，而是符合现实采用面。([Raft][6])

需要特别强调的一点是：2PC 不是共识算法。它解决“多个参与者是否共同提交”的原子提交问题，但自身不能在异步网络故障中提供副本共识，所以 Spanner 和 F1 都明确是把 2PC 放在 Paxos 之上。把 2PC 和 Paxos/Raft 混为一谈，是分布式系统里最常见的概念性错误之一。([research.google.com][2])

---

## 六、问题 5：一致性共识算法的落地工程实践有哪些，它又如何演进

工程落地的第一阶段，是“把共识做成一个独立协调服务”。Chubby 和 ZooKeeper 都属于这一阶段：业务系统不直接实现共识，而是通过锁、命名、配置、leader election 这些高层能力间接获得一致性。这个阶段的优点是复用高，缺点是协调服务会成为基础设施里的关键依赖。([research.google.com][14])

第二阶段，是“把共识内嵌进存储系统本体”。Megastore、Spanner、F1 代表的是 Paxos 内嵌数据库核心；TiKV 则把每个 Region 都变成一个 Raft Group。这一步的意义非常大：一致性不再只是控制面的锁服务，而是直接参与数据页、事务、索引和分片迁移的正确性维护。([cidrdb.org][12])

第三阶段，是“从单个共识组走向大规模分片共识”。Spanner 用多个 Paxos group 管目录；TiKV 用多个 Region/Raft Group；Kubernetes 则把集群状态统一放进 etcd；Kafka 又把元数据一致性从外部 ZooKeeper 收回到内部 KRaft controller quorum。演进方向非常清楚：从“集中协调器”到“系统内部处处是共识分组”，从“局部主从”到“全局元数据、事务与调度都靠共识维护”。([research.google.com][2])

第四阶段，是“在共识之上叠加更细的工程机制”。例如 Chubby 用锁获取计数防止延迟包污染新主写入；Raft 用 joint consensus 处理安全成员变更；etcd 通过多数派配置防止 split-brain；Spanner 用 TrueTime 支撑 external consistency；TiDB 用 MVCC 与 Percolator 风格事务补全数据库语义。共识算法从来都不是完整系统，它只是底座。真正可靠的系统，靠的是底座之上的那一层层约束。([research.google.com][14])

---

## 七、问题 6：强一致性和弱一致性的区别、优点和缺点

强一致性不是一个单一名词，而是一组从弱到强的语义：对象层面常说 linearizability；事务系统里常说 strict serializability 或 external consistency。etcd 官方写得很清楚：其 KV API 保证 strict serializability，并默认提供 linearizability；Spanner 官方则说明自己默认提供 external consistency，而且这是比 linearizability 更强的属性。把这些统称为“强一致性”没问题，但若在论文里不区分层次，就是概念偷懒。([etcd][18])

强一致性的优点非常直接：系统更容易推理，读到旧值、读到不存在过的中间状态、跨对象顺序颠倒这类问题会显著减少。etcd 文档甚至明确说 linearizability makes reasoning easily；Spanner 进一步说明 strong reads 能看到操作开始前已提交的所有事务效果。这类语义对金融、元数据、调度、权限、控制面尤其重要。([etcd][18])

强一致性的代价同样直接：写入通常需要走 quorum 或完整共识流程，读也可能要走线性化路径；一旦多数派不可达，就必须牺牲部分可用性。etcd 文档说明 linearized requests 必须经过 Raft，而切到 serializable read 则可降低延迟但可能读到 stale data；它在多数派故障下也明确不能继续接受写入。这个代价不是实现不好，而是理论边界决定的。([etcd][18])

弱一致性则是一个更宽的伞形概念，常见代表是 eventual consistency 与 tunable consistency。Dynamo 直接承认自己为了 always-on 在某些故障场景下牺牲一致性；Cassandra 官方明确使用 eventually consistent semantics，并通过 consistency levels 提供按操作权衡。优点是低延迟、高可用、跨地域更灵活；缺点是冲突解决、修复、业务补偿、幂等控制会被推给应用层，最终复杂度并没有消失，只是转移了。([All Things Distributed][8])

我的判断是：关键状态必须默认强一致，尤其是元数据、锁、租约、主节点身份、账户余额、资源调度；弱一致适合高吞吐、可容忍短暂陈旧、冲突可合并的业务，如购物车、推荐画像、社交 feed、统计指标。把所有业务都做成强一致，是过度设计；把关键状态做成弱一致，则通常是事故预告。([sre.google][4])

---

## 八、问题 7：网络不确定引发的脑裂问题，会带来哪些麻烦

Google SRE 给了非常直白的案例：两个副本节点仅靠心跳和超时来做主备切换时，一旦网络变慢或丢包，双方都可能认为对方失联，于是都提升自己为主；也可能双方都互相“打死”，最后双主或双挂。结果只有两种：数据损坏，或者服务不可用。SRE 明确指出，用简单心跳解决 leader election 在概念上就是不成立的。([sre.google][4])

脑裂最直接的麻烦有四类。第一，双主写入，导致日志分叉、元数据冲突、最后只能人工对账或强制回滚；第二，旧主复活后继续接受旧租约或旧锁，造成重复调度、重复执行或覆盖写；第三，组成员视图不一致，两个分区各自选主并接收写删请求，最终造成不可恢复的数据偏差；第四，为了避免脑裂，系统可能反过来自我阉割可用性，让主节点在无法确认对端状态时直接停写并等待人工介入。SRE 的三个案例把这些麻烦基本说透了。([sre.google][4])

正确的工程解法从来不是“更勤快地发心跳”，而是多数派法定人数、持久化日志、显式成员变更和栅栏机制。etcd 官方明确说明，网络分区后只有多数派一侧可用，少数派不可用；若 leader 落在少数派，它会主动 step down，且系统不会发生 split-brain。ZooKeeper 官方则强调部署奇数副本，6 台的容错收益并不比 5 台更高，而奇数规模有利于避免 brain split 问题。([etcd][19])

此外，还需要“旧主不能凭旧身份继续写”的栅栏思想。Chubby 论文给出的工程技巧非常经典：主节点获得锁后，在写 RPC 中附带锁获取计数，接收方拒绝比当前计数更小的请求，以抵御延迟包和旧主残留写入。本质上，这就是今天广泛使用的 fencing token 思想。没有这层防线，即便你完成了重新选主，旧主的迟到写入仍然会污染系统。([research.google.com][14])

所以，脑裂不是“偶发 Bug”，而是分布式系统的基础宿命。能不能避免，不取决于运气，而取决于系统是否从一开始就把主节点身份、成员视图、日志顺序和写入资格都绑定在一个经过证明的共识协议上。靠脚本切主、靠心跳拉活、靠人工拍脑袋决策，迟早会出事故。([sre.google][4])

---

## 九、结论

分布式一致性问题的真正解决路径，可以概括为四层。

第一层，承认边界。FLP 说明纯异步下无法保证确定性共识总有活性；CAP 与 PACELC 说明系统不可能白拿一致性、可用性与低延迟。([MIT CSAIL][3])

第二层，保证安全。Viewstamped Replication、Paxos、Raft、Zab、PBFT 等协议通过 quorum、日志复制与状态机复制，把“多个副本看到同一顺序”这件事做对。([普林斯顿大学计算机科学系][5])

第三层，补足活性与事务语义。partial synchrony、failure detectors、leader lease、joint consensus、2PC、MVCC、TrueTime 这些机制，把“副本顺序一致”扩展成“数据库事务正确、全局时序可解释、成员变更不脑裂”。([MIT CSAIL][9])

第四层，工程化验证。Chubby、ZooKeeper、Dynamo、Megastore、Percolator、Spanner、F1、etcd、TiDB/TiKV、Kafka KRaft 证明了这些方案不只是论文里的美学，而是生产级基础设施。([research.google.com][14])

一句话收束全文：**分布式一致性不是被“消灭”的，而是被“约束”出来的。真正靠得住的办法，从来都是多数派、持久化日志、正确的成员变更、事务控制与明确的一致性语义；靠心跳、脚本和侥幸心理维护的一致性，根本不配叫方案。** ([sre.google][4])

---

## 参考文献（核心论文与官方资料）

1. Fischer, Lynch, Paterson. *Impossibility of Distributed Consensus with One Faulty Process* (1985). ([MIT CSAIL][3])
2. Gilbert, Lynch. *Brewer’s Conjecture and the Feasibility of Consistent, Available, Partition-Tolerant Web Services* (2002). ([普林斯顿大学计算机科学系][20])
3. Abadi. *Consistency Tradeoffs in Modern Distributed Database System Design: CAP is Only Part of the Story* (PACELC, 2012). ([UMD计算机科学系][21])
4. Oki, Liskov. *Viewstamped Replication: A New Primary Copy Method to Support Highly-Available Distributed Systems* (1988). ([普林斯顿大学计算机科学系][5])
5. Dwork, Lynch, Stockmeyer. *Consensus in the Presence of Partial Synchrony* (1988). ([MIT CSAIL][9])
6. Chandra, Toueg. *Unreliable Failure Detectors for Reliable Distributed Systems* (1996). ([cs.utexas.edu][22])
7. Lamport. *The Part-Time Parliament* (1998)；*Paxos Made Simple* (2001). ([Leslie Lamport's Home Page][23])
8. Castro, Liskov. *Practical Byzantine Fault Tolerance* (1999). ([css.csail.mit.edu][13])
9. Lamport. *Fast Paxos* (2005/2006). ([Microsoft][10])
10. Mao, Junqueira, Marzullo. *Mencius* (2008). ([USENIX][24])
11. Junqueira, Reed, Serafini. *Zab: High-performance broadcast for primary-backup systems* (2011). ([CSE Labs][25])
12. Ongaro, Ousterhout. *In Search of an Understandable Consensus Algorithm (Raft)* (2014). ([Raft][6])
13. Moraru, Andersen, Kaminsky. *EPaxos* (2013). ([CMU School of Computer Science][26])
14. Howard, Malkhi, Spiegelman. *Flexible Paxos* (2016). ([EECS Department][27])
15. Burrows. *The Chubby lock service for loosely-coupled distributed systems* (2006). ([research.google.com][14])
16. Hunt et al. *ZooKeeper: Wait-free coordination for Internet-scale systems* (2010). ([USENIX][11])
17. DeCandia et al. *Dynamo: Amazon’s Highly Available Key-value Store* (2007). ([All Things Distributed][8])
18. Baker et al. *Megastore* (2011). ([cidrdb.org][12])
19. Peng, Dabek. *Percolator* (2010). ([USENIX][7])
20. Corbett et al. *Spanner* (2012). ([research.google.com][2])
21. Shute et al. *F1: A Distributed SQL Database That Scales* (2013). ([research.google.com][28])
22. Google Cloud Spanner 官方文档：external consistency、strong reads、stale reads。([Google Cloud Documentation][29])
23. etcd 官方文档：strict serializability、linearizability、network partition 下的 majority/minority 行为。([etcd][18])
24. Kubernetes 官方文档：etcd 是所有集群数据的 backing store。([Kubernetes][30])
25. TiDB/TiKV 官方文档：Raft replication、Region/Raft Group、Percolator 风格事务。([PingCAP Docs][16])
26. Cassandra 官方文档：eventually consistent semantics、`R + W > N` consistency levels。([Apache Cassandra][31])
27. Kafka 官方文档与公告：KRaft mode uses Raft to manage metadata；Kafka 4.0 默认不再依赖 ZooKeeper。([kafka.apache.org][17])
28. Google SRE Book：脑裂、错误主备切换与为什么心跳/脚本不是正确答案。([sre.google][4])

[1]: https://research.google.com/archive/gfs-sosp2003.pdf "https://research.google.com/archive/gfs-sosp2003.pdf"
[2]: https://research.google.com/archive/spanner-osdi2012.pdf "https://research.google.com/archive/spanner-osdi2012.pdf"
[3]: https://groups.csail.mit.edu/tds/papers/Lynch/jacm85.pdf "https://groups.csail.mit.edu/tds/papers/Lynch/jacm85.pdf"
[4]: https://sre.google/sre-book/managing-critical-state/ "https://sre.google/sre-book/managing-critical-state/"
[5]: https://www.cs.princeton.edu/courses/archive/fall09/cos518/papers/viewstamped.pdf "https://www.cs.princeton.edu/courses/archive/fall09/cos518/papers/viewstamped.pdf"
[6]: https://raft.github.io/raft.pdf "https://raft.github.io/raft.pdf"
[7]: https://www.usenix.org/legacy/event/osdi10/tech/full_papers/Peng.pdf "https://www.usenix.org/legacy/event/osdi10/tech/full_papers/Peng.pdf"
[8]: https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf "https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf"
[9]: https://groups.csail.mit.edu/tds/papers/Lynch/jacm88.pdf "https://groups.csail.mit.edu/tds/papers/Lynch/jacm88.pdf"
[10]: https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/tr-2005-112.pdf "https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/tr-2005-112.pdf"
[11]: https://www.usenix.org/legacy/event/atc10/tech/full_papers/Hunt.pdf "https://www.usenix.org/legacy/event/atc10/tech/full_papers/Hunt.pdf"
[12]: https://www.cidrdb.org/cidr2011/Papers/CIDR11_Paper32.pdf "https://www.cidrdb.org/cidr2011/Papers/CIDR11_Paper32.pdf"
[13]: https://css.csail.mit.edu/6.824/2014/papers/castro-practicalbft.pdf "https://css.csail.mit.edu/6.824/2014/papers/castro-practicalbft.pdf"
[14]: https://research.google.com/archive/chubby-osdi06.pdf "https://research.google.com/archive/chubby-osdi06.pdf"
[15]: https://etcd.io/ "https://etcd.io/"
[16]: https://docs.pingcap.com/tidb/stable/tidb-storage "https://docs.pingcap.com/tidb/stable/tidb-storage"
[17]: https://kafka.apache.org/42/getting-started/zk2kraft/ "https://kafka.apache.org/42/getting-started/zk2kraft/"
[18]: https://etcd.io/docs/v3.5/learning/api_guarantees/ "https://etcd.io/docs/v3.5/learning/api_guarantees/"
[19]: https://etcd.io/docs/v3.5/op-guide/failures/ "https://etcd.io/docs/v3.5/op-guide/failures/"
[20]: https://www.cs.princeton.edu/courses/archive/spr22/cos418/papers/cap.pdf "https://www.cs.princeton.edu/courses/archive/spr22/cos418/papers/cap.pdf"
[21]: https://www.cs.umd.edu/~abadi/papers/abadi-pacelc.pdf "https://www.cs.umd.edu/~abadi/papers/abadi-pacelc.pdf"
[22]: https://www.cs.utexas.edu/~lorenzo/corsi/cs380d/papers/p225-chandra.pdf "https://www.cs.utexas.edu/~lorenzo/corsi/cs380d/papers/p225-chandra.pdf"
[23]: https://lamport.azurewebsites.net/pubs/lamport-paxos.pdf "https://lamport.azurewebsites.net/pubs/lamport-paxos.pdf"
[24]: https://www.usenix.org/event/osdi08/tech/full_papers/mao/mao.pdf "https://www.usenix.org/event/osdi08/tech/full_papers/mao/mao.pdf"
[25]: https://classpages.cselabs.umn.edu/Fall-2017/csci8211/Papers/Distributed%20Systems%20Zab-%20High-performance%20broadcast%20for%20primary-backup%20systems.pdf "https://classpages.cselabs.umn.edu/Fall-2017/csci8211/Papers/Distributed%20Systems%20Zab-%20High-performance%20broadcast%20for%20primary-backup%20systems.pdf"
[26]: https://www.cs.cmu.edu/~dga/papers/epaxos-sosp2013.pdf "https://www.cs.cmu.edu/~dga/papers/epaxos-sosp2013.pdf"
[27]: https://web.eecs.umich.edu/~manosk/assets/papers/flexible_paxos_opodis2016.pdf "https://web.eecs.umich.edu/~manosk/assets/papers/flexible_paxos_opodis2016.pdf"
[28]: https://research.google.com/pubs/archive/41344.pdf "https://research.google.com/pubs/archive/41344.pdf"
[29]: https://docs.cloud.google.com/spanner/docs/true-time-external-consistency "https://docs.cloud.google.com/spanner/docs/true-time-external-consistency"
[30]: https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/ "https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/"
[31]: https://cassandra.apache.org/doc/4.1/cassandra/architecture/overview.html "https://cassandra.apache.org/doc/4.1/cassandra/architecture/overview.html"
