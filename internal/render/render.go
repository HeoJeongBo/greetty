// Package render turns a Config into a colorful ASCII banner. It is defensive:
// bad fonts or colors fall back to defaults and it never panics, because its
// output runs during shell startup.
package render

import (
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"

	"github.com/HeoJeongBo/greetty/internal/config"
)

var colorByName = map[string]*color.Color{
	"black":   color.New(color.FgBlack),
	"red":     color.New(color.FgRed),
	"green":   color.New(color.FgGreen),
	"yellow":  color.New(color.FgYellow),
	"blue":    color.New(color.FgBlue),
	"magenta": color.New(color.FgMagenta),
	"cyan":    color.New(color.FgCyan),
	"white":   color.New(color.FgWhite),
}

// Render returns the full greeting string (emoji header, ASCII banner, and a
// dotted divider) ready to print to stdout.
func Render(cfg config.Config) string {
	banner := figure.NewFigure(cfg.Text, cfg.Font, true).String()
	// go-figure returns an empty/odd result for unknown fonts; fall back.
	if strings.TrimSpace(banner) == "" {
		banner = figure.NewFigure(cfg.Text, "standard", true).String()
	}

	c, ok := colorByName[strings.ToLower(cfg.Color)]
	if !ok {
		c = colorByName["cyan"]
	}

	width := bannerWidth(banner)
	divider := dottedDivider(width)

	var b strings.Builder
	b.WriteByte('\n')
	if cfg.Emoji != "" {
		b.WriteString("  " + cfg.Emoji + "\n")
	}
	b.WriteString(c.Sprint(banner))
	b.WriteString(color.New(color.Faint).Sprint(divider))
	b.WriteByte('\n')
	return b.String()
}

// bannerWidth returns the length of the widest line in the banner.
func bannerWidth(banner string) int {
	max := 0
	for _, line := range strings.Split(banner, "\n") {
		if n := len([]rune(line)); n > max {
			max = n
		}
	}
	if max < 8 {
		max = 8
	}
	return max
}

// dottedDivider builds a "· · ·" style developer divider sized to the banner.
func dottedDivider(width int) string {
	n := width / 2
	if n < 4 {
		n = 4
	}
	return strings.TrimRight(strings.Repeat("· ", n), " ") + "\n"
}
