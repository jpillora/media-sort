package mediasearch

import (
	"fmt"
	"log"
	"sync"

	"github.com/fatih/color"
)

var Debug = false
var Info = true

//search function interface
type search func(string, string, MediaType) ([]Result, error)

//various searches based on media-type
var defaultSearches = []search{searchOMDB, searchGoogle}
var tvSearches = append([]search{searchTVMaze}, defaultSearches...)
var movieSearches = defaultSearches /*TODO moviedb,*/

//thread-safe global search cache
//lock protects the cache/inflight maps
var lock sync.Mutex
var cache = map[string]Result{} //TODO(jpillora) global cache will grow forever, should convert to LRU cache
var inflight = map[string]*sync.WaitGroup{}

//Fuzzy search for IMDB data (query is required, year and media type are optional)
func Search(query, year, mediatype string) (Result, error) {
	if year != "" && !onlyYear.MatchString(year) {
		return Result{}, fmt.Errorf("Invalid year (%s)", year)
	}
	mt := MediaType(mediatype)
	if mediatype != "" && mt != Movie && mt != Series {
		return Result{}, fmt.Errorf("Invalid media type (%s)", mediatype)
	}

	lock.Lock()
	//cached searches are served instantly
	r, exists := cache[query]
	if exists {
		lock.Unlock()
		return r, nil
	}

	//duplicate searchs wait on the first
	w, inf := inflight[query]
	if inf {
		lock.Unlock()
		w.Wait()
		return cache[query], nil
	}

	w = &sync.WaitGroup{}
	w.Add(1)
	inflight[query] = w
	lock.Unlock()
	//queue up removal of inflight search since cached should be set
	defer func() {
		lock.Lock()
		w.Done()
		delete(inflight, query)
		lock.Unlock()
	}()
	//show searches
	if Info {
		msg := fmt.Sprintf("Searching %s", color.CyanString(query))
		if m := string(mediatype); m != "" {
			msg += " (" + m + ")"
		}
		if year != "" {
			msg += " from " + year
		}
		log.Print(msg)
	}
	//search various search engines
	var searches = defaultSearches
	if MediaType(mediatype) == Series {
		searches = tvSearches
	}

	//search returns results
	var results []Result
	var err error
	for _, s := range searches {
		results, err = s(query, year, mt)
		if len(results) > 0 {
			break
		}
	}
	if len(results) == 0 {
		return Result{}, fmt.Errorf("No results (%s)", err)
	}

	//matcher picks result (r)
	m := matcher{query: query, year: year}
	otherTypes := []Result{}
	for _, result := range results {
		//only consider tv/movies
		if string(result.Type) != "" && result.Type != Series && result.Type != Movie {
			continue
		}
		//if media type set, ensure match
		if mediatype != "" && result.Type != mt {
			otherTypes = append(otherTypes, result)
		} else {
			m.add(result)
		}
	}
	if len(m.resultSlice) == 0 {
		//if nothing was added, use mismatched types
		for _, result := range otherTypes {
			m.add(result)
		}
	}
	if r, err = m.bestMatch(); err != nil {
		return Result{}, err
	}

	lock.Lock()
	cache[query] = r
	lock.Unlock()

	return r, nil
}
