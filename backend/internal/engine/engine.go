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

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}
