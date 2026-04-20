# Stellar Axis 基础应用环境变量规范

本文档定义 `Stellar Axis（星轴）` 体系下所有中间件、所有业务应用统一遵循的基础应用环境变量规范。

本规范用于解决以下问题：

- 平台统一注入应用身份与部署拓扑信息
- 所有中间件自动识别当前应用名称、版本、实例和环境
- 业务应用无需在每个 SDK 中重复传递相同元数据
- Kubernetes、物理机、虚拟机与混合部署场景下保持统一语义

## 规范定位

本规范是全局基础运行规范的一部分。

适用对象如下：

- 所有业务服务
- 所有中间件 SDK
- 所有 Agent、Sidecar、Collector
- 所有平台注入系统
- 所有部署脚本与镜像装机脚本

## 命名原则

- 全局统一前缀使用 `STELLAR_`
- 全部使用大写英文与下划线命名
- 变量语义必须长期稳定，不允许在单个产品内重定义
- 基础元数据与产品级配置分层，不混用命名空间

## 优先级原则

所有中间件统一建议采用以下优先级：

1. 代码显式配置
2. 产品级环境变量
3. 全局 `STELLAR_*` 环境变量
4. 中间件默认值

说明如下：

- `STELLAR_*` 是平台统一注入的全局基础元数据
- 产品级变量只允许覆盖本产品局部行为
- 不允许产品级变量改变全局基础元数据的定义

## 必选基础变量

以下变量建议作为企业级运行环境的最小必选集合：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_APP_NAME` | `user-service` | 应用或服务名称 |
| `STELLAR_APP_NAMESPACE` | `stellar.trade` | 逻辑业务域或应用命名空间 |
| `STELLAR_APP_VERSION` | `1.4.2` | 当前发布版本 |
| `STELLAR_APP_INSTANCE_ID` | `user-service-7f6d9d6d7b-2xk9p` | 当前运行实例唯一标识 |
| `STELLAR_ENV` | `dev` / `test` / `prod` | 部署环境 |

## 推荐拓扑变量

以下变量建议由平台统一注入，供日志、链路、指标、配置、网关等能力共同使用：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_CLUSTER` | `cluster-sh-prod-01` | 集群标识 |
| `STELLAR_REGION` | `cn-east-1` | 区域标识 |
| `STELLAR_ZONE` | `cn-east-1a` | 可用区标识 |
| `STELLAR_IDC` | `sh-a` | 机房或园区标识 |
| `STELLAR_HOST_NAME` | `node-01` | 主机名 |
| `STELLAR_HOST_IP` | `10.10.0.11` | 主机 IP |

## Kubernetes 推荐变量

如果运行在 Kubernetes 中，建议额外统一注入以下变量：

| 环境变量 | 示例 | 说明 |
| :--- | :--- | :--- |
| `STELLAR_NODE_NAME` | `worker-node-01` | 节点名称 |
| `STELLAR_K8S_NAMESPACE` | `trade` | Kubernetes Namespace |
| `STELLAR_POD_NAME` | `user-service-7f6d9d6d7b-2xk9p` | Pod 名称 |
| `STELLAR_POD_IP` | `172.20.10.23` | Pod IP |
| `STELLAR_CONTAINER_NAME` | `app` | 容器名称 |

## 统一语义约束

- `APP_NAME`
  表示业务应用或工作负载的稳定名称，不使用主机名、Pod 名或镜像名替代
- `APP_NAMESPACE`
  表示逻辑业务域，不等同于 Kubernetes Namespace
- `APP_VERSION`
  表示本次发布版本，不使用 Git 分支名替代
- `APP_INSTANCE_ID`
  表示单个运行实例唯一标识，允许使用 Pod 名、实例 ID 或平台生成值
- `ENV`
  只表示部署环境，不混入区域、租户、集群或机房信息
- `CLUSTER / REGION / ZONE / IDC`
  只表示基础设施拓扑，不与环境语义混用

## 产品级扩展约束

在 `STELLAR_*` 之外，各产品可以保留自己的产品级环境变量，例如：

- 日志平台：`STELLSPEC_*`
- 链路平台：`STELLTRACE_*`
- 配置中心：`STELLNULA_*`
- 服务治理：`STELLORBIT_*`

但产品级变量只允许描述：

- 本产品的接入地址
- 本产品的协议与格式
- 本产品的采样、级别、超时、批量参数
- 对基础元数据的局部覆盖

不允许产品级变量重新发明：

- 应用名语义
- 环境语义
- 区域语义
- 集群语义

## 平台注入建议

推荐统一由平台在如下阶段注入：

- Kubernetes Admission Webhook
- Helm / Kustomize 模板
- 宿主机装机脚本
- 容器运行时入口脚本
- 发布平台注入流程

推荐来源包括：

- `metadata.name`
- `metadata.namespace`
- `spec.nodeName`
- `status.podIP`
- `status.hostIP`
- 发布系统生成的版本号
- CMDB / 资源平台中的区域、机房、集群信息

## 对中间件 SDK 的约束

所有中间件 SDK 应默认遵循本规范：

- 默认读取 `STELLAR_*` 基础元数据
- 在业务方未显式配置时自动补全应用身份
- 将基础元数据映射到各自的资源模型、日志字段、链路属性与指标标签
- 允许产品级环境变量做本地覆盖，但不得破坏全局语义
