package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port      string
	DBURL     string // when set (DB_URL), preferred over discrete DB_* fields
	DBHost    string
	DBPort    string
	DBName    string
	DBUser    string
	DBPass    string
	DBSSLMode string

	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	CORSOrigins []string

	OAuthGoogleClientID     string
	OAuthGoogleClientSecret string
	OAuthRedirectURL        string
	OAuthAllowedProviders   []string

	RedisURL string
	Cache    CacheConfig

	// AutoMigrate runs migrations and seeders on API start (AUTO_MIGRATE, default true).
	// Set false in multi-replica deploys and run ./cmd/migrate + ./cmd/seed as a release step.
	AutoMigrate bool

	// UploadDir is the local filesystem root for uploaded media (avatars). Owned by be/.env.
	UploadDir string

	// Logging: LOG_FORMAT text|json (default text), LOG_LEVEL debug|info|warn|error.
	LogFormat string
	LogLevel  string

	// TrustedProxies are the proxy IPs/CIDRs whose X-Forwarded-For is believed
	// (TRUSTED_PROXIES, comma-separated). Default: loopback + private networks.
	TrustedProxies []string

	Auth AuthConfig

	// AppName labels outgoing mail (APP_NAME, default "System App").
	AppName string
	// AppURL is the public FE origin used to build links in emails
	// (APP_URL, e.g. https://app.example.com). It is a URL, not a code dependency.
	AppURL string

	Mail MailConfig

	// Modules lists the enabled optional modules (MODULES env, comma-separated).
	// ModulesSet is false when MODULES is unset, meaning "use registry defaults".
	Modules    []string
	ModulesSet bool
}

// AuthConfig hardens the browser auth flow.
type AuthConfig struct {
	// RefreshCookie sends the refresh token as an httpOnly cookie (path /api/auth)
	// instead of the JSON body (AUTH_REFRESH_COOKIE, default true).
	RefreshCookie bool
	// CookieDomain is optional (AUTH_COOKIE_DOMAIN); empty = host-only cookie.
	CookieDomain string
	// CookieSameSite: lax (default), strict, or none (cross-site FE; needs HTTPS).
	CookieSameSite string
	// CookieSecure: auto (default — Secure when the request is HTTPS), true, false.
	CookieSecure string
	// RateLimitPerMinute caps auth requests per client IP and route (0 disables).
	RateLimitPerMinute int
	// MaxFailedLogins locks an account after this many consecutive failures (0 disables).
	MaxFailedLogins int
	// LockoutDuration is how long a locked account stays locked.
	LockoutDuration time.Duration
	// RequireEmailVerification blocks password login until the email is verified
	// (AUTH_REQUIRE_EMAIL_VERIFICATION, default false).
	RequireEmailVerification bool
}

// MailConfig selects the outgoing mail transport.
type MailConfig struct {
	// Driver: log (default — write mail to the log, for dev) or smtp.
	Driver string
	// From is the RFC 5322 sender, e.g. "System App <no-reply@example.com>".
	From         string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	// SMTPTLS: starttls (default), tls (implicit, usually port 465) or none.
	SMTPTLS string
}

type CacheConfig struct {
	Enabled    bool
	Driver     string
	DefaultTTL time.Duration
	FileDir    string
	RedisURL   string
}

type fileConfig struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Name    string `yaml:"name"`
		User    string `yaml:"user"`
		SSLMode string `yaml:"sslMode"`
	} `yaml:"database"`
	JWT struct {
		AccessTTL  string `yaml:"accessTtl"`
		RefreshTTL string `yaml:"refreshTtl"`
	} `yaml:"jwt"`
	CORS struct {
		Origins []string `yaml:"origins"`
	} `yaml:"cors"`
	Redis struct {
		URL string `yaml:"url"`
	} `yaml:"redis"`
	Cache struct {
		Enabled    bool   `yaml:"enabled"`
		Driver     string `yaml:"driver"`
		DefaultTTL string `yaml:"defaultTtl"`
		File       struct {
			Dir string `yaml:"dir"`
		} `yaml:"file"`
		Redis struct {
			URL string `yaml:"url"`
		} `yaml:"redis"`
	} `yaml:"cache"`
	OAuth struct {
		RedirectURL      string   `yaml:"redirectUrl"`
		AllowedProviders []string `yaml:"allowedProviders"`
	} `yaml:"oauth"`
}

// Load reads public settings from config.yaml, then applies env overrides.
// Env wins over YAML. Secrets (DB password, JWT, OAuth client credentials) must come from env.
func Load() Config {
	fc := loadFileConfig()

	cfg := Config{
		Port:      firstNonEmpty(os.Getenv("PORT"), fc.Server.Port, "8080"),
		DBURL:     strings.TrimSpace(os.Getenv("DB_URL")),
		DBHost:    firstNonEmpty(os.Getenv("DB_HOST"), fc.Database.Host, "localhost"),
		DBPort:    firstNonEmpty(os.Getenv("DB_PORT"), fc.Database.Port, "5432"),
		DBName:    firstNonEmpty(os.Getenv("DB_NAME"), fc.Database.Name, "myapp_pg"),
		DBUser:    firstNonEmpty(os.Getenv("DB_USER"), fc.Database.User, "postgres"),
		DBPass:    os.Getenv("DB_PASSWORD"),
		DBSSLMode: firstNonEmpty(os.Getenv("DB_SSLMODE"), fc.Database.SSLMode, "disable"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		CORSOrigins: firstNonEmptySlice(splitEnv("CORS_ORIGINS", ""), fc.CORS.Origins, []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		}),

		OAuthGoogleClientID:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
		OAuthGoogleClientSecret: os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
		OAuthRedirectURL: firstNonEmpty(
			os.Getenv("OAUTH_REDIRECT_URL"),
			fc.OAuth.RedirectURL,
			"http://localhost:3000/auth/callback",
		),
		OAuthAllowedProviders: firstNonEmptySlice(
			splitEnv("OAUTH_ALLOWED_PROVIDERS", ""),
			fc.OAuth.AllowedProviders,
			[]string{"google"},
		),

		RedisURL: firstNonEmpty(os.Getenv("REDIS_URL"), fc.Redis.URL, ""),
	}

	cfg.Cache = CacheConfig{
		Enabled: parseBoolEnv(os.Getenv("CACHE_ENABLED"), fc.Cache.Enabled),
		Driver: firstNonEmpty(
			os.Getenv("CACHE_DRIVER"),
			fc.Cache.Driver,
			"file",
		),
		DefaultTTL: parseDuration(
			firstNonEmpty(os.Getenv("CACHE_DEFAULT_TTL"), fc.Cache.DefaultTTL, "5m"),
			5*time.Minute,
		),
		FileDir: firstNonEmpty(
			os.Getenv("CACHE_FILE_DIR"),
			fc.Cache.File.Dir,
			"storage/cache",
		),
		RedisURL: firstNonEmpty(
			os.Getenv("CACHE_REDIS_URL"),
			fc.Cache.Redis.URL,
			cfg.RedisURL,
		),
	}

	cfg.JWTAccessTTL = parseDuration(
		firstNonEmpty(os.Getenv("JWT_ACCESS_TTL"), fc.JWT.AccessTTL, "15m"),
		15*time.Minute,
	)
	cfg.JWTRefreshTTL = parseDuration(
		firstNonEmpty(os.Getenv("JWT_REFRESH_TTL"), fc.JWT.RefreshTTL, "168h"),
		7*24*time.Hour,
	)

	cfg.AutoMigrate = parseBoolEnv(os.Getenv("AUTO_MIGRATE"), true)

	cfg.UploadDir = firstNonEmpty(os.Getenv("UPLOAD_DIR"), "data/uploads")

	cfg.LogFormat = strings.ToLower(firstNonEmpty(os.Getenv("LOG_FORMAT"), "text"))
	cfg.LogLevel = strings.ToLower(firstNonEmpty(os.Getenv("LOG_LEVEL"), "info"))
	cfg.TrustedProxies = firstNonEmptySlice(splitEnv("TRUSTED_PROXIES", ""), nil, []string{
		"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
	})

	cfg.Auth = AuthConfig{
		RefreshCookie:      parseBoolEnv(os.Getenv("AUTH_REFRESH_COOKIE"), true),
		CookieDomain:       strings.TrimSpace(os.Getenv("AUTH_COOKIE_DOMAIN")),
		CookieSameSite:     strings.ToLower(firstNonEmpty(os.Getenv("AUTH_COOKIE_SAMESITE"), "lax")),
		CookieSecure:       strings.ToLower(firstNonEmpty(os.Getenv("AUTH_COOKIE_SECURE"), "auto")),
		RateLimitPerMinute: parseNonNegativeIntEnv(os.Getenv("AUTH_RATE_LIMIT_PER_MINUTE"), 20),
		MaxFailedLogins:    parseNonNegativeIntEnv(os.Getenv("AUTH_MAX_FAILED_LOGINS"), 5),
		LockoutDuration:    parseDuration(os.Getenv("AUTH_LOCKOUT_DURATION"), 15*time.Minute),

		RequireEmailVerification: parseBoolEnv(os.Getenv("AUTH_REQUIRE_EMAIL_VERIFICATION"), false),
	}

	cfg.AppName = firstNonEmpty(os.Getenv("APP_NAME"), "System App")
	cfg.AppURL = strings.TrimRight(firstNonEmpty(os.Getenv("APP_URL"), "http://localhost:3000"), "/")
	cfg.Mail = MailConfig{
		Driver:       strings.ToLower(firstNonEmpty(os.Getenv("MAIL_DRIVER"), "log")),
		From:         firstNonEmpty(os.Getenv("MAIL_FROM"), cfg.AppName+" <no-reply@localhost>"),
		SMTPHost:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:     parseNonNegativeIntEnv(os.Getenv("SMTP_PORT"), 587),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPTLS:      strings.ToLower(firstNonEmpty(os.Getenv("SMTP_TLS"), "starttls")),
	}

	if _, ok := os.LookupEnv("MODULES"); ok {
		cfg.ModulesSet = true
		cfg.Modules = splitEnv("MODULES", "")
	}

	// Discrete DB_* password default only when not using a full connection URL.
	if cfg.DBURL == "" && cfg.DBPass == "" {
		cfg.DBPass = "postgres"
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "dev-jwt-secret-change-in-production"
	}

	return cfg
}

func loadFileConfig() fileConfig {
	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fileConfig{}
	}

	var fc fileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		fmt.Fprintf(os.Stderr, "config: failed to parse %s: %v\n", path, err)
		return fileConfig{}
	}
	return fc
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func firstNonEmptySlice(candidates ...[]string) []string {
	for _, list := range candidates {
		if len(list) > 0 {
			return list
		}
	}
	return nil
}

func splitEnv(key, fallback string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		raw = fallback
	}
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	if seconds, err := strconv.Atoi(raw); err == nil {
		return time.Duration(seconds) * time.Second
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	return fallback
}

func parseBoolEnv(raw string, yamlDefault bool) bool {
	raw = strings.TrimSpace(strings.ToLower(raw))
	switch raw {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return yamlDefault
	}
}

func parseIntEnv(raw string, yamlDefault int, fallback int) int {
	raw = strings.TrimSpace(raw)
	if raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			return n
		}
	}
	if yamlDefault > 0 {
		return yamlDefault
	}
	return fallback
}

// parseNonNegativeIntEnv returns raw as an int >= 0, or fallback when unset/invalid.
// Unlike parseIntEnv, an explicit 0 is kept (it disables the feature).
func parseNonNegativeIntEnv(raw string, fallback int) int {
	if n, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil && n >= 0 {
		return n
	}
	return fallback
}
