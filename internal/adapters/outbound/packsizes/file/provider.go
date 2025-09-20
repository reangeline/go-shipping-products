package file

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/reangeline/go-shipping-products/internal/core/ports/outbound/packsizes"
)

const DefaultPathEnv = "PACK_SIZES_FILE"

var (
	ErrPathNotSet      = errors.New("pack sizes file path not set")
	ErrNoValidPack     = errors.New("no valid pack sizes parsed from file")
	ErrInvalidPackSize = errors.New("pack size must be > 0")
)

type provider struct {
	sizes []int // ordenado asc e sem duplicados
}

// compile-time check
var _ packsizes.Provider = (*provider)(nil)

// New creates a Provider by reading and parsing the file pointed to by path.
// The file can contain values ​​separated by commas, semicolons, spaces, or newlines.
// Ex.: "250,500,1000\n2000,5000"
func New(path string) (packsizes.Provider, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, ErrPathNotSet
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}

	sizes, err := ParsePackSizes(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %w", path, err)
	}
	if len(sizes) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrNoValidPack, path)
	}

	return &provider{sizes: sizes}, nil
}

func (p *provider) List() ([]int, error) {
	out := make([]int, len(p.sizes))
	copy(out, p.sizes)
	return out, nil
}

// ParsePackSizes parses textual content containing sizes separated by commas,
// semicolons, spaces, or line breaks (“loose” CSV).
// Rules: remove spaces, reject <= 0, remove duplicates, and sort asc.
func ParsePackSizes(s string) ([]int, error) {
	sep := func(r rune) bool {
		switch r {
		case ',', ';', ' ', '\n', '\r', '\t':
			return true
		default:
			return false
		}
	}
	tokens := strings.FieldsFunc(s, sep)

	seen := make(map[int]struct{}, len(tokens))
	out := make([]int, 0, len(tokens))

	for _, tok := range tokens {
		n, err := strconv.Atoi(strings.TrimSpace(tok))
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", tok, err)
		}
		if n <= 0 {
			return nil, ErrInvalidPackSize
		}
		if _, dup := seen[n]; dup {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}

	sort.Ints(out)
	return out, nil
}
