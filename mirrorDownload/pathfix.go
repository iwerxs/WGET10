package mirrorDownload

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// fixFilePaths updates file paths in HTML and CSS files
func fixFilePaths(filePath string) error {
	// Read the file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	text := string(content)

	// Fix paths in HTML files
	if strings.HasSuffix(filePath, ".html") {
		// Update image paths in <img> tags
		text = regexp.MustCompile(`(<img[^>]+src=["'])([^"']+)`).ReplaceAllString(text, `$1img/$2`)

		// Update CSS file paths in <link> tags
		text = regexp.MustCompile(`(<link[^>]+href=["'])([^"']+\.css)`).ReplaceAllString(text, `$1css/$2`)

		// Update inline CSS styles
		text = regexp.MustCompile(`(background-image:\s*url\(["']?)([^"')]+)`).ReplaceAllString(text, `$1img/$2`)
	}

	// Fix paths in CSS files
	if strings.HasSuffix(filePath, ".css") {
		// Update image references inside CSS
		text = regexp.MustCompile(`(url\(["']?)([^"')]+)`).ReplaceAllString(text, `$1../img/$2`)
	}

	// Write the modified content back to the file
	err = ioutil.WriteFile(filePath, []byte(text), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %v", filePath, err)
	}

	fmt.Printf("Updated file paths in %s\n", filePath)
	return nil
}

// ProcessDownloadedFiles walks through the download directory and fixes file paths
func ProcessDownloadedFiles(rootDir string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".css")) {
			return fixFilePaths(path)
		}
		return nil
	})
	return err
}
