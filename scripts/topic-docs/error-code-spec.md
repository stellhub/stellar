## Abstract

A usable enterprise error-code model should not fight HTTP, gRPC, or observability standards. It should work with them. The right design is a layered one:

1. transport status
2. canonical base error category
3. business-specific reason
4. structured error details
5. observability mapping

The key principle is to keep the base error space small and stable while allowing business variation in structured extensions rather than inventing endless protocol-level private codes.

## 1. Why Error Design Usually Goes Wrong

Most broken API error systems suffer from the same problems:

- transport status and business meaning are mixed together
- teams create large private integer code spaces with weak semantics
- responses are optimized for humans reading strings instead of machines making decisions

That makes client behavior unstable and observability aggregation noisy.

## 2. Standards Baseline

The recommended baseline is:

- RFC 9110 for HTTP semantics
- RFC 9457 Problem Details for HTTP error payloads
- canonical gRPC status codes
- `google.rpc.Status` and related detail models
- OpenTelemetry error recording conventions

The point is not to copy every standard mechanically. The point is to avoid inventing a parallel universe for error semantics.

## 3. The Layered Error Model

### 3.1 Transport Status

Transport status expresses protocol semantics:

- HTTP status code for HTTP APIs
- gRPC status code for gRPC APIs

These should not be redefined to encode business-specific meaning.

### 3.2 Canonical Base Error Category

This is the cross-protocol stable layer. It answers:

```text
What broad class of failure is this?
```

Recommended categories align with canonical gRPC / `google.rpc.Code` style naming, such as:

- `INVALID_ARGUMENT`
- `UNAUTHENTICATED`
- `PERMISSION_DENIED`
- `NOT_FOUND`
- `ALREADY_EXISTS`
- `FAILED_PRECONDITION`
- `RESOURCE_EXHAUSTED`
- `DEADLINE_EXCEEDED`
- `UNAVAILABLE`
- `INTERNAL`

### 3.3 Business Reason

This is the extension point. It answers:

```text
Why did this happen in business terms?
```

Examples:

- `ORDER_NOT_CANCELLABLE`
- `INSUFFICIENT_INVENTORY`
- `POLICY_VERSION_MISMATCH`

Business variation belongs here, not in custom transport status codes.

### 3.4 Structured Details

This layer holds machine-usable context such as:

- field validation violations
- precondition failures
- retry hints
- resource identifiers
- missing dependencies

### 3.5 Observability Mapping

Logs, traces, and metrics need low-cardinality error classification. They should aggregate on stable categories and selected reason fields, not raw human-readable messages.

## 4. Canonical Base Code Set

The most useful enterprise rule is to keep the base code set intentionally small.

| Base Code | Typical Meaning | Common HTTP Mapping |
| --- | --- | --- |
| `INVALID_ARGUMENT` | malformed or invalid input | `400` |
| `UNAUTHENTICATED` | authentication missing or invalid | `401` |
| `PERMISSION_DENIED` | authenticated but not allowed | `403` |
| `NOT_FOUND` | target resource missing | `404` |
| `ALREADY_EXISTS` | creation target already exists | `409` |
| `FAILED_PRECONDITION` | current state blocks the operation | `409`, `412`, or `428` |
| `RESOURCE_EXHAUSTED` | quota or limit exceeded | `429` |
| `DEADLINE_EXCEEDED` | operation timed out | `504` |
| `UNAVAILABLE` | service is temporarily unavailable | `503` |
| `INTERNAL` | internal failure | `500` |

The selection rule is simple: return the most specific stable category available, not the vaguest one.

## 5. HTTP Error Model

### 5.1 Status-Line Rules

HTTP APIs should:

- return an RFC-consistent HTTP status code
- never hide failure behind `200 OK`
- choose the most specific 4xx or 5xx status available

### 5.2 Response Body

The recommended payload model is RFC 9457 `application/problem+json` plus a few stable enterprise extensions:

- `code` for the canonical base category
- `domain` for business domain
- `reason` for stable business reason
- `retryable` for retry guidance
- `request_id` for request correlation
- `violations` for field or precondition details

Example:

```json
{
  "type": "https://errors.stellhub.top/trade/order-not-cancellable",
  "title": "Order cannot be cancelled",
  "status": 409,
  "detail": "The order has already entered settlement and can no longer be cancelled.",
  "code": "FAILED_PRECONDITION",
  "domain": "trade.order",
  "reason": "ORDER_NOT_CANCELLABLE",
  "retryable": false,
  "request_id": "req-20260426-000001"
}
```

## 6. gRPC Error Model

gRPC APIs should:

- return canonical gRPC status codes
- avoid inventing new protocol-level status categories
- use `google.rpc.Status` and related detail messages where structured context is needed

This keeps transport semantics standard while still allowing rich machine-readable detail.

## 7. Business Extension Rules

Business-specific reasons should be:

- stable
- machine-readable
- low enough cardinality to remain useful
- scoped by business domain where needed

The purpose of a business reason is not to mirror every possible message string. It is to provide a stable branch point for client logic, platform governance, and analytics.

## 8. Observability Mapping

Error models should align with traces, metrics, and logs.

Recommended practice:

- transport status is recorded in the protocol-specific status field
- canonical base category is recorded as a low-cardinality error classification
- selected business reason is included where governance or analysis requires it
- free-form messages remain diagnostic, not primary aggregation keys

This allows platforms to answer questions like:

- which failure classes dominate a service?
- which business reasons are trending?
- which APIs produce retryable versus non-retryable failures?

## 9. Implementation Guidance

### 9.1 For HTTP APIs

- keep transport semantics honest
- use Problem Details consistently
- keep the base code set small

### 9.2 For gRPC APIs

- use canonical status codes correctly
- add structured detail rather than new error categories

### 9.3 For SDKs and Clients

- branch on stable base codes and selected reasons
- avoid parsing human-readable strings for control logic

## 10. Conclusion

The strongest error systems are layered, not overloaded. Transport status should describe protocol truth. A small canonical base code space should describe failure class. Business reasons should describe domain-specific causes. Structured details should carry repair context. Observability mapping should keep the whole model aggregatable.

That approach is more sustainable than inventing a giant private code universe, and it works across HTTP, gRPC, client SDKs, and observability pipelines with far less long-term friction.
