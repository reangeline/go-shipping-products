package order

import (
	"context"
	"errors"
	"reflect"
	"testing"

	uc "github.com/reangeline/go-shipping-products/internal/core/ports/inbound/order"
)

// ---------- fakes use cases ----------

type fakeCalc struct {
	out    uc.CalculatePacksOutput
	err    error
	lastIn uc.CalculatePacksInput
}

func (f *fakeCalc) Execute(ctx context.Context, in uc.CalculatePacksInput) (uc.CalculatePacksOutput, error) {
	f.lastIn = in
	return f.out, f.err
}

type fakeGet struct {
	out uc.GetPackSizesOutput
	err error
}

func (f *fakeGet) Execute(ctx context.Context) (uc.GetPackSizesOutput, error) {
	return f.out, f.err
}

// ---------- tests ----------

func TestController_HandleCalculate_Success_NoOverride(t *testing.T) {
	fc := &fakeCalc{
		out: uc.CalculatePacksOutput{
			ItemsByPack: map[int]int{5000: 2, 2000: 1, 250: 1},
			TotalItems:  12250,
			TotalPacks:  4,
			Leftover:    249,
		},
	}
	ctrl := NewController(fc, &fakeGet{})

	req := CalculateRequest{Quantity: 12001}
	res, err := ctrl.HandleCalculate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := CalculateResponse{
		ItemsByPack: map[int]int{5000: 2, 2000: 1, 250: 1},
		TotalItems:  12250,
		TotalPacks:  4,
		Leftover:    249,
	}
	if !reflect.DeepEqual(res, want) {
		t.Fatalf("response mismatch:\n got=%v\nwant=%v", res, want)
	}

	if fc.lastIn.Quantity != 12001 {
		t.Fatalf("quantity mismatch: got=%d want=%d", fc.lastIn.Quantity, 12001)
	}
	if len(fc.lastIn.PacksOverride) > 0 {
		t.Fatalf("override should be empty; got=%v", fc.lastIn.PacksOverride)
	}
}

func TestController_HandleCalculate_Success_WithOverride(t *testing.T) {
	fc := &fakeCalc{
		out: uc.CalculatePacksOutput{
			ItemsByPack: map[int]int{7: 1, 3: 1},
			TotalItems:  10,
			TotalPacks:  2,
			Leftover:    0,
		},
	}
	ctrl := NewController(fc, &fakeGet{})

	req := CalculateRequest{Quantity: 10, PacksOverride: []int{3, 7}}
	res, err := ctrl.HandleCalculate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.TotalItems != 10 || res.TotalPacks != 2 {
		t.Fatalf("unexpected totals: %+v", res)
	}

	if !reflect.DeepEqual(fc.lastIn.PacksOverride, []int{3, 7}) || fc.lastIn.Quantity != 10 {
		t.Fatalf("use case received wrong input: %+v", fc.lastIn)
	}
}

func TestController_HandleCalculate_ErrorIsPropagated(t *testing.T) {
	wantErr := errors.New("boom")
	fc := &fakeCalc{err: wantErr}
	ctrl := NewController(fc, &fakeGet{})

	_, err := ctrl.HandleCalculate(context.Background(), CalculateRequest{Quantity: 1})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error to be propagated; got=%v", err)
	}
}

func TestController_HandleGetPackSizes_Success(t *testing.T) {
	fg := &fakeGet{out: uc.GetPackSizesOutput{Sizes: []int{250, 500, 1000}}}
	ctrl := NewController(&fakeCalc{}, fg)

	res, err := ctrl.HandleGetPackSizes(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := PackSizesResponse{Sizes: []int{250, 500, 1000}}
	if !reflect.DeepEqual(res, want) {
		t.Fatalf("response mismatch: got=%v want=%v", res, want)
	}
}

func TestController_HandleGetPackSizes_ErrorIsPropagated(t *testing.T) {
	wantErr := errors.New("provider fail")
	fg := &fakeGet{err: wantErr}
	ctrl := NewController(&fakeCalc{}, fg)

	_, err := ctrl.HandleGetPackSizes(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error to be propagated; got=%v", err)
	}
}
