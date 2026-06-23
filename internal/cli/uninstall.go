package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HeoJeongBo/greetty/internal/shell"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove greetty's hook from zsh (config is left in place)",
	RunE: func(cmd *cobra.Command, _ []string) error {
		out := cmd.OutOrStdout()
		removed, rc, err := shell.Uninstall()
		if err != nil {
			return err
		}
		if removed {
			fmt.Fprintf(out, "🧹 Removed greetty hook from %s\n", rc)
			fmt.Fprintln(out, "✅ Restart your terminal to take effect. Config kept at ~/.config/greetty.")
		} else {
			fmt.Fprintf(out, "Nothing to remove — no greetty hook found in %s.\n", rc)
		}
		return nil
	},
}
