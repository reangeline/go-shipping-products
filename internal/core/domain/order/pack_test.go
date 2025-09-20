package order

import "testing"

func TestNewPack_OK(t *testing.T) {
	p, err := NewPack(250)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Size != 250 {
		t.Fatalf("size mismatch: got=%d want=%d", p.Size, 250)
	}
}

func TestNewPack_Invalid(t *testing.T) {
	if _, err := NewPack(0); err == nil {
		t.Fatalf("expected error for size=0")
	}
	if _, err := NewPack(-10); err == nil {
		t.Fatalf("expected error for negative size")
	}
}
