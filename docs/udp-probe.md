# UDP

The UDP probe will attempt to connect to the host and port specified. If the connection can be established successfully and at least one zero-length packet can be sent, the probe will exit successfully.

If the connection cannot be established, the probe will retry until either the timeout is reached or the resource becomes available.

By default, if the protocol isn't specified when providing a `--host` flag, TCP is assumed. To use UDP, you must prefix the host with `udp://`:

```bash
wait-for --host "udp://localhost:53"
```

UDP probes make no guarantees the response received from the server has any sort of validity. They merely check if the server is accepting connections on the specified port and if a zero-length packet can be sent.
