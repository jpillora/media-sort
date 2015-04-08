package ms

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jpillora/mediasort/ms/search"
)

type fileSorter struct {
	name, dir string
	info      os.FileInfo

	tvSe, tvEp int

	moviePts, tvPts int
}

func (f *fileSorter) run() error {

	f.name = strings.ToLower(f.name)
	f.name = strings.Replace(f.name, ".", " ", -1)

	f.setEpisodeSeason()

	//ready to be sorted
	// results, err := search.Do(f.name)
	// if err != nil {
	// 	return nil
	// }

	if f.tvPts > 0 {
		results, err := search.Do(f.name + " tv show site:imdb.com")
		if err != nil {
			return nil
		}
		log.Printf("SUCCESS %s #%d", f.name, len(results))
		for _, r := range results {
			fmt.Printf("  %s\n", r.Title)
		}
	} else {
		log.Println("FAIL", f.name)
	}

	return nil
}

// var epSeTypes = []regexp.Regexp{
// 	,
// 	regexp.MustCompile(`\bseason(\d{1,2})episode(\d{1,3})\b`),
// 	regexp.MustCompile(`\bseason(\d{1,2})episode(\d{1,3})\b`),
// }

var epSeType = regexp.MustCompile(`^(.+?)\bs?(eason)?(\d{1,2})(e|\ |\ e|x|xe)(pisode)?(\d{1,2})\b`)

func (f *fileSorter) setEpisodeSeason() {
	m := epSeType.FindStringSubmatch(f.name)
	if len(m) == 0 {
		return
	}
	f.tvPts += 10
	f.name = m[1]
	f.tvSe, _ = strconv.Atoi(m[3])
	f.tvEp, _ = strconv.Atoi(m[6])
}

var junk = regexp.MustCompile(`\b(720p|1080p|HDTV|)\b`)

func (f *fileSorter) stripJunk() {

}
