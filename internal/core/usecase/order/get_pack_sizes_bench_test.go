package order

import (
	"context"
	"testing"
)

type benchProvider2 struct{ sizes []int }

func (f *benchProvider2) List() ([]int, error) {
	out := make([]int, len(f.sizes))
	copy(out, f.sizes)
	return out, nil
}

func BenchmarkGetPackSizes(b *testing.B) {
	prov := &benchProvider2{sizes: []int{250, 500, 1000, 2000, 5000}}
	ucase, err := NewGetPackSizes(prov)
	if err != nil {
		b.Fatalf("wire use case: %v", err)
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ucase.Execute(ctx); err != nil {
			b.Fatal(err)
		}
	}
}
