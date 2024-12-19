# MySQL

The MySQL probe will attempt to connect to the host and port specified. Once connected, it will attempt to perform a "ping" with a 1 second timeout from establishing the connection. If the connection can be established successfully and the database responds to the ping, the probe will exit successfully.

If the connection cannot be established or the ping fails, the probe will retry until either the timeout is reached or the resource becomes available.

The probe makes no guarantees about the existence of a table or the validity of the data in the database. It merely checks if the server is accepting connections on the specified port and if the database responds to the ping.

Internally, [the probe uses the `github.com/go-sql-driver/mysql` package](https://github.com/go-sql-driver/mysql) to connect to the database and perform the ping. This means that the connection string must be in the format `mysql://user:password@host:port/dbname`.

## Security

Since the credentials have to be provided plain text in the command line, it's recommended to use this probe in a secure environment or dynamically create the configuration file for such purpose (we recommend [using something like `tgen` to generate the configuration file](https://github.com/patrickdappollonio/tgen) on-the-fly).

It is not possible to provide *just the password* via environment variables, since the application supports specifying multiple hosts and ports, and each one can have a different password. However, nothing prevents you from using environment variables as part of the connection string when providing command-line flags, like:

```bash
wait-for --host "mysql://$MYSQL_USER:$MYSQL_PASSWORD@localhost:3306"
```

## TLS Support

To perform TLS connections, the container or host running the probe must have the necessary certificates to establish the connection. The probe will not attempt to validate the certificate chain or the hostname, so it's recommended to use this probe in a secure environment or to validate the certificates in another way.

By default, any certs stored in `/etc/ssl/certs/ca-certificates.crt` will be used to validate the connection.
