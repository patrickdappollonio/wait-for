package probes

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestTCPPinger_Bootstrap(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			urlStr:  "tcp://example.com:80",
			wantErr: false,
		},
		{
			name:    "No host specified",
			urlStr:  "tcp://",
			wantErr: true,
		},
		{
			name:    "Invalid scheme",
			urlStr:  "http://example.com:80",
			wantErr: true,
		},
		{
			name:    "No scheme",
			urlStr:  "example.com:80",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pinger := &TCPPinger{}
			if err := pinger.Bootstrap(tt.urlStr); (err != nil) != tt.wantErr {
				t.Errorf("TCPPinger.Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTCPPinger_Ping(t *testing.T) {
	// Launch a local server to test the pinger
	chPort := make(chan string, 1)

	go func() {
		srv, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			t.Errorf("net.Listen() error = %v", err)
		}
		defer srv.Close()

		chPort <- srv.Addr().String()

		for {
			conn, err := srv.Accept()
			if err != nil {
				return
			}
			conn.Close()
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
			pinger := &TCPPinger{Host: tt.host}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := pinger.Ping(ctx); (err != nil) != tt.wantErr {
				t.Errorf("TCPPinger.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
