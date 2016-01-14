package mediasort

import (
	"regexp"
	"strings"
)

var yearstr = `(19\d\d|20[0,1]\d)`
var sample = regexp.MustCompile(`\bsample\b`)
var encodings = regexp.MustCompile(`\b(720p|1080p|hdtv|x264|dts|bluray)\b.*`) //strip all junk
var nonalpha = regexp.MustCompile(`[^A-Za-z0-9]`)
var spaces = regexp.MustCompile(`\s+`)
var episeason = regexp.MustCompile(`^(.+?)\bs?(eason)?(\d{1,2})(e|\ |\ e|x|xe)(pisode)?(\d{1,2})\b`)
var epidate = regexp.MustCompile(`^(.+?\b)(` + yearstr + ` \d{2} \d{2}|\d{2} \d{2} ` + yearstr + `)\b`)
var year = regexp.MustCompile(`^(.+?\b)` + yearstr + `\b`)
var joinedepiseason = regexp.MustCompile(`^(.+?\b)(\d)(\d{2})\b`)
var partnum = regexp.MustCompile(`^(.+?\b)(\d{1,2})\b`)
var extRe = regexp.MustCompile(`\.\w+$`)

var apost = regexp.MustCompile(`'`)
var colon = regexp.MustCompile(`:`)
var invalidChars = regexp.MustCompile(`[^a-zA-Z0-9\_\-\.\ \(\)]`)
var doubleDash = regexp.MustCompile(`-(\s*-\s*)+ `)
var doubleSpace = regexp.MustCompile(`\s+`)

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
