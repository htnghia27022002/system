package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"be/common/reqctx"
)

const (
	// RefreshCookieName holds the refresh token when cookie mode is on.
	RefreshCookieName = "refresh_token"
	// refreshCookiePath limits the cookie to the auth endpoints that need it.
	refreshCookiePath = "/api/auth"
	// CSRFHeader must accompany cookie-authenticated refresh/logout. Browsers only
	// send custom headers cross-site after a CORS preflight, which CORS_ORIGINS gates.
	CSRFHeader = "X-Requested-With"
)

// RefreshCookieConfig controls how the refresh token travels to browsers.
type RefreshCookieConfig struct {
	// Enabled sends the refresh token as an httpOnly cookie instead of JSON.
	Enabled bool
	Domain  string
	// SameSite: lax (default), strict, none (cross-site FE; forces Secure).
	SameSite string
	// Secure: auto (Secure when the request is HTTPS), true, false.
	Secure string
	MaxAge time.Duration
}

func (cfg RefreshCookieConfig) sameSite() http.SameSite {
	switch strings.ToLower(cfg.SameSite) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func (cfg RefreshCookieConfig) secure(c *gin.Context) bool {
	if cfg.sameSite() == http.SameSiteNoneMode {
		return true // browsers reject SameSite=None without Secure
	}
	switch strings.ToLower(cfg.Secure) {
	case "true":
		return true
	case "false":
		return false
	default:
		return reqctx.IsHTTPS(c)
	}
}

func (cfg RefreshCookieConfig) set(c *gin.Context, token string) {
	c.SetSameSite(cfg.sameSite())
	c.SetCookie(RefreshCookieName, token, int(cfg.MaxAge.Seconds()), refreshCookiePath, cfg.Domain, cfg.secure(c), true)
}

func (cfg RefreshCookieConfig) clear(c *gin.Context) {
	c.SetSameSite(cfg.sameSite())
	c.SetCookie(RefreshCookieName, "", -1, refreshCookiePath, cfg.Domain, cfg.secure(c), true)
}

// read returns the refresh token from the cookie, if cookie mode is on and it is present.
func (cfg RefreshCookieConfig) read(c *gin.Context) string {
	if !cfg.Enabled {
		return ""
	}
	token, err := c.Cookie(RefreshCookieName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(token)
}
