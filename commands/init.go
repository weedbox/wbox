package commands

import (
	"fmt"
	"go/token"
	"os"

	"github.com/rogpeppe/go-internal/modfile"
	"github.com/spf13/cobra"
	"github.com/weedbox/wbox/lib"
)

var InitCmd = &cobra.Command{
	Use:   "init [project name] [module name]",
	Short: "Initialize a new project based on weedbox template",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		moduleName := args[1]
		err := initProject(projectName, moduleName)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Project initialized successfully.")
	},
}

func initProject(projectName string, moduleName string) error {

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

	// Update go.mod
	fmt.Println("Initializing go.mod...")
	err = initGoMod(moduleName)
	if err != nil {
		return err
	}

	return nil
}

func initGoMod(moduleName string) error {

	data, err := os.ReadFile("go.mod")
	if err != nil {
		return err
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return err
	}

	if err := f.AddModuleStmt(moduleName); err != nil {
		return err
	}

	newData, err := f.Format()
	if err != nil {
		return err
	}

	if err := os.WriteFile("go.mod", newData, 0644); err != nil {
		return err
	}

	return nil
}
