package main

import (
	"fmt"
	"time"

	"github.com/patrickdappollonio/wait-for/wait"
	"github.com/spf13/cobra"
)

var version = "development"

const (
	helpShort = "wait-for allows you to wait for a TCP resource to respond to requests."

	helpLong = `wait-for allows you to wait for a TCP resource to respond to requests.

It does this by performing a TCP connection to the specified host and port. If there's
no resource behind it and the connection cannot be established, the request is retried
until either the timeout is reached or the resource becomes available.

By default, the standard timeout is 10 seconds.

For documentation, visit: https://github.com/patrickdappollonio/wait-for.`
)

func root() *cobra.Command {
	var (
		hosts   []string
		timeout time.Duration
		step    time.Duration
		verbose bool
	)

	rootCommand := &cobra.Command{
		Use:           "wait-for",
		Short:         helpShort,
		Long:          helpLong,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			w, err := wait.New(hosts, step, timeout, verbose)
			if err != nil {
				return err
			}

			fmt.Println(w.String())
			if err := w.PingAll(); err != nil {
				return err
			}

			fmt.Println("All hosts are up and responding.")
			return nil
		},
	}

	rootCommand.Flags().StringSliceVarP(&hosts, "host", "s", []string{}, "hosts to connect to in the format \"host:port\" with optional protocol prefix (tcp:// or udp://)")
	rootCommand.Flags().DurationVarP(&timeout, "timeout", "t", time.Second*10, "maximum time to wait for the endpoints to respond before giving up")
	rootCommand.Flags().DurationVarP(&step, "every", "e", time.Second*1, "time to wait between each request attempt against the host")
	rootCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output -- will print every time a request is made")

	return rootCommand
}
