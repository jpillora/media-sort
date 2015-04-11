package mediasort

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
)

//Sorter is a media sorter
type Sorter struct {
	c                       Config
	exts                    map[string]bool
	files                   map[string]*fileSorter //abs-path -> file
	Checked, Matched, Moved int
}

//NewSorter creates a new Sorter
func New(c Config) (*Sorter, error) {

	if len(c.Targets) == 0 {
		return nil, fmt.Errorf("Please provide at least one file or directory")
	}

	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	if c.MovieDir == "" {
		c.MovieDir = path.Join(home, "movies")
	}
	if c.TVDir == "" {
		c.TVDir = path.Join(home, "tv")
	}

	exts := map[string]bool{}
	for _, e := range strings.Split(c.Exts, ",") {
		exts["."+e] = true
	}

	s := &Sorter{
		c:     c,
		exts:  exts,
		files: map[string]*fileSorter{},
	}

	//convert all targets into file sorter objects
	for _, path := range c.Targets {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		err = s.addFile(path, info, 0)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Sorter) addFile(path string, info os.FileInfo, depth int) error {

	//limit recursion depth
	if depth > s.c.Depth {
		return nil
	}

	if info.Mode().IsRegular() {
		s.Checked++
		//add single file
		f, err := newFileSorter(s, path, info)
		if f == nil || err != nil {
			return err
		}
		s.files[path] = f
		s.Matched++
		return nil
	}

	if info.IsDir() {
		//TODO watch directory
		// if s.c.Watch {
		// }

		//add all files in dir
		infos, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, info := range infos {
			p := filepath.Join(path, info.Name())
			//recurse
			err := s.addFile(p, info, depth+1)
			if err != nil {
				return err
			}
		}
	}

	//skip links,pipes,etc
	return nil
}

func (s *Sorter) sortFiles() []error {
	var errors []error
	//parallel sort all files
	wg := &sync.WaitGroup{}
	for _, f := range s.files {
		wg.Add(1)
		go f.goRun(wg)
	}
	wg.Wait()

	//collect errors
	for _, f := range s.files {
		if f.err != nil {
			err := fmt.Errorf("%s: %s", f.name, f.err)
			errors = append(errors, err)
		} else {
			s.Moved++
		}
	}

	return errors
}

func (s *Sorter) Run() []error {

	if s.c.DryRun {
		log.Println("[Dryrun Mode]")
	}

	//just run once
	if !s.c.Watch {
		return s.sortFiles()
	}

	//run watcher
	// for {
	// 	s.sortFiles()
	// }

	return []error{fmt.Errorf("Watch not implemented")}
}
