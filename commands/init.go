package commands

import (
	"archive/zip"
	"fmt"
	"go/token"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/weedbox/cli/lib"
)

const (
	GitHubOwner  = "weedbox"
	GitHubRepo   = "template-project"
	GitHubBranch = "main"
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

	zipPath, err := downloadAndExtractRepo(GitHubOwner, GitHubRepo, GitHubBranch)
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
	err = unzip(zipPath, ".")
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

func downloadAndExtractRepo(owner, repo, branch string) (string, error) {

	fmt.Println("Downloading WeedBox template...")

	url := fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", owner, repo, branch)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download template: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch: status code %d", resp.StatusCode)
	}

	zipPath := filepath.Join(os.TempDir(), repo+".zip")
	out, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return zipPath, nil
}

// unzip extracts a ZIP file to a destination directory and removes .git-related files
func unzip(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	var baseFolder string
	for _, f := range r.File {
		// Get base folder from the first file
		if baseFolder == "" && f.FileInfo().IsDir() {
			baseFolder = f.Name
		}

		// Skip .git files
		if strings.Contains(f.Name, ".git") {
			continue
		}

		fpath := filepath.Join(dest, strings.TrimPrefix(f.Name, baseFolder))

		/*
			// Check for ZipSlip vulnerability
			if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
				return fmt.Errorf("illegal file path: %s", fpath)
			}
		*/
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create directory for file
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Extract file
		outFile, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open zip content: %w", err)
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("failed to write extracted file: %w", err)
		}
	}

	return nil
}
