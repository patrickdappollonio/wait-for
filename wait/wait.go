package wait

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Wait struct {
	hosts   []string
	timeout time.Duration
	step    time.Duration
	log     *log.Logger
	padding int
}

func New(hosts []string, step, timeout time.Duration, verbose bool) (*Wait, error) {
	w := &Wait{
		hosts:   hosts,
		timeout: timeout,
		step:    step,
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts specified")
	}

	for _, v := range hosts {
		if len(v) > w.padding {
			w.padding = len(v)
		}

		if _, _, err := net.SplitHostPort(v); err != nil {
			return nil, fmt.Errorf("invalid host format: %q -- must be in the format \"host:port\"", v)
		}
	}

	w.log = log.New(ioutil.Discard, "", 0)

	if verbose {
		w.log.SetOutput(os.Stdout)
	}

	return w, nil
}

func (w *Wait) String() string {
	return fmt.Sprintf(
		"Waiting for hosts: %s (timeout: %s, attempting every %s)",
		strings.Join(w.hosts, ", "),
		w.timeout,
		w.step,
	)
}
