package commands

import (
	"fmt"
	"go/token"
	"os"

	"github.com/spf13/cobra"
	"github.com/weedbox/cli/lib"
)

var InitModuleCmd = &cobra.Command{
	Use:   "init-module [module-name]",
	Short: "Initialize a new module based on weedbox module template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		err := initModule(moduleName)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Module initialized successfully.")
	},
}

func initModule(moduleName string) error {

	GitHubOwner := "weedbox"
	GitHubRepo := "template-module"
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

	fmt.Println("Initializing module...")

	// Extract ZIP file
	err = lib.ExtractFile(zipPath, ".")
	if err != nil {
		return err
	}

	// Update values in files
	gt, err := lib.OpenGolangTemplate("module.go")
	if err != nil {
		return err
	}

	gt.SetConstValue("ModuleName", token.STRING, moduleName)
	gt.RenameType("TemplateModule", moduleName)
	gt.RenameReceiver("TemplateModule", moduleName)
	gt.RenameVariableType("TemplateModule", moduleName)
	gt.RenameAllocationType("TemplateModule", moduleName)
	gt.RenameFunctionResult("TemplateModule", moduleName)
	gt.RenameFunctionResultInCallExpr("TemplateModule", moduleName)

	err = gt.Save()
	if err != nil {
		return err
	}

	return nil
}
