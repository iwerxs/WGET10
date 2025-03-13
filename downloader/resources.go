package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// DownloadResources extracts and downloads all resources from a given HTML page.
func DownloadResources(htmlContent, baseURL, saveDir string, excludeDirs []string) error {
	links := extractLinks(htmlContent, baseURL)

	// Extract and download images from inline <style> blocks
	cssImages := extractImagesFromCSSContent(htmlContent, baseURL)
	links = append(links, cssImages...)

	for _, link := range links {
		// Skip excluded directories
		if shouldExclude(link, excludeDirs) {
			fmt.Println("Skipping excluded directory:", link)
			continue
		}

		err := downloadResource(link, saveDir)
		if err != nil {
			fmt.Println("Error downloading:", link, "-", err)
		}
	}

	// Update local references in downloaded HTML
	htmlContent = adjustCSSLinks(htmlContent, baseURL, saveDir)
	htmlContent = adjustLinks(htmlContent, baseURL, saveDir)

	// Save updated HTML file
	indexPath := filepath.Join(saveDir, "index.html")
	err := os.WriteFile(indexPath, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to update index.html: %v", err)
	}

	fmt.Println("Updated index.html with correct offline references.")
	return nil
}

// shouldExclude checks if a URL belongs to any excluded directory
func shouldExclude(resourceURL string, excludeDirs []string) bool {
	for _, dir := range excludeDirs {
		if strings.Contains(resourceURL, dir) {
			return true
		}
	}
	return false
}

// extractLinks finds all valid resource links in the HTML
func extractLinks(htmlContent, baseURL string) []string {
	var links []string
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "link" || token.Data == "script" || token.Data == "img" {
				for _, attr := range token.Attr {
					if attr.Key == "href" || attr.Key == "src" {
						link := resolveURL(attr.Val, baseURL)
						if link != "" && !strings.HasPrefix(link, "mailto:") {
							links = append(links, link)
						}
					}
				}
			} else if token.Data == "style" {
				// Extract background images from <style> blocks
				cssContent := extractStyleContent(tokenizer)
				cssImages := extractImagesFromCSSContent(cssContent, baseURL)
				links = append(links, cssImages...)
			}
		}
	}
	return links
}

// extractStyleContent extracts raw CSS from a <style> block
func extractStyleContent(tokenizer *html.Tokenizer) string {
	var cssContent strings.Builder
	for {
		tt := tokenizer.Next()
		if tt == html.TextToken {
			cssContent.WriteString(tokenizer.Token().Data)
		} else if tt == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "style" {
				break
			}
		}
	}
	return cssContent.String()
}

// extractImagesFromCSSContent extracts background-image URLs from CSS content
func extractImagesFromCSSContent(cssContent, baseURL string) []string {
	var images []string
	
	// Match all background-image declarations
	re := regexp.MustCompile(`background-image\s*:\s*[^;]*`)
	matches := re.FindAllString(cssContent, -1)
	
	for _, match := range matches {
		// Extract URLs from within url() functions
		urlRe := regexp.MustCompile(`url\(['"]?([^'")]+)['"]?\)`)
		urlMatches := urlRe.FindAllStringSubmatch(match, -1)
		
		for _, urlMatch := range urlMatches {
			if len(urlMatch) > 1 {
				// Clean and resolve each URL
				imgURL := strings.TrimSpace(urlMatch[1])
				resolvedURL := resolveURL(imgURL, baseURL)
				if resolvedURL != "" {
					images = append(images, resolvedURL)
				}
			}
		}
	}
	return images
}

// splitAndResolveURLs splits comma-separated URLs and resolves them
func splitAndResolveURLs(urlList, baseURL string) []string {
	var resolvedURLs []string
	// Split the URLs by comma, remove any leading/trailing spaces
	URLs := strings.Split(urlList, ",")
	for _, u := range URLs {
		trimmed := strings.TrimSpace(u)
		if trimmed != "" {
			resolved := resolveURL(trimmed, baseURL)
			if resolved != "" {
				resolvedURLs = append(resolvedURLs, resolved)
			}
		}
	}
	return resolvedURLs
}

// resolveURL converts relative URLs to absolute URLs
func resolveURL(imgURL, baseURL string) string {
	// If it's already an absolute URL, return it
	if strings.HasPrefix(imgURL, "http://") || strings.HasPrefix(imgURL, "https://") {
		return imgURL
	}
	
	// Parse the base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	
	// Resolve the relative URL against the base
	absURL, err := base.Parse(imgURL)
	if err != nil {
		return ""
	}
	
	return absURL.String()
}


// downloadResource downloads CSS, JS, and image files
func downloadResource(fileURL, saveDir string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		return nil // Avoid downloading extra HTML pages
	}

	fileName := filepath.Base(fileURL)
	filePath := filepath.Join(saveDir, fileName)

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded:", filePath)
	return nil
}

// adjustLinks modifies links for offline browsing
func adjustLinks(htmlContent, baseURL, saveDir string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))
	var modifiedHTML strings.Builder

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			// End of document, return the modified HTML
			return modifiedHTML.String()
		case html.TextToken:
			// Append text token to the modified HTML
			modifiedHTML.WriteString(tokenizer.Token().Data)
		case html.StartTagToken, html.SelfClosingTagToken:
			// Start tags and self-closing tags
			token := tokenizer.Token()
			modifiedHTML.WriteString("<" + token.Data)
			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					// Resolve URL for the href or src attributes
					attr.Val = "./" + filepath.Base(resolveURL(attr.Val, baseURL))
				}
				modifiedHTML.WriteString(fmt.Sprintf(` %s="%s"`, attr.Key, attr.Val))
			}
			modifiedHTML.WriteString(">")
		case html.EndTagToken:
			// End tags
			token := tokenizer.Token()
			modifiedHTML.WriteString("</" + token.Data + ">")
		}
	}
}


// adjustCSSLinks modifies inline <style> blocks to reference local files
func adjustCSSLinks(htmlContent, baseURL, saveDir string) string {
	re := regexp.MustCompile(`background-image\s*:\s*url\(['"]?([^'")]+)['"]?\)`)
	return re.ReplaceAllStringFunc(htmlContent, func(match string) string {
		originalURL := re.FindStringSubmatch(match)[1]
		localFile := "./" + filepath.Base(resolveURL(originalURL, baseURL))
		return strings.Replace(match, originalURL, localFile, 1)
	})
}
