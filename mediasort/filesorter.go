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

	"github.com/jpillora/media-sort/search"
)

type fileSorter struct {
	name, path string
	info       os.FileInfo
	err        error

	s               *Sorter
	query           string
	ext             string
	mtype           string
	season, episode string
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
		query:   name,
		ext:     ext,
		season:  "1",
		episode: "1",
	}, nil
}

var junk = regexp.MustCompile(`\b(720p|1080p|hdtv|x264|dts|bluray)\b`)
var nonalpha = regexp.MustCompile(`[^A-Za-z0-9]`)
var spaces = regexp.MustCompile(`\s+`)
var episeason = regexp.MustCompile(`^(.+?)\bs?(eason)?(\d{1,2})(e|\ |\ e|x|xe)(pisode)?(\d{1,2})\b`)
var year = regexp.MustCompile(`^(.+?[^\d])(\d{4})[^\d]`)
var partnum = regexp.MustCompile(`^(.+?\b)(\d{1,2})\b`)

//TODO var romannumerals...

//always in a goroutine
func (f *fileSorter) goRun(wg *sync.WaitGroup) {
	f.err = f.run()
	wg.Done()
}

func (f *fileSorter) run() error {

	//normalize name
	f.query = strings.ToLower(f.query)
	f.query = strings.Replace(f.query, ".", " ", -1)
	f.query = nonalpha.ReplaceAllString(f.query, " ")
	f.query = junk.ReplaceAllString(f.query, "")
	f.query = spaces.ReplaceAllString(f.query, " ")

	//extract episde season numbers if they exist
	m := episeason.FindStringSubmatch(f.query)
	if len(m) > 0 {
		f.query = m[1] //trim name
		f.mtype = "series"
		f.season = m[3]
		f.episode = m[6]
	}

	//extract movie year
	if f.mtype == "" {
		m = year.FindStringSubmatch(f.query)
		if len(m) > 0 {
			f.query = m[1] //trim name
			f.mtype = "movie"
			f.year = m[2]
		}
	}

	//if the above fails, extract "Part 1/2/3..."
	if f.mtype == "" {
		m = partnum.FindStringSubmatch(f.query)
		if len(m) > 0 {
			f.query = m[1] //trim name
			f.mtype = "series"
			f.episode = m[2]
		}
	}

	// if f.mtype == "" {
	// 	return fmt.Errorf("No season/episode or year found")
	// }

	//trim spaces
	f.query = strings.TrimSpace(f.query)

	//search for normalized name
	r, err := search.Do(f.query, f.mtype)
	if err != nil {
		//not found
		return err
	}

	//calculate destination path
	dest := ""
	if r.Type == "series" {
		s, _ := strconv.Atoi(f.season)
		e, _ := strconv.Atoi(f.episode)
		filename := fmt.Sprintf("%s S%02dE%02d%s", r.Title, s, e, f.ext)
		dest = filepath.Join(f.s.c.TVDir, r.Title, filename)
	} else {
		filename := fmt.Sprintf("%s (%s)%s", r.Title, f.year, f.ext)
		dest = filepath.Join(f.s.c.MovieDir, filename)
	}

	//DEBUG
	// log.Printf("SUCCESS = D%d #%d\n  %s\n  %s", r.Distance, len(f.query), f.query, r.Title)
	log.Printf("Moving\n  '%s'\n  └─> '%s'", f.path, dest)

	if f.s.c.DryRun {
		return nil
	}

	//check already exists
	if _, err := os.Stat(dest); os.IsExist(err) {
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
