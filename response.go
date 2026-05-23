package actuator

import (
	"encoding/json"
	"net/http"
)

const (
	statusUP   = "UP"
	statusDOWN = "DOWN"
)

type HealthResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks,omitempty"`
}

// IsHealthy reports whether the aggregated status is UP.
func (r HealthResponse) IsHealthy() bool {
	return r.Status == statusUP
}

// CheckIsHealthy reports whether an individual check status is UP.
func CheckIsHealthy(status string) bool {
	return status == statusUP
}

type errorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}

func recoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				writeJSON(w, http.StatusInternalServerError, errorResponse{
					Status: statusDOWN,
					Error:  "internal error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
