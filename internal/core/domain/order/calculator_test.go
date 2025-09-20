package order

import (
	"reflect"
	"sort"
	"testing"
)

// helpers
func mkPacks(t *testing.T, sizes ...int) []Pack {
	t.Helper()
	out := make([]Pack, 0, len(sizes))
	for _, s := range sizes {
		p, err := NewPack(s)
		if err != nil {
			t.Fatalf("mkPacks: invalid pack size %d: %v", s, err)
		}
		out = append(out, p)
	}
	return out
}

func sortedKeys(m map[int]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func assertInvariants(t *testing.T, qty int, comb Combination, packs []int) {
	t.Helper()

	// TotalItems >= qty
	if comb.TotalItems < qty {
		t.Fatalf("TotalItems(%d) < qty(%d)", comb.TotalItems, qty)
	}

	if comb.Leftover != comb.TotalItems-qty {
		t.Fatalf("Leftover inconsistente: got=%d want=%d", comb.Leftover, comb.TotalItems-qty)
	}
	// soma(size*count) == TotalItems
	sum := 0
	for size, cnt := range comb.ItemsByPack {
		sum += size * cnt
	}
	if sum != comb.TotalItems {
		t.Fatalf("soma=%d != TotalItems=%d", sum, comb.TotalItems)
	}
}

func TestPackCalculator_CanonicalCases(t *testing.T) {
	type want struct {
		totalItems int
		totalPacks int
		byPack     map[int]int
	}
	tests := []struct {
		name  string
		qty   int
		packs []int
		want  want
	}{
		{
			name:  "qty=1 -> 1x250",
			qty:   1,
			packs: []int{250, 500, 1000, 2000, 5000},
			want:  want{250, 1, map[int]int{250: 1}},
		},
		{
			name:  "qty=250 -> 1x250",
			qty:   250,
			packs: []int{500, 250, 1000},
			want:  want{250, 1, map[int]int{250: 1}},
		},
		{
			name:  "qty=251 -> 1x500 (less overpack, less packages)",
			qty:   251,
			packs: []int{250, 500, 1000},
			want:  want{500, 1, map[int]int{500: 1}},
		},
		{
			name:  "qty=501 -> 500+250=750 em vez de 1000",
			qty:   501,
			packs: []int{250, 500, 1000},
			want:  want{750, 2, map[int]int{500: 1, 250: 1}},
		},
		{
			name:  "qty=12001 -> 2x5000+1x2000+1x250",
			qty:   12001,
			packs: []int{250, 500, 1000, 2000, 5000},
			want:  want{12250, 4, map[int]int{5000: 2, 2000: 1, 250: 1}},
		},
	}

	pc := NewPackCalculator()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			packs := mkPacks(t, tc.packs...)
			comb, err := pc.Calculate(tc.qty, packs)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertInvariants(t, tc.qty, comb, tc.packs)

			if comb.TotalItems != tc.want.totalItems {
				t.Fatalf("TotalItems got=%d want=%d", comb.TotalItems, tc.want.totalItems)
			}
			if comb.TotalPacks != tc.want.totalPacks {
				t.Fatalf("TotalPacks got=%d want=%d", comb.TotalPacks, tc.want.totalPacks)
			}

			gotK := sortedKeys(comb.ItemsByPack)
			wantK := sortedKeys(tc.want.byPack)
			if !reflect.DeepEqual(gotK, wantK) {
				t.Fatalf("packs diferentes: got %v want %v", gotK, wantK)
			}
			for _, k := range wantK {
				if comb.ItemsByPack[k] != tc.want.byPack[k] {
					t.Fatalf("count %d: got=%d want=%d", k, comb.ItemsByPack[k], tc.want.byPack[k])
				}
			}
		})
	}
}

func TestPackCalculator_ErrorsAndEdges(t *testing.T) {
	pc := NewPackCalculator()

	if _, err := pc.Calculate(0, mkPacks(t, 250)); err == nil {
		t.Fatalf("expected error for qty<=0")
	}
	if _, err := pc.Calculate(10, []Pack{}); err == nil {
		t.Fatalf("expected error for empty packs")
	}
	if _, err := pc.Calculate(10, []Pack{{Size: -1}}); err == nil {
		t.Fatalf("expected error for invalid pack size")
	}
}
