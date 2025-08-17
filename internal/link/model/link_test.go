package model

import "testing"

func TestNewLinkGeneratesUniqueHash(t *testing.T) {
	generated := make(map[string]bool)
	for i := 0; i < 100; i++ {
		l, err := NewLink("https://example.com", func(hash string) bool {
			return generated[hash]
		})
		if err != nil {
			t.Fatalf("NewLink error: %v", err)
		}
		if l.Hash == "" || len(l.Hash) != hashLength {
			t.Fatalf("unexpected hash: %s", l.Hash)
		}
		if generated[l.Hash] {
			t.Fatalf("hash collision: %s", l.Hash)
		}
		generated[l.Hash] = true
	}
}

func TestNewLinkCollisionFail(t *testing.T) {
	_, err := NewLink("https://example.com", func(string) bool { return true })
	if err == nil {
		t.Fatalf("expected error due to collisions")
	}
}