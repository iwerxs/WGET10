package inputDownload

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"wget/fileDownload"
)

// Start handles downloading multiple files listed in the input file.
func Start(inputFile string) {
	// Open the input file
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", inputFile, err)
	}
	defer file.Close()

	// Read URLs from the input file
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()

		// Notify the user about the download starting asynchronously
		fmt.Printf("Starting download for: %s\n", url)

		// Add the download task to the WaitGroup
		wg.Add(1)

		// Start the download in a goroutine
		go func(url string) {
			defer wg.Done()

			// Call the download function (fileDownload.Start handles the actual download)
			fileDownload.Start(url, "", "")
			if err != nil {
				log.Printf("Error downloading %s: %v", url, err)
			} else {
				fmt.Printf("Download complete: %s\n", url)
			}
		}(url)
	}

	// Wait for all downloads to finish
	wg.Wait()

	// Notify the user once all downloads are complete
	fmt.Println("All downloads complete.")
}
