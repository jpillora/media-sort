package mediasort

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	mediasearch "github.com/jpillora/media-sort/search"
)

func Sort(path string) (*Result, error) {
	return SortThreshold(path, 95)
}

func SortThreshold(path string, threshold int) (*Result, error) {
	return runPathSort(path, threshold)
}

type Result struct {
	Query                         string
	Name, Path                    string
	Ext                           string
	MType                         string
	Season, Episode, ExtraEpisode int
	EpisodeDate                   string //weekly series
	Year                          string
	Accuracy                      int
}

var (
	DefaultTVTemplate = `{{ .Name }} S{{ printf "%02d" .Season }}E{{ printf "%02d" .Episode }}` +
		`{{ if ne .ExtraEpisode -1 }}-{{ printf "%02d" .ExtraEpisode }}{{end}}.{{ .Ext }}`
	DefaultMovieTemplate = "{{ .Name }} ({{ .Year }}).{{ .Ext }}"
)

type PathConfig struct {
	TVTemplate    string `help:"tv series path template"`
	MovieTemplate string `help:"movie path template"`
}

var prettyPathFuncs = template.FuncMap{}

//PrettyPath converts the provided "messy" path into a
//"pretty" cleanly formatted path using the media result
func (result *Result) PrettyPath(config PathConfig) (string, error) {
	//config
	if config.TVTemplate == "" {
		config.TVTemplate = DefaultTVTemplate
	}
	if config.MovieTemplate == "" {
		config.MovieTemplate = DefaultMovieTemplate
	}
	//find template
	tmpl := ""
	switch mediasearch.MediaType(result.MType) {
	case mediasearch.Series:
		tmpl = config.TVTemplate
	case mediasearch.Movie:
		tmpl = config.MovieTemplate
	default:
		return "", fmt.Errorf("Invalid result type: %s", result.MType)
	}
	//run template
	str := bytes.Buffer{}
	if t, err := template.New("t").Funcs(prettyPathFuncs).Parse(tmpl); err != nil {
		return "", err
	} else if err := t.Execute(&str, result); err != nil {
		return "", err
	}

	prettyPath := fixPath(str.String())
	return prettyPath, nil
}

func runPathSort(path string, threshold int) (*Result, error) {
	if sample.MatchString(strings.ToLower(path)) {
		return nil, fmt.Errorf("Skipped sample media")
	}
	_, name := filepath.Split(path)
	ext := getExtension(name)
	name = strings.TrimSuffix(name, ext)
	result := &Result{
		Name:         name,
		Path:         path,
		Ext:          strings.TrimPrefix(ext, "."),
		Season:       1,
		Episode:      -1,
		ExtraEpisode: -1,
	}
	//normalize name
	query := mediasearch.Normalize(result.Name)
	//extract episode date (weekly show)
	if result.MType == "" {
		m := epidate.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			result.MType = "series"
			result.EpisodeDate = strings.Replace(m[2], " ", "-", 2)
		}
	}
	//extract double episde season numbers
	if result.MType == "" {
		m := doubleepiseason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			result.MType = "series"
			result.Season, _ = strconv.Atoi(m[2])
			result.Episode, _ = strconv.Atoi(m[4])
			result.ExtraEpisode, _ = strconv.Atoi(m[6])
		}
	}
	//extract episde season numbers
	if result.MType == "" {
		m := episeason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			result.MType = "series"
			result.Season, _ = strconv.Atoi(m[3])
			result.Episode, _ = strconv.Atoi(m[6])
		}
	}
	//extract *joined* episde season numbers
	if result.MType == "" {
		m := joinedepiseason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			result.MType = "series"
			result.Season, _ = strconv.Atoi(m[2])
			result.Episode, _ = strconv.Atoi(m[3])
		}
	}
	//extract release year
	m := year.FindStringSubmatch(query)
	if len(m) > 0 {
		query = m[1] //trim name
		if result.MType == "" {
			result.MType = "movie" //set type to "movie", if not already set
		}
		result.Year = m[2]
	}
	//if the above fails, extract "Part 1/2/3..."
	if result.MType == "" {
		m := partnum.FindStringSubmatch(query)
		if len(m) > 0 {
			result.Episode, _ = strconv.Atoi(m[2])
		}
	}
	//trim spaces
	result.Query = strings.TrimSpace(query)
	//search for normalized name
	searchResult, err := mediasearch.SearchThreshold(result.Query, result.Year, result.MType, threshold)
	if err != nil {
		return nil, err //search failed
	}
	//use results
	result.Name = searchResult.Title
	if searchResult.Type == mediasearch.Series && searchResult.IsDupe { //differentiate duplicates by year
		result.Name += " (" + searchResult.Year + ")"
	}
	result.Year = searchResult.Year
	result.MType = string(searchResult.Type)
	result.Accuracy = searchResult.Accuracy
	return result, nil
}
