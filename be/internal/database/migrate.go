package database

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"be/internal/config"
	"be/migrations"
)

// CoreMigrationSet is the name of the core schema set (table schema_migrations).
const CoreMigrationSet = "core"

// MigrationSet is one independently versioned group of SQL migrations. Core uses
// golang-migrate's default table; each module records its version in
// schema_migrations_<name>, so modules can be added or removed without touching
// another set's history.
type MigrationSet struct {
	Name string
	FS   fs.FS
}

// CoreMigrations returns the core schema set (identity, RBAC, sessions).
func CoreMigrations() MigrationSet {
	return MigrationSet{Name: CoreMigrationSet, FS: migrations.FS}
}

func (s MigrationSet) table() string {
	if s.Name == CoreMigrationSet {
		return "schema_migrations"
	}
	return "schema_migrations_" + s.Name
}

// MigrationURL builds a postgres URL for golang-migrate.
// When cfg.DBURL is set (from DB_URL), it is used as-is.
func MigrationURL(cfg config.Config) string {
	if cfg.DBURL != "" {
		return cfg.DBURL
	}

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

func setURL(cfg config.Config, set MigrationSet) (string, error) {
	u, err := url.Parse(MigrationURL(cfg))
	if err != nil {
		return "", fmt.Errorf("parse database url: %w", err)
	}
	q := u.Query()
	q.Set("x-migrations-table", set.table())
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// SQL is embedded in the binary (go:embed), so deploys need no migrations folder
// and `go run` always picks up the files on disk.
func newMigrator(cfg config.Config, set MigrationSet) (*migrate.Migrate, error) {
	source, err := iofs.New(set.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("migrations %s: source: %w", set.Name, err)
	}
	dbURL, err := setURL(cfg, set)
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, dbURL)
	if err != nil {
		return nil, fmt.Errorf("migrations %s: %w", set.Name, err)
	}
	return m, nil
}

func withMigrator(cfg config.Config, set MigrationSet, fn func(*migrate.Migrate) error) error {
	m, err := newMigrator(cfg, set)
	if err != nil {
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			fmt.Fprintf(os.Stderr, "migrate %s source close: %v\n", set.Name, srcErr)
		}
		if dbErr != nil {
			fmt.Fprintf(os.Stderr, "migrate %s database close: %v\n", set.Name, dbErr)
		}
	}()
	return fn(m)
}

// RunMigrations applies pending migrations for core, then each given module set.
func RunMigrations(cfg config.Config, modules ...MigrationSet) error {
	sets := append([]MigrationSet{CoreMigrations()}, modules...)
	for _, set := range sets {
		err := withMigrator(cfg, set, func(m *migrate.Migrate) error {
			if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("migrate %s up: %w", set.Name, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// MigrateDown rolls back the last migration of one set.
func MigrateDown(cfg config.Config, set MigrationSet) error {
	return MigrateSteps(cfg, set, -1)
}

// MigrateDrop drops every object in the database (all sets) and their version tables.
func MigrateDrop(cfg config.Config, set MigrationSet) error {
	return withMigrator(cfg, set, func(m *migrate.Migrate) error {
		if err := m.Drop(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate %s drop: %w", set.Name, err)
		}
		return nil
	})
}

// MigrateVersion returns the current version and dirty flag of one set.
func MigrateVersion(cfg config.Config, set MigrationSet) (uint, bool, error) {
	var (
		version uint
		dirty   bool
	)
	err := withMigrator(cfg, set, func(m *migrate.Migrate) error {
		v, d, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				return nil
			}
			return fmt.Errorf("migrate %s version: %w", set.Name, err)
		}
		version, dirty = v, d
		return nil
	})
	return version, dirty, err
}

// MigrateForce sets one set's version without running SQL (recovery helper).
func MigrateForce(cfg config.Config, set MigrationSet, version int) error {
	return withMigrator(cfg, set, func(m *migrate.Migrate) error {
		if err := m.Force(version); err != nil {
			return fmt.Errorf("migrate %s force: %w", set.Name, err)
		}
		return nil
	})
}

// MigrateSteps runs n up (positive) or down (negative) migrations of one set.
func MigrateSteps(cfg config.Config, set MigrationSet, steps int) error {
	return withMigrator(cfg, set, func(m *migrate.Migrate) error {
		if err := m.Steps(steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate %s steps (%d): %w", set.Name, steps, err)
		}
		return nil
	})
}

// ParseForceVersion parses the force version argument for the CLI.
func ParseForceVersion(raw string) (int, error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid force version %q: %w", raw, err)
	}
	return v, nil
}
