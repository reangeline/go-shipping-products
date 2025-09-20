package order

import "context"

// GetPackSizes exposes the current list of package sizes
// (provided by an outbound Provider).
type GetPackSizes interface {
	Execute(ctx context.Context) (GetPackSizesOutput, error)
}
