package georeview

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

const usage = `usage: geo review <business-name> [--json|-j]

Analyze customer reviews for GEO content insights.

Examples:
  geo review "Brewok Prime"
  geo review "Brewok Prime" --json
`

func Run(args []string) {
	output := "text"
	nameArgs := []string{}

	for _, a := range args {
		switch a {
		case "--json", "-j":
			output = "json"
		case "--help", "-h":
			fmt.Print(usage)
			os.Exit(0)
		default:
			nameArgs = append(nameArgs, a)
		}
	}

	if len(nameArgs) < 1 {
		fmt.Print(usage)
		os.Exit(1)
	}

	business := strings.Join(nameArgs, " ")

	result, err := Analyze(business, nil)
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

func printText(r *ReviewResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Business:\t%s\n", r.Business)
	fmt.Fprintf(w, "Reviews Analyzed:\t%d\n", r.ReviewCount)
	fmt.Fprintf(w, "Review Score:\t%s (%d/%d)\n", r.Label, r.Score, r.MaxScore)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Sentiment:")
	total := r.Sentiment.Positive + r.Sentiment.Negative + r.Sentiment.Neutral
	fmt.Fprintf(w, "  Positive:\t%d/%d (%.0f%%)\n", r.Sentiment.Positive, total, float64(r.Sentiment.Positive)/float64(total)*100)
	fmt.Fprintf(w, "  Negative:\t%d/%d\n", r.Sentiment.Negative, total)
	fmt.Fprintf(w, "  Neutral:\t%d/%d\n", r.Sentiment.Neutral, total)
	fmt.Fprintln(w)

	if len(r.Entities) > 0 {
		fmt.Fprintln(w, "Top Entities:")
		sort.Slice(r.Entities, func(i, j int) bool {
			return r.Entities[i].Count > r.Entities[j].Count
		})
		for _, e := range r.Entities {
			if e.Count > 1 || len(r.Entities) <= 10 {
				fmt.Fprintf(w, "  %s (%s, %dx)\n", e.Word, e.Type, e.Count)
			}
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "Content Ideas:")
	for _, idea := range r.ContentIdeas {
		fmt.Fprintf(w, "  • %s\n", idea)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Review Schema:")
	fmt.Fprintln(w, r.Schema)
	w.Flush()
}
