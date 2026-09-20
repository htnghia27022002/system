package ingest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	mapsdto "be/internal/dto/maps"
	categorymodel "be/internal/models/category"
)

// MappedItem is one list item after field mapping onto the fixed schema.
type MappedItem struct {
	LocationName        string
	LocationKey         string
	CountryCode         string
	AdminCode           string
	Street              string
	PostalCode          string
	Formatted           string
	LocationLat         *float64
	LocationLng         *float64
	PlaceName           string
	PlaceCategoryRaw    string
	PlaceKey            string
	Unit                string
	PlaceLat            *float64
	PlaceLng            *float64
	NewsTitle           string
	OriginalURL         string
	Details             map[string]any
	CategoryMap         map[string]string
}

func extractItems(root any, mapping map[string]any) ([]any, error) {
	listPath := mapsdto.MappingString(mapping, "listPath")
	node := root
	if listPath != "" {
		var ok bool
		node, ok = walkPath(root, listPath)
		if !ok {
			return nil, fmt.Errorf("listPath %q not found", listPath)
		}
	}
	switch v := node.(type) {
	case []any:
		return v, nil
	default:
		return nil, fmt.Errorf("payload is not a JSON list")
	}
}

func mapItem(raw any, mapping map[string]any) (MappedItem, error) {
	item := MappedItem{Details: map[string]any{}, CategoryMap: map[string]string{}}
	obj, ok := raw.(map[string]any)
	if !ok {
		return item, fmt.Errorf("item is not an object")
	}

	loc := asObject(mapping["location"])
	place := asObject(mapping["place"])
	news := asObject(mapping["news"])
	details := asObject(mapping["details"])

	item.LocationName = readPathString(obj, mapsdto.MappingString(loc, "name"))
	item.LocationKey = readPathString(obj, mapsdto.MappingString(loc, "locationKey"))
	item.CountryCode = strings.ToUpper(readPathString(obj, mapsdto.MappingString(loc, "countryCode")))
	item.AdminCode = readPathString(obj, mapsdto.MappingString(loc, "adminCode"))
	item.Street = readPathString(obj, mapsdto.MappingString(loc, "street"))
	item.PostalCode = readPathString(obj, mapsdto.MappingString(loc, "postalCode"))
	item.Formatted = readPathString(obj, mapsdto.MappingString(loc, "formatted"))
	if item.Formatted == "" {
		item.Formatted = readPathString(obj, mapsdto.MappingString(loc, "address"))
	}
	item.LocationLat = readPathFloat(obj, mapsdto.MappingString(loc, "lat"))
	item.LocationLng = readPathFloat(obj, mapsdto.MappingString(loc, "lng"))

	item.PlaceName = readPathString(obj, mapsdto.MappingString(place, "name"))
	item.PlaceCategoryRaw = readPathString(obj, mapsdto.MappingString(place, "category"))
	item.PlaceKey = readPathString(obj, mapsdto.MappingString(place, "placeKey"))
	item.Unit = readPathString(obj, mapsdto.MappingString(place, "unit"))
	item.PlaceLat = readPathFloat(obj, mapsdto.MappingString(place, "lat"))
	item.PlaceLng = readPathFloat(obj, mapsdto.MappingString(place, "lng"))

	item.NewsTitle = readPathString(obj, mapsdto.MappingString(news, "title"))
	item.OriginalURL = readPathString(obj, mapsdto.MappingString(news, "originalUrl"))

	for destKey, pathRaw := range details {
		if _, ok := mapsdto.DetailsAllowList[destKey]; !ok {
			continue
		}
		path, _ := pathRaw.(string)
		if value, ok := walkPath(obj, path); ok && value != nil {
			item.Details[destKey] = value
		}
	}
	item.Details = mapsdto.FilterDetails(item.Details)

	if rawMap, ok := mapping["categoryMap"].(map[string]any); ok {
		for k, v := range rawMap {
			if s, ok := v.(string); ok {
				item.CategoryMap[k] = s
			}
		}
	}
	return item, nil
}

func resolveCategoryKey(item MappedItem) string {
	raw := strings.TrimSpace(item.PlaceCategoryRaw)
	if mapped, ok := item.CategoryMap[raw]; ok && strings.TrimSpace(mapped) != "" {
		raw = mapped
	}
	key := categorymodel.FromJSONKey(strings.TrimSpace(raw))
	switch key {
	case categorymodel.KeyRoomRental, categorymodel.KeyRestaurant, categorymodel.KeyHotel, categorymodel.KeyEatery, categorymodel.KeyUncategorized:
		return key
	default:
		return categorymodel.KeyUncategorized
	}
}

func derivePlaceKey(item MappedItem) string {
	if key := strings.TrimSpace(item.PlaceKey); key != "" {
		return key
	}
	name := strings.ToLower(strings.TrimSpace(item.PlaceName))
	unit := strings.ToLower(strings.TrimSpace(item.Unit))
	if name == "" && unit == "" {
		return ""
	}
	return name + "|" + unit
}

func asObject(raw any) map[string]any {
	obj, _ := raw.(map[string]any)
	if obj == nil {
		return map[string]any{}
	}
	return obj
}

func walkPath(root any, path string) (any, bool) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, false
	}
	cur := root
	for _, part := range strings.Split(path, ".") {
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		next, exists := obj[part]
		if !exists {
			return nil, false
		}
		cur = next
	}
	return cur, true
}

func readPathString(root map[string]any, path string) string {
	value, ok := walkPath(root, path)
	if !ok || value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return v.String()
	case float64:
		return strings.TrimSpace(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		return strconv.FormatBool(v)
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func readPathFloat(root map[string]any, path string) *float64 {
	value, ok := walkPath(root, path)
	if !ok || value == nil {
		return nil
	}
	switch v := value.(type) {
	case float64:
		return &v
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return nil
		}
		return &f
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return nil
		}
		return &f
	default:
		return nil
	}
}
