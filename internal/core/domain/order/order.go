package order

import "fmt"

// Order is the entity root to agragate a customer request.
type Order struct {
	Quantity int
}

// We have a minimum quantity for the request
func NewOrder(qty int) (Order, error) {
	if qty <= 0 {
		return Order{}, fmt.Errorf("quantity must be > 0")
	}
	return Order{Quantity: qty}, nil
}
