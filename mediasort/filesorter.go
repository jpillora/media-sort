package mediasort

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
	"github.com/fatih/color"
	"github.com/jpillora/media-sort/search"
)

type fileSorter struct {
	name, path string
	info       os.FileInfo
	err        error

	s               *Sorter
	ext             string
	mtype           string
	season, episode string
	episodeDate     string //weekly series
	year            string
}

func newFileSorter(s *Sorter, path string, info os.FileInfo) (*fileSorter, error) {
	//attempt to rule out file
	if !info.Mode().IsRegular() {
		return nil, nil
	}

	name := info.Name()
	ext := filepath.Ext(name)
	if _, exists := s.exts[ext]; !exists {
		return nil, nil
	}
	name = strings.TrimSuffix(name, ext)

	//setup
	return &fileSorter{
		s:       s,
		name:    name,
		path:    path,
		info:    info,
		ext:     ext,
		season:  "1",
		episode: "",
	}, nil
}

//always in a goroutine
func (f *fileSorter) goRun(wg *sync.WaitGroup) {
	f.err = f.run()
	wg.Done()
}

func (f *fileSorter) run() error {

	//normalize name
	query := normalize(f.name)

	//extract episode date (weekly show)
	if f.mtype == "" {
		m := epidate.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			f.mtype = "series"
			f.episodeDate = strings.Replace(m[2], " ", "-", 2)
		}
	}

	//extract episde season numbers
	if f.mtype == "" {
		m := episeason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			f.mtype = "series"
			f.season = m[3]
			f.episode = m[6]
		}
	}

	//extract *joined* episde season numbers
	if f.mtype == "" {
		m := joinedepiseason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			f.mtype = "series"
			f.season = m[2]
			f.episode = m[3]
		}
	}

	//extract release year
	m := year.FindStringSubmatch(query)
	if len(m) > 0 {
		query = m[1] //trim name
		if f.mtype == "" {
			f.mtype = "movie" //set type to "movie", if not already set
		}
		f.year = m[2]
	}

	//if the above fails, extract "Part 1/2/3..."
	if f.mtype == "" {
		m := partnum.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			f.mtype = "series"
			f.episode = m[2]
		}
	}

	//trim spaces
	query = strings.TrimSpace(query)

	//search for normalized name
	r, err := search.Do(query, f.year, f.mtype)
	if err != nil {
		return err //search failed
	}

	// normalize so the title is query-like
	title := normalize(r.Title)

	ldiff := abs(len(query) - len(title))

	dist, _ := levenshtein.ComputeDistance(query, title)

	// log.Printf("search found: '%s' -> '%s' %s => %d (#%d)", query, title, r.Type, dist, ldiff)

	// incorrect by more than the query!
	if dist-ldiff > 5 {
		return fmt.Errorf("Best match for '%s' was '%s'", f.name, r.Title)
	}

	//calculate destination path
	dest := ""
	if r.Type == "series" && f.episodeDate != "" {
		filename := fmt.Sprintf("%s %s%s", r.Title, f.episodeDate, f.ext)
		dest = filepath.Join(f.s.c.TVDir, r.Title, filename)
	} else if r.Type == "series" && f.episode != "" {
		s, _ := strconv.Atoi(f.season)
		e, _ := strconv.Atoi(f.episode)
		filename := fmt.Sprintf("%s S%02dE%02d%s", r.Title, s, e, f.ext)
		dest = filepath.Join(f.s.c.TVDir, r.Title, filename)
	} else if r.Type == "movie" {
		filename := fmt.Sprintf("%s (%s)%s", r.Title, r.Year, f.ext)
		dest = filepath.Join(f.s.c.MovieDir, filename)
	} else {
		return fmt.Errorf("TV Series with no episode found '%s'", query)
	}

	//DEBUG
	// log.Printf("SUCCESS = D%d #%d\n  %s\n  %s", r.Distance, len(query), query, r.Title)
	log.Printf("Moving\n  '%s'\n  └─> '%s'", f.path, color.GreenString(dest))

	if f.s.c.DryRun {
		return nil
	}

	//check already exists
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("File already exists '%s'", dest)
	}

	//mkdir -p
	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err //failed to mkdir
	}

	//mv
	err = os.Rename(f.path, dest)
	if err != nil {
		return err //failed to move
	}
	return nil
}

var yearstr = `(19\d\d|20[0,1]\d)`

var encodings = regexp.MustCompile(`\b(720p|1080p|hdtv|x264|dts|bluray)\b.*`) //strip all junk
var nonalpha = regexp.MustCompile(`[^A-Za-z0-9]`)
var spaces = regexp.MustCompile(`\s+`)
var episeason = regexp.MustCompile(`^(.+?)\bs?(eason)?(\d{1,2})(e|\ |\ e|x|xe)(pisode)?(\d{1,2})\b`)
var epidate = regexp.MustCompile(`^(.+?\b)(` + yearstr + ` \d{2} \d{2}|\d{2} \d{2} ` + yearstr + `)\b`)
var year = regexp.MustCompile(`^(.+?\b)` + yearstr + `\b`)
var joinedepiseason = regexp.MustCompile(`^(.+?\b)(\d)(\d{2})\b`)
var partnum = regexp.MustCompile(`^(.+?\b)(\d{1,2})\b`)

//TODO var romannumerals...

func normalize(s string) string {
	s = strings.ToLower(s)
	s = nonalpha.ReplaceAllString(s, " ")
	s = encodings.ReplaceAllString(s, "")
	s = spaces.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func abs(n int) int {
	if n < 0 {
		return n * -1
	}
	return n
}
