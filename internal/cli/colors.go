package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HeoJeongBo/greetty/internal/config"
	"github.com/HeoJeongBo/greetty/internal/render"
)

var colorsCmd = &cobra.Command{
	Use:   "colors",
	Short: "List available colors and gradient presets",
	Long: "List the single colors and gradient presets you can set with " +
		"'greetty set color <name>'. Each name is shown in its own color.\n" +
		"Preview one with your banner via 'greetty preview --color <name>'.",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Best-effort: mark the current color. A failed load just means no marker.
		current := ""
		if cfg, err := config.Load(); err == nil {
			current = cfg.Color
		}

		out := cmd.OutOrStdout()
		marked := false
		mark := func(name string) string {
			if name == current {
				marked = true
				return " *"
			}
			return ""
		}

		fmt.Fprintln(out, "Colors:")
		for _, name := range render.NamedColors() {
			fmt.Fprintf(out, "  %s%s\n", render.Colorize(name, name), mark(name))
		}

		fmt.Fprintln(out, "\nPresets (gradients):")
		for _, name := range render.Presets() {
			fmt.Fprintf(out, "  %s%s\n", render.Colorize(name, name), mark(name))
		}

		if marked {
			fmt.Fprintln(out, "\n* = current color")
		}
		return nil
	},
}
