package ms

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

//Sorter is a media sorter
type Sorter struct {
	c    *Config
	exts map[string]bool
}

//NewSorter creates a new Sorter
func NewSorter(c *Config) (*Sorter, error) {

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

	return &Sorter{
		c:    c,
		exts: exts,
	}, nil
}

func (s *Sorter) Run() error {

	t := s.c.Target

	info, err := os.Stat(t)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if err := s.sortDir(t); err != nil {
			return err
		}
	} else if info.Mode().IsRegular() {
		dir := filepath.Dir(t)
		if err := s.sortFile(dir, info); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Invalid file: %s", t)
	}

	return nil
}

func (s *Sorter) sortDir(dir string) error {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	//parallel sort all files
	wg := sync.WaitGroup{}
	wg.Add(len(infos))
	for _, i := range infos {
		go func(info os.FileInfo) {
			defer wg.Done()
			err := s.sortFile(dir, info)
			if err != nil {
				log.Println(err)
			}
		}(i)
	}
	wg.Wait()

	return nil
}

func (s *Sorter) sortFile(dir string, info os.FileInfo) error {
	//attempt to rule out file
	if !info.Mode().IsRegular() {
		return nil
	}
	name := info.Name()
	if _, exists := s.exts[filepath.Ext(name)]; !exists {
		return nil
	}

	f := fileSorter{
		name: name,
		dir:  dir,
		info: info,
	}

	return f.run()
}
