package mediasearch

import (
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
	v.Set("q", query)
	v.Set("btnI", "I'm feeling lucky")
	urlstr := "https://www.google.com/search?" + v.Encode()
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	//I'm a browser... :)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36")
	//roundtripper doesn't follow redirects
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	//assume redirection
	if resp.StatusCode/100 != 3 {
		return nil, fmt.Errorf("Google search expected redirect, got %d", resp.StatusCode)
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
