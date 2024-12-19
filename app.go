package main

import (
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/patrickdappollonio/wait-for/wait"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "development"

const (
	helpShort = "wait-for allows you to wait for a resource to respond to requests."

	helpLong = `wait-for allows you to wait for a resource to respond to requests.

It does this by performing a connection to the specified host and port. If there's no resource behind it and the connection cannot be established, the request is retried until either the timeout is reached or the resource becomes available.

Each protocol defines its own way of checking for the resource. For example, a TCP connection will attempt to connect to the host and port specified, while a MySQL connection will attempt to connect to the host and port, and then ping the database.

By default, the standard timeout is 10 seconds but it can be customized for all requests. The time between each request is 1 second, but this can also be customized.

For documentation, visit: https://github.com/patrickdappollonio/wait-for.`
)

func root() *cobra.Command {
	var cfgFile string
	var hosts []string

	rootCommand := &cobra.Command{
		Use:     "wait-for",
		Short:   helpShort,
		Long:    wrap(helpLong, 80),
		Version: version,
		Example: exampleCommands("wait-for", []example{
			{command: "-s localhost:80", helper: "wait for a web server to accept connections"},
			{command: "-s mysql.example.local:3306", helper: "wait for a MySQL database to accept connections"},
			{command: "-s udp://localhost:53", helper: "wait for a DNS server to accept connections"},
			{command: "--host localhost:80 --host localhost:81", helper: "wait for multiple resources to accept connections"},
			{command: "--host mysql://localhost:3306", helper: "wait until a MySQL database is ready to accept connections and responds to pings"},
			{command: "--host postgres://localhost:5432", helper: "wait until a PostgreSQL database is ready to accept connections and responds to pings"},
			{command: "--host http://localhost:8080", helper: "wait until an HTTP server is ready to accept connections and responds to requests with a 200-299 status code"},
			{command: "--host https://localhost:443", helper: "wait until an HTTPS server is ready to accept connections and responds to requests with a 200-299 status code and a valid certificate"},
			{command: "--config targets.yaml", helper: "load hosts and settings from a YAML file"},
		}),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			// Read config file if available
			viper.SetConfigFile(cfgFile)
			if err := viper.ReadInConfig(); err != nil {
				// If config not found, it's not fatal unless we rely on it
				if errors.Is(err, &fs.PathError{}) {
					return fmt.Errorf("error reading config file: %w", err)
				}
			}

			// Merge hosts flag with viper flags
			hosts = append(hosts, viper.GetStringSlice("host")...)

			// Retrieve final values from viper after merging CLI flags and config
			app := &wait.App{
				Hosts:   hosts,
				Timeout: viper.GetDuration("timeout"),
				Every:   viper.GetDuration("every"),
				Verbose: viper.GetBool("verbose"),
			}

			// Run the application
			if err := app.Run(); err != nil {
				return err
			}

			fmt.Println("All hosts are up and responding.")
			return nil
		},
	}

	// Flags for the program
	rootCommand.Flags().StringSliceVarP(&hosts, "host", "s", []string{}, `hosts to connect to in the format "host:port" or with protocol prefix for one of the supported protocols (e.g. "udp://host:port")`)
	rootCommand.Flags().DurationP("timeout", "t", 10*time.Second, "maximum time to wait for the endpoints to respond before giving up")
	rootCommand.Flags().DurationP("every", "e", 1*time.Second, "time to wait between each request attempt against the host")
	rootCommand.Flags().BoolP("verbose", "v", false, "enable verbose output -- will print every time a request is made")
	rootCommand.Flags().StringVar(&cfgFile, "config", "targets.yaml", "config file to load hosts and settings from")

	// Bind flags to viper except hosts and config file
	viper.BindPFlag("timeout", rootCommand.Flags().Lookup("timeout"))
	viper.BindPFlag("every", rootCommand.Flags().Lookup("every"))
	viper.BindPFlag("verbose", rootCommand.Flags().Lookup("verbose"))

	return rootCommand
}
