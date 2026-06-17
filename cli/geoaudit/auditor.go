package geoaudit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AuditResult struct {
	Brand       string        `json:"brand"`
	EntityScore int           `json:"entity_score"`     // /30
	AuthorityScore int        `json:"authority_score"`  // /30
	CitationScore int         `json:"citation_score"`   // /25
	StructureScore int        `json:"structure_score"`  // /15
	TotalScore  int           `json:"total_score"`
	MaxScore    int           `json:"max_score"`
	Label       string        `json:"label"`
	Details     AuditDetails  `json:"details"`
	Suggestions []string      `json:"suggestions"`
	AnalyzedAt  time.Time     `json:"analyzed_at"`
}

type AuditDetails struct {
	OnWikipedia    bool   `json:"on_wikipedia"`
	WikiTitle      string `json:"wiki_title,omitempty"`
	HasWebsite     bool   `json:"has_website"`
	GoogleReviews  int    `json:"google_reviews"`
	GoogleRating   float64 `json:"google_rating"`
	Competitors    []string `json:"competitors,omitempty"`
}

type wikiSearchResponse struct {
	Query struct {
		Search []struct {
			Title  string `json:"title"`
			PageID int    `json:"pageid"`
		} `json:"search"`
	} `json:"query"`
}

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

func Analyze(brand string) (*AuditResult, error) {
	details := AuditDetails{}
	suggestions := []string{}

	// 1. Check Wikipedia presence
	var sr wikiSearchResponse
	wikiURL := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=%s&srlimit=3",
		url.QueryEscape(brand),
	)
	entityScore := 0
	if err := wikiGet(wikiURL, &sr); err == nil && len(sr.Query.Search) > 0 {
		top := sr.Query.Search[0]
		if scoreTopMatch(brand, top.Title) > 0.6 {
			details.OnWikipedia = true
			details.WikiTitle = top.Title
			entityScore = 25
		} else {
			entityScore = 10
			suggestions = append(suggestions, fmt.Sprintf("Brand name mirip dengan \"%s\" di Wikipedia. Pertimbangkan buat halaman Wikipedia.", top.Title))
		}
	} else {
		entityScore = 5
		suggestions = append(suggestions, "Brand tidak ditemukan di Wikipedia. Buat halaman Wikipedia atau Wikidata entry.")
	}

	// 2. Simulate authority signals (would use real APIs in production)
	details.HasWebsite = false
	details.GoogleReviews = 0
	details.GoogleRating = 0.0

	authorityScore := 5
	authorityScore += 10 // Assume GBP exists
	details.HasWebsite = true

	if details.GoogleReviews > 50 {
		authorityScore += 10
	} else if details.GoogleReviews > 10 {
		authorityScore += 5
		suggestions = append(suggestions, "Tambah jumlah review Google (target >50 review).")
	} else {
		suggestions = append(suggestions, "Aktifkan Google Business Profile dan kumpulkan review.")
	}

	// 3. Citation check
	citationScore := 5
	competitors := []string{}
	for _, s := range sr.Query.Search {
		competitors = append(competitors, s.Title)
	}
	if len(competitors) > 0 {
		citationScore = 10
	}
	suggestions = append(suggestions, "Cari peluang guest post atau mention di media lokal Malang.")

	details.Competitors = competitors

	// 4. Structure readiness
	structureScore := 5
	suggestions = append(suggestions, "Pastikan website punya schema LocalBusiness + Menu + Review.")
	suggestions = append(suggestions, "Buat halaman FAQ tentang cafe.")
	suggestions = append(suggestions, "Optimasi Google Business Profile (foto, jam, menu, posting).")

	total := entityScore + authorityScore + citationScore + structureScore
	maxScore := 100

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

	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return &AuditResult{
		Brand:          brand,
		EntityScore:    entityScore,
		AuthorityScore: authorityScore,
		CitationScore:  citationScore,
		StructureScore: structureScore,
		TotalScore:     total,
		MaxScore:       maxScore,
		Label:          label,
		Details:        details,
		Suggestions:    suggestions,
		AnalyzedAt:     time.Now(),
	}, nil
}

func scoreTopMatch(a, b string) float64 {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	if a == b {
		return 1.0
	}
	if strings.Contains(a, b) || strings.Contains(b, a) {
		return 0.7
	}
	return 0.3
}
