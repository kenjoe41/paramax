package aggregator

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/corpix/uarand"
	"github.com/hashicorp/go-retryablehttp"
)

type WebArchive struct {
}

func (wa WebArchive) FetchURLs(domain string, subs bool, excludeExtensions []string) ([]string, error) {
	var urls []string

	// Fetch the URL from the Web Archive
	var urlToFetch string
	if subs {
		urlToFetch = fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=*.%s/*&output=txt&fl=original&collapse=urlkey", domain)
	} else {
		urlToFetch = fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/*&output=txt&fl=original&collapse=urlkey", domain)
	}

	// Create a new retryable HTTP client
	client := retryablehttp.NewClient()
	client.RetryMax = 3                   // Maximum number of retries
	client.RetryWaitMax = 5 * time.Second // Maximum wait time between retries
	client.Backoff = retryablehttp.LinearJitterBackoff

	// Disable debug logging
	client.Logger = nil

	// Create a new request
	req, err := retryablehttp.NewRequest("GET", urlToFetch, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", uarand.GetRandom())

	// Make the request using the retryable client
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Web Archive: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Split the response into lines
	lines := strings.Split(string(body), "\n")

	// Extract the URLs from each line
	uniqueURLs := make(map[string]bool)
	for _, rawURL := range lines {

		if strings.HasPrefix(rawURL, "http") {
			u, err := url.Parse(rawURL)
			if err != nil {
				continue
			}

			// Check if the URL has an excluded file extension
			if excludeExtensions != nil {
				excluded := false
				for _, ext := range excludeExtensions {
					if !strings.HasPrefix(ext, ".") {
						ext = "." + ext
					}
					if strings.HasSuffix(u.Path, ext) {
						excluded = true
						break
					}
				}
				if excluded {
					continue
				}
			}
			uniqueURLs[u.String()] = true // TODO: This isn't workinig well for some reason maybe varying parameter values. Re-Deduplicating later after i change value to <placeholder> keyword.
		}
	}

	// Convert the map to a slice and return it
	for url := range uniqueURLs {
		urls = append(urls, url)
	}
	return urls, nil
}
