package order

import (
	"context"
	"errors"
	"testing"

	domain "github.com/reangeline/go-shipping-products/internal/core/domain/order"
	uc "github.com/reangeline/go-shipping-products/internal/core/ports/inbound/order"
	"github.com/reangeline/go-shipping-products/internal/core/ports/outbound/packsizes"
)

type fakeProvider struct {
	sizes []int
	err   error
}

func (f *fakeProvider) List() ([]int, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]int(nil), f.sizes...), nil
}

func TestCalculatePacks_Execute(t *testing.T) {
	calc := domain.NewPackCalculator()
	var errBoom = errors.New("boom")

	tests := []struct {
		name      string
		input     uc.CalculatePacksInput
		provider  packsizes.Provider
		wantErr   error
		wantTotal int
		wantLeft  int
	}{
		{
			name:      "valid quantity from provider",
			input:     uc.CalculatePacksInput{Quantity: 12001},
			provider:  &fakeProvider{sizes: []int{250, 500, 1000, 2000, 5000}},
			wantErr:   nil,
			wantLeft:  249, // 12001 -> 5000+5000+1000+250 = 12250
			wantTotal: 12250,
		},
		{
			name:     "invalid quantity <= 0",
			input:    uc.CalculatePacksInput{Quantity: 0},
			provider: &fakeProvider{sizes: []int{250, 500}},
			wantErr:  ErrInvalidQuantity,
		},
		{
			name:      "override is used",
			input:     uc.CalculatePacksInput{Quantity: 10, PacksOverride: []int{3, 7}},
			provider:  &fakeProvider{sizes: []int{9999}},
			wantErr:   nil,
			wantTotal: 10, // 7+3
		},
		{
			name:     "override invalid (<=0)",
			input:    uc.CalculatePacksInput{Quantity: 5, PacksOverride: []int{-1, 2}},
			provider: &fakeProvider{sizes: []int{2}},
			wantErr:  ErrInvalidPackInOverride,
		},
		{
			name:     "provider empty",
			input:    uc.CalculatePacksInput{Quantity: 5},
			provider: &fakeProvider{sizes: []int{}},
			wantErr:  ErrNoPackSizes,
		},
		{
			name:     "provider error",
			input:    uc.CalculatePacksInput{Quantity: 5},
			provider: &fakeProvider{err: errBoom},
			wantErr:  errBoom,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ucase, err := NewCalculatePacks(calc, tt.provider)
			if err != nil {
				t.Fatalf("unexpected NewCalculatePacks error: %v", err)
			}

			out, err := ucase.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}

				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out.TotalItems != tt.wantTotal {
				t.Errorf("TotalItems got %d, want %d", out.TotalItems, tt.wantTotal)
			}
			if out.Leftover != tt.wantLeft {
				t.Errorf("Leftover got %d, want %d", out.Leftover, tt.wantLeft)
			}
		})
	}
}
