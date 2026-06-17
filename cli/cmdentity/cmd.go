package cmdentity

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
		fmt.Fprintln(os.Stderr, "usage: geo entity <content> | geo entity --stdin")
		fmt.Fprintln(os.Stderr, "       geo entity --stdin < file.txt")
		fmt.Fprintln(os.Stderr, "       geo entity --json <content>")
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

func printText(r *EntityResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Found %d linked entities\n\n", len(r.Entities))
	for _, e := range r.Entities {
		fmt.Fprintf(w, "  %s (%s)\n", e.Name, e.Confidence)
		fmt.Fprintf(w, "    %s\n", truncateText(e.Description, 120))
		fmt.Fprintf(w, "    %s\n", e.WikipediaURL)
		fmt.Fprintln(w)
	}
	w.Flush()
}
