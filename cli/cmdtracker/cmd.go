package cmdtracker

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
	showHistory := false

	for _, a := range args {
		switch a {
		case "--json", "-j":
			output = "json"
		case "--history", "-h":
			showHistory = true
		default:
			topicArgs = append(topicArgs, a)
		}
	}

	if len(topicArgs) < 1 {
		fmt.Fprintln(os.Stderr, "usage: geo tracker <topic> [--json|-j] [--history|-h]")
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
		printText(result, showHistory)
	}
}

func printText(r *TrackerResult, showHistory bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Topic:\t%s\n", r.Topic)
	fmt.Fprintf(w, "Trend:\t%s\n", r.Trend)
	fmt.Fprintf(w, "Last Checked:\t%s\n", r.LastChecked.Format("2006-01-02 15:04"))
	fmt.Fprintf(w, "Current Daily Views:\t%d\n", r.Current.PageViews)
	fmt.Fprintln(w)

	if r.Current.Description != "" {
		fmt.Fprintf(w, "Description:\t%s\n", truncateText(r.Current.Description, 150))
		fmt.Fprintln(w)
	}

	if showHistory && len(r.History) > 0 {
		fmt.Fprintln(w, "Tracking History:")
		for _, h := range r.History {
			fmt.Fprintf(w, "  %s\t%d views\n", h.Date.Format("2006-01-02"), h.PageViews)
		}
	}

	w.Flush()
}

func truncateText(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}
