package mediasearch

import (
	"fmt"
	"log"
	"sync"

	"github.com/fatih/color"
)

const debugMode = false

//search function interface
type search func(string, string, MediaType) ([]Result, error)

//various searches based on media-type
var tvSearches = []search{searchTVMaze, searchMovieDB, searchGoogle}
var movieSearches = []search{searchMovieDB, searchGoogle}

//thread-safe global search cache
//lock protects the cache/inflight maps
var lock sync.Mutex
var cache = map[string]Result{} //TODO(jpillora) global cache will grow forever, should convert to LRU cache
var inflight = map[string]*sync.WaitGroup{}

//Search for IMDB data (query is required, year and media type are optional)
func Search(query, year, mediatype string) (Result, error) {
	return SearchThreshold(query, year, mediatype, DefaultThreshold)
}

//SearchThreshold for IMDB data with a specific match threshoold
func SearchThreshold(query, year, mediatype string, threshold int) (Result, error) {
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
	msg := fmt.Sprintf("Searching %s", color.CyanString(query))
	if m := string(mediatype); m != "" {
		msg += " (" + color.CyanString(m) + ")"
	}
	if year != "" {
		msg += " from " + color.CyanString(year)
	}
	log.Print(msg)
	//search various search engines
	var searcheEngines = movieSearches
	if MediaType(mediatype) == Series {
		searcheEngines = tvSearches
	}
	//search returns results
	var results []Result
	var err error
	for _, s := range searcheEngines {
		results, err = s(query, year, mt)
		if len(results) > 0 {
			break
		}
	}
	if len(results) == 0 && err != nil {
		return Result{}, fmt.Errorf("No results (%s)", err)
	}
	if len(results) == 0 {
		return Result{}, fmt.Errorf("No results")
	}
	//matcher picks result (r)
	m := matcher{query: query, year: year, threshold: threshold}
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
