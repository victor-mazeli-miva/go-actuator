package actuator

import (
	"context"
	"net/http"
)

type HealthCheck interface {
	Name() string
	Check(ctx context.Context) error
}

// EvaluateHealth runs all registered health checks and returns the aggregated result.
func (a *Actuator) EvaluateHealth(ctx context.Context) HealthResponse {
	return a.aggregateChecks(ctx)
}

func (a *Actuator) aggregateChecks(ctx context.Context) HealthResponse {
	if len(a.checks) == 0 {
		return HealthResponse{Status: statusUP}
	}

	checks := make(map[string]string, len(a.checks))
	overall := statusUP

	for _, h := range a.checks {
		name := h.Name()
		if err := h.Check(ctx); err != nil {
			checks[name] = statusDOWN
			overall = statusDOWN
		} else {
			checks[name] = statusUP
		}
	}

	return HealthResponse{
		Status: overall,
		Checks: checks,
	}
}

func (a *Actuator) handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := a.aggregateChecks(r.Context())
	statusCode := http.StatusOK
	if resp.Status == statusDOWN {
		statusCode = http.StatusServiceUnavailable
	}
	writeJSON(w, statusCode, resp)
}

func (a *Actuator) handleReady(w http.ResponseWriter, r *http.Request) {
	a.handleHealth(w, r)
}

func (a *Actuator) handleLive(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: statusUP})
}
