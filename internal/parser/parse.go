package parser

import (
	"bytes"
	"errors"
	"github.com/aliskhannn/wget-go/internal/crawler"
	"golang.org/x/net/html"
	"io"
	"net/url"
	"path/filepath"
	"strings"
)

var (
	ErrMaxDepthExceeded = errors.New("max depth exceeded")
)

// ParseAllLinks recursively parses HTML content, extracts links, and returns
// all local links up to the specified depth. It also returns the potentially
// rewritten HTML data for local usage.
func ParseAllLinks(c *crawler.Crawler, data []byte, baseURL *url.URL, curDepth int) ([]*url.URL, []byte, error) {
	// Stop recursion if the current depth exceeds the crawler's maximum depth.
	if curDepth > c.Depth {
		return nil, nil, ErrMaxDepthExceeded
	}

	var links []*url.URL
	var err error

	// Only parse HTML content.
	if IsHTML(baseURL.Path, data) {
		links, _, err = ParseAndRewriteLinks(bytes.NewReader(data), baseURL)
		if err != nil {
			return nil, nil, err
		}
	}

	var allLinks []*url.URL
	for _, link := range links {
		// Skip external links.
		if link.Host != baseURL.Host {
			continue
		}

		// Skip already visited URLs.
		if _, loaded := c.Visited.LoadOrStore(link.String(), true); loaded {
			continue
		}

		allLinks = append(allLinks, link)

		// Recursively fetch and parse links if depth allows.
		if curDepth < c.Depth {
			res, err := c.Fetch(link, 2)
			if err != nil {
				continue
			}

			subLinks, _, _ := ParseAllLinks(c, res, link, curDepth+1)
			allLinks = append(allLinks, subLinks...)
		}
	}

	return allLinks, nil, nil
}

// ParseAndRewriteLinks parses HTML content from the provided reader,
// extracts all links (href/src/data attributes), rewrites local links
// to relative paths for offline usage, and returns the list of extracted links
// and the modified HTML data.
func ParseAndRewriteLinks(r io.Reader, baseURL *url.URL) ([]*url.URL, []byte, error) {
	var links []*url.URL
	seen := make(map[string]struct{})

	// Map of HTML elements to their attribute containing URLs.
	attrMap := map[string]string{
		"a":      "href",
		"link":   "href",
		"img":    "src",
		"script": "src",
		"iframe": "src",
		"source": "src",
		"video":  "src",
		"audio":  "src",
		"embed":  "src",
		"object": "data",
	}

	// Parse HTML document.
	doc, err := html.Parse(r)
	if err != nil {
		return nil, nil, err
	}

	// Determine current HTML file path for relative rewriting.
	currentPath := filepath.Join(baseURL.Host, baseURL.Host, baseURL.Path)
	if strings.HasSuffix(baseURL.Path, "/") || filepath.Ext(baseURL.Path) == "" {
		currentPath = filepath.Join(currentPath, "index.html")
	}

	// Recursive function to traverse HTML nodes.
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if key, ok := attrMap[n.Data]; ok {
				for i, attr := range n.Attr {
					if attr.Key != key || strings.TrimSpace(attr.Val) == "" {
						continue
					}

					// Parse the attribute URL.
					parsedURL, err := url.Parse(strings.TrimSpace(attr.Val))
					if err != nil {
						continue
					}

					// Resolve relative URLs against the base URL.
					absURL := baseURL.ResolveReference(parsedURL)

					// Skip duplicates.
					if _, ok := seen[absURL.String()]; !ok {
						seen[absURL.String()] = struct{}{}
						links = append(links, absURL)
					}

					// Skip external links.
					if absURL.Host != baseURL.Host {
						continue
					}

					// Determine a local target path.
					targetPath := filepath.Join(baseURL.Host, baseURL.Host, absURL.Path)
					if strings.HasSuffix(absURL.Path, "/") || filepath.Ext(absURL.Path) == "" {
						targetPath = filepath.Join(targetPath, "index.html")
					}

					// Rewrite a link to relative path for offline use.
					relPath, err := filepath.Rel(filepath.Dir(currentPath), targetPath)
					if err != nil {
						continue
					}

					n.Attr[i].Val = filepath.ToSlash(relPath)

				}
			}
		}

		// Recurse into child nodes.
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// Start traversing from the root node.
	f(doc)

	// Render the modified HTML document to bytes.
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, nil, err
	}

	return links, buf.Bytes(), nil
}
