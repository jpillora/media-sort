package mediasearch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type googleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title         string   `json:"title"`
			Authors       []string `json:"authors"`
			PublishedDate string   `json:"publishedDate"`
			Description   string   `json:"description"`
			Categories    []string `json:"categories"`
			ReadingModes  struct {
				Audio bool `json:"audio"`
			} `json:"readingModes"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			SeriesInfo struct {
				Title      string `json:"title"`
				BookNumber int    `json:"bookNumber"`
			} `json:"seriesInfo"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

func searchGoogleBooks(query, year string, mediatype MediaType) ([]Result, error) {
	if mediatype != Audiobook {
		return nil, fmt.Errorf("GoogleBooks search only supports audiobooks")
	}

	// Construct search URL
	searchURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s",
		url.QueryEscape(query+" audiobook"))

	// Add year if provided
	if year != "" {
		searchURL += "&publishedDate=" + url.QueryEscape(year)
	}

	// Make request
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var data googleBooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Convert to Results
	var results []Result
	for _, item := range data.Items {
		info := item.VolumeInfo

		// Skip if not an audiobook
		if !info.ReadingModes.Audio {
			continue
		}

		// Create result
		r := Result{
			Title:    info.Title,
			Type:     Audiobook,
			Year:     info.PublishedDate[:4], // Extract year from date
			Accuracy: accuracy(query, info.Title),
		}

		// Add authors
		if len(info.Authors) > 0 {
			r.Author = info.Authors[0]
			if len(info.Authors) > 1 {
				r.Author = "Various Authors"
			}
		}

		// Add series info if available
		if info.SeriesInfo.Title != "" {
			r.Series = info.SeriesInfo.Title
			r.SeriesPart = info.SeriesInfo.BookNumber
		}

		results = append(results, r)
	}

	return results, nil
}

func init() {
	// Register the GoogleBooks search
	audiobookSearches = append(audiobookSearches, searchGoogleBooks)
}
