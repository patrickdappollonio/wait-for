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
func (t *TCPPinger) Bootstrap(u *url.URL) error {
	if u.Host == "" {
		return fmt.Errorf("no host specified for tcp scheme")
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
