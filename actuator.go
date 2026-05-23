// Package actuator provides lightweight operational HTTP endpoints for Go services,
// including health checks, Kubernetes liveness/readiness probes, and runtime statistics.
//
// Use New to create an Actuator, register HealthCheck implementations, and expose
// Router as an http.Handler. Framework integrations live in the adapters subpackage.
package actuator

import "net/http"

type Actuator struct {
	mux    *http.ServeMux
	checks []HealthCheck
}

func New() *Actuator {
	a := &Actuator{
		mux: http.NewServeMux(),
	}
	a.registerRoutes()
	return a
}

// RegisterHealthCheck adds a health check to the actuator.
// Checks run sequentially on each /health and /ready request.
// RegisterHealthCheck is not safe for concurrent use.
func (a *Actuator) RegisterHealthCheck(h HealthCheck) {
	a.checks = append(a.checks, h)
}

func (a *Actuator) Router() http.Handler {
	return recoverHandler(a.mux)
}
