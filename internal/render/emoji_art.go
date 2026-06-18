package render

import "strings"

// emojiArt maps common emoji to small hand-drawn ASCII-art blocks. Lookup is by
// the emoji's base form (variation selectors, ZWJ, and skin-tone modifiers are
// stripped first via normalizeEmoji), so "❤️" and "❤" both resolve here.
//
// Each entry is a slice of equal-length lines. Keep them compact (a handful of
// rows) so they sit nicely next to figlet letters; joinBlocks bottom-aligns and
// pads them. Emoji not listed here fall back to a repeated-glyph block.
var emojiArt = map[string][]string{
	"🚀": {
		`   /\`,
		`  |==|`,
		`  |  |`,
		` /|  |\`,
		`/_|__|_\`,
		`  *  *`,
	},
	"🔥": {
		`  )`,
		` ((`,
		` )\)`,
		`(  ( \`,
		`( () )`,
		` \__/`,
	},
	"⭐": {
		`   __`,
		`  /  \`,
		`-< () >-`,
		`  \__/`,
		`  /  \`,
	},
	"✨": {
		` .  *  .`,
		`  \ | /`,
		`* -- + --*`,
		`  / | \`,
		` '  *  '`,
	},
	"❤": {
		` _  _`,
		`( \/ )`,
		` \  /`,
		`  \/`,
	},
	"☕": {
		` ( (`,
		`  ) )`,
		` ______`,
		`|      |]`,
		`\      /`,
		` \____/`,
	},
	"🍺": {
		` _____`,
		`|~~~~~|_`,
		`|booze| |`,
		`|booze| |`,
		`|_____|_|`,
	},
	"🎉": {
		`  \o/`,
		`  .|.`,
		` ./ \.`,
		`* . ' .`,
		`. ' * '`,
	},
	"👍": {
		`   _`,
		`  | |`,
		` _| |__`,
		`|      |`,
		`|______|`,
	},
	"💀": {
		`  ____`,
		` /    \`,
		`| o  o |`,
		`|  ||  |`,
		` \ -- /`,
		`  |||| `,
	},
	"🌙": {
		`   _`,
		`  / )`,
		` ( (`,
		`  \ \`,
		`   \_)`,
	},
	"☀": {
		` \ | /`,
		`  .-.`,
		`-( O )-`,
		`  '-'`,
		` / | \`,
	},
	"🐱": {
		` /\_/\`,
		`( o.o )`,
		` > ^ <`,
	},
	"🐶": {
		` / \__`,
		`(    @\__`,
		` /         O`,
		`/   (_____/`,
		`/_____/   U`,
	},
	"🎸": {
		`  ___`,
		` ||_|"|`,
		` ||_| |`,
		` (.   )`,
		`  '--'`,
	},
	"💻": {
		` ________`,
		`|  ____  |`,
		`| |    | |`,
		`| |____| |`,
		`|________|`,
		` \______/`,
	},
	"🍕": {
		`  /\`,
		` /  \`,
		`/ o  \`,
		`\  o /`,
		` \  /`,
		`  \/`,
	},
	"🌈": {
		`   ___`,
		`  / _ \`,
		` / / \ \`,
		`/ /   \ \`,
	},
	"⚡": {
		`   /`,
		`  /`,
		` /__`,
		`   /`,
		`  /`,
	},
	"💡": {
		`  ___`,
		` /   \`,
		`( bulb )`,
		` \   /`,
		`  | |`,
		`  |_|`,
	},
	"🎯": {
		`  ___`,
		` / _ \`,
		`( (o) )`,
		` \ _ /`,
		`  '-'`,
	},
	"🏆": {
		` \    /`,
		`  \__/`,
		`  |  |`,
		` _|__|_`,
		`|______|`,
	},
	"😎": {
		`  ____`,
		` / oo \`,
		`| -__- |`,
		` \ -- /`,
		`  '--'`,
	},
	"😀": {
		`  ____`,
		` / ^^ \`,
		`| o  o |`,
		` \ \_/ /`,
		`  '--'`,
	},
	"🤖": {
		` [o_o]`,
		`/|___|\`,
		` | | |`,
		` d   b`,
	},
	"💩": {
		`  ___`,
		` /~~~\`,
		`( o o )`,
		`(_____)`,
		`(_____)`,
	},
	"🎂": {
		`  i i i`,
		` |||||||`,
		`|~~~~~~~|`,
		`|  cake |`,
		`|_______|`,
	},
	"🍀": {
		`  _ _`,
		` ( | )`,
		`(__|__)`,
		`  /|`,
		`   |`,
	},
	"🌊": {
		`  __    __`,
		` /  \  /  \`,
		`/    \/    \`,
		`~~~~~~~~~~~~`,
	},
	"🦄": {
		`   /)`,
		`  //  ___`,
		` //  / o \`,
		`(((  \___/`,
		` \\\___|`,
	},
	"👀": {
		` ___  ___`,
		`/o  \/  o\`,
		`\___/\___/`,
	},
	"✅": {
		`      _`,
		`     //`,
		` _  //`,
		`( \//`,
		` \/`,
	},
	"❌": {
		`\\   //`,
		` \\ //`,
		`  \X/`,
		` // \\`,
		`//   \\`,
	},
}

// normalizeEmoji strips variation selectors, zero-width joiners, and skin-tone
// modifiers so that styled variants resolve to their base glyph in emojiArt.
func normalizeEmoji(emoji string) string {
	var b strings.Builder
	for _, r := range emoji {
		if isJoiner(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// artFor returns the ASCII-art block for an emoji if one is registered, padded
// to equal-width lines. The bool reports whether a mapping was found.
func artFor(emoji string) ([]string, bool) {
	art, ok := emojiArt[normalizeEmoji(emoji)]
	if !ok {
		return nil, false
	}
	// Defensive copy with equal-width padding so joinBlocks aligns columns.
	width := 0
	for _, line := range art {
		if w := displayWidth(line); w > width {
			width = w
		}
	}
	out := make([]string, len(art))
	for i, line := range art {
		if pad := width - displayWidth(line); pad > 0 {
			line += strings.Repeat(" ", pad)
		}
		out[i] = line
	}
	return out, true
}
