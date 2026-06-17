package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kucnigplaygame/geo-project/backend/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	log.Printf("GEO Dashboard starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
