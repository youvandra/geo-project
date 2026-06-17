package cmdsitemap

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SitemapIndex struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type SitemapOptions struct {
	BaseURL     string
	ContentDir  string
	ChangeFreq  string
	Priority    string
	FileExts    []string
	ExcludeDirs []string
}

type SitemapResult struct {
	BaseURL    string        `json:"base_url"`
	TotalURLs  int           `json:"total_urls"`
	URLs       []SitemapURL  `json:"urls"`
	GeneratedAt time.Time    `json:"generated_at"`
}

func Generate(opts SitemapOptions) (*SitemapResult, error) {
	if opts.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if opts.ContentDir == "" {
		return nil, fmt.Errorf("content directory is required")
	}

	if len(opts.FileExts) == 0 {
		opts.FileExts = []string{".md", ".html", ".htm"}
	}

	if opts.ChangeFreq == "" {
		opts.ChangeFreq = "weekly"
	}

	if opts.Priority == "" {
		opts.Priority = "0.5"
	}

	excludeMap := map[string]bool{}
	for _, d := range opts.ExcludeDirs {
		excludeMap[d] = true
	}

	var urls []SitemapURL

	err := filepath.Walk(opts.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if excludeMap[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		validExt := false
		for _, e := range opts.FileExts {
			if ext == e {
				validExt = true
				break
			}
		}
		if !validExt {
			return nil
		}

		relPath, _ := filepath.Rel(opts.ContentDir, path)
		relPath = strings.TrimSuffix(relPath, ext)

		if strings.HasPrefix(relPath, "_") || strings.HasPrefix(filepath.Base(relPath), ".") {
			return nil
		}

		urlPath := "/" + strings.ReplaceAll(relPath, string(os.PathSeparator), "/")
		if strings.HasSuffix(urlPath, "/index") {
			urlPath = strings.TrimSuffix(urlPath, "/index")
		}
		if urlPath == "/"+strings.TrimPrefix(opts.ContentDir, ".") {
			urlPath = "/"
		}

		loc := strings.TrimRight(opts.BaseURL, "/") + urlPath
		if ext == ".md" {
			loc = strings.TrimSuffix(loc, ".md")
		}

		priority := calculatePriority(relPath)

		urls = append(urls, SitemapURL{
			Loc:        loc,
			LastMod:    info.ModTime().Format("2006-01-02"),
			ChangeFreq: opts.ChangeFreq,
			Priority:   priority,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk content dir: %w", err)
	}

	return &SitemapResult{
		BaseURL:    opts.BaseURL,
		TotalURLs:  len(urls),
		URLs:       urls,
		GeneratedAt: time.Now(),
	}, nil
}

func GenerateXML(opts SitemapOptions) (string, error) {
	result, err := Generate(opts)
	if err != nil {
		return "", err
	}

	index := SitemapIndex{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  result.URLs,
	}

	data, err := xml.MarshalIndent(index, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(data), nil
}

func calculatePriority(path string) string {
	depth := len(strings.Split(path, string(os.PathSeparator)))

	switch {
	case depth <= 1:
		return "0.9"
	case depth == 2:
		return "0.7"
	case depth == 3:
		return "0.5"
	default:
		return "0.3"
	}
}
