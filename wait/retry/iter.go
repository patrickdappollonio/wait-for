package retry

import (
	"context"
	"errors"
	"time"
)

type Retry struct {
	interval time.Duration
	context  context.Context
}

var ErrTryAgain = errors.New("retry: again")
var ErrTimedOut = errors.New("retry: timed out")

// type ErrStop struct {
// 	Parent error
// }

// func (e ErrStop) Error() string {
// 	return "stop: " + e.Parent.Error()
// }

// func (e ErrStop) Unwrap() error {
// 	return e.Parent
// }

// func (e ErrStop) Is(target error) bool {
// 	_, is := target.(*ErrStop)
// 	return is
// }

func New(interval time.Duration, ctx context.Context) *Retry {
	return &Retry{interval: interval, context: ctx}
}

func (it *Retry) Run(fn func() error) error {
	for {
		select {
		case <-it.context.Done():
			return ErrTimedOut

		default:
			err := fn()
			if err != nil {
				if errors.Is(err, ErrTryAgain) {
					time.Sleep(it.interval)
					continue
				}

				// if errors.Is(err, &ErrStop{}) {
				// 	return err.(*ErrStop).Unwrap()
				// }

				return err
			}

			return nil
		}
	}
}
