package cmdscore

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

func Run(args []string) {
	output := "text"
	contentArgs := []string{}
	fromStdin := false

	for _, a := range args {
		switch a {
		case "--json", "-j":
			output = "json"
		case "--stdin", "-s":
			fromStdin = true
		default:
			contentArgs = append(contentArgs, a)
		}
	}

	var content string
	if fromStdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			os.Exit(1)
		}
		content = string(data)
	} else if len(contentArgs) > 0 {
		content = strings.Join(contentArgs, " ")
	} else {
		fmt.Fprintln(os.Stderr, "usage: geo score <content> | geo score --stdin")
		fmt.Fprintln(os.Stderr, "       geo score --stdin < file.txt")
		fmt.Fprintln(os.Stderr, "       geo score --json <content>")
		os.Exit(1)
	}

	result := Analyze(content)

	switch output {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
	default:
		printText(result)
	}
}

func printText(r *ScoreResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "GEO Score:\t%s (%d/%d)\n", r.Overall.Label, r.Overall.Score, r.Overall.Max)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Breakdown:")
	fmt.Fprintf(w, "  Structure\t\t%s (%d/%d)\n", r.Breakdown.Structure.Label, r.Breakdown.Structure.Score, r.Breakdown.Structure.Max)
	fmt.Fprintf(w, "  Q&A Coverage\t\t%s (%d/%d)\n", r.Breakdown.QACoverage.Label, r.Breakdown.QACoverage.Score, r.Breakdown.QACoverage.Max)
	fmt.Fprintf(w, "  Entity Richness\t%s (%d/%d)\n", r.Breakdown.EntityRichness.Label, r.Breakdown.EntityRichness.Score, r.Breakdown.EntityRichness.Max)
	fmt.Fprintf(w, "  Citation Quality\t%s (%d/%d)\n", r.Breakdown.CitationQuality.Label, r.Breakdown.CitationQuality.Score, r.Breakdown.CitationQuality.Max)
	fmt.Fprintf(w, "  Schema Readiness\t%s (%d/%d)\n", r.Breakdown.SchemaReadiness.Label, r.Breakdown.SchemaReadiness.Score, r.Breakdown.SchemaReadiness.Max)
	fmt.Fprintf(w, "  Readability\t\t%s (%d/%d)\n", r.Breakdown.Readability.Label, r.Breakdown.Readability.Score, r.Breakdown.Readability.Max)

	if len(r.Suggestion) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Suggestions:")
		for _, s := range r.Suggestion {
			fmt.Fprintf(w, "  • %s\n", s)
		}
	}

	w.Flush()
}
