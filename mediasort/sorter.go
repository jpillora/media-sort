package mediasort

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

//Sorter is a media sorter
type Sorter struct {
	c                       *Config
	exts                    map[string]bool
	files                   map[string]*fileSorter //abs-path -> file
	Checked, Matched, Moved int
}

//NewSorter creates a new Sorter
func New(c *Config) (*Sorter, error) {

	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	if c.MovieDir == "" {
		c.MovieDir = path.Join(u.HomeDir, "movies")
	}
	if c.TVDir == "" {
		c.TVDir = path.Join(u.HomeDir, "tv")
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

	//convert targets into file sorter objects
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

	//only go 1 level deep
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
			err := s.addFile(p, info, depth+1)
			if err != nil {
				return err
			}
		}
	}

	//skip links,pipes,etc
	return nil
}

func (s *Sorter) Run() []error {
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
