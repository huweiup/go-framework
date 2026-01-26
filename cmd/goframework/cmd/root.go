package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goframework",
	Short: "A CLI tool for generating Go projects",
	Long:  `goframework is a CLI tool that helps you generate a standard Go project structure based on the go-framework library.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
