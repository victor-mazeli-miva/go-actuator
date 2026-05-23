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
