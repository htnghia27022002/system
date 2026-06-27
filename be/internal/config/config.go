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

	if cfg.DBPass == "" {
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
