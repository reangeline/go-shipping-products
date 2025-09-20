package order

// Combination represents the result of the calculation: how many packages of each size,
// total items, total packages, and leftovers.
type Combination struct {
	ItemsByPack map[int]int // size -> count
	TotalItems  int
	TotalPacks  int
	Leftover    int
}

func (c Combination) IsZero() bool {
	return c.TotalItems == 0 && c.TotalPacks == 0
}
