package mediasort

import (
	"log"
	"strings"
	"testing"

	mediasearch "github.com/jpillora/media-sort/search"
)

func TestPathParse(t *testing.T) {
	for _, tc := range []struct {
		Input  string
		Depth  int
		Expect Result
	}{
		{
			"/a/long/path/foo/bar s01e02.mp4",
			0,
			Result{
				Query:   "bar",
				Name:    "bar s01e02",
				Ext:     "mp4",
				MType:   string(mediasearch.Series),
				Episode: 2,
			},
		},
		{
			"/a/long/path/foo/bar s01e02.mp4",
			1,
			Result{
				Query:   "foo bar",
				Name:    "foo bar s01e02",
				Ext:     "mp4",
				MType:   string(mediasearch.Series),
				Episode: 2,
			},
		},
		{
			"/a/path/bazz (2020).mkv",
			0,
			Result{
				Query: "bazz",
				Name:  "bazz (2020)",
				Ext:   "mkv",
				MType: string(mediasearch.Movie),
				Year:  "2020",
			},
		},
		{
			"/another/path/My Cool Show Season 2/S02E05-a-title.mp4",
			1,
			Result{
				Query:   "my cool show",
				Name:    "My Cool Show Season 2 S02E05-a-title",
				Ext:     "mp4",
				MType:   string(mediasearch.Series),
				Season:  2,
				Episode: 5,
			},
		},
		{
			"/another/path/My Cool Show Season 2 720p-aBcD/S02E05-a-title.mp4",
			1,
			Result{
				Query:   "my cool show",
				Name:    "My Cool Show Season 2 720p-aBcD S02E05-a-title",
				Ext:     "mp4",
				MType:   string(mediasearch.Series),
				Season:  2,
				Episode: 5,
			},
		},
		{
			"xyz 2012.mp4",
			42, //42 does nothing here
			Result{
				Query: "xyz",
				Name:  "xyz 2012",
				Ext:   "mp4",
				MType: string(mediasearch.Movie),
				Year:  "2012",
			},
		},
		{
			"/my/movie/xyz 2012.mp4",
			-1, //all dirs
			Result{
				Query: "my movie xyz",
				Name:  "my movie xyz 2012",
				Ext:   "mp4",
				MType: string(mediasearch.Movie),
				Year:  "2012",
			},
		},
		{
			"/audiobooks/Stephen King/The Stand.m4b",
			0,
			Result{
				Query:  "the stand",
				Name:   "The Stand",
				Author: "Stephen King",
				Ext:    "m4b",
				MType:  string(mediasearch.Audiobook),
			},
		},
		{
			"/downloads/The Way of Kings by Brandon Sanderson (Narrated by Michael Kramer).m4b",
			0,
			Result{
				Query:    "the way of kings",
				Name:     "The Way of Kings",
				Author:   "Brandon Sanderson",
				Narrator: "Michael Kramer",
				Ext:      "m4b",
				MType:    string(mediasearch.Audiobook),
			},
		},
		{
			"/books/Terry Pratchett/Discworld/Book 1 - The Color of Magic.mp3",
			1,
			Result{
				Query:      "discworld book 1 the color of magic",
				Name:       "The Color of Magic",
				Author:     "Terry Pratchett",
				Series:     "Discworld",
				SeriesPart: 1,
				Ext:        "mp3",
				MType:      string(mediasearch.Audiobook),
			},
		},
	} {
		//support windows
		path := strings.ReplaceAll(tc.Input, "/", sep)
		//defaults
		exp := tc.Expect
		exp.Path = path
		if exp.Season == 0 {
			exp.Season = 1
		}
		if exp.Episode == 0 {
			exp.Episode = -1
		}
		if exp.ExtraEpisode == 0 {
			exp.ExtraEpisode = -1
		}
		//execute test case
		got, err := runPathParse(path, tc.Depth)
		if err != nil {
			t.Fatal(err)
		}
		if got != exp {
			log.Fatalf("input: %s (depth %d)\ngot: %#v\nexp: %#v",
				tc.Input, tc.Depth, got, exp)
		}
	}
}
