package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HeoJeongBo/greetty/internal/config"
	"github.com/HeoJeongBo/greetty/internal/render"
)

var setCmd = &cobra.Command{
	Use:   "set <field> <value>",
	Short: "Update a config field (text, emoji, font, color)",
	Long:  "Update a config field without editing config.toml by hand.\nFields: text, emoji, font, color.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		field, value := args[0], args[1]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		switch field {
		case "text":
			cfg.Text = value
		case "emoji":
			cfg.Emoji = value
		case "font":
			if !render.FontExists(value) {
				return fmt.Errorf("unknown font %q — run 'greetty fonts' to see options", value)
			}
			cfg.Font = value
		case "color":
			cfg.Color = value
		default:
			return fmt.Errorf("unknown field %q (use: text, emoji, font, color)", field)
		}

		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "✅ Set %s = %q\n", field, value)
		return nil
	},
}
