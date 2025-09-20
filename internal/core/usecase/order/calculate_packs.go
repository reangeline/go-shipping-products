package order

import (
	"context"
	"errors"
	"sort"

	domain "github.com/reangeline/go-shipping-products/internal/core/domain/order"
	uc "github.com/reangeline/go-shipping-products/internal/core/ports/inbound/order"
	"github.com/reangeline/go-shipping-products/internal/core/ports/outbound/packsizes"
)

var (
	ErrInvalidQuantity       = errors.New("quantity must be > 0")
	ErrNoPackSizes           = errors.New("no pack sizes available")
	ErrInvalidPackInOverride = errors.New("override contains non-positive pack size")
)

type calculatePacks struct {
	calc     domain.PackCalculator
	provider packsizes.Provider
}

// compile-time check to keep my cohesion with my conctact
var _ uc.CalculatePacks = (*calculatePacks)(nil)

func NewCalculatePacks(calc domain.PackCalculator, provider packsizes.Provider) (uc.CalculatePacks, error) {
	if calc == nil {
		return nil, errors.New("nil PackCalculator")
	}
	if provider == nil {
		return nil, errors.New("nil packsizes.Provider")
	}
	return &calculatePacks{calc: calc, provider: provider}, nil
}

func (c *calculatePacks) Execute(ctx context.Context, in uc.CalculatePacksInput) (uc.CalculatePacksOutput, error) {
	_ = ctx // (no-op for now; kept for future cancellation/telemetry)

	if in.Quantity <= 0 {
		return uc.CalculatePacksOutput{}, ErrInvalidQuantity
	}

	// Validate the size
	var sizes []int
	if len(in.PacksOverride) > 0 {
		norm, err := normalizeOverride(in.PacksOverride)
		if err != nil {
			return uc.CalculatePacksOutput{}, err
		}
		sizes = norm
	} else {
		list, err := c.provider.List()
		if err != nil {
			return uc.CalculatePacksOutput{}, err
		}
		if len(list) == 0 {
			return uc.CalculatePacksOutput{}, ErrNoPackSizes
		}
		sizes = list
	}

	// Converte []int -> []domain.Pack
	packs := make([]domain.Pack, 0, len(sizes))
	for _, s := range sizes {
		p, err := domain.NewPack(s)
		if err != nil {
			return uc.CalculatePacksOutput{}, err
		}
		packs = append(packs, p)
	}

	comb, err := c.calc.Calculate(in.Quantity, packs)
	if err != nil {
		return uc.CalculatePacksOutput{}, err
	}

	out := uc.CalculatePacksOutput{
		ItemsByPack: comb.ItemsByPack,
		TotalItems:  comb.TotalItems,
		TotalPacks:  comb.TotalPacks,
		Leftover:    comb.Leftover,
	}
	return out, nil
}

// normalizeOverride applies minimal rules to the override coming from the caller:
// - rejects values ​​<= 0 (explicit error)
// - removes duplicates
// - sorts ascending
func normalizeOverride(in []int) ([]int, error) {
	if len(in) == 0 {
		return nil, nil
	}
	seen := make(map[int]struct{}, len(in))
	out := make([]int, 0, len(in))
	for _, v := range in {
		if v <= 0 {
			return nil, ErrInvalidPackInOverride
		}
		if _, dup := seen[v]; dup {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Ints(out)
	return out, nil
}
