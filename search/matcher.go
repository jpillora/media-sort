package mediasearch

import (
	"fmt"
	"log"
	"sort"
)

//matcher collects and finds the best match
type matcher struct {
	query, year string
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
	r.strdist = dist(m.query, r.Title)
	m.resultMap[r.Title] = &r
	m.resultSlice = append(m.resultSlice, &r)
}

func (m *matcher) bestMatch() (Result, error) {
	if len(m.resultSlice) == 0 {
		return Result{}, fmt.Errorf("No results")
	}
	sort.Sort(m)
	log.Println("best result from:")
	for i, r := range m.resultSlice {
		log.Printf("[%d] %+v", i, r)
	}
	r := m.resultSlice[0]
	if r.strdist > 5 {
		return Result{}, fmt.Errorf("No results (closest result was '%s')", r.Title)
	}
	return *r, nil
}

func (m *matcher) Len() int { return len(m.resultSlice) }
func (m *matcher) Swap(i, j int) {
	m.resultSlice[i], m.resultSlice[j] = m.resultSlice[j], m.resultSlice[i]
}
func (m *matcher) Less(i, j int) bool {
	//sort by string dist
	if m.resultSlice[i].strdist != m.resultSlice[j].strdist {
		return m.resultSlice[i].strdist < m.resultSlice[j].strdist
	}
	//specific year
	if m.year != "" {
		return m.resultSlice[i].Year == m.year
	}
	//sort by newest
	return m.resultSlice[i].Year > m.resultSlice[j].Year
}
