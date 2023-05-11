// internal/modifier/modifier.go
package modifier

import (
	"net/url"
)

// ModifyParameters replaces the values of parameters in the given URL with "FUZZ".
// Returns a list of modified URLs.
func ModifyParameters(urlStr string, parameters []string, uniqueURLs []string, placeholder string) ([]string, []string) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return []string{}, uniqueURLs
	}

	// Get the query parameters
	query := parsedURL.Query()

	// Modify the parameter values
	for _, param := range parameters {
		query.Set(param, placeholder)
	}

	// Update the URL with modified query parameters
	parsedURL.RawQuery = query.Encode()

	modifiedURL := parsedURL.String()

	if !containsURL(uniqueURLs, modifiedURL) {
		uniqueURLs = append(uniqueURLs, modifiedURL)

		return []string{modifiedURL}, uniqueURLs
	}

	return []string{}, uniqueURLs
}

// Helper function to check if a URL exists in a slice
func containsURL(urls []string, url string) bool {
	for _, u := range urls {
		if u == url {
			return true
		}
	}
	return false
}
