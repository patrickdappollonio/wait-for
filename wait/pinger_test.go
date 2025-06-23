package wait

import (
	"io"
	"os"
	"testing"
)

func TestParseHost(t *testing.T) {
	tests := []struct {
		name    string
		hostStr string
		want    *matchedURLItem
		wantErr bool
	}{
		{
			name:    "Valid TCP URL",
			hostStr: "example.com",
			want: &matchedURLItem{
				Raw:    "tcp://example.com",
				Pinger: pingerRegistry["tcp"](),
			},
			wantErr: false,
		},
		{
			name:    "Valid HTTP URL",
			hostStr: "http://example.com",
			want: &matchedURLItem{
				Raw:    "http://example.com",
				Pinger: pingerRegistry["http"](),
			},
			wantErr: false,
		},
		{
			name:    "Invalid URL without scheme",
			hostStr: "://example.com",
			wantErr: true,
		},
		{
			name:    "Invalid URL with unknown scheme",
			hostStr: "unknown://example.com",
			wantErr: true,
		},
		{
			name:    "Invalid URL with no scheme and no host",
			hostStr: "://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHost(tt.hostStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.Raw != tt.want.Raw {
				t.Errorf("parseHost() got = %v, want %v", got.Raw, tt.want.Raw)
			}

			if !tt.wantErr && got.Pinger == nil {
				t.Errorf("parseHost() got Pinger = nil, want non-nil")
			}
		})
	}
}

func TestStringifyHosts(t *testing.T) {
	tests := []struct {
		name string
		urls []matchedURLItem
		want string
	}{
		{
			name: "Single URL",
			urls: []matchedURLItem{
				{
					Raw: "http://example.com",
				},
			},
			want: `"http://example.com"`,
		},
		{
			name: "Multiple URLs",
			urls: []matchedURLItem{
				{
					Raw: "http://example.com",
				},
				{
					Raw: "https://example.org",
				},
			},
			want: `"http://example.com", "https://example.org"`,
		},
		{
			name: "No URLs",
			urls: []matchedURLItem{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringifyHosts(tt.urls); got != tt.want {
				t.Errorf("stringifyHosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppPad(t *testing.T) {
	a := &App{padding: 5}
	if got := a.pad("a"); got != "a    " {
		t.Errorf("pad() = %q", got)
	}
}

func TestPrintOnVerbose(t *testing.T) {
	a := &App{Verbose: true}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	old := os.Stdout
	os.Stdout = w
	a.printOnVerbose("msg")
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	if string(out) != "msg\n" {
		t.Errorf("unexpected output: %q", string(out))
	}

	// now ensure nothing printed when verbose false
	r, w, _ = os.Pipe()
	os.Stdout = w
	a.Verbose = false
	a.printOnVerbose("nope")
	w.Close()
	os.Stdout = old
	out, _ = io.ReadAll(r)
	if len(out) != 0 {
		t.Errorf("expected no output, got %q", string(out))
	}
}
