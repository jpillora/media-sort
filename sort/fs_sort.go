package mediasort

import (
	"errors"
	"fmt"
	"io"
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
	NumDirs           int           `opts:"help=number of directories to include in search (default 0 where -1 means all dirs)"`
	AccuracyThreshold int           `opts:"help=filename match accuracy threshold" default:"is 95, perfect match is 100"`
	MinFileSize       sizestr.Bytes `opts:"help=minimum file size"`
	Recursive         bool          `opts:"help=also search through subdirectories"`
	DryRun            bool          `opts:"help=perform sort but don't actually move any files"`
	SkipHidden        bool          `opts:"help=skip dot files"`
	SkipSubs          bool          `opts:"help=skip subtitles (srt files)"`
	Action            Action        `opts:"help=filesystem action used to sort files (copy|link|move)"`
	HardLink          bool          `opts:"help=use hardlinks instead of symlinks (forces --action link)"`
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
	linkType linkType
}

type fileSort struct {
	id     int
	path   string
	info   os.FileInfo
	result *Result
	err    error
}

// Action used to sort files
type Action string

const (
	// MoveAction sorts by moving
	MoveAction Action = "move"
	// LinkAction sorts by linking
	LinkAction Action = "link"
	// CopyAction sorts by copying
	CopyAction Action = "copy"
)

type linkType string

const (
	hardLink linkType = "hardLink"
	symLink  linkType = "symLink"
)

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
	if c.Action == LinkAction && c.Overwrite {
		return errors.New("Link is already specified, Overwrite won't do anything")
	}
	switch c.Action {
	case MoveAction, LinkAction, CopyAction:
		break
	default:
		return errors.New("Provided action is not available")
	}
	//init fs sort
	fs := &fsSort{
		Config:    c,
		validExts: map[string]bool{},
		linkType:  symLink,
	}
	if c.HardLink {
		fs.Action = LinkAction
		fs.linkType = hardLink
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
	//skip "hidden" files and directories
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
	result, err := SortDepthThreshold(file.path, fs.NumDirs, fs.AccuracyThreshold)
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
	hasSubs := false
	subsExt := ""
	pathSubs := strings.TrimSuffix(result.Path, filepath.Ext(result.Path)) + ".srt"
	if fs.SkipSubs == false {
		_, err = os.Stat(pathSubs)
		hasSubs = err == nil
		if hasSubs {
			subsExt = "," + color.GreenString("srt")
		}
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
	// mkdir -p
	err = os.MkdirAll(filepath.Dir(newPath), 0755)
	if err != nil {
		return err //failed to mkdir
	}
	// action the file
	err = fs.action(result.Path, newPath)
	if err != nil {
		return err //failed to move
	}
	//if .srt file exists for the file, action it too
	if hasSubs {
		newPathSubs := strings.TrimSuffix(newPath, filepath.Ext(newPath)) + ".srt"
		fs.action(pathSubs, newPathSubs) //best-effort
	}
	return nil
}

func (fs *fsSort) verbf(f string, args ...interface{}) {
	if fs.Verbose {
		log.Printf(f, args...)
	}
}

func (fs *fsSort) action(src, dst string) error {
	switch fs.Action {
	case MoveAction:
		return move(src, dst)
	case CopyAction:
		return copy(src, dst)
	case LinkAction:
		return link(src, dst, fs.linkType)
	}
	return errors.New("unknown action")
}

func move(src, dst string) error {
	err := os.Rename(src, dst)
	// cross device move
	if err != nil && strings.Contains(err.Error(), "cross-device") {
		if err := copy(src, dst); err != nil {
			return err
		}
		if err := os.Remove(src); err != nil {
			return err
		}
	}
	return nil
}

func copy(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return nil
}

func link(src, dst string, linkType linkType) error {
	switch linkType {
	case hardLink:
		return os.Link(src, dst)
	case symLink:
		return os.Symlink(src, dst)
	}
	panic("wrong link type, please open an issue")
}
