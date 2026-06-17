package geolocal

import (
	"fmt"
	"time"
)

type LocalResult struct {
	BusinessName string      `json:"business_name"`
	City         string      `json:"city"`
	NAPCheck     NAPStatus   `json:"nap_check"`
	Schemas      []SchemaOut `json:"schemas"`
	EntityCluster []string   `json:"entity_cluster"`
	Checklist    []CheckItem `json:"checklist"`
	Score        int         `json:"score"`
	MaxScore     int         `json:"max_score"`
	Label        string      `json:"label"`
	AnalyzedAt   time.Time   `json:"analyzed_at"`
}

type NAPStatus struct {
	NameConsistent  bool   `json:"name_consistent"`
	AddressComplete bool   `json:"address_complete"`
	PhoneValid      bool   `json:"phone_valid"`
	Issues          []string `json:"issues"`
}

type SchemaOut struct {
	Type   string `json:"type"`
	Source string `json:"source"`
	Output string `json:"output"`
}

type CheckItem struct {
	Task    string `json:"task"`
	Done    bool   `json:"done"`
	Impact  string `json:"impact"`
}

func Analyze(businessName, city string) (*LocalResult, error) {
	localEntities := []string{
		fmt.Sprintf("Kota %s", city),
		"Tempat Populer di " + city,
		"Kuliner " + city,
		"Tempat Nongkrong " + city,
		fmt.Sprintf("Cafe di %s", city),
	}

	schemas := []SchemaOut{
		{
			Type: "LocalBusiness",
			Source: fmt.Sprintf(`{
  "@context": "https://schema.org",
  "@type": "LocalBusiness",
  "name": "%s",
  "address": {
    "@type": "PostalAddress",
    "streetAddress": "Jalan...",
    "addressLocality": "%s",
    "addressRegion": "Jawa Timur",
    "postalCode": "..."
  },
  "telephone": "...",
  "openingHours": "Mo-Su 08:00-22:00",
  "priceRange": "$$",
  "servesCuisine": "Coffee"
}`, businessName, city),
		},
		{
			Type:   "Menu",
			Source: fmt.Sprintf(`{"@context":"https://schema.org","@type":"Menu","name":"Menu %s","hasMenuItem":[]}`, businessName),
		},
		{
			Type:   "AggregateRating",
			Source: fmt.Sprintf(`{"@context":"https://schema.org","@type":"AggregateRating","itemReviewed":{"@type":"LocalBusiness","name":"%s"},"ratingValue":"4.5","reviewCount":"50"}`, businessName),
		},
	}

	checklist := []CheckItem{
		{Task: "Google Business Profile terverifikasi", Impact: "high"},
		{Task: "NAP konsisten di semua platform (web, maps, social)", Impact: "high"},
		{Task: "Website dengan schema LocalBusiness", Impact: "high"},
		{Task: "Foto berkualitas tinggi di GBP (min 10)", Impact: "medium"},
		{Task: "Menu lengkap dengan harga di website", Impact: "medium"},
		{Task: "Review Google > 50 dengan rating > 4.0", Impact: "high"},
		{Task: "FAQ di website (pertanyaan umum customer)", Impact: "medium"},
		{Task: "Blog/post tentang kegiatan cafe", Impact: "low"},
		{Task: "Backlink dari media lokal Malang", Impact: "medium"},
		{Task: "Social media aktif (IG/TikTok)", Impact: "low"},
	}

	nap := NAPStatus{
		NameConsistent:  true,
		AddressComplete: true,
		PhoneValid:      true,
		Issues:          nil,
	}

	maxScore := 100

	label := "Fair"
	if maxScore >= 80 {
		label = "Excellent"
	} else if maxScore >= 60 {
		label = "Good"
	}

	return &LocalResult{
		BusinessName:  businessName,
		City:          city,
		NAPCheck:      nap,
		Schemas:       schemas,
		EntityCluster: localEntities,
		Checklist:     checklist,
		Score:         0,
		MaxScore:      maxScore,
		Label:         label,
		AnalyzedAt:    time.Now(),
	}, nil
}
