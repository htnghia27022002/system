package maps

import (
	"context"
	"strings"
	"unicode/utf8"

	"be/internal/config"
	mapsdto "be/internal/dto/maps"
	placemodel "be/internal/models/place"
	locationsvc "be/internal/services/location"
	placesvc "be/internal/services/place"
)

const (
	minQueryRunes = 2
	localLimit    = 8
)

type factory func(config.Config, HTTPDoer) Provider

var factories = map[string]factory{}

// Register adds a geocode adapter. Called from providers/ init.
func Register(id string, fn factory) {
	id = strings.ToLower(strings.TrimSpace(id))
	if id == "" || fn == nil {
		return
	}
	factories[id] = fn
}

// New builds the configured geocode adapter. Default is OSM Nominatim.
func New(cfg config.Config, httpDoer HTTPDoer) Provider {
	key := strings.ToLower(strings.TrimSpace(cfg.MapGeocodeProvider))
	if key == "" {
		key = "osm"
	}
	if fn := factories[key]; fn != nil {
		return fn(cfg, httpDoer)
	}
	if fn := factories["osm"]; fn != nil {
		return fn(cfg, httpDoer)
	}
	return nil
}

// Service searches saved Places/Locations and the configured geocode adapter.
type Service struct {
	cfg       config.Config
	places    *placesvc.Service
	locations *locationsvc.Service
}

func NewService(cfg config.Config, places *placesvc.Service, locations *locationsvc.Service) *Service {
	return &Service{cfg: cfg, places: places, locations: locations}
}

func (s *Service) Search(ctx context.Context, q mapsdto.SearchQuery) (*mapsdto.SearchResponse, error) {
	query := strings.TrimSpace(q.Q)
	if utf8.RuneCountInString(query) < minQueryRunes {
		return &mapsdto.SearchResponse{Items: []mapsdto.SearchHit{}}, nil
	}
	country := strings.ToUpper(strings.TrimSpace(q.CountryCode))

	items := make([]mapsdto.SearchHit, 0, localLimit*2)
	seen := map[string]struct{}{}

	placeHits, err := s.places.List(ctx, mapsdto.PlaceListQuery{
		Q:           query,
		CountryCode: country,
		Page:        1,
		Limit:       localLimit,
	})
	if err != nil {
		return nil, err
	}
	for _, row := range placeHits.Items {
		if !placemodel.ValidCoords(row.Lat, row.Lng) {
			continue
		}
		id := "place:" + row.ID
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		subtitle := strings.TrimSpace(row.LocationName)
		if unit := strings.TrimSpace(row.Unit); unit != "" {
			if subtitle != "" {
				subtitle += " · " + unit
			} else {
				subtitle = unit
			}
		}
		items = append(items, mapsdto.SearchHit{
			ID:         id,
			Kind:       "place",
			Title:      row.Name,
			Subtitle:   subtitle,
			Lat:        *row.Lat,
			Lng:        *row.Lng,
			PlaceID:    row.ID,
			LocationID: row.LocationID,
		})
	}

	locHits, err := s.locations.List(ctx, mapsdto.LocationListQuery{
		Q:           query,
		CountryCode: country,
		Page:        1,
		Limit:       localLimit,
	})
	if err != nil {
		return nil, err
	}
	for _, row := range locHits.Items {
		if !placemodel.ValidCoords(row.Lat, row.Lng) {
			continue
		}
		id := "location:" + row.ID
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		subtitle := strings.TrimSpace(row.Formatted)
		if subtitle == "" {
			subtitle = strings.TrimSpace(row.Street)
		}
		items = append(items, mapsdto.SearchHit{
			ID:         id,
			Kind:       "location",
			Title:      row.Name,
			Subtitle:   subtitle,
			Lat:        *row.Lat,
			Lng:        *row.Lng,
			LocationID: row.ID,
		})
	}

	if geocode := New(s.cfg, nil); geocode != nil {
		geoHits, geoErr := geocode.Search(ctx, query, country)
		if geoErr == nil {
			providerID := geocode.ID()
			for _, hit := range geoHits {
				id := "geocode:" + providerID + ":" + hit.ExternalID
				if _, ok := seen[id]; ok {
					continue
				}
				seen[id] = struct{}{}
				items = append(items, mapsdto.SearchHit{
					ID:       id,
					Kind:     "geocode",
					Title:    hit.Title,
					Subtitle: hit.Subtitle,
					Lat:      hit.Lat,
					Lng:      hit.Lng,
				})
			}
		}
	}

	return &mapsdto.SearchResponse{Items: items}, nil
}
