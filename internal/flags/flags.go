package flags

import (
	"time"

	"github.com/spf13/pflag"
)

// Flags holds all command-line options for the program.
type Flags struct {
	Depth          *int           // recursion depth
	MaxConcurrency *int           // number of parallel downloads
	Timeout        *time.Duration // timeout on HTTP requests
	Robots         *bool          // whether to respect robots.txt
	Help           *bool          // show help and exit
}

// InitFlags initializes and parses command-line flags.
// It returns a Flags struct with pointers to the parsed values.
func InitFlags() Flags {
	return Flags{
		Depth:          pflag.IntP("depth", "d", 0, "Maximum recursion depth (0 = only the given page)"),
		MaxConcurrency: pflag.IntP("concurrency", "c", 4, "Maximum number of parallel downloads"),
		Timeout:        pflag.DurationP("timeout", "t", 10*time.Second, "Timeout for HTTP requests (e.g. 10s, 2m)"),
		Robots:         pflag.Bool("robots", false, "Respect robots.txt rules"),
		Help:           pflag.BoolP("help", "h", false, "Show this help message"),
	}
}
