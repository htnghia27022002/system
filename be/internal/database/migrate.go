package database

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"be/internal/config"
)

// MigrationURL builds a postgres URL for golang-migrate.
func MigrationURL(cfg config.Config) string {
	user := url.UserPassword(cfg.DBUser, cfg.DBPass)
	u := url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}
	q := u.Query()
	q.Set("sslmode", cfg.DBSSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func migrationsSourceURL() (string, error) {
	dir, err := findMigrationsDir()
	if err != nil {
		return "", err
	}
	return "file://" + filepath.ToSlash(dir), nil
}

func findMigrationsDir() (string, error) {
	if dir := os.Getenv("MIGRATIONS_PATH"); dir != "" {
		return filepath.Abs(dir)
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	dir := wd
	for {
		candidate := filepath.Join(dir, "migrations")
		if st, statErr := os.Stat(candidate); statErr == nil && st.IsDir() {
			return filepath.Abs(candidate)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("migrations directory not found (set MIGRATIONS_PATH)")
}

func newMigrator(cfg config.Config) (*migrate.Migrate, error) {
	sourceURL, err := migrationsSourceURL()
	if err != nil {
		return nil, err
	}

	m, err := migrate.New(sourceURL, MigrationURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}
	return m, nil
}

// RunMigrations applies all pending SQL migrations.
func RunMigrations(cfg config.Config) error {
	m, err := newMigrator(cfg)
	if err != nil {
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			fmt.Fprintf(os.Stderr, "migrate source close: %v\n", srcErr)
		}
		if dbErr != nil {
			fmt.Fprintf(os.Stderr, "migrate database close: %v\n", dbErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// MigrateDown rolls back the last migration.
func MigrateDown(cfg config.Config) error {
	m, err := newMigrator(cfg)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate down: %w", err)
	}
	return nil
}

// MigrateDrop removes all migration state and drops schema objects from down scripts.
func MigrateDrop(cfg config.Config) error {
	m, err := newMigrator(cfg)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Drop(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate drop: %w", err)
	}
	return nil
}

// MigrateVersion returns the current migration version and dirty flag.
func MigrateVersion(cfg config.Config) (uint, bool, error) {
	m, err := newMigrator(cfg)
	if err != nil {
		return 0, false, err
	}
	defer func() { _, _ = m.Close() }()

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("migrate version: %w", err)
	}
	return version, dirty, nil
}

// MigrateForce sets the migration version without running SQL (recovery helper).
func MigrateForce(cfg config.Config, version int) error {
	m, err := newMigrator(cfg)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Force(version); err != nil {
		return fmt.Errorf("migrate force: %w", err)
	}
	return nil
}

// MigrateSteps runs n up (positive) or down (negative) migrations.
func MigrateSteps(cfg config.Config, steps int) error {
	m, err := newMigrator(cfg)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Steps(steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate steps (%d): %w", steps, err)
	}
	return nil
}

// ParseForceVersion parses the force version argument for the CLI.
func ParseForceVersion(raw string) (int, error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid force version %q: %w", raw, err)
	}
	return v, nil
}
