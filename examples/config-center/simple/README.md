# Config Center Simple Example

This example starts Stellar with `config_center.adapter: stellnula` and exposes a small HTTP API for checking the merged runtime config.

Start StellNula on `localhost:8060`, then create a config item whose key or id is `application.yaml` and whose content is a normal Stellar config file:

```yaml
http:
  server:
    enabled: true
    port: 18087
mq:
  enabled: false
```

Run:

```bash
go run ./examples/config-center/simple
```

Try the API:

```text
GET http://localhost:18087/api/v1/config-center/status
GET http://localhost:18087/api/v1/config-center/config
GET http://localhost:18087/api/v1/config-center/sources
GET http://localhost:18087/metrics
```

To try Nacos, change `application.yml` to:

```yaml
config_center:
  enabled: true
  adapter: nacos
  endpoint: http://localhost:8848
  namespace: public
  group: DEFAULT_GROUP
  data_id: application.yaml
```

