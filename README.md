
# media-sort

A command-line tool which categorizes provided files and directories by moving them into to a structured directory tree, using the [Open Movie Database API](http://www.omdbapi.com/).

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
* Plex compatible directory structure
* Integration with uTorrent and qbittorrent "Run on Completion" option

### Quick use

``` sh
$ curl i.jpillora.com/media-sort | sh
Downloading: media-sort_1.1.0_darwin_amd64
######################################### 100.0%
$ ./media-sort --dryrun .
```

Optionally move into `$PATH`

```
mv media-sort /usr/local/bin/
```

### Usage

```
$ media-sort --help
```

<tmpl,code: go run main.go --help>
```

	Usage: mediasort [options] [file/directory]

	MediaSort performs a simple categorization on the provided file or directory.

	Movies are moved to:
		<movie-dir>/<title> (<year>).ext
	TV Shows are moved to:
		<tv-dir>/<title>/<title> S<season>E<episode>.ext

	Version: 0.0.0

	Options:
	--movie-dir -m  The destination movie directory (defaults to $HOME/movies).
	--tv-dir -t     The destination TV directory (defaults to $HOME/tv).
	--ext -e        Extensions considered (defaults to "mp4,avi,mkv").
	--dryrun -d     Runs in read-only mode.
	--depth         Directory depth to search for files.
	--version -v    Display version.
	--help -h       This help text.

	Read more:
	  https://github.com/jpillora/mediasort

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