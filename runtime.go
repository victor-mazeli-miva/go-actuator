package actuator

import (
	"net/http"
	"runtime"
)

type RuntimeResponse struct {
	Goroutines uint64 `json:"goroutines"`
	HeapAlloc  uint64 `json:"heapAlloc"`
	GCCount    uint32 `json:"gcCount"`
	CPUCount   int    `json:"cpuCount"`
}

func collectRuntimeStats() RuntimeResponse {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return RuntimeResponse{
		Goroutines: uint64(runtime.NumGoroutine()),
		HeapAlloc:  mem.HeapAlloc,
		GCCount:    mem.NumGC,
		CPUCount:   runtime.NumCPU(),
	}
}

func (a *Actuator) handleRuntime(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, collectRuntimeStats())
}
