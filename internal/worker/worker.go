package worker

import (
	"bytes"
	"log"
	"net/url"
	"sync"

	"github.com/aliskhannn/wget-go/internal/crawler"
	"github.com/aliskhannn/wget-go/internal/parser"
)

// Job represents a single crawling task with its URL and current recursion depth.
type Job struct {
	URL   *url.URL
	Depth int
}

// Worker manages a pool of goroutines to process crawling jobs concurrently.
type Worker struct {
	Crawler    *crawler.Crawler // crawler instance
	BaseURL    *url.URL         // the starting point of the crawl
	Jobs       chan *Job        // the channel of Job tasks
	MaxWorkers int              // the number of concurrent worker goroutines
	wg         *sync.WaitGroup
}

// New creates a new Worker instance with the given crawler, base URL,
// maximum number of workers, and WaitGroup for synchronization.
func New(crawler *crawler.Crawler, baseURL *url.URL, maxWorkers int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		Crawler:    crawler,
		BaseURL:    baseURL,
		Jobs:       make(chan *Job, maxWorkers*2), // buffered channel
		MaxWorkers: maxWorkers,
		wg:         wg,
	}
}

// Start launches the worker pool and begins processing jobs from.
// It waits until all workers have finished.
func (w *Worker) Start() {
	// Start MaxWorkers worker goroutines.
	for i := 0; i < w.MaxWorkers; i++ {
		w.wg.Add(1)
		go w.worker()
	}

	// Start the goroutine that seeds the initial URLs.
	w.wg.Add(1)
	go w.processAllURL()

	w.wg.Wait() // wait for all goroutines to finish
}

// worker continuously processes jobs from the Jobs channel until it is closed.
func (w *Worker) worker() {
	defer w.wg.Done()

	for job := range w.Jobs {
		w.processJob(job)
	}
}

// processJob handles a single Job: checks robots.txt, fetches the URL,
// parses HTML links if applicable, and saves the data locally.
func (w *Worker) processJob(job *Job) {
	// Skip the URL if disallowed by robots.txt.
	if w.Crawler.RobotsMap != nil {
		uaGroup := w.Crawler.RobotsMap.FindGroup("*") // applies to all user-agents
		if !uaGroup.Test(job.URL.Path) {
			log.Printf("robots.txt disallow: skipping %s", job.URL)
			return
		}
	}

	// Fetch the URL content.
	data, err := w.Crawler.Fetch(job.URL, 2)
	if err != nil {
		log.Printf("fetch error %s: %v", job.URL, err)
		return
	}

	// Parse HTML content and rewrite local links.
	if parser.IsHTML(job.URL.Path, data) {
		_, data, err = parser.ParseAndRewriteLinks(bytes.NewReader(data), job.URL)
		if err != nil {
			log.Printf("parse error %s: %v", job.URL, err)
			return
		}
	}

	// Save the fetched (and possibly rewritten) data to disk.
	if err := w.Crawler.Save(job.URL, data); err != nil {
		log.Printf("save error %s: %v", job.URL, err)
	}
}

// processAllURL fetches the initial base URL, parses all links, and
// seeds the Jobs channel with crawl tasks for the worker pool.
func (w *Worker) processAllURL() {
	defer w.wg.Done()

	// Fetch the base URL content.
	data, err := w.Crawler.Fetch(w.BaseURL, 2)
	if err != nil {
		close(w.Jobs)
		return
	}

	// Parse all links in the base page recursively up to the maximum depth.
	allLinks, _, err := parser.ParseAllLinks(w.Crawler, data, w.BaseURL, 0)
	if err != nil {
		close(w.Jobs)
		return
	}

	// Send all discovered links to the Jobs channel.
	for _, link := range allLinks {
		w.Jobs <- &Job{
			URL:   link,
			Depth: w.Crawler.Depth,
		}
	}

	// Close the Jobs channel to signal workers that no more jobs will be sent.
	close(w.Jobs)
}
