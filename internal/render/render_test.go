package render

import (
	"regexp"
	"strings"
	"testing"

	"github.com/fatih/color"

	"github.com/HeoJeongBo/greetty/internal/config"
)

var ansiRE = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string { return ansiRE.ReplaceAllString(s, "") }

func TestRenderPureASCIIUnchanged(t *testing.T) {
	out := stripANSI(Render(config.Config{Text: "heo", Font: "slant", Color: "cyan"}))
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected non-empty banner")
	}
	if strings.Contains(out, "?") {
		t.Errorf("ascii banner should not contain '?':\n%s", out)
	}
	if !strings.Contains(out, "·") {
		t.Errorf("expected dotted divider:\n%s", out)
	}
}

func TestRenderRainbow(t *testing.T) {
	// go test's stdout is not a TTY, so color.NoColor is true and the rainbow
	// guard would no-op. Force color on for this test and restore after.
	prev := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = prev }()

	out := Render(config.Config{Text: "heo", Font: "slant", Color: "rainbow"})
	if !strings.Contains(out, "\x1b[38;2;") {
		t.Errorf("rainbow should emit truecolor escapes:\n%q", out)
	}

	stripped := stripANSI(out)
	if strings.TrimSpace(stripped) == "" {
		t.Fatal("expected non-empty banner")
	}
	if strings.Contains(stripped, "?") {
		t.Errorf("rainbow banner should not contain '?':\n%s", stripped)
	}
	if !strings.Contains(stripped, "·") {
		t.Errorf("expected dotted divider:\n%s", stripped)
	}
}

func TestRenderPresets(t *testing.T) {
	// go test's stdout is not a TTY, so color.NoColor is true and gradients
	// would no-op. Force color on for this test and restore after.
	prev := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = prev }()

	for _, name := range Presets() {
		out := Render(config.Config{Text: "heo", Font: "slant", Color: name})
		if !strings.Contains(out, "\x1b[38;2;") {
			t.Errorf("preset %q should emit truecolor escapes:\n%q", name, out)
		}
		stripped := stripANSI(out)
		if strings.TrimSpace(stripped) == "" {
			t.Errorf("preset %q produced an empty banner", name)
		}
		if strings.Contains(stripped, "?") {
			t.Errorf("preset %q banner should not contain '?':\n%s", name, stripped)
		}
		if !strings.Contains(stripped, "·") {
			t.Errorf("preset %q should keep the dotted divider:\n%s", name, stripped)
		}
	}
}

func TestIsColor(t *testing.T) {
	for _, name := range []string{"cyan", "CYAN", "red", "rainbow", "Fire", "ocean"} {
		if !IsColor(name) {
			t.Errorf("IsColor(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"", "bogus", "neonish", "purple"} {
		if IsColor(name) {
			t.Errorf("IsColor(%q) = true, want false", name)
		}
	}
}

func TestColorize(t *testing.T) {
	prev := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = prev }()

	// A gradient preset emits truecolor escapes around the text.
	if got := Colorize("rainbow", "hi"); !strings.Contains(got, "\x1b[38;2;") {
		t.Errorf("Colorize gradient should emit truecolor escapes: %q", got)
	}
	// A single named color still wraps the text in an ANSI escape.
	if got := Colorize("red", "hi"); !strings.Contains(got, "\x1b[") {
		t.Errorf("Colorize named color should emit an ANSI escape: %q", got)
	}
	// An unknown name is returned unchanged.
	if got := Colorize("bogus", "hi"); got != "hi" {
		t.Errorf("Colorize(unknown) = %q, want %q", got, "hi")
	}
	// Either way the visible text is preserved.
	if got := stripANSI(Colorize("ocean", "hi")); got != "hi" {
		t.Errorf("Colorize should not alter the text: %q", got)
	}
}

func TestRenderEmojiOnlyMapped(t *testing.T) {
	// 🚀 is mapped, so it renders as multi-line ASCII art (not the glyph).
	out := stripANSI(Render(config.Config{Text: "🚀", Font: "slant", Color: "cyan"}))
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected non-empty banner")
	}
	if strings.Contains(out, "?") {
		t.Errorf("emoji should not degrade to '?':\n%s", out)
	}
	if !strings.Contains(out, "/\\") {
		t.Errorf("expected rocket ASCII art:\n%s", out)
	}
}

func TestRenderEmojiOnlyFallback(t *testing.T) {
	// 🟦 is NOT mapped, so it falls back to a repeated-glyph block.
	out := stripANSI(Render(config.Config{Text: "🟦", Font: "slant", Color: "cyan"}))
	if !strings.Contains(out, "🟦🟦") {
		t.Errorf("unmapped emoji should repeat as a glyph block:\n%s", out)
	}
}

func TestRenderMixed(t *testing.T) {
	out := stripANSI(Render(config.Config{Text: "heo 🚀", Font: "slant", Color: "cyan"}))
	// figlet slant uses '/' and '_' heavily — a proxy for "the letters rendered".
	if !strings.Contains(out, "/") {
		t.Errorf("expected ASCII art for letters in mixed banner:\n%s", out)
	}
	// rocket art tail.
	if !strings.Contains(out, "*  *") {
		t.Errorf("expected rocket ASCII art next to letters:\n%s", out)
	}
	if strings.Contains(out, "?") {
		t.Errorf("mixed banner should not contain '?':\n%s", out)
	}
}

func TestRenderNeverPanicsOnEmoji(t *testing.T) {
	// Regression guard: pre-fix, emoji in text reached go-figure with strict=true
	// and called log.Fatal during shell startup. This simply must return.
	for _, txt := range []string{"🚀", "heo 🚀", "👨‍💻", "🇰🇷 heo", "✈️"} {
		out := Render(config.Config{Text: txt, Font: "slant", Color: "cyan"})
		if strings.TrimSpace(out) == "" {
			t.Errorf("Render(%q) returned blank", txt)
		}
	}
}

func TestRenderBogusFontNeverPanics(t *testing.T) {
	// An unknown font would panic inside go-figure (newFont→MustAsset); the
	// buildBanner guard must substitute a safe font instead.
	out := stripANSI(Render(config.Config{Text: "heo", Font: "totally-not-a-font", Color: "cyan"}))
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected non-empty banner with a bogus font")
	}
	if strings.Contains(out, "?") {
		t.Errorf("bogus font should fall back cleanly, not emit '?':\n%s", out)
	}
}

func TestSegments(t *testing.T) {
	tests := []struct {
		in   string
		want []segment
	}{
		{"heo", []segment{{"heo", false}}},
		{"🚀", []segment{{"🚀", true}}},
		{"heo 🚀", []segment{{"heo ", false}, {"🚀", true}}},
		{"👨‍💻", []segment{{"👨‍💻", true}}},   // ZWJ stays one segment
		{"✈️", []segment{{"✈️", true}}},     // VS16 stays attached
		{"🇰🇷", []segment{{"🇰🇷", true}}},     // flag = two regional indicators
		{"🚀heo🔥", []segment{{"🚀", true}, {"heo", false}, {"🔥", true}}},
	}
	for _, tt := range tests {
		got := segments(tt.in)
		if len(got) != len(tt.want) {
			t.Errorf("segments(%q) = %d segs %v, want %d %v", tt.in, len(got), got, len(tt.want), tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("segments(%q)[%d] = %+v, want %+v", tt.in, i, got[i], tt.want[i])
			}
		}
	}
}

func TestEmojiBlock(t *testing.T) {
	blk := emojiBlock("🚀", 6)
	if len(blk) != 6 {
		t.Fatalf("emojiBlock height = %d, want 6", len(blk))
	}
	for i, row := range blk {
		if !strings.Contains(row, "🚀") {
			t.Errorf("row %d missing emoji: %q", i, row)
		}
	}
	if emojiBlock("", 6) != nil {
		t.Error("empty emoji should yield nil block")
	}
	if emojiBlock("🚀", 0) != nil {
		t.Error("zero height should yield nil block")
	}
}

func TestJoinBlocks(t *testing.T) {
	blocks := [][]string{
		{"AA", "AA", "AA"},
		{"B", "B"}, // shorter — bottom-aligned, top-padded
	}
	out := joinBlocks(blocks)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("joinBlocks rows = %d, want 3:\n%s", len(lines), out)
	}
	// First row: only the taller block has content; "B" appears on rows 2-3.
	if strings.Contains(lines[0], "B") {
		t.Errorf("row 0 should not contain shorter block (bottom-align):\n%s", out)
	}
	if !strings.Contains(lines[2], "B") {
		t.Errorf("row 2 should contain shorter block:\n%s", out)
	}
}

func TestEmojiArtMapping(t *testing.T) {
	// A mapped emoji renders its hand-drawn art, not a repeated glyph block.
	out := stripANSI(Render(config.Config{Text: "🚀", Font: "slant", Color: "cyan"}))
	if strings.Contains(out, "🚀🚀") {
		t.Errorf("mapped emoji should use ASCII art, not repeated glyphs:\n%s", out)
	}
	if !strings.Contains(out, "/\\") {
		t.Errorf("expected rocket ASCII art:\n%s", out)
	}
}

func TestArtForNormalization(t *testing.T) {
	// ❤ (bare) and ❤️ (with VS16) must both resolve to the heart art.
	if _, ok := artFor("❤"); !ok {
		t.Error("bare heart should be mapped")
	}
	if _, ok := artFor("❤️"); !ok {
		t.Error("heart + VS16 should normalize to the mapped heart")
	}
	if _, ok := artFor("🟦"); ok {
		t.Error("unmapped emoji should report not found")
	}
}

func TestSplitEmoji(t *testing.T) {
	got := splitEmoji("🚀🔥")
	if len(got) != 2 || got[0] != "🚀" || got[1] != "🔥" {
		t.Errorf("splitEmoji(🚀🔥) = %v, want [🚀 🔥]", got)
	}
	zwj := splitEmoji("👨‍💻")
	if len(zwj) != 1 {
		t.Errorf("splitEmoji(👨‍💻) = %v, want a single glyph", zwj)
	}
}

func TestDisplayWidth(t *testing.T) {
	if got := displayWidth("ab"); got != 2 {
		t.Errorf("displayWidth(ascii) = %d, want 2", got)
	}
	if got := displayWidth("🚀"); got != 2 {
		t.Errorf("displayWidth(emoji) = %d, want 2", got)
	}
	if got := displayWidth("a🚀"); got != 3 {
		t.Errorf("displayWidth(mixed) = %d, want 3", got)
	}
}
