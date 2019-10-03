# media-sort

[![GoDoc](https://godoc.org/github.com/jpillora/media-sort?status.svg)](https://godoc.org/github.com/jpillora/media-sort) [![CI](https://github.com/jpillora/media-sort/workflows/CI/badge.svg)](https://github.com/jpillora/media-sort/actions?workflow=CI)

A command-line tool and Go (golang) library which categorizes provided files and directories by moving them into to a structured directory tree, using various live sources.

### Install

**Binaries**

[![Releases](https://img.shields.io/github/release/jpillora/media-sort.svg)](https://github.com/jpillora/media-sort/releases) [![Releases](https://img.shields.io/github/downloads/jpillora/media-sort/total.svg)](https://github.com/jpillora/media-sort/releases)

See [the latest release](https://github.com/jpillora/media-sort/releases/latest) or download and install it now with `curl https://i.jpillora.com/media-sort! | bash`

**Source**

``` sh
$ go get -v github.com/jpillora/media-sort
```

### Features

* Cross platform single binary
* No dependencies
* Easily create a [Plex](https://plex.tv)-compatible directory structure
* Integration with uTorrent and qbittorrent "Run on Completion" option

### Quick use

``` sh
$ curl https://i.jpillora.com/media-sort | bash
Downloading media-sort...
Latest version is 2.X.X
######################################### 100.0%
$ ./media-sort --dryrun --recursive .
```

Optionally move into `$PATH`

```
$ mv media-sort /usr/local/bin/
```

Test run `media-sort`

```
$ cd my-media/
$ media-sort --dry-run --recursive .
2016/01/30 09:35:47 [Dryrun]
2016/01/30 09:35:47 Searching dick van dyke show (series)
2016/01/30 09:35:47 [#1/1] dick-van-dyke-show.s01e10.[Awesome-Audio]-[Super-Quality]-[Name-of-Encoder].mp4
  └─> The Dick Van Dyke Show S01E10.mp4
```

### CLI Usage

```
$ media-sort --help
```

<tmpl,code: go run main.go --help>
``` plain

  Usage: media-sort [options] targets...

  media-sort categorizes the provided files and directories (targets) by
  moving them into to a structured directory tree, sorting is currently
  performed using TVMaze, OMDB and Google.

  Options:
  --tv-dir, -t           tv series base directory (defaults to current
                         directory)
  --movie-dir, -m        movie base directory (defaults to current directory)
  --tv-template          tv series path template
  --movie-template       movie path template
  --extensions, -e       types of files that should be sorted (default
                         mp4,avi,mkv)
  --concurrency, -c      search concurrency [warning] setting this too high
                         can cause rate-limiting errors (default 6)
  --file-limit, -f       maximum number of files to search (default 1000)
  --min-file-size        minimum file size (default 25MB)
  --recursive, -r        also search through subdirectories
  --dry-run, -d          perform sort but don't actually move any files
  --skip-hidden, -s      skip dot files
  --overwrite, -o        overwrites duplicates
  --overwrite-if-larger  overwrites duplicates if the new file is larger
  --watch, -w            watch the specified directories for changes and
                         re-sort on change
  --watch-delay          delay before next sort after a change (default 3s)
  --verbose, -v          verbose logs
  --help, -h
  --version

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

  Version:
    0.0.0-src

  Read more:
    github.com/jpillora/media-sort

```
</tmpl>

#### Programmatic Use

See https://godoc.org/github.com/jpillora/media-sort

The API has 3 layers:

1. An explicit search: `mediasearch.Search(query, year, mediatype string) (mediasearch.Result, error)`
    Returns search result
2. A path string correction (using `Search`): `mediasort.Sort(path string) (*mediasort.Result, error)`
    Attempts to extract search query information from the path string, returns result which can be used to format a new path or `result.PrettyPath()` can be used.
3. A filesystem correction (using `Sort`): `mediasort.FileSystemSort(config mediasort.Config) error`
    Attempts to sort all paths provided in `config.Targets`, when successful - results are formatted and renamed to use the newly formatted path.

#### MIT License

Copyright © 2016 Jaime Pillora &lt;dev@jpillora.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
