// Package settings holds runtime-editable application settings. Every setting
// is declared in code (key, type, default, visibility); the database only stores
// overrides. Unknown keys are rejected, so the admin UI can be generated from
// the declarations and a typo can never create a stray value.
package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	apperrors "be/common/errors"
	"be/internal/audit"
)

// Type is the value type of a setting.
type Type string

const (
	TypeString Type = "string"
	TypeBool   Type = "bool"
	TypeInt    Type = "int"
)

// Definition declares one setting.
type Definition struct {
	// Key is "<group>.<name>", e.g. "auth.registrationEnabled".
	Key         string `json:"key"`
	Type        Type   `json:"type"`
	Default     any    `json:"default"`
	Group       string `json:"group"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	// Public settings are exposed to anonymous clients via GET /api/app/info.
	// Never mark anything sensitive public.
	Public bool `json:"public"`
	// Min / Max bound TypeInt values (inclusive) when Max > Min.
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
	// MaxLength bounds TypeString values (0 = 500).
	MaxLength int `json:"maxLength,omitempty"`
}

// Value is a definition plus its effective value.
type Value struct {
	Definition
	Value any `json:"value"`
	// Overridden is true when the value comes from the database.
	Overridden bool `json:"overridden"`
}

// Store persists overrides (implemented by repository.SettingRepository).
type Store interface {
	All(ctx context.Context) (map[string]json.RawMessage, error)
	Upsert(ctx context.Context, values map[string]json.RawMessage, updatedBy string) error
}

// Core setting keys.
const (
	KeyRegistrationEnabled = "auth.registrationEnabled"
	KeySupportEmail        = "app.supportEmail"
	KeyAuditRetentionDays  = "audit.retentionDays"
)

// CoreDefinitions are the settings every project has. Modules add their own
// through module.SettingsProvider.
func CoreDefinitions() []Definition {
	return []Definition{
		{
			Key: KeyRegistrationEnabled, Type: TypeBool, Default: true, Group: "auth", Public: true,
			Label:       "Allow self sign-up",
			Description: "When off, only invited users can join (password and OAuth sign-up are refused).",
		},
		{
			Key: KeySupportEmail, Type: TypeString, Default: "", Group: "app", Public: true, MaxLength: 255,
			Label:       "Support email",
			Description: "Shown to users who need help. Leave empty to hide.",
		},
		{
			Key: KeyAuditRetentionDays, Type: TypeInt, Default: 180, Group: "audit", Min: 7, Max: 3650,
			Label:       "Audit log retention (days)",
			Description: "Entries older than this are deleted by the daily cleanup.",
		},
	}
}

// Service reads and updates settings with a short in-process cache, so a value
// changed on one replica reaches the others within cacheTTL.
type Service struct {
	store Store
	defs  map[string]Definition
	order []string

	mu       sync.Mutex
	cache    map[string]json.RawMessage
	loadedAt time.Time
	cacheTTL time.Duration
	now      func() time.Time
	auditor  audit.Recorder
}

// SetAuditor records setting changes (default: discard).
func (s *Service) SetAuditor(r audit.Recorder) { s.auditor = r }

// NewService validates the declarations (duplicate keys and bad defaults fail
// startup).
func NewService(store Store, defs ...[]Definition) (*Service, error) {
	s := &Service{store: store, defs: map[string]Definition{}, cacheTTL: 30 * time.Second, now: time.Now}
	for _, group := range defs {
		for _, d := range group {
			if _, dup := s.defs[d.Key]; dup {
				return nil, fmt.Errorf("settings: duplicate key %q", d.Key)
			}
			if _, err := d.normalize(d.Default); err != nil {
				return nil, fmt.Errorf("settings: default of %q: %w", d.Key, err)
			}
			s.defs[d.Key] = d
			s.order = append(s.order, d.Key)
		}
	}
	sort.SliceStable(s.order, func(i, j int) bool {
		gi, gj := s.defs[s.order[i]].Group, s.defs[s.order[j]].Group
		if gi != gj {
			return gi < gj
		}
		return false
	})
	return s, nil
}

// List returns every setting with its effective value, grouped by Group.
func (s *Service) List(ctx context.Context) ([]Value, error) {
	stored, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Value, 0, len(s.order))
	for _, key := range s.order {
		out = append(out, s.resolve(s.defs[key], stored))
	}
	return out, nil
}

// Public returns the public settings as key → value (safe for anonymous clients).
func (s *Service) Public(ctx context.Context) (map[string]any, error) {
	stored, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	for _, key := range s.order {
		if d := s.defs[key]; d.Public {
			out[key] = s.resolve(d, stored).Value
		}
	}
	return out, nil
}

// Bool returns a bool setting, falling back to its default on any error.
func (s *Service) Bool(ctx context.Context, key string) bool {
	v, _ := s.get(ctx, key).(bool)
	return v
}

// Int returns an int setting, falling back to its default on any error.
func (s *Service) Int(ctx context.Context, key string) int {
	v, _ := s.get(ctx, key).(int)
	return v
}

// String returns a string setting, falling back to its default on any error.
func (s *Service) String(ctx context.Context, key string) string {
	v, _ := s.get(ctx, key).(string)
	return v
}

// Update validates and stores values (partial update). It returns the keys
// that changed.
func (s *Service) Update(ctx context.Context, values map[string]any, updatedBy string) ([]string, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("%w: no settings to update", apperrors.ErrBadRequest)
	}
	encoded := make(map[string]json.RawMessage, len(values))
	keys := make([]string, 0, len(values))
	for key, raw := range values {
		d, ok := s.defs[key]
		if !ok {
			return nil, fmt.Errorf("%w: unknown setting %q", apperrors.ErrBadRequest, key)
		}
		v, err := d.normalize(raw)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %v", apperrors.ErrBadRequest, key, err)
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		encoded[key] = b
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if err := s.store.Upsert(ctx, encoded, updatedBy); err != nil {
		return nil, err
	}
	s.invalidate()
	if s.auditor != nil {
		s.auditor.Record(ctx, audit.Entry{Action: "settings.update", ResourceType: "settings",
			Metadata: map[string]any{"keys": keys}})
	}
	return keys, nil
}

func (s *Service) get(ctx context.Context, key string) any {
	d, ok := s.defs[key]
	if !ok {
		return nil
	}
	stored, err := s.load(ctx)
	if err != nil {
		return d.mustDefault()
	}
	return s.resolve(d, stored).Value
}

func (s *Service) resolve(d Definition, stored map[string]json.RawMessage) Value {
	value := Value{Definition: d, Value: d.mustDefault()}
	raw, ok := stored[d.Key]
	if !ok {
		return value
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return value
	}
	if v, err := d.normalize(decoded); err == nil {
		value.Value = v
		value.Overridden = true
	}
	return value
}

func (s *Service) load(ctx context.Context) (map[string]json.RawMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cache != nil && s.now().Sub(s.loadedAt) < s.cacheTTL {
		return s.cache, nil
	}
	stored, err := s.store.All(ctx)
	if err != nil {
		return nil, err
	}
	s.cache, s.loadedAt = stored, s.now()
	return stored, nil
}

func (s *Service) invalidate() {
	s.mu.Lock()
	s.cache = nil
	s.mu.Unlock()
}

func (d Definition) mustDefault() any {
	v, _ := d.normalize(d.Default)
	return v
}

// normalize coerces a JSON-decoded value to the definition's Go type.
func (d Definition) normalize(v any) (any, error) {
	switch d.Type {
	case TypeBool:
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("must be true or false")
		}
		return b, nil
	case TypeInt:
		var n int
		switch x := v.(type) {
		case int:
			n = x
		case float64:
			if x != float64(int(x)) {
				return nil, fmt.Errorf("must be a whole number")
			}
			n = int(x)
		default:
			return nil, fmt.Errorf("must be a number")
		}
		if d.Max > d.Min && (n < d.Min || n > d.Max) {
			return nil, fmt.Errorf("must be between %d and %d", d.Min, d.Max)
		}
		return n, nil
	case TypeString:
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("must be text")
		}
		str = strings.TrimSpace(str)
		max := d.MaxLength
		if max == 0 {
			max = 500
		}
		if len(str) > max {
			return nil, fmt.Errorf("must be at most %d characters", max)
		}
		return str, nil
	default:
		return nil, fmt.Errorf("unsupported type %q", d.Type)
	}
}

// SetClock replaces the time source (tests).
func (s *Service) SetClock(now func() time.Time) { s.now = now }
