package mediasearch

type MediaType string

const (
	Series MediaType = "series"
	Movie  MediaType = "movie"
)

type Result struct {
	Title   string
	Year    string
	Type    MediaType
	IsDupe  bool
	strdist int
}

func (r Result) String() string {
	return r.Title + " ("+r.Year+")"
}
