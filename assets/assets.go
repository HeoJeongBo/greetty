// Package assets embeds static files shipped with greetty.
package assets

import _ "embed"

// InitZsh is the zsh hook snippet written to the config dir by `greetty init`.
//
//go:embed init.zsh
var InitZsh string
