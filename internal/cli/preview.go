package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HeoJeongBo/greetty/internal/config"
	"github.com/HeoJeongBo/greetty/internal/render"
)

var previewColor string

var previewCmd = &cobra.Command{
	Use:   "preview [font]",
	Short: "Render your banner with a font and/or color without saving it",
	Long: "Render your banner without saving it. Pass a font name to try a font, " +
		"and/or --color to try a color or gradient preset.\n" +
		"Examples:\n" +
		"  greetty preview small\n" +
		"  greetty preview --color rainbow\n" +
		"  greetty preview small --color fire",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(args) == 1 {
			font := args[0]
			if !render.FontExists(font) {
				return fmt.Errorf("unknown font %q — run 'greetty fonts' to see options", font)
			}
			cfg.Font = font // override only for this render; nothing is saved
		}
		if previewColor != "" {
			if !render.IsColor(previewColor) {
				return fmt.Errorf("unknown color %q — run 'greetty colors' to see options", previewColor)
			}
			cfg.Color = previewColor
		}

		fmt.Fprint(cmd.OutOrStdout(), render.Render(cfg))
		return nil
	},
}

func init() {
	previewCmd.Flags().StringVarP(&previewColor, "color", "c", "", "color or gradient preset to preview")
}
