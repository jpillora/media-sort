package mediasearch

import (
	"errors"
	"fmt"
	"log"

	"net/http"
	"net/url"
	"regexp"
)

var imdbIDRe = regexp.MustCompile(`\/(tt\d+)\/`)

//uses im feeling lucky and grabs the "Location"
//header from the 302, which contains the IMDB ID
func searchGoogle(query, year string, mediatype MediaType) ([]Result, error) {
	if year != "" {
		query += " " + year
	}
	if string(mediatype) != "" {
		query += " " + string(mediatype)
	}
	query += " site:imdb.com"
	if debugMode {
		log.Printf("Searching Google for '%s'", query)
	}
	v := url.Values{}
	v.Set("btnI", "") //I'm feeling lucky
	v.Set("q", query)
	urlstr := "https://www.google.com.au/search?" + v.Encode()
	req, err := http.NewRequest("HEAD", urlstr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	//I'm a browser... :)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36")
	//roundtripper doesn't follow redirects
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//assume redirection
	if resp.StatusCode != 302 {
		return nil, errors.New("Google search failed")
	}
	//extract Location header URL
	loc := resp.Header.Get("Location")
	//extract imdb ID
	m := imdbIDRe.FindStringSubmatch(loc)
	if len(m) == 0 {
		return nil, fmt.Errorf("No IMDB match (%s)", loc)
	}
	//lookup imdb ID using OMDB
	r, err := imdbGet(imdbID(m[1]), mediatype)
	if err != nil {
		return nil, err
	}
	return []Result{r}, nil
}
