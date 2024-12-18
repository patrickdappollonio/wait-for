package probes

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLPinger is a pinger for MySQL connections.
type MySQLPinger struct {
	DSN string
}

// Bootstrap sets up the pinger with the URL.
func (m *MySQLPinger) Bootstrap(u *url.URL) error {
	// Construct DSN from URL. For example:
	// mysql://user:pass@tcp(host:port)/dbname
	// We'll assume the URL is provided in a usable format.

	host := u.Host
	user := u.User.Username()
	pass, _ := u.User.Password()

	if user == "" {
		user = "root"
	}

	// We use the "tcp(host:port)" format for MySQL driver.
	m.DSN = fmt.Sprintf("%s:%s@tcp(%s)/", user, pass, host)
	return nil
}

// Ping attempts to connect to the host and ping the database.
func (m *MySQLPinger) Ping(ctx context.Context) error {
	db, err := sql.Open("mysql", m.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
