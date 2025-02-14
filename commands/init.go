package commands

import (
	"fmt"
	"go/token"
	"os"

	"github.com/spf13/cobra"
	"github.com/weedbox/wbox/lib"
)

var InitCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new project based on weedbox template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		err := initProject(projectName)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Project initialized successfully.")
	},
}

func initProject(projectName string) error {

	GitHubOwner := "weedbox"
	GitHubRepo := "template-project"
	GitHubBranch := "main"

	zipPath, err := lib.DownloadRepo(GitHubOwner, GitHubRepo, GitHubBranch)
	defer func(zipPath string) {
		if err := os.Remove(zipPath); err != nil {
			fmt.Errorf("failed to remove file: %w", err)
		}
	}(zipPath)

	if err != nil {
		return err
	}

	fmt.Println("Initializing project ...")

	// Extract ZIP file
	err = lib.ExtractFile(zipPath, ".")
	if err != nil {
		return err
	}

	// Update values in files
	gt, err := lib.OpenGolangTemplate("main.go")
	if err != nil {
		return err
	}

	gt.SetConstValue("appName", token.STRING, projectName)
	gt.SetConstValue("appDescription", token.STRING, fmt.Sprintf("%s is a general service.", projectName))
	err = gt.Save()
	if err != nil {
		return err
	}

	return nil
}
