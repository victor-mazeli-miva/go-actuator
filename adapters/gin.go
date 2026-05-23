package adapters

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uLesson-Education/go-actuator"
)

// MountGin registers actuator routes on a Gin engine at the given prefix.
// For example, MountGin(r, "/actuator", act) serves /actuator/health, /actuator/live, etc.
func MountGin(r gin.IRoutes, prefix string, a *actuator.Actuator) {
	prefix = strings.TrimSuffix(prefix, "/")
	r.Any(prefix+"/*path", gin.WrapH(http.StripPrefix(prefix, a.Router())))
}
