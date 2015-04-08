package search

import (
	"encoding/json"
	"net/http"
	"sync"
)

var lock sync.Mutex
var cache map[string]Results
var inflight map[string]sync.WaitGroup

func init() {
	cache = map[string]Results{}
	inflight = map[string]sync.WaitGroup{}
}

// Results is an ordered list of search results.
type Results []Result

// A Result contains the title and URL of a search result.
type Result struct {
	Title, URL string
}

// Do sends query to Google search and returns the results.
func Do(query string) (Results, error) {

	lock.Lock()

	r, exists := cache[query]
	if exists {
		lock.Unlock()
		return r, nil
	}

	w, inf := inflight[query]
	if inf {
		lock.Unlock()
		w.Wait()
		return cache[query], nil
	}

	w = sync.WaitGroup{}
	w.Add(1)
	inflight[query] = w
	lock.Unlock()

	// Prepare the Google Search API request.
	req, err := http.NewRequest("GET", "https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)
	// q.Set("userip", "108.170.219.33")
	req.URL.RawQuery = q.Encode()

	var results Results
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the JSON search result.
	// https://developers.google.com/web-search/docs/#fonje
	var data struct {
		ResponseData struct {
			Results []struct {
				TitleNoFormatting string
				URL               string
			}
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	for _, res := range data.ResponseData.Results {
		results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
	}

	lock.Lock()
	cache[query] = results
	lock.Unlock()

	return results, nil
}
