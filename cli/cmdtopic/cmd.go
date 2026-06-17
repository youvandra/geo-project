package cmdtopic

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
		fmt.Fprintln(os.Stderr, "usage: geo topic <topic> [--json|-j]")
		os.Exit(1)
	}

	topic := strings.Join(topicArgs, " ")
	result, err := Analyze(topic)
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

func printText(r *TopicResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Topic:\t%s\n", r.Topic)
	fmt.Fprintf(w, "Description:\t%s\n", truncate(r.Description, 200))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Subtopics:")
	for _, s := range r.Subtopics {
		fmt.Fprintf(w, "  - %s\n", s)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Related Entities:")
	for _, e := range r.Entities {
		fmt.Fprintf(w, "  - %s (%s)\n", e.Name, e.Relevance)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Suggested Questions for GEO:")
	for _, q := range r.Questions {
		fmt.Fprintf(w, "  • %s\n", q)
	}
	w.Flush()
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}
