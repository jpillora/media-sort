package main

import (
	"log"
	"time"

	"github.com/jpillora/media-sort/sort"
	"github.com/jpillora/opts"
)

var VERSION string = "0.0.0-src" //set via ldflags

const (
	info = `
media-sort categorizes the provided files and directories (targets) by
moving them into to a structured directory tree, sorting is currently
performed using TVMaze, OMDB and Google.
`
	pathTemplates = `
by default, tv series are moved to:
  ./<title> S<season>E<episode>.<ext>
and movies are moved to:
  ./<title> (<year>).<ext>

to modify the these paths, you can use the --tv-template and
--movie-template options. These options describe the new file path for
tv series and movies using Go template syntax. You can find the
default values here:
  https://godoc.org/github.com/jpillora/media-sort/sort#pkg-variables
and you can view all possible template variables here:
  https://godoc.org/github.com/jpillora/media-sort/sort#Result
`
)

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
		DocAfter("options", "pathtemplates", pathTemplates).
		Version(VERSION).
		Parse()

	if err := mediasort.FileSystemSort(c); err != nil {
		log.Fatal(err)
	}
}
