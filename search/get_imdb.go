package mediasearch

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type imdbID string

//imdbGet actually uses omdb because it accepts IMDB IDs
func imdbGet(id imdbID) (Result, error) {

	v := url.Values{}
	v.Set("i", string(id))

	resp, err := omdbRequest(v)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	r := &omdbResult{}
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return Result{}, fmt.Errorf("omdb Get: Failed to decode: %s", err)
	}
	if r.Error != "" {
		return Result{}, fmt.Errorf("omdb Error: %s: %s", id, r.Error)
	}
	//dont allow episode respose, return series instead
	if r.Type == "episode" {
		return imdbGet(imdbID(r.SeriesID))
	}
	m := getYear.FindStringSubmatch(r.Year)
	if len(m) == 0 {
		return Result{}, fmt.Errorf("omdb Get: No year: %+v", r)
	}
	r.Result.Year = m[1]
	return r.Result, nil
}
