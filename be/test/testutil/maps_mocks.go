package testutil

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	categorymodel "be/internal/models/category"
	dsmodel "be/internal/models/datasource"
	geomodel "be/internal/models/geo"
	locationmodel "be/internal/models/location"
	newsmodel "be/internal/models/news"
	placemodel "be/internal/models/place"
	"be/internal/repository/interfaces"
)

type MemoryCountryRepo struct {
	mu   sync.Mutex
	ByID map[string]*geomodel.Country
}

func NewMemoryCountryRepo() *MemoryCountryRepo {
	vn := &geomodel.Country{ID: "country-vn", Code: "VN", Code3: "VNM", Name: "Vietnam", NameLocal: "Viet Nam", IsActive: true}
	return &MemoryCountryRepo{ByID: map[string]*geomodel.Country{vn.ID: vn}}
}

func (m *MemoryCountryRepo) ListActive(_ context.Context) ([]geomodel.Country, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]geomodel.Country, 0, len(m.ByID))
	for _, row := range m.ByID {
		if row.IsActive {
			out = append(out, *row)
		}
	}
	return out, nil
}

func (m *MemoryCountryRepo) GetByCode(_ context.Context, code string) (*geomodel.Country, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if strings.EqualFold(row.Code, code) {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *MemoryCountryRepo) GetByID(_ context.Context, id string) (*geomodel.Country, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByID[id]
	if row == nil {
		return nil, nil
	}
	cp := *row
	return &cp, nil
}

type MemoryDivisionRepo struct {
	mu   sync.Mutex
	ByID map[string]*geomodel.Division
}

func NewMemoryDivisionRepo() *MemoryDivisionRepo {
	return &MemoryDivisionRepo{ByID: map[string]*geomodel.Division{}}
}

func (m *MemoryDivisionRepo) ListActive(_ context.Context, filter interfaces.DivisionListFilter) ([]geomodel.Division, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []geomodel.Division{}
	for _, row := range m.ByID {
		if !row.IsActive {
			continue
		}
		if filter.CountryID != "" && row.CountryID != filter.CountryID {
			continue
		}
		out = append(out, *row)
	}
	return out, nil
}

func (m *MemoryDivisionRepo) GetByID(_ context.Context, id string) (*geomodel.Division, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByID[id]
	if row == nil {
		return nil, nil
	}
	cp := *row
	return &cp, nil
}

func (m *MemoryDivisionRepo) GetByCountryAndCode(_ context.Context, countryID, code string) (*geomodel.Division, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if row.CountryID == countryID && row.Code == code {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

type MemoryLocationRepo struct {
	mu   sync.Mutex
	ByID map[string]*locationmodel.Location
}

func NewMemoryLocationRepo() *MemoryLocationRepo {
	return &MemoryLocationRepo{ByID: map[string]*locationmodel.Location{}}
}

func (m *MemoryLocationRepo) GetByID(_ context.Context, id string) (*locationmodel.Location, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByID[id]
	if row == nil {
		return nil, nil
	}
	cp := *row
	return &cp, nil
}

func (m *MemoryLocationRepo) FindByLocationKey(_ context.Context, key string) (*locationmodel.Location, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if row.LocationKey != nil && *row.LocationKey == key {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *MemoryLocationRepo) FindByCountryAdminStreet(_ context.Context, countryID, adminDivisionID, street string) (*locationmodel.Location, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if row.CountryID != nil && *row.CountryID == countryID &&
			row.AdminDivisionID != nil && *row.AdminDivisionID == adminDivisionID &&
			strings.EqualFold(strings.TrimSpace(row.Street), strings.TrimSpace(street)) {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *MemoryLocationRepo) FindByFormatted(_ context.Context, formatted string) (*locationmodel.Location, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if strings.EqualFold(strings.TrimSpace(row.Formatted), strings.TrimSpace(formatted)) && strings.TrimSpace(row.Formatted) != "" {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *MemoryLocationRepo) FindByRoundedCoordsAndFormatted(_ context.Context, lat, lng float64, formatted string) (*locationmodel.Location, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if row.Lat == nil || row.Lng == nil {
			continue
		}
		if *row.Lat == lat && *row.Lng == lng && strings.EqualFold(strings.TrimSpace(row.Formatted), strings.TrimSpace(formatted)) {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *MemoryLocationRepo) Create(_ context.Context, loc *locationmodel.Location) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if loc.ID == "" {
		loc.ID = uuid.NewString()
	}
	cp := *loc
	m.ByID[loc.ID] = &cp
	return nil
}

func (m *MemoryLocationRepo) Update(_ context.Context, loc *locationmodel.Location) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *loc
	m.ByID[loc.ID] = &cp
	return nil
}

func (m *MemoryLocationRepo) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ByID, id)
	return nil
}

func (m *MemoryLocationRepo) List(_ context.Context, filter interfaces.LocationListFilter) ([]interfaces.LocationListRow, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]interfaces.LocationListRow, 0, len(m.ByID))
	q := strings.ToLower(strings.TrimSpace(filter.Q))
	for _, row := range m.ByID {
		if q != "" &&
			!strings.Contains(strings.ToLower(row.Name), q) &&
			!strings.Contains(strings.ToLower(row.Street), q) &&
			!strings.Contains(strings.ToLower(row.Formatted), q) {
			continue
		}
		out = append(out, interfaces.LocationListRow{Location: *row})
	}
	return out, int64(len(out)), nil
}

type MemoryCategoryRepo struct {
	mu     sync.Mutex
	ByKey  map[string]*categorymodel.Category
	Inserts int
}

func NewMemoryCategoryRepo() *MemoryCategoryRepo {
	now := time.Now().UTC()
	keys := []struct{ key, name string }{
		{categorymodel.KeyRoomRental, "Room rental"},
		{categorymodel.KeyRestaurant, "Restaurant"},
		{categorymodel.KeyHotel, "Hotel"},
		{categorymodel.KeyEatery, "Eatery"},
		{categorymodel.KeyUncategorized, "Uncategorized"},
	}
	byKey := map[string]*categorymodel.Category{}
	for i, item := range keys {
		byKey[item.key] = &categorymodel.Category{
			ID:        "cat-" + item.key,
			Key:       item.key,
			Name:      item.name,
			SortOrder: i + 1,
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return &MemoryCategoryRepo{ByKey: byKey}
}

func (m *MemoryCategoryRepo) ListActive(_ context.Context) ([]categorymodel.Category, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]categorymodel.Category, 0, len(m.ByKey))
	for _, row := range m.ByKey {
		if row.IsActive {
			out = append(out, *row)
		}
	}
	return out, nil
}

func (m *MemoryCategoryRepo) GetByKey(_ context.Context, key string) (*categorymodel.Category, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByKey[key]
	if row == nil {
		return nil, nil
	}
	cp := *row
	return &cp, nil
}

func (m *MemoryCategoryRepo) GetByID(_ context.Context, id string) (*categorymodel.Category, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByKey {
		if row.ID == id {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

type MemoryPlaceRepo struct {
	mu   sync.Mutex
	ByID map[string]*placemodel.Place
}

func NewMemoryPlaceRepo() *MemoryPlaceRepo {
	return &MemoryPlaceRepo{ByID: map[string]*placemodel.Place{}}
}

func (m *MemoryPlaceRepo) GetByID(_ context.Context, id string) (*placemodel.Place, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByID[id]
	if row == nil {
		return nil, nil
	}
	cp := *row
	return &cp, nil
}

func (m *MemoryPlaceRepo) FindByLocationAndKey(_ context.Context, locationID, placeKey string) (*placemodel.Place, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if row.LocationID == locationID && row.PlaceKey == placeKey {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *MemoryPlaceRepo) Create(_ context.Context, p *placemodel.Place) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	cp := *p
	if p.Details != nil {
		cp.Details = copyMap(p.Details)
	}
	m.ByID[p.ID] = &cp
	return nil
}

func (m *MemoryPlaceRepo) Update(_ context.Context, p *placemodel.Place) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *p
	if p.Details != nil {
		cp.Details = copyMap(p.Details)
	}
	m.ByID[p.ID] = &cp
	return nil
}

func (m *MemoryPlaceRepo) ListPinnable(_ context.Context, filter interfaces.PlaceListFilter) ([]interfaces.PlaceListRow, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	q := strings.ToLower(strings.TrimSpace(filter.Q))
	out := []interfaces.PlaceListRow{}
	for _, row := range m.ByID {
		if row.Status != placemodel.StatusActive {
			continue
		}
		if !placemodel.ValidCoords(row.Lat, row.Lng) {
			continue
		}
		if q != "" &&
			!strings.Contains(strings.ToLower(row.Name), q) &&
			!strings.Contains(strings.ToLower(row.Unit), q) {
			continue
		}
		out = append(out, interfaces.PlaceListRow{
			Place:        *row,
			LocationID:   row.LocationID,
			LocationName: row.LocationID,
		})
	}
	return out, int64(len(out)), nil
}

func (m *MemoryPlaceRepo) ListAdmin(_ context.Context, filter interfaces.PlaceAdminListFilter) ([]interfaces.PlaceListRow, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []interfaces.PlaceListRow{}
	for _, row := range m.ByID {
		if filter.LocationID != "" && row.LocationID != filter.LocationID {
			continue
		}
		if filter.Status != nil && row.Status != *filter.Status {
			continue
		}
		out = append(out, interfaces.PlaceListRow{
			Place:      *row,
			LocationID: row.LocationID,
		})
	}
	return out, int64(len(out)), nil
}

func (m *MemoryPlaceRepo) CountByLocationID(_ context.Context, locationID string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var n int64
	for _, row := range m.ByID {
		if row.LocationID == locationID {
			n++
		}
	}
	return n, nil
}

func (m *MemoryPlaceRepo) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ByID, id)
	return nil
}

func (m *MemoryPlaceRepo) GetDetail(_ context.Context, id string) (*interfaces.PlaceDetailRow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByID[id]
	if row == nil {
		return nil, nil
	}
	return &interfaces.PlaceDetailRow{
		Place:       *row,
		CategoryKey: "eatery",
		Location:    interfaces.PlaceLocationRow{ID: row.LocationID},
		News:        []newsmodel.News{},
	}, nil
}

type MemoryNewsRepo struct {
	mu   sync.Mutex
	ByID map[string]*newsmodel.News
}

func NewMemoryNewsRepo() *MemoryNewsRepo {
	return &MemoryNewsRepo{ByID: map[string]*newsmodel.News{}}
}

func (m *MemoryNewsRepo) GetByOriginalURL(_ context.Context, originalURL string) (*newsmodel.News, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if row.OriginalURL == originalURL {
			cp := *row
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *MemoryNewsRepo) ListByPlaceID(_ context.Context, placeID string) ([]newsmodel.News, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []newsmodel.News{}
	for _, row := range m.ByID {
		if row.PlaceID == placeID {
			out = append(out, *row)
		}
	}
	return out, nil
}

func (m *MemoryNewsRepo) Create(_ context.Context, item *newsmodel.News) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	cp := *item
	m.ByID[item.ID] = &cp
	return nil
}

func (m *MemoryNewsRepo) Update(_ context.Context, item *newsmodel.News) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *item
	m.ByID[item.ID] = &cp
	return nil
}

type MemoryDataSourceRepo struct {
	mu   sync.Mutex
	ByID map[string]*dsmodel.Source
}

func NewMemoryDataSourceRepo() *MemoryDataSourceRepo {
	return &MemoryDataSourceRepo{ByID: map[string]*dsmodel.Source{}}
}

func (m *MemoryDataSourceRepo) GetByID(_ context.Context, id string) (*dsmodel.Source, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByID[id]
	if row == nil {
		return nil, nil
	}
	cp := *row
	return &cp, nil
}

func (m *MemoryDataSourceRepo) List(_ context.Context, _, _ int) ([]dsmodel.Source, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]dsmodel.Source, 0, len(m.ByID))
	for _, row := range m.ByID {
		out = append(out, *row)
	}
	return out, int64(len(out)), nil
}

func (m *MemoryDataSourceRepo) ListEnabled(_ context.Context) ([]dsmodel.Source, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []dsmodel.Source{}
	for _, row := range m.ByID {
		if row.Enabled {
			out = append(out, *row)
		}
	}
	return out, nil
}

func (m *MemoryDataSourceRepo) Create(_ context.Context, src *dsmodel.Source) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if src.ID == "" {
		src.ID = uuid.NewString()
	}
	cp := *src
	m.ByID[src.ID] = &cp
	return nil
}

func (m *MemoryDataSourceRepo) Update(_ context.Context, src *dsmodel.Source) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *src
	m.ByID[src.ID] = &cp
	return nil
}

func (m *MemoryDataSourceRepo) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ByID, id)
	return nil
}

type MemoryIngestRunRepo struct {
	mu      sync.Mutex
	ByID    map[string]*dsmodel.IngestRun
	Sources map[string]*dsmodel.IngestRunSource
}

func NewMemoryIngestRunRepo() *MemoryIngestRunRepo {
	return &MemoryIngestRunRepo{
		ByID:    map[string]*dsmodel.IngestRun{},
		Sources: map[string]*dsmodel.IngestRunSource{},
	}
}

func (m *MemoryIngestRunRepo) HasActive(_ context.Context) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, row := range m.ByID {
		if row.Status == dsmodel.RunQueued || row.Status == dsmodel.RunRunning {
			return true, nil
		}
	}
	return false, nil
}

func (m *MemoryIngestRunRepo) Insert(_ context.Context, run *dsmodel.IngestRun) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if run.ID == "" {
		run.ID = uuid.NewString()
	}
	cp := *run
	m.ByID[run.ID] = &cp
	return nil
}

func (m *MemoryIngestRunRepo) GetByID(_ context.Context, id string) (*dsmodel.IngestRun, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByID[id]
	if row == nil {
		return nil, nil
	}
	cp := *row
	return &cp, nil
}

func (m *MemoryIngestRunRepo) Latest(_ context.Context) (*dsmodel.IngestRun, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var latest *dsmodel.IngestRun
	for _, row := range m.ByID {
		if latest == nil || row.CreatedAt.After(latest.CreatedAt) {
			cp := *row
			latest = &cp
		}
	}
	return latest, nil
}

func (m *MemoryIngestRunRepo) ClaimQueued(_ context.Context, id string) (*dsmodel.IngestRun, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	row := m.ByID[id]
	if row == nil || row.Status != dsmodel.RunQueued {
		return nil, nil
	}
	row.Status = dsmodel.RunRunning
	now := time.Now().UTC()
	row.StartedAt = &now
	cp := *row
	return &cp, nil
}

func (m *MemoryIngestRunRepo) ListQueued(_ context.Context, _ int) ([]dsmodel.IngestRun, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []dsmodel.IngestRun{}
	for _, row := range m.ByID {
		if row.Status == dsmodel.RunQueued {
			out = append(out, *row)
		}
	}
	return out, nil
}

func (m *MemoryIngestRunRepo) Update(_ context.Context, run *dsmodel.IngestRun) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *run
	m.ByID[run.ID] = &cp
	return nil
}

func (m *MemoryIngestRunRepo) InsertSource(_ context.Context, row *dsmodel.IngestRunSource) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if row.ID == "" {
		row.ID = uuid.NewString()
	}
	cp := *row
	m.Sources[row.ID] = &cp
	return nil
}

func (m *MemoryIngestRunRepo) UpdateSource(_ context.Context, row *dsmodel.IngestRunSource) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *row
	m.Sources[row.ID] = &cp
	return nil
}

func (m *MemoryIngestRunRepo) ListSources(_ context.Context, runID string) ([]dsmodel.IngestRunSource, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []dsmodel.IngestRunSource{}
	for _, row := range m.Sources {
		if row.RunID == runID {
			out = append(out, *row)
		}
	}
	return out, nil
}

func copyMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

var (
	_ interfaces.CountryRepository     = (*MemoryCountryRepo)(nil)
	_ interfaces.DivisionRepository    = (*MemoryDivisionRepo)(nil)
	_ interfaces.LocationRepository    = (*MemoryLocationRepo)(nil)
	_ interfaces.CategoryRepository    = (*MemoryCategoryRepo)(nil)
	_ interfaces.PlaceRepository       = (*MemoryPlaceRepo)(nil)
	_ interfaces.NewsRepository        = (*MemoryNewsRepo)(nil)
	_ interfaces.DataSourceRepository  = (*MemoryDataSourceRepo)(nil)
	_ interfaces.IngestRunRepository   = (*MemoryIngestRunRepo)(nil)
)
