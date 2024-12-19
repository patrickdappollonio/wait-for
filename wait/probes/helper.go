package probes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// oneOf returns true if the first argument is equal to any of the
// following arguments.
func oneOf[T comparable](s T, values ...T) bool {
	for _, v := range values {
		if s == v {
			return true
		}
	}

	return false
}

// unwrapError recursively unwraps the error to get the root cause.
func unwrapError(err error) error {
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		return unwrapError(unwrapped)
	}

	return err
}

// doGet performs a GET request to the given URL with the provided client
// and context, then checks the status code to ensure it is in the 2xx range.
func doGet(ctx context.Context, client *http.Client, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return unwrapError(err)
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("received non-2xx status code: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}

// extractProtocol extracts the protocol from the host string.
// If no protocol is found, an empty string is returned.
func extractProtocol(host string) string {
	// Find if there's a "://" in the host string.
	// If there is, extract the protocol.
	// If there isn't, assume it's a hostname and return an empty string.
	if i := strings.Index(host, "://"); i >= 0 {
		return host[:i]
	}

	return ""
}
