package wait

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func (w *Wait) PingAll() error {
	var wg sync.WaitGroup

	startTime := time.Now()
	finished := make(chan struct{}, 1)

	go func() {
		for _, host := range w.hosts {
			wg.Add(1)
			go w.ping(startTime, host, &wg)
		}

		wg.Wait()
		finished <- struct{}{}
	}()

	select {
	case <-finished:
		return nil
	case <-time.After(w.timeout):
		return fmt.Errorf("%s timeout reached before all hosts were up", w.timeout)
	}
}

func (w *Wait) ping(startTime time.Time, host host, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		conn, err := net.Dial(host.GetProtocol(), host.String())
		if err == nil {
			conn.Close()
			w.log.Printf("> up:   %s (after %s)", w.pad(host.String()), time.Since(startTime))
			return
		}

		w.log.Printf("> down: %s", w.pad(host.String()))
		time.Sleep(w.step)
	}
}

func (w *Wait) pad(str string) string {
	format := fmt.Sprintf("%%-%ds", w.padding)
	return fmt.Sprintf(format, str)
}
