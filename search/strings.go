package mediasearch

import (
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
)

var nonalpha = regexp.MustCompile(`[^a-z0-9]`)
var yearstr = `(19\d\d|20[0,1]\d)`
var onlyYear = regexp.MustCompile(`^` + yearstr + `$`)
var getYear = regexp.MustCompile(`\b` + yearstr + `\b`)
var sample = regexp.MustCompile(`\bsample\b`)
var encodings = regexp.MustCompile(`\b(720p|1080p|hdtv|x264|dts|bluray)\b.*`) //strip all junk
var spaces = regexp.MustCompile(`\s+`)
var episeason = regexp.MustCompile(`^(.+?)\bs?(eason)?(\d{1,2})(e|\ |\ e|x|xe)(pisode)?(\d{1,2})\b`)
var epidate = regexp.MustCompile(`^(.+?\b)(` + yearstr + ` \d{2} \d{2}|\d{2} \d{2} ` + yearstr + `)\b`)
var year = regexp.MustCompile(`^(.+?\b)` + yearstr + `\b`)
var joinedepiseason = regexp.MustCompile(`^(.+?\b)(\d)(\d{2})\b`)
var partnum = regexp.MustCompile(`^(.+?\b)(\d{1,2})\b`)

// Normalize strings to become search terms
// TODO var romannumerals...
func Normalize(s string) string {
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

func accuracy(a, b string) int {
	return 100 - dist(a, b)
}

func dist(a, b string) int {
	a = Normalize(a)
	b = Normalize(b)
	return levenshtein.ComputeDistance(a, b)
}

func isNear(a, b string) bool {
	lendiff := abs(len(a) - len(b))
	return dist(a, b)-lendiff <= 5
}

// return nil, fmt.Errorf("Failed to match '%s' (closest result was '%s')", ps.name, r.Title)
