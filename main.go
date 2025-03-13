package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"wget/bckgrdDownload"
	"wget/fileDownload"
	"wget/inputDownload"
	"wget/mirrorDownload"
	"wget/rateDownload"
)

func main() {
	// Command-line flags
	background := flag.Bool("B", false, "Run download in background and log output")
	inputFile := flag.String("i", "", "Download multiple files from an input file")
	rateLimit := flag.String("rate-limit", "", "Limit download speed (e.g., 300k, 700k, 2M)")
	mirror := flag.Bool("mirror", false, "Mirror a website")
	convertLinks := flag.Bool("convert-links", false, "Convert links for offline browsing")
	reject := flag.String("reject", "", "Comma-separated list of file extensions to reject")
	exclude := flag.String("X", "", "Comma-separated list of paths to exclude")
	output := flag.String("O", "", "Save as different filename")
	saveDir := flag.String("P", "", "Save file in specific directory")

	flag.Parse()

	// Ensure URL or input file is provided
	if flag.NArg() == 0 && *inputFile == "" {
		fmt.Println("Usage: wget [options] <URL>")
		flag.Usage()
		os.Exit(1)
	}

	url := flag.Arg(0)

	// Convert reject and exclude flags to slices
	var rejectExtensions, excludeDirs []string
	if *reject != "" {
		rejectExtensions = strings.Split(*reject, ",")
	}
	if *exclude != "" {
		excludeDirs = strings.Split(*exclude, ",")
	}

	// Background download
	if *background {
		log.Println("Starting background download...")
		bckgrdDownload.Start(url)
		return
	}

	// Download multiple files from input list
	if *inputFile != "" {
		inputDownload.Start(*inputFile)
		return
	}

	// Rate-limited download
	if *rateLimit != "" {
		rateDownload.Start(url, *rateLimit)
		return
	}

	// Mirror a website
	if *mirror {
		mirrorDownload.Start(url, *convertLinks, rejectExtensions, excludeDirs)
		return
	}

	// Normal file download
	fileDownload.Start(url, *output, *saveDir)

	//strings.Join(rejectExtensions, ","), strings.Join(excludeDirs, ",")
}
