package config

import "time"

// Config represents the configuration options for the crawler.
type Config struct {
	Depth      int           // maximum recursion depth
	MaxWorkers int           // the number of concurrent workers
	Timeout    time.Duration // HTTP request timeout duration
	Robots     bool          // whether to respect robots.txt rules
	ShowHelp   bool          // whether to display the help message
}

// New creates and returns a new Config instance with the given parameters.
func New(depth int, maxWorkers int, timeout time.Duration, robots bool, showHelp bool) Config {
	return Config{
		Depth:      depth,
		MaxWorkers: maxWorkers,
		Timeout:    timeout,
		Robots:     robots,
		ShowHelp:   showHelp,
	}
}
