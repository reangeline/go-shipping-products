package presenter

import (
	"errors"
	"net/http"

	usecases "github.com/reangeline/go-shipping-products/internal/core/usecase/order"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// MapError converts use case errors to (status, ErrorBody).
func MapError(err error) (int, ErrorBody) {
	switch {
	case errors.Is(err, usecases.ErrInvalidQuantity):
		return http.StatusBadRequest, ErrorBody{Code: "invalid_quantity", Message: "quantity must be > 0"}
	case errors.Is(err, usecases.ErrInvalidPackInOverride):
		return http.StatusBadRequest, ErrorBody{Code: "invalid_pack", Message: "packsOverride must contain positive integers"}
	case errors.Is(err, usecases.ErrNoPackSizes):
		return http.StatusUnprocessableEntity, ErrorBody{Code: "no_pack_sizes", Message: "no pack sizes available"}
	default:

		return http.StatusInternalServerError, ErrorBody{Code: "internal_error", Message: "unexpected error"}
	}
}
