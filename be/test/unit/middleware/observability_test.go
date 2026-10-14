package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"be/common/httpx"
	"be/common/reqctx"
	"be/internal/middleware"
	"be/public/handlers"
)

func newEngine(t *testing.T, logs *bytes.Buffer, trusted []string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewJSONHandler(logs, nil))
	r := gin.New()
	if err := r.SetTrustedProxies(trusted); err != nil {
		t.Fatalf("trusted proxies: %v", err)
	}
	r.Use(middleware.RequestID(), middleware.Recovery(logger), middleware.AccessLog(logger), middleware.SecurityHeaders())
	return r
}

func TestRequestIDIsGeneratedOrReusedAndEchoed(t *testing.T) {
	var logs bytes.Buffer
	r := newEngine(t, &logs, nil)
	var seen string
	r.GET("/x", func(c *gin.Context) { seen = reqctx.RequestID(c); c.Status(http.StatusNoContent) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	generated := w.Header().Get(middleware.RequestIDHeader)
	if generated == "" || generated != seen {
		t.Fatalf("expected generated request id echoed and in context, header=%q ctx=%q", generated, seen)
	}

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(middleware.RequestIDHeader, "edge-123")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if got := w.Header().Get(middleware.RequestIDHeader); got != "edge-123" {
		t.Fatalf("inbound request id should be reused, got %q", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(middleware.RequestIDHeader, "bad id\nwith newline")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if got := w.Header().Get(middleware.RequestIDHeader); got == "" || strings.ContainsAny(got, " \n") {
		t.Fatalf("malformed inbound id must be replaced, got %q", got)
	}
}

func TestAccessLogIsStructuredWithoutQueryString(t *testing.T) {
	var logs bytes.Buffer
	r := newEngine(t, &logs, nil)
	r.GET("/items/:id", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/items/42?token=secret", nil))

	var line map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(logs.Bytes()), &line); err != nil {
		t.Fatalf("log line is not JSON: %v (%s)", err, logs.String())
	}
	if line["msg"] != "http request" || line["level"] != "WARN" || line["route"] != "/items/:id" ||
		line["path"] != "/items/42" || line["status"].(float64) != http.StatusTeapot {
		t.Fatalf("unexpected access log: %v", line)
	}
	if strings.Contains(logs.String(), "secret") {
		t.Fatalf("query string leaked into logs: %s", logs.String())
	}
}

func TestRecoveryReturnsGeneric500AndLogs(t *testing.T) {
	var logs bytes.Buffer
	r := newEngine(t, &logs, nil)
	r.GET("/boom", func(*gin.Context) { panic("db password=hunter2") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))
	if w.Code != http.StatusInternalServerError || strings.Contains(w.Body.String(), "hunter2") {
		t.Fatalf("expected generic 500, got %d %s", w.Code, w.Body.String())
	}
	if !strings.Contains(logs.String(), "panic recovered") {
		t.Fatalf("panic not logged: %s", logs.String())
	}
}

func TestSecurityHeadersAndHSTSOnlyOverHTTPS(t *testing.T) {
	var logs bytes.Buffer
	r := newEngine(t, &logs, nil)
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	if w.Header().Get("X-Content-Type-Options") != "nosniff" || w.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatalf("missing security headers: %v", w.Header())
	}
	if w.Header().Get("Strict-Transport-Security") != "" {
		t.Fatal("HSTS must not be sent over plain HTTP")
	}

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !strings.HasPrefix(w.Header().Get("Strict-Transport-Security"), "max-age=") {
		t.Fatalf("expected HSTS over https, got %v", w.Header())
	}
}

func TestClientIPIgnoresForwardedHeaderFromUntrustedPeer(t *testing.T) {
	var logs bytes.Buffer
	var got string
	handler := func(c *gin.Context) { got = httpx.ClientIP(c); c.Status(http.StatusOK) }

	direct := newEngine(t, &logs, []string{"10.0.0.0/8"})
	direct.GET("/ip", handler)
	req := httptest.NewRequest(http.MethodGet, "/ip", nil)
	req.RemoteAddr = "203.0.113.7:5555"
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	direct.ServeHTTP(httptest.NewRecorder(), req)
	if got != "203.0.113.7" {
		t.Fatalf("spoofed X-Forwarded-For from untrusted peer accepted: %q", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/ip", nil)
	req.RemoteAddr = "10.1.2.3:5555" // trusted proxy (e.g. nginx on the docker network)
	req.Header.Set("X-Forwarded-For", "198.51.100.9")
	direct.ServeHTTP(httptest.NewRecorder(), req)
	if got != "198.51.100.9" {
		t.Fatalf("X-Forwarded-For from trusted proxy ignored: %q", got)
	}
}

func TestHealthAndReadiness(t *testing.T) {
	gin.SetMode(gin.TestMode)
	run := func(checks ...handlers.ReadinessCheck) (int, map[string]any, int) {
		h := handlers.NewHealthHandler(checks...)
		r := gin.New()
		r.GET("/health", h.Health)
		r.GET("/ready", h.Ready)
		live := httptest.NewRecorder()
		r.ServeHTTP(live, httptest.NewRequest(http.MethodGet, "/health", nil))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ready", nil))
		var body map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &body)
		return w.Code, body, live.Code
	}

	ok := handlers.ReadinessCheck{Name: "database", Check: func(context.Context) error { return nil }}
	down := handlers.ReadinessCheck{Name: "redis", Check: func(context.Context) error { return errors.New("dial tcp: refused") }}

	code, body, live := run(ok)
	if code != http.StatusOK || live != http.StatusOK || body["status"] != "ok" {
		t.Fatalf("healthy: got %d %v (live %d)", code, body, live)
	}

	code, body, live = run(ok, down)
	checks := body["checks"].(map[string]any)
	if code != http.StatusServiceUnavailable || live != http.StatusOK || checks["redis"] != "down" || checks["database"] != "ok" {
		t.Fatalf("degraded: got %d %v (live %d)", code, body, live)
	}
	if strings.Contains(strings.ToLower(body["status"].(string)+checks["redis"].(string)), "refused") {
		t.Fatal("dependency error details must not leak in the public readiness body")
	}
}
