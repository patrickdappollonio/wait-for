# HTTP & HTTPS

The HTTP and HTTPS probes are used to send an HTTP or HTTPS `GET` request to a server and check the response. A request is successful not only if the HTTP server was able to provide a connection but also if the response status code is within the range of 200 to 299. If the request responds within this range, the probe will exit successfully.

If the connection cannot be established or the response status code is outside the range of 200 to 299, the probe will retry until either the timeout is reached or the resource becomes available.

An example request to `http://localhost:80` would look like this:

```bash
wait-for --host "http://localhost:80"
```

An example request to `https://localhost:443` would look like this:

```bash
wait-for --host "https://localhost:443"
```

## Certificate Validation

The HTTPS probe (that is, where a target host is configured to use `https://` protocol) will attempt to validate the certificate chain and the hostname. If the certificate chain is invalid or the hostname doesn't match, the probe will exit with an error and the resource will be considered unavailable.

A valid HTTPS request on resources with custom certificates will require you to provide the CA certificate to the probe. By default, any certs stored in `/etc/ssl/certs/ca-certificates.crt` will be used to validate the connection.
