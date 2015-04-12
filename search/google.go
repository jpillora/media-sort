package search

import (
	"errors"

	"net/http"
	"net/url"
	"regexp"
)

var endpoint = "https://www.google.com.au/search"

// var endpoint = "https://echo.jpillora.com/search"

var imdbIDRe = regexp.MustCompile(`\/(tt\d+)\/`)

//uses im feeling lucky and grabs the "Location"
//header from the 302, which contains the IMDB ID
func googleSearch(query, year, mediatype string) (imdbID, error) {

	if year != "" {
		query += " " + year
	}
	if mediatype != "" {
		query += " " + mediatype
	}

	v := url.Values{}
	v.Set("btnI", "") //I'm feeling lucky
	v.Set("q", query+" site:imdb.com")
	// q.Set("userip", "108.170.219.33")

	urlstr := endpoint + "?" + v.Encode()
	req, err := http.NewRequest("HEAD", urlstr, nil)
	if err != nil {
		return missingID, err
	}
	req.Header.Set("Accept", "*/*")
	//I'm a browser... :)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36")

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return missingID, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 302 {
		// b, _ := ioutil.ReadAll(resp.Body)
		// log.Printf("\nHEAD %s\nRESPONSE: %d\n%+v\n%s", urlstr, resp.StatusCode, resp.Header, b)
		//didnt work
		return missingID, errors.New("Google search failed")
	}

	url, _ := url.Parse(resp.Header.Get("Location"))
	if url.Host != "www.imdb.com" {
		return missingID, errors.New("Google IMDB redirection failed")
	}

	m := imdbIDRe.FindStringSubmatch(url.Path)
	if len(m) == 0 {
		return missingID, errors.New("Invalid ID")
	}

	return imdbID(m[1]), nil
}
