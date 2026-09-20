package maps

import (
	"fmt"
	"strings"
)

var locationMappingKeys = map[string]struct{}{
	"name":        {},
	"locationKey": {},
	"countryCode": {},
	"adminCode":   {},
	"street":      {},
	"postalCode":  {},
	"formatted":   {},
	"address":     {},
	"lat":         {},
	"lng":         {},
}

var placeMappingKeys = map[string]struct{}{
	"name":     {},
	"category": {},
	"placeKey": {},
	"unit":     {},
	"lat":      {},
	"lng":      {},
}

var newsMappingKeys = map[string]struct{}{
	"title":       {},
	"originalUrl": {},
}

// DetailsAllowList is the only keys stored on places.details.
var DetailsAllowList = map[string]struct{}{
	"description": {},
	"phone":       {},
	"hours":       {},
	"website":     {},
	"priceRange":  {},
}

var topLevelMappingKeys = map[string]struct{}{
	"listPath":    {},
	"location":    {},
	"place":       {},
	"news":        {},
	"details":     {},
	"categoryMap": {},
}

// ValidateFieldMapping rejects unknown stored-property targets on a Source mapping.
func ValidateFieldMapping(mapping map[string]any) error {
	if mapping == nil {
		return fmt.Errorf("fieldMapping is required")
	}
	for key := range mapping {
		if _, ok := topLevelMappingKeys[key]; !ok {
			return fmt.Errorf("unknown fieldMapping key %q", key)
		}
	}
	if err := validateStringMapKeys(mapping, "location", locationMappingKeys); err != nil {
		return err
	}
	if err := validateStringMapKeys(mapping, "place", placeMappingKeys); err != nil {
		return err
	}
	if err := validateStringMapKeys(mapping, "news", newsMappingKeys); err != nil {
		return err
	}
	if err := validateStringMapKeys(mapping, "details", DetailsAllowList); err != nil {
		return err
	}
	if raw, ok := mapping["categoryMap"]; ok && raw != nil {
		if _, ok := raw.(map[string]any); !ok {
			return fmt.Errorf("fieldMapping.categoryMap must be an object")
		}
	}
	if raw, ok := mapping["listPath"]; ok && raw != nil {
		if _, ok := raw.(string); !ok {
			return fmt.Errorf("fieldMapping.listPath must be a string")
		}
	}
	return nil
}

func validateStringMapKeys(mapping map[string]any, section string, allowed map[string]struct{}) error {
	raw, ok := mapping[section]
	if !ok || raw == nil {
		return nil
	}
	obj, ok := raw.(map[string]any)
	if !ok {
		return fmt.Errorf("fieldMapping.%s must be an object", section)
	}
	for key := range obj {
		if _, ok := allowed[key]; !ok {
			return fmt.Errorf("unknown fieldMapping.%s key %q", section, key)
		}
	}
	return nil
}

// FilterDetails keeps only allow-listed detail keys.
func FilterDetails(in map[string]any) map[string]any {
	out := map[string]any{}
	for key, value := range in {
		if _, ok := DetailsAllowList[key]; ok {
			out[key] = value
		}
	}
	return out
}

// MappingString reads a mapping path value as a trimmed string.
func MappingString(obj map[string]any, key string) string {
	if obj == nil {
		return ""
	}
	raw, ok := obj[key]
	if !ok || raw == nil {
		return ""
	}
	s, ok := raw.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}
