package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/reangeline/go-shipping-products/internal/app/config"
)

func doRequest(h http.Handler, method, path string, body []byte) (int, []byte) {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.Bytes()
}

func TestWire_WithFileProvider_Smoke(t *testing.T) {

	dir := t.TempDir()
	path := filepath.Join(dir, "packs.csv")
	if err := os.WriteFile(path, []byte("250,500,1000\n2000,5000"), 0o600); err != nil {
		t.Fatalf("write packs file: %v", err)
	}

	cfg := config.Config{
		ProviderType: "file",
		FilePath:     path,
		HTTPAddr:     ":0",
	}

	container, err := Wire(cfg)
	if err != nil {
		t.Fatalf("Wire failed: %v", err)
	}

	// GET /v1/packsizes
	status, body := doRequest(container.HTTP, http.MethodGet, "/v1/packsizes", nil)
	if status != http.StatusOK {
		t.Fatalf("GET /v1/packsizes status=%d want=200 body=%s", status, string(body))
	}
	var packsResp struct {
		Sizes []int `json:"sizes"`
	}
	if err := json.Unmarshal(body, &packsResp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(packsResp.Sizes) != 5 || packsResp.Sizes[0] != 250 || packsResp.Sizes[4] != 5000 {
		t.Fatalf("unexpected sizes: %v", packsResp.Sizes)
	}

	// POST /v1/calculate
	status, body = doRequest(container.HTTP, http.MethodPost, "/v1/calculate", []byte(`{"quantity":251}`))
	if status != http.StatusOK {
		t.Fatalf("POST /v1/calculate status=%d want=200 body=%s", status, string(body))
	}
	var calcResp struct {
		TotalItems int `json:"totalItems"`
		TotalPacks int `json:"totalPacks"`
		Leftover   int `json:"leftover"`
	}
	if err := json.Unmarshal(body, &calcResp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if calcResp.TotalItems != 500 || calcResp.TotalPacks != 1 || calcResp.Leftover != 249 {
		t.Fatalf("unexpected calc resp: %+v", calcResp)
	}
}

func TestWire_UnknownProvider_ReturnsError(t *testing.T) {
	cfg := config.Config{
		ProviderType: "unknown",
		HTTPAddr:     ":0",
	}
	_, err := Wire(cfg)
	if err == nil {
		t.Fatalf("expected error for unknown provider type")
	}
}
