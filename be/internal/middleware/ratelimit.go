package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter is a fixed-window request counter keyed by client IP + route.
//
// It is in-memory, so the limit applies per API instance; behind N replicas the
// effective limit is N× the configured value. That is enough to blunt password
// guessing (account lockout is DB-backed and global); swap in a shared store if
// you need an exact cluster-wide limit.
type RateLimiter struct {
	limit  int
	window time.Duration
	now    func() time.Time

	mu       sync.Mutex
	counters map[string]*windowCounter
}

type windowCounter struct {
	start time.Time
	count int
}

// sweepThreshold bounds memory: expired windows are dropped once the map grows past it.
const sweepThreshold = 10_000

// NewRateLimiter allows `limit` requests per `window` per key. limit <= 0 disables it.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{limit: limit, window: window, now: time.Now, counters: map[string]*windowCounter{}}
}

// Allow records one request for key and reports whether it is within the limit,
// plus how long until the current window resets.
func (l *RateLimiter) Allow(key string) (bool, time.Duration) {
	if l == nil || l.limit <= 0 {
		return true, 0
	}
	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.counters) > sweepThreshold {
		for k, wc := range l.counters {
			if now.Sub(wc.start) >= l.window {
				delete(l.counters, k)
			}
		}
	}

	wc, ok := l.counters[key]
	if !ok || now.Sub(wc.start) >= l.window {
		wc = &windowCounter{start: now}
		l.counters[key] = wc
	}
	wc.count++
	return wc.count <= l.limit, wc.start.Add(l.window).Sub(now)
}

// Middleware rejects requests over the limit with 429 and a Retry-After header.
func (l *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, resetIn := l.Allow(c.ClientIP() + " " + c.FullPath())
		if allowed {
			c.Next()
			return
		}
		c.Header("Retry-After", strconv.Itoa(int(math.Ceil(resetIn.Seconds()))))
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests, try again later"})
	}
}

// SetClock replaces the time source (tests).
func (l *RateLimiter) SetClock(now func() time.Time) {
	l.now = now
}
