package file

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "packs.csv")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestParsePackSizes_OK(t *testing.T) {
	got, err := ParsePackSizes("500, 250\n1000; 250 \n 2000  5000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{250, 500, 1000, 2000, 5000}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestParsePackSizes_InvalidToken(t *testing.T) {
	if _, err := ParsePackSizes("abc,100"); err == nil {
		t.Fatalf("expected error for invalid token")
	}
}

func TestParsePackSizes_NonPositive(t *testing.T) {
	if _, err := ParsePackSizes("0,250"); err == nil {
		t.Fatalf("expected error for non-positive value")
	}
	if _, err := ParsePackSizes("-5,250"); err == nil {
		t.Fatalf("expected error for non-positive value")
	}
}

func TestNew_FromFile(t *testing.T) {
	path := writeTemp(t, "250,500,1000\n2000,5000")
	prov, err := New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := prov.List()
	want := []int{250, 500, 1000, 2000, 5000}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestNew_DuplicatesAndSpaces(t *testing.T) {
	path := writeTemp(t, " 500 , 250 , 250 ; 1000 ")
	prov, err := New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := prov.List()
	want := []int{250, 500, 1000}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestNew_PathErrors(t *testing.T) {
	if _, err := New(""); err == nil {
		t.Fatalf("expected error for empty path")
	}
	if _, err := New("no/such/file.csv"); err == nil {
		t.Fatalf("expected error for not found")
	}
}

func TestList_ReturnsCopy(t *testing.T) {
	path := writeTemp(t, "250,500")
	prov, _ := New(path)
	a, _ := prov.List()
	a[0] = 999
	b, _ := prov.List()
	if b[0] != 250 {
		t.Fatalf("List must return a defensive copy")
	}
}
