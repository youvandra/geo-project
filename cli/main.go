package main

import (
	"fmt"
	"os"

	"github.com/kucnigplaygame/geo-project/cli/cmdentity"
	"github.com/kucnigplaygame/geo-project/cli/cmdschema"
	"github.com/kucnigplaygame/geo-project/cli/cmdscore"
	"github.com/kucnigplaygame/geo-project/cli/cmdsitemap"
	"github.com/kucnigplaygame/geo-project/cli/cmdtopic"
	"github.com/kucnigplaygame/geo-project/cli/cmdtracker"
	"github.com/kucnigplaygame/geo-project/cli/geoaudit"
	"github.com/kucnigplaygame/geo-project/cli/geocrawl"
	"github.com/kucnigplaygame/geo-project/cli/geolocal"
	"github.com/kucnigplaygame/geo-project/cli/georeview"
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
	fmt.Println("  audit     AI answer auditor")
	fmt.Println("  crawl     Digital presence crawler")
	fmt.Println("  local     Local business GEO optimizer")
	fmt.Println("  review    Review analyzer")
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
		cmdentity.Run(args)
	case "tracker":
		cmdtracker.Run(args)
	case "sitemap":
		cmdsitemap.Run(args)
	case "audit":
		geoaudit.Run(args)
	case "crawl":
		geocrawl.Run(args)
	case "local":
		geolocal.Run(args)
	case "review":
		georeview.Run(args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
