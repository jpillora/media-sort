package mediasearch

import (
	"fmt"
	"log"
	"sort"
)

// DefaultThreshold for matching names. 100 is a perfect match.
// This represents the string distance from the query (or filename)
// to the found movie/series title..
const DefaultThreshold = 95

//matcher collects search results and finds the best match
type matcher struct {
	query, year string
	threshold   int
	resultSlice []*Result
	resultMap   map[string]*Result
}

func (m *matcher) add(r Result) {
	if m.resultMap == nil {
		m.resultMap = map[string]*Result{}
	} else if rother, ok := m.resultMap[r.Title]; ok {
		r.IsDupe = true
		rother.IsDupe = true
	}
	r.Accuracy = accuracy(m.query, r.Title)
	m.resultMap[r.Title] = &r
	m.resultSlice = append(m.resultSlice, &r)
}

func (m *matcher) bestMatch() (Result, error) {
	if len(m.resultSlice) == 0 {
		return Result{}, fmt.Errorf("No results")
	}
	sort.Sort(m)
	if debugMode {
		log.Println("Matched results:")
		for i, r := range m.resultSlice {
			log.Printf("#%d: %s (acc: %d)", i, r, r.Accuracy)
		}
	}
	r := m.resultSlice[0]
	if r.Accuracy < m.threshold {
		return Result{}, fmt.Errorf("No results (closest result was '%s' with an accuracy score of %d)", r.Title, r.Accuracy)
	}
	return *r, nil
}

func (m *matcher) Len() int { return len(m.resultSlice) }
func (m *matcher) Swap(i, j int) {
	m.resultSlice[i], m.resultSlice[j] = m.resultSlice[j], m.resultSlice[i]
}
func (m *matcher) Less(i, j int) bool {
	//sort by string dist
	if m.resultSlice[i].Accuracy != m.resultSlice[j].Accuracy {
		return m.resultSlice[i].Accuracy > m.resultSlice[j].Accuracy
	}
	//specific year
	if m.year != "" {
		return m.resultSlice[i].Year == m.year
	}
	//sort by newest
	return m.resultSlice[i].Year > m.resultSlice[j].Year
}
