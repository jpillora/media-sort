package mediasort

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	mediasearch "github.com/jpillora/media-sort/search"
	"github.com/jpillora/sizestr"

	"gopkg.in/fsnotify.v1"
)

//Config is a sorter configuration
type Config struct {
	Targets           []string `opts:"mode=arg,min=1"`
	TVDir             string   `opts:"help=tv series base directory (defaults to current directory)"`
	MovieDir          string   `opts:"help=movie base directory (defaults to current directory)"`
	PathConfig        `mode:"embedded"`
	Extensions        string        `opts:"help=types of files that should be sorted"`
	Concurrency       int           `opts:"help=search concurrency [warning] setting this too high can cause rate-limiting errors"`
	FileLimit         int           `opts:"help=maximum number of files to search"`
	AccuracyThreshold int           `opts:"help=filename match accuracy threshold" default:"is 95, perfect match is 100"`
	MinFileSize       sizestr.Bytes `opts:"help=minimum file size"`
	Recursive         bool          `opts:"help=also search through subdirectories"`
	DryRun            bool          `opts:"help=perform sort but don't actually move any files"`
	SkipHidden        bool          `opts:"help=skip dot files"`
	Hardlink          bool          `opts:"help=hardlink files to the new location instead of moving"`
	Overwrite         bool          `opts:"help=overwrites duplicates"`
	OverwriteIfLarger bool          `opts:"help=overwrites duplicates if the new file is larger"`
	Watch             bool          `opts:"help=watch the specified directories for changes and re-sort on change"`
	WatchDelay        time.Duration `opts:"help=delay before next sort after a change"`
	Verbose           bool          `opts:"help=verbose logs"`
}

//fsSort is a media sorter
type fsSort struct {
	Config
	validExts map[string]bool
	sorts     map[string]*fileSort
	dirs      map[string]bool
	stats     struct {
		found, matched, moved int
	}
}

type fileSort struct {
	id     int
	path   string
	info   os.FileInfo
	result *Result
	err    error
}

//FileSystemSort performs a media sort
//against the file system using the provided
//configuration
func FileSystemSort(c Config) error {
	if c.MovieDir == "" {
		c.MovieDir = "."
	}
	if c.TVDir == "" {
		c.TVDir = "."
	}
	if c.Watch && !c.Recursive {
		return errors.New("Recursive mode is required to watch directories")
	}
	if c.Overwrite && c.OverwriteIfLarger {
		return errors.New("Overwrite is already specified, overwrite-if-larger is redundant")
	}
	if c.Hardlink && c.Overwrite {
		return errors.New("Hardlink is already specified, Overwrite won't do anything")
	}
	//init fs sort
	fs := &fsSort{
		Config:    c,
		validExts: map[string]bool{},
	}
	for _, e := range strings.Split(c.Extensions, ",") {
		fs.validExts["."+e] = true
	}
	//sort loop
	for {
		//reset state
		fs.sorts = map[string]*fileSort{}
		fs.dirs = map[string]bool{}
		//look for files
		if err := fs.scan(); err != nil {
			return err
		}
		//ensure we have dirs to watch
		if fs.Watch && len(fs.dirs) == 0 {
			return errors.New("No directories to watch")
		}
		if len(fs.sorts) > 0 {
			//moment of truth - sort all files!
			if err := fs.sortAllFiles(); err != nil {
				return err
			}
		}
		//watch directories
		if !c.Watch {
			break
		}
		if err := fs.watch(); err != nil {
			return err
		}
	}
	return nil
}

func (fs *fsSort) scan() error {
	fs.verbf("scanning targets...")
	//scan targets for media files
	for _, path := range fs.Targets {
		fs.verbf("scanning: %s", path)
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if err = fs.add(path, info); err != nil {
			return err
		}
	}
	//ensure we found something
	if len(fs.sorts) == 0 && (!fs.Watch || len(fs.dirs) == 0) {
		return fmt.Errorf("No sortable files found (%d files checked)", fs.stats.found)
	}
	fs.verbf("scanned targets. found #%d", fs.stats.found)
	return nil
}

func (fs *fsSort) sortAllFiles() error {
	fs.verbf("sorting files...")
	//perform sort
	if fs.DryRun {
		log.Println(color.CyanString("[Dryrun]"))
	}
	//sort concurrency-many files at a time,
	//wait for all to complete and show errors
	queue := make(chan bool, fs.Concurrency)
	wg := &sync.WaitGroup{}
	sortFile := func(file *fileSort) {
		if err := fs.sortFile(file); err != nil {
			log.Printf("[#%d/%d] %s\n  └─> %s\n", file.id, len(fs.sorts), color.RedString(file.path), err)
		}
		<-queue
		wg.Done()
	}
	for _, file := range fs.sorts {
		wg.Add(1)
		queue <- true
		go sortFile(file)
	}
	wg.Wait()
	return nil
}

func (fs *fsSort) watch() error {
	if len(fs.dirs) == 0 {
		return errors.New("No directories to watch")
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Failed to create file watcher: %s", err)
	}
	for dir := range fs.dirs {
		if err := watcher.Add(dir); err != nil {
			return fmt.Errorf("Failed to watch directory: %s", err)
		}
		log.Printf("Watching %s for changes...", color.CyanString(dir))
	}
	select {
	case <-watcher.Events:
	case err := <-watcher.Errors:
		fs.verbf("watch error detected: %s", err)
	}
	go watcher.Close()
	log.Printf("Change detected, re-sorting in %s...", fs.WatchDelay)
	time.Sleep(fs.WatchDelay)
	return nil
}

func (fs *fsSort) add(path string, info os.FileInfo) error {
	//skip hidden files and directories
	if fs.SkipHidden && strings.HasPrefix(info.Name(), ".") {
		fs.verbf("skip hidden file: %s", path)
		return nil
	}
	//limit recursion depth
	if len(fs.sorts) >= fs.FileLimit {
		fs.verbf("skip file: %s. surpassed file limit: %d", path, fs.FileLimit)
		return nil
	}
	//add regular files (non-symlinks)
	if info.Mode().IsRegular() {
		fs.stats.found++
		//skip unmatched file types
		if !fs.validExts[filepath.Ext(path)] {
			fs.verbf("skip unmatched file ext: %s", path)
			return nil
		}
		//skip small files
		if info.Size() < int64(fs.MinFileSize) {
			fs.verbf("skip small file: %s", path)
			return nil
		}
		fs.sorts[path] = &fileSort{id: len(fs.sorts) + 1, path: path, info: info}
		fs.stats.matched++
		return nil
	}
	//recurse into directories
	if info.IsDir() {
		if !fs.Recursive {
			return errors.New("Recursive mode (-r) is required to sort directories")
		}
		//note directory
		fs.dirs[path] = true
		//add all files in dir
		infos, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, info := range infos {
			p := filepath.Join(path, info.Name())
			//recurse
			if err := fs.add(p, info); err != nil {
				return err
			}
		}
	}
	fs.verbf("skip non-regular file: %s", path)
	//skip links,pipes,etc
	return nil
}

func (fs *fsSort) sortFile(file *fileSort) error {
	result, err := SortThreshold(file.path, fs.AccuracyThreshold)
	if err != nil {
		return err
	}
	newPath, err := result.PrettyPath(fs.PathConfig)
	if err != nil {
		return err
	}
	baseDir := ""
	switch mediasearch.MediaType(result.MType) {
	case mediasearch.Series:
		baseDir = fs.TVDir
	case mediasearch.Movie:
		baseDir = fs.MovieDir
	default:
		return fmt.Errorf("Invalid result type: %s", result.MType)
	}
	newPath = filepath.Join(baseDir, newPath)
	//check for subs.srt file
	pathSubs := strings.TrimSuffix(result.Path, filepath.Ext(result.Path)) + ".srt"
	_, err = os.Stat(pathSubs)
	hasSubs := err == nil
	subsExt := ""
	if hasSubs {
		subsExt = "," + color.GreenString("srt")
	}
	//found sort path
	log.Printf("[#%d/%d] %s\n  └─> %s", file.id, len(fs.sorts), color.GreenString(result.Path)+subsExt, color.GreenString(newPath)+subsExt)
	if fs.DryRun {
		return nil //don't actually move
	}
	if result.Path == newPath {
		return nil //already sorted
	}
	//check already exists
	if newInfo, err := os.Stat(newPath); err == nil {
		fileIsLarger := file.info.Size() > newInfo.Size()
		overwrite := fs.Overwrite || (fs.OverwriteIfLarger && fileIsLarger)
		//check if it the same file
		if !os.SameFile(file.info, newInfo) {
			if !overwrite {
				return fmt.Errorf("File already exists '%s' (try setting --overwrite)", newPath)
			}
		} else {
			return nil // File are the same
		}
	}

	//mkdir -p
	err = os.MkdirAll(filepath.Dir(newPath), 0755)
	if err != nil {
		return err //failed to mkdir
	}
	//mv or hardlink
	err = move(fs.Hardlink, result.Path, newPath)
	if err != nil {
		return err //failed to move
	}
	//if .srt file exists for the file, mv it too
	if hasSubs {
		newPathSubs := strings.TrimSuffix(newPath, filepath.Ext(newPath)) + ".srt"
		move(fs.Hardlink, pathSubs, newPathSubs) //best-effort
	}
	return nil
}

func (fs *fsSort) verbf(f string, args ...interface{}) {
	if fs.Verbose {
		log.Printf(f, args...)
	}
}

func move(hard bool, src, dst string) (err error) {
	if hard {
		err = os.Link(src, dst)
	} else {
		err = os.Rename(src, dst)
		//cross-device? shell out to mv
		if err != nil && strings.Contains(err.Error(), "cross-device") && canSysMove {
			err = sysMove(src, dst)
		}
	}
	return
}
