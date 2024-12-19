# `wait-for` documentation

`wait-for` allows you to wait for a resource to respond to requests. It does this by performing a connection to the specified host provided either by a configuration file or by command line arguments. If there's no resource behind it and the connection cannot be established, the request is retried until either the timeout is reached or the resource becomes available.

## Configuration

The application can be configured either by command line arguments or by a configuration file. The configuration file is a YAML file that can be passed to the application using the `--config` flag.

Host flags (those with `--host` or `-s`) can be specified both via the command line and the configuration file. If a host is specified in both, it will be pinged twice.

For more information on how to use the configuration file, please refer to the [configuration file documentation](configuration-file.md).

## Supported probes

"Probes" are the way `wait-for` checks for the availability of a resource. Each probe maps to a specific protocol and checks for the availability of a resource in a specific way.

The following probes are supported:

* [TCP probe](tcp-probe.md)
* [UDP probe](udp-probe.md)
* [HTTP & HTTPS probe](http-https-probe.md)
* [MySQL probe](mysql-probe.md) *(experimental)*
* [PostgreSQL probe](postgres-probe.md) *(experimental)*

## Adding new probes

If you want to add a new probe, you can do so by implementing the `Probe` interface. The interface is defined as follows:

```go
// Pinger defines the interface for a pinger.
type Pinger interface {
	Bootstrap(host string) error
	Ping(ctx context.Context) error
}
```

Then, the probe has to be matched to a protocol stored in the `pingerRegistry` variable in the `wait` package. This is done by adding a new entry to the map, where the key is the protocol and the value is the probe implementation. Currently, the following protocols are supported:

```go
// pingerRegistry holds the mapping from protocol to pinger handler.
// Add your own pinger here.
var pingerRegistry = map[string]func() Pinger{
	"tcp":      func() Pinger { return &probes.TCPPinger{} },
	"udp":      func() Pinger { return &probes.UDPPinger{} },
	"mysql":    func() Pinger { return &probes.MySQLPinger{} },
	"postgres": func() Pinger { return &probes.PostgresPinger{} },
	"http":     func() Pinger { return &probes.HTTPPinger{} },
	"https":    func() Pinger { return &probes.HTTPSPinger{} },
}
```

When creating your own probes, the following rules apply:

* A protocol must be unique.
* The protocol must be lowercase.
* A probe `struct` should accept no parameters.
* A probe `Bootstrap` method should accept a `host` parameter which should validate the host and set up the probe, or return an error if the host is invalid.
* A probe `Ping` method should accept a `context.Context` parameter and return an error if the ping fails or the context is canceled.
* I reserve the discretion to accept or reject any pull request that adds a new probe.
