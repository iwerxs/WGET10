package rateDownload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Start handles downloading a file with rate limiting
func Start(url string, rateLimit string) error {
	// Parse the rate limit (e.g., "300k", "2M")
	parsedRate, err := parseRateLimit(rateLimit)
	if err != nil {
		return err
	}

	// Make a GET request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Get file name from URL
	fileName := getFileNameFromURL(url)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Buffered read for controlled speed
	buffer := make([]byte, 4096) // Read in chunks
	var totalBytesDownloaded int64
	startTime := time.Now()

	// Download loop with rate limit enforcement
	for {
		// Read from response body
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file: %v", err)
		}
		if n == 0 {
			break
		}

		// Write to file
		_, err = file.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf("failed to write data to file: %v", err)
		}

		// Update total downloaded
		totalBytesDownloaded += int64(n)

		// Throttle speed
		throttleDownload(n, parsedRate)

		// Display progress
		displayProgress(totalBytesDownloaded, resp.ContentLength, startTime, parsedRate)
	}
	// Calculate total time taken
	totalTime := time.Since(startTime).Seconds()

	// Format total time to two decimal places
	fmt.Printf("\nTime taken: %.2f seconds\n", totalTime)
	fmt.Println("Download complete.")

	return nil
}

// parseRateLimit converts rate strings (e.g., "300k", "2M") into bytes per second
func parseRateLimit(rateLimit string) (int64, error) {
	var bytesPerSecond int64
	if strings.HasSuffix(rateLimit, "k") {
		bytesPerSecond = parseRateUnit(rateLimit, 1024)
	} else if strings.HasSuffix(rateLimit, "M") {
		bytesPerSecond = parseRateUnit(rateLimit, 1024*1024)
	} else {
		return 0, fmt.Errorf("unsupported rate limit unit: %s", rateLimit)
	}
	return bytesPerSecond, nil
}

// parseRateUnit extracts the numerical value and multiplies it by the given unit
func parseRateUnit(rateLimit string, unit int64) int64 {
	rate := strings.TrimSuffix(rateLimit, "k")
	rate = strings.TrimSuffix(rate, "M")
	var result int64
	fmt.Sscanf(rate, "%d", &result)
	return result * unit
}

// throttleDownload sleeps to maintain the rate limit
func throttleDownload(bytesDownloaded int, rateLimit int64) {
	// Time required for the downloaded bytes to match rate limit
	timeToSleep := float64(bytesDownloaded) / float64(rateLimit) // Seconds
	time.Sleep(time.Duration(timeToSleep * float64(time.Second))) // Convert to time.Duration
}

// displayProgress shows download status and speed
func displayProgress(totalBytes int64, totalSize int64, startTime time.Time, rateLimit int64) {
	// Calculate percentage
	percentage := float64(totalBytes) / float64(totalSize) * 100

	// Calculate elapsed time
	elapsedTime := time.Since(startTime).Seconds()

	// Ensure speed calculation does not divide by zero
	var speed float64
	if elapsedTime > 0 {
		speed = float64(totalBytes) / elapsedTime // Bytes per second
	}

	// Convert speed to KB/s or MB/s
	var speedDisplay string
	if speed >= 1024*1024 {
		speedDisplay = fmt.Sprintf("%.2f MB/s", speed/(1024*1024))
	} else {
		speedDisplay = fmt.Sprintf("%.2f KB/s", speed/1024)
	}

	// Estimate remaining time
	remainingTime := time.Duration((float64(totalSize)-float64(totalBytes))/speed) * time.Second

	// Print progress
	fmt.Printf("\r[%-50s] %.2f%% (Speed: %s, ETA: %v)",
		strings.Repeat("=", int(percentage/2)), percentage, speedDisplay, remainingTime)
}

// getFileNameFromURL extracts the filename from the URL
func getFileNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
