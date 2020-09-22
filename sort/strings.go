package mediasort

import (
	"path/filepath"
	"regexp"
)

//NOTE strings have been mediasearch.Normalized before these regexps run over them
var (
	yearstr         = `(19\d\d|20\d\d)`
	sample          = regexp.MustCompile(`\bsample\b`)
	encodings       = regexp.MustCompile(`\b(720p|1080p|hdtv|x264|dts|bluray)\b.*`) //strip all junk
	nonalpha        = regexp.MustCompile(`[^A-Za-z0-9]`)
	spaces          = regexp.MustCompile(`\s+`)
	doubleepiseason = regexp.MustCompile(`^(.+?)\bs?(\d{1,2})(e||x|xe)(\d{2}).?(e||x|xe)(\d{2})\b`)
	episeason       = regexp.MustCompile(`(?i)^(.+?)\bs?(eason)?(\d{1,2})(e|\ |\ e|x|xe)(pisode)?(\d{1,2})\b`)
	epidate         = regexp.MustCompile(`^(.+?\b)(` + yearstr + ` \d{2} \d{2}|\d{2} \d{2} ` + yearstr + `)\b`)
	season          = regexp.MustCompile(`(?i)\bseason[\s\.\-]*\d{1,2}\b`)
	year            = regexp.MustCompile(`^(.+?\b)` + yearstr + `\b`)
	joinedepiseason = regexp.MustCompile(`^(.+?\b)(\d)(\d{2})\b`)
	partnum         = regexp.MustCompile(`^(.+?\b)(\d{1,2})\b`)
	partof          = regexp.MustCompile(`(?i)^(.+?\b)(\d{1,3})\s*of\s*\d{1,3}\b`)
	extRe           = regexp.MustCompile(`\.\w+$`)
	apost           = regexp.MustCompile(`'`)
	colon           = regexp.MustCompile(`:`)
	invalidChars    = regexp.MustCompile(`[^\p{Greek}\pP\pN\p{L}\_\-\.\ \(\)\/\\]`)
	doubleDash      = regexp.MustCompile(`-(\s*-\s*)+ `)
	doubleSpace     = regexp.MustCompile(`\s+`)
	sep             = string(filepath.Separator)
)

func fixPath(s string) string {
	s = apost.ReplaceAllString(s, "")
	s = colon.ReplaceAllString(s, " -")
	s = invalidChars.ReplaceAllString(s, "-")
	s = doubleDash.ReplaceAllString(s, "- ")
	s = doubleSpace.ReplaceAllString(s, " ")
	return s
}

func getExtension(s string) string {
	return extRe.FindString(s)
}
