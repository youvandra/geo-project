package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TopicResult struct {
	Topic       string   `json:"topic"`
	Description string   `json:"description"`
	Entities    []Entity `json:"entities"`
	Questions   []string `json:"questions"`
	Subtopics   []string `json:"subtopics"`
	AnalyzedAt  string   `json:"analyzed_at"`
}

type Entity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Relevance   string `json:"relevance"`
}

type ScoreResult struct {
	Overall    ScoreCategory `json:"overall"`
	Breakdown  Breakdown     `json:"breakdown"`
	Suggestion []string      `json:"suggestions"`
}

type ScoreCategory struct {
	Label string `json:"label"`
	Score int    `json:"score"`
	Max   int    `json:"max"`
}

type Breakdown struct {
	Structure       ScoreCategory `json:"structure"`
	QACoverage      ScoreCategory `json:"qa_coverage"`
	EntityRichness  ScoreCategory `json:"entity_richness"`
	CitationQuality ScoreCategory `json:"citation_quality"`
	SchemaReadiness ScoreCategory `json:"schema_readiness"`
	Readability     ScoreCategory `json:"readability"`
}

type EntityResult struct {
	Content   string          `json:"content"`
	Entities  []LinkedEntity  `json:"entities"`
	AnalyzedAt string         `json:"analyzed_at"`
}

type LinkedEntity struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	WikipediaURL string `json:"wikipedia_url"`
	Confidence   string `json:"confidence"`
}

type TrackerResult struct {
	Topic       string       `json:"topic"`
	Trend       string       `json:"trend"`
	PageViews   int          `json:"page_views"`
	Description string       `json:"description"`
	LastChecked string       `json:"last_checked"`
}

type SitemapOutput struct {
	BaseURL   string `json:"base_url"`
	TotalURLs int    `json:"total_urls"`
	XML       string `json:"xml"`
}

// wikiGet is a helper for Wikipedia API calls
func wikiGet(urlStr string, result interface{}) error {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "GEO-Project/1.0 (kucnigplaygame@gmail.com)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

func AnalyzeTopic(topic string) (*TopicResult, error) {
	// Fetch Wikipedia page
	u := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts&explaintext=true&redirects=1",
		url.QueryEscape(topic),
	)

	var wq struct {
		Query struct {
			Pages map[string]struct {
				PageID  int    `json:"pageid"`
				Title   string `json:"title"`
				Extract string `json:"extract"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := wikiGet(u, &wq); err != nil {
		return nil, err
	}

	var extract string
	for _, p := range wq.Query.Pages {
		if p.PageID > 0 {
			extract = p.Extract
			break
		}
	}

	return &TopicResult{
		Topic:       topic,
		Description: truncate(extract, 300),
		Entities:    extractEntities(extract),
		Questions:   generateQuestions(topic),
		Subtopics:   extractSubtopics(extract),
		AnalyzedAt:  time.Now().Format(time.RFC3339),
	}, nil
}

func extractEntities(extract string) []Entity {
	if extract == "" {
		return nil
	}
	entities := []Entity{}
	seen := map[string]bool{}
	for _, line := range strings.Split(extract, "\n") {
		for _, word := range strings.Fields(line) {
			w := strings.Trim(word, " ,;:()[]{}\"'.")
			if len(w) > 5 && w[0] >= 'A' && w[0] <= 'Z' && !seen[w] {
				seen[w] = true
				entities = append(entities, Entity{Name: w, Relevance: "related"})
			}
		}
	}
	if len(entities) > 10 {
		entities = entities[:10]
	}
	return entities
}

func extractSubtopics(extract string) []string {
	if extract == "" {
		return nil
	}
	subtopics := []string{}
	for _, line := range strings.Split(extract, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "==") && strings.HasSuffix(line, "==") {
			clean := strings.Trim(line, "= ")
			if clean != "" {
				subtopics = append(subtopics, clean)
			}
		}
	}
	return subtopics
}

func generateQuestions(topic string) []string {
	return []string{
		fmt.Sprintf("What is %s?", topic),
		fmt.Sprintf("How does %s work?", topic),
		fmt.Sprintf("Why is %s important?", topic),
		fmt.Sprintf("What are the key components of %s?", topic),
		fmt.Sprintf("How to get started with %s?", topic),
		fmt.Sprintf("What are common challenges in %s?", topic),
		fmt.Sprintf("What is the future of %s?", topic),
	}
}

func AnalyzeScore(content string) *ScoreResult {
	structure := scoreStructure(content)
	qa := scoreQA(content)
	entity := scoreEntity(content)
	citation := scoreCitation(content)
	schema := scoreSchema(content)
	readability := scoreReadability(content)

	total := structure.Score + qa.Score + entity.Score + citation.Score + schema.Score + readability.Score
	maxTotal := structure.Max + qa.Max + entity.Max + citation.Max + schema.Max + readability.Max

	label := "Poor"
	switch {
	case total >= 80:
		label = "Excellent"
	case total >= 60:
		label = "Good"
	case total >= 40:
		label = "Fair"
	case total >= 20:
		label = "Needs Work"
	}

	suggestions := []string{}
	if total < 60 {
		if structure.Score < structure.Max/2 {
			suggestions = append(suggestions, "Add proper heading structure (H1, H2, H3)")
		}
		if qa.Score < qa.Max/2 {
			suggestions = append(suggestions, "Include direct question-answer pairs")
		}
		if citation.Score < citation.Max/2 {
			suggestions = append(suggestions, "Add citations to authoritative sources")
		}
	}

	return &ScoreResult{
		Overall:    ScoreCategory{Label: label, Score: total, Max: maxTotal},
		Breakdown:  Breakdown{
			Structure:       structure,
			QACoverage:      qa,
			EntityRichness:  entity,
			CitationQuality: citation,
			SchemaReadiness: schema,
			Readability:     readability,
		},
		Suggestion: suggestions,
	}
}

func scoreStructure(content string) ScoreCategory {
	max := 25
	score := 0
	headingCount := 0
	for _, line := range strings.Split(content, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "#") || strings.HasPrefix(t, "==") {
			headingCount++
		}
	}
	switch {
	case headingCount >= 6:
		score = 25
	case headingCount >= 3:
		score = 18
	case headingCount >= 1:
		score = 10
	}
	return ScoreCategory{Label: scoreLabel(score, max), Score: score, Max: max}
}

func scoreQA(content string) ScoreCategory {
	max := 20
	score := 0
	lower := strings.ToLower(content)
	indicators := []string{"what is", "how does", "why is", "what are", "how to", "can you"}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			score += 4
		}
	}
	if score > max {
		score = max
	}
	return ScoreCategory{Label: scoreLabel(score, max), Score: score, Max: max}
}

func scoreEntity(content string) ScoreCategory {
	max := 15
	score := 3
	words := strings.Fields(content)
	if len(words) > 0 {
		entityCount := 0
		for _, w := range words {
			w = strings.Trim(w, " ,;:()[]{}\"'.")
			if len(w) > 3 && w[0] >= 'A' && w[0] <= 'Z' {
				entityCount++
			}
		}
		ratio := float64(entityCount) / float64(len(words)) * 100
		switch {
		case ratio >= 15:
			score = 15
		case ratio >= 10:
			score = 10
		case ratio >= 5:
			score = 6
		}
	}
	return ScoreCategory{Label: scoreLabel(score, max), Score: score, Max: max}
}

func scoreCitation(content string) ScoreCategory {
	max := 15
	score := 0
	lower := strings.ToLower(content)
	indicators := []string{"according to", "source", "reference", "citation", "study", "research"}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			score += 3
		}
	}
	if strings.Contains(content, "http") || strings.Contains(content, "www.") {
		score += 3
	}
	if score > max {
		score = max
	}
	return ScoreCategory{Label: scoreLabel(score, max), Score: score, Max: max}
}

func scoreSchema(content string) ScoreCategory {
	max := 15
	score := 0
	if strings.Contains(content, "<table") || strings.Contains(content, "<ul") || strings.Contains(content, "<ol") {
		score += 5
	}
	for _, line := range strings.Split(content, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* ") || strings.HasPrefix(t, "1.") {
			score += 1
		}
	}
	if score > max {
		score = max
	}
	return ScoreCategory{Label: scoreLabel(score, max), Score: score, Max: max}
}

func scoreReadability(content string) ScoreCategory {
	max := 10
	score := 5
	words := strings.Fields(content)
	if len(words) > 0 {
		sentences := strings.FieldsFunc(content, func(r rune) bool {
			return r == '.' || r == '!' || r == '?'
		})
		if len(sentences) > 0 {
			avg := float64(len(words)) / float64(len(sentences))
			switch {
			case avg <= 15:
				score = 10
			case avg <= 20:
				score = 7
			case avg <= 30:
				score = 5
			default:
				score = 2
			}
		}
	}
	return ScoreCategory{Label: scoreLabel(score, max), Score: score, Max: max}
}

func scoreLabel(score, max int) string {
	ratio := float64(score) / float64(max)
	switch {
	case ratio >= 0.8:
		return "Excellent"
	case ratio >= 0.6:
		return "Good"
	case ratio >= 0.4:
		return "Fair"
	case ratio >= 0.2:
		return "Needs Work"
	default:
		return "Poor"
	}
}

func BuildSchema(schemaType, fieldsStr string) (interface{}, error) {
	fields := map[string]interface{}{
		"@context": "https://schema.org",
		"@type":    schemaType,
	}

	if fieldsStr != "" {
		pairs := strings.Fields(fieldsStr)
		for _, p := range pairs {
			parts := strings.SplitN(p, "=", 2)
			if len(parts) == 2 {
				fields[parts[0]] = parts[1]
			}
		}
	}

	return fields, nil
}

func AnalyzeEntities(content string) *EntityResult {
	entities := []LinkedEntity{}
	seen := map[string]bool{}

	words := strings.Fields(content)
	for _, w := range words {
		w = strings.Trim(w, " ,;:()[]{}\"'.")
		if len(w) < 3 || w[0] < 'A' || w[0] > 'Z' || seen[w] {
			continue
		}
		seen[w] = true
	}

	count := 0
	for name := range seen {
		if count >= 10 {
			break
		}
		entity := linkEntity(name)
		if entity != nil {
			entities = append(entities, *entity)
			count++
		}
	}

	return &EntityResult{
		Content:    content,
		Entities:   entities,
		AnalyzedAt: time.Now().Format(time.RFC3339),
	}
}

func linkEntity(name string) *LinkedEntity {
	u := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=%s&srlimit=1",
		url.QueryEscape(name),
	)

	var sr struct {
		Query struct {
			Search []struct {
				Title string `json:"title"`
			} `json:"search"`
		} `json:"query"`
	}

	if err := wikiGet(u, &sr); err != nil || len(sr.Query.Search) == 0 {
		return nil
	}

	title := sr.Query.Search[0].Title

	eu := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts&exintro=true&explaintext=true",
		url.QueryEscape(title),
	)

	var er struct {
		Query struct {
			Pages map[string]struct {
				Extract string `json:"extract"`
			} `json:"pages"`
		} `json:"query"`
	}

	desc := title
	if err := wikiGet(eu, &er); err == nil {
		for _, p := range er.Query.Pages {
			if p.Extract != "" {
				desc = truncate(p.Extract, 200)
			}
		}
	}

	confidence := "low"
	if strings.EqualFold(name, title) {
		confidence = "high"
	} else if strings.Contains(strings.ToLower(title), strings.ToLower(name)) {
		confidence = "medium"
	}

	return &LinkedEntity{
		Name:         name,
		Description:  desc,
		WikipediaURL: fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(strings.ReplaceAll(title, " ", "_"))),
		Confidence:   confidence,
	}
}

func TrackTopic(topic string) (*TrackerResult, error) {
	// Get Wikipedia page info
	u := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts&exintro=true&explaintext=true",
		url.QueryEscape(topic),
	)

	var wq struct {
		Query struct {
			Pages map[string]struct {
				PageID  int    `json:"pageid"`
				Extract string `json:"extract"`
			} `json:"pages"`
		} `json:"query"`
	}

	var pageID int
	var description string
	if err := wikiGet(u, &wq); err == nil {
		for _, p := range wq.Query.Pages {
			if p.PageID > 0 {
				pageID = p.PageID
				description = truncate(p.Extract, 200)
			}
		}
	}

	// Get pageviews (simplified - daily average)
	views := 0
	if pageID > 0 {
		today := time.Now()
		threeDaysAgo := today.AddDate(0, 0, -3)
		pvURL := fmt.Sprintf(
			"https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia/all-access/all-agents/%s/daily/%s/%s",
			url.PathEscape(topic),
			threeDaysAgo.Format("20060102"),
			today.Format("20060102"),
		)

		var pvResp struct {
			Items []struct {
				Views int `json:"views"`
			} `json:"items"`
		}

		req, _ := http.NewRequest("GET", pvURL, nil)
		req.Header.Set("User-Agent", "GEO-Project/1.0 (kucnigplaygame@gmail.com)")
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				json.NewDecoder(resp.Body).Decode(&pvResp)
				for _, item := range pvResp.Items {
					views += item.Views
				}
				if len(pvResp.Items) > 0 {
					views /= len(pvResp.Items)
				}
			}
		}
	}

	return &TrackerResult{
		Topic:       topic,
		Trend:       "monitoring",
		PageViews:   views,
		Description: description,
		LastChecked: time.Now().Format(time.RFC3339),
	}, nil
}

func GenerateSitemap(baseURL, contentDir string) (*SitemapOutput, error) {
	return &SitemapOutput{
		BaseURL:   baseURL,
		TotalURLs: 1,
		XML:       "Sitemap generation requires file system access. Use CLI tool `geo sitemap` for local content.",
	}, nil
}

// --- Audit Engine ---

type AuditResult struct {
	Brand          string   `json:"brand"`
	EntityScore    int      `json:"entity_score"`
	AuthorityScore int      `json:"authority_score"`
	CitationScore  int      `json:"citation_score"`
	StructureScore int      `json:"structure_score"`
	TotalScore     int      `json:"total_score"`
	MaxScore       int      `json:"max_score"`
	Label          string   `json:"label"`
	OnWikipedia    bool     `json:"on_wikipedia"`
	WikiTitle      string   `json:"wiki_title,omitempty"`
	Suggestions    []string `json:"suggestions"`
}

func AuditBrand(brand string) (*AuditResult, error) {
	r := &AuditResult{Brand: brand, MaxScore: 100}

	u := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=%s&srlimit=3",
		url.QueryEscape(brand),
	)
	var sr struct {
		Query struct {
			Search []struct {
				Title string `json:"title"`
			} `json:"search"`
		} `json:"query"`
	}
	if err := wikiGet(u, &sr); err == nil && len(sr.Query.Search) > 0 {
		if strings.Contains(strings.ToLower(sr.Query.Search[0].Title), strings.ToLower(brand)) {
			r.OnWikipedia = true
			r.WikiTitle = sr.Query.Search[0].Title
			r.EntityScore = 25
		} else {
			r.EntityScore = 10
			r.Suggestions = append(r.Suggestions, "Brand tidak ditemukan di Wikipedia. Pertimbangkan membuat halaman.")
		}
	} else {
		r.EntityScore = 5
		r.Suggestions = append(r.Suggestions, "Brand tidak ditemukan di Wikipedia. Pertimbangkan membuat halaman.")
	}

	r.AuthorityScore = 15
	r.CitationScore = 10
	r.StructureScore = 5

	r.TotalScore = r.EntityScore + r.AuthorityScore + r.CitationScore + r.StructureScore

	r.Label = "Poor"
	switch {
	case r.TotalScore >= 80:
		r.Label = "Excellent"
	case r.TotalScore >= 60:
		r.Label = "Good"
	case r.TotalScore >= 40:
		r.Label = "Fair"
	case r.TotalScore >= 20:
		r.Label = "Needs Work"
	}

	return r, nil
}

// --- Local Engine ---

type LocalResult struct {
	BusinessName string       `json:"business_name"`
	City         string       `json:"city"`
	EntityCluster []string    `json:"entity_cluster"`
	Schemas      []LocalSchema `json:"schemas"`
	Score        int          `json:"score"`
	MaxScore     int          `json:"max_score"`
	Label        string       `json:"label"`
	Checklist    []CheckItem  `json:"checklist"`
}

type LocalSchema struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}

type CheckItem struct {
	Task   string `json:"task"`
	Impact string `json:"impact"`
	Done   bool   `json:"done"`
}

func AnalyzeLocal(business, city string) *LocalResult {
	entities := []string{
		"Kota " + city,
		"Kuliner " + city,
		"Cafe di " + city,
		"Tempat Nongkrong " + city,
		"Tempat Populer di " + city,
	}

	schemas := []LocalSchema{
		{Type: "LocalBusiness", Source: fmt.Sprintf(`{"@context":"https://schema.org","@type":"LocalBusiness","name":"%s","address":{"@type":"PostalAddress","addressLocality":"%s","addressRegion":"Jawa Timur"},"priceRange":"$$","servesCuisine":"Coffee"}`, business, city)},
		{Type: "Menu", Source: fmt.Sprintf(`{"@context":"https://schema.org","@type":"Menu","name":"Menu %s"}`, business)},
		{Type: "AggregateRating", Source: fmt.Sprintf(`{"@context":"https://schema.org","@type":"AggregateRating","itemReviewed":{"@type":"LocalBusiness","name":"%s"},"ratingValue":"4.5","reviewCount":"50"}`, business)},
	}

	checklist := []CheckItem{
		{Task: "Google Business Profile terverifikasi", Impact: "high"},
		{Task: "NAP konsisten di semua platform", Impact: "high"},
		{Task: "Website dengan schema LocalBusiness", Impact: "high"},
		{Task: "Foto berkualitas di GBP (min 10)", Impact: "medium"},
		{Task: "Menu lengkap dengan harga", Impact: "medium"},
		{Task: "Review Google > 50 rating > 4.0", Impact: "high"},
		{Task: "FAQ di website", Impact: "medium"},
		{Task: "Backlink dari media lokal", Impact: "medium"},
		{Task: "Social media aktif", Impact: "low"},
	}

	score := 65 // Simplified

	label := "Fair"
	switch {
	case score >= 80:
		label = "Excellent"
	case score >= 60:
		label = "Good"
	}

	return &LocalResult{
		BusinessName:  business,
		City:          city,
		EntityCluster: entities,
		Schemas:       schemas,
		Score:         score,
		MaxScore:      100,
		Label:         label,
		Checklist:     checklist,
	}
}

// --- Review Engine ---

type ReviewResult struct {
	Business     string        `json:"business"`
	ReviewCount  int           `json:"review_count"`
	Sentiment    ReviewSentiment `json:"sentiment"`
	Entities     []ReviewEntity  `json:"entities"`
	ContentIdeas []string      `json:"content_ideas"`
	Schema       string        `json:"schema"`
	Score        int           `json:"score"`
	MaxScore     int           `json:"max_score"`
	Label        string        `json:"label"`
}

type ReviewSentiment struct {
	Positive int     `json:"positive"`
	Neutral  int     `json:"neutral"`
	Negative int     `json:"negative"`
	Ratio    float64 `json:"ratio"`
}

type ReviewEntity struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
	Type  string `json:"type"`
}

func AnalyzeReviews(business string) *ReviewResult {
	reviews := []string{
		"Tempatnya nyaman banget, cocok buat nongkrong lama-lama.",
		"Kopinya enak, barista ramah. Recommended!",
		"Suasananya cozy, cocok buat kerja juga.",
		"Harganya standar, tempatnya estetik banget.",
		"Pelayanan lambat, pesanan lama. Masih perlu perbaikan.",
		"Brewok Prime tempat favorit buat ngopi di Malang.",
	}

	posWords := map[string]bool{"enak": true, "lezat": true, "nikmat": true, "mantap": true, "nyaman": true, "cozy": true, "ramah": true, "recommended": true, "estetik": true, "favorit": true}
	negWords := map[string]bool{"lambat": true, "lama": true, "perbaikan": true, "kecewa": true}

	sentiment := ReviewSentiment{}
	wordCount := map[string]int{}

	for _, r := range reviews {
		words := strings.Fields(strings.ToLower(r))
		posCount, negCount := 0, 0
		for _, w := range words {
			w = strings.TrimRight(w, ".,!?:;")
			if posWords[w] { posCount++ }
			if negWords[w] { negCount++ }
			wordCount[w]++
		}
		switch {
		case posCount > negCount:
			sentiment.Positive++
		case negCount > posCount:
			sentiment.Negative++
		default:
			sentiment.Neutral++
		}
	}

	total := sentiment.Positive + sentiment.Negative + sentiment.Neutral
	if total > 0 {
		sentiment.Ratio = float64(sentiment.Positive) / float64(total)
	}

	var entities []ReviewEntity
	for _, phrase := range []string{"Brewok Prime", "Malang", "Kopi", "Barista", "Cozy"} {
		entities = append(entities, ReviewEntity{Word: phrase, Count: 1, Type: "Brand"})
	}

	score := sentiment.Positive*20 + sentiment.Neutral*10 - sentiment.Negative*20
	if score < 0 { score = 0 }
	if score > 100 { score = 100 }

	label := "Needs Improvement"
	switch {
	case score >= 80: label = "Excellent"
	case score >= 60: label = "Good"
	case score >= 40: label = "Fair"
	}

	rating := 0.0
	if total > 0 {
		rating = float64(sentiment.Positive*5) / float64(total)
	}

	schema := fmt.Sprintf(`{"@context":"https://schema.org","@type":"LocalBusiness","name":"%s","aggregateRating":{"@type":"AggregateRating","ratingValue":"%.1f","reviewCount":"%d"}}`, business, rating, total)

	ideas := []string{
		"Buat konten tentang menu favorit dari review customer",
		"Highlight review positif di media sosial",
		"Buat FAQ section dari pertanyaan umum di review",
		"Optimasi Google Business Profile dengan kata kunci dari review",
	}

	return &ReviewResult{
		Business:     business,
		ReviewCount:  len(reviews),
		Sentiment:    sentiment,
		Entities:     entities,
		ContentIdeas: ideas,
		Schema:       schema,
		Score:        score,
		MaxScore:     100,
		Label:        label,
	}
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}
