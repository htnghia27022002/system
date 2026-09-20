package address_test

import (
	"context"
	"errors"
	"testing"

	apperrors "be/common/errors"
	addressdto "be/internal/dto/address"
	addresssvc "be/internal/services/address"
	"be/test/testutil"
)

func TestListCountriesReturnsActiveCatalog(t *testing.T) {
	t.Parallel()
	svc := addresssvc.NewService(testutil.NewMemoryCountryRepo(), testutil.NewMemoryDivisionRepo())
	out, err := svc.ListCountries(context.Background())
	if err != nil {
		t.Fatalf("countries: %v", err)
	}
	if len(out.Items) != 1 || out.Items[0].Code != "VN" {
		t.Fatalf("items: %#v", out.Items)
	}
}

func TestListDivisionsRequiresCountry(t *testing.T) {
	t.Parallel()
	svc := addresssvc.NewService(testutil.NewMemoryCountryRepo(), testutil.NewMemoryDivisionRepo())
	_, err := svc.ListDivisions(context.Background(), addressdto.DivisionListQuery{})
	if err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected bad request, got %v", err)
	}
}
