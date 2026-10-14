package settings_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	apperrors "be/common/errors"
	"be/internal/settings"
)

type memStore struct {
	values map[string]json.RawMessage
	loads  int
}

func (m *memStore) All(context.Context) (map[string]json.RawMessage, error) {
	m.loads++
	out := map[string]json.RawMessage{}
	for k, v := range m.values {
		out[k] = v
	}
	return out, nil
}

func (m *memStore) Upsert(_ context.Context, values map[string]json.RawMessage, _ string) error {
	for k, v := range values {
		m.values[k] = v
	}
	return nil
}

func newService(t *testing.T, store *memStore) *settings.Service {
	t.Helper()
	svc, err := settings.NewService(store, settings.CoreDefinitions(), []settings.Definition{
		{Key: "billing.currency", Type: settings.TypeString, Default: "USD", Group: "billing", Public: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestDefaultsOverridesAndPublicView(t *testing.T) {
	store := &memStore{values: map[string]json.RawMessage{}}
	svc := newService(t, store)
	ctx := context.Background()

	if !svc.Bool(ctx, settings.KeyRegistrationEnabled) || svc.Int(ctx, settings.KeyAuditRetentionDays) != 180 {
		t.Fatal("defaults not applied")
	}
	if _, err := svc.Update(ctx, map[string]any{
		settings.KeyRegistrationEnabled: false,
		settings.KeyAuditRetentionDays:  float64(30), // JSON numbers decode as float64
		"billing.currency":              "  VND ",
	}, "admin"); err != nil {
		t.Fatal(err)
	}
	if svc.Bool(ctx, settings.KeyRegistrationEnabled) || svc.Int(ctx, settings.KeyAuditRetentionDays) != 30 || svc.String(ctx, "billing.currency") != "VND" {
		t.Fatal("overrides not applied")
	}

	public, _ := svc.Public(ctx)
	if _, leaked := public[settings.KeyAuditRetentionDays]; leaked {
		t.Fatal("non-public setting exposed")
	}
	if public[settings.KeyRegistrationEnabled] != false || public["billing.currency"] != "VND" {
		t.Fatalf("public view: %v", public)
	}

	list, _ := svc.List(ctx)
	if len(list) != 4 || list[0].Group > list[len(list)-1].Group {
		t.Fatalf("list should hold every setting grouped: %+v", list)
	}
}

func TestUpdateRejectsUnknownAndInvalidValues(t *testing.T) {
	svc := newService(t, &memStore{values: map[string]json.RawMessage{}})
	ctx := context.Background()
	cases := []map[string]any{
		{"nope.key": true},
		{settings.KeyRegistrationEnabled: "yes"},
		{settings.KeyAuditRetentionDays: float64(1)},
		{settings.KeyAuditRetentionDays: 7.5},
		{settings.KeySupportEmail: string(make([]byte, 300))},
		{},
	}
	for _, values := range cases {
		if _, err := svc.Update(ctx, values, ""); !errors.Is(err, apperrors.ErrBadRequest) {
			t.Fatalf("%v: expected bad request, got %v", values, err)
		}
	}
}

func TestDuplicateKeysFailStartup(t *testing.T) {
	_, err := settings.NewService(&memStore{}, settings.CoreDefinitions(), settings.CoreDefinitions())
	if err == nil {
		t.Fatal("duplicate setting keys must fail")
	}
}

func TestCacheExpires(t *testing.T) {
	store := &memStore{values: map[string]json.RawMessage{}}
	svc := newService(t, store)
	now := time.Now()
	svc.SetClock(func() time.Time { return now })
	ctx := context.Background()
	svc.Bool(ctx, settings.KeyRegistrationEnabled)
	svc.Bool(ctx, settings.KeyRegistrationEnabled)
	if store.loads != 1 {
		t.Fatalf("reads within the TTL should hit the cache, loads=%d", store.loads)
	}
	store.values[settings.KeyRegistrationEnabled] = json.RawMessage("false") // written by another replica
	now = now.Add(time.Minute)
	if svc.Bool(ctx, settings.KeyRegistrationEnabled) {
		t.Fatal("cache must refresh after the TTL")
	}
}
