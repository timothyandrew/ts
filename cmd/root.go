package cmd

import (
	"fmt"
	"os"

	cli "github.com/spf13/cobra"
)

var (
	rootCmd = &cli.Command{
		Use:   "totalspaces",
		Short: "CLI interface to `totalspaces2-api`",
	}
)

func init() {
	rootCmd.AddCommand(spacesCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
