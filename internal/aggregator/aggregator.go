// internal/aggregator/aggregator.go
package aggregator

// Aggregator interface defines the contract for fetching URLs from an aggregator.
type Aggregator interface {
	FetchURLs(domain string, fetchSubdomains bool, excludeExtensions []string) ([]string, error)
}
