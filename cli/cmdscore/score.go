package cmdscore

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

type ScoreResult struct {
	Overall    ScoreCategory  `json:"overall"`
	Breakdown  Breakdown      `json:"breakdown"`
	Suggestion []string       `json:"suggestions"`
}

type ScoreCategory struct {
	Label string `json:"label"`
	Score int    `json:"score"`
	Max   int    `json:"max"`
}

type Breakdown struct {
	Structure     ScoreCategory `json:"structure"`
	QACoverage    ScoreCategory `json:"qa_coverage"`
	EntityRichness ScoreCategory `json:"entity_richness"`
	CitationQuality ScoreCategory `json:"citation_quality"`
	SchemaReadiness ScoreCategory `json:"schema_readiness"`
	Readability   ScoreCategory `json:"readability"`
}

func Analyze(content string) *ScoreResult {
	structure := scoreStructure(content)
	qa := scoreQACoverage(content)
	entity := scoreEntityRichness(content)
	citation := scoreCitationQuality(content)
	schema := scoreSchemaReadiness(content)
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
			suggestions = append(suggestions, "Include direct question-answer pairs for common queries")
		}
		if entity.Score < entity.Max/2 {
			suggestions = append(suggestions, "Increase use of named entities and key terminology")
		}
		if citation.Score < citation.Max/2 {
			suggestions = append(suggestions, "Add citations, references, or authoritative sources")
		}
		if schema.Score < schema.Max/2 {
			suggestions = append(suggestions, "Use structured content like lists, tables, and definitions")
		}
	}

	return &ScoreResult{
		Overall: ScoreCategory{Label: label, Score: total, Max: maxTotal},
		Breakdown: Breakdown{
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

	lines := strings.Split(content, "\n")
	headingCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "==") {
			headingCount++
		}
	}
	if headingCount == 0 {
		score += 0
	} else if headingCount < 3 {
		score += 10
	} else if headingCount < 6 {
		score += 18
	} else {
		score += 25
	}

	wordCount := len(strings.Fields(content))
	paraCount := countParagraphs(content)
	if paraCount > 0 {
		avgPara := wordCount / paraCount
		if avgPara >= 30 && avgPara <= 100 {
			score += 7
		} else if avgPara > 100 {
			score += 3
		}
	}

	if score > max {
		score = max
	}

	label := scoreLabel(score, max)
	return ScoreCategory{Label: label, Score: score, Max: max}
}

func scoreQACoverage(content string) ScoreCategory {
	max := 20
	score := 0

	lower := strings.ToLower(content)
	qaIndicators := []string{"what is", "how does", "why is", "what are", "how to", "what does", "can you", "what's"}
	for _, indicator := range qaIndicators {
		if strings.Contains(lower, indicator) {
			score += 3
		}
	}

	if score > max {
		score = max
	}

	label := scoreLabel(score, max)
	return ScoreCategory{Label: label, Score: score, Max: max}
}

func scoreEntityRichness(content string) ScoreCategory {
	max := 15
	score := 0

	words := strings.Fields(content)
	if len(words) == 0 {
		return ScoreCategory{Label: "None", Score: 0, Max: max}
	}

	entityCount := 0
	for _, w := range words {
		w = strings.Trim(w, " ,;:()[]{}\"'.")
		if len(w) > 3 && isUpperCase(rune(w[0])) {
			entityCount++
		}
	}

	ratio := float64(entityCount) / float64(len(words)) * 100
	switch {
	case ratio >= 15:
		score = max
	case ratio >= 10:
		score = 10
	case ratio >= 5:
		score = 6
	default:
		score = 3
	}

	label := scoreLabel(score, max)
	return ScoreCategory{Label: label, Score: score, Max: max}
}

func scoreCitationQuality(content string) ScoreCategory {
	max := 15
	score := 0

	lower := strings.ToLower(content)
	indicators := []string{"according to", "source", "reference", "citation", "study", "research", "report", "data shows", "statistics", "doi:"}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			score += 2
		}
	}

	if strings.Contains(content, "http") || strings.Contains(content, "www.") {
		score += 5
	}

	if score > max {
		score = max
	}

	label := scoreLabel(score, max)
	return ScoreCategory{Label: label, Score: score, Max: max}
}

func scoreSchemaReadiness(content string) ScoreCategory {
	max := 15
	score := 0

	schemaMarkers := []string{"<table", "<ul", "<ol", "<dl", "|", "- "}
	for _, m := range schemaMarkers {
		if strings.Contains(content, m) {
			score += 3
		}
	}

	listCount := countBulletLists(content)
	if listCount > 3 {
		score += 3
	} else if listCount > 0 {
		score += 2
	}

	if score > max {
		score = max
	}

	label := scoreLabel(score, max)
	return ScoreCategory{Label: label, Score: score, Max: max}
}

func scoreReadability(content string) ScoreCategory {
	max := 10
	score := 0

	words := strings.Fields(content)
	if len(words) == 0 {
		return ScoreCategory{Label: "None", Score: 0, Max: max}
	}

	sentences := strings.FieldsFunc(content, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})

	var avgWordsPerSentence float64
	if len(sentences) > 0 {
		avgWordsPerSentence = float64(len(words)) / float64(len(sentences))
	}

	switch {
	case avgWordsPerSentence <= 15:
		score = 10
	case avgWordsPerSentence <= 20:
		score = 7
	case avgWordsPerSentence <= 30:
		score = 5
	default:
		score = 2
	}

	longWords := 0
	for _, w := range words {
		if utf8.RuneCountInString(w) > 8 {
			longWords++
		}
	}
	complexRatio := float64(longWords) / float64(len(words))
	if complexRatio > 0.3 {
		score = int(math.Max(0, float64(score)-3))
	}

	label := scoreLabel(score, max)
	return ScoreCategory{Label: label, Score: score, Max: max}
}

func countParagraphs(content string) int {
	paras := strings.Split(content, "\n\n")
	count := 0
	for _, p := range paras {
		if len(strings.TrimSpace(p)) > 0 {
			count++
		}
	}
	return count
}

func countBulletLists(content string) int {
	count := 0
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			count++
		}
	}
	return count
}

func isUpperCase(r rune) bool {
	return unicode.IsUpper(r)
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
