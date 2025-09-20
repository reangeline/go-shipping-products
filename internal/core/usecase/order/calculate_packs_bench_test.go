package order

import (
	"context"
	"testing"

	domain "github.com/reangeline/go-shipping-products/internal/core/domain/order"
	uc "github.com/reangeline/go-shipping-products/internal/core/ports/inbound/order"
)

type benchProvider struct{ sizes []int }

func (f *benchProvider) List() ([]int, error) {

	out := make([]int, len(f.sizes))
	copy(out, f.sizes)
	return out, nil
}

func BenchmarkCalculatePacks_Various(b *testing.B) {
	calc := domain.NewPackCalculator()
	prov := &benchProvider{sizes: []int{250, 500, 1000, 2000, 5000}}

	ucase, err := NewCalculatePacks(calc, prov)
	if err != nil {
		b.Fatalf("wire use case: %v", err)
	}

	ctx := context.Background()

	cases := []struct {
		name string
		qty  int
	}{
		{"Small", 12_001},
		{"Medium", 100_000},
		{"Large", 500_000},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			in := uc.CalculatePacksInput{Quantity: tc.qty}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := ucase.Execute(ctx, in); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCalculatePacks_WithOverride(b *testing.B) {
	calc := domain.NewPackCalculator()
	prov := &benchProvider{sizes: []int{1}}
	ucase, err := NewCalculatePacks(calc, prov)
	if err != nil {
		b.Fatalf("wire use case: %v", err)
	}

	ctx := context.Background()
	in := uc.CalculatePacksInput{
		Quantity:      12_001,
		PacksOverride: []int{250, 500, 1000, 2000, 5000},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ucase.Execute(ctx, in); err != nil {
			b.Fatal(err)
		}
	}
}
