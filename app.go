package main

import (
	"fmt"
	"time"

	"github.com/patrickdappollonio/wait-for/wait"
	"github.com/spf13/cobra"
)

var version = "development"

const (
	helpShort = "wait-for allows you to wait for a resource to respond to requests."

	helpLong = `wait-for allows you to wait for a resource to respond to requests.

It does this by performing a connection to the specified host and port. If there's
no resource behind it and the connection cannot be established, the request is retried
until either the timeout is reached or the resource becomes available.

Each protocol defines its own way of checking for the resource. For example, a TCP
connection will attempt to connect to the host and port specified, while a MySQL
connection will attempt to connect to the host and port and then ping the database.

By default, the standard timeout is 10 seconds.

For documentation, visit: https://github.com/patrickdappollonio/wait-for.`
)

func root() *cobra.Command {
	var app wait.App

	rootCommand := &cobra.Command{
		Use:           "wait-for",
		Short:         helpShort,
		Long:          helpLong,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := app.Run(); err != nil {
				return err
			}

			fmt.Println("All hosts are up and responding.")
			return nil
		},
	}

	rootCommand.Flags().StringSliceVarP(&app.Hosts, "host", "s", []string{}, "hosts to connect to in the format \"host:port\" with optional protocol prefix (tcp:// or udp://)")
	rootCommand.Flags().DurationVarP(&app.Timeout, "timeout", "t", time.Second*10, "maximum time to wait for the endpoints to respond before giving up")
	rootCommand.Flags().DurationVarP(&app.Every, "every", "e", time.Second*1, "time to wait between each request attempt against the host")
	rootCommand.Flags().BoolVarP(&app.Verbose, "verbose", "v", false, "enable verbose output -- will print every time a request is made")

	return rootCommand
}
