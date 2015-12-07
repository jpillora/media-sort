package mediasearch

import (
	"errors"

	"net/http"
	"net/url"
	"regexp"
)

var endpoint = ""

// var endpoint = "https://echo.jpillora.com/search"

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

	v := url.Values{}
	v.Set("btnI", "") //I'm feeling lucky
	v.Set("q", query+" site:imdb.com")
	urlstr := "https://www.google.com.au/search?" + v.Encode()
	req, err := http.NewRequest("HEAD", urlstr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	//I'm a browser... :)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 302 {
		return nil, errors.New("Google search failed")
	}
	url, _ := url.Parse(resp.Header.Get("Location"))
	if url.Host != "www.imdb.com" {
		return nil, errors.New("Google IMDB redirection failed")
	}
	m := imdbIDRe.FindStringSubmatch(url.Path)
	if len(m) == 0 {
		return nil, errors.New("No IMDB match")
	}
	r, err := imdbGet(imdbID(m[1]))
	if err != nil {
		return nil, err
	}
	return []Result{r}, nil
}
