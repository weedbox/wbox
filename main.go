package main

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
	"github.com/weedbox/cli/commands"
)

const (
	appName        = "wbox"
	appDescription = "wbox is a command line tool to manage weedbox project"
)

// Command options
var verbose bool

func main() {

	rootCmd := &cobra.Command{
		Use:  appName,
		Long: appDescription,
	}

	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Display detailed logs")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
