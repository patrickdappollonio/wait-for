package wait

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/patrickdappollonio/wait-for/wait/probes"
	"golang.org/x/sync/errgroup"
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
	Bootstrap(host string) error
	Ping(ctx context.Context) error
}

// pingerRegistry holds the mapping from protocol to pinger handler.
// Add your own pinger here.
var pingerRegistry = map[string]func() Pinger{
	"tcp":      func() Pinger { return &probes.TCPPinger{} },
	"udp":      func() Pinger { return &probes.UDPPinger{} },
	"mysql":    func() Pinger { return &probes.MySQLPinger{} },
	"postgres": func() Pinger { return &probes.PostgresPinger{} },
	"http":     func() Pinger { return &probes.HTTPPinger{} },
	"https":    func() Pinger { return &probes.HTTPSPinger{} },
}

// matchedURLItem is a helper struct to hold the URL and the raw string.
type matchedURLItem struct {
	Raw    string
	Pinger Pinger
}

// String returns the string representation of the URL.
func (u *matchedURLItem) String() string {
	return u.Raw
}

// stringifyHosts returns a string representation of the hosts, with all
// the URLs quoted and separated by commas.
func stringifyHosts(urls []matchedURLItem) string {
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

	hostItems := make([]matchedURLItem, 0, len(app.Hosts))
	for _, rawURL := range app.Hosts {
		// Parse the host URL
		matched, err := parseHost(rawURL)
		if err != nil {
			return fmt.Errorf("failed to parse host %q: %v", rawURL, err)
		}

		// Calculate the padding for the output
		if len(rawURL) > app.padding {
			app.padding = len(matched.Raw)
		}

		// Bootstrap the pinger and validate its URL
		if err := matched.Pinger.Bootstrap(rawURL); err != nil {
			return fmt.Errorf("failed to bootstrap host %q: %v", rawURL, err)
		}

		// Append the host to the list
		hostItems = append(hostItems, *matched)
	}

	// Register signal handlers for early termination.
	sigterm, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer done()

	// Create a context with a timeout for the maximum time allowed.
	ctx, cancel := context.WithTimeout(context.Background(), app.Timeout)
	defer cancel()

	// Document what are we doing to the user.
	fmt.Fprintln(
		os.Stdout,
		"Waiting for hosts:", stringifyHosts(hostItems),
		fmt.Sprintf("(timeout: %s, attempting every %s)", app.Timeout, app.Every),
	)

	// Create an error group.
	var eg errgroup.Group

	// Iterate over all hosts and ping them.
	for _, host := range hostItems {
		eg.Go(app.handlePing(ctx, sigterm, host))
	}

	// Create a channel to signal when all goroutines are done.
	doneChan := make(chan error)

	go func() {
		// Wait for all goroutines or for the first error.
		if err := eg.Wait(); err != nil {
			doneChan <- err
		}

		// All goroutines finished successfully.
		close(doneChan)
	}()

	select {
	case err := <-doneChan:
		// Immediately return the first error encountered.
		return err
	case <-ctx.Done():
		// Global timeout triggered.
		return fmt.Errorf("%s timeout reached before all hosts were up", app.Timeout)
	}
}

// handlePing pings the host asynchronously and returns an error if the host
// is not reachable.
func (app *App) handlePing(ctx, sigterm context.Context, h matchedURLItem) func() error {
	return func() error {
		startTime := time.Now()

		// Ping right away the first time
		if err := h.Pinger.Ping(ctx); err == nil {
			app.printOnVerbose("> up:   %s (after %s)", app.pad(h.Raw), time.Since(startTime))
			return nil // Host is reachable, break the loop.
		} else {
			app.printOnVerbose("> down: %s -- %s", app.pad(h.Raw), err.Error())
		}

		// Create a ticker to ping the host every `app.Every` duration.
		ticker := time.NewTicker(app.Every)
		defer ticker.Stop()

		for {
			select {
			case <-sigterm.Done():
				// User requested early termination.
				return fmt.Errorf("user requested early termination")
			case <-ctx.Done():
				// Timeout reached.
				return fmt.Errorf("timeout reached while waiting for %q", h)
			case <-ticker.C:
				// Ping the host and check if it's reachable.
				if err := h.Pinger.Ping(ctx); err == nil {
					app.printOnVerbose("> up:   %s (after %s)", app.pad(h.Raw), time.Since(startTime))
					return nil // Host is reachable, break the loop.
				} else {
					app.printOnVerbose("> down: %s -- %s", app.pad(h.Raw), err.Error())
				}
			}
		}
	}
}

// printOnVerbose prints the message to the standard output if the verbose
// flag is enabled.
func (app *App) printOnVerbose(format string, args ...interface{}) {
	if app.Verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

// parseHost parses the host string and returns a URL.
func parseHost(hostStr string) (*matchedURLItem, error) {
	// If no scheme, assume tcp
	if !strings.Contains(hostStr, "://") {
		hostStr = "tcp://" + hostStr
	}

	// Parse the URL without url.Parse
	if !strings.Contains(hostStr, "://") {
		return nil, fmt.Errorf("invalid URL: %s", hostStr)
	}

	// Find the scheme
	schemeEnd := strings.Index(hostStr, "://")
	if schemeEnd == -1 {
		return nil, fmt.Errorf("invalid URL: %s", hostStr)
	}

	scheme := hostStr[:schemeEnd]

	// Find the Pinger
	pingerCtor, ok := pingerRegistry[scheme]
	if !ok {
		return nil, fmt.Errorf("no handler registered for scheme %q (host: %q)", scheme, hostStr)
	}

	// Return the URL and the Pinger
	return &matchedURLItem{
		Raw:    hostStr,
		Pinger: pingerCtor(),
	}, nil
}

// pad pads the string to the configured padding based on the longest host
// full string URL representation (including protocol).
func (app *App) pad(str string) string {
	format := fmt.Sprintf("%%-%ds", app.padding)
	return fmt.Sprintf(format, str)
}
