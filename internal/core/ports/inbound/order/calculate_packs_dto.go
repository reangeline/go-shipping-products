package order

// CalculatePacksInput is the input DTO.
// - Quantity: required (> 0)
// - PacksOverride: optional; when provided, overrides the Provider's default list.
// Must contain only positive values; duplicates will be ignored by the implementation.
type CalculatePacksInput struct {
	Quantity      int   `json:"quantity"`
	PacksOverride []int `json:"packsOverride,omitempty"`
}

// CalculatePacksOutput is the output DTO.
// - ItemsByPack: map "package size" -> "package quantity"
// - TotalItems: sum(size*count) of ItemsByPack
// - TotalPacks: sum of counts
// - Leftover: TotalItems - Quantity
type CalculatePacksOutput struct {
	ItemsByPack map[int]int `json:"itemsByPack"`
	TotalItems  int         `json:"totalItems"`
	TotalPacks  int         `json:"totalPacks"`
	Leftover    int         `json:"leftover"`
}
