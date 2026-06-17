// Package config loads, creates, and persists greetty's user configuration.
package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds the greeting content and presentation options.
type Config struct {
	Text  string `toml:"text"`  // banner text
	Emoji string `toml:"emoji"` // emoji shown above the banner
	Font  string `toml:"font"`  // go-figure font name
	Color string `toml:"color"` // banner color name
}

// Defaults returns a Config with sensible values. Text defaults to the current
// login name so a fresh install already shows something personal.
func Defaults() Config {
	text := "hello"
	if u, err := user.Current(); err == nil && u.Username != "" {
		text = u.Username
	}
	return Config{
		Text:  text,
		Emoji: "🚀",
		Font:  "slant",
		Color: "cyan",
	}
}

// Dir returns the greetty config directory. It honors $XDG_CONFIG_HOME and
// otherwise uses ~/.config/greetty on every platform, so the path stays the
// familiar XDG-style location (rather than ~/Library/Application Support on
// macOS).
func Dir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "greetty"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "greetty"), nil
}

// Path returns the full path to config.toml.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

// Load reads config.toml, falling back to defaults for any missing field. A
// missing file is not an error — defaults are returned so `greetty greet`
// always produces output.
func Load() (Config, error) {
	cfg := Defaults()
	path, err := Path()
	if err != nil {
		return cfg, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, fmt.Errorf("reading %s: %w", path, err)
	}
	if cfg.Text == "" {
		cfg.Text = Defaults().Text
	}
	if cfg.Font == "" {
		cfg.Font = Defaults().Font
	}
	if cfg.Color == "" {
		cfg.Color = Defaults().Color
	}
	return cfg, nil
}

// Save writes the config to config.toml, creating the directory if needed.
func Save(cfg Config) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "config.toml")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

// EnsureDefault writes a default config.toml only if none exists yet. It never
// overwrites an existing file. Returns true if it created the file.
func EnsureDefault() (bool, error) {
	path, err := Path()
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(path); err == nil {
		return false, nil // already exists, leave it alone
	}
	return true, Save(Defaults())
}
