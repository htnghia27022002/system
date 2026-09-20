package ingest

import (
	"context"
	"strings"
	"time"

	locationmodel "be/internal/models/location"
	newsmodel "be/internal/models/news"
	placemodel "be/internal/models/place"
)

func (s *Service) persistItem(ctx context.Context, actor *string, sourceID *string, sourceName string, item MappedItem) error {
	if strings.TrimSpace(item.PlaceName) == "" {
		return errSkipItem
	}

	loc, err := s.findOrCreateLocation(ctx, actor, item)
	if err != nil {
		return err
	}
	if loc == nil {
		return errSkipItem
	}

	placeKey := derivePlaceKey(item)
	if placeKey == "" {
		placeKey = strings.ToLower(strings.TrimSpace(item.PlaceName)) + "|"
	}

	catKey := resolveCategoryKey(item)
	cat, err := s.categories.GetByKey(ctx, catKey)
	if err != nil {
		return err
	}
	if cat == nil {
		cat, err = s.categories.GetByKey(ctx, "uncategorized")
		if err != nil {
			return err
		}
		if cat == nil {
			return errSkipItem
		}
	}

	lat, lng := item.PlaceLat, item.PlaceLng
	if !placemodel.ValidCoords(lat, lng) {
		lat, lng = loc.Lat, loc.Lng
	}

	existing, err := s.places.FindByLocationAndKey(ctx, loc.ID, placeKey)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	if existing == nil {
		status := placemodel.StatusPending
		if placemodel.ValidCoords(lat, lng) {
			status = placemodel.StatusActive
		}
		existing = &placemodel.Place{
			LocationID: loc.ID,
			CategoryID: cat.ID,
			PlaceKey:   placeKey,
			Name:       strings.TrimSpace(item.PlaceName),
			Unit:       strings.TrimSpace(item.Unit),
			Status:     status,
			Lat:        lat,
			Lng:        lng,
			Details:    item.Details,
			CreatedBy:  actor,
			UpdatedBy:  actor,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := s.places.Create(ctx, existing); err != nil {
			return err
		}
	} else {
		existing.Name = strings.TrimSpace(item.PlaceName)
		if item.Unit != "" {
			existing.Unit = strings.TrimSpace(item.Unit)
		}
		existing.CategoryID = cat.ID
		if len(item.Details) > 0 {
			if existing.Details == nil {
				existing.Details = map[string]any{}
			}
			for k, v := range item.Details {
				existing.Details[k] = v
			}
		}
		if existing.Status != placemodel.StatusHidden {
			if placemodel.ValidCoords(lat, lng) {
				existing.Lat = lat
				existing.Lng = lng
				if existing.Status == placemodel.StatusPending {
					existing.Status = placemodel.StatusActive
				}
			}
		}
		existing.UpdatedBy = actor
		existing.UpdatedAt = now
		if err := s.places.Update(ctx, existing); err != nil {
			return err
		}
	}

	url := strings.TrimSpace(item.OriginalURL)
	if url == "" {
		return nil
	}
	article, err := s.news.GetByOriginalURL(ctx, url)
	if err != nil {
		return err
	}
	if article == nil {
		article = &newsmodel.News{
			PlaceID:     existing.ID,
			CategoryID:  existing.CategoryID,
			SourceID:    sourceID,
			SourceName:  sourceName,
			OriginalURL: url,
			Title:       item.NewsTitle,
			Details:     map[string]any{},
			CreatedBy:   actor,
			UpdatedBy:   actor,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		return s.news.Create(ctx, article)
	}
	article.Title = item.NewsTitle
	article.SourceName = sourceName
	article.SourceID = sourceID
	article.UpdatedBy = actor
	article.UpdatedAt = now
	return s.news.Update(ctx, article)
}

func (s *Service) findOrCreateLocation(ctx context.Context, actor *string, item MappedItem) (*locationmodel.Location, error) {
	if key := strings.TrimSpace(item.LocationKey); key != "" {
		found, err := s.locations.FindByLocationKey(ctx, key)
		if err != nil {
			return nil, err
		}
		if found != nil {
			return found, nil
		}
	}

	var countryID, adminID *string
	if item.CountryCode != "" {
		country, err := s.countries.GetByCode(ctx, item.CountryCode)
		if err != nil {
			return nil, err
		}
		if country != nil {
			countryID = &country.ID
			if item.AdminCode != "" {
				div, err := s.divisions.GetByCountryAndCode(ctx, country.ID, item.AdminCode)
				if err != nil {
					return nil, err
				}
				if div != nil && div.IsActive {
					adminID = &div.ID
				}
			}
		}
	}

	if countryID != nil && adminID != nil && strings.TrimSpace(item.Street) != "" {
		found, err := s.locations.FindByCountryAdminStreet(ctx, *countryID, *adminID, item.Street)
		if err != nil {
			return nil, err
		}
		if found != nil {
			return found, nil
		}
	}

	if strings.TrimSpace(item.Formatted) != "" {
		found, err := s.locations.FindByFormatted(ctx, item.Formatted)
		if err != nil {
			return nil, err
		}
		if found != nil {
			return found, nil
		}
	}

	if placemodel.ValidCoords(item.LocationLat, item.LocationLng) {
		found, err := s.locations.FindByRoundedCoordsAndFormatted(ctx, *item.LocationLat, *item.LocationLng, item.Formatted)
		if err != nil {
			return nil, err
		}
		if found != nil {
			return found, nil
		}
	}

	hasIdentity := strings.TrimSpace(item.LocationKey) != "" ||
		strings.TrimSpace(item.Formatted) != "" ||
		strings.TrimSpace(item.Street) != "" ||
		strings.TrimSpace(item.LocationName) != ""
	if !hasIdentity {
		return nil, nil
	}

	name := strings.TrimSpace(item.LocationName)
	if name == "" {
		name = strings.TrimSpace(item.Formatted)
	}
	if name == "" {
		name = strings.TrimSpace(item.Street)
	}
	now := time.Now().UTC()
	loc := &locationmodel.Location{
		Name:            name,
		CountryID:       countryID,
		AdminDivisionID: adminID,
		Street:          strings.TrimSpace(item.Street),
		PostalCode:      strings.TrimSpace(item.PostalCode),
		Formatted:       strings.TrimSpace(item.Formatted),
		Lat:             item.LocationLat,
		Lng:             item.LocationLng,
		CreatedBy:       actor,
		UpdatedBy:       actor,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if key := strings.TrimSpace(item.LocationKey); key != "" {
		loc.LocationKey = &key
	}
	if err := s.locations.Create(ctx, loc); err != nil {
		return nil, err
	}
	return loc, nil
}
