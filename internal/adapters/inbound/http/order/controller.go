package order

import (
	"context"
	"errors"

	uc "github.com/reangeline/go-shipping-products/internal/core/ports/inbound/order"
)

// Controller contains only orchestration logic (transport ↔ use cases).
// It does not depend on the HTTP framework.
type Controller struct {
	Calc uc.CalculatePacks
	Get  uc.GetPackSizes
}

func NewController(calc uc.CalculatePacks, get uc.GetPackSizes) *Controller {
	return &Controller{Calc: calc, Get: get}
}

// HandleCalculate maps request → use case → response.
func (c *Controller) HandleCalculate(ctx context.Context, req CalculateRequest) (CalculateResponse, error) {
	out, err := c.Calc.Execute(ctx, uc.CalculatePacksInput{
		Quantity:      req.Quantity,
		PacksOverride: req.PacksOverride,
	})
	if err != nil {
		return CalculateResponse{}, err
	}
	return CalculateResponse{
		ItemsByPack: out.ItemsByPack,
		TotalItems:  out.TotalItems,
		TotalPacks:  out.TotalPacks,
		Leftover:    out.Leftover,
	}, nil
}

// HandleGetPackSizes simply delegates to the use case.
func (c *Controller) HandleGetPackSizes(ctx context.Context) (PackSizesResponse, error) {
	out, err := c.Get.Execute(ctx)
	if err != nil {
		return PackSizesResponse{}, err
	}
	return PackSizesResponse{Sizes: out.Sizes}, nil
}

// Helpers for known errors (optional; avoids importing implementation details).
var (
	ErrInvalidQuantity       = errors.New("quantity must be > 0")
	ErrNoPackSizes           = errors.New("no pack sizes available")
	ErrInvalidPackInOverride = errors.New("override contains non-positive pack size")
)
