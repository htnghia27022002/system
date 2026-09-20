package utils_test

import (
	"testing"

	"be/common/utils"
)

func TestStringPtr(t *testing.T) {
	if utils.StringPtr("   ") != nil {
		t.Fatal("expected nil for blank string")
	}
	got := utils.StringPtr("  hello  ")
	if got == nil || *got != "hello" {
		t.Fatalf("StringPtr = %v, want hello", got)
	}
}
