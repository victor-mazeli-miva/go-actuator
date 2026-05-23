# go-actuator


A lightweight, idiomatic actuator package for Go services. It exposes operational HTTP endpoints for health checks, Kubernetes probes, and basic runtime statistics—without turning into a full observability platform.

## Goals

- **Simple** — small API surface, easy to read and maintain
- **Stdlib-first** — core package depends only on `net/http`
- **Framework-agnostic** — Chi and Gin integrations live in `adapters/`
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

### Health response

```json
{
  "status": "UP",
  "checks": {
    "postgres": "UP",
    "redis": "UP"
  }
}