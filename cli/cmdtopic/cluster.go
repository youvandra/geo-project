package cmdtopic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TopicResult struct {
	Topic       string    `json:"topic"`
	Description string    `json:"description"`
	Entities    []Entity  `json:"entities"`
	Questions   []string  `json:"questions"`
	Subtopics   []string  `json:"subtopics"`
	AnalyzedAt  time.Time `json:"analyzed_at"`
}

type Entity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Relevance   string `json:"relevance"`
}

type wikiPage struct {
	PageID  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

type wikiQuery struct {
	Query struct {
		Pages map[string]wikiPage `json:"pages"`
	} `json:"query"`
}

type wikiLinks struct {
	Query struct {
		Pages map[string]struct {
			Links []struct {
				Title string `json:"title"`
			} `json:"links"`
		} `json:"pages"`
	} `json:"query"`
}

type wikiCategories struct {
	Query struct {
		Pages map[string]struct {
			Categories []struct {
				Title string `json:"title"`
			} `json:"categories"`
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

func Analyze(topic string) (*TopicResult, error) {
	page, err := fetchPage(topic)
	if err != nil {
		return nil, fmt.Errorf("fetch topic: %w", err)
	}

	entities := fetchEntities(page.Title)
	questions := generateQuestions(topic, entities)
	subtopics := extractSubtopics(page.Extract)

	return &TopicResult{
		Topic:       topic,
		Description: page.Extract,
		Entities:    entities,
		Questions:   questions,
		Subtopics:   subtopics,
		AnalyzedAt:  time.Now(),
	}, nil
}

func fetchPage(topic string) (*wikiPage, error) {
	u := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts&explaintext=true&redirects=1",
		url.QueryEscape(topic),
	)
	var wq wikiQuery
	if err := wikiGet(u, &wq); err != nil {
		return nil, err
	}
	for _, p := range wq.Query.Pages {
		if p.PageID == 0 {
			continue
		}
		return &p, nil
	}
	return nil, fmt.Errorf("no Wikipedia page found for %q", topic)
}

func fetchEntities(title string) []Entity {
	entities := []Entity{}
	seen := map[string]bool{}

	var wl wikiLinks
	linksURL := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=links&pllimit=20",
		url.QueryEscape(title),
	)
	if err := wikiGet(linksURL, &wl); err == nil {
		for _, p := range wl.Query.Pages {
			for _, link := range p.Links {
				if !seen[link.Title] {
					seen[link.Title] = true
					entities = append(entities, Entity{
						Name:      link.Title,
						Relevance: "related",
					})
				}
			}
		}
	}

	var wc wikiCategories
	catURL := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=categories&cllimit=20",
		url.QueryEscape(title),
	)
	if err := wikiGet(catURL, &wc); err == nil {
		for _, p := range wc.Query.Pages {
			for _, cat := range p.Categories {
				name := strings.TrimPrefix(cat.Title, "Category:")
				if !seen[name] {
					seen[name] = true
					entities = append(entities, Entity{
						Name:      name,
						Relevance: "category",
					})
				}
			}
		}
	}

	if len(entities) > 15 {
		entities = entities[:15]
	}
	return entities
}

func extractSubtopics(extract string) []string {
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

func generateQuestions(topic string, entities []Entity) []string {
	questions := []string{
		fmt.Sprintf("What is %s?", topic),
		fmt.Sprintf("How does %s work?", topic),
		fmt.Sprintf("Why is %s important?", topic),
		fmt.Sprintf("What are the key components of %s?", topic),
		fmt.Sprintf("How to get started with %s?", topic),
		fmt.Sprintf("What are common challenges in %s?", topic),
		fmt.Sprintf("What is the future of %s?", topic),
	}
	for _, e := range entities {
		questions = append(questions, fmt.Sprintf("What is the relationship between %s and %s?", topic, e.Name))
	}
	return questions
}
