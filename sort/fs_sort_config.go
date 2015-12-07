package mediasort

//Config is a sorter configuration
type Config struct {
	Targets  []string
	TVDir    string
	MovieDir string
	Exts     string
	Depth    int
	DryRun   bool
	Watch    bool
}
