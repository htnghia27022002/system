package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"be/internal/middleware"
)

func TestRateLimiterFixedWindow(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	l := middleware.NewRateLimiter(2, time.Minute)
	l.SetClock(func() time.Time { return now })

	for i := 0; i < 2; i++ {
		if ok, _ := l.Allow("a"); !ok {
			t.Fatalf("request %d should pass", i+1)
		}
	}
	ok, resetIn := l.Allow("a")
	if ok || resetIn != time.Minute {
		t.Fatalf("third request must be limited with reset in 1m, got ok=%v reset=%v", ok, resetIn)
	}
	if ok, _ := l.Allow("b"); !ok {
		t.Fatal("keys are limited independently")
	}

	now = now.Add(time.Minute)
	if ok, _ := l.Allow("a"); !ok {
		t.Fatal("new window must reset the count")
	}
}

func TestRateLimiterDisabledWhenLimitIsZero(t *testing.T) {
	l := middleware.NewRateLimiter(0, time.Minute)
	for i := 0; i < 100; i++ {
		if ok, _ := l.Allow("a"); !ok {
			t.Fatal("limit 0 must disable the limiter")
		}
	}
}

func TestRateLimiterMiddlewareReturns429WithRetryAfter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	l := middleware.NewRateLimiter(1, time.Minute)
	r := gin.New()
	r.POST("/login", l.Middleware(), func(c *gin.Context) { c.Status(http.StatusOK) })
	r.POST("/register", l.Middleware(), func(c *gin.Context) { c.Status(http.StatusOK) })

	do := func(path, ip string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.RemoteAddr = ip + ":1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	if w := do("/login", "203.0.113.1"); w.Code != http.StatusOK {
		t.Fatalf("first request: %d", w.Code)
	}
	w := do("/login", "203.0.113.1")
	if w.Code != http.StatusTooManyRequests || w.Header().Get("Retry-After") != "60" {
		t.Fatalf("expected 429 + Retry-After 60, got %d %q", w.Code, w.Header().Get("Retry-After"))
	}
	if w := do("/login", "203.0.113.2"); w.Code != http.StatusOK {
		t.Fatalf("other IP must not share the limit: %d", w.Code)
	}
	if w := do("/register", "203.0.113.1"); w.Code != http.StatusOK {
		t.Fatalf("other route must not share the limit: %d", w.Code)
	}
}
