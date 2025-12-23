package handler

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

// Probes holds the state for Kubernetes probes.
type Probes struct {
	ready   atomic.Bool
	healthy atomic.Bool
}

// NewProbes creates a new Probes instance.
func NewProbes() *Probes {
	p := &Probes{}
	p.healthy.Store(true)
	return p
}

// SetReady marks the application as ready to receive traffic.
func (p *Probes) SetReady(ready bool) {
	p.ready.Store(ready)
}

// SetHealthy marks the application as healthy.
func (p *Probes) SetHealthy(healthy bool) {
	p.healthy.Store(healthy)
}

// IsReady returns true if the application is ready.
func (p *Probes) IsReady() bool {
	return p.ready.Load()
}

// IsHealthy returns true if the application is healthy.
func (p *Probes) IsHealthy() bool {
	return p.healthy.Load()
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// ReadinessResponse represents the readiness check response.
type ReadinessResponse struct {
	Status string `json:"status"`
	Ready  bool   `json:"ready"`
}

// LivenessResponse represents the liveness check response.
type LivenessResponse struct {
	Status  string `json:"status"`
	Healthy bool   `json:"healthy"`
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

// readinessProbe handles Kubernetes readiness probe.
// Returns 200 OK if the application is ready to receive traffic.
func (h *Handler) readinessProbe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.probes == nil || !h.probes.IsReady() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ReadinessResponse{
			Status: "not ready",
			Ready:  false,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ReadinessResponse{
		Status: "ready",
		Ready:  true,
	})
}

// livenessProbe handles Kubernetes liveness probe.
// Returns 200 OK if the application is alive and healthy.
func (h *Handler) livenessProbe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.probes == nil || !h.probes.IsHealthy() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(LivenessResponse{
			Status:  "unhealthy",
			Healthy: false,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LivenessResponse{
		Status:  "healthy",
		Healthy: true,
	})
}
