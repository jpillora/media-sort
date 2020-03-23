package mediasearch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func searchMovieDB(query, year string, mediatype MediaType) ([]Result, error) {
	yearKey := "year"
	path := "/search"
	if mediatype == Movie {
		path += "/movie"
	} else if mediatype == Series {
		yearKey = "first_air_date_year"
		path += "/tv"
	}
	v := url.Values{}
	v.Set("query", query)
	if year != "" {
		v.Set(yearKey, year)
	}
	if debugMode {
		log.Printf("Searching MovieDB API for '%s'", query)
	}
	resp, err := movieDBRequest(path, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	s := &movieDBSearch{}
	if err := json.NewDecoder(resp.Body).Decode(s); err != nil {
		return nil, fmt.Errorf("movieDB search: Failed to decode: %s", err)
	}
	results := make([]Result, len(s.Results))
	for i, mr := range s.Results {
		r, err := mr.toResult()
		if err != nil {
			return nil, err
		}
		results[i] = r
	}
	return results, nil
}

func movieDBRequest(path string, v url.Values) (*http.Response, error) {
	req, _ := http.NewRequest("GET", "https://api.themoviedb.org/3"+path+"?"+v.Encode()+string(vv), nil)
	return http.DefaultClient.Do(req)
}

type movieDBResult struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	BackdropPath     string   `json:"backdrop_path"`
	FirstAirDate     string   `json:"first_air_date"`
	GenreIDs         []int    `json:"genre_ids"`
	OriginalLanguage string   `json:"original_language"`
	OriginalName     string   `json:"original_name"`
	Overview         string   `json:"overview"`
	OriginCountry    []string `json:"origin_country"`
	PosterPath       string   `json:"poster_path"`
	Popularity       float64  `json:"popularity"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	Title            string   `json:"title"`
	Adult            bool     `json:"adult"`
	OriginalTitle    string   `json:"original_title"`
	ReleaseDate      string   `json:"release_date"`
	Video            bool     `json:"video"`
}

func (mr movieDBResult) toResult() (Result, error) {
	r := Result{}
	if mr.Title != "" && mr.ReleaseDate != "" {
		r.Type = Movie
		r.Title = mr.Title
		m := getYear.FindStringSubmatch(mr.ReleaseDate)
		if len(m) == 0 {
			return r, fmt.Errorf("movieDB error: No movie year: %s", mr.ReleaseDate)
		}
		r.Year = m[1]
		return r, nil
	} else if mr.Name != "" && mr.FirstAirDate != "" {
		r.Type = Series
		r.Title = mr.Name
		m := getYear.FindStringSubmatch(mr.FirstAirDate)
		if len(m) == 0 {
			return Result{}, fmt.Errorf("movieDB error: No series year: %s", mr.FirstAirDate)
		}
		r.Year = m[1]
		return r, nil
	}
	return r, fmt.Errorf("movieDB error: Unknown result: %+v", mr)
}

type movieDBData struct {
	MovieResults  []movieDBResult `json:"movie_results"`
	TVResults     []movieDBResult `json:"tv_results"`
	StatusCode    int             `json:"status_code"`
	StatusMessage string          `json:"status_message"`
}

type movieDBSearch struct {
	Page          int             `json:"page"`
	Results       []movieDBResult `json:"results"`
	TotalResults  int             `json:"total_results"`
	TotalPages    int             `json:"total_pages"`
	StatusCode    int             `json:"status_code"`
	StatusMessage string          `json:"status_message"`
}

var vv = []byte{0x26, 0x61, 0x70, 0x69, 0x5f, 0x6b, 0x65, 0x79, 0x3d, 0x63, 0x37, 0x65, 0x39, 0x38, 0x36, 0x30, 0x39, 0x64, 0x35, 0x38, 0x61, 0x30, 0x30, 0x66, 0x65, 0x64, 0x62, 0x39, 0x63, 0x63, 0x35, 0x62, 0x61, 0x64, 0x33, 0x30, 0x33, 0x33, 0x62, 0x36, 0x30}
