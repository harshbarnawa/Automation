package auth

import "testing"

func TestBcryptHasherHashAndCompare(t *testing.T) {
	hasher := NewBcryptHasher(4)

	hash, err := hasher.Hash("supersecret")
	if err != nil {
		t.Fatalf("expected hash, got error: %v", err)
	}
	if hash == "supersecret" {
		t.Fatal("expected hash to differ from plaintext")
	}

	if err := hasher.Compare(hash, "supersecret"); err != nil {
		t.Fatalf("expected password to match, got %v", err)
	}

	if err := hasher.Compare(hash, "wrong"); err == nil {
		t.Fatal("expected mismatch error for wrong password")
	}
}

func TestNewBcryptHasherDefaultsCost(t *testing.T) {
	hasher := NewBcryptHasher(0)
	if _, err := hasher.Hash("supersecret"); err != nil {
		t.Fatalf("expected default-cost hasher to work, got %v", err)
	}
}
