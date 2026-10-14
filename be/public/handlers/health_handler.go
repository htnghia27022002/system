package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"be/common/response"
)

// ReadinessCheck reports whether one dependency is usable.
type ReadinessCheck struct {
	Name  string
	Check func(ctx context.Context) error
}

// HealthHandler serves liveness (process up) and readiness (dependencies up).
type HealthHandler struct {
	checks  []ReadinessCheck
	timeout time.Duration
}

func NewHealthHandler(checks ...ReadinessCheck) *HealthHandler {
	return &HealthHandler{checks: checks, timeout: 2 * time.Second}
}

// Health is the liveness probe: 200 while the process serves HTTP.
func (h *HealthHandler) Health(c *gin.Context) {
	response.JSON(c, http.StatusOK, HealthResponse{Status: "ok"})
}

// Ready is the readiness probe: 200 when every check passes, else 503. Failure
// details go to the log, not the (public) response.
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	status := http.StatusOK
	results := make(map[string]string, len(h.checks))
	for _, check := range h.checks {
		if err := check.Check(ctx); err != nil {
			status = http.StatusServiceUnavailable
			results[check.Name] = "down"
			slog.WarnContext(ctx, "readiness check failed", slog.String("check", check.Name), slog.String("error", err.Error()))
			continue
		}
		results[check.Name] = "ok"
	}

	overall := "ok"
	if status != http.StatusOK {
		overall = "unavailable"
	}
	response.JSON(c, status, ReadyResponse{Status: overall, Checks: results})
}

// HealthResponse is GET /api/health.
type HealthResponse struct {
	Status string `json:"status"`
}

// ReadyResponse is GET /api/ready: status "ok" or "unavailable" (503) and
// each dependency as "ok" / "down" (no error details).
type ReadyResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}
