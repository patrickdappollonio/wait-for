package wait

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type proto string

const (
	tcp proto = "tcp"
	udp proto = "udp"
)

type host struct {
	host     string
	port     string
	protocol proto
}

func (h host) String() string {
	return fmt.Sprintf("%s:%s", h.host, h.port)
}

func (h host) GetProtocol() string {
	if h.protocol == "" {
		return string(tcp)
	}

	return string(h.protocol)
}

func stringifyHosts(hosts []host) string {
	var sb strings.Builder

	for i, v := range hosts {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(`"` + fmt.Sprintf("%s://%s:%s", v.GetProtocol(), v.host, v.port) + `"`)
	}

	return sb.String()
}

type Wait struct {
	hosts   []host
	timeout time.Duration
	step    time.Duration
	log     *log.Logger
	padding int
}

var reLooksLikeProtocol = regexp.MustCompile(`^(\w+)://`)

func New(hosts []string, step, timeout time.Duration, verbose bool) (*Wait, error) {
	w := &Wait{
		timeout: timeout,
		step:    step,
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts specified")
	}

	full := make([]host, 0, len(hosts))
	for _, v := range hosts {
		if len(v) > w.padding {
			w.padding = len(v)
		}

		var proto proto

		if strings.HasPrefix(v, "tcp://") {
			proto = tcp
			v = strings.TrimPrefix(v, "tcp://")
		}

		if strings.HasPrefix(v, "udp://") {
			proto = udp
			v = strings.TrimPrefix(v, "udp://")
		}

		if proto == "" && reLooksLikeProtocol.MatchString(v) {
			return nil, fmt.Errorf("invalid protocol specified: %q -- only \"tcp\" and \"udp\" are supported", v)
		}

		parsedHost, parsedPort, err := net.SplitHostPort(v)
		if err != nil {
			return nil, fmt.Errorf("invalid host format: %q -- must be in the format \"host:port\" or \"(tcp|udp)://host:port\"", v)
		}

		full = append(full, host{
			host:     parsedHost,
			port:     parsedPort,
			protocol: proto,
		})
	}

	w.hosts = full
	w.log = log.New(io.Discard, "", 0)

	if verbose {
		w.log.SetOutput(os.Stdout)
	}

	return w, nil
}

func (w *Wait) String() string {
	return fmt.Sprintf(
		"Waiting for hosts: %s (timeout: %s, attempting every %s)",
		stringifyHosts(w.hosts),
		w.timeout,
		w.step,
	)
}
