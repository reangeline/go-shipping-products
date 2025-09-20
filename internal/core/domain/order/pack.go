// internal/core/domain/order/pack.go
package order

import "fmt"

// Pack é um Value Object imutável que representa o tamanho de um pacote.
type Pack struct {
	Size int
}

func NewPack(size int) (Pack, error) {
	if size <= 0 {
		return Pack{}, fmt.Errorf("invalid pack size: %d", size)
	}
	return Pack{Size: size}, nil
}
