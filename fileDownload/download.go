package fileDownload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Start handles the file download
func Start(url, output, saveDir string) {
	startTime := time.Now()
	fmt.Println("Start time:", startTime.Format("2006-01-02 15:04:05"))

	// Send HTTP request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("HTTP Response:", resp.Status)

	// Get filename from URL if not provided
	if output == "" {
		output = filepath.Base(url)
	}

	// Handle `-P` flag and expand `~`
	if saveDir != "" {
		if strings.HasPrefix(saveDir, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				return
			}
			saveDir = filepath.Join(homeDir, saveDir[1:])
		}

		// Create directory if it does not exist
		err := os.MkdirAll(saveDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}

		output = filepath.Join(saveDir, output)
	}

	// Create file
	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Display content length
	contentLength := resp.ContentLength
	fmt.Printf("Content Length: %.2f MB (%d bytes)\n", float64(contentLength)/(1024*1024), contentLength)

	// Download file with progress
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	endTime := time.Now()
	fmt.Println("End time:", endTime.Format("2006-01-02 15:04:05"))
	fmt.Println("File saved as:", output)
	fmt.Printf("Time taken: %.2f seconds\n", endTime.Sub(startTime).Seconds())
}
