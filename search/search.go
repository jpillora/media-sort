package search

import (
	"fmt"
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
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

//Fuzzy search for IMDB data
func Do(query string, mediatype string) (*OMDBResult, error) {

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
	id, err := omdbSearch(query, mediatype)
	if err != nil {
		// then fallback to google
		id, err = googleSearch(query, mediatype)
		if err != nil {
			return nil, fmt.Errorf("OMDB and Google searches failed")
		}
	}

	// pull information for that IMDB ID
	r, err = omdbGet(id)
	if err != nil {
		return nil, err
	}

	// santiy check
	r.Distance, _ = levenshtein.ComputeDistance(query, strings.ToLower(r.Title))
	// incorrect by more than the query!
	if len(query)-r.Distance < 0 {
		return nil, fmt.Errorf("Best match was '%s' (%s)", r.Title, mediatype)
	}

	lock.Lock()
	cache[query] = r
	lock.Unlock()
	w.Done()

	return r, nil
}
