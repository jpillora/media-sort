package mediasort

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jpillora/media-sort/search"
)

func Sort(s string) (search.Result, error) {
	return ps.run()
}

//PathSort converts the provided path into a
//well-formatted path based of live media data
func PathSort(path string) (string, error) {

	ps, err := runPathSort(path)
	if err != nil {
		return "", err
	}
	if ps.mtype == "series" && ps.episodeDate != "" {
		return filepath.Join(ps.name, fmt.Sprintf("%s %s%s", ps.name, ps.episodeDate, ps.ext)), nil
	} else if ps.mtype == "series" && ps.episode != "" {
		s, _ := strconv.Atoi(ps.season)
		e, _ := strconv.Atoi(ps.episode)
		return filepath.Join(ps.name, fmt.Sprintf("%s S%02dE%02d%s", ps.name, s, e, ps.ext)), nil
	} else if ps.mtype == "series" {
		return "", fmt.Errorf("TV Series with no episode found '%s'", ps.query)
	} else if ps.mtype == "movie" {
		return fmt.Sprintf("%s (%s)%s", ps.name, ps.year, ps.ext), nil
	}
	return "", fmt.Errorf("Search failed")
}

type pathSort struct {
	query           string
	name, path      string
	ext             string
	mtype           string
	season, episode string
	episodeDate     string //weekly series
	year            string
}

func runPathSort(path string) (*pathSort, error) {
	if sample.MatchString(strings.ToLower(path)) {
		return nil, fmt.Errorf("Skipped sample media")
	}
	_, name := filepath.Split(path)
	ext := getExtension(name)
	name = strings.TrimSuffix(name, ext)
	ps := &pathSort{
		name:    name,
		path:    path,
		ext:     ext,
		season:  "1",
		episode: "",
	}
	//normalize name
	query := normalize(ps.name)
	//extract episode date (weekly show)
	if ps.mtype == "" {
		m := epidate.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			ps.mtype = "series"
			ps.episodeDate = strings.Replace(m[2], " ", "-", 2)
		}
	}
	//extract episde season numbers
	if ps.mtype == "" {
		m := episeason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			ps.mtype = "series"
			ps.season = m[3]
			ps.episode = m[6]
		}
	}
	//extract *joined* episde season numbers
	if ps.mtype == "" {
		m := joinedepiseason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			ps.mtype = "series"
			ps.season = m[2]
			ps.episode = m[3]
		}
	}
	//extract release year
	m := year.FindStringSubmatch(query)
	if len(m) > 0 {
		query = m[1] //trim name
		if ps.mtype == "" {
			ps.mtype = "movie" //set type to "movie", if not already set
		}
		ps.year = m[2]
	}
	//if the above fails, extract "Part 1/2/3..."
	if ps.mtype == "" {
		m := partnum.FindStringSubmatch(query)
		if len(m) > 0 {
			ps.episode = m[2]
		}
	}
	//trim spaces
	ps.query = strings.TrimSpace(query)
	//search for normalized name
	r, err := mediasearch.Search(ps.query, ps.year, ps.mtype)
	if err != nil {
		return nil, err //search failed
	}
	//use results
	ps.name = r.Title
	if r.Type == mediasearch.Series && r.IsDupe { //differentiate duplicates by year
		ps.name += " (" + r.Year + ")"
	}
	ps.year = r.Year
	ps.mtype = string(r.Type)
	return ps, nil
}
