package mediasearch

import (
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
)

var (
	nonalpha        = regexp.MustCompile(`[^a-z0-9]`)
	yearstr         = `(19\d\d|20\d\d)`
	onlyYear        = regexp.MustCompile(`^` + yearstr + `$`)
	getYear         = regexp.MustCompile(`\b` + yearstr + `\b`)
	getDate         = regexp.MustCompile(`\b` + yearstr + `-(\d\d)-(\d\d)\b`)
	sample          = regexp.MustCompile(`\bsample\b`)
	encodings       = regexp.MustCompile(`\b(720p|1080p|hdtv|x264|dts|bluray)\b.*`) //strip all junk
	spaces          = regexp.MustCompile(`\s+`)
	episeason       = regexp.MustCompile(`^(.+?)\bs?(eason)?(\d{1,2})(e|\ |\ e|x|xe)(pisode)?(\d{1,2})\b`)
	epidate         = regexp.MustCompile(`^(.+?\b)(` + yearstr + ` \d{2} \d{2}|\d{2} \d{2} ` + yearstr + `)\b`)
	year            = regexp.MustCompile(`^(.+?\b)` + yearstr + `\b`)
	joinedepiseason = regexp.MustCompile(`^(.+?\b)(\d)(\d{2})\b`)
	partnum         = regexp.MustCompile(`^(.+?\b)(\d{1,2})\b`)
)

// Normalize strings to become search terms
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
