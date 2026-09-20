package utils_test

import (
	"testing"
	"time"

	"be/common/utils"
)

func TestFormatRFC3339UsesUTC(t *testing.T) {
	loc := time.FixedZone("UTC+7", 7*3600)
	got := utils.FormatRFC3339(time.Date(2026, 9, 12, 22, 0, 0, 0, loc))
	want := "2026-09-12T15:00:00Z"
	if got != want {
		t.Fatalf("FormatRFC3339 = %q, want %q", got, want)
	}
}

func TestFormatRFC3339PtrNil(t *testing.T) {
	if utils.FormatRFC3339Ptr(nil) != nil {
		t.Fatal("expected nil for nil time")
	}
}

func TestFormatDate(t *testing.T) {
	if utils.FormatDate(nil) != nil {
		t.Fatal("expected nil for nil date")
	}
	day := time.Date(2026, 9, 12, 22, 0, 0, 0, time.UTC)
	got := utils.FormatDate(&day)
	if got == nil || *got != "2026-09-12" {
		t.Fatalf("FormatDate = %v, want 2026-09-12", got)
	}
}
