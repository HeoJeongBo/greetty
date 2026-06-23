package render

// Gradient color presets. A preset is an ordered list of RGB stops; a glyph's
// diagonal position (column + row), normalized to [0,1], selects a color by
// linearly interpolating between adjacent stops. Cyclic presets (rainbow) wrap
// from the last stop back to the first; linear presets clamp at the ends.

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type rgb struct{ r, g, b uint8 }

type gradient struct {
	stops  []rgb
	cyclic bool // true: wrap last→first (a color wheel); false: clamp at the ends
}

// presets are the multi-color gradients selectable as a "color". Names are
// lowercase and must not collide with colorByName (the single ANSI colors).
var presets = map[string]gradient{
	"rainbow": {cyclic: true, stops: []rgb{
		{255, 0, 0}, {255, 127, 0}, {255, 255, 0},
		{0, 200, 0}, {0, 120, 255}, {75, 0, 130}, {148, 0, 211},
	}},
	"fire":   {stops: []rgb{{255, 236, 89}, {255, 140, 0}, {200, 0, 0}}},
	"ocean":  {stops: []rgb{{0, 255, 213}, {0, 140, 255}, {20, 40, 150}}},
	"sunset": {stops: []rgb{{255, 205, 96}, {255, 94, 98}, {120, 47, 194}}},
	"forest": {stops: []rgb{{183, 255, 130}, {56, 176, 0}, {0, 110, 90}}},
	"neon":   {stops: []rgb{{0, 255, 255}, {188, 19, 254}, {255, 16, 160}}},
	"mono":   {stops: []rgb{{210, 210, 210}, {80, 80, 80}}},
}

// at returns the gradient's color at position t in [0,1].
func (g gradient) at(t float64) rgb {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	stops := g.stops
	if len(stops) == 1 {
		return stops[0]
	}
	if g.cyclic {
		scaled := t * float64(len(stops))
		i := int(scaled) % len(stops)
		return lerp(stops[i], stops[(i+1)%len(stops)], scaled-float64(int(scaled)))
	}
	scaled := t * float64(len(stops)-1)
	i := int(scaled)
	if i >= len(stops)-1 {
		return stops[len(stops)-1]
	}
	return lerp(stops[i], stops[i+1], scaled-float64(i))
}

// lerp linearly interpolates between two colors at t in [0,1].
func lerp(a, b rgb, t float64) rgb {
	return rgb{
		r: uint8(float64(a.r) + (float64(b.r)-float64(a.r))*t),
		g: uint8(float64(a.g) + (float64(b.g)-float64(a.g))*t),
		b: uint8(float64(a.b) + (float64(b.b)-float64(a.b))*t),
	}
}

// applyGradient paints banner with preset g: a 24-bit truecolor diagonal sweep
// where color tracks (column + row). It returns banner unchanged when color is
// disabled (NO_COLOR / non-TTY) so no escape sequences leak into piped output —
// matching how fatih/color's Sprint self-disables. Spaces are left uncolored and
// each non-empty line is reset at its end.
func applyGradient(banner string, g gradient) string {
	if color.NoColor {
		return banner
	}
	lines := strings.Split(banner, "\n")
	maxD := float64(maxDiagonal(lines))
	var b strings.Builder
	for row, line := range lines {
		if row > 0 {
			b.WriteByte('\n')
		}
		col := 0
		for _, r := range line {
			if r == ' ' {
				b.WriteByte(' ')
				col += runeWidth(r) // spaces still advance the gradient position
				continue
			}
			c := g.at(float64(col+row) / maxD)
			fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm%c", c.r, c.g, c.b, r)
			col += runeWidth(r)
		}
		if line != "" {
			b.WriteString("\x1b[0m")
		}
	}
	return b.String()
}

// maxDiagonal returns the largest (column + row) position across the banner,
// used to normalize gradient position into [0,1]. Never returns < 1.
func maxDiagonal(lines []string) int {
	max := 0
	for row, line := range lines {
		col := 0
		for _, r := range line {
			if d := col + row; d > max {
				max = d
			}
			col += runeWidth(r)
		}
	}
	if max < 1 {
		max = 1
	}
	return max
}

// Presets returns the sorted names of the gradient color presets.
func Presets() []string {
	names := make([]string, 0, len(presets))
	for name := range presets {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// NamedColors returns the sorted names of the single (non-gradient) colors.
func NamedColors() []string {
	names := make([]string, 0, len(colorByName))
	for name := range colorByName {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// IsColor reports whether name is a known single color or gradient preset.
func IsColor(name string) bool {
	name = strings.ToLower(name)
	if _, ok := colorByName[name]; ok {
		return true
	}
	_, ok := presets[name]
	return ok
}

// Colorize renders a single line of text in the named color or gradient preset.
// Unknown names return text unchanged. Like the banner path it self-disables
// when color is off, so it is safe for swatches in piped output.
func Colorize(name, text string) string {
	name = strings.ToLower(name)
	if g, ok := presets[name]; ok {
		return applyGradient(text, g)
	}
	if c, ok := colorByName[name]; ok {
		return c.Sprint(text)
	}
	return text
}
