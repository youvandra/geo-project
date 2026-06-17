package cmdtracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type TrackerResult struct {
	Topic       string       `json:"topic"`
	Current     TopicStats   `json:"current"`
	History     []TopicStats `json:"history"`
	Trend       string       `json:"trend"`
	LastChecked time.Time    `json:"last_checked"`
}

type TopicStats struct {
	Date         time.Time `json:"date"`
	PageViews    int       `json:"page_views"`
	WikiPageID   int       `json:"wiki_page_id,omitempty"`
	Description  string    `json:"description,omitempty"`
}

type pageviewResponse struct {
	Items []struct {
		Article  string `json:"article"`
		Views    int    `json:"views"`
		Date     string `json:"timestamp"`
	} `json:"items"`
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

var dataDir string

func init() {
	home, _ := os.UserHomeDir()
	dataDir = filepath.Join(home, ".geo-tracker")
	os.MkdirAll(dataDir, 0755)
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

func Analyze(topic string) (*TrackerResult, error) {
	pageID, description, err := getWikiPageInfo(topic)
	if err != nil {
		return nil, fmt.Errorf("get page info: %w", err)
	}

	views, err := getPageViews(topic)
	if err != nil {
		return nil, fmt.Errorf("get page views: %w", err)
	}

	current := TopicStats{
		Date:        time.Now(),
		PageViews:   views,
		WikiPageID:  pageID,
		Description: description,
	}

	history := loadHistory(topic)
	history = append(history, current)
	saveHistory(topic, history)

	trend := analyzeTrend(history)

	sort.Slice(history, func(i, j int) bool {
		return history[i].Date.Before(history[j].Date)
	})

	if len(history) > 30 {
		history = history[len(history)-30:]
	}

	return &TrackerResult{
		Topic:       topic,
		Current:     current,
		History:     history,
		Trend:       trend,
		LastChecked: time.Now(),
	}, nil
}

func getWikiPageInfo(topic string) (int, string, error) {
	u := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=extracts&exintro=true&explaintext=true",
		url.QueryEscape(topic),
	)

	var resp wikiPageResponse
	if err := wikiGet(u, &resp); err != nil {
		return 0, "", err
	}

	for _, p := range resp.Query.Pages {
		if p.PageID > 0 {
			desc := p.Extract
			if len([]rune(desc)) > 200 {
				desc = string([]rune(desc)[:200]) + "..."
			}
			return p.PageID, desc, nil
		}
	}

	return 0, "", fmt.Errorf("no Wikipedia page for %q", topic)
}

func getPageViews(topic string) (int, error) {
	today := time.Now()
	twoDaysAgo := today.AddDate(0, 0, -2)

	dateStr := twoDaysAgo.Format("20060102")
	u := fmt.Sprintf(
		"https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia/all-access/all-agents/%s/daily/%s/%s",
		url.PathEscape(topic),
		dateStr,
		today.Format("20060102"),
	)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "GEO-Project/1.0 (kucnigplaygame@gmail.com)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, nil
	}

	var pv pageviewResponse
	if err := json.NewDecoder(resp.Body).Decode(&pv); err != nil {
		return 0, err
	}

	total := 0
	for _, item := range pv.Items {
		total += item.Views
	}

	if len(pv.Items) > 0 {
		return total / len(pv.Items), nil
	}

	return 0, nil
}

func loadHistory(topic string) []TopicStats {
	safeName := strings.ReplaceAll(strings.ToLower(topic), " ", "_")
	path := filepath.Join(dataDir, safeName+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return []TopicStats{}
	}

	var history []TopicStats
	json.Unmarshal(data, &history)
	return history
}

func saveHistory(topic string, history []TopicStats) {
	safeName := strings.ReplaceAll(strings.ToLower(topic), " ", "_")
	path := filepath.Join(dataDir, safeName+".json")

	data, _ := json.MarshalIndent(history, "", "  ")
	os.WriteFile(path, data, 0644)
}

func analyzeTrend(history []TopicStats) string {
	if len(history) < 2 {
		return "insufficient data"
	}

	recent := history[len(history)-1].PageViews
	if len(history) >= 2 {
		previous := history[len(history)-2].PageViews
		diff := recent - previous

		switch {
		case diff > 100:
			return fmt.Sprintf("increasing (+%d views)", diff)
		case diff > 0:
			return fmt.Sprintf("slightly increasing (+%d views)", diff)
		case diff == 0:
			return "stable"
		case diff > -100:
			return fmt.Sprintf("slightly decreasing (%d views)", diff)
		default:
			return fmt.Sprintf("decreasing (%d views)", diff)
		}
	}

	return "stable"
}
