package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kucnigplaygame/geo-project/backend/internal/engine"
)

type PageData struct {
	Title      string
	Content    string
	Result     interface{}
	ResultJSON string
	Error      string
	Time       string
}

func render(w http.ResponseWriter, page string, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		fmt.Sprintf("web/templates/%s.html", page),
	)
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "template error", 500)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("template execute error: %v", err)
		http.Error(w, "internal error", 500)
	}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, map[string]string{"status": "ok", "service": "geo-backend"})
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	render(w, "index", PageData{Title: "GEO Dashboard"})
}

func HandleTopicForm(w http.ResponseWriter, r *http.Request) {
	render(w, "topic", PageData{Title: "Topic Cluster Analyzer"})
}

func HandleTopicAnalyze(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	topic := r.FormValue("topic")
	if topic == "" {
		render(w, "topic", PageData{Title: "Topic Cluster Analyzer", Error: "Topic is required"})
		return
	}

	result, err := engine.AnalyzeTopic(topic)
	if err != nil {
		render(w, "topic", PageData{Title: "Topic Cluster Analyzer", Error: err.Error()})
		return
	}

	data := PageData{Title: "Topic Cluster Analyzer", Result: result, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/topic.html")
		tmpl.ExecuteTemplate(w, "topic-result", data)
		return
	}
	render(w, "topic", data)
}

func HandleScoreForm(w http.ResponseWriter, r *http.Request) {
	render(w, "score", PageData{Title: "GEO Readability Score"})
}

func HandleScoreAnalyze(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	if content == "" {
		body, _ := io.ReadAll(r.Body)
		content = string(body)
	}
	if content == "" {
		render(w, "score", PageData{Title: "GEO Readability Score", Error: "Content is required"})
		return
	}

	result := engine.AnalyzeScore(content)
	data := PageData{Title: "GEO Readability Score", Result: result, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/score.html")
		tmpl.ExecuteTemplate(w, "score-result", data)
		return
	}
	render(w, "score", data)
}

func HandleSchemaForm(w http.ResponseWriter, r *http.Request) {
	render(w, "schema", PageData{Title: "JSON-LD Schema Builder"})
}

func HandleSchemaBuild(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	schemaType := r.FormValue("type")
	fields := r.FormValue("fields")

	if schemaType == "" {
		render(w, "schema", PageData{Title: "JSON-LD Schema Builder", Error: "Schema type is required"})
		return
	}

	result, err := engine.BuildSchema(schemaType, fields)
	if err != nil {
		render(w, "schema", PageData{Title: "JSON-LD Schema Builder", Error: err.Error()})
		return
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	data := PageData{Title: "JSON-LD Schema Builder", Result: result, ResultJSON: string(b), Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/schema.html")
		tmpl.ExecuteTemplate(w, "schema-result", data)
		return
	}
	render(w, "schema", data)
}

func HandleEntityForm(w http.ResponseWriter, r *http.Request) {
	render(w, "entity", PageData{Title: "Entity Linker"})
}

func HandleEntityAnalyze(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	if content == "" {
		body, _ := io.ReadAll(r.Body)
		content = string(body)
	}
	if content == "" {
		render(w, "entity", PageData{Title: "Entity Linker", Error: "Content is required"})
		return
	}

	result := engine.AnalyzeEntities(content)
	data := PageData{Title: "Entity Linker", Result: result, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/entity.html")
		tmpl.ExecuteTemplate(w, "entity-result", data)
		return
	}
	render(w, "entity", data)
}

func HandleTrackerForm(w http.ResponseWriter, r *http.Request) {
	render(w, "tracker", PageData{Title: "AI Answer Tracker"})
}

func HandleTrackerAnalyze(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	topic := r.FormValue("topic")
	if topic == "" {
		render(w, "tracker", PageData{Title: "AI Answer Tracker", Error: "Topic is required"})
		return
	}

	result, err := engine.TrackTopic(topic)
	if err != nil {
		render(w, "tracker", PageData{Title: "AI Answer Tracker", Error: err.Error()})
		return
	}

	data := PageData{Title: "AI Answer Tracker", Result: result, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/tracker.html")
		tmpl.ExecuteTemplate(w, "tracker-result", data)
		return
	}
	render(w, "tracker", data)
}

func HandleSitemapForm(w http.ResponseWriter, r *http.Request) {
	render(w, "sitemap", PageData{Title: "Sitemap Generator"})
}

func HandleSitemapGenerate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	baseURL := r.FormValue("base_url")
	contentDir := r.FormValue("content_dir")

	if baseURL == "" || contentDir == "" {
		render(w, "sitemap", PageData{Title: "Sitemap Generator", Error: "Base URL and content directory are required"})
		return
	}

	result, err := engine.GenerateSitemap(baseURL, contentDir)
	if err != nil {
		render(w, "sitemap", PageData{Title: "Sitemap Generator", Error: err.Error()})
		return
	}

	data := PageData{Title: "Sitemap Generator", Result: result, ResultJSON: result.XML, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/sitemap.html")
		tmpl.ExecuteTemplate(w, "sitemap-result", data)
		return
	}
	render(w, "sitemap", data)
}
