package probes

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"
)

// UDPPinger is a pinger for UDP connections.
type UDPPinger struct {
	Host string
}

// Bootstrap sets up the pinger with the URL.
func (u *UDPPinger) Bootstrap(host string) error {
	url, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("failed to parse host %q: %v", host, err)
	}

	if url.Host == "" {
		return fmt.Errorf("no host specified for udp scheme")
	}

	if !oneOf(url.Scheme, "udp", "udp4", "udp6") {
		return fmt.Errorf("invalid scheme for udp probe: %s", url.Scheme)
	}

	u.Host = url.Host
	return nil
}

// Ping attempts to send a datagram to the host.
func (u *UDPPinger) Ping(ctx context.Context) error {
	// For UDP "ping", we can attempt to send a datagram and check for error.
	// Unlike TCP, we don't get a "connected" state just by dialing.
	conn, err := net.DialTimeout("udp", u.Host, 1*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send a zero-length packet just to see if it errors out.
	_, err = conn.Write([]byte{})
	return err
}
