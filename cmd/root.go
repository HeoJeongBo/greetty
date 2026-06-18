package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is overridable at build time:
//
//	go build -ldflags "-X github.com/HeoJeongBo/greetty/cmd.version=1.2.3"
var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "greetty",
	Short:   "A pretty developer greeting for your terminal",
	Version: version,
	// Errors are surfaced by Execute; don't dump usage on every runtime error.
	SilenceUsage:  true,
	SilenceErrors: true,
	// With no subcommand, show the greeting — handy for piping/testing.
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGreet(cmd, args)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "greetty:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd, greetCmd, setCmd, uninstallCmd, fontsCmd, previewCmd)
}
