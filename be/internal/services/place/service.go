package place

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	apperrors "be/common/errors"
	"be/common/utils"
	mapsdto "be/internal/dto/maps"
	categorymodel "be/internal/models/category"
	placemodel "be/internal/models/place"
	"be/internal/repository/interfaces"
	addresssvc "be/internal/services/address"
	locationsvc "be/internal/services/location"
)

// Service lists Places, returns detail, and applies operator CRUD.
type Service struct {
	places     interfaces.PlaceRepository
	locations  *locationsvc.Service
	categories interfaces.CategoryRepository
	address    *addresssvc.Service
}

func NewService(
	places interfaces.PlaceRepository,
	locations *locationsvc.Service,
	categories interfaces.CategoryRepository,
	address *addresssvc.Service,
) *Service {
	return &Service{
		places:     places,
		locations:  locations,
		categories: categories,
		address:    address,
	}
}

func (s *Service) CountByLocationID(ctx context.Context, locationID string) (int64, error) {
	return s.places.CountByLocationID(ctx, strings.TrimSpace(locationID))
}

func (s *Service) List(ctx context.Context, q mapsdto.PlaceListQuery) (*mapsdto.PlaceListResponse, error) {
	if q.Manage {
		return s.listAdmin(ctx, q)
	}
	page, limit := mapsdto.NormalizeList(q.Page, q.Limit)
	filter := interfaces.PlaceListFilter{
		Q:           strings.TrimSpace(q.Q),
		CountryCode: strings.ToUpper(strings.TrimSpace(q.CountryCode)),
		Page:        page,
		Limit:       limit,
	}
	if cat := strings.TrimSpace(q.Category); cat != "" {
		if !categorymodel.IsP1FilterKey(cat) {
			return nil, fmt.Errorf("%w: invalid category", apperrors.ErrBadRequest)
		}
		filter.CategoryKey = categorymodel.FromJSONKey(cat)
	}
	if id := strings.TrimSpace(q.AdminDivisionID); id != "" {
		div, err := s.address.DivisionByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if div == nil {
			return nil, fmt.Errorf("%w: unknown adminDivisionId", apperrors.ErrBadRequest)
		}
		filter.AdminPathPrefix = div.Path
	}

	rows, total, err := s.places.ListPinnable(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]mapsdto.PlacePin, 0, len(rows))
	for _, row := range rows {
		if !placemodel.ValidCoords(row.Place.Lat, row.Place.Lng) {
			continue
		}
		items = append(items, mapsdto.PlacePin{
			ID:           row.Place.ID,
			Name:         row.Place.Name,
			CategoryID:   row.Place.CategoryID,
			Category:     categorymodel.ToJSONKey(row.CategoryKey),
			Status:       placemodel.StatusToJSON(row.Place.Status),
			Lat:          row.Place.Lat,
			Lng:          row.Place.Lng,
			LocationID:   row.LocationID,
			LocationName: row.LocationName,
			Unit:         row.Place.Unit,
			PlaceKey:     row.Place.PlaceKey,
		})
	}
	return &mapsdto.PlaceListResponse{
		Items:   items,
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasMore: int64(page*limit) < total,
	}, nil
}

func (s *Service) Get(ctx context.Context, id string) (*mapsdto.PlaceDetail, error) {
	detail, err := s.places.GetDetail(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, apperrors.ErrNotFound
	}
	return toPlaceDetail(detail), nil
}

func (s *Service) Create(ctx context.Context, userID string, req mapsdto.CreatePlaceRequest) (*mapsdto.PlaceDetail, error) {
	loc, err := s.requireLocation(ctx, req.LocationID)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 200 {
		return nil, fmt.Errorf("%w: name is required (max 200)", apperrors.ErrBadRequest)
	}
	cat, err := s.requireCategory(ctx, req.Category)
	if err != nil {
		return nil, err
	}
	if err := validateOptionalCoords(req.Lat, req.Lng); err != nil {
		return nil, err
	}
	placeKey := placemodel.DerivePlaceKey(name, req.Unit, req.PlaceKey)
	if placeKey == "|" || placeKey == "" {
		return nil, fmt.Errorf("%w: placeKey could not be derived", apperrors.ErrBadRequest)
	}
	existing, err := s.places.FindByLocationAndKey(ctx, loc.ID, placeKey)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: placeKey already exists at this location", apperrors.ErrConflict)
	}
	status := placemodel.StatusPending
	if req.Status != "" {
		parsed, ok := placemodel.StatusFromJSON(req.Status)
		if !ok {
			return nil, fmt.Errorf("%w: invalid status", apperrors.ErrBadRequest)
		}
		status = parsed
	} else if placemodel.ValidCoords(req.Lat, req.Lng) {
		status = placemodel.StatusActive
	}
	if status == placemodel.StatusActive && !placemodel.ValidCoords(req.Lat, req.Lng) {
		return nil, fmt.Errorf("%w: cannot set active without valid coordinates", apperrors.ErrBadRequest)
	}
	details := req.Details
	if details == nil {
		details = map[string]any{}
	}
	now := time.Now().UTC()
	actor := utils.StringPtr(userID)
	place := &placemodel.Place{
		LocationID: loc.ID,
		CategoryID: cat.ID,
		PlaceKey:   placeKey,
		Name:       name,
		Unit:       strings.TrimSpace(req.Unit),
		Status:     status,
		Lat:        req.Lat,
		Lng:        req.Lng,
		Details:    details,
		CreatedBy:  actor,
		UpdatedBy:  actor,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.places.Create(ctx, place); err != nil {
		return nil, err
	}
	return s.Get(ctx, place.ID)
}

func (s *Service) Patch(ctx context.Context, id, userID string, req mapsdto.PatchPlaceRequest) (*mapsdto.PlaceDetail, error) {
	place, err := s.places.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if place == nil {
		return nil, apperrors.ErrNotFound
	}
	if req.LocationID != nil {
		loc, err := s.requireLocation(ctx, *req.LocationID)
		if err != nil {
			return nil, err
		}
		place.LocationID = loc.ID
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" || len(name) > 200 {
			return nil, fmt.Errorf("%w: name is required (max 200)", apperrors.ErrBadRequest)
		}
		place.Name = name
	}
	if req.Unit != nil {
		place.Unit = strings.TrimSpace(*req.Unit)
	}
	if req.Category != nil {
		cat, err := s.requireCategory(ctx, *req.Category)
		if err != nil {
			return nil, err
		}
		place.CategoryID = cat.ID
	}
	if req.PlaceKey != nil || req.Name != nil || req.Unit != nil || req.LocationID != nil {
		explicit := place.PlaceKey
		if req.PlaceKey != nil {
			explicit = *req.PlaceKey
		}
		place.PlaceKey = placemodel.DerivePlaceKey(place.Name, place.Unit, explicit)
		dup, err := s.places.FindByLocationAndKey(ctx, place.LocationID, place.PlaceKey)
		if err != nil {
			return nil, err
		}
		if dup != nil && dup.ID != place.ID {
			return nil, fmt.Errorf("%w: placeKey already exists at this location", apperrors.ErrConflict)
		}
	}
	if req.Lat != nil || req.Lng != nil {
		lat, lng := place.Lat, place.Lng
		if req.Lat != nil {
			lat = req.Lat
		}
		if req.Lng != nil {
			lng = req.Lng
		}
		if err := validateOptionalCoords(lat, lng); err != nil {
			return nil, err
		}
		place.Lat = lat
		place.Lng = lng
	}
	if req.Details != nil {
		place.Details = *req.Details
		if place.Details == nil {
			place.Details = map[string]any{}
		}
	}
	if req.Status != nil {
		want, ok := placemodel.StatusFromJSON(strings.TrimSpace(*req.Status))
		if !ok {
			return nil, fmt.Errorf("%w: invalid status", apperrors.ErrBadRequest)
		}
		if want == placemodel.StatusPending && req.Status != nil && strings.TrimSpace(*req.Status) == "pending" {
			// Operators may keep pending when coords are still missing.
		}
		if want == placemodel.StatusActive && !placemodel.ValidCoords(place.Lat, place.Lng) {
			return nil, fmt.Errorf("%w: cannot set active without valid coordinates", apperrors.ErrBadRequest)
		}
		place.Status = want
	}
	place.UpdatedBy = utils.StringPtr(userID)
	place.UpdatedAt = time.Now().UTC()
	if err := s.places.Update(ctx, place); err != nil {
		return nil, err
	}
	return s.Get(ctx, place.ID)
}

func (s *Service) PatchStatus(ctx context.Context, id, userID, status string) (*mapsdto.PlaceDetail, error) {
	return s.Patch(ctx, id, userID, mapsdto.PatchPlaceRequest{Status: &status})
}

func (s *Service) Delete(ctx context.Context, id string) error {
	place, err := s.places.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	if place == nil {
		return apperrors.ErrNotFound
	}
	return s.places.Delete(ctx, place.ID)
}

func (s *Service) listAdmin(ctx context.Context, q mapsdto.PlaceListQuery) (*mapsdto.PlaceListResponse, error) {
	page, limit := mapsdto.NormalizeList(q.Page, q.Limit)
	filter := interfaces.PlaceAdminListFilter{
		LocationID: strings.TrimSpace(q.LocationID),
		Q:          strings.TrimSpace(q.Q),
		Page:       page,
		Limit:      limit,
	}
	if status := strings.TrimSpace(q.Status); status != "" {
		parsed, ok := placemodel.StatusFromJSON(status)
		if !ok {
			return nil, fmt.Errorf("%w: invalid status", apperrors.ErrBadRequest)
		}
		filter.Status = &parsed
	}
	rows, total, err := s.places.ListAdmin(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]mapsdto.PlacePin, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapsdto.PlacePin{
			ID:           row.Place.ID,
			Name:         row.Place.Name,
			CategoryID:   row.Place.CategoryID,
			Category:     categorymodel.ToJSONKey(row.CategoryKey),
			Status:       placemodel.StatusToJSON(row.Place.Status),
			Lat:          row.Place.Lat,
			Lng:          row.Place.Lng,
			LocationID:   row.LocationID,
			LocationName: row.LocationName,
			Unit:         row.Place.Unit,
			PlaceKey:     row.Place.PlaceKey,
		})
	}
	return &mapsdto.PlaceListResponse{
		Items:   items,
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasMore: int64(page*limit) < total,
	}, nil
}

func (s *Service) requireLocation(ctx context.Context, id string) (*interfaces.PlaceLocationRow, error) {
	loc, err := s.locations.Get(ctx, strings.TrimSpace(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, fmt.Errorf("%w: unknown locationId", apperrors.ErrBadRequest)
		}
		return nil, err
	}
	return &interfaces.PlaceLocationRow{ID: loc.ID, Name: loc.Name}, nil
}

func (s *Service) requireCategory(ctx context.Context, key string) (*categorymodel.Category, error) {
	stored := categorymodel.FromJSONKey(strings.TrimSpace(key))
	if !categorymodel.IsCatalogKey(stored) {
		return nil, fmt.Errorf("%w: invalid category", apperrors.ErrBadRequest)
	}
	cat, err := s.categories.GetByKey(ctx, stored)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, fmt.Errorf("%w: unknown category", apperrors.ErrBadRequest)
	}
	return cat, nil
}

func validateOptionalCoords(lat, lng *float64) error {
	if lat == nil && lng == nil {
		return nil
	}
	if lat == nil || lng == nil {
		return fmt.Errorf("%w: lat and lng must be set together", apperrors.ErrBadRequest)
	}
	if *lat < -90 || *lat > 90 || *lng < -180 || *lng > 180 {
		return fmt.Errorf("%w: coordinates are out of range", apperrors.ErrBadRequest)
	}
	return nil
}

func (s *Service) ListCategories(ctx context.Context) (*mapsdto.CategoryListResponse, error) {
	rows, err := s.categories.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]mapsdto.CategoryItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapsdto.CategoryItem{
			ID:        row.ID,
			Key:       categorymodel.ToJSONKey(row.Key),
			Name:      row.Name,
			NameLocal: row.NameLocal,
			SortOrder: row.SortOrder,
		})
	}
	return &mapsdto.CategoryListResponse{Items: items}, nil
}

func toPlaceDetail(detail *interfaces.PlaceDetailRow) *mapsdto.PlaceDetail {
	news := make([]mapsdto.NewsItem, 0, len(detail.News))
	for _, item := range detail.News {
		catKey := detail.NewsCategoryKey[item.ID]
		if catKey == "" {
			catKey = detail.CategoryKey
		}
		news = append(news, mapsdto.NewsItem{
			ID:          item.ID,
			Title:       item.Title,
			OriginalURL: item.OriginalURL,
			CategoryID:  item.CategoryID,
			Category:    categorymodel.ToJSONKey(catKey),
			SourceName:  item.SourceName,
			SourceID:    item.SourceID,
			CreatedAt:   utils.FormatRFC3339(item.CreatedAt),
			UpdatedAt:   utils.FormatRFC3339(item.UpdatedAt),
		})
	}
	details := detail.Place.Details
	if details == nil {
		details = map[string]any{}
	}
	return &mapsdto.PlaceDetail{
		ID:         detail.Place.ID,
		Name:       detail.Place.Name,
		Unit:       detail.Place.Unit,
		PlaceKey:   detail.Place.PlaceKey,
		CategoryID: detail.Place.CategoryID,
		Category:   categorymodel.ToJSONKey(detail.CategoryKey),
		Status:     placemodel.StatusToJSON(detail.Place.Status),
		Lat:        detail.Place.Lat,
		Lng:        detail.Place.Lng,
		Details:    details,
		Location: mapsdto.LocationSummary{
			ID:              detail.Location.ID,
			Name:            detail.Location.Name,
			LocationKey:     detail.Location.LocationKey,
			CountryCode:     detail.Location.CountryCode,
			AdminDivisionID: detail.Location.AdminDivisionID,
			AdminPath:       detail.Location.AdminPath,
			Street:          detail.Location.Street,
			PostalCode:      detail.Location.PostalCode,
			Formatted:       detail.Location.Formatted,
			Lat:             detail.Location.Lat,
			Lng:             detail.Location.Lng,
		},
		News:      news,
		CreatedAt: utils.FormatRFC3339(detail.Place.CreatedAt),
		UpdatedAt: utils.FormatRFC3339(detail.Place.UpdatedAt),
	}
}
