package crawler

import (
	"errors"
	"fmt"
	"github.com/aliskhannn/wget-go/internal/fetcher"
	"github.com/aliskhannn/wget-go/internal/files"
	"github.com/temoto/robotstxt"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Crawler represents a web crawler with an HTTP client, crawling depth,
// a map of visited URLs, and optional robots.txt rules.
type Crawler struct {
	Client    *http.Client
	Depth     int
	Visited   sync.Map
	RobotsMap *robotstxt.RobotsData
}

// New creates a new Crawler instance with the specified depth.
// The HTTP client uses a default timeout of 5 seconds.
func New(depth int) *Crawler {
	return &Crawler{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Depth:   depth,
		Visited: sync.Map{},
	}
}

// Fetch downloads the content from the given URL using the crawler's HTTP client.
// It retries the request up to the specified number of times in case of a timeout.
// Returns the response body as a byte slice, or an error if the request fails permanently.
func (c *Crawler) Fetch(u *url.URL, retries int) ([]byte, error) {
	var data []byte
	var err error

	for i := 0; i <= retries; i++ {
		// Attempt to download data using the fetcher.
		data, err = fetcher.Fetch(u, c.Client)
		if err == nil {
			// Successful download, return data.
			return data, nil
		}
		// Handle timeout errors separately with a retry.
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("timeout fetching %s (attempt %d/%d)", u, i+1, retries+1)
			time.Sleep(time.Second * 2) // small pause before retrying
			continue
		}

		// For other errors, do not retry and break the loop.
		break
	}

	return nil, err
}

// Save writes the downloaded data to a file corresponding to the URL.
// Returns an error if saving fails.
func (c *Crawler) Save(u *url.URL, data []byte) error {
	// Save downloaded data to a file.
	err := files.SaveFile(u, data)
	if err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

// LoadRobots fetches and parses robots.txt for the given base URL.
// It populates the Crawler's RobotsMap field with the parsed rules.
// Returns an error if fetching or parsing fails.
func (c *Crawler) LoadRobots(baseURL *url.URL) error {
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", baseURL.Scheme, baseURL.Host)

	// Parse robots.txt URL.
	parsedURL, err := url.Parse(robotsURL)
	if err != nil {
		return err
	}

	// Fetch robots.txt content.
	data, err := fetcher.Fetch(parsedURL, c.Client)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	// Parse robots.txt rules.
	c.RobotsMap, err = robotstxt.FromBytes(data)
	if err != nil {
		return err
	}

	return nil
}
