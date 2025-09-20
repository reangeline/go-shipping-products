package order

import "testing"

func TestNewOrder_OK(t *testing.T) {
	o, err := NewOrder(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if o.Quantity != 10 {
		t.Fatalf("qty mismatch: got=%d want=%d", o.Quantity, 10)
	}
}

func TestNewOrder_Invalid(t *testing.T) {
	if _, err := NewOrder(0); err == nil {
		t.Fatalf("expected error for qty=0")
	}
	if _, err := NewOrder(-5); err == nil {
		t.Fatalf("expected error for negative qty")
	}
}
