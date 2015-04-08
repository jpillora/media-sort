package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jpillora/mediasort/ms"
)

var VERSION string = "0.0.0" //set via ldflags

var help = `
	Usage: mediasort [options] [file/directory]

	MediaSort performs a simple categorization on the
	provided file or directory.

	Movies are moved to:
		<movie-dir>/<title> (<year>)
	TV Shows are moved to:
		<tv-dir>/<title>/<title> S<season>E<episode>

	Version: ` + VERSION + `

	Options:
	--config-file, A JSON file describing these options.
	--movie-dir -m, The destination movie directory
	(defaults to $HOME/movies).
	--tv-dir -t, The destination TV directory (defaults
	to $HOME/tv).
	--ext, Extensions considered (defaults to "mp4,avi,
	mkv")
	--watch, Watches the provided directory for changes.
	and sorts new files as they arrive.
	--version -v, Display version.
	--help -h, This help text.

	Read more:
	  https://github.com/jpillora/serve
`

func main() {

	// cpath := flag.String(&c.Config, "config-file", "", "")

	//fill sorter config
	c := &ms.Config{}
	flag.StringVar(&c.MovieDir, "movie-dir", "", "")
	flag.StringVar(&c.MovieDir, "m", "", "")
	flag.StringVar(&c.TVDir, "tv-dir", "", "")
	flag.StringVar(&c.TVDir, "t", "", "")
	flag.StringVar(&c.Exts, "ext", "mp4,avi,mkv", "")

	//meta cli
	h := false
	flag.BoolVar(&h, "h", false, "")
	flag.BoolVar(&h, "help", false, "")
	v := false
	flag.BoolVar(&v, "v", false, "")
	flag.BoolVar(&v, "version", false, "")

	//parse
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, help)
		os.Exit(1)
	}
	if v {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	if h {
		flag.Usage()
	}

	//get directory
	args := flag.Args()
	if len(args) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		c.Target = cwd
	} else {
		c.Target = args[0]
	}

	//ready!
	s, err := ms.NewSorter(c)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
