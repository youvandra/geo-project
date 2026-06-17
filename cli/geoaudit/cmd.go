package geoaudit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func Run(args []string) {
	output := "text"
	topicArgs := []string{}

	for _, a := range args {
		switch a {
		case "--json", "-j":
			output = "json"
		default:
			topicArgs = append(topicArgs, a)
		}
	}

	if len(topicArgs) < 1 {
		fmt.Fprintln(os.Stderr, "usage: geo audit <brand-name> [--json|-j]")
		fmt.Fprintln(os.Stderr, "       geo audit \"Brewok Prime Malang\"")
		os.Exit(1)
	}

	brand := strings.Join(topicArgs, " ")
	result, err := Analyze(brand)
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

func printText(r *AuditResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Brand:\t%s\n", r.Brand)
	fmt.Fprintf(w, "GEO Readiness:\t%s (%d/%d)\n", r.Label, r.TotalScore, r.MaxScore)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Breakdown:")
	fmt.Fprintf(w, "  Entity Score:\t%d/30 — %s\n", r.EntityScore, labelShort(r.EntityScore, 30))
	fmt.Fprintf(w, "  Authority Score:\t%d/30 — %s\n", r.AuthorityScore, labelShort(r.AuthorityScore, 30))
	fmt.Fprintf(w, "  Citation Score:\t%d/25 — %s\n", r.CitationScore, labelShort(r.CitationScore, 25))
	fmt.Fprintf(w, "  Structure Score:\t%d/15 — %s\n", r.StructureScore, labelShort(r.StructureScore, 15))
	fmt.Fprintln(w)

	if r.Details.OnWikipedia {
		fmt.Fprintf(w, "Wikipedia:\tyes (%s)\n", r.Details.WikiTitle)
	} else {
		fmt.Fprintln(w, "Wikipedia:\tnot found")
	}
	if len(r.Details.Competitors) > 0 {
		fmt.Fprintln(w, "Related entities:")
		for _, c := range r.Details.Competitors {
			fmt.Fprintf(w, "  - %s\n", c)
		}
	}
	fmt.Fprintln(w)

	if len(r.Suggestions) > 0 {
		fmt.Fprintln(w, "Suggestions:")
		for _, s := range r.Suggestions {
			fmt.Fprintf(w, "  • %s\n", s)
		}
	}
	w.Flush()
}

func labelShort(score, max int) string {
	r := float64(score) / float64(max)
	switch {
	case r >= 0.8:
		return "Excellent"
	case r >= 0.6:
		return "Good"
	case r >= 0.4:
		return "Fair"
	default:
		return "Needs Work"
	}
}
