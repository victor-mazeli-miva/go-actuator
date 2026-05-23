package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gin-gonic/gin"
	"github.com/uLesson-Education/go-actuator"
)

type healthResponse struct {
	Status string `json:"status"`
}

func TestMountNetHTTP(t *testing.T) {
	act := actuator.New()
	mux := http.NewServeMux()
	Mount(mux, "/actuator", act)

	req := httptest.NewRequest(http.MethodGet, "/actuator/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != "UP" {
		t.Fatalf("status = %q, want UP", resp.Status)
	}
}

func TestMountChi(t *testing.T) {
	act := actuator.New()
	r := chi.NewRouter()
	MountChi(r, "/actuator", act)

	req := httptest.NewRequest(http.MethodGet, "/actuator/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != "UP" {
		t.Fatalf("status = %q, want UP", resp.Status)
	}
}

func TestMountGin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	act := actuator.New()
	r := gin.New()
	MountGin(r, "/actuator", act)

	req := httptest.NewRequest(http.MethodGet, "/actuator/live", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != "UP" {
		t.Fatalf("status = %q, want UP", resp.Status)
	}
}
