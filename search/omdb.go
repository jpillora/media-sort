package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
)

type Search struct {
	Search OMDBResults
}

type OMDBResult struct {
	Title  string
	Year   string
	Type   string
	ImdbID string
	//meta
	Distance int
}

type OMDBResults []*OMDBResult

func (rs OMDBResults) Len() int      { return len(rs) }
func (rs OMDBResults) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }
func (rs OMDBResults) Less(i, j int) bool {
	//sort by string dist
	if rs[i].Distance != rs[j].Distance {
		return rs[i].Distance < rs[j].Distance
	}
	//sort by newest
	return rs[i].Year > rs[j].Year
}

func omdbRequest(v url.Values) (*http.Response, error) {
	req, _ := http.NewRequest("GET", "http://www.omdbapi.com/?"+v.Encode(), nil)
	return http.DefaultClient.Do(req)
}

func omdbSearch(query string, mediatype string) (imdbID, error) {

	v := url.Values{}
	v.Set("s", query)
	if mediatype != "" {
		v.Set("type", mediatype)
	}

	resp, err := omdbRequest(v)
	if err != nil {
		return missingID, err
	}
	defer resp.Body.Close()

	s := &Search{}
	if err := json.NewDecoder(resp.Body).Decode(s); err != nil {
		return missingID, fmt.Errorf("OMDB Search: Failed to decode: %s", err)
	}

	results := s.Search
	if len(results) == 0 {
		return missingID, errors.New("No results")
	}

	for _, r := range results {
		r.Distance, _ = levenshtein.ComputeDistance(query, strings.ToLower(r.Title))
	}

	sort.Sort(results)

	// inspect...
	// for i, r := range results {
	// 	fmt.Printf("%d %s\n", i, r.Title)
	// }

	return imdbID(results[0].ImdbID), nil
}

var releaseYear = regexp.MustCompile(`^(\d{4})`)

func omdbGet(id imdbID) (*OMDBResult, error) {

	v := url.Values{}
	v.Set("i", string(id))

	resp, err := omdbRequest(v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &OMDBResult{}
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, fmt.Errorf("OMDB Get: Failed to decode: %s", err)
	}

	m := releaseYear.FindStringSubmatch(r.Year)
	if len(m) == 0 {
		return nil, fmt.Errorf("OMDB Get: Invalid year")
	}
	r.Year = m[1]

	return r, nil
}
