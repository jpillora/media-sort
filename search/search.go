package mediasearch

import (
	"fmt"
	"log"
	"sync"
)

type search func(string, string, MediaType) ([]Result, error)

//lock protects the cache/inflight maps
//TODO(jpillora) global cache will grow forever, should convert to LRU cache
var lock sync.Mutex
var cache = map[string]Result{}
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

	log.Printf("searching for '%s' (%s) from %s", query, string(mediatype), year)

	//search various search engines
	var searches []search
	if MediaType(mediatype) == Series {
		searches = []search{searchTVMaze, searchOMDB, searchGoogle}
	} else if MediaType(mediatype) == Movie {
		searches = []search{ /*TODO moviedb,*/ searchOMDB, searchGoogle}
	} else {
		searches = []search{searchOMDB, searchGoogle}
	}

	//search returns results
	var results []Result
	var err error
	for _, s := range searches {
		if results, err = s(query, year, mt); err != nil {
			log.Printf("search error: %s", err)
		} else if len(results) > 0 {
			break
		}
	}
	if len(results) == 0 {
		return Result{}, fmt.Errorf("No results")
	}

	//matcher picks result (r)
	m := matcher{query: query, year: year}
	for _, result := range results {
		//only consider tv/movies
		if string(result.Type) != "" && result.Type != Series && result.Type != Movie {
			continue
		}
		//if media type set, ensure match
		if mediatype != "" && result.Type != mt {
			continue
		}
		m.add(result)
	}
	if r, err = m.bestMatch(); err != nil {
		return Result{}, err
	}

	lock.Lock()
	cache[query] = r
	lock.Unlock()
	w.Done()

	return r, nil
}
