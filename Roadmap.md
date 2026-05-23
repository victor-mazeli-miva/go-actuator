# Go Actuator — Roadmap

## Vision

**go-actuator** is a lightweight, idiomatic operational layer for Go services: health aggregation, Kubernetes probes, and basic runtime insight—without growing into a full observability platform.

We want teams to wire operational endpoints in minutes, keep the core on `net/http`, and add depth only when real production pain appears. The north star is a small, readable package that feels like standard library ergonomics, not an enterprise framework.

### Guiding principles

| Principle | What we optimize for |
|-----------|----------------------|
| **Lightweight** | Small API surface; easy to read in one sitting |
| **Stdlib-first** | Core depends only on `net/http` |
| **Framework-agnostic** | Chi, Gin, gRPC, and others live in thin `adapters/` |
| **Operational-first** | Health and runtime before metrics and tracing |
| **Incremental** | Complexity only when usage justifies it |
| **Idiomatic Go** | Explicit code, small interfaces, composition over hierarchy |

### Non-goals (for now)

We are **not** building an enterprise platform, plugin ecosystem, or distributed health mesh in the early phases. Prometheus, OpenTelemetry, DI containers, and async orchestration stay out of scope until Phase 1 is stable and demand is clear.

---

## Phase 1 — Operational foundation

**Goal:** Give every Go service a dependable day-one actuator: register checks, expose standard routes, mount on common routers.

### Delivered

- **HTTP endpoints:** `/health`, `/ready`, `/live`, `/runtime`
- **Health model:** small `HealthCheck` interface; aggregated JSON (`UP` / `DOWN`)
- **Runtime snapshot:** goroutines, heap, GC count, CPU count
- **Adapters:** `net/http`, Chi, Gin
- **gRPC:** standard `grpc.health.v1` registration for probes and load balancers
- **Examples** for each integration path

### Stabilization (in progress)

- Harden handler behavior (no panics, consistent JSON, correct status codes)
- Expand integration tests across adapters
- Polish README and examples so onboarding matches the vision above

**Done when:** A service can register dependencies, mount the actuator, and get trustworthy health/readiness responses in Kubernetes and local dev—with no surprise dependencies or magic.

---

## Phase 2 — Deeper operations (planned)

**Goal:** Optional, opt-in tooling for teams that outgrow “health + runtime” but still want a thin library.

| Direction | Intent |
|-----------|--------|
| **Prometheus metrics** | Expose common process/request counters without owning a metrics stack |
| **pprof** | Safe, documented hooks for on-demand profiling |
| **Middleware helpers** | Thin instrumentation wrappers for Chi/Gin/`net/http` |
| **Request timing / counters** | Lightweight HTTP observability at the edge |

Phase 2 features ship as separate, optional surfaces—never as required core complexity.

---

## Phase 3 — Advanced observability (exploratory)

**Goal:** Explore integrations only if Phase 1–2 stay simple and adoption warrants it.

- OpenTelemetry bridges (traces/metrics) as adapters, not core
- Richer dependency probes (caches, queues, external APIs)
- Security-oriented actuator policies (auth on admin routes, IP allowlists)
- Distributed or federated health views (only if a concrete use case emerges)

Each item needs a clear “why now” before design work starts.

---

## Long-term direction

Over time, **go-actuator** may grow into a reusable operational utility for cloud-native Go services: the place you mount health, readiness, runtime, and—when you choose—metrics and profiling.

That growth must stay **organic**: same readability, same stdlib-first core, same refusal to become a framework people fight against.

---

## How to influence the roadmap

Open an issue or discussion with:

- What you run in production (K8s, mesh, gRPC-only, etc.)
- What hurt without this library
- What you explicitly do **not** want us to build

We add phases from real usage, not hypothetical enterprise requirements.
