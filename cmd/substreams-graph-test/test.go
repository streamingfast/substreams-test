package main

import (
	"github.com/spf13/cobra"
)

// runCmd represents the command to run substreams remotely
var testCmd = &cobra.Command{
	Use:          "test",
	Short:        "Test subcommand",
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(testCmd)
}
