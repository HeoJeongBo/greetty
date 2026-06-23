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

// defaultEmojiHeight is the block height used when the banner has no ASCII
// segment to measure against (emoji-only text). It roughly matches the height
// of the default "slant" font.
const defaultEmojiHeight = 6

// Render returns the full greeting string (emoji header, ASCII banner, and a
// dotted divider) ready to print to stdout.
func Render(cfg config.Config) string {
	banner := buildBanner(cfg.Text, cfg.Font)

	width := bannerWidth(banner)
	divider := dottedDivider(width)

	var b strings.Builder
	b.WriteByte('\n')
	if cfg.Emoji != "" {
		b.WriteString("  " + cfg.Emoji + "\n")
	}
	if g, ok := presets[strings.ToLower(cfg.Color)]; ok {
		b.WriteString(applyGradient(banner, g))
	} else {
		c, ok := colorByName[strings.ToLower(cfg.Color)]
		if !ok {
			c = colorByName["cyan"]
		}
		b.WriteString(c.Sprint(banner))
	}
	b.WriteString(color.New(color.Faint).Sprint(divider))
	b.WriteByte('\n')
	return b.String()
}

// runeWidth approximates the terminal column width of a single rune, matching
// displayWidth's accounting so the gradient stays aligned across wide emoji.
func runeWidth(r rune) int {
	switch {
	case isJoiner(r):
		return 0
	case r > 0x7E:
		return 2
	default:
		return 1
	}
}

// buildBanner renders text into the banner body. Pure-ASCII text takes the
// original single go-figure path (byte-identical to before). When text contains
// emoji, each emoji is drawn as a large block sized to the letters' height and
// joined horizontally with the ASCII art.
func buildBanner(text, font string) string {
	// go-figure panics on an unknown font name (newFont→MustAsset). Resolve to a
	// safe font once, up front, so neither render path can crash shell startup.
	if !FontExists(font) {
		font = fallbackFont
	}
	segs := segments(text)

	hasEmoji := false
	for _, s := range segs {
		if s.isEmoji {
			hasEmoji = true
			break
		}
	}

	// Fast path: no emoji — keep the exact original behavior.
	if !hasEmoji {
		banner := figure.NewFigure(text, font, true).String()
		if strings.TrimSpace(banner) == "" {
			banner = figure.NewFigure(text, "standard", true).String()
		}
		return banner
	}

	// Render each non-emoji segment first so we know the banner's height, then
	// build emoji blocks to match.
	asciiBlocks := make(map[int][]string, len(segs))
	height := 0
	for i, s := range segs {
		if s.isEmoji {
			continue
		}
		blk := asciiBlock(s.text, font)
		asciiBlocks[i] = blk
		if len(blk) > height {
			height = len(blk)
		}
	}
	if height == 0 {
		height = defaultEmojiHeight
	}

	blocks := make([][]string, 0, len(segs))
	for i, s := range segs {
		if s.isEmoji {
			blocks = append(blocks, emojiSegmentBlock(s.text, height))
		} else {
			blocks = append(blocks, asciiBlocks[i])
		}
	}

	banner := joinBlocks(blocks)
	if strings.TrimSpace(banner) == "" {
		// Last-resort fallback so output is never blank during shell startup.
		banner = text + "\n"
	}
	return banner
}

// segment is a run of the banner text that is either all emoji or all
// figure-renderable (ASCII) characters.
type segment struct {
	text    string
	isEmoji bool
}

// segments splits text into emoji and ASCII runs. A rune is treated as emoji
// when it falls outside the printable-ASCII range go-figure can render
// (' '..'~'); everything else is ASCII. Joiners and modifiers (ZWJ, variation
// selectors, skin-tone modifiers) attach to the preceding emoji so composed
// glyphs like 👨‍💻, ✈️, and 🇰🇷 stay in a single segment.
//
// This is an approximation, not full UAX #29 grapheme segmentation; it is
// dependency-free and correct for common cases. If exact emoji clustering is
// ever needed, github.com/rivo/uniseg is the robust alternative.
func segments(text string) []segment {
	var segs []segment
	var cur strings.Builder
	curEmoji := false
	started := false

	flush := func() {
		if cur.Len() > 0 {
			segs = append(segs, segment{text: cur.String(), isEmoji: curEmoji})
			cur.Reset()
		}
	}

	for _, r := range text {
		emoji := isEmojiRune(r)
		// Joiners/modifiers continue whatever the current run is (they only ever
		// follow emoji in practice, keeping composed glyphs together).
		if started && isJoiner(r) {
			cur.WriteRune(r)
			continue
		}
		if !started {
			curEmoji = emoji
			started = true
		} else if emoji != curEmoji {
			flush()
			curEmoji = emoji
		}
		cur.WriteRune(r)
	}
	flush()
	return segs
}

// isEmojiRune reports whether r is outside the printable-ASCII range that
// go-figure can render.
func isEmojiRune(r rune) bool {
	return r > 0x7E || r < 0x20
}

// isJoiner reports whether r is a zero-width joiner, variation selector, or
// skin-tone modifier that should stay attached to the preceding emoji.
func isJoiner(r rune) bool {
	switch {
	case r == 0x200D: // zero-width joiner
		return true
	case r >= 0xFE00 && r <= 0xFE0F: // variation selectors
		return true
	case r >= 0x1F3FB && r <= 0x1F3FF: // skin-tone modifiers
		return true
	}
	return false
}

// asciiBlock renders an ASCII text segment to figlet rows. strict=false so a
// stray non-ASCII rune degrades to '?' rather than calling log.Fatal during
// shell startup. An empty result falls back to the "standard" font.
func asciiBlock(text, font string) []string {
	rows := figure.NewFigure(text, font, false).Slicify()
	if blankRows(rows) {
		rows = figure.NewFigure(text, "standard", false).Slicify()
	}
	return rows
}

// blankRows reports whether rows is empty or all-whitespace.
func blankRows(rows []string) bool {
	for _, r := range rows {
		if strings.TrimSpace(r) != "" {
			return false
		}
	}
	return true
}

// emojiSegmentBlock renders an emoji segment. A registered emoji uses its
// hand-drawn ASCII art; anything else falls back to a repeated-glyph block sized
// to the banner height. A segment with multiple emoji is joined horizontally.
func emojiSegmentBlock(emoji string, height int) []string {
	glyphs := splitEmoji(emoji)
	if len(glyphs) <= 1 {
		if art, ok := artFor(emoji); ok {
			return art
		}
		return emojiBlock(emoji, height)
	}
	blocks := make([][]string, 0, len(glyphs))
	for _, g := range glyphs {
		if art, ok := artFor(g); ok {
			blocks = append(blocks, art)
		} else {
			blocks = append(blocks, emojiBlock(g, height))
		}
	}
	return strings.Split(strings.TrimRight(joinBlocks(blocks), "\n"), "\n")
}

// splitEmoji breaks an emoji run into individual glyphs, keeping joiners and
// modifiers attached to the base they follow (so 👨‍💻 stays one glyph).
func splitEmoji(emoji string) []string {
	var glyphs []string
	var cur strings.Builder
	prevJoiner := false
	for _, r := range emoji {
		if isJoiner(r) {
			cur.WriteRune(r)
			prevJoiner = true
			continue
		}
		// A base rune starts a new glyph unless it directly follows a joiner
		// (ZWJ/modifier), in which case it stays part of the composed glyph.
		if cur.Len() > 0 && !prevJoiner {
			glyphs = append(glyphs, cur.String())
			cur.Reset()
		}
		cur.WriteRune(r)
		prevJoiner = false
	}
	if cur.Len() > 0 {
		glyphs = append(glyphs, cur.String())
	}
	return glyphs
}

// emojiBlock builds a roughly square block of the emoji glyph at the given
// height. Each emoji is ~2 terminal columns wide, so it repeats height/2 times
// per row to stay visually square next to the big letters.
func emojiBlock(emoji string, height int) []string {
	if emoji == "" || height <= 0 {
		return nil
	}
	cols := (height + 1) / 2
	if cols < 1 {
		cols = 1
	}
	row := strings.Repeat(emoji, cols)
	rows := make([]string, height)
	for i := range rows {
		rows[i] = row
	}
	return rows
}

// joinBlocks concatenates per-segment line-blocks horizontally. Blocks are
// bottom-aligned to a common height (shorter blocks get blank rows on top so
// figlet baselines line up), each row is padded to its block's display width,
// and a single-space gutter separates blocks. Color is applied by the caller.
func joinBlocks(blocks [][]string) string {
	height := 0
	for _, b := range blocks {
		if len(b) > height {
			height = len(b)
		}
	}
	if height == 0 {
		return ""
	}

	// Per-block display width, for column alignment.
	widths := make([]int, len(blocks))
	for i, b := range blocks {
		w := 0
		for _, row := range b {
			if dw := displayWidth(row); dw > w {
				w = dw
			}
		}
		widths[i] = w
	}

	var out strings.Builder
	for r := 0; r < height; r++ {
		var line strings.Builder
		for i, b := range blocks {
			if i > 0 {
				line.WriteByte(' ') // gutter
			}
			// Bottom-align: this block's row r maps to a row offset by the
			// difference between the common height and the block's own height.
			idx := r - (height - len(b))
			cell := ""
			if idx >= 0 && idx < len(b) {
				cell = b[idx]
			}
			line.WriteString(cell)
			if pad := widths[i] - displayWidth(cell); pad > 0 {
				line.WriteString(strings.Repeat(" ", pad))
			}
		}
		out.WriteString(strings.TrimRight(line.String(), " "))
		out.WriteByte('\n')
	}
	return out.String()
}

// bannerWidth returns the display width of the widest line in the banner.
func bannerWidth(banner string) int {
	max := 0
	for _, line := range strings.Split(banner, "\n") {
		if n := displayWidth(line); n > max {
			max = n
		}
	}
	if max < 8 {
		max = 8
	}
	return max
}

// displayWidth approximates the terminal column width of a line: emoji and
// other non-ASCII runes occupy ~2 columns, ASCII runes 1. Joiners/modifiers
// add no width.
func displayWidth(line string) int {
	w := 0
	for _, r := range line {
		w += runeWidth(r)
	}
	return w
}

// dottedDivider builds a "· · ·" style developer divider sized to the banner.
func dottedDivider(width int) string {
	n := width / 2
	if n < 4 {
		n = 4
	}
	return strings.TrimRight(strings.Repeat("· ", n), " ") + "\n"
}
