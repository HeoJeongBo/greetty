package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/HeoJeongBo/greetty/internal/config"
	"github.com/HeoJeongBo/greetty/internal/shell"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up greetty: create config and hook into zsh (run once)",
	RunE: func(cmd *cobra.Command, _ []string) error {
		out := cmd.OutOrStdout()

		created, err := config.EnsureDefault()
		if err != nil {
			return err
		}
		path, _ := config.Path()
		if created {
			fmt.Fprintf(out, "📝 Created default config at %s\n", path)
		} else {
			fmt.Fprintf(out, "📝 Using existing config at %s\n", path)
		}

		hookPath, err := shell.WriteHook()
		if err != nil {
			return err
		}

		added, rc, err := shell.Install(hookPath)
		if err != nil {
			return err
		}
		if added {
			fmt.Fprintf(out, "🔗 Hooked greetty into %s\n", rc)
			fmt.Fprintln(out, "✅ Done! Restart your terminal (or run `exec zsh`) to see it.")
		} else {
			fmt.Fprintf(out, "✅ Already installed in %s — nothing to do.\n", rc)
		}
		return nil
	},
}
