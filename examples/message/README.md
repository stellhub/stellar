# Message Quickstarts

This directory contains two message queue quickstarts with the same HTTP API:

- `stellflow/quickstart`: uses `mq.adapter: stellflow`.
- `kafka/quickstart`: uses `mq.adapter: kafka` and requires CGO/librdkafka for Confluent Kafka Go SDK.

## Stellflow

Start a Stellflow broker on `localhost:9092`, then run:

```bash
go run ./examples/message/stellflow/quickstart
```

Send a message:

```bash
curl -X POST http://localhost:18085/api/v1/messages/send \
  -H "Content-Type: application/json" \
  -d '{"key":"order-1001","value":"hello stellflow","headers":{"source":"quickstart"}}'
```

Receive and commit messages:

```bash
curl -X POST http://localhost:18085/api/v1/messages/receive \
  -H "Content-Type: application/json" \
  -d '{"commit":true}'
```

You can also receive from a custom topic by adding `topic` or `topics` to the request body.

## Kafka

Start a Kafka broker on `localhost:9092`, make sure the local toolchain can build CGO code, then run:

```bash
CGO_ENABLED=1 go run ./examples/message/kafka/quickstart
```

Send a message:

```bash
curl -X POST http://localhost:18086/api/v1/messages/send \
  -H "Content-Type: application/json" \
  -d '{"key":"order-1001","value":"hello kafka","headers":{"source":"quickstart"}}'
```

Receive and commit messages:

```bash
curl -X POST http://localhost:18086/api/v1/messages/receive \
  -H "Content-Type: application/json" \
  -d '{"commit":true}'
```

You can also receive from a custom topic by adding `topic` or `topics` to the request body.

Both examples also expose:

```text
GET http://localhost:18085/api/v1/messages/status
GET http://localhost:18086/api/v1/messages/status
GET http://localhost:18085/metrics
GET http://localhost:18086/metrics
```
