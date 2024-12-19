# Configuring `wait-for` with a configuration file

`wait-for` supports reading configuration from a file from a file called `targets.yaml` (configurable with `--config`). This is useful when you have a lot of hosts to check and you don't want to pass them all as command-line arguments.

The following is an example YAML configuration file:

```yaml
# file: targets.yaml
hosts:
  - "tcp://localhost:8080"
  - "udp://localhost:53"
timeout: 30s
every: 2s
verbose: true
```

This is equal to calling the CLI with the following arguments:

```bash
wait-for \
  --host "tcp://localhost:8080" \
  --host "udp://localhost:53" \
  --timeout 30s \
  --every 2s \
  --verbose
```

You can mix-and-match hosts: any hosts provided via the configuration file will be merged with the hosts provided via the command line argument `--host` or `-s`, for example, the following config file and the following command will ping all endpoints (both from the config file and the command line):

```bash
$ cat targets.yaml
# file: targets.yaml
hosts:
  - "tcp://localhost:8080"
  - "udp://localhost:53"
timeout: 30s
every: 2s
verbose: true

$ wait-for \
  --host "localhost:80" \
  --host "localhost:81" \
  --timeout 10s
```

The above command will ping the following hosts by merging the two sources (configuration file and command-line flags):

```text
tcp://localhost:8080
udp://localhost:53
tcp://localhost:80
tcp://localhost:81
```

A host present both in the command-line arguments and in the configuration file will be pinged twice.
