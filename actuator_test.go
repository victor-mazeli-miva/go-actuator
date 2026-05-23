package actuator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterEndpoints(t *testing.T) {
	act := New()
	handler := act.Router()

	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/health", http.StatusOK, statusUP},
		{"/ready", http.StatusOK, statusUP},
		{"/live", http.StatusOK, statusUP},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", ct)
			}

			var resp HealthResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if resp.Status != tt.wantBody {
				t.Fatalf("status = %q, want %q", resp.Status, tt.wantBody)
			}
		})
	}
}

func TestRuntimeEndpoint(t *testing.T) {
	act := New()
	req := httptest.NewRequest(http.MethodGet, "/runtime", nil)
	rec := httptest.NewRecorder()

	act.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp RuntimeResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Goroutines < 1 {
		t.Fatalf("goroutines = %d, want >= 1", resp.Goroutines)
	}
}

func TestRecoverHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, _ *http.Request) {
		panic("boom")
	})

	handler := recoverHandler(mux)
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	var resp errorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != statusDOWN || resp.Error != "internal error" {
		t.Fatalf("resp = %+v, want DOWN with internal error", resp)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	act := New()
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	act.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
