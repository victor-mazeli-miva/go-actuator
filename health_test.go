package actuator

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubCheck struct {
	name string
	err  error
}

func (s stubCheck) Name() string { return s.name }

func (s stubCheck) Check(_ context.Context) error { return s.err }

func TestEvaluateHealth(t *testing.T) {
	act := New()
	act.RegisterHealthCheck(stubCheck{name: "postgres"})

	resp := act.EvaluateHealth(context.Background())
	if !resp.IsHealthy() {
		t.Fatalf("status = %q, want UP", resp.Status)
	}
	if resp.Checks["postgres"] != statusUP {
		t.Fatalf("postgres = %q, want UP", resp.Checks["postgres"])
	}
}

func TestHealthNoChecks(t *testing.T) {
	act := New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	act.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != statusUP {
		t.Fatalf("status = %q, want %q", resp.Status, statusUP)
	}
	if resp.Checks != nil {
		t.Fatalf("checks = %v, want nil", resp.Checks)
	}
}

func TestHealthAllPass(t *testing.T) {
	act := New()
	act.RegisterHealthCheck(stubCheck{name: "postgres"})
	act.RegisterHealthCheck(stubCheck{name: "redis"})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	act.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != statusUP {
		t.Fatalf("status = %q, want %q", resp.Status, statusUP)
	}
	if resp.Checks["postgres"] != statusUP || resp.Checks["redis"] != statusUP {
		t.Fatalf("checks = %v, want all UP", resp.Checks)
	}
}

func TestHealthOneFail(t *testing.T) {
	act := New()
	act.RegisterHealthCheck(stubCheck{name: "postgres", err: errors.New("connection refused")})
	act.RegisterHealthCheck(stubCheck{name: "redis"})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	act.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != statusDOWN {
		t.Fatalf("status = %q, want %q", resp.Status, statusDOWN)
	}
	if resp.Checks["postgres"] != statusDOWN {
		t.Fatalf("postgres = %q, want DOWN", resp.Checks["postgres"])
	}
}

func TestReadyMatchesHealth(t *testing.T) {
	act := New()
	act.RegisterHealthCheck(stubCheck{name: "postgres", err: errors.New("down")})

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()
	act.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestHealthContextCancellation(t *testing.T) {
	act := New()
	act.RegisterHealthCheck(stubCheck{
		name: "slow",
		err:  context.Canceled,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := httptest.NewRequest(http.MethodGet, "/health", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	act.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}
