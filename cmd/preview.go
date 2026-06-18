package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HeoJeongBo/greetty/internal/config"
	"github.com/HeoJeongBo/greetty/internal/render"
)

var previewCmd = &cobra.Command{
	Use:   "preview <font>",
	Short: "Render your banner with a font without saving it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		font := args[0]
		if !render.FontExists(font) {
			return fmt.Errorf("unknown font %q — run 'greetty fonts' to see options", font)
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.Font = font // override only the font; nothing is saved

		fmt.Fprint(cmd.OutOrStdout(), render.Render(cfg))
		return nil
	},
}
