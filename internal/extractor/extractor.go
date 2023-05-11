// internal/extractor/extractor.go
package extractor

import (
	"net/url"
)

// ExtractParameters extracts the parameter names from the given URL.
// Returns a slice of parameter names.
func ExtractParameters(urlStr string) []string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		// Handle error appropriately
	}

	parameters := make([]string, 0)

	// Get the query parameters
	query := parsedURL.Query()

	// Extract the parameter names
	for param := range query {
		parameters = append(parameters, param)
	}

	return parameters
}
