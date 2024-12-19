package probes

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"
)

// TCPPinger is a pinger for TCP connections.
type TCPPinger struct {
	Host string
}

// Bootstrap sets up the pinger with the URL.
func (t *TCPPinger) Bootstrap(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("failed to parse host %q: %v", host, err)
	}

	if u.Host == "" {
		return fmt.Errorf("no host specified for tcp scheme")
	}

	if !oneOf(u.Scheme, "", "tcp", "tcp4", "tcp6") {
		return fmt.Errorf("invalid scheme for tcp probe: %s", u.Scheme)
	}

	t.Host = u.Host
	return nil
}

// Ping attempts to connect to the host.
func (t *TCPPinger) Ping(ctx context.Context) error {
	d := net.Dialer{Timeout: 1 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", t.Host)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
