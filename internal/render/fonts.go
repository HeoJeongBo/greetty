package render

import (
	"sort"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

// fallbackFont is always present in go-figure's embedded assets and is used as a
// safe substitute whenever a requested font does not exist.
const fallbackFont = "standard"

// Fonts returns the sorted list of available go-figure font names. go-figure
// embeds its fonts as assets named "fonts/<name>.flf"; AssetNames returns those
// paths in arbitrary order, so we strip the prefix/suffix and sort.
func Fonts() []string {
	var names []string
	for _, asset := range figure.AssetNames() {
		if strings.HasPrefix(asset, "fonts/") && strings.HasSuffix(asset, ".flf") {
			name := strings.TrimSuffix(strings.TrimPrefix(asset, "fonts/"), ".flf")
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// FontExists reports whether name is a known go-figure font. It uses
// figure.Asset (which returns an error for a missing font) rather than
// figure.NewFigure (which panics), so it is safe to call during shell startup.
func FontExists(name string) bool {
	if name == "" {
		return false
	}
	_, err := figure.Asset("fonts/" + name + ".flf")
	return err == nil
}
