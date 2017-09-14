package mediasearch

//MediaType can be Movie or Series
type MediaType string

const (
	//Series represents a tv series
	Series MediaType = "series"
	//Movie represents a movie
	Movie MediaType = "movie"
)

// Result is a single search result
type Result struct {
	Title    string
	Year     string
	Type     MediaType
	IsDupe   bool
	Accuracy int
}

func (r Result) String() string {
	return r.Title + " (" + r.Year + ")"
}
