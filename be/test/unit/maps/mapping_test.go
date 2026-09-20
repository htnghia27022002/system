package maps_test

import (
	"testing"

	mapsdto "be/internal/dto/maps"
)

func TestValidateFieldMappingAllowsDocumentedKeys(t *testing.T) {
	t.Parallel()
	mapping := map[string]any{
		"listPath": "data.items",
		"location": map[string]any{"name": "a", "address": "b", "formatted": "c"},
		"place":    map[string]any{"name": "n", "placeKey": "k"},
		"news":     map[string]any{"title": "t", "originalUrl": "u"},
		"details":  map[string]any{"phone": "p", "priceRange": "r"},
		"categoryMap": map[string]any{
			"nha-hang": "restaurant",
		},
	}
	if err := mapsdto.ValidateFieldMapping(mapping); err != nil {
		t.Fatalf("expected valid mapping, got %v", err)
	}
}

func TestValidateFieldMappingRejectsUnknownTargets(t *testing.T) {
	t.Parallel()
	cases := []map[string]any{
		{"unknownRoot": "x"},
		{"location": map[string]any{"invented": "x"}},
		{"place": map[string]any{"newColumn": "x"}},
		{"news": map[string]any{"bodyHtml": "x"}},
		{"details": map[string]any{"wifiPassword": "x"}},
	}
	for _, mapping := range cases {
		if err := mapsdto.ValidateFieldMapping(mapping); err == nil {
			t.Fatalf("expected error for mapping %#v", mapping)
		}
	}
}

func TestFilterDetailsDropsUnknownKeys(t *testing.T) {
	t.Parallel()
	got := mapsdto.FilterDetails(map[string]any{
		"phone":        "1",
		"wifiPassword": "secret",
	})
	if _, ok := got["wifiPassword"]; ok {
		t.Fatal("unknown detail key should be dropped")
	}
	if got["phone"] != "1" {
		t.Fatalf("expected phone kept, got %#v", got)
	}
}
