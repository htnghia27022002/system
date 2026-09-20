package maps_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"be/internal/config"
	mapsdto "be/internal/dto/maps"
	locationmodel "be/internal/models/location"
	placemodel "be/internal/models/place"
	addresssvc "be/internal/services/address"
	locationsvc "be/internal/services/location"
	mapssvc "be/internal/services/maps"
	mapsproviders "be/internal/services/maps/providers"
	placesvc "be/internal/services/place"
	"be/test/testutil"
)

func TestSearchMergesSavedAndGeocodeHits(t *testing.T) {
	t.Parallel()
	lat, lng := 16.07, 108.22
	places := testutil.NewMemoryPlaceRepo()
	_ = places.Create(context.Background(), &placemodel.Place{
		ID:         "p1",
		LocationID: "loc-1",
		Name:       "Com Tam Demo",
		Status:     placemodel.StatusActive,
		Lat:        &lat,
		Lng:        &lng,
	})
	locations := testutil.NewMemoryLocationRepo()
	_ = locations.Create(context.Background(), &locationmodel.Location{
		ID:        "l1",
		Name:      "Cho Han Market",
		Formatted: "Hai Chau, Da Nang",
		Lat:       &lat,
		Lng:       &lng,
	})

	nominatim := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(strings.ToLower(r.URL.Query().Get("q")), "dragon") {
			_, _ = w.Write([]byte("[]"))
			return
		}
		_, _ = w.Write([]byte(`[{"place_id":12,"name":"Dragon Bridge","display_name":"Dragon Bridge, Da Nang","lat":"16.061","lon":"108.227"}]`))
	}))
	t.Cleanup(nominatim.Close)

	svc := mapssvc.NewService(config.Config{
		MapNominatimURL:       nominatim.URL,
		MapNominatimUserAgent: "system-maps-test",
	}, newTestPlaceService(places), newTestLocationService(locations))

	placeHits, err := svc.Search(context.Background(), mapsdto.SearchQuery{Q: "tam"})
	if err != nil {
		t.Fatalf("place search: %v", err)
	}
	if len(placeHits.Items) != 1 || placeHits.Items[0].Kind != "place" || placeHits.Items[0].PlaceID != "p1" {
		t.Fatalf("place: %#v", placeHits.Items)
	}

	locHits, err := svc.Search(context.Background(), mapsdto.SearchQuery{Q: "han"})
	if err != nil {
		t.Fatalf("location search: %v", err)
	}
	if len(locHits.Items) != 1 || locHits.Items[0].Kind != "location" || locHits.Items[0].LocationID != "l1" {
		t.Fatalf("location: %#v", locHits.Items)
	}

	geoHits, err := svc.Search(context.Background(), mapsdto.SearchQuery{Q: "dragon"})
	if err != nil {
		t.Fatalf("geocode search: %v", err)
	}
	if len(geoHits.Items) != 1 || geoHits.Items[0].Kind != "geocode" {
		t.Fatalf("geocode: %#v", geoHits.Items)
	}
}

func TestSearchIgnoresShortQuery(t *testing.T) {
	t.Parallel()
	svc := mapssvc.NewService(config.Config{}, newTestPlaceService(testutil.NewMemoryPlaceRepo()), newTestLocationService(testutil.NewMemoryLocationRepo()))
	out, err := svc.Search(context.Background(), mapsdto.SearchQuery{Q: "a"})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(out.Items) != 0 {
		t.Fatalf("want empty, got %#v", out.Items)
	}
}

func TestOSMProviderMapsNominatimJSON(t *testing.T) {
	t.Parallel()
	provider := mapsproviders.NewOSMProvider(config.Config{
		MapNominatimURL:       "https://nominatim.test",
		MapNominatimUserAgent: "system-maps-test",
	}, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/search" {
			t.Fatalf("path: %s", req.URL.Path)
		}
		if req.Header.Get("User-Agent") != "system-maps-test" {
			t.Fatalf("ua: %s", req.Header.Get("User-Agent"))
		}
		if !strings.Contains(req.URL.RawQuery, "countrycodes=vn") {
			t.Fatalf("query: %s", req.URL.RawQuery)
		}
		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(strings.NewReader(`[
				{"place_id":12,"name":"Dragon Bridge","display_name":"Dragon Bridge, Da Nang","lat":"16.061","lon":"108.227"},
				{"place_id":13,"name":"Broken","lat":"x","lon":"y"}
			]`)),
			Header: make(http.Header),
		}, nil
	}))

	hits, err := provider.Search(context.Background(), "Dragon Bridge", "VN")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(hits) != 1 || hits[0].Title != "Dragon Bridge" || hits[0].Lat != 16.061 {
		t.Fatalf("hits: %#v", hits)
	}
}

func TestGeocodeRegistryDefaultsToOSM(t *testing.T) {
	t.Parallel()
	if mapssvc.New(config.Config{}, nil).ID() != "osm" {
		t.Fatal("default provider")
	}
	if mapssvc.New(config.Config{MapGeocodeProvider: "google"}, nil).ID() != "google" {
		t.Fatal("google provider")
	}
}

func newTestPlaceService(places *testutil.MemoryPlaceRepo) *placesvc.Service {
	address := addresssvc.NewService(testutil.NewMemoryCountryRepo(), testutil.NewMemoryDivisionRepo())
	locations := locationsvc.NewService(testutil.NewMemoryLocationRepo(), address)
	return placesvc.NewService(places, locations, testutil.NewMemoryCategoryRepo(), address)
}

func newTestLocationService(locations *testutil.MemoryLocationRepo) *locationsvc.Service {
	return locationsvc.NewService(
		locations,
		addresssvc.NewService(testutil.NewMemoryCountryRepo(), testutil.NewMemoryDivisionRepo()),
	)
}
