# `wait-for`

A tiny Go application with zero dependencies. Given a number of TCP `host:port` pairs, the app will wait until either all are available or a timeout is reached.

Kudos to @vishnubob for the [original implementation in Bash](https://github.com/vishnubob/wait-for-it).

### Usage

The easiest way is to provide a list of `host:port` pairs and configure a timeout:

```bash
wait-for \
  --host "google.com:443" \
  --host "mysql.example.com:3306" \
  --timeout 10s
```

This will ping both `google.com` on port `443` and `mysql.example.com` on port `3306`. If they both start accepting connections within 10 seconds, the app will exit with a `0` exit code. If either one does not start accepting connections within 10 seconds, the app will exit with a `1` exit code, which will allow you to catch the error in CI/CD environments.

All the parameters accepted by the application are shown in the help section, as shown below.

```
wait-for allows you to wait for a TCP resource to respond to requests.

It does this by performing a TCP connection to the specified host and port. If there's
no resource behind it and the connection cannot be established, the request is retried
until either the timeout is reached or the resource becomes available.

By default, the standard timeout is 10 seconds.

For documentation, visit: https://github.com/patrickdappollonio/wait-for.

Usage:
  wait-for [flags]

Flags:
  -e, --every duration     time to wait between each request attempt against the host (default 1s)
  -h, --help               help for wait-for
  -s, --host strings       hosts to connect to in the format "host:port"
  -t, --timeout duration   maximum time to wait for the endpoints to respond before giving up (default 10s)
  -v, --verbose            enable verbose output -- will print every time a request is made
      --version            version for wait-for
```

### Usage with Kubernetes

Simply use this tool as an `initContainer` before your application runs, and validate whether your databases or any TCP-accessible resource (such as websites, too) are up and running, or fail early with proper knowledge of the situation.

In the example above, before we run the `nginx` container, we use `initContainers` to ensure that `google.com` responds on port `443`, `mysql.example.com` responds on port `3306` -- since it's the default port for MySQL -- and we also enable the verbose mode, which allows you to see the output of each probe against these hosts.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-container-demo
spec:
  initContainers:
  - name: wait-for
    image: ghcr.io/patrickdappollonio/wait-for:latest
    env:
    - name: POSTGRES_HOST
      value: "postgres.default.svc.cluster.local:5432"
    command:
      - /wait-for
    args:
      - --host="google.com:443"
      - --host="mysql.example.com:3306"
      - --host="$(POSTGRES_HOST)"
      - --verbose
  containers:
  - name: nginx-container
    image: nginx
```
