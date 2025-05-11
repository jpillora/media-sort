package mediasort

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	mediasearch "github.com/jpillora/media-sort/search"
)

// Sort the given path, creates a search query,
// performing the search and returning a Result
func Sort(path string) (*Result, error) {
	return SortThreshold(path, 95)
}

// SortThreshold sorts the given path, creates a search query,
// performing the search with the given threshold and returning a Result
func SortThreshold(path string, threshold int) (*Result, error) {
	return SortDepthThreshold(path, 0, threshold)
}

// SortDepthThreshold sorts the given path, includes <depth>
// parent directories, creates a search query,
// performing the search with the given threshold and returning a Result
func SortDepthThreshold(path string, depth, threshold int) (*Result, error) {
	r, err := runPathSort(path, threshold, depth)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Result holds both the results from parsing path, and the from
// performing the search
type Result struct {
	Query                         string
	Name, Path                    string
	Ext                           string
	MType                         string
	Season, Episode, ExtraEpisode int
	EpisodeDate                   string //weekly series
	Year                          string
	Accuracy                      int
	// Audiobook specific fields
	Author     string
	Narrator   string
	Series     string
	SeriesPart int
}

var (
	//DefaultTVTemplate defines the default TV path format
	DefaultTVTemplate = `{{ .Name }} S{{ printf "%02d" .Season }}E{{ printf "%02d" .Episode }}` +
		`{{ if ne .ExtraEpisode -1 }}-{{ printf "%02d" .ExtraEpisode }}{{end}}.{{ .Ext }}`
	//DefaultMovieTemplate defines the default movie path format
	DefaultMovieTemplate = "{{ .Name }} ({{ .Year }}).{{ .Ext }}"
	//DefaultAudiobookTemplate defines the default audiobook path format
	DefaultAudiobookTemplate = `{{ .Author }}/{{ if .Series }}{{ .Series }} - {{ printf "%02d" .SeriesPart }} - {{end}}{{ .Name }}{{ if .Narrator }} (Narrated by {{ .Narrator }}){{end}}.{{ .Ext }}`
)

// PathConfig customises the path templates
type PathConfig struct {
	TVTemplate        string `help:"tv series path template"`
	MovieTemplate     string `help:"movie path template"`
	AudiobookTemplate string `help:"audiobook path template"`
}

var prettyPathFuncs = template.FuncMap{}

// PrettyPath converts the provided "messy" path into a
// "pretty" cleanly formatted path using the media result
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
	case mediasearch.Audiobook:
		tmpl = config.AudiobookTemplate
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

func runPathParse(path string, depth int) (Result, error) {
	result := Result{
		Path:         path,
		Season:       1,
		Episode:      -1,
		ExtraEpisode: -1,
	}
	if sample.MatchString(strings.ToLower(path)) {
		return result, fmt.Errorf("Skipped sample media")
	}
	dir, name := filepath.Split(path)
	ext := getExtension(name)
	name = strings.TrimSuffix(name, ext)
	//add depth*parts of dir onto name
	dir = strings.Trim(dir, sep)
	parts := []string{}
	if dir != "" {
		parts = strings.Split(dir, sep)
	}
	l := len(parts)
	if depth < 0 || depth > l {
		depth = l
	}

	// For audiobooks, check if the parent directory is the author and extract title
	if ext != "" && (strings.EqualFold(ext, ".m4b") ||
		strings.EqualFold(ext, ".mp3") ||
		strings.EqualFold(ext, ".aac") ||
		strings.EqualFold(ext, ".flac") ||
		strings.EqualFold(ext, ".aa") ||
		strings.EqualFold(ext, ".aax")) {
		result.MType = string(mediasearch.Audiobook)
		if len(parts) > 0 {
			result.Author = parts[len(parts)-1] // Last directory is author
		}
		result.Query = mediasearch.Normalize(name) // Use just the normalized filename as the query
	}

	name = strings.Join(append(parts[l-depth:], name), " ")
	//split name/ext
	result.Name = name
	result.Ext = strings.TrimPrefix(ext, ".")
	//query is normalized name
	query := mediasearch.Normalize(name)
	log.Printf("'%s' -> '%s'", name, query)

	if result.MType != string(mediasearch.Audiobook) {
		result.Query = query // Set query only for non-audiobook types
	}

	// If it's an audiobook, query should be the normalized title
	if result.MType == string(mediasearch.Audiobook) {
		query = mediasearch.Normalize(result.Name)
	}

	// Try to identify audiobook first
	if tryAudiobookParse(&result, query) {
		return result, nil
	}

	//extract episode date (weekly show)
	if result.MType == "" {
		m := epidate.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			result.MType = string(mediasearch.Series)
			result.EpisodeDate = strings.Replace(m[2], " ", "-", 2)
		}
	}
	//extract double episode season numbers
	if result.MType == "" {
		m := doubleepiseason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			result.MType = string(mediasearch.Series)
			result.Season, _ = strconv.Atoi(m[2])
			result.Episode, _ = strconv.Atoi(m[4])
			result.ExtraEpisode, _ = strconv.Atoi(m[6])
		}
	}
	//extract episode season numbers
	if result.MType == "" {
		m := episeason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			result.MType = string(mediasearch.Series)
			result.Season, _ = strconv.Atoi(m[3])
			result.Episode, _ = strconv.Atoi(m[6])
		}
		//second chance, search untrimmed name
		if len(m) == 0 {
			m = episeason.FindStringSubmatch(name)
		}
		if len(m) > 0 {
			//cant't trim query
			result.MType = string(mediasearch.Series)
			result.Season, _ = strconv.Atoi(m[3])
			result.Episode, _ = strconv.Atoi(m[6])
		}
	}
	//remove phrase "season X" from tv series queries
	if result.MType == string(mediasearch.Series) && season.MatchString(query) {
		//and re-noramlise
		query = mediasearch.Normalize(season.ReplaceAllString(query, ""))
	}
	//extract *joined* episode season numbers
	if result.MType == "" {
		m := joinedepiseason.FindStringSubmatch(query)
		if len(m) > 0 {
			query = m[1] //trim name
			result.MType = string(mediasearch.Series)
			result.Season, _ = strconv.Atoi(m[2])
			result.Episode, _ = strconv.Atoi(m[3])
		}
	}
	//extract release year
	m := year.FindStringSubmatch(query)
	if len(m) > 0 {
		query = m[1] //trim name
		if result.MType == "" {
			result.MType = string(mediasearch.Movie) //set type to "movie", if not already set
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
	//ready for search
	return result, nil
}

func tryAudiobookParse(result *Result, query string) bool {
	// Try Author - Title format
	if m := authorTitle.FindStringSubmatch(query); len(m) > 0 {
		result.Author = m[1]
		result.Name = m[2]
		result.MType = string(mediasearch.Audiobook)
		return true
	}

	// Try Title by Author format
	if m := authorBy.FindStringSubmatch(query); len(m) > 0 {
		result.Name = m[1]
		result.Author = m[2]
		result.MType = string(mediasearch.Audiobook)
		return true
	}

	// Look for narrator information
	if m := narratorParen.FindStringSubmatch(query); len(m) > 0 {
		result.Narrator = m[1]
		// Remove narrator info from query for further parsing
		query = narratorParen.ReplaceAllString(query, "")
	}

	// Try Series Book format
	if m := seriesBook.FindStringSubmatch(query); len(m) > 0 {
		result.SeriesPart, _ = strconv.Atoi(m[1])
		if m[2] != "" {
			result.Series = m[2]
		}
		result.Name = m[3]
		result.MType = string(mediasearch.Audiobook)
		return true
	}

	// Common audiobook extensions
	switch strings.ToLower(result.Ext) {
	case "m4b", "aa", "aax":
		result.MType = string(mediasearch.Audiobook)
		return true
	}

	return false
}

func runPathSort(path string, threshold, depth int) (Result, error) {
	result, err := runPathParse(path, depth)
	if err != nil {
		return result, err
	}
	//search for normalized name
	searchResult, err := mediasearch.SearchThreshold(result.Query, result.Year, result.MType, threshold)
	if err != nil {
		return result, err
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
