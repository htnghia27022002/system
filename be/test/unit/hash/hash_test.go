package hash_test

import (
	"testing"

	"be/pkg/hash"
)

func TestSHA256HexStable(t *testing.T) {
	got := hash.SHA256Hex("refresh-token")
	if len(got) != 64 {
		t.Fatalf("digest length = %d, want 64", len(got))
	}
	if hash.SHA256Hex("refresh-token") != got {
		t.Fatal("digest must be deterministic")
	}
	if hash.SHA256Hex("") == got {
		t.Fatal("empty input must not match non-empty digest")
	}

	empty := hash.SHA256Hex("")
	const emptySHA256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if empty != emptySHA256 {
		t.Fatalf("empty digest = %s", empty)
	}
}
