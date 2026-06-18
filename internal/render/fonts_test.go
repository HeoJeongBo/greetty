package render

import (
	"sort"
	"testing"
)

func TestFonts(t *testing.T) {
	fonts := Fonts()
	if len(fonts) < 100 {
		t.Fatalf("expected many fonts, got %d", len(fonts))
	}
	if !sort.StringsAreSorted(fonts) {
		t.Error("Fonts() should be sorted")
	}
	for _, want := range []string{"slant", "standard"} {
		if !contains(fonts, want) {
			t.Errorf("expected font %q in list", want)
		}
	}
}

func TestFontExists(t *testing.T) {
	tests := map[string]bool{
		"slant":    true,
		"standard": true,
		"nope":     false,
		"":         false,
	}
	for name, want := range tests {
		if got := FontExists(name); got != want {
			t.Errorf("FontExists(%q) = %v, want %v", name, got, want)
		}
	}
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
