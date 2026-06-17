package cmd

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

func runGreet(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	fmt.Fprint(cmd.OutOrStdout(), render.Render(cfg))
	return nil
}
