package dependency

import (
	addresssvc "be/internal/services/address"
	locationsvc "be/internal/services/location"
	mapssvc "be/internal/services/maps"
	ingestsvc "be/internal/services/maps/ingest"
	_ "be/internal/services/maps/providers"
	sourcesvc "be/internal/services/maps/source"
	placesvc "be/internal/services/place"
)

// MapsServices groups Maps HTTP services.
type MapsServices struct {
	Places    *placesvc.Service
	Locations *locationsvc.Service
	Sources   *sourcesvc.Service
	Ingest    *ingestsvc.Service
	Search    *mapssvc.Service
}

func NewMapsServices(infra *Infra, address *addresssvc.Service) *MapsServices {
	locationsRepo := newLocationRepository(infra.DB)
	placesRepo := newPlaceRepository(infra.DB)
	categories := newCategoryRepository(infra.DB)
	news := newNewsRepository(infra.DB)
	sources := newDataSourceRepository(infra.DB)
	runs := newIngestRunRepository(infra.DB)

	locations := locationsvc.NewService(locationsRepo, address)
	places := placesvc.NewService(placesRepo, locations, categories, address)
	locations.SetPlaceCounter(places)

	return &MapsServices{
		Places:    places,
		Locations: locations,
		Search:    mapssvc.NewService(infra.Config, places, locations),
		Sources:   sourcesvc.NewService(sources),
		Ingest: ingestsvc.NewService(
			sources,
			runs,
			locationsRepo,
			placesRepo,
			news,
			categories,
			newCountryRepository(infra.DB),
			newDivisionRepository(infra.DB),
			infra.Publisher,
			nil,
		),
	}
}
