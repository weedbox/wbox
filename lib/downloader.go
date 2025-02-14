package lib

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadRepo(owner, repo, branch string) (string, error) {

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

func ExtractFile(zipPath, dest string) error {
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
