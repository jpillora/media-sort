package mediasort

import (
	"testing"

	mediasearch "github.com/jpillora/media-sort/search"
)

func TestAudiobookPrettyPath(t *testing.T) {
	config := PathConfig{
		AudiobookTemplate: DefaultAudiobookTemplate,
	}

	for _, tc := range []struct {
		name     string
		result   Result
		expected string
	}{
		{
			"simple audiobook",
			Result{
				Name:   "Project Hail Mary",
				Author: "Andy Weir",
				Ext:    "m4b",
				MType:  string(mediasearch.Audiobook),
			},
			"Andy Weir/Project Hail Mary.m4b",
		},
		{
			"audiobook with narrator",
			Result{
				Name:     "Dune",
				Author:   "Frank Herbert",
				Narrator: "Scott Brick",
				Ext:      "mp3",
				MType:    string(mediasearch.Audiobook),
			},
			"Frank Herbert/Dune (Narrated by Scott Brick).mp3",
		},
		{
			"audiobook in series",
			Result{
				Name:       "Gardens of the Moon",
				Author:     "Steven Erikson",
				Series:     "Malazan Book of the Fallen",
				SeriesPart: 1,
				Ext:        "m4b",
				MType:      string(mediasearch.Audiobook),
			},
			"Steven Erikson/Malazan Book of the Fallen - 01 - Gardens of the Moon.m4b",
		},
		{
			"complete audiobook metadata",
			Result{
				Name:       "The Way of Kings",
				Author:     "Brandon Sanderson",
				Series:     "The Stormlight Archive",
				SeriesPart: 1,
				Narrator:   "Michael Kramer",
				Ext:        "m4b",
				MType:      string(mediasearch.Audiobook),
			},
			"Brandon Sanderson/The Stormlight Archive - 01 - The Way of Kings (Narrated by Michael Kramer).m4b",
		},
		{
			"various authors collection",
			Result{
				Name:   "Best Science Fiction Stories of 2024",
				Author: "Various Authors",
				Ext:    "mp3",
				MType:  string(mediasearch.Audiobook),
			},
			"Various Authors/Best Science Fiction Stories of 2024.mp3",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.result.PrettyPath(config)
			if err != nil {
				t.Fatalf("PrettyPath failed: %v", err)
			}
			if got != tc.expected {
				t.Errorf("PrettyPath = %q, want %q", got, tc.expected)
			}
		})
	}
}
