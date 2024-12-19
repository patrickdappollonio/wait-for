package probes

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestUDPPinger_Bootstrap(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			urlStr:  "udp://example.com:80",
			wantErr: false,
		},
		{
			name:    "No host specified",
			urlStr:  "udp://",
			wantErr: true,
		},
		{
			name:    "Invalid scheme",
			urlStr:  "http://example.com:80",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pinger := &UDPPinger{}
			if err := pinger.Bootstrap(tt.urlStr); (err != nil) != tt.wantErr {
				t.Errorf("UDPPinger.Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUDPPinger_Ping(t *testing.T) {
	// Launch a local server to test the pinger
	chPort := make(chan string, 1)

	go func() {
		srv, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
		if err != nil {
			t.Errorf("net.Listen() error = %v", err)
		}
		defer srv.Close()

		chPort <- srv.LocalAddr().String()

		for {
			buf := make([]byte, 1024)
			n, _, err := srv.ReadFromUDP(buf)
			if err != nil {
				t.Errorf("srv.ReadFrom() error = %v", err)
			}

			if n >= 0 {
				srv.Close()
				break
			}
		}
	}()

	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{
			name:    "Valid host",
			host:    <-chPort,
			wantErr: false,
		},
		{
			name:    "Invalid host",
			host:    "invalidhost:80",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pinger := &UDPPinger{Host: tt.host}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := pinger.Ping(ctx); (err != nil) != tt.wantErr {
				t.Errorf("UDPPinger.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
