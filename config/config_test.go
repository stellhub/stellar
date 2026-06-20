package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadApplicationYMLConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "application.yml")
	content := `
app:
  name: order-service
  env: uat
  zone: zone-a
http:
  server:
    enabled: true
    port: 18080
    adapter: chi
    observability:
      trace: true
      metrics: true
      logs: true
  client:
    enabled: true
    timeout: 3s
    max_idle_conns: 100
    max_idle_conns_per_host: 10
    idle_conn_timeout: 90s
    discovery:
      enabled: true
      adapter: stellmap
      endpoints:
        - http://localhost:18090
      namespace: default
      load_balance: p2c
    observability:
      trace: true
      metrics: true
      logs: false
    clients:
      user-service:
        timeout: 2s
        discovery:
          enabled: true
          service: user-service
          protocol: http
          endpoint_name: http
      order-service:
        base_url: http://localhost:8082
        timeout: 5s
grpc:
  server:
    enabled: true
    port: 19090
    adapter: grpc-go
    observability:
      trace: true
      metrics: true
      logs: true
  client:
    enabled: true
    timeout: 3s
    insecure: true
    discovery:
      enabled: true
      adapter: stellmap
      endpoints:
        - http://localhost:18090
      namespace: default
      load_balance: p2c
    observability:
      trace: true
      metrics: true
      logs: false
    clients:
      user-service:
        timeout: 2s
        discovery:
          enabled: true
          service: user-service
          protocol: grpc
          endpoint_name: grpc
      order-service:
        target: dns:///localhost:19092
        timeout: 5s
redis:
  enabled: true
  addr: localhost:6379
  db: 1
  pool_size: 16
  dial_timeout: 2s
  read_timeout: 1s
  write_timeout: 1s
  debug_api:
    enabled: true
    prefix: /redis
  observability:
    trace: true
    metrics: true
    logs: true
mysql:
  enabled: true
  dsn: user:pass@tcp(localhost:3306)/orders?parseTime=true
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 30m
  conn_max_idle_time: 5m
  debug_api:
    enabled: true
    prefix: /mysql
  observability:
    trace: true
    metrics: true
    logs: true
postgresql:
  enabled: true
  dsn: postgres://user:pass@localhost:5432/orders?sslmode=disable
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 30m
  conn_max_idle_time: 5m
  debug_api:
    enabled: true
    prefix: /postgresql
  observability:
    trace: true
    metrics: true
    logs: true
cache:
  enabled: true
  adapter: freecache
  ttl: 5m
  clean_window: 30s
  size_bytes: 1048576
  debug_api:
    enabled: true
    prefix: /cache
  observability:
    trace: false
    metrics: true
    logs: false
mq:
  enabled: true
  adapter: stellflow
  brokers:
    - stellflow://localhost:9092
  client_id: order-service
  observability:
    trace: true
    metrics: true
    logs: false
  producer:
    enabled: true
    default_topic: orders.created
    delivery_timeout: 30s
    retry_max_attempts: 3
    auto_create_topic_partition_count: 4
    observability:
      metrics: true
  consumer:
    enabled: true
    group_id: order-worker
    topics:
      - orders.created
    auto_offset_reset: earliest
    poll_timeout: 1s
    observability:
      metrics: true
config_center:
  enabled: true
  adapter: stellnula
  endpoint: http://localhost:8060
  namespace: platform
  group: app
  app_id: order-service
  client_id: order-service-local
  config_key: application.yaml
  labels:
    region: local
  sources:
    - config_key: application.yaml
      format: yaml
registry:
  enabled: true
  adapter: consul
  endpoints:
    - http://localhost:8500
  namespace: platform
  group: DEFAULT_GROUP
  cluster: DEFAULT
  service: order-service
  instance_id: order-service-1
  zone: zone-a
  ttl: 30s
  heartbeat_interval: 10s
  timeout: 3s
  labels:
    version: v1
  metadata:
    owner: platform
  observability:
    metrics: true
  service_endpoints:
    - name: http
      protocol: http
      host: 127.0.0.1
      port: 18080
      path: /
      weight: 100
discovery:
  enabled: true
  adapter: stellmap
  endpoints:
    - http://localhost:18090
  namespace: default
  refresh_interval: 10s
  stale_ttl: 1m
  observability:
    metrics: true
opentelemetry:
  log: true
  trace: true
  metrics: true
  trace_output: none
  metrics_output: prometheus
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AppName != "order-service" {
		t.Fatalf("unexpected app name %q", cfg.AppName)
	}
	if cfg.Environment != EnvUAT {
		t.Fatalf("unexpected env %q", cfg.Environment)
	}
	if cfg.HTTP.Server == nil || cfg.HTTP.Server.Addr != ":18080" {
		t.Fatalf("unexpected http server config %#v", cfg.HTTP.Server)
	}
	if cfg.HTTP.Server.Adapter != "chi" {
		t.Fatalf("unexpected http adapter %q", cfg.HTTP.Server.Adapter)
	}
	if cfg.HTTP.Server.Observability.Trace == nil || !*cfg.HTTP.Server.Observability.Trace {
		t.Fatalf("expected http server trace observability")
	}
	if cfg.HTTP.Client == nil || cfg.HTTP.Client.Timeout != "3s" {
		t.Fatalf("unexpected http client config %#v", cfg.HTTP.Client)
	}
	if cfg.HTTP.Client.MaxIdleConns != 100 {
		t.Fatalf("unexpected max idle conns %d", cfg.HTTP.Client.MaxIdleConns)
	}
	if cfg.HTTP.Client.Observability.Logs == nil || *cfg.HTTP.Client.Observability.Logs {
		t.Fatalf("expected http client logs observability disabled")
	}
	if cfg.HTTP.Client.Discovery == nil || cfg.HTTP.Client.Discovery.Adapter != "stellmap" {
		t.Fatalf("unexpected http client discovery config %#v", cfg.HTTP.Client.Discovery)
	}
	userClient := cfg.HTTP.Client.Clients["user-service"]
	if userClient.Discovery == nil || userClient.Discovery.Service != "user-service" || userClient.Discovery.EndpointName != "http" {
		t.Fatalf("unexpected user-service discovery config %#v", userClient.Discovery)
	}
	if userClient.Timeout != "2s" {
		t.Fatalf("unexpected user-service client config %#v", userClient)
	}
	if cfg.GRPC.Server == nil || cfg.GRPC.Server.Addr != ":19090" {
		t.Fatalf("unexpected grpc server config %#v", cfg.GRPC.Server)
	}
	if cfg.GRPC.Server.Adapter != "grpc-go" {
		t.Fatalf("unexpected grpc adapter %q", cfg.GRPC.Server.Adapter)
	}
	if cfg.GRPC.Server.Observability.Trace == nil || !*cfg.GRPC.Server.Observability.Trace {
		t.Fatalf("expected grpc server trace observability")
	}
	if cfg.GRPC.Client == nil || cfg.GRPC.Client.Timeout != "3s" {
		t.Fatalf("unexpected grpc client config %#v", cfg.GRPC.Client)
	}
	if cfg.GRPC.Client.Insecure == nil || !*cfg.GRPC.Client.Insecure {
		t.Fatalf("expected grpc client insecure")
	}
	if cfg.GRPC.Client.Observability.Logs == nil || *cfg.GRPC.Client.Observability.Logs {
		t.Fatalf("expected grpc client logs observability disabled")
	}
	if cfg.GRPC.Client.Discovery == nil || cfg.GRPC.Client.Discovery.Adapter != "stellmap" {
		t.Fatalf("unexpected grpc client discovery config %#v", cfg.GRPC.Client.Discovery)
	}
	grpcUserClient := cfg.GRPC.Client.Clients["user-service"]
	if grpcUserClient.Discovery == nil || grpcUserClient.Discovery.Service != "user-service" || grpcUserClient.Discovery.EndpointName != "grpc" {
		t.Fatalf("unexpected grpc user-service discovery config %#v", grpcUserClient.Discovery)
	}
	if grpcUserClient.Timeout != "2s" {
		t.Fatalf("unexpected grpc user-service client config %#v", grpcUserClient)
	}
	if cfg.Redis == nil || cfg.Redis.Addr != "localhost:6379" || cfg.Redis.DB != 1 {
		t.Fatalf("unexpected redis config %#v", cfg.Redis)
	}
	if cfg.Redis.DebugAPI == nil || cfg.Redis.DebugAPI.Enabled == nil || !*cfg.Redis.DebugAPI.Enabled || cfg.Redis.DebugAPI.Prefix != "/redis" {
		t.Fatalf("unexpected redis debug api config %#v", cfg.Redis.DebugAPI)
	}
	if cfg.Redis.Observability.Logs == nil || !*cfg.Redis.Observability.Logs {
		t.Fatalf("expected redis logs observability")
	}
	if cfg.MySQL == nil || cfg.MySQL.Driver != "mysql" || cfg.MySQL.MaxOpenConns != 25 {
		t.Fatalf("unexpected mysql config %#v", cfg.MySQL)
	}
	if cfg.MySQL.DebugAPI == nil || cfg.MySQL.DebugAPI.Enabled == nil || !*cfg.MySQL.DebugAPI.Enabled || cfg.MySQL.DebugAPI.Prefix != "/mysql" {
		t.Fatalf("unexpected mysql debug api config %#v", cfg.MySQL.DebugAPI)
	}
	if cfg.MySQL.Observability.Metrics == nil || !*cfg.MySQL.Observability.Metrics {
		t.Fatalf("expected mysql metrics observability")
	}
	if cfg.PostgreSQL == nil || cfg.PostgreSQL.Driver != "pgx" || cfg.PostgreSQL.MaxOpenConns != 25 {
		t.Fatalf("unexpected postgresql config %#v", cfg.PostgreSQL)
	}
	if cfg.PostgreSQL.DebugAPI == nil || cfg.PostgreSQL.DebugAPI.Enabled == nil || !*cfg.PostgreSQL.DebugAPI.Enabled || cfg.PostgreSQL.DebugAPI.Prefix != "/postgresql" {
		t.Fatalf("unexpected postgresql debug api config %#v", cfg.PostgreSQL.DebugAPI)
	}
	if cfg.PostgreSQL.Observability.Trace == nil || !*cfg.PostgreSQL.Observability.Trace {
		t.Fatalf("expected postgresql trace observability")
	}
	if cfg.Cache == nil || cfg.Cache.Adapter != "freecache" || cfg.Cache.TTL != "5m" || cfg.Cache.SizeBytes != 1048576 {
		t.Fatalf("unexpected cache config %#v", cfg.Cache)
	}
	if cfg.Cache.DebugAPI == nil || cfg.Cache.DebugAPI.Enabled == nil || !*cfg.Cache.DebugAPI.Enabled || cfg.Cache.DebugAPI.Prefix != "/cache" {
		t.Fatalf("unexpected cache debug api config %#v", cfg.Cache.DebugAPI)
	}
	if cfg.Cache.Observability.Metrics == nil || !*cfg.Cache.Observability.Metrics {
		t.Fatalf("expected cache metrics observability")
	}
	if cfg.Cache.Observability.Logs == nil || *cfg.Cache.Observability.Logs {
		t.Fatalf("expected cache logs observability disabled")
	}
	if cfg.MQ == nil || cfg.MQ.Adapter != "stellflow" || cfg.MQ.ClientID != "order-service" {
		t.Fatalf("unexpected mq config %#v", cfg.MQ)
	}
	if len(cfg.MQ.Brokers) != 1 || cfg.MQ.Brokers[0] != "stellflow://localhost:9092" {
		t.Fatalf("unexpected mq brokers %#v", cfg.MQ.Brokers)
	}
	if cfg.MQ.Producer.DefaultTopic != "orders.created" || cfg.MQ.Producer.AutoCreateTopicPartitionCount != 4 {
		t.Fatalf("unexpected mq producer config %#v", cfg.MQ.Producer)
	}
	if cfg.MQ.Consumer.GroupID != "order-worker" || cfg.MQ.Consumer.AutoOffsetReset != "earliest" || len(cfg.MQ.Consumer.Topics) != 1 {
		t.Fatalf("unexpected mq consumer config %#v", cfg.MQ.Consumer)
	}
	if cfg.MQ.Producer.Observability.Metrics == nil || !*cfg.MQ.Producer.Observability.Metrics {
		t.Fatalf("expected mq producer metrics observability")
	}
	if cfg.MQ.Consumer.Observability.Metrics == nil || !*cfg.MQ.Consumer.Observability.Metrics {
		t.Fatalf("expected mq consumer metrics observability")
	}
	if cfg.ConfigCenter == nil || cfg.ConfigCenter.Adapter != "stellnula" || cfg.ConfigCenter.Endpoint != "http://localhost:8060" {
		t.Fatalf("unexpected config center config %#v", cfg.ConfigCenter)
	}
	if len(cfg.ConfigCenter.Endpoints) != 1 || cfg.ConfigCenter.Endpoints[0] != "http://localhost:8060" {
		t.Fatalf("unexpected config center endpoints %#v", cfg.ConfigCenter.Endpoints)
	}
	if cfg.ConfigCenter.AppID != "order-service" || cfg.ConfigCenter.ClientID != "order-service-local" || cfg.ConfigCenter.Env != "uat" {
		t.Fatalf("unexpected config center identity %#v", cfg.ConfigCenter)
	}
	if cfg.ConfigCenter.Labels["region"] != "local" || len(cfg.ConfigCenter.Sources) != 1 {
		t.Fatalf("unexpected config center labels/sources %#v", cfg.ConfigCenter)
	}
	if cfg.Registry == nil || cfg.Registry.Adapter != "consul" || cfg.Registry.Namespace != "platform" {
		t.Fatalf("unexpected registry config %#v", cfg.Registry)
	}
	if len(cfg.Registry.Endpoints) != 1 || cfg.Registry.Endpoints[0] != "http://localhost:8500" {
		t.Fatalf("unexpected registry endpoints %#v", cfg.Registry.Endpoints)
	}
	if cfg.Registry.Service != "order-service" || cfg.Registry.InstanceID != "order-service-1" {
		t.Fatalf("unexpected registry instance config %#v", cfg.Registry)
	}
	if cfg.Registry.Labels["version"] != "v1" || cfg.Registry.Metadata["owner"] != "platform" {
		t.Fatalf("unexpected registry metadata labels=%#v metadata=%#v", cfg.Registry.Labels, cfg.Registry.Metadata)
	}
	if cfg.Registry.Observability.Metrics == nil || !*cfg.Registry.Observability.Metrics {
		t.Fatalf("expected registry metrics observability")
	}
	if len(cfg.Registry.ServiceEndpoints) != 1 || cfg.Registry.ServiceEndpoints[0].Port != 18080 {
		t.Fatalf("unexpected registry service endpoints %#v", cfg.Registry.ServiceEndpoints)
	}
	if cfg.Discovery == nil || cfg.Discovery.Adapter != "stellmap" || cfg.Discovery.RefreshInterval != "10s" {
		t.Fatalf("unexpected discovery config %#v", cfg.Discovery)
	}
	if cfg.Discovery.Observability.Metrics == nil || !*cfg.Discovery.Observability.Metrics {
		t.Fatalf("expected discovery metrics observability")
	}
	if cfg.Starter.OpenTelemetry == nil || !cfg.Starter.OpenTelemetry.Log.Enabled {
		t.Fatalf("expected opentelemetry log starter")
	}
}

func TestNormalizeGovernanceEnvInheritsConfigCenterWhenAppEnvOmitted(t *testing.T) {
	cfg := Config{
		AppName: "order-service",
		ConfigCenter: &ConfigCenterConfig{
			Endpoint: "http://localhost:8060",
			Env:      "uat",
		},
		Governance: &GovernanceConfig{
			ConfigCenter: GovernanceConfigCenterConfig{
				Endpoint: "http://localhost:8060",
			},
		},
	}.Normalize()

	if cfg.Governance == nil || cfg.Governance.Env != "uat" {
		t.Fatalf("expected governance env to inherit config center env, got %#v", cfg.Governance)
	}
}

func TestNormalizeGovernanceEnvPrefersExplicitAppEnv(t *testing.T) {
	cfg := Config{
		AppName:     "order-service",
		Environment: EnvProd,
		ConfigCenter: &ConfigCenterConfig{
			Endpoint: "http://localhost:8060",
			Env:      "uat",
		},
		Governance: &GovernanceConfig{
			ConfigCenter: GovernanceConfigCenterConfig{
				Endpoint: "http://localhost:8060",
			},
		},
	}.Normalize()

	if cfg.Governance == nil || cfg.Governance.Env != string(EnvProd) {
		t.Fatalf("expected governance env to prefer app env, got %#v", cfg.Governance)
	}
}

func TestLoadBytesAndMergeConfig(t *testing.T) {
	enabled := true
	base := Config{
		AppName:     "local-service",
		Environment: EnvDev,
		HTTP: HTTPConfig{
			Server: &HTTPServerConfig{
				Enabled: &enabled,
				Port:    8080,
			},
		},
		ConfigCenter: &ConfigCenterConfig{
			Adapter:  "stellnula",
			Endpoint: "http://localhost:8060",
		},
	}.Normalize()

	remote, err := LoadBytes("remote-config.yaml", []byte(`
http:
  server:
    port: 18080
    adapter: chi
mq:
  adapter: kafka
  brokers:
    - localhost:9092
  producer:
    default_topic: orders.created
`))
	if err != nil {
		t.Fatalf("load remote config bytes: %v", err)
	}

	merged := Merge(base, remote)
	if merged.AppName != "local-service" {
		t.Fatalf("expected local app name to be preserved, got %q", merged.AppName)
	}
	if merged.HTTP.Server == nil || merged.HTTP.Server.Enabled == nil || !*merged.HTTP.Server.Enabled {
		t.Fatalf("expected local http enabled flag to be preserved, got %#v", merged.HTTP.Server)
	}
	if merged.HTTP.Server.Addr != ":18080" || merged.HTTP.Server.Adapter != "chi" {
		t.Fatalf("unexpected merged http server %#v", merged.HTTP.Server)
	}
	if merged.MQ == nil || merged.MQ.Adapter != "kafka" || merged.MQ.Producer.DefaultTopic != "orders.created" {
		t.Fatalf("unexpected merged mq config %#v", merged.MQ)
	}
	if merged.ConfigCenter == nil || merged.ConfigCenter.Adapter != "stellnula" {
		t.Fatalf("expected local config center bootstrap to be preserved, got %#v", merged.ConfigCenter)
	}
}

func TestLoadApplicationYAMLConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "application.yaml")
	content := `
app:
  name: user-service
  env: dev
http:
  server:
    port: 8081
  client:
    timeout: 3s
    clients:
      user-service:
        base_url: http://localhost:8081
        timeout: 2s
opentelemetry:
  log:
    enabled: false
    output: file
    dir: logs
    file_name: app.log
    max_size_bytes: 104857600
    max_backups: 5
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AppName != "user-service" {
		t.Fatalf("unexpected app name %q", cfg.AppName)
	}
	if cfg.HTTP.Server == nil || cfg.HTTP.Server.Addr != ":8081" {
		t.Fatalf("unexpected http server config %#v", cfg.HTTP.Server)
	}
	if cfg.HTTP.Client == nil || cfg.HTTP.Client.Clients["user-service"].BaseURL != "http://localhost:8081" {
		t.Fatalf("unexpected http client config %#v", cfg.HTTP.Client)
	}
	if cfg.Starter.OpenTelemetry == nil || cfg.Starter.OpenTelemetry.Log.Output != "file" {
		t.Fatalf("unexpected opentelemetry config: %#v", cfg.Starter.OpenTelemetry)
	}
	if cfg.Starter.OpenTelemetry.Log.FileName != "app.log" {
		t.Fatalf("unexpected log file name %q", cfg.Starter.OpenTelemetry.Log.FileName)
	}
}

func TestLoadGovernanceConfigDefaultsAndInheritance(t *testing.T) {
	path := filepath.Join(t.TempDir(), "application.yaml")
	content := `
app:
  name: order-service
  env: uat
  zone: zone-a
config_center:
  enabled: true
  adapter: stellnula
  endpoint: http://localhost:8060
  app_id: order-service
  client_id: order-service-local
  env: uat
  labels:
    region: local
governance:
  enabled: true
  route:
    enabled: true
  rate_limit:
    enabled: true
    distributed:
      enabled: true
      address: 127.0.0.1:19091
  auth:
    enabled: true
    key_provider:
      placeholder: true
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Governance == nil {
		t.Fatalf("expected governance config")
	}
	if cfg.Governance.Adapter != "stellorbit" {
		t.Fatalf("unexpected governance adapter %q", cfg.Governance.Adapter)
	}
	if cfg.Governance.AppID != "order-service" || cfg.Governance.ClientID != "order-service-local" {
		t.Fatalf("unexpected governance identity %#v", cfg.Governance)
	}
	if cfg.Governance.ConfigCenter.Adapter != "stellnula" || cfg.Governance.ConfigCenter.Endpoint != "http://localhost:8060" {
		t.Fatalf("unexpected governance config center %#v", cfg.Governance.ConfigCenter)
	}
	if cfg.Governance.ConfigCenter.Namespace != "governance" || cfg.Governance.ConfigCenter.Group != "service-governance" {
		t.Fatalf("unexpected governance rule scope %#v", cfg.Governance.ConfigCenter)
	}
	if cfg.Governance.RateLimit.Mode != "local" || cfg.Governance.RateLimit.Behavior != "reject" {
		t.Fatalf("unexpected rate limit defaults %#v", cfg.Governance.RateLimit)
	}
	if cfg.Governance.RateLimit.Local.DefaultRate != 100 || cfg.Governance.RateLimit.Local.DefaultBurst != 100 {
		t.Fatalf("unexpected local rate limit defaults %#v", cfg.Governance.RateLimit.Local)
	}
	if cfg.Governance.RateLimit.Distributed.Adapter != "stellpulsar" || cfg.Governance.RateLimit.Distributed.Fallback != "fail_open" {
		t.Fatalf("unexpected distributed rate limit defaults %#v", cfg.Governance.RateLimit.Distributed)
	}
	if cfg.Governance.Auth.KeyProvider.Adapter != "stellguard-agent" || !cfg.Governance.Auth.KeyProvider.Placeholder {
		t.Fatalf("unexpected auth key provider %#v", cfg.Governance.Auth.KeyProvider)
	}
}

func TestLoadUsesCommandLineConfigPath(t *testing.T) {
	path := writeTestConfig(t, filepath.Join(t.TempDir(), "application.yaml"), "cli-service")
	setArgs(t, "stellar", "--config", path)
	clearConfigEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AppName != "cli-service" {
		t.Fatalf("unexpected app name %q", cfg.AppName)
	}
}

func TestLoadUsesCommandLineConfigDirectory(t *testing.T) {
	dir := t.TempDir()
	writeTestConfig(t, filepath.Join(dir, "application.yml"), "cli-dir-service")
	setArgs(t, "stellar", "--config.file="+dir)
	clearConfigEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AppName != "cli-dir-service" {
		t.Fatalf("unexpected app name %q", cfg.AppName)
	}
}

func TestLoadUsesEnvConfigPath(t *testing.T) {
	path := writeTestConfig(t, filepath.Join(t.TempDir(), "application.yaml"), "env-service")
	setArgs(t, "stellar")
	clearConfigEnv(t)
	t.Setenv(EnvConfigFile, path)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AppName != "env-service" {
		t.Fatalf("unexpected app name %q", cfg.AppName)
	}
}

func TestLoadCommandLineConfigPathOverridesEnv(t *testing.T) {
	envPath := writeTestConfig(t, filepath.Join(t.TempDir(), "application.yaml"), "env-service")
	cliPath := writeTestConfig(t, filepath.Join(t.TempDir(), "application.yaml"), "cli-service")
	setArgs(t, "stellar", "--stellar.config", cliPath)
	clearConfigEnv(t)
	t.Setenv(EnvConfigFile, envPath)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.AppName != "cli-service" {
		t.Fatalf("unexpected app name %q", cfg.AppName)
	}
}

func TestLoadRejectsMissingCommandLineConfigValue(t *testing.T) {
	setArgs(t, "stellar", "--config")
	clearConfigEnv(t)

	if _, err := Load(); err == nil {
		t.Fatalf("expected missing command line config value to be rejected")
	}
}

func TestLoadFileRejectsUnsupportedConfigName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stellar.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: order-service\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := LoadFile(path); err == nil {
		t.Fatalf("expected unsupported config file name to be rejected")
	}
}

func TestConfigPathsPreferMainDirBeforeWorkingDir(t *testing.T) {
	mainDir := filepath.Join("repo", "examples", "http", "server", "simple")
	workingDir := filepath.Join("repo")

	paths := configPathsInDirs(mainDir, workingDir)

	if paths[0] != filepath.Join(mainDir, "application.yml") {
		t.Fatalf("expected main dir application.yml first, got %q", paths[0])
	}
	if paths[1] != filepath.Join(mainDir, "application.yaml") {
		t.Fatalf("expected main dir application.yaml second, got %q", paths[1])
	}
	if paths[2] != filepath.Join(workingDir, "application.yml") {
		t.Fatalf("expected working dir after main dir candidates, got %q", paths[2])
	}
}

func writeTestConfig(t *testing.T, path string, appName string) string {
	t.Helper()
	content := "app:\n  name: " + appName + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func setArgs(t *testing.T, args ...string) {
	t.Helper()
	originalArgs := os.Args
	os.Args = args
	t.Cleanup(func() {
		os.Args = originalArgs
	})
}

func clearConfigEnv(t *testing.T) {
	t.Helper()
	t.Setenv(EnvConfigFile, "")
	t.Setenv(EnvConfig, "")
	t.Setenv(EnvApplicationConfig, "")
}
