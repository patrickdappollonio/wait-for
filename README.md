# `wait-for`

[![Github Downloads](https://img.shields.io/github/downloads/patrickdappollonio/wait-for/total?color=orange&label=github%20downloads)](https://github.com/patrickdappollonio/wait-for/releases)

A Go application with zero dependencies. Given a number of hosts, the app will wait until either all are available or a timeout is reached. `wait-for` supports pinging several host types (see [supported probes](#supported-probes)), by prefixing the host with a specific protocol. If no prefix is provided, the app will default to TCP.

Kudos to @vishnubob for the [original implementation in Bash](https://github.com/vishnubob/wait-for-it).

### Installation

There are three ways to use `wait-for`: as a [Docker container](#container-image), as a [standalone binary](#standalone-binary), or [via Homebrew](#homebrew). The Docker container is the easiest way to use it, especially in Kubernetes environments.

#### Container image

The container has multiple versions: there's always a specific tag in the format `major.minor.patch` version, but it's recommended to use the `major` version tag, which will always point to the latest stable version with no risk of breaking changes:

```bash
ghcr.io/patrickdappollonio/wait-for:latest # maps to latest version
ghcr.io/patrickdappollonio/wait-for:v2     # maps to version 2 (recommended)
```

For all available images, see [Packages](https://github.com/users/patrickdappollonio/packages/container/package/wait-for).

#### Standalone binary

Download the latest release from our [releases page](https://github.com/patrickdappollonio/wait-for/releases/latest) and place it in your `$PATH`.

#### Homebrew

Install it using a Homebrew tap:

```bash
brew install patrickdappollonio/tap/wait-for
```

### Usage

The easiest way is to provide a list of `host:port` pairs and configure a timeout:

```bash
wait-for \
  --host "google.com:443" \
  --host "mysql.example.com:3306" \
  --timeout 10s
```

This will ping both `google.com` on port `443` and `mysql.example.com` on port `3306` via TCP. If they both start accepting connections within 10 seconds, the app will exit with a `0` exit code. If either one does not start accepting connections within 10 seconds, the app will exit with a `1` exit code, which will allow you to catch the error in CI/CD environments.

All the parameters accepted by the application are shown in the help section, as shown below.

### Command-line help

```text
wait-for allows you to wait for a resource to respond to requests.

It does this by performing a connection to the specified host and port. If
there's no resource behind it and the connection cannot be established, the
request is retried until either the timeout is reached or the resource becomes
available.

Each protocol defines its own way of checking for the resource. For example, a
TCP connection will attempt to connect to the host and port specified, while a
MySQL connection will attempt to connect to the host and port, and then ping the
database.

By default, the standard timeout is 10 seconds but it can be customized for all
requests. The time between each request is 1 second, but this can also be
customized.

For documentation, visit: https://github.com/patrickdappollonio/wait-for.

Usage:
  wait-for [flags]

Examples:
  wait-for -s localhost:80                             wait for a web server to accept connections
  wait-for -s mysql.example.local:3306                 wait for a MySQL database to accept connections
  wait-for -s udp://localhost:53                       wait for a DNS server to accept connections
  wait-for --host localhost:80 --host localhost:81     wait for multiple resources to accept connections
  wait-for --host mysql://localhost:3306               wait until a MySQL database is ready to accept connections and responds to pings

Flags:
  -e, --every duration     time to wait between each request attempt against the host (default 1s)
  -h, --help               help for wait-for
  -s, --host strings       hosts to connect to in the format "host:port" with optional protocol prefix (tcp:// or udp://)
  -t, --timeout duration   maximum time to wait for the endpoints to respond before giving up (default 10s)
  -v, --verbose            enable verbose output -- will print every time a request is made
      --version            version for wait-for
```

### Supported probes

The following probes are supported:

* [TCP probe](docs/tcp-probe.md)
* [UDP probe](docs/udp-probe.md)
* [HTTP & HTTPS probe](docs/http-https-probe.md)
* [MySQL probe](docs/mysql-probe.md) *(experimental)*
* [PostgreSQL probe](docs/postgres-probe.md) *(experimental)*

If you're interested in adding a new probe, please refer to the [Adding new probes documentation](docs/readme.md#adding-new-probes).

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

### Validating connectivity to a MySQL or Postgres database

If you want to validate that a MySQL database is up and running, you can use the `mysql://` or `postgres://` prefix. This will attempt to connect to the host and port specified, and then ping the database as well. This is different than the default TCP probe, which only checks if the server is accepting connections on the specified port.

For more details, check the [MySQL probe documentation](docs/mysql-probe.md) and the [PostgreSQL probe documentation](docs/postgres-probe.md).

### Validating connectivity to an HTTP or HTTPS endpoint

If you want to validate that an HTTP or HTTPS endpoint is up and running, you can use the `http://` or `https://` prefix. This will attempt to connect to the host and port specified, and then perform an HTTP GET request to the root path (`/`) of the server where the server must respond within 1 second. This is different than the default TCP probe, which only checks if the server is accepting connections on the specified port.

For HTTPS requests, the certificate is also validated. For more details, check the [HTTP & HTTPS probe documentation](docs/http-https-probe.md).
