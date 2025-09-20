// internal/app/wiring.go
package app

import (
	"fmt"
	"net/http"

	"github.com/reangeline/go-shipping-products/internal/app/config"
	domain "github.com/reangeline/go-shipping-products/internal/core/domain/order"
	inbound "github.com/reangeline/go-shipping-products/internal/core/ports/inbound/order"
	usecases "github.com/reangeline/go-shipping-products/internal/core/usecase/order"

	fileProv "github.com/reangeline/go-shipping-products/internal/adapters/outbound/packsizes/file"

	ginadapter "github.com/reangeline/go-shipping-products/internal/adapters/inbound/http/gin"
	ctr "github.com/reangeline/go-shipping-products/internal/adapters/inbound/http/order"
	"github.com/reangeline/go-shipping-products/internal/core/ports/outbound/packsizes"
)

type Container struct {
	Calc inbound.CalculatePacks
	Get  inbound.GetPackSizes
	HTTP http.Handler
}

func Wire(cfg config.Config) (*Container, error) {
	var prov packsizes.Provider
	var err error

	switch cfg.ProviderType {
	case "file":
		prov, err = fileProv.New(cfg.FilePath)
	default:
		return nil, fmt.Errorf("unknown provider type: %s", cfg.ProviderType)
	}

	if err != nil {
		return nil, fmt.Errorf("init provider: %w", err)
	}

	calcDomain := domain.NewPackCalculator()

	calcUC, err := usecases.NewCalculatePacks(calcDomain, prov)
	if err != nil {
		return nil, err
	}
	getUC, err := usecases.NewGetPackSizes(prov)
	if err != nil {
		return nil, err
	}

	controller := ctr.NewController(calcUC, getUC)
	handler := ginadapter.BuildHandler(controller)

	return &Container{
		Calc: calcUC,
		Get:  getUC,
		HTTP: handler,
	}, nil
}
