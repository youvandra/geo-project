package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kucnigplaygame/geo-project/backend/internal/api"
	"github.com/kucnigplaygame/geo-project/backend/internal/db"
)

func init() {
	loadEnv()
}

func loadEnv() {
	paths := []string{".env", "../.env", filepath.Join(os.Getenv("HOME"), ".geo-env")}
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
		break
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Try PostgreSQL (non-fatal if unavailable)
	if err := db.Connect(db.DefaultConfig()); err != nil {
		log.Printf("PostgreSQL not available: %v", err)
		log.Println("Tracker will use file-based storage")
	} else {
		defer db.DB.Close()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", api.HandleIndex)
	mux.HandleFunc("GET /health", api.HandleHealth)

	mux.HandleFunc("GET /topic", api.HandleTopicForm)
	mux.HandleFunc("POST /topic", api.HandleTopicAnalyze)

	mux.HandleFunc("GET /score", api.HandleScoreForm)
	mux.HandleFunc("POST /score", api.HandleScoreAnalyze)

	mux.HandleFunc("GET /schema", api.HandleSchemaForm)
	mux.HandleFunc("POST /schema", api.HandleSchemaBuild)

	mux.HandleFunc("GET /entity", api.HandleEntityForm)
	mux.HandleFunc("POST /entity", api.HandleEntityAnalyze)

	mux.HandleFunc("GET /tracker", api.HandleTrackerForm)
	mux.HandleFunc("POST /tracker", api.HandleTrackerAnalyze)

	mux.HandleFunc("GET /sitemap", api.HandleSitemapForm)
	mux.HandleFunc("POST /sitemap", api.HandleSitemapGenerate)

	mux.HandleFunc("GET /crawl", api.HandleCrawlForm)
	mux.HandleFunc("POST /crawl", api.HandleCrawlAnalyze)

	mux.HandleFunc("GET /audit", api.HandleAuditForm)
	mux.HandleFunc("POST /audit", api.HandleAuditAnalyze)

	mux.HandleFunc("GET /local", api.HandleLocalForm)
	mux.HandleFunc("POST /local", api.HandleLocalAnalyze)

	mux.HandleFunc("GET /review", api.HandleReviewForm)
	mux.HandleFunc("POST /review", api.HandleReviewAnalyze)

	mux.HandleFunc("GET /report", api.HandleReportForm)
	mux.HandleFunc("POST /report", api.HandleReportGenerate)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	log.Printf("GEO Dashboard starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
