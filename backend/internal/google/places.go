package gplaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type PlaceResult struct {
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	Rating        float64 `json:"rating"`
	UserRatings   int     `json:"user_ratings"`
	Phone         string  `json:"phone"`
	Website       string  `json:"website"`
	PlaceID       string  `json:"place_id"`
	Reviews       []Review `json:"reviews,omitempty"`
	Found         bool    `json:"found"`
	ErrorMessage  string  `json:"error_message,omitempty"`
}

type Review struct {
	Author  string  `json:"author_name"`
	Rating  int     `json:"rating"`
	Text    string  `json:"text"`
	Time    int64   `json:"time"`
}

type placesResponse struct {
	Candidates []struct {
		Name     string `json:"name"`
		PlaceID  string `json:"place_id"`
		FormattedAddress string `json:"formatted_address"`
		Rating   float64 `json:"rating"`
		UserRatingsTotal int `json:"user_ratings_total"`
		FormattedPhone   string `json:"formatted_phone_number"`
		Website  string `json:"website"`
		Reviews  []struct {
			AuthorName string `json:"author_name"`
			Rating     int    `json:"rating"`
			Text       string `json:"text"`
			Time       int64  `json:"time"`
		} `json:"reviews"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"candidates"`
	Status string `json:"status"`
}

func FindPlace(query, apiKey string) *PlaceResult {
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if apiKey == "" {
		return simulate(query)
	}

	return fetchFromAPI(query, apiKey)
}

func fetchFromAPI(query, apiKey string) *PlaceResult {
	endpoint := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=%s&inputtype=textquery&fields=name,place_id,formatted_address,rating,user_ratings_total,formatted_phone_number,website,reviews,geometry&key=%s",
		url.QueryEscape(query), apiKey,
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return &PlaceResult{Found: false, ErrorMessage: err.Error()}
	}
	defer resp.Body.Close()

	var pr placesResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return &PlaceResult{Found: false, ErrorMessage: err.Error()}
	}

	if pr.Status != "OK" || len(pr.Candidates) == 0 {
		return &PlaceResult{Found: false, ErrorMessage: fmt.Sprintf("Google API: %s", pr.Status)}
	}

	c := pr.Candidates[0]
	result := &PlaceResult{
		Name:        c.Name,
		Address:     c.FormattedAddress,
		Rating:      c.Rating,
		UserRatings: c.UserRatingsTotal,
		Phone:       c.FormattedPhone,
		Website:     c.Website,
		PlaceID:     c.PlaceID,
		Found:       true,
	}

	for _, r := range c.Reviews {
		result.Reviews = append(result.Reviews, Review{
			Author: r.AuthorName,
			Rating: r.Rating,
			Text:   r.Text,
			Time:   r.Time,
		})
	}

	return result
}

func simulate(query string) *PlaceResult {
	parts := strings.Fields(query)
	name := strings.Join(parts, " ")
	if len(parts) > 0 {
		name = strings.Title(parts[0])
		if len(parts) > 1 {
			name += " " + strings.Title(parts[1])
		}
	}

	city := "Malang"
	for _, p := range parts {
		lower := strings.ToLower(p)
		if lower == "malang" || lower == "surabaya" || lower == "batu" || lower == "jakarta" {
			city = strings.Title(lower)
		}
	}

	return &PlaceResult{
		Name:        name,
		Address:     fmt.Sprintf("Jl. Contoh No. 123, %s", city),
		Rating:      4.5,
		UserRatings: 87,
		Phone:       "+62 812-3456-7890",
		Website:     fmt.Sprintf("https://%s.example.com", strings.ToLower(strings.ReplaceAll(name, " ", ""))),
		PlaceID:     "simulated",
		Found:       true,
		Reviews:     generateSampleReviews(name),
	}
}

func generateSampleReviews(name string) []Review {
	return []Review{
		{Author: "Budi", Rating: 5, Text: fmt.Sprintf("%s tempatnya nyaman banget, kopinya enak!", name)},
		{Author: "Sari", Rating: 4, Text: "Suasananya cozy, cocok buat nongkrong kerja."},
		{Author: "Dimas", Rating: 5, Text: "Pelayanan ramah, recommended banget!"},
		{Author: "Rina", Rating: 3, Text: "Harganya standar, tempat lumayan."},
	}
}
