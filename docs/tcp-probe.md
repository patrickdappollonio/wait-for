# TCP

The TCP probe will attempt to connect to the host and port specified. If the connection can be established successfully, the probe will exit successfully.

If the connection cannot be established, the probe will retry until either the timeout is reached or the resource becomes available.

By default, if the protocol isn't specified when providing a `--host` flag, TCP is assumed. This means that the following two commands are equivalent:

```bash
wait-for --host "localhost:80"
wait-for --host "tcp://localhost:80"
```

TCP probes make no guarantees the response received from the server has any sort of validity. They merely check if the server is accepting connections on the specified port.
