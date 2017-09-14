package mediasearch

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func searchTVMaze(query, year string, mediatype MediaType) ([]Result, error) {
	v := url.Values{}
	v.Set("q", query)
	if debugMode {
		log.Printf("Searching TVMaze for '%s'", query)
	}
	urlstr := "http://api.tvmaze.com/search/shows?" + v.Encode()
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tvMazeResults := []*tvMazeResult{}
	if err := json.NewDecoder(resp.Body).Decode(&tvMazeResults); err != nil {
		return nil, err
	}
	rs := []Result{}
	for _, tvMazeResult := range tvMazeResults {
		m := getYear.FindStringSubmatch(tvMazeResult.Show.Premiered)
		if len(m) == 0 {
			continue //skip no year
		}
		rs = append(rs, Result{
			Title:    tvMazeResult.Show.Name,
			Year:     m[1],
			Type:     Series,
			Accuracy: accuracy(query, tvMazeResult.Show.Name),
		})
	}
	return rs, nil
}

type tvMazeResult struct {
	Score float64 `json:"score"`
	Show  struct {
		Links struct {
			Previousepisode struct {
				Href string `json:"href"`
			} `json:"previousepisode"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
		Externals struct {
			Thetvdb int `json:"thetvdb"`
			Tvrage  int `json:"tvrage"`
		} `json:"externals"`
		Genres []string `json:"genres"`
		ID     int      `json:"id"`
		Image  struct {
			Medium   string `json:"medium"`
			Original string `json:"original"`
		} `json:"image"`
		Language string `json:"language"`
		Name     string `json:"name"`
		Network  struct {
			Country struct {
				Code     string `json:"code"`
				Name     string `json:"name"`
				Timezone string `json:"timezone"`
			} `json:"country"`
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"network"`
		Premiered string `json:"premiered"`
		Rating    struct {
			Average float64 `json:"average"`
		} `json:"rating"`
		Runtime  int `json:"runtime"`
		Schedule struct {
			Days []interface{} `json:"days"`
			Time string        `json:"time"`
		} `json:"schedule"`
		Status     string      `json:"status"`
		Summary    string      `json:"summary"`
		Type       string      `json:"type"`
		Updated    int         `json:"updated"`
		URL        string      `json:"url"`
		WebChannel interface{} `json:"webChannel"`
		Weight     int         `json:"weight"`
	} `json:"show"`
}
