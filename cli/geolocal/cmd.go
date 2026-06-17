package geolocal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func Run(args []string) {
	output := "text"
	nameArgs := []string{}
	city := ""

	skipNext := false
	for i, a := range args {
		if skipNext {
			skipNext = false
			continue
		}
		switch {
		case a == "--json" || a == "-j":
			output = "json"
		case a == "--city" || a == "-c":
			if i+1 < len(args) {
				city = args[i+1]
				skipNext = true
			}
		case strings.HasPrefix(a, "--city="):
			city = strings.TrimPrefix(a, "--city=")
		case strings.HasPrefix(a, "-c="):
			city = strings.TrimPrefix(a, "-c=")
		case !strings.HasPrefix(a, "-"):
			nameArgs = append(nameArgs, a)
		}
	}

	if len(nameArgs) < 1 {
		fmt.Fprintln(os.Stderr, "usage: geo local <business-name> --city <city> [--json|-j]")
		fmt.Fprintln(os.Stderr, "       geo local \"Brewok Prime\" --city \"Malang\"")
		os.Exit(1)
	}

	businessName := strings.Join(nameArgs, " ")
	if city == "" {
		city = "Malang"
	}

	result, err := Analyze(businessName, city)
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

func printText(r *LocalResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Business:\t%s (%s)\n", r.BusinessName, r.City)
	fmt.Fprintf(w, "Local GEO Score:\t%s (%d/%d)\n", r.Label, r.Score, r.MaxScore)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Local Entity Cluster:")
	for _, e := range r.EntityCluster {
		fmt.Fprintf(w, "  - %s\n", e)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Checklist:")
	for _, c := range r.Checklist {
		status := "⬜"
		mark := " "
		if c.Done {
			status = "✅"
			mark = "x"
		}
		fmt.Fprintf(w, "  %s [%s] %s (%s)\n", status, mark, c.Task, c.Impact)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Schema Examples:")
	for _, s := range r.Schemas {
		fmt.Fprintf(w, "\n  [%s]\n", s.Type)
		fmt.Fprintf(w, "  %s\n", truncate(s.Source, 200))
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
