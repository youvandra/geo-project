package geocrawl

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func Run(args []string) {
	output := "text"
	queryArgs := []string{}

	for _, a := range args {
		switch a {
		case "--json", "-j":
			output = "json"
		default:
			queryArgs = append(queryArgs, a)
		}
	}

	if len(queryArgs) < 1 {
		fmt.Fprintln(os.Stderr, "usage: geo crawl <query> [--json|-j]")
		fmt.Fprintln(os.Stderr, "       geo crawl \"Brewok Prime Malang\"")
		os.Exit(1)
	}

	query := strings.Join(queryArgs, " ")
	result, err := Crawl(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch output {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
	default:
		printText(result)
	}
}

func printText(r *CrawlResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Query:\t%s\n", r.Query)
	fmt.Fprintf(w, "Digital Presence:\t%s (%d/%d)\n", r.Label, r.Score, r.MaxScore)
	fmt.Fprintf(w, "Sources Hit:\t%d/%d\n", r.TotalHits, len(r.Sources))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Source Results:")
	for _, s := range r.Sources {
		status := "\u274c"
		if s.Found {
			status = "\u2705"
		}
		fmt.Fprintf(w, "  %s [%s]\n", s.Source, status)
		if s.Found {
			fmt.Fprintf(w, "     Title: %s\n", s.Title)
			fmt.Fprintf(w, "     URL:   %s\n", s.URL)
			if s.Snippet != "" {
				fmt.Fprintf(w, "     Preview: %s\n", truncate(s.Snippet, 100))
			}
		}
	}
	w.Flush()
}
