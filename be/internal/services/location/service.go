package location

import (
	"context"
	"fmt"
	"strings"
	"time"

	apperrors "be/common/errors"
	"be/common/utils"
	mapsdto "be/internal/dto/maps"
	locationmodel "be/internal/models/location"
	"be/internal/repository/interfaces"
	addresssvc "be/internal/services/address"
)

// PlaceCounter is the Place feature API Location Delete needs.
type PlaceCounter interface {
	CountByLocationID(ctx context.Context, locationID string) (int64, error)
}

// Service implements Location CRUD.
type Service struct {
	locations interfaces.LocationRepository
	address   *addresssvc.Service
	places    PlaceCounter
}

func NewService(locations interfaces.LocationRepository, address *addresssvc.Service) *Service {
	return &Service{locations: locations, address: address}
}

// SetPlaceCounter wires Place after both services exist (avoids an import cycle).
func (s *Service) SetPlaceCounter(places PlaceCounter) {
	s.places = places
}

func (s *Service) List(ctx context.Context, q mapsdto.LocationListQuery) (*mapsdto.LocationListResponse, error) {
	page, limit := mapsdto.NormalizeList(q.Page, q.Limit)
	rows, total, err := s.locations.List(ctx, interfaces.LocationListFilter{
		Q:               strings.TrimSpace(q.Q),
		CountryCode:     strings.ToUpper(strings.TrimSpace(q.CountryCode)),
		AdminDivisionID: strings.TrimSpace(q.AdminDivisionID),
		Page:            page,
		Limit:           limit,
	})
	if err != nil {
		return nil, err
	}
	items := make([]mapsdto.LocationRecord, 0, len(rows))
	for _, row := range rows {
		items = append(items, toLocationRecord(row))
	}
	return &mapsdto.LocationListResponse{
		Items:   items,
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasMore: int64(page*limit) < total,
	}, nil
}

func (s *Service) Get(ctx context.Context, id string) (*mapsdto.LocationRecord, error) {
	row, err := s.loadRecord(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, apperrors.ErrNotFound
	}
	rec := toLocationRecord(*row)
	return &rec, nil
}

func (s *Service) Create(ctx context.Context, userID string, req mapsdto.CreateLocationRequest) (*mapsdto.LocationRecord, error) {
	loc, err := s.buildCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	actor := utils.StringPtr(userID)
	now := time.Now().UTC()
	loc.CreatedBy = actor
	loc.UpdatedBy = actor
	loc.CreatedAt = now
	loc.UpdatedAt = now
	if err := s.locations.Create(ctx, loc); err != nil {
		return nil, err
	}
	return s.Get(ctx, loc.ID)
}

func (s *Service) Patch(ctx context.Context, id, userID string, req mapsdto.PatchLocationRequest) (*mapsdto.LocationRecord, error) {
	loc, err := s.locations.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if loc == nil {
		return nil, apperrors.ErrNotFound
	}
	if err := s.applyPatch(ctx, loc, req); err != nil {
		return nil, err
	}
	loc.UpdatedBy = utils.StringPtr(userID)
	loc.UpdatedAt = time.Now().UTC()
	if err := s.locations.Update(ctx, loc); err != nil {
		return nil, err
	}
	return s.Get(ctx, loc.ID)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	loc, err := s.locations.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if loc == nil {
		return apperrors.ErrNotFound
	}
	if s.places == nil {
		return fmt.Errorf("place service is not wired")
	}
	count, err := s.places.CountByLocationID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("%w: location has places", apperrors.ErrConflict)
	}
	return s.locations.Delete(ctx, id)
}

func (s *Service) buildCreate(ctx context.Context, req mapsdto.CreateLocationRequest) (*locationmodel.Location, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 200 {
		return nil, fmt.Errorf("%w: name is required (max 200)", apperrors.ErrBadRequest)
	}
	if err := validateOptionalCoords(req.Lat, req.Lng); err != nil {
		return nil, err
	}
	countryID, err := s.resolveCountryID(ctx, req.CountryCode)
	if err != nil {
		return nil, err
	}
	adminID, err := s.resolveAdminID(ctx, countryID, req.AdminDivisionID)
	if err != nil {
		return nil, err
	}
	key, err := s.uniqueLocationKey(ctx, "", req.LocationKey)
	if err != nil {
		return nil, err
	}
	return &locationmodel.Location{
		Name:            name,
		LocationKey:     key,
		CountryID:       countryID,
		AdminDivisionID: adminID,
		Street:          strings.TrimSpace(req.Street),
		PostalCode:      strings.TrimSpace(req.PostalCode),
		Formatted:       strings.TrimSpace(req.Formatted),
		Lat:             req.Lat,
		Lng:             req.Lng,
	}, nil
}

func (s *Service) applyPatch(ctx context.Context, loc *locationmodel.Location, req mapsdto.PatchLocationRequest) error {
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" || len(name) > 200 {
			return fmt.Errorf("%w: name is required (max 200)", apperrors.ErrBadRequest)
		}
		loc.Name = name
	}
	if req.Street != nil {
		loc.Street = strings.TrimSpace(*req.Street)
	}
	if req.PostalCode != nil {
		loc.PostalCode = strings.TrimSpace(*req.PostalCode)
	}
	if req.Formatted != nil {
		loc.Formatted = strings.TrimSpace(*req.Formatted)
	}
	if req.Lat != nil || req.Lng != nil {
		lat, lng := loc.Lat, loc.Lng
		if req.Lat != nil {
			lat = req.Lat
		}
		if req.Lng != nil {
			lng = req.Lng
		}
		if err := validateOptionalCoords(lat, lng); err != nil {
			return err
		}
		loc.Lat = lat
		loc.Lng = lng
	}
	if req.CountryCode != nil {
		countryID, err := s.resolveCountryID(ctx, req.CountryCode)
		if err != nil {
			return err
		}
		loc.CountryID = countryID
	}
	if req.AdminDivisionID != nil {
		adminID, err := s.resolveAdminID(ctx, loc.CountryID, req.AdminDivisionID)
		if err != nil {
			return err
		}
		loc.AdminDivisionID = adminID
	}
	if req.LocationKey != nil {
		key, err := s.uniqueLocationKey(ctx, loc.ID, req.LocationKey)
		if err != nil {
			return err
		}
		loc.LocationKey = key
	}
	return nil
}

func (s *Service) resolveCountryID(ctx context.Context, code *string) (*string, error) {
	if code == nil {
		return nil, nil
	}
	trimmed := strings.ToUpper(strings.TrimSpace(*code))
	if trimmed == "" {
		return nil, nil
	}
	country, err := s.address.CountryByCode(ctx, trimmed)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, fmt.Errorf("%w: unknown countryCode", apperrors.ErrBadRequest)
	}
	return &country.ID, nil
}

func (s *Service) resolveAdminID(ctx context.Context, countryID *string, adminID *string) (*string, error) {
	if adminID == nil {
		return nil, nil
	}
	id := strings.TrimSpace(*adminID)
	if id == "" {
		return nil, nil
	}
	div, err := s.address.DivisionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if div == nil {
		return nil, fmt.Errorf("%w: unknown adminDivisionId", apperrors.ErrBadRequest)
	}
	if countryID != nil && div.CountryID != *countryID {
		return nil, fmt.Errorf("%w: adminDivisionId does not belong to country", apperrors.ErrBadRequest)
	}
	return &id, nil
}

func (s *Service) uniqueLocationKey(ctx context.Context, currentID string, key *string) (*string, error) {
	if key == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*key)
	if trimmed == "" {
		return nil, nil
	}
	existing, err := s.locations.FindByLocationKey(ctx, trimmed)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.ID != currentID {
		return nil, fmt.Errorf("%w: locationKey already exists", apperrors.ErrConflict)
	}
	return &trimmed, nil
}

func (s *Service) loadRecord(ctx context.Context, id string) (*interfaces.LocationListRow, error) {
	loc, err := s.locations.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		return nil, nil
	}
	row := interfaces.LocationListRow{Location: *loc}
	if loc.CountryID != nil {
		country, err := s.address.CountryByID(ctx, *loc.CountryID)
		if err != nil {
			return nil, err
		}
		if country != nil {
			code := country.Code
			row.CountryCode = &code
		}
	}
	if loc.AdminDivisionID != nil {
		div, err := s.address.DivisionByID(ctx, *loc.AdminDivisionID)
		if err != nil {
			return nil, err
		}
		if div != nil {
			path := div.Path
			row.AdminPath = &path
		}
	}
	return &row, nil
}

func toLocationRecord(row interfaces.LocationListRow) mapsdto.LocationRecord {
	loc := row.Location
	return mapsdto.LocationRecord{
		LocationSummary: mapsdto.LocationSummary{
			ID:              loc.ID,
			Name:            loc.Name,
			LocationKey:     loc.LocationKey,
			CountryCode:     row.CountryCode,
			AdminDivisionID: loc.AdminDivisionID,
			AdminPath:       row.AdminPath,
			Street:          loc.Street,
			PostalCode:      loc.PostalCode,
			Formatted:       loc.Formatted,
			Lat:             loc.Lat,
			Lng:             loc.Lng,
		},
		CreatedAt: utils.FormatRFC3339(loc.CreatedAt),
		UpdatedAt: utils.FormatRFC3339(loc.UpdatedAt),
	}
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

