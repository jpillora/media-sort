# media-sort

[![GoDoc](https://godoc.org/github.com/jpillora/media-sort?status.svg)](https://godoc.org/github.com/jpillora/media-sort)

A command-line tool which categorizes provided files and directories by moving them into to a structured directory tree, using various live sources.

### Install

**Binaries**

See [the latest release](https://github.com/jpillora/media-sort/releases/latest) and one-line downloader `curl i.jpillora.com/media-sort | bash`

**Source**

``` sh
$ go get -v github.com/jpillora/media-sort
```

### Features

* Cross platform single binary
* No dependencies
* Easily create [Plex](https://plex.tv) compatible directory structure
* Integration with uTorrent and qbittorrent "Run on Completion" option

### Quick use

``` sh
$ curl i.jpillora.com/media-sort | bash
Downloading media-sort...
Latest version is 2.0.0
######################################### 100.0%
$ ./media-sort --dryrun --recursive .
```

Optionally move into `$PATH`

```
mv media-sort /usr/local/bin/
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
  --recursive, -r        also search through subdirectories
  --dry-run, -d          perform sort but don't actually move any files
  --skip-hidden, -s      skip dot files
  --overwrite, -o        overwrites duplicates
  --overwrite-if-larger  overwrites duplicates if the new file is larger
  --watch, -w            watch the specified directories for changes and
                         re-sort on change
  --watch-delay          delay before next sort after a change (default 3s)
  --help, -h
  --version, -v

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

#### MIT License

Copyright © 2015 Jaime Pillora &lt;dev@jpillora.com&gt;

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
