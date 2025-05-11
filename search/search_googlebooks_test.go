package mediasearch

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchGoogleBooks(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test the URL construction
		if r.URL.Path != "/books/v1/volumes" {
			t.Errorf("Expected path /books/v1/volumes, got %s", r.URL.Path)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"items": [
				{
					"volumeInfo": {
						"title": "Project Hail Mary",
						"authors": ["Andy Weir"],
						"publishedDate": "2021-05-04",
						"readingModes": {
							"audio": true
						},
						"seriesInfo": {
							"title": "",
							"bookNumber": 0
						}
					}
				}
			]
		}`))
	}))
	defer server.Close()

	// Store the real URL and restore it after the test
	realURL := "https://www.googleapis.com/books/v1/volumes"
	searchURL = server.URL + "/books/v1/volumes"
	defer func() { searchURL = realURL }()

	// Test cases
	tests := []struct {
		name      string
		query     string
		year      string
		mediatype MediaType
		want      Result
		wantErr   bool
	}{
		{
			name:      "basic audiobook search",
			query:     "Project Hail Mary",
			year:      "2021",
			mediatype: Audiobook,
			want: Result{
				Title:  "Project Hail Mary",
				Year:   "2021",
				Type:   Audiobook,
				Author: "Andy Weir",
			},
			wantErr: false,
		},
		{
			name:      "non-audiobook media type",
			query:     "test",
			mediatype: Movie,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchGoogleBooks(tt.query, tt.year, tt.mediatype)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchGoogleBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) == 0 {
					t.Fatal("expected at least one result")
				}
				result := got[0]
				if result.Title != tt.want.Title {
					t.Errorf("Title = %v, want %v", result.Title, tt.want.Title)
				}
				if result.Author != tt.want.Author {
					t.Errorf("Author = %v, want %v", result.Author, tt.want.Author)
				}
				if result.Type != tt.want.Type {
					t.Errorf("Type = %v, want %v", result.Type, tt.want.Type)
				}
			}
		})
	}
}
