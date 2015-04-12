package search

import (
	"fmt"
	"regexp"
	"sync"
)

type imdbID string

const missingID = imdbID("")

//lock protects the cache/inflight maps
var lock sync.Mutex
var cache map[string]*OMDBResult
var inflight map[string]*sync.WaitGroup

func init() {
	cache = map[string]*OMDBResult{}
	inflight = map[string]*sync.WaitGroup{}
}

var nonalpha = regexp.MustCompile(`[^a-z0-9]`)

//Fuzzy search for IMDB data
func Do(query, year, mediatype string) (*OMDBResult, error) {

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

	// log.Printf("searching for '%s' (%s)", query, mediatype)

	//	since google has strict throttling,
	//  we first try omdb search
	id, err := omdbSearch(query, year, mediatype)
	if err != nil {
		// then fallback to google
		id, err = googleSearch(query, year, mediatype)
		if err != nil {
			return nil, fmt.Errorf("OMDB and Google searches failed")
		}
	}

	// pull information for that IMDB ID
	r, err = omdbGet(id)
	if err != nil {
		return nil, err
	}

	lock.Lock()
	cache[query] = r
	lock.Unlock()
	w.Done()

	return r, nil
}
