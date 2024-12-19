package probes

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPPinger_Bootstrap(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			urlStr:  "http://example.com",
			wantErr: false,
		},
		{
			name:    "No host specified",
			urlStr:  "http://",
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			urlStr:  "://example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pinger := &HTTPPinger{}
			if err := pinger.Bootstrap(tt.urlStr); (err != nil) != tt.wantErr {
				t.Errorf("HTTPPinger.Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPPinger_Ping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "Status OK",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Status Not Found",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that returns the specified status code
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer ts.Close()

			pinger := &HTTPPinger{}
			if err := pinger.Bootstrap(ts.URL); err != nil {
				t.Fatalf("HTTPPinger.Bootstrap() error = %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			if err := pinger.Ping(ctx); (err != nil) != tt.wantErr {
				t.Errorf("HTTPPinger.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func generateSelfSignedCert() (tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Example Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return tls.X509KeyPair(certPEM, keyPEM)
}

func TestHTTPSPinger_Bootstrap(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			urlStr:  "https://example.com",
			wantErr: false,
		},
		{
			name:    "No host specified",
			urlStr:  "https://",
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			urlStr:  "://example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pinger := &HTTPSPinger{}
			if err := pinger.Bootstrap(tt.urlStr); (err != nil) != tt.wantErr {
				t.Errorf("HTTPSPinger.Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPSPinger_Ping(t *testing.T) {
	cert, err := generateSelfSignedCert()
	if err != nil {
		t.Fatalf("Failed to generate self-signed certificate: %v", err)
	}

	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "Status OK",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Status Not Found",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			ts.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
			ts.StartTLS()
			defer ts.Close()

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]}))

			pinger := &HTTPSPinger{}
			if err := pinger.Bootstrap(ts.URL); err != nil {
				t.Fatalf("HTTPSPinger.Bootstrap() error = %v", err)
			}

			// Override the bootstrapped client
			pinger.httpClient = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: caCertPool,
					},
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			if err := pinger.Ping(ctx); (err != nil) != tt.wantErr {
				t.Errorf("HTTPSPinger.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
