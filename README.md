# go-actuator


A lightweight, idiomatic actuator package for Go services. It exposes operational HTTP endpoints for health checks, Kubernetes probes, and basic runtime statistics—without turning into a full observability platform.

## Goals

- **Simple** — small API surface, easy to read and maintain
- **Stdlib-first** — core package depends only on `net/http`
- **Framework-agnostic** — Chi, Gin, and gRPC integrations live under `adapters/`
- **Operational-first** — health and runtime before metrics/tracing
- **Incremental** — add complexity only when real usage needs it

This is intentionally **not** an enterprise platform. It focuses on what most services need on day one.

## Endpoints

| Route      | Purpose                                      |
|------------|----------------------------------------------|
| `GET /health`  | Aggregated health checks (200 UP, 503 DOWN) |
| `GET /ready`   | Kubernetes readiness (runs registered checks) |
| `GET /live`    | Kubernetes liveness (process is running)      |
| `GET /runtime` | Go runtime stats (goroutines, heap, GC, CPUs)   |

### gRPC health (server)

Register the standard [`grpc.health.v1.Health`](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) service on your gRPC server (for Kubernetes gRPC probes and load balancers):

```go
srv := grpc.NewServer()
grpcadapter.Register(srv, act)
```

| RPC | `service` field | Behavior |
|-----|-----------------|----------|
| `Check` | `""` | Overall health (all registered checks) |
| `Check` | `"postgres"` | Single check matching `HealthCheck.Name()` |
| `List` | — | Status for `""` and each registered check |
| `Watch` | — | Not supported |

Import: `github.com/victor-mazeli-miva/go-actuator/adapters/grpc` (package `grpcadapter`).

### Health response

```json
{
  "status": "UP",
  "checks": {
    "postgres": "UP",
    "redis": "UP"
  }
}