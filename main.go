package main

import (
	"log"
	"time"

	"github.com/jpillora/media-sort/sort"
	"github.com/jpillora/opts"
)

var VERSION string = "0.0.0-src" //set via ldflags

const info = `
media-sort categorizes the provided files and directories by moving
them into to a structured directory tree, using various live sources.

by default, tv series are moved to:
  <tv-dir>/<title> S<season>E<episode>.<ext>
and movies are moved to:
  <movie-dir>/<title> (<year>).<ext>
`

func main() {

	c := mediasort.Config{
		Extensions:  "mp4,avi,mkv",
		Concurrency: 6,
		FileLimit:   1000,
		WatchDelay:  3 * time.Second,
	}

	opts.New(&c).
		Name("media-sort").
		Repo("github.com/jpillora/media-sort").
		DocAfter("usage", "info", info).
		Version(VERSION).
		Parse()

	if err := mediasort.FileSystemSort(c); err != nil {
		log.Fatal(err)
	}
}
