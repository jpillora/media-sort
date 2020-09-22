package main

import (
	"log"
	"time"

	mediasort "github.com/jpillora/media-sort/sort"
	"github.com/jpillora/opts"
	"github.com/jpillora/sizestr"
)

var version = "0.0.0-src" //set via ldflags

const (
	info = `
media-sort categorizes the provided files and directories (targets) by
moving them into to a structured directory tree, sorting is currently
performed using TVMaze, MovieDB and Google.
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
		Extensions:        "mp4,m4v,avi,mkv,mpeg,mpg,mov,webm",
		Concurrency:       6,
		FileLimit:         1000,
		MinFileSize:       sizestr.Bytes(sizestr.MustParse("25MB")),
		WatchDelay:        3 * time.Second,
		AccuracyThreshold: 95, //100 is perfect match,
		Action:            mediasort.MoveAction,
	}

	opts.New(&c).
		Name("media-sort").
		Repo("github.com/jpillora/media-sort").
		DocAfter("usage", "info", info).
		DocBefore("version", "pathtemplates", pathTemplates).
		SetLineWidth(128).
		Version(version).
		Parse()

	if err := mediasort.FileSystemSort(c); err != nil {
		log.Fatal(err)
	}
}
