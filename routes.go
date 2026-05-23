package actuator

func (a *Actuator) registerRoutes() {
	a.mux.HandleFunc("GET /health", a.handleHealth)
	a.mux.HandleFunc("GET /ready", a.handleReady)
	a.mux.HandleFunc("GET /live", a.handleLive)
	a.mux.HandleFunc("GET /runtime", a.handleRuntime)
}
