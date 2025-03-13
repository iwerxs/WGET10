package mirrorDownload

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"wget/downloader" // Handles downloading resources
)

// Start begins mirroring a website
func Start(siteURL string, convertLinks bool, rejectExtensions []string, excludeDirs []string) {
	startTime := time.Now()
	fmt.Printf("Start time: %s\n", startTime.Format("2006-01-02 15:04:05"))

	parsedURL, err := url.Parse(siteURL)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return
	}

	// Define save directory for the mirrored site
	domain := parsedURL.Hostname()
	saveDir := filepath.Join("mirrored_sites", domain)
	err = os.MkdirAll(saveDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	fmt.Println("Mirroring:", siteURL)

	// Fetch the HTML content
	resp, err := http.Get(siteURL)
	if err != nil {
		fmt.Println("Error fetching site:", err)
		return
	}
	defer resp.Body.Close()

	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading HTML content:", err)
		return
	}

	// Download resources (CSS, images, JS, etc.)
	err = downloader.DownloadResources(string(htmlContent), siteURL, saveDir, excludeDirs)
	if err != nil {
		fmt.Println("Error downloading resources:", err)
		return
	}

	// Fix file paths inside HTML and CSS files
	err = ProcessDownloadedFiles(saveDir)
	if err != nil {
		fmt.Println("Error updating file paths:", err)
	}

	endTime := time.Now()
	fmt.Printf("End time: %s\n", endTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Time taken: %.2f seconds\n", endTime.Sub(startTime).Seconds())
	fmt.Println("Website successfully mirrored to:", saveDir)
}
