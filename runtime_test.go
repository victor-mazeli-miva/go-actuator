package actuator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRuntimeHandler(t *testing.T) {
	act := New()
	req := httptest.NewRequest(http.MethodGet, "/runtime", nil)
	rec := httptest.NewRecorder()

	act.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}

	var resp RuntimeResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Goroutines < 1 {
		t.Fatalf("goroutines = %d, want >= 1", resp.Goroutines)
	}
	if resp.CPUCount < 1 {
		t.Fatalf("cpuCount = %d, want >= 1", resp.CPUCount)
	}
}
