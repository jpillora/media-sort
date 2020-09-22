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
$ curl https://i.jpillora.com/media-sort! | bash
Installing jpillora/media-sort v2.4.3.....
######################################################################## 100.0%
Installed at /usr/local/bin/media-sort
```

Test run `media-sort` (read-only mode)

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

``` plain

  Usage: media-sort [options] <target> [target] ...

  media-sort categorizes the provided files and directories (targets) by
  moving them into to a structured directory tree, sorting is currently
  performed using TVMaze, MovieDB and Google.

  Options:
  --tv-dir, -t              tv series base directory (defaults to current directory)
  --movie-dir, -m           movie base directory (defaults to current directory)
  --tv-template             tv series path template
  --movie-template          movie path template
  --extensions, -e          types of files that should be sorted (default mp4,m4v,avi,mkv,mpeg,mpg,mov,webm)
  --concurrency, -c         search concurrency [warning] setting this too high can cause rate-limiting errors (default 6)
  --file-limit, -f          maximum number of files to search (default 1000)
  --num-dirs, -n            number of directories to include in search (default 0 where -1 means all dirs)
  --accuracy-threshold, -a  filename match accuracy threshold (default 95)
  --min-file-size           minimum file size (default 25MB)
  --recursive, -r           also search through subdirectories
  --dry-run, -d             perform sort but don't actually move any files
  --skip-hidden, -s         skip dot files
  --skip-subs               skip subtitles (srt files)
  --action                  filesystem action used to sort files (copy|link|move, default move)
  --hard-link, -h           use hardlinks instead of symlinks (forces --action link)
  --overwrite, -o           overwrites duplicates
  --overwrite-if-larger     overwrites duplicates if the new file is larger
  --watch, -w               watch the specified directories for changes and re-sort on change
  --watch-delay             delay before next sort after a change (default 3s)
  --verbose, -v             verbose logs
  --version                 display version
  --help                    display help

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
    X.Y.Z

  Read more:
    github.com/jpillora/media-sort

```

#### Programmatic Use

See https://godoc.org/github.com/jpillora/media-sort

The API has 3 layers:

1. An explicit search: `mediasearch.Search(query, year, mediatype string) (mediasearch.Result, error)`
    Returns search result
2. A path string correction (using `Search`): `mediasort.Sort(path string) (*mediasort.Result, error)`
    Attempts to extract search query information from the path string, returns result which can be used to format a new path or `result.PrettyPath()` can be used.
3. A filesystem correction (using `Sort`): `mediasort.FileSystemSort(config mediasort.Config) error`
    Attempts to sort all paths provided in `config.Targets`, when successful - results are formatted and renamed to use the newly formatted path.
