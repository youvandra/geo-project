package georeview

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

type ReviewResult struct {
	Business   string    `json:"business"`
	ReviewCount int      `json:"review_count"`
	Sentiment  Sentiment `json:"sentiment"`
	Entities   []Entity     `json:"entities"`
	Keywords   []Keyword   `json:"keywords"`
	ContentIdeas []string `json:"content_ideas"`
	Schema     string    `json:"schema"`
	Score      int       `json:"score"`
	MaxScore   int       `json:"max_score"`
	Label      string    `json:"label"`
	AnalyzedAt time.Time `json:"analyzed_at"`
}

type Sentiment struct {
	Positive int     `json:"positive"`
	Negative int     `json:"negative"`
	Neutral  int     `json:"neutral"`
	Ratio    float64 `json:"ratio"`
}

type Entity struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
	Type  string `json:"type"`
}

type Keyword struct {
	Phrase string `json:"phrase"`
	Count  int    `json:"count"`
	Score  int    `json:"score"`
}

var positiveWords = map[string]bool{
	"enak": true, "lezat": true, "nikmat": true, "mantap": true, "keren": true,
	"nyaman": true, "cozy": true, "asik": true, "seru": true, "ramah": true,
	"cepat": true, "murah": true, "worth": true, "recommended": true, "terbaik": true,
	"good": true, "great": true, "amazing": true, "delicious": true, "friendly": true,
	"fast": true, "affordable": true, "comfortable": true, "best": true, "love": true,
	"recommend": true, "perfect": true, "excellent": true, "awesome": true, "nice": true,
	"bersih": true, "wangi": true, "estetik": true, "instagramable": true, "cocok": true,
}

var negativeWords = map[string]bool{
	"mahal": true, "lambat": true, "lama": true, "kotor": true, "bau": true,
	"jelek": true, "buruk": true, "parah": true, "kecewa": true, "mengecewakan": true,
	"tidak": true, "nggak": true, "ga": true, "gak": true, " worst": true,
	"bad": true, "terrible": true, "horrible": true, "slow": true, "expensive": true,
	"dirty": true, "disappointed": true, "disappointing": true, "rude": true, "overprice": true,
	"menyesal": true, "gagal": true, "boring": true, "membosankan": true, "panas": true,
}

func Analyze(business string, reviews []string) (*ReviewResult, error) {
	if len(reviews) == 0 {
		reviews = getSampleReviews(business)
	}

	sentiment := analyzeSentiment(reviews)
	entities := extractEntities(reviews)
	keywords := extractKeywords(reviews)
	contentIdeas := generateContentIdea(entities, keywords)

	schema := buildReviewSchema(business, reviews, sentiment)

	total := sentiment.Positive*2 + sentiment.Neutral - sentiment.Negative*2
	switch {
	case total < 0:
		total = 0
	case total > 100:
		total = 100
	}

	label := "Needs Improvement"
	switch {
	case total >= 80:
		label = "Excellent"
	case total >= 60:
		label = "Good"
	case total >= 40:
		label = "Fair"
	}

	return &ReviewResult{
		Business:     business,
		ReviewCount:  len(reviews),
		Sentiment:    sentiment,
		Entities:     entities,
		Keywords:     keywords,
		ContentIdeas: contentIdeas,
		Schema:       schema,
		Score:        total,
		MaxScore:     100,
		Label:        label,
		AnalyzedAt:   time.Now(),
	}, nil
}

func analyzeSentiment(reviews []string) Sentiment {
	s := Sentiment{}
	for _, r := range reviews {
		words := strings.Fields(strings.ToLower(r))
		posCount := 0
		negCount := 0
		for _, w := range words {
			w = strings.TrimRight(w, ".,!?;:")
			if positiveWords[w] {
				posCount++
			}
			if negativeWords[w] {
				negCount++
			}
		}
		if posCount > negCount {
			s.Positive++
		} else if negCount > posCount {
			s.Negative++
		} else {
			s.Neutral++
		}
	}
	total := s.Positive + s.Negative + s.Neutral
	if total > 0 {
		s.Ratio = float64(s.Positive) / float64(total)
	}
	return s
}

func extractEntities(reviews []string) []Entity {
	wordCount := map[string]int{}
	for _, r := range reviews {
		words := strings.Fields(r)
		for i, w := range words {
			w = strings.TrimRight(w, ".,!?;:")
			if len([]rune(w)) < 3 {
				continue
			}
			if isCapitalized(w) {
				wordCount[w]++
			}
			if i+1 < len(words) {
				bigram := w + " " + words[i+1]
				bigram = strings.TrimRight(bigram, ".,!?;:")
				if isCapitalized(bigram) && len([]rune(bigram)) > 5 {
					wordCount[bigram]++
				}
			}
		}
	}

	return mapToEntities(wordCount)
}

func mapToEntities(m map[string]int) []Entity {
	var result []Entity
	for word, count := range m {
		if count >= 1 {
			eType := "Brand"
			if isPlace(word) {
				eType = "Place"
			} else if isMenu(word) {
				eType = "Menu"
			}
			result = append(result, Entity{Word: word, Count: count, Type: eType})
		}
	}
	return result
}

func isPlace(w string) bool {
	l := strings.ToLower(w)
	places := []string{"malang", "cafe", "kafe", "resto", "restoran", "warung", "kopi", "coffee",
		"batu", "surabaya", "area", "parkir", "toilet", "musala", "lantai"}
	for _, p := range places {
		if strings.Contains(l, p) {
			return true
		}
	}
	return false
}

func isMenu(w string) bool {
	l := strings.ToLower(w)
	menus := []string{"kopi", "coffee", "latte", "cappuccino", "espresso", "matcha", "teh",
		"makanan", "minuman", "snack", "roti", "kue", "nasi", "mie", "ayam"}
	for _, m := range menus {
		if strings.Contains(l, m) {
			return true
		}
	}
	return false
}

func extractKeywords(reviews []string) []Keyword {
	wordCount := map[string]int{}
	stopWords := map[string]bool{
		"dan": true, "di": true, "ke": true, "dari": true, "yang": true, "ini": true,
		"itu": true, "dengan": true, "untuk": true, "pada": true, "adalah": true,
		"ada": true, "bisa": true, "juga": true, "saya": true, "kami": true,
		"the": true, "and": true, "is": true, "in": true, "to": true, "of": true,
		"a": true, "an": true, "it": true, "its": true, "very": true, "so": true,
		"really": true, "sudah": true, "telah": true, "akan": true, "tidak": true,
	}

	for _, r := range reviews {
		words := strings.Fields(strings.ToLower(r))
		for _, w := range words {
			w = strings.TrimRight(w, ".,!?;:")
			if len([]rune(w)) < 4 || stopWords[w] {
				continue
			}
			wordCount[w]++
		}
	}

	var keywords []Keyword
	for w, c := range wordCount {
		if c >= 1 {
			kw := Keyword{Phrase: w, Count: c, Score: c * len(w)}
			if kw.Score > 0 {
				keywords = append(keywords, kw)
			}
		}
	}
	return keywords
}

func generateContentIdea(entities []Entity, keywords []Keyword) []string {
	var ideas []string
	for _, e := range entities {
		if e.Type == "Menu" && len(ideas) < 5 {
			ideas = append(ideas, fmt.Sprintf("Buat artikel tentang \"%s\" — varian, harga, dan rekomendasi", e.Word))
		}
	}
	for _, kw := range keywords {
		if kw.Count >= 2 && len(ideas) < 8 {
			ideas = append(ideas, fmt.Sprintf("Bahas \"%s\" dalam konteks review customer", kw.Phrase))
		}
	}
	if len(ideas) < 3 {
		ideas = append(ideas, "Highlight review positif di media sosial")
		ideas = append(ideas, "Buat FAQ dari pertanyaan umum customer")
		ideas = append(ideas, "Optimasi menu dengan schema.org/Menu")
	}
	return ideas
}

func buildReviewSchema(business string, reviews []string, s Sentiment) string {
	total := s.Positive + s.Negative + s.Neutral
	rating := 0.0
	if total > 0 {
		rating = float64(s.Positive*5) / float64(total)
	}
	return fmt.Sprintf(`{
  "@context": "https://schema.org",
  "@type": "LocalBusiness",
  "name": "%s",
  "aggregateRating": {
    "@type": "AggregateRating",
    "ratingValue": "%.1f",
    "reviewCount": "%d",
    "bestRating": "5"
  }
}`, business, rating, total)
}

func isCapitalized(s string) bool {
	if len(s) == 0 {
		return false
	}
	r := []rune(s)
	return unicode.IsUpper(r[0])
}

func getSampleReviews(business string) []string {
	return []string{
		"Tempatnya nyaman banget, cocok buat nongkrong lama-lama.",
		"Kopinya enak, barista ramah. Recommended!",
		"Suasananya cozy, cocok buat kerja juga.",
		"Harganya standar, tempatnya estetik banget.",
		"Pelayanan lambat, pesanan lama. Masih perlu perbaikan.",
		"Brewok Prime tempat favorit buat ngopi di Malang.",
	}
}
