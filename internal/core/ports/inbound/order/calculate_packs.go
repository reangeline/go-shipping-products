package order

import "context"

// CalculatePacks defines the contract for the main use case:
// given a requested quantity, calculate the best combination of packs,
// prioritizing (1) minimizing the total number of items shipped and (2) in a tie,
// minimizing the number of packs.
type CalculatePacks interface {
	Execute(ctx context.Context, in CalculatePacksInput) (CalculatePacksOutput, error)
}
