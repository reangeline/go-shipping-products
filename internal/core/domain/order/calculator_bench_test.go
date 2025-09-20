package order

import "testing"

func BenchmarkPackCalculator_Various(b *testing.B) {
	pc := NewPackCalculator()
	packs := []Pack{{250}, {500}, {1000}, {2000}, {5000}}

	cases := []struct {
		name string
		qty  int
	}{
		{"Small", 12001},
		{"Medium", 100_000},
		{"Large", 500_000},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				if _, err := pc.Calculate(tc.qty, packs); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
