package mediasearch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type omdbSearch struct {
	Search []Result
}

type omdbResult struct {
	Result
	ImdbID   string
	SeriesID string
	Error    string
}

func omdbRequest(v url.Values) (*http.Response, error) {
	req, _ := http.NewRequest("GET", "http://www.omdbapi.com/?"+v.Encode(), nil)
	return http.DefaultClient.Do(req)
}

func searchOMDB(query, year string, mediatype MediaType) ([]Result, error) {
	v := url.Values{}
	v.Set("s", query)
	// we want to include other matches so we can mark dupes
	// if year != "" { v.Set("y", year) }
	// if string(mediatype) != "" {
	// 	v.Set("type", string(mediatype))
	// }
	if debugMode {
		log.Printf("Searching OMDB API for '%s'", query)
	}
	resp, err := omdbRequest(v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	s := &omdbSearch{}
	if err := json.NewDecoder(resp.Body).Decode(s); err != nil {
		return nil, fmt.Errorf("omdb Search: Failed to decode: %s", err)
	}
	return s.Search, nil
}
