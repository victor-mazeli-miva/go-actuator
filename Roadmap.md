# Go Actuator — Road map Plan

## Vision

Build a lightweight, idiomatic, extensible actuator package for Go services that provides operational endpoints such as health checks, runtime information, and readiness/liveness probes.

The project should prioritize:

- Simplicity
- Idiomatic Go
- Minimal dependencies
- Excellent developer experience
- Framework agnostic design
- Incremental evolution

This is NOT intended to be an enterprise platform initially.

---

# Initial Scope

## Supported Frameworks

Phase 1 support:

- net/http
- Chi
- Gin

---

# Design Philosophy

## Core Principles

| Principle | Description |
|---|---|
| Lightweight | Small API surface |
| Stdlib First | Prefer standard library |
| Composition | Avoid inheritance-style abstractions |
| Minimal Interfaces | Small focused interfaces |
| Framework Agnostic | Core package independent from frameworks |
| Incremental Growth | Add complexity only when needed |
| Operational First | Health/runtime before observability |

---

# Phase 1 Goals

The first milestone should support:

| Endpoint | Purpose |
|---|---|
| `/health` | Aggregated health checks |
| `/live` | Kubernetes liveness |
| `/ready` | Kubernetes readiness |
| `/runtime` | Go runtime statistics |

---

# Initial Project Structure

```text
actuator/
├── actuator.go
├── health.go
├── runtime.go
├── response.go
├── routes.go
├── adapters/
│   ├── chi.go
│   ├── gin.go
│   └── nethttp.go
├── examples/
│   ├── chi/
│   ├── gin/
│   └── nethttp/
└── go.mod
```

---

# Core Architecture

## Primary Type

```go
type Actuator struct {
    mux *http.ServeMux

    checks []HealthCheck
}
```

The actuator acts as:

- route registry
- health registry
- operational runtime container

---

# Health System

## Health Check Interface

```go
type HealthCheck interface {
    Name() string
    Check(ctx context.Context) error
}
```

Keep the interface intentionally small.

Do NOT add:
- metadata
- priorities
- generics
- async pipelines
- plugin hooks

in the first implementation.

---

# Health Aggregation Flow

```text
HTTP Request
    ↓
Execute health checks
    ↓
Collect errors
    ↓
Aggregate status
    ↓
Return JSON response
```

---

# Health Response

## Successful Response

```json
{
  "status": "UP",
  "checks": {
    "postgres": "UP",
    "redis": "UP"
  }
}
```

## Failure Response

```json
{
  "status": "DOWN",
  "checks": {
    "postgres": "DOWN"
  }
}
```

---

# Runtime Endpoint

## Responsibilities

Expose lightweight runtime information using the Go runtime package.

## Initial Runtime Metrics

| Metric | Source |
|---|---|
| Goroutines | runtime.NumGoroutine |
| Heap Alloc | runtime.MemStats |
| GC Count | runtime.MemStats |
| CPU Count | runtime.NumCPU |

---

# Routing Strategy

The core package should rely ONLY on:

```go
net/http
```

Framework adapters should only map framework routers to the actuator handler.

---

# net/http Integration

## Example

```go
act := actuator.New()

http.ListenAndServe(":8080", act.Router())
```

---

# Chi Integration

## Example

```go
r := chi.NewRouter()

act := actuator.New()

r.Mount("/actuator", act.Router())
```

---

# Gin Integration

## Example

```go
r := gin.Default()

act := actuator.New()

r.Any("/actuator/*path", gin.WrapH(act.Router()))
```

---

# Initial API Design

## Constructor

```go
func New() *Actuator
```

---

## Register Health Check

```go
func (a *Actuator) RegisterHealthCheck(h HealthCheck)
```

---

## Router Exposure

```go
func (a *Actuator) Router() http.Handler
```

---

# Internal Route Design

## Routes

| Route | Method |
|---|---|
| `/health` | GET |
| `/live` | GET |
| `/ready` | GET |
| `/runtime` | GET |

---

# Error Handling Rules

## Principles

- Never panic inside handlers
- Always return JSON
- Avoid leaking internal errors
- Use proper HTTP status codes

---

# JSON Response Rules

All responses should:

- use `application/json`
- include consistent structure
- avoid unnecessary nesting

---

# Concurrency Rules

## Health Checks

Initial implementation may execute checks sequentially.

Only introduce concurrency later if needed.

Avoid premature optimization.

---

# Dependency Rules

## Initial Dependencies

Allowed:

- stdlib
- chi
- gin

Avoid:
- config frameworks
- DI frameworks
- observability frameworks
- reflection-heavy packages

---

# What NOT To Build Yet

Avoid implementing:

- Prometheus
- OpenTelemetry
- plugin systems
- module registries
- distributed aggregation
- dependency injection containers
- middleware pipelines
- advanced security layers
- async schedulers

These belong in later stages.

---

# Phase 2 (Future)

Only after Phase 1 stabilizes.

## Planned Features

- Prometheus metrics
- Middleware instrumentation
- pprof support
- Request timing
- Request counters
- Framework middleware helpers

---

# Phase 3 (Future)

## Potential Features

- OpenTelemetry
- Plugin architecture
- Dependency probes
- Security policies
- Distributed health aggregation

---

# Testing Strategy

## Unit Tests

- health aggregation
- response serialization
- runtime collection

## Integration Tests

- chi mounting
- gin integration
- net/http integration

---

# Success Criteria

Phase 1 is successful when:

- health checks register correctly
- endpoints return valid JSON
- chi integration works
- gin integration works
- net/http integration works
- runtime metrics are exposed
- the package remains lightweight and readable

---

# Initial Milestone

The following should work cleanly:

```go
db := sql.Open(...)

act := actuator.New()

act.RegisterHealthCheck(
    postgres.New(db),
)

r := chi.NewRouter()

r.Mount("/actuator", act.Router())

http.ListenAndServe(":8080", r)
```

Then:

```bash
curl localhost:8080/actuator/health
```

Returns:

```json
{
  "status": "UP"
}
```

---

# Long-Term Goal

The actuator should evolve organically into:

- a reusable operational layer
- a lightweight observability runtime
- a cloud-native service utility

without sacrificing simplicity or idiomatic Go design.
