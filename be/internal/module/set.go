package module

import (
	"github.com/gin-gonic/gin"

	"be/internal/apidoc"
	"be/internal/database"
	"be/internal/database/seeders"
	"be/internal/queue"
	"be/internal/settings"
)

// Set is the list of enabled modules with helpers that collect their capabilities.
type Set []Module

// Names returns the enabled module names in registry order.
func (s Set) Names() []string {
	names := make([]string, 0, len(s))
	for _, m := range s {
		names = append(names, m.Name())
	}
	return names
}

// RegisterRoutes mounts every module's routes.
func (s Set) RegisterRoutes(r Routes) {
	for _, m := range s {
		m.RegisterRoutes(r)
	}
}

// GlobalMiddleware returns engine-level middleware from modules, in order.
func (s Set) GlobalMiddleware() []gin.HandlerFunc {
	var out []gin.HandlerFunc
	for _, m := range s {
		if p, ok := m.(MiddlewareProvider); ok {
			out = append(out, p.GlobalMiddleware()...)
		}
	}
	return out
}

// SeedContributions returns each module's permissions and seeders.
func (s Set) SeedContributions() []seeders.Contribution {
	var out []seeders.Contribution
	for _, m := range s {
		var c seeders.Contribution
		if p, ok := m.(PermissionProvider); ok {
			c.Permissions = p.Permissions()
		}
		if p, ok := m.(SeederProvider); ok {
			c.Seeders = p.Seeders()
		}
		if len(c.Permissions) > 0 || len(c.Seeders) > 0 {
			out = append(out, c)
		}
	}
	return out
}

// MigrationSets returns one versioned migration set per module that has SQL.
func (s Set) MigrationSets() []database.MigrationSet {
	var out []database.MigrationSet
	for _, m := range s {
		if p, ok := m.(MigrationProvider); ok {
			out = append(out, database.MigrationSet{Name: m.Name(), FS: p.Migrations()})
		}
	}
	return out
}

// Streams returns every JetStream stream declared by modules.
func (s Set) Streams() []queue.Stream {
	var out []queue.Stream
	for _, m := range s {
		if p, ok := m.(StreamProvider); ok {
			out = append(out, p.Streams()...)
		}
	}
	return out
}

// Pollers returns modules that can process queued work without NATS.
func (s Set) Pollers() []Poller {
	var out []Poller
	for _, m := range s {
		if p, ok := m.(Poller); ok {
			out = append(out, p)
		}
	}
	return out
}

// Close releases module resources in reverse order.
func (s Set) Close() {
	for i := len(s) - 1; i >= 0; i-- {
		if c, ok := s[i].(Closer); ok {
			c.Close()
		}
	}
}

// Settings collects the settings declared by modules.
func (s Set) Settings() []settings.Definition {
	var out []settings.Definition
	for _, m := range s {
		if p, ok := m.(SettingsProvider); ok {
			out = append(out, p.Settings()...)
		}
	}
	return out
}

// APIDoc collects the route docs of modules that provide them.
func (s Set) APIDoc() []apidoc.Op {
	var out []apidoc.Op
	for _, m := range s {
		if p, ok := m.(APIDocProvider); ok {
			out = append(out, p.APIDoc()...)
		}
	}
	return out
}
