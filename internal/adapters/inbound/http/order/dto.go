package order

// Transport DTOs (used only in the HTTP layer; different from use case DTOs).
type CalculateRequest struct {
	Quantity      int   `json:"quantity"`
	PacksOverride []int `json:"packsOverride,omitempty"`
}

type CalculateResponse struct {
	ItemsByPack map[int]int `json:"itemsByPack"`
	TotalItems  int         `json:"totalItems"`
	TotalPacks  int         `json:"totalPacks"`
	Leftover    int         `json:"leftover"`
}

type PackSizesResponse struct {
	Sizes []int `json:"sizes"`
}
