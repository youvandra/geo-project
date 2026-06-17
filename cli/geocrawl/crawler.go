package geocrawl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CrawlResult struct {
	Query      string      `json:"query"`
	Sources    []SourceHit `json:"sources"`
	TotalHits  int         `json:"total_hits"`
	Score      int         `json:"score"`
	MaxScore   int         `json:"max_score"`
	Label      string      `json:"label"`
	CrawledAt  time.Time   `json:"crawled_at"`
}

type SourceHit struct {
	Source   string `json:"source"`
	Found    bool   `json:"found"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
}

type wikiSearchResponse struct {
	Query struct {
		Search []struct {
			Title  string `json:"title"`
			Snippet string `json:"snippet"`
		} `json:"search"`
	} `json:"query"`
}

type wikiPageResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID  int    `json:"pageid"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

type wikidataResponse struct {
	Results []struct {
		Item struct {
			ID    string `json:"id"`
			Label string `json:"label"`
		} `json:"item"`
	} `json:"results"`
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

func Crawl(query string) (*CrawlResult, error) {
	result := &CrawlResult{
		Query:     query,
		MaxScore:  100,
		CrawledAt: time.Now(),
	}

	// 1. Wikipedia search
	wikiHit := checkWikipedia(query)
	result.Sources = append(result.Sources, wikiHit)

	// 2. Wikipedia direct page
	wikiPageHit := checkWikipediaPage(query)
	result.Sources = append(result.Sources, wikiPageHit)

	// 3. Wikidata search
	wdHit := checkWikidata(query)
	result.Sources = append(result.Sources, wdHit)

	// 4. Title-based scoring (break down query into parts)
	partsHit := checkParts(query)
	result.Sources = append(result.Sources, partsHit)

	// Calculate score
	score := 0
	hits := 0
	for _, s := range result.Sources {
		if s.Found {
			hits++
		}
	}
	result.TotalHits = hits

	switch {
	case hits >= 4:
		score = 90
	case hits == 3:
		score = 70
	case hits == 2:
		score = 50
	case hits == 1:
		score = 25
	}
	result.Score = score

	label := "Not Found"
	switch {
	case score >= 80:
		label = "Strong Presence"
	case score >= 60:
		label = "Moderate Presence"
	case score >= 40:
		label = "Weak Presence"
	case score >= 20:
		label = "Minimal Presence"
	}
	result.Label = label

	return result, nil
}

func checkWikipedia(query string) SourceHit {
	u := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=%s&srlimit=3&srprop=snippet",
		url.QueryEscape(query),
	)
	var sr wikiSearchResponse
	if err := wikiGet(u, &sr); err != nil {
		return SourceHit{Source: "wikipedia-search", Found: false}
	}
	for _, s := range sr.Query.Search {
		if scoreMatch(query, s.Title) {
			return SourceHit{
				Source:  "wikipedia-search",
				Found:   true,
				Title:   s.Title,
				URL:     fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(strings.ReplaceAll(s.Title, " ", "_"))),
				Snippet: truncate(s.Snippet, 150),
			}
		}
	}
	return SourceHit{Source: "wikipedia-search", Found: false}
}

func checkWikipediaPage(query string) SourceHit {
	u := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts&exintro=true&explaintext=true",
		url.QueryEscape(query),
	)
	var pr wikiPageResponse
	if err := wikiGet(u, &pr); err != nil {
		return SourceHit{Source: "wikipedia-direct", Found: false}
	}
	for _, p := range pr.Query.Pages {
		if p.PageID > 0 {
			return SourceHit{
				Source:  "wikipedia-direct",
				Found:   true,
				Title:   p.Title,
				URL:     fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(strings.ReplaceAll(p.Title, " ", "_"))),
				Snippet: truncate(p.Extract, 150),
			}
		}
	}
	return SourceHit{Source: "wikipedia-direct", Found: false}
}

func checkWikidata(query string) SourceHit {
	sparql := url.QueryEscape(fmt.Sprintf(
		`SELECT ?item ?itemLabel WHERE { ?item wdt:P31 wd:Q5; rdfs:label "%s"@en. SERVICE wikibase:label { bd:serviceParam wikibase:language "en". } } LIMIT 5`,
		query,
	))
	u := fmt.Sprintf(
		"https://query.wikidata.org/sparql?format=json&query=%s", sparql,
	)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "GEO-Project/1.0 (kucnigplaygame@gmail.com)")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SourceHit{Source: "wikidata", Found: false}
	}
	defer resp.Body.Close()

	var wd wikidataResponse
	if err := json.NewDecoder(resp.Body).Decode(&wd); err != nil {
		return SourceHit{Source: "wikidata", Found: false}
	}
	for _, r := range wd.Results {
		if r.Item.Label != "" && scoreMatch(query, r.Item.Label) {
			return SourceHit{
				Source: "wikidata",
				Found:  true,
				Title:  r.Item.Label,
				URL:    fmt.Sprintf("https://www.wikidata.org/wiki/%s", r.Item.ID),
			}
		}
	}
	return SourceHit{Source: "wikidata", Found: false}
}

func checkParts(query string) SourceHit {
	parts := strings.Fields(query)
	found := 0
	for _, p := range parts {
		u := fmt.Sprintf(
			"https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=%s&srlimit=1",
			url.QueryEscape(p),
		)
		var sr wikiSearchResponse
		if err := wikiGet(u, &sr); err == nil && len(sr.Query.Search) > 0 {
			if scoreMatch(p, sr.Query.Search[0].Title) {
				found++
			}
		}
	}
	if found > 0 {
		return SourceHit{
			Source:  "part-match",
			Found:   true,
			Title:   fmt.Sprintf("%d/%d parts found on Wikipedia", found, len(parts)),
			Snippet: fmt.Sprintf("Query parts matched %d of %d terms on Wikipedia", found, len(parts)),
		}
	}
	return SourceHit{Source: "part-match", Found: false}
}

func scoreMatch(query, title string) bool {
	ql := strings.ToLower(query)
	tl := strings.ToLower(title)
	return strings.Contains(tl, ql) || strings.Contains(ql, tl)
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}
