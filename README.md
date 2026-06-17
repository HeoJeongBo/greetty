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
go install github.com/HeoJeongBo/greetty@latest
```

### From source

```sh
git clone https://github.com/HeoJeongBo/greetty && cd greetty
go build -o greetty .
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

| Field   | Description                          | Default                  |
| ------- | ------------------------------------ | ------------------------ |
| `text`  | The banner text                      | your login name          |
| `emoji` | Emoji shown above the banner         | `🚀`                     |
| `font`  | go-figure ASCII font                 | `slant`                  |
| `color` | Banner color                         | `cyan`                   |

Edit the file directly, or use the `set` command:

```sh
greetty set text  "ship it"
greetty set emoji 🔥
greetty set font  small
greetty set color magenta
```

**Colors:** `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`.
**Fonts:** any [go-figure](https://github.com/common-nighthawk/go-figure) font, e.g. `slant`, `standard`, `small`, `big`, `banner`, `block`. An unknown font falls back to `standard`.

## How it works

`greetty init` does **not** rewrite your `.zshrc`. It appends a single marker-guarded block:

```zsh
# >>> greetty >>>
[ -f "~/.config/greetty/init.zsh" ] && source "~/.config/greetty/init.zsh"
# <<< greetty <<<
```

Your own `.zshrc` content is left untouched, a one-time backup (`.zshrc.greetty.bak`) is made before the first change, and re-running `init` is idempotent (the block is added at most once). The sourced hook runs `greetty greet` exactly once per session, guarded by a `GREETTY_SHOWN` flag so it doesn't repeat on every prompt.

If `$ZDOTDIR` is set, greetty targets `$ZDOTDIR/.zshrc`; otherwise `~/.zshrc`.

## Commands

| Command              | Description                                                        |
| -------------------- | ------------------------------------------------------------------ |
| `greetty init`       | Create the default config and hook greetty into your shell.        |
| `greetty greet`      | Print the banner to stdout (this is what the shell hook calls).    |
| `greetty set <k> <v>`| Update a config field: `text`, `emoji`, `font`, or `color`.        |
| `greetty uninstall`  | Remove the shell hook. Your config under `~/.config/greetty` stays.|

## Uninstall

```sh
greetty uninstall
```

Removes only greetty's marker block from your `.zshrc`, leaving the rest of the file exactly as it was. Restart your terminal to take effect.

## Notes

- Currently targets **zsh**. The shell-hook layer is isolated, so bash/fish support can be added later.
- The greeting is rendered defensively (bad fonts/colors fall back to defaults, errors are swallowed) so it can never break shell startup.
