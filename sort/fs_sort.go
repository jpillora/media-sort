package mediasort

// //fsSort is a media sorter
// type fsSort struct {
// 	c                       Config
// 	exts                    map[string]bool
// 	files                   map[string]*pathSort //abs-path -> file
// 	Checked, Matched, Moved int
// }
//
// //NewfsSort creates a new fsSort
// func FileSystemSort(c Config) error {
// 	if len(c.Targets) == 0 {
// 		return nil, fmt.Errorf("Please provide at least one file or directory")
// 	}
// 	if c.MovieDir == "" {
// 		c.MovieDir = "."
// 	}
// 	if c.TVDir == "" {
// 		c.TVDir = "."
// 	}
// 	exts := map[string]bool{}
// 	for _, e := range strings.Split(c.Exts, ",") {
// 		exts["."+e] = true
// 	}
//
// 	fs := &fsSort{
// 		c:     c,
// 		exts:  exts,
// 		files: map[string]*pathSort{},
// 	}
//
// 	//convert all targets into file sorter objects
// 	for _, path := range c.Targets {
// 		info, err := os.Stat(path)
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = s.addFile(path, info, 0)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return s, nil
// }
//
// // func (s *fsSort) addFile(path string, info os.FileInfo, depth int) error {
// // }
//
// func (s *fsSort) addFile(path string, info os.FileInfo, depth int) error {
//
// 	//skip hidden files and directories
// 	if strings.HasPrefix(info.Name(), ".") {
// 		return nil
// 	}
//
// 	//limit recursion depth
// 	if depth > s.c.Depth {
// 		return nil
// 	}
//
// 	if info.Mode().IsRegular() {
// 		s.Checked++
// 		//add single file
// 		f, err := newFilefsSort(s, path, info)
// 		if f == nil || err != nil {
// 			return err
// 		}
// 		s.files[path] = f
// 		s.Matched++
// 		return nil
// 	}
//
// 	if info.IsDir() {
//
// 		//TODO watch directory
// 		// if s.c.Watch {
// 		// }
//
// 		//add all files in dir
// 		infos, err := ioutil.ReadDir(path)
// 		if err != nil {
// 			return err
// 		}
// 		for _, info := range infos {
// 			p := filepath.Join(path, info.Name())
// 			//recurse
// 			err := s.addFile(p, info, depth+1)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
//
// 	//skip links,pipes,etc
// 	return nil
// }
//
// func (s *fsSort) sortFiles() []error {
// 	var errors []error
// 	//parallel sort all files
// 	wg := &sync.WaitGroup{}
// 	for _, f := range s.files {
// 		wg.Add(1)
// 		go f.goRun(wg)
// 	}
// 	wg.Wait()
//
// 	//collect errors
// 	for _, f := range s.files {
// 		if f.err != nil {
// 			err := fmt.Errorf("%s: %s", f.name, f.err)
// 			errors = append(errors, err)
// 		} else {
// 			s.Moved++
// 		}
// 	}
//
// 	return errors
// }
//
// func (s *fsSort) Run() []error {
//
// 	if s.c.DryRun {
// 		log.Println(color.CyanString("[Dryrun Mode]"))
// 	}
//
// 	//just run once
// 	if !s.c.Watch {
// 		return s.sortFiles()
// 	}
//
// 	//run watcher
// 	// for {
// 	// 	s.sortFiles()
// 	// }
//
// 	return []error{fmt.Errorf("Watch not implemented")}
// }

// //DEBUG
// // log.Printf("SUCCESS = D%d #%d\n  %s\n  %s", r.Distance, len(query), query, r.Title)
// log.Printf("Moving\n  %s\n  └─> %s", ps.path, color.GreenString(dest))
//
// if ps.s.c.DryRun {
//     return nil
// }
//
// //check already exists
// if _, err := os.Stat(dest); err == nil {
//     return fmt.Errorf("File already exists '%s'", dest)
// }
//
// //mkdir -p
// err = os.MkdirAll(filepath.Dir(dest), 0755)
// if err != nil {
//     return err //failed to mkdir
// }
//
// //mv
// err = os.Rename(ps.path, dest)
// if err != nil {
//     return err //failed to move
// }
