package adapters

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/uLesson-Education/go-actuator"
)

// MountChi registers actuator routes on a chi router at the given prefix.
// For example, MountChi(r, "/actuator", act) serves /actuator/health, /actuator/live, etc.
func MountChi(r chi.Router, prefix string, a *actuator.Actuator) {
	prefix = strings.TrimSuffix(prefix, "/")
	r.Mount(prefix, http.StripPrefix(prefix, a.Router()))
}
