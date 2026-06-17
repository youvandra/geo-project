package cmdentity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type EntityResult struct {
	Content   string           `json:"content"`
	Entities  []LinkedEntity   `json:"entities"`
	AnalyzedAt time.Time       `json:"analyzed_at"`
}

type LinkedEntity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WikipediaURL string `json:"wikipedia_url"`
	WikidataID  string `json:"wikidata_id,omitempty"`
	Confidence  string `json:"confidence"`
}

type wikiSearchResponse struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			PageID  int    `json:"pageid"`
			Snippet string `json:"snippet"`
		} `json:"search"`
	} `json:"query"`
}

type wikiExtractResponse struct {
	Query struct {
		Pages map[string]struct {
			Extract  string `json:"extract"`
			PageID   int    `json:"pageid"`
			Title    string `json:"title"`
		} `json:"pages"`
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

func Analyze(content string) *EntityResult {
	candidates := extractCandidateEntities(content)

	var mu sync.Mutex
	var entities []LinkedEntity
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for _, c := range candidates {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			entity := linkEntity(name)
			if entity != nil {
				mu.Lock()
				entities = append(entities, *entity)
				mu.Unlock()
			}
		}(c)
	}
	wg.Wait()

	if len(entities) > 20 {
		entities = entities[:20]
	}

	return &EntityResult{
		Content:    content,
		Entities:   entities,
		AnalyzedAt: time.Now(),
	}
}

func extractCandidateEntities(content string) []string {
	seen := map[string]bool{}
	candidates := []string{}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		words := strings.Fields(line)
		for i := 0; i < len(words); i++ {
			w := strings.Trim(words[i], " ,;:()[]{}\"'.")
			if len(w) < 3 {
				continue
			}

			if isUpperCase(rune(w[0])) {
				phrase := w
				for j := i + 1; j < len(words); j++ {
					next := strings.Trim(words[j], " ,;:()[]{}\"'.")
					if len(next) > 0 && isUpperCase(rune(next[0])) {
						phrase += " " + next
						i = j
					} else {
						break
					}
				}
				if !seen[phrase] && len(phrase) > 2 {
					seen[phrase] = true
					candidates = append(candidates, phrase)
				}
			}
		}
	}

	if len(candidates) > 30 {
		candidates = candidates[:30]
	}
	return candidates
}

func linkEntity(name string) *LinkedEntity {
	searchURL := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=%s&srlimit=1&srprop=snippet",
		url.QueryEscape(name),
	)

	var sr wikiSearchResponse
	if err := wikiGet(searchURL, &sr); err != nil {
		return nil
	}

	if len(sr.Query.Search) == 0 {
		return nil
	}

	topResult := sr.Query.Search[0]
	titleSimilarity := scoreSimilarity(name, topResult.Title)

	if titleSimilarity < 0.3 {
		return nil
	}

	confidence := "low"
	switch {
	case titleSimilarity >= 0.9:
		confidence = "high"
	case titleSimilarity >= 0.6:
		confidence = "medium"
	}

	extractURL := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts&exintro=true&explaintext=true",
		url.QueryEscape(topResult.Title),
	)

	var er wikiExtractResponse
	description := topResult.Title
	if err := wikiGet(extractURL, &er); err == nil {
		for _, p := range er.Query.Pages {
			if p.Extract != "" {
				description = truncateText(p.Extract, 200)
			}
		}
	}

	wikiURL := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", url.PathEscape(strings.ReplaceAll(topResult.Title, " ", "_")))

	return &LinkedEntity{
		Name:         name,
		Description:  description,
		WikipediaURL: wikiURL,
		Confidence:   confidence,
	}
}

func scoreSimilarity(a, b string) float64 {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	if a == b {
		return 1.0
	}

	if strings.Contains(a, b) || strings.Contains(b, a) {
		return 0.85
	}

	aWords := strings.Fields(a)
	bWords := strings.Fields(b)

	if len(aWords) == 0 || len(bWords) == 0 {
		return 0
	}

	matchCount := 0
	for _, wa := range aWords {
		for _, wb := range bWords {
			if wa == wb {
				matchCount++
				break
			}
		}
	}

	maxWords := len(aWords)
	if len(bWords) > maxWords {
		maxWords = len(bWords)
	}

	return float64(matchCount) / float64(maxWords)
}

func truncateText(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

func isUpperCase(r rune) bool {
	return r >= 'A' && r <= 'Z'
}
