// Package adapters provides thin helpers for mounting an actuator on popular Go HTTP routers.
package adapters

import (
	"net/http"
	"strings"

	"github.com/uLesson-Education/go-actuator"
)

// Mount registers actuator routes on a ServeMux at the given prefix.
// For example, Mount(mux, "/actuator", act) serves /actuator/health, /actuator/live, etc.
func Mount(mux *http.ServeMux, prefix string, a *actuator.Actuator) {
	prefix = strings.TrimSuffix(prefix, "/")
	mux.Handle(prefix+"/", http.StripPrefix(prefix, a.Router()))
}
