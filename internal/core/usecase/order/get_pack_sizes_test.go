package order

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/reangeline/go-shipping-products/internal/core/ports/outbound/packsizes"
)

type fakeProvider2 struct {
	sizes []int
	err   error
}

func (f *fakeProvider2) List() ([]int, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]int(nil), f.sizes...), nil
}

func TestGetPackSizes_Execute(t *testing.T) {
	tests := []struct {
		name     string
		provider packsizes.Provider
		want     []int
		wantErr  bool
	}{
		{
			name:     "happy path",
			provider: &fakeProvider2{sizes: []int{250, 500, 1000}},
			want:     []int{250, 500, 1000},
		},
		{
			name:     "provider error",
			provider: &fakeProvider2{err: errors.New("fail")},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ucase, err := NewGetPackSizes(tt.provider)
			if err != nil {
				t.Fatalf("NewGetPackSizes unexpected error: %v", err)
			}

			out, err := ucase.Execute(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(out.Sizes, tt.want) {
				t.Fatalf("got %v, want %v", out.Sizes, tt.want)
			}
		})
	}
}
