// File: probes/https_pinger.go
package probes

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// validateURL checks if the URL is valid. There are no checks
// to validate if the URL is nil.
func validateURL(scheme string, u *url.URL) error {
	// Validate the URL scheme
	if u.Scheme != scheme {
		return fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	if !u.IsAbs() {
		return fmt.Errorf("invalid URL: %s", u.String())
	}

	if u.Hostname() == "" {
		return fmt.Errorf("no host specified: %s", u.String())
	}

	return nil
}

// HTTPSPinger is a pinger for HTTPS connections.
type HTTPSPinger struct {
	url        *url.URL
	httpClient *http.Client
}

// Bootstrap sets up the pinger with the HTTPS URL.
func (h *HTTPSPinger) Bootstrap(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("failed to parse host %q: %v", host, err)
	}

	if err := validateURL("https", u); err != nil {
		return err
	}

	h.url = u

	// Initialize HTTPS client with timeout and TLS configuration
	h.httpClient = &http.Client{
		Timeout: 1 * time.Second, // 1 second timeout per request

		// Ensure TLS certificate verification is enabled
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Do not skip TLS verification
			},
		},
	}

	return nil
}

// Ping performs an HTTPS GET request and checks the status code.
func (h *HTTPSPinger) Ping(ctx context.Context) error {
	return doGet(ctx, h.httpClient, h.url.String())
}

// HTTPPinger is a pinger for HTTP connections.
type HTTPPinger struct {
	url        *url.URL
	HTTPClient *http.Client
}

// Bootstrap sets up the pinger with the HTTP URL.
func (h *HTTPPinger) Bootstrap(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("failed to parse host %q: %v", host, err)
	}

	if err := validateURL("http", u); err != nil {
		return err
	}

	h.url = u

	// Initialize HTTP client with timeout for each request
	h.HTTPClient = &http.Client{
		Timeout: 1 * time.Second, // 1 second timeout per request
	}

	return nil
}

// Ping performs an HTTP GET request and checks the status code.
func (h *HTTPPinger) Ping(ctx context.Context) error {
	return doGet(ctx, h.HTTPClient, h.url.String())
}
