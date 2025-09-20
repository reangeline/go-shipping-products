package ginadapter

import (
	"bytes"
	"context"
	"encoding/json"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	ctr "github.com/reangeline/go-shipping-products/internal/adapters/inbound/http/order"
	uc "github.com/reangeline/go-shipping-products/internal/core/ports/inbound/order"
	usecases "github.com/reangeline/go-shipping-products/internal/core/usecase/order"
)

// ---- fakes usecases ----

type fakeCalc struct {
	out uc.CalculatePacksOutput
	err error
}

func (f *fakeCalc) Execute(_ context.Context, in uc.CalculatePacksInput) (uc.CalculatePacksOutput, error) {
	return f.out, f.err
}

type fakeGet struct {
	out uc.GetPackSizesOutput
	err error
}

func (f *fakeGet) Execute(_ context.Context) (uc.GetPackSizesOutput, error) {
	return f.out, f.err
}

func newTestHandler(calc *fakeCalc, get *fakeGet) http.Handler {
	controller := ctr.NewController(calc, get)
	return BuildHandler(controller)
}

// ---- tests ----

func TestGET_PackSizes_OK(t *testing.T) {
	h := newTestHandler(
		&fakeCalc{},
		&fakeGet{out: uc.GetPackSizesOutput{Sizes: []int{250, 500, 1000}}},
	)

	req := httptest.NewRequest(http.MethodGet, "/v1/packsizes", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status got=%d want=%d", rec.Code, http.StatusOK)
	}

	var body struct {
		Sizes []int `json:"sizes"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	want := []int{250, 500, 1000}
	if len(body.Sizes) != len(want) || body.Sizes[0] != 250 || body.Sizes[1] != 500 || body.Sizes[2] != 1000 {
		t.Fatalf("sizes got=%v want=%v", body.Sizes, want)
	}
}

func TestGET_PackSizes_ProviderError_500(t *testing.T) {
	h := newTestHandler(
		&fakeCalc{},
		&fakeGet{err: errors.New("boom")},
	)

	req := httptest.NewRequest(http.MethodGet, "/v1/packsizes", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status got=%d want=%d", rec.Code, http.StatusInternalServerError)
	}
	var body struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Code == "" {
		t.Fatalf("expected error body with code, got=%s", rec.Body.String())
	}
}

func TestPOST_Calculate_OK(t *testing.T) {
	out := uc.CalculatePacksOutput{
		ItemsByPack: map[int]int{5000: 2, 2000: 1, 250: 1},
		TotalItems:  12250,
		TotalPacks:  4,
		Leftover:    249,
	}
	h := newTestHandler(
		&fakeCalc{out: out},
		&fakeGet{},
	)

	payload := `{"quantity":12001}`
	req := httptest.NewRequest(http.MethodPost, "/v1/calculate", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status got=%d want=%d", rec.Code, http.StatusOK)
	}

	var body struct {
		ItemsByPack map[int]int `json:"itemsByPack"`
		TotalItems  int         `json:"totalItems"`
		TotalPacks  int         `json:"totalPacks"`
		Leftover    int         `json:"leftover"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body.TotalItems != 12250 || body.TotalPacks != 4 || body.Leftover != 249 {
		t.Fatalf("mismatch body=%+v", body)
	}
	if body.ItemsByPack[5000] != 2 || body.ItemsByPack[2000] != 1 || body.ItemsByPack[250] != 1 {
		t.Fatalf("itemsByPack mismatch: %v", body.ItemsByPack)
	}
}

func TestPOST_Calculate_InvalidJSON_400(t *testing.T) {
	h := newTestHandler(&fakeCalc{}, &fakeGet{})

	req := httptest.NewRequest(http.MethodPost, "/v1/calculate", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status got=%d want=%d", rec.Code, http.StatusBadRequest)
	}
	var body struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Code != "invalid_request" {
		t.Fatalf("expected invalid_request, got=%s", body.Code)
	}
}

func TestPOST_Calculate_InvalidQuantity_400(t *testing.T) {
	h := newTestHandler(
		&fakeCalc{err: usecases.ErrInvalidQuantity},
		&fakeGet{},
	)

	payload := `{"quantity":0}`
	req := httptest.NewRequest(http.MethodPost, "/v1/calculate", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status got=%d want=%d", rec.Code, http.StatusBadRequest)
	}
	var body struct{ Code string }
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Code != "invalid_quantity" {
		t.Fatalf("expected invalid_quantity, got=%s (resp=%s)", body.Code, rec.Body.String())
	}
}

func TestPOST_Calculate_NoPackSizes_422(t *testing.T) {
	h := newTestHandler(
		&fakeCalc{err: usecases.ErrNoPackSizes},
		&fakeGet{},
	)

	payload := `{"quantity":5}`
	req := httptest.NewRequest(http.MethodPost, "/v1/calculate", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status got=%d want=%d", rec.Code, http.StatusUnprocessableEntity)
	}
	var body struct{ Code string }
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Code != "no_pack_sizes" {
		t.Fatalf("expected no_pack_sizes, got=%s (resp=%s)", body.Code, rec.Body.String())
	}
}

func TestOPTIONS_CORS_Preflight(t *testing.T) {
	h := newTestHandler(&fakeCalc{}, &fakeGet{})

	req := httptest.NewRequest(http.MethodOptions, "/v1/calculate", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status got=%d want=%d", rec.Code, http.StatusNoContent)
	}

	if rec.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("expected CORS headers")
	}
}
