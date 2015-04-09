package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jpillora/mediasort/mediasort"
)

var VERSION string = "0.0.0" //set via ldflags

var help = `
	Usage: mediasort [options] [file/directory]

	MediaSort performs a simple categorization on the provided file or directory.

	Movies are moved to:
		<movie-dir>/<title> (<year>).ext
	TV Shows are moved to:
		<tv-dir>/<title>/<title> S<season>E<episode>.ext

	Version: ` + VERSION + `

	Options:
	--movie-dir -m, The destination movie directory (defaults to $HOME/movies).
	--tv-dir -t,    The destination TV directory (defaults to $HOME/tv).
	--ext,          Extensions considered (defaults to "mp4,avi,mkv")
	--dryrun,       Runs in read-only mode
	--version -v,   Display version.
	--help -h,      This help text.

	Read more:
	  https://github.com/jpillora/mediasort
`

var todo = `
	--config-file, A JSON file describing these options.
	--watch, Watches the provided directory for changes.
	and sorts new files as they arrive.
`

func main() {

	// cpath := flag.String(&c.Config, "config-file", "", "")

	//fill sorter config
	c := &mediasort.Config{}
	flag.StringVar(&c.MovieDir, "movie-dir", "", "")
	flag.StringVar(&c.MovieDir, "m", "", "")
	flag.StringVar(&c.TVDir, "tv-dir", "", "")
	flag.StringVar(&c.TVDir, "t", "", "")
	flag.StringVar(&c.Exts, "ext", "mp4,avi,mkv", "")
	flag.IntVar(&c.Depth, "depth", 1, "")
	flag.BoolVar(&c.DryRun, "dryrun", false, "")
	flag.BoolVar(&c.Watch, "watch", false, "")

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
	c.Targets = flag.Args()
	if len(c.Targets) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		c.Targets = []string{cwd}
	}

	//ready!
	s, err := mediasort.New(c)
	if err != nil {
		log.Fatal(err)
	}
	errs := s.Run()

	log.Printf("Checked #%d, Matched #%d, Moved #%d", s.Checked, s.Matched, s.Moved)

	if len(errs) > 0 {
		log.Printf("Encountered #%d errors:", len(errs))
		for i, err := range errs {
			fmt.Printf(" [%2d] %s\n", i+1, err)
		}
		os.Exit(1)
	}
}
