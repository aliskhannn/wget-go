package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/pflag"

	"github.com/aliskhannn/wget-go/internal/config"
	"github.com/aliskhannn/wget-go/internal/crawler"
	"github.com/aliskhannn/wget-go/internal/flags"
	"github.com/aliskhannn/wget-go/internal/parser"
	"github.com/aliskhannn/wget-go/internal/worker"
)

func main() {
	// Initialize flags.
	options := flags.InitFlags()
	pflag.Parse()

	// If --help was requested, show usage and exit
	if *options.Help {
		pflag.Usage()
		os.Exit(1)
	}

	args := pflag.Args()
	if len(args) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "error: url is required")
		pflag.Usage()
		os.Exit(1)
	}

	cfg := config.New(
		*options.Depth,
		*options.MaxConcurrency,
		*options.Timeout,
		*options.Robots,
		*options.Help,
	)

	rawURL := args[0]

	// Normalize a raw URL string into a *url.URL, adding a scheme if missing.
	baseURL, err := parser.NormalizeURL(rawURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: invalid url: %s\n", rawURL)
		os.Exit(1)
	}

	c := crawler.New(cfg.Depth)

	if cfg.Robots {
		if err := c.LoadRobots(baseURL); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to load robots.txt: %v", err)
		}
	}

	wg := new(sync.WaitGroup)
	w := worker.New(c, baseURL, cfg.MaxWorkers, wg)

	w.Start()
}
