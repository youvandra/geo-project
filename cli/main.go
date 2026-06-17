package main

import (
	"fmt"
	"os"

	"github.com/kucnigplaygame/geo-project/cli/cmdschema"
	"github.com/kucnigplaygame/geo-project/cli/cmdscore"
	"github.com/kucnigplaygame/geo-project/cli/cmdtopic"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("GEO Toolbelt - CLI tools for Generative Engine Optimization")
		fmt.Println()
		fmt.Println("Usage: geo <command> [args...]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  topic     Topic cluster analyzer")
		fmt.Println("  score     GEO readability score")
		fmt.Println("  schema    JSON-LD schema builder")
		fmt.Println("  entity    Entity linker")
		fmt.Println("  tracker   AI answer monitor")
		fmt.Println("  sitemap   Sitemap optimizer")
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "topic":
		cmdtopic.Run(args)
	case "score":
		cmdscore.Run(args)
	case "schema":
		cmdschema.Run(args)
	case "entity":
		fmt.Println("geo-entity: not yet implemented")
	case "tracker":
		fmt.Println("geo-tracker: not yet implemented")
	case "sitemap":
		fmt.Println("geo-sitemap: not yet implemented")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
