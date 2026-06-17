package cmdsitemap

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Run(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: geo sitemap <baseurl> <content-dir> [--json|-j]")
		fmt.Fprintln(os.Stderr, "       geo sitemap <baseurl> <content-dir> [--output file.xml]")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  geo sitemap https://example.com ./blog/content")
		fmt.Fprintln(os.Stderr, "  geo sitemap https://example.com ./blog/content --output sitemap.xml")
		fmt.Fprintln(os.Stderr, "  geo sitemap https://example.com ./blog/content --json")
		os.Exit(1)
	}

	baseURL := args[0]
	contentDir := args[1]

	outputFormat := "xml"
	outputFile := ""
	restArgs := args[2:]

	for _, a := range restArgs {
		switch {
		case a == "--json" || a == "-j":
			outputFormat = "json"
		case strings.HasPrefix(a, "--output="):
			outputFile = strings.TrimPrefix(a, "--output=")
		case strings.HasPrefix(a, "-o"):
			// handled later
		default:
			// skip
		}
	}

	// Check for -o <filename>
	for i, a := range restArgs {
		if a == "-o" && i+1 < len(restArgs) {
			outputFile = restArgs[i+1]
		}
	}

	opts := SitemapOptions{
		BaseURL:    baseURL,
		ContentDir: contentDir,
		ExcludeDirs: []string{"node_modules", ".git", "vendor", "public", "resources"},
	}

	if outputFile != "" {
		xmlStr, err := GenerateXML(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		os.WriteFile(outputFile, []byte(xmlStr), 0644)
		fmt.Printf("sitemap written to %s (%d URLs)\n", outputFile, strings.Count(xmlStr, "<url>"))
		return
	}

	switch outputFormat {
	case "json":
		result, err := Generate(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
	default:
		xmlStr, err := GenerateXML(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(xmlStr)
	}
}
