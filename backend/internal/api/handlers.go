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

func HandleCrawlForm(w http.ResponseWriter, r *http.Request) {
	render(w, "crawl", PageData{Title: "Digital Presence Crawler"})
}

func HandleCrawlAnalyze(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := r.FormValue("query")
	if query == "" {
		render(w, "crawl", PageData{Title: "Digital Presence Crawler", Error: "Query is required"})
		return
	}

	result := engine.CrawlBrand(query)
	data := PageData{Title: "Digital Presence Crawler", Result: result, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/crawl.html")
		tmpl.ExecuteTemplate(w, "crawl-result", data)
		return
	}
	render(w, "crawl", data)
}

func HandleAuditForm(w http.ResponseWriter, r *http.Request) {
	render(w, "audit", PageData{Title: "AI Answer Auditor"})
}

func HandleAuditAnalyze(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	brand := r.FormValue("brand")
	if brand == "" {
		render(w, "audit", PageData{Title: "AI Answer Auditor", Error: "Brand is required"})
		return
	}

	result, err := engine.AuditBrand(brand)
	if err != nil {
		render(w, "audit", PageData{Title: "AI Answer Auditor", Error: err.Error()})
		return
	}

	data := PageData{Title: "AI Answer Auditor", Result: result, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/audit.html")
		tmpl.ExecuteTemplate(w, "audit-result", data)
		return
	}
	render(w, "audit", data)
}

func HandleLocalForm(w http.ResponseWriter, r *http.Request) {
	render(w, "local", PageData{Title: "Local Business GEO"})
}

func HandleLocalAnalyze(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	business := r.FormValue("business")
	city := r.FormValue("city")
	if business == "" {
		render(w, "local", PageData{Title: "Local Business GEO", Error: "Business name is required"})
		return
	}
	if city == "" {
		city = "Malang"
	}

	result := engine.AnalyzeLocal(business, city)
	data := PageData{Title: "Local Business GEO", Result: result, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/local.html")
		tmpl.ExecuteTemplate(w, "local-result", data)
		return
	}
	render(w, "local", data)
}

func HandleReviewForm(w http.ResponseWriter, r *http.Request) {
	render(w, "review", PageData{Title: "Review Analyzer"})
}

func HandleReviewAnalyze(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	business := r.FormValue("business")
	if business == "" {
		render(w, "review", PageData{Title: "Review Analyzer", Error: "Business name is required"})
		return
	}

	result := engine.AnalyzeReviews(business)
	data := PageData{Title: "Review Analyzer", Result: result, Time: time.Now().Format(time.RFC3339)}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/review.html")
		tmpl.ExecuteTemplate(w, "review-result", data)
		return
	}
	render(w, "review", data)
}

func HandleReportForm(w http.ResponseWriter, r *http.Request) {
	render(w, "report", PageData{Title: "GEO Report"})
}

func HandleReportGenerate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	brand := r.FormValue("brand")
	if brand == "" {
		render(w, "report", PageData{Title: "GEO Report", Error: "Brand name is required"})
		return
	}

	audit, err := engine.AuditBrand(brand)
	if err != nil {
		audit = &engine.AuditResult{Brand: brand, Label: "Error"}
	}

	crawl := engine.CrawlBrand(brand)
	review := engine.AnalyzeReviews(brand)
	city := r.FormValue("city")
	if city == "" {
		city = "Malang"
	}
	local := engine.AnalyzeLocal(brand, city)

	report := map[string]interface{}{
		"brand":  brand,
		"city":   city,
		"audit":  audit,
		"crawl":  crawl,
		"review": review,
		"local":  local,
	}

	data := PageData{
		Title:  "GEO Report — " + brand,
		Result: report,
		Time:   time.Now().Format(time.RFC3339),
	}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, _ := template.ParseFiles("web/templates/report.html")
		tmpl.ExecuteTemplate(w, "report-result", data)
		return
	}
	render(w, "report", data)
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
