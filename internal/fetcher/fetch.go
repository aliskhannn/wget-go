package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Fetch fetches the content of the given URL and returns the response body.
// Caller is responsible for saving the data to a file.
func Fetch(u *url.URL, client *http.Client) ([]byte, error) {
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", u, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// check for non-200 response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s: status code: %d", u, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body %s: %w", u, err)
	}

	return data, nil
}
