package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/HeoJeongBo/greetty/internal/config"
	"github.com/HeoJeongBo/greetty/internal/render"
)

const fontColumns = 4

var fontsCmd = &cobra.Command{
	Use:   "fonts",
	Short: "List available banner fonts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		names := render.Fonts()

		// Best-effort: mark the current font. A failed load just means no marker.
		current := ""
		if cfg, err := config.Load(); err == nil {
			current = cfg.Font
		}

		out := cmd.OutOrStdout()
		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		marked := false
		for i, name := range names {
			label := name
			if name == current {
				label += " *"
				marked = true
			}
			fmt.Fprint(w, label)
			if (i+1)%fontColumns == 0 {
				fmt.Fprintln(w)
			} else {
				fmt.Fprint(w, "\t")
			}
		}
		if len(names)%fontColumns != 0 {
			fmt.Fprintln(w)
		}
		w.Flush()

		if marked {
			fmt.Fprintln(out, "\n* = current font")
		}
		return nil
	},
}
