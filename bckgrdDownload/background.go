package bckgrdDownload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

// Start handles the background download and logs the process to wget-log
func Start(url string) error {
	// Open log file to append logs
	logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return err
	}
	defer logFile.Close()

	// Set up logging to file
	log.SetOutput(logFile)

	// Start download and log details
	startTime := time.Now()
	logMessage(fmt.Sprintf("Start at %s", startTime.Format(time.RFC1123)))
	logMessage(fmt.Sprintf("Sending request to download %s...", url))

	// Make the HTTP request to the provided URL
	resp, err := http.Get(url)
	if err != nil {
		logMessage(fmt.Sprintf("Error: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logMessage(fmt.Sprintf("Failed to download HTTP request: %s, Status Code: %d", url, resp.StatusCode))
		return fmt.Errorf("HTTP request failed: %s: %d", url, resp.StatusCode)
	}

	// Determine the file path to save the content
	fileName := path.Base(url)
	file, err := os.Create(fileName)
	if err != nil {
		logMessage(fmt.Sprintf("Error creating file: %v", err))
		return err
	}
	defer file.Close()

	// Copy content to the file
	bytesWritten, err := io.Copy(file, resp.Body)
	if err != nil {
		logMessage(fmt.Sprintf("Error writing to file: %v", err))
		return err
	}

	// Log the download size and completion details
	fileSizeMB := float64(bytesWritten) / (1024 * 1024)
	logMessage(fmt.Sprintf("Content size: %d bytes [%.fMB]", bytesWritten, fileSizeMB))
	logMessage(fmt.Sprintf("Saving file to: ./%s", fileName))
	logMessage(fmt.Sprintf("Downloaded [%s]", url))
	logMessage(fmt.Sprintf("Finished at %s", time.Now().Format(time.RFC1123)))

	return nil
}

// logMessage is a helper function to log messages
func logMessage(msg string) {
	log.Println(msg)
}
