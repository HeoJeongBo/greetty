# greetty Usage

A practical guide to installing, configuring, and using greetty day to day.
For publishing/release steps, see [DEPLOY.md](DEPLOY.md).

---

## 1. Install

### Homebrew (recommended)

```sh
brew install HeoJeongBo/greetty/greetty
```

### Go

```sh
go install github.com/HeoJeongBo/greetty@latest
```

### From source

```sh
git clone https://github.com/HeoJeongBo/greetty && cd greetty
go build -o greetty .
mv greetty /usr/local/bin/      # anywhere on your PATH
```

> The shell hook calls the bare command `greetty`, so the binary **must be on
> your `PATH`**. Homebrew and `go install` handle this for you.

---

## 2. Set up (once)

```sh
greetty init
exec zsh        # or just open a new terminal
```

`greetty init`:
- creates a default config at `~/.config/greetty/config.toml` (never overwrites an existing one),
- writes the shell hook to `~/.config/greetty/init.zsh`,
- appends **one** marker-guarded line to your `.zshrc` (with a one-time `.zshrc.greetty.bak` backup),
- is safe to run twice — the block is added at most once.

Your own `.zshrc` content is never modified.

---

## 3. Customize

Edit `~/.config/greetty/config.toml` directly:

```toml
text  = "hello"
emoji = "🚀"
font  = "slant"
color = "cyan"
```

…or use the `set` command (no file editing needed):

```sh
greetty set text  "ship it"
greetty set emoji 🔥
greetty set font  small
greetty set color magenta
```

Run `greetty greet` to preview without opening a new shell.

| Field   | What it does                  | Default         |
| ------- | ----------------------------- | --------------- |
| `text`  | The banner text               | your login name |
| `emoji` | Emoji above the banner        | `🚀`            |
| `font`  | ASCII font (go-figure)        | `slant`         |
| `color` | Banner color                  | `cyan`          |

**Colors:** `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`.
**Fonts:** any go-figure font — e.g. `slant`, `standard`, `small`, `big`, `banner`, `block`. Unknown fonts fall back to `standard`.

---

## 4. Commands

| Command                | Description                                                  |
| ---------------------- | ------------------------------------------------------------ |
| `greetty init`         | Create config and hook greetty into your shell.              |
| `greetty greet`        | Print the banner (what the shell hook runs each session).    |
| `greetty set <k> <v>`  | Update a field: `text`, `emoji`, `font`, `color`.            |
| `greetty uninstall`    | Remove the shell hook. Config under `~/.config/greetty` stays.|
| `greetty --version`    | Print the version.                                           |

---

## 5. Try it safely (no changes to your real shell)

Test in a throwaway sandbox before installing for real:

```sh
TMP=$(mktemp -d)
ZDOTDIR=$TMP XDG_CONFIG_HOME=$TMP/config greetty init
ZDOTDIR=$TMP XDG_CONFIG_HOME=$TMP/config zsh -i   # banner shows once; type `exit`
rm -rf $TMP
```

This isolates both the `.zshrc` (`ZDOTDIR`) and the config dir (`XDG_CONFIG_HOME`),
so your real `~/.zshrc` and `~/.config` are untouched.

---

## 6. Uninstall

```sh
greetty uninstall
exec zsh
```

Removes only greetty's marker block from `.zshrc`; the rest of the file is left
exactly as it was.

---

## How it works

The hook is sourced from your `.zshrc` and runs the greeting **once per session**,
guarded by a `GREETTY_SHOWN` flag so it doesn't repeat on every prompt:

```zsh
if [[ -z "$GREETTY_SHOWN" ]]; then
  export GREETTY_SHOWN=1
  command greetty greet 2>/dev/null
fi
```

If `$ZDOTDIR` is set, greetty targets `$ZDOTDIR/.zshrc`; otherwise `~/.zshrc`.
The greeting is rendered defensively (bad fonts/colors fall back, errors are
swallowed) so it can never break shell startup.

### Works in VS Code

The VS Code integrated terminal is a normal zsh shell, so once greetty is
installed and `greetty init` has run, the banner appears there automatically —
no extension needed.
