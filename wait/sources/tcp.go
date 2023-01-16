package sources

import (
	"context"
	"net"
	"time"

	"github.com/patrickdappollonio/wait-for/wait/retry"
)

type TCP struct {
	hostname string
	interval time.Duration
	logger   Logger
}

func (t *TCP) SetLogger(logger Logger) {
	t.logger = logger
}

func (t *TCP) SetInterval(d time.Duration) {
	t.interval = d
}

func (t *TCP) Ping(ctx context.Context) error {
	startTime := time.Now()
	return retry.New(t.interval, ctx).Run(func() error {
		conn, err := net.Dial("tcp", t.hostname)
		if err != nil {
			t.logger.Printf("%s %s", downPrefix, t.hostname)
			return retry.ErrTryAgain
		}

		conn.Close()
		t.logger.Printf("%s %s (after %s)", upPrefix, t.hostname, time.Since(startTime))
		return nil
	})
}
