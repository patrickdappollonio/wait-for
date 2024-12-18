package wait

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/patrickdappollonio/wait-for/wait/probes"
)

// App represents the application configuration.
type App struct {
	Hosts   []string
	Timeout time.Duration
	Every   time.Duration
	Verbose bool

	padding int
}

// Pinger defines the interface for a pinger.
type Pinger interface {
	Bootstrap(u *url.URL) error
	Ping(ctx context.Context) error
}

// pingerRegistry holds the mapping from protocol to pinger handler.
var pingerRegistry = map[string]func() Pinger{
	"tcp":   func() Pinger { return &probes.TCPPinger{} },
	"udp":   func() Pinger { return &probes.UDPPinger{} },
	"mysql": func() Pinger { return &probes.MySQLPinger{} },
}

// urlItem is a helper struct to hold the URL and the raw string.
type urlItem struct {
	URL *url.URL
	Raw string
}

// String returns the string representation of the URL.
func (u *urlItem) String() string {
	return fmt.Sprintf("%s://%s", u.URL.Scheme, u.URL.Host)
}

// stringifyHosts returns a string representation of the hosts, with all
// the URLs quoted and separated by commas.
func stringifyHosts(urls []urlItem) string {
	var sb strings.Builder

	for i, v := range urls {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(`"` + v.String() + `"`)
	}

	return sb.String()
}

// Run executes the application.
func (app *App) Run() error {
	if len(app.Hosts) == 0 {
		return fmt.Errorf("no hosts specified")
	}

	hostItems := make([]urlItem, 0, len(app.Hosts))
	for _, v := range app.Hosts {
		hostURL, err := parseHost(v)
		if err != nil {
			return fmt.Errorf("failed to parse host %q: %v", v, err)
		}

		if len(v) > app.padding {
			app.padding = len(hostURL.String())
		}

		hostItems = append(hostItems, urlItem{hostURL, v})
	}

	ctx, cancel := context.WithTimeout(context.Background(), app.Timeout)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, len(hostItems))

	fmt.Fprintln(
		os.Stdout,
		"Waiting for hosts:", stringifyHosts(hostItems),
		fmt.Sprintf("(timeout: %s, attempting every %s)", app.Timeout, app.Every),
	)

	startTime := time.Now()
	for _, host := range hostItems {
		wg.Add(1)

		go func(h urlItem, startTime time.Time) {
			defer wg.Done()

			pingerCtor, ok := pingerRegistry[h.URL.Scheme]
			if !ok {
				errChan <- fmt.Errorf("no handler registered for scheme %q (host: %q)", h.URL.Scheme, h.Raw)
				return
			}

			pinger := pingerCtor()
			if err := pinger.Bootstrap(h.URL); err != nil {
				errChan <- fmt.Errorf("failed to bootstrap configuration for host %q: %v", h, err)
				return
			}

			if err := pinger.Ping(ctx); err == nil {
				// Host is reachable, break the loop.
				if app.Verbose {
					fmt.Fprintf(os.Stdout, "> up:   %s (after %s)\n", app.pad(host.URL.String()), time.Since(startTime))
				}
				return
			} else {
				if app.Verbose {
					fmt.Fprintf(os.Stdout, "> down: %s -- %s\n", app.pad(host.URL.String()), err.Error())
				}
			}

			ticker := time.NewTicker(app.Every)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					errChan <- fmt.Errorf("timeout reached while waiting for %q", h)
					return
				case <-ticker.C:
					if err := pinger.Ping(ctx); err == nil {
						// Host is reachable, break the loop.
						if app.Verbose {
							fmt.Fprintf(os.Stdout, "> up:   %s (after %s)\n", app.pad(host.URL.String()), time.Since(startTime))
						}
						return
					} else {
						if app.Verbose {
							fmt.Fprintf(os.Stdout, "> down: %s -- %s\n", app.pad(host.URL.String()), err.Error())
						}
					}
				}
			}
		}(host, startTime)
	}

	// Wait for all goroutines or for the first error.
	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-doneChan:
		// All goroutines finished. Check if any error was reported.
		select {
		case err := <-errChan:
			if err != nil {
				return err
			}
		default:
			// No errors, success!
			return nil
		}
	case err := <-errChan:
		// Immediately return the first error encountered.
		return err
	case <-ctx.Done():
		// Global timeout triggered.
		return fmt.Errorf("%s timeout reached before all hosts were up", app.Timeout)
	}

	return nil
}

// parseHost parses the host string and returns a URL.
func parseHost(hostStr string) (*url.URL, error) {
	// If no scheme, assume tcp
	if !strings.Contains(hostStr, "://") {
		hostStr = "tcp://" + hostStr
	}
	return url.Parse(hostStr)
}

// pad pads the string to the configured padding based on the longest host
// full string URL representation (including protocol).
func (app *App) pad(str string) string {
	format := fmt.Sprintf("%%-%ds", app.padding)
	return fmt.Sprintf(format, str)
}
