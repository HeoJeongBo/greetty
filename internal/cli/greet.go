package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HeoJeongBo/greetty/internal/config"
	"github.com/HeoJeongBo/greetty/internal/render"
)

var greetCmd = &cobra.Command{
	Use:   "greet",
	Short: "Print the greeting banner (called by the shell hook)",
	RunE:  runGreet,
}

func init() {
	// One-shot overrides for the greet path (root and `greet`). They change only
	// this render and are never saved — handy for a quick `greetty --color fire`.
	for _, c := range []*cobra.Command{rootCmd, greetCmd} {
		c.Flags().StringP("color", "c", "", "color or gradient preset for this run (not saved)")
		c.Flags().StringP("font", "f", "", "font for this run (not saved)")
	}
}

func runGreet(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if c, _ := cmd.Flags().GetString("color"); c != "" {
		if !render.IsColor(c) {
			return fmt.Errorf("unknown color %q — run 'greetty colors' to see options", c)
		}
		cfg.Color = c
	}
	if f, _ := cmd.Flags().GetString("font"); f != "" {
		if !render.FontExists(f) {
			return fmt.Errorf("unknown font %q — run 'greetty fonts' to see options", f)
		}
		cfg.Font = f
	}
	fmt.Fprint(cmd.OutOrStdout(), render.Render(cfg))
	return nil
}
