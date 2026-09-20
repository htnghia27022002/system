package maps_test

import (
	"context"
	"errors"
	"testing"

	apperrors "be/common/errors"
	mapsdto "be/internal/dto/maps"
	addresssvc "be/internal/services/address"
	locationsvc "be/internal/services/location"
	placesvc "be/internal/services/place"
	"be/test/testutil"
)

func TestLocationAndPlaceCRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	countries := testutil.NewMemoryCountryRepo()
	divisions := testutil.NewMemoryDivisionRepo()
	locations := testutil.NewMemoryLocationRepo()
	places := testutil.NewMemoryPlaceRepo()
	categories := testutil.NewMemoryCategoryRepo()

	address := addresssvc.NewService(countries, divisions)
	locSvc := locationsvc.NewService(locations, address)
	placeSvc := placesvc.NewService(places, locSvc, categories, address)
	locSvc.SetPlaceCounter(placeSvc)

	created, err := locSvc.Create(ctx, "user-1", mapsdto.CreateLocationRequest{
		Name:        "Chung cu Example",
		CountryCode: strPtr("VN"),
		Street:      "12 Nguyen Hue",
		Lat:         floatPtr(10.7769),
		Lng:         floatPtr(106.7009),
	})
	if err != nil {
		t.Fatalf("create location: %v", err)
	}
	if created.Name != "Chung cu Example" {
		t.Fatalf("location name: %s", created.Name)
	}

	place, err := placeSvc.Create(ctx, "user-1", mapsdto.CreatePlaceRequest{
		LocationID: created.ID,
		Name:       "Pho Shop",
		Category:   "eatery",
		Lat:        floatPtr(10.7769),
		Lng:        floatPtr(106.7009),
	})
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	if place.Status != "active" {
		t.Fatalf("expected active, got %s", place.Status)
	}

	if err := locSvc.Delete(ctx, created.ID); err == nil || !apperrors.IsConflict(err) {
		t.Fatalf("expected conflict deleting location with places, got %v", err)
	}

	hidden := "hidden"
	patched, err := placeSvc.Patch(ctx, place.ID, "user-1", mapsdto.PatchPlaceRequest{Status: &hidden})
	if err != nil {
		t.Fatalf("patch place: %v", err)
	}
	if patched.Status != "hidden" {
		t.Fatalf("expected hidden, got %s", patched.Status)
	}

	if err := placeSvc.Delete(ctx, place.ID); err != nil {
		t.Fatalf("delete place: %v", err)
	}
	if err := locSvc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete location: %v", err)
	}
}

func TestCreatePlaceRejectsUnknownLocation(t *testing.T) {
	t.Parallel()
	address := addresssvc.NewService(testutil.NewMemoryCountryRepo(), testutil.NewMemoryDivisionRepo())
	locSvc := locationsvc.NewService(testutil.NewMemoryLocationRepo(), address)
	placeSvc := placesvc.NewService(
		testutil.NewMemoryPlaceRepo(),
		locSvc,
		testutil.NewMemoryCategoryRepo(),
		address,
	)
	_, err := placeSvc.Create(context.Background(), "user-1", mapsdto.CreatePlaceRequest{
		LocationID: "missing",
		Name:       "Shop",
		Category:   "hotel",
	})
	if err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected bad request, got %v", err)
	}
}

func strPtr(v string) *string { return &v }

func floatPtr(v float64) *float64 { return &v }
