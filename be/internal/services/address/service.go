package address

import (
	"context"
	"fmt"
	"strings"

	apperrors "be/common/errors"
	addressdto "be/internal/dto/address"
	geomodel "be/internal/models/geo"
	"be/internal/repository/interfaces"
)

// Service reads the shared country and administrative-division catalogs.
type Service struct {
	countries interfaces.CountryRepository
	divisions interfaces.DivisionRepository
}

func NewService(
	countries interfaces.CountryRepository,
	divisions interfaces.DivisionRepository,
) *Service {
	return &Service{countries: countries, divisions: divisions}
}

func (s *Service) ListCountries(ctx context.Context) (*addressdto.CountryListResponse, error) {
	rows, err := s.countries.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]addressdto.CountryItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, addressdto.CountryItem{
			ID:        row.ID,
			Code:      row.Code,
			Code3:     row.Code3,
			Name:      row.Name,
			NameLocal: row.NameLocal,
		})
	}
	return &addressdto.CountryListResponse{Items: items}, nil
}

func (s *Service) ListDivisions(ctx context.Context, q addressdto.DivisionListQuery) (*addressdto.DivisionListResponse, error) {
	code := strings.ToUpper(strings.TrimSpace(q.CountryCode))
	if code == "" {
		return nil, fmt.Errorf("%w: countryCode is required", apperrors.ErrBadRequest)
	}
	country, err := s.countries.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, fmt.Errorf("%w: unknown countryCode", apperrors.ErrBadRequest)
	}
	filter := interfaces.DivisionListFilter{
		CountryID: country.ID,
		Q:         strings.TrimSpace(q.Q),
	}
	parent := strings.TrimSpace(q.ParentID)
	filter.ParentID = &parent
	rows, err := s.divisions.ListActive(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]addressdto.AdminDivision, 0, len(rows))
	for _, row := range rows {
		items = append(items, addressdto.AdminDivision{
			ID:          row.ID,
			CountryCode: country.Code,
			ParentID:    row.ParentID,
			Level:       row.Level,
			Code:        row.Code,
			Name:        row.Name,
			NameEn:      row.NameEn,
			FullName:    row.FullName,
			Type:        row.Type,
			Path:        row.Path,
			Lat:         row.Lat,
			Lng:         row.Lng,
		})
	}
	return &addressdto.DivisionListResponse{Items: items}, nil
}

func (s *Service) CountryByCode(ctx context.Context, code string) (*geomodel.Country, error) {
	return s.countries.GetByCode(ctx, strings.ToUpper(strings.TrimSpace(code)))
}

func (s *Service) CountryByID(ctx context.Context, id string) (*geomodel.Country, error) {
	return s.countries.GetByID(ctx, strings.TrimSpace(id))
}

func (s *Service) DivisionByID(ctx context.Context, id string) (*geomodel.Division, error) {
	return s.divisions.GetByID(ctx, strings.TrimSpace(id))
}
