# TODO: Add Audiobook Sorting Support

This document outlines the steps required to add audiobook sorting functionality to the `media-sort` tool.

## 1. Define Audiobook Media Type
-   **File:** `search/result.go`
    -   [x] Add `Audiobook MediaType = "audiobook"` to the `MediaType` constants.
-   **Throughout codebase:**
    -   [x] Review and update all logic that switches or checks `MediaType` (e.g., in `search/search.go`, `search/get_imdb.go`, `sort/path_sort.go`) to correctly handle the new `Audiobook` type.

## 2. Configuration Updates
-   **File:** `sort/fs_sort.go` (Config struct)
    -   [x] Add `AudiobookDir string \`opts:"help=audiobook base directory (defaults to current directory)"\`` field for specifying the output directory for audiobooks.
-   **File:** `sort/path_sort.go` (PathConfig struct and defaults)
    -   [x] Add `AudiobookTemplate string \`help:"audiobook path template"\`` to `PathConfig`.
    -   [x] Define `DefaultAudiobookTemplate`. A possible example:
        `{{ .Author }}/{{ if .Series }}{{ .Series }} - {{ printf "%02d" .SeriesPart }} - {{end}}{{ .Title }}{{ if .Narrator }} (Narrated by {{ .Narrator }}){{end}}.{{ .Ext }}`
        (Consider which fields like Author, Series, SeriesPart, Narrator will be available in the `Result` struct).
-   **File:** `main.go`
    -   [x] Add new command-line flags `--audiobook-dir` and `--audiobook-template` linked to the new `Config` fields.
    -   [x] Update the `info` and `pathTemplates` help messages to include audiobooks.

## 3. Search Functionality
-   **Research:**
    -   [x] Identify and evaluate potential metadata sources/APIs for audiobooks (e.g., LibriVox, Audible (if API access is feasible), Google Books API, Open Library, etc.).
-   **File(s):** `search/` directory (e.g., `search_googlebooks.go`, `search_librivox.go`)
    -   [x] Implement new search functions for one or more chosen audiobook metadata providers. These functions should be similar in signature to `searchTVMaze` or `searchMovieDB` (e.g., `func searchAudiobookSource(query, year string, mediatype MediaType) ([]Result, error)`).
    -   [x] Ensure these functions populate the `Result` struct with audiobook-specific data (Author, Narrator, Series info, etc.).
-   **File:** `search/search.go`
    -   [x] Create an `audiobookSearches` slice of type `[]search` and add the new audiobook search functions to it.
    -   [x] Modify the `Search` and `SearchThreshold` functions:
        -   When `mediatype` (or `mt`) is `Audiobook`, iterate over `audiobookSearches` instead of `tvSearches` or `movieSearches`.

## 4. Path Parsing and Information Extraction
-   **File:** `sort/path_sort.go` (Result struct)
    -   [x] Add new fields to the `Result` struct to store audiobook-specific metadata. For example:
        -   `Author string`
        -   `Narrator string`
        -   `BookSeries string` (if it's part of a series)
        -   `SeriesPart int` (e.g., book number in series)
-   **File(s):** `sort/strings.go` (and potentially `search/strings.go` if normalization needs to be audiobook-aware)
    -   [x] Develop and add new regular expressions to parse audiobook-specific details from filenames/paths (e.g., "Author - Series Book 1 - Title", "Title by Author Narrated by Narrator").
    -   [x] Consider common audiobook naming conventions and file structures (e.g., often in subfolders by Author/Series).
-   **File:** `sort/path_sort.go` (path parsing logic, likely within `SortDepthThreshold` or a helper function it calls)
    -   [x] Update the path parsing logic to attempt to identify files as audiobooks.
    -   [x] Extract audiobook-specific information (Author, Title, Series, Part, Narrator) using the new regexes.
    -   [x] Set `r.MType = string(mediasearch.Audiobook)` for identified audiobooks.
    -   [x] Populate the new audiobook-specific fields in the `Result` struct.

## 5. Sorting Logic
-   **File:** `sort/fs_sort.go` (primarily `FileSystemSort` and related functions)
    -   [x] Modify the main sorting logic to handle the `Audiobook` `MType`.
    -   [x] When a result is an audiobook, use `c.AudiobookDir` (or current dir if not set) as the base output directory.
    -   [x] Use `c.AudiobookTemplate` (or `DefaultAudiobookTemplate`) for formatting the final path.
-   **File:** `sort/path_sort.go` (`PrettyPath` function)
    -   [x] Ensure `PrettyPath` can correctly apply the `AudiobookTemplate` using the audiobook-specific fields in the `Result` struct.

## 6. Testing
-   **File:** `sort/path_sort_test.go`
    -   [x] Add comprehensive test cases for parsing various audiobook filenames and directory structures.
        -   Include examples with author, title, series, book number, narrator.
        -   Test different delimiters and orders of information.
    -   [x] Add tests to verify the correct application of the `AudiobookTemplate`.
-   **New Files (Optional):** `search/search_googlebooks_test.go`, etc.
    -   [x] If feasible, add unit tests for the new search functions, potentially mocking API responses.

## 7. Documentation
-   **File:** `README.md`
    -   [x] Add a section describing the new audiobook sorting feature.
    -   [x] Document the new CLI flags (`--audiobook-dir`, `--audiobook-template`).
    -   [x] Provide an example of the default audiobook path structure.
-   **File:** `main.go`
    -   [x] Update the `info` constant (general description) to mention audiobook support.
    -   [x] Update the `pathTemplates` constant to show the default audiobook template and relevant variables.

## 8. Refinements & Edge Cases
-   [ ] Consider how to handle multi-file audiobooks (e.g., one file per chapter). The current model seems to be one output file per input file. This might need significant thought if a "grouping" concept is desired.
-   [ ] How to handle "Various Authors" or "Collections" for audiobooks.
-   [ ] Disambiguation if a title could be a movie, TV show, or audiobook. (The current system might rely on file naming conventions or might need explicit type hints if it becomes an issue).
-   [ ] Update any file extension logic (`c.Extensions`) if audiobooks commonly use extensions not already covered (e.g., `.m4b`, `.aax`). The current default `mp4,m4v,avi,mkv,mpeg,mpg,mov,webm` might miss common audiobook formats like `.mp3`, `.aac`, `.flac`. This should be configurable or have sensible defaults.
