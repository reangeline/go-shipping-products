package order

import (
	"context"
	"errors"

	uc "github.com/reangeline/go-shipping-products/internal/core/ports/inbound/order"
	"github.com/reangeline/go-shipping-products/internal/core/ports/outbound/packsizes"
)

type getPackSizes struct {
	provider packsizes.Provider
}

// compile-time check to keep my cohesion with my conctact
var _ uc.GetPackSizes = (*getPackSizes)(nil)

func NewGetPackSizes(provider packsizes.Provider) (uc.GetPackSizes, error) {
	if provider == nil {
		return nil, errors.New("nil packsizes.Provider")
	}
	return &getPackSizes{provider: provider}, nil
}

func (g *getPackSizes) Execute(ctx context.Context) (uc.GetPackSizesOutput, error) {
	_ = ctx // (no-op for now; kept for future cancellation/telemetry)

	sizes, err := g.provider.List()
	if err != nil {
		return uc.GetPackSizesOutput{}, err
	}
	return uc.GetPackSizesOutput{Sizes: sizes}, nil
}
