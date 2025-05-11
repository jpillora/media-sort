package mediasearch

// MediaType can be Movie or Series
type MediaType string

const (
	//Series represents a tv series
	Series MediaType = "series"
	//Movie represents a movie
	Movie MediaType = "movie"
	//Audiobook represents an audiobook
	Audiobook MediaType = "audiobook"
)

// Result is a single search result
type Result struct {
	Title      string
	Year       string
	Type       MediaType
	IsDupe     bool
	Accuracy   int
	Author     string // Author of the audiobook
	Narrator   string // Narrator of the audiobook
	Series     string // Series name for the audiobook
	SeriesPart int    // Position in the series
}

func (r Result) String() string {
	return r.Title + " (" + r.Year + ")"
}
