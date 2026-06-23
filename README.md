# greetty

A pretty, developer-flavored greeting for your terminal. Set your text and emoji once, and `greetty` renders them as a colorful ASCII banner every time you open a new shell.

```
  🚀
    __                      _                                   __
   / /_   ___   ____       (_)  ___   ____    ____    ____ _   / /_   ____
  / __ \ / _ \ / __ \     / /  / _ \ / __ \  / __ \  / __ `/  / __ \ / __ \
 / / / //  __// /_/ /    / /  /  __// /_/ / / / / / / /_/ /  / /_/ // /_/ /
/_/ /_/ \___/ \____/  __/ /   \___/ \____/ /_/ /_/  \__, /  /_.___/ \____/
                     /___/                         /____/
· · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · ·
```

## Why

Most shell greeters make you edit your `.zshrc` by hand. `greetty` doesn't. Run `greetty init` once and it wires itself into your shell — no manual editing, sensible defaults, and a clean uninstall.

## Install

### Homebrew (recommended)

```sh
brew install HeoJeongBo/greetty/greetty
```

This taps `HeoJeongBo/homebrew-greetty` and installs the binary onto your `PATH`.

### Go

```sh
go install github.com/HeoJeongBo/greetty/cmd/greetty@latest
```

### From source

```sh
git clone https://github.com/HeoJeongBo/greetty && cd greetty
go build -o greetty ./cmd/greetty
mv greetty /usr/local/bin/   # or anywhere on your PATH
```

## Quick start

```sh
greetty init      # creates config + hooks into your shell (run once)
exec zsh          # or just open a new terminal
```

That's it. The next interactive shell will print your banner.

## Configuration

`greetty init` writes a default config to `~/.config/greetty/config.toml` (it respects `$XDG_CONFIG_HOME`). It **never** overwrites an existing config.

```toml
text  = "hello"
emoji = "🚀"
font  = "slant"
color = "cyan"
```

| Field   | Description                                          | Default                  |
| ------- | --------------------------------------------------- | ------------------------ |
| `text`  | The banner text (emoji become ASCII art — see below) | your login name          |
| `emoji` | Emoji shown above the banner                         | `🚀`                     |
| `font`  | go-figure ASCII font                                 | `slant`                  |
| `color` | Banner color                                         | `cyan`                   |

Edit the file directly, or use the `set` command:

```sh
greetty set text  "ship it"
greetty set emoji 🔥
greetty set font  small
greetty set color magenta
```

**Colors:** single colors `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`.
**Gradient presets:** `rainbow`, `fire`, `ocean`, `sunset`, `forest`, `neon`, `mono`. A preset paints the banner as a diagonal 24-bit gradient (lolcat-style) and needs a truecolor-capable terminal (e.g. iTerm2); it falls back to plain text when color is disabled (`NO_COLOR` or piped output). Run `greetty colors` to see every color and preset rendered in its own colors, and `greetty preview --color <name>` to try one on your banner before committing. An unknown color is rejected by `set`.

```sh
greetty colors                       # list colors + presets
greetty --color rainbow              # one-shot: render now, don't save (-c/-f also work)
greetty preview --color rainbow      # try a preset (current font)
greetty preview small --color fire   # try a font and a preset together
greetty set color sunset             # save it
```

**Fonts:** any [go-figure](https://github.com/common-nighthawk/go-figure) font (149 in total). Run `greetty fonts` to list them and `greetty preview <font>` to try one before committing. An unknown font is rejected by `set` and safely falls back to `standard` at render time.

## Emoji in the banner text

Put an emoji directly in `text` and greetty renders it as **ASCII art** sized to sit next to the big letters:

```sh
greetty set text "heo 🚀"
greetty greet
```

```
                         /\
    __                  |==|
   / /_   ___   ____    |  |
  / __ \ / _ \ / __ \  /|  |\
 / / / //  __// /_/ / /_|__|_\
/_/ /_/ \___/ \____/    *  *
· · · · · · · · · · · · · · ·
```

Common emoji (🚀 🔥 ⭐ ❤️ ☕ 🐱 💻 🎉 …) have hand-drawn art; any other emoji falls back to a large repeated-glyph block, so nothing ever breaks. This is separate from the `emoji` field, which still prints above the banner.

## How it works

`greetty init` does **not** rewrite your `.zshrc`. It appends a single marker-guarded block:

```zsh
# >>> greetty >>>
[ -f "$HOME/.config/greetty/init.zsh" ] && source "$HOME/.config/greetty/init.zsh"
# <<< greetty <<<
```

(greetty writes the **absolute** path it resolved, so the line works regardless of how your shell expands `~`.)

Your own `.zshrc` content is left untouched, a one-time backup (`.zshrc.greetty.bak`) is made before the first change, and re-running `init` is idempotent (the block is added at most once). The sourced hook runs `greetty greet` exactly once per session, guarded by a `GREETTY_SHOWN` flag so it doesn't repeat on every prompt.

If `$ZDOTDIR` is set, greetty targets `$ZDOTDIR/.zshrc`; otherwise `~/.zshrc`.

## Commands

| Command              | Description                                                        |
| -------------------- | ------------------------------------------------------------------ |
| `greetty init`       | Create the default config and hook greetty into your shell.        |
| `greetty greet`      | Print the banner to stdout (the shell hook calls this). `--color/-c` and `--font/-f` override just this run without saving; the same flags work on bare `greetty`. |
| `greetty set <k> <v>`| Update a config field: `text`, `emoji`, `font`, or `color`.        |
| `greetty fonts`      | List all available banner fonts (the current one is marked `*`).   |
| `greetty colors`     | List colors and gradient presets (the current one is marked `*`).  |
| `greetty preview [font] [--color <c>]`| Render your banner with a font and/or color without saving it. |
| `greetty uninstall`  | Remove the shell hook. Your config under `~/.config/greetty` stays.|

## Uninstall

```sh
greetty uninstall
```

Removes only greetty's marker block from your `.zshrc`, leaving the rest of the file exactly as it was. Restart your terminal to take effect.

## Notes

- Currently targets **zsh**. The shell-hook layer is isolated, so bash/fish support can be added later.
- The greeting is rendered defensively (bad fonts/colors fall back to defaults, errors are swallowed) so it can never break shell startup.
