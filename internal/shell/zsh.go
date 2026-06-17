// Package shell installs and removes greetty's startup hook in the user's
// .zshrc using a marker-guarded block, so the user's own content is never
// touched.
package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HeoJeongBo/greetty/assets"
	"github.com/HeoJeongBo/greetty/internal/config"
)

const (
	markerStart = "# >>> greetty >>>"
	markerEnd   = "# <<< greetty <<<"
)

// Rcfile returns the path to the active .zshrc, honoring $ZDOTDIR.
func Rcfile() (string, error) {
	if zd := os.Getenv("ZDOTDIR"); zd != "" {
		return filepath.Join(zd, ".zshrc"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".zshrc"), nil
}

// HookPath returns the path of the sourced hook file inside the config dir.
func HookPath() (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "init.zsh"), nil
}

// WriteHook writes the embedded init.zsh into the config directory.
func WriteHook() (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "init.zsh")
	if err := os.WriteFile(path, []byte(assets.InitZsh), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// block returns the marker-guarded snippet to inject into .zshrc.
func block(hookPath string) string {
	return fmt.Sprintf("%s\n[ -f %q ] && source %q\n%s\n",
		markerStart, hookPath, hookPath, markerEnd)
}

// Install ensures the marker block is present in .zshrc exactly once. It backs
// up the file once before its first modification. Returns (added, rcfile, err)
// where added is false if the block was already present.
func Install(hookPath string) (bool, string, error) {
	rc, err := Rcfile()
	if err != nil {
		return false, "", err
	}

	content := ""
	if data, err := os.ReadFile(rc); err == nil {
		content = string(data)
	} else if !os.IsNotExist(err) {
		return false, rc, err
	}

	if strings.Contains(content, markerStart) {
		return false, rc, nil // already installed
	}

	if content != "" {
		if err := os.WriteFile(rc+".greetty.bak", []byte(content), 0o644); err != nil {
			return false, rc, err
		}
	}

	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += "\n" + block(hookPath)

	if err := os.WriteFile(rc, []byte(content), 0o644); err != nil {
		return false, rc, err
	}
	return true, rc, nil
}

// Uninstall removes the marker block from .zshrc, leaving all other lines
// intact. Returns (removed, rcfile, err).
func Uninstall() (bool, string, error) {
	rc, err := Rcfile()
	if err != nil {
		return false, "", err
	}
	data, err := os.ReadFile(rc)
	if err != nil {
		if os.IsNotExist(err) {
			return false, rc, nil
		}
		return false, rc, err
	}
	content := string(data)
	if !strings.Contains(content, markerStart) {
		return false, rc, nil
	}

	lines := strings.Split(content, "\n")
	var out []string
	inBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == markerStart {
			inBlock = true
			continue
		}
		if trimmed == markerEnd {
			inBlock = false
			continue
		}
		if inBlock {
			continue
		}
		out = append(out, line)
	}

	cleaned := strings.TrimRight(strings.Join(out, "\n"), "\n") + "\n"
	if err := os.WriteFile(rc, []byte(cleaned), 0o644); err != nil {
		return false, rc, err
	}
	return true, rc, nil
}
