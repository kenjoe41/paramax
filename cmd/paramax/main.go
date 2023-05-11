// cmd/paramax/main.go
package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/kenjoe41/paramax/internal/aggregator"
	"github.com/kenjoe41/paramax/internal/extractor"
	"github.com/kenjoe41/paramax/internal/modifier"
	"github.com/spf13/cobra"
)

var (
	domainFlag        string
	subsFlag          bool
	excludeFlag       string
	excludeExtensions []string
	outputFileFlag    string
	placeholderFlag   string
	silentFlag        bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "paramax",
		Short: "A program for analyzing URL parameters",
		Run: func(cmd *cobra.Command, args []string) {
			runPassiveMode(cmd)
		},
	}

	passiveCmd := &cobra.Command{
		Use:   "passive",
		Short: "Run the program in passive mode",
		Run: func(cmd *cobra.Command, args []string) {
			runPassiveMode(cmd)
		},
	}

	activeCmd := &cobra.Command{
		Use:   "active",
		Short: "Run the program in active mode",
		Run: func(cmd *cobra.Command, args []string) {
			// runActiveMode()
			fmt.Fprintln(os.Stderr, "Active mode has yet to be implemented. Consider using passive mode for now.")
			cmd.Usage()
		},
	}

	rootCmd.AddCommand(passiveCmd, activeCmd)

	rootCmd.PersistentFlags().StringVar(&domainFlag, "domain", "", "Domain to analyze (required; can be a full URL or just the domain)")
	rootCmd.PersistentFlags().BoolVar(&subsFlag, "subs", false, "Include subdomains when fetching from aggregators that support it (optional)")
	rootCmd.PersistentFlags().StringVar(&excludeFlag, "exclude", "", "Comma-separated list of file extensions to exclude")
	rootCmd.Flags().StringVarP(&outputFileFlag, "output", "o", "", "Output file name")
	rootCmd.PersistentFlags().StringVar(&placeholderFlag, "placeholder", "FUZZ", "The string to add as a placeholder after the parameter name.")
	rootCmd.Flags().BoolVarP(&silentFlag, "silent", "s", false, "Do not print the results to the screen if Output flag is specified.")

	passiveCmd.Flags().StringVar(&domainFlag, "domain", "", "Domain to analyze (required; can be a full URL or just the domain)")
	passiveCmd.Flags().BoolVar(&subsFlag, "subs", false, "Include subdomains when fetching from aggregators that support it (optional)")
	passiveCmd.Flags().StringVar(&excludeFlag, "exclude", "", "Comma-separated list of file extensions to exclude")
	passiveCmd.Flags().StringVarP(&outputFileFlag, "output", "o", "", "Output file name")
	passiveCmd.Flags().StringVar(&placeholderFlag, "placeholder", "FUZZ", "The string to add as a placeholder after the parameter name.")
	passiveCmd.Flags().BoolVarP(&silentFlag, "silent", "s", false, "Do not print the results to the screen if Output flag is specified.")

	activeCmd.Flags().StringVar(&domainFlag, "domain", "", "Domain to analyze (required; can be a full URL or just the domain)")
	activeCmd.Flags().BoolVar(&subsFlag, "subs", false, "Include subdomains when fetching from aggregators that support it (optional)")
	activeCmd.Flags().StringVar(&excludeFlag, "exclude", "", "Comma-separated list of file extensions to exclude")
	activeCmd.Flags().StringVarP(&outputFileFlag, "output", "o", "", "Output file name")
	activeCmd.Flags().StringVar(&placeholderFlag, "placeholder", "FUZZ", "The string to add as a placeholder after the parameter name.")
	activeCmd.Flags().BoolVarP(&silentFlag, "silent", "s", false, "Do not print the results to the screen if Output flag is specified.")

	rootCmd.SetUsageTemplate(
		`Usage: paramax [mode] [flags]
Modes:
  passive       Run the program in passive mode (default)
  active        Run the program in active mode

Flags:
{{.Flags.FlagUsages | trimTrailingWhitespaces}}

Use "paramax [command] --help" for more information about a command.`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runPassiveMode(cmd *cobra.Command) {
	domain := extractDomain(domainFlag)
	if domain == "" {
		fmt.Fprintf(os.Stderr, "Domain format error: %s, please enter a valid domain.\n", domain)
		cmd.Usage()
		os.Exit(1)
	}

	urlChan := make(chan string)
	modifiedURLsChan := make(chan string)
	done := make(chan struct{})

	go func() {
		defer close(urlChan)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			fetchAndProcessURLs(domain, urlChan, subsFlag, excludeExtensions)
		}()

		wg.Wait()
	}()

	go func() {
		defer close(modifiedURLsChan)

		var uniqueURLs []string
		for url := range urlChan {
			uniqueURLs = modifyURLParameters(url, modifiedURLsChan, uniqueURLs, placeholderFlag, false)
		}

		if outputFileFlag != "" {
			err := writeURLsToFile(uniqueURLs, outputFileFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write URLs to file: %v", err)
			}
		}

	}()

	go func() {
		defer close(done)
		for modifiedURL := range modifiedURLsChan {
			// Print the modified URL if the silent flag is not set or the output file flag is not set
			// This ensures that the modified URL is printed in all situations except when both the silent flag and the output file flag are set
			if !(silentFlag && outputFileFlag != "") {
				fmt.Println(modifiedURL)
			}

		}
	}()

	<-done
}

func fetchAndProcessURLs(domain string, urlChan chan<- string, subsFlag bool, excludeExtensions []string) {
	var urls []string

	aggregators := []aggregator.Aggregator{
		&aggregator.WebArchive{},
		// &aggregator.AlienVault{},
		// Add more aggregator implementations here
	}
	for _, aggregator := range aggregators {
		var err error
		urls, err = aggregator.FetchURLs(domain, subsFlag, excludeExtensions)
		if err != nil {
			log.Printf("Failed to fetch URLs for %s: %v", domain, err)
			return
		}

		for _, u := range urls {
			urlChan <- u
		}
	}
}

func modifyURLParameters(url string, modifiedURLsChan chan<- string, uniqueURLs []string, placeholder string, activeMode bool) []string {
	parameters := extractor.ExtractParameters(url)
	if len(parameters) == 0 {
		return uniqueURLs
	}

	if activeMode {
		// Uncomment the following code if using bruteforcer
		// bruteforcer := bruteforcer.NewBruteforcer()
		// modifiedURLs := bruteforcer.BruteForceParameters(url, parameters, uniqueURLs)
	} else { // Passive Mode
		var modifiedURLs []string
		modifiedURLs, uniqueURLs = modifier.ModifyParameters(url, parameters, uniqueURLs, placeholder)
		for _, modifiedURL := range modifiedURLs {
			modifiedURLsChan <- modifiedURL
		}
	}

	return uniqueURLs
}

func runActiveMode() {
	domain := extractDomain(domainFlag)

	urlChan := make(chan string)
	modifiedURLsChan := make(chan string)
	done := make(chan struct{})

	go func() {
		defer close(urlChan)
		defer close(modifiedURLsChan)
		defer close(done)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			fetchAndProcessURLs(domain, urlChan, subsFlag, excludeExtensions)
		}()

		wg.Wait()
	}()

	go func() {

		var uniqueURLs []string
		for url := range urlChan {
			uniqueURLs = modifyURLParameters(url, modifiedURLsChan, uniqueURLs, placeholderFlag, true)
		}

		if outputFileFlag != "" {
			err := writeURLsToFile(uniqueURLs, outputFileFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write URLs to file: %v", err)
			}
		}

	}()

	go func() {
		for modifiedURL := range modifiedURLsChan {
			fmt.Println(modifiedURL)
		}
	}()

	<-done
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func extractDomain(input string) string {
	if input == "" {
		return ""
	}

	// Clean domain if provided as a URL
	if strings.HasPrefix(input, "http") {
		parsedURL, err := url.Parse(input)
		if err != nil {
			exitWithError(fmt.Sprintf("Failed to parse domain: %v", err))
		}

		return parsedURL.Hostname()
	} else if strings.Contains(input, "/") {
		domain := strings.Split(input, "/")[0]

		return domain
	}
	return input
}

func writeURLsToFile(urls []string, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, url := range urls {
		_, err := file.WriteString(url + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
