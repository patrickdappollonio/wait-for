package probes

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// PostgresPinger is a pinger for PostgreSQL connections.
type PostgresPinger struct {
	DSN string
}

// Bootstrap sets up the pinger with the PostgreSQL URL.
// Expected URL format: postgres://user:password@host:port/dbname
func (p *PostgresPinger) Bootstrap(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("failed to parse host %q: %v", host, err)
	}

	// Extract user credentials
	user := u.User.Username()
	pass, _ := u.User.Password()

	// Extract host (includes port if specified)
	hostname := u.Host

	// Extract database name, trimming the leading '/'
	dbname := strings.TrimPrefix(u.Path, "/")
	if dbname == "" {
		return fmt.Errorf("no database name specified in the URL")
	}

	// Handle query parameters
	queryParams := u.Query()

	// Construct the DSN (Data Source Name)
	// Example: postgres://user:password@host:port/dbname
	p.DSN = fmt.Sprintf("postgres://%s:%s@%s/%s?%s",
		user, pass, hostname, dbname, queryParams.Encode())

	return nil
}

// Ping attempts to connect to the PostgreSQL database and ping it.
func (p *PostgresPinger) Ping(ctx context.Context) error {
	// Open a connection to the database
	db, err := pgx.Connect(ctx, p.DSN)
	if err != nil {
		return fmt.Errorf("error opening PostgreSQL connection: %w", err)
	}
	defer db.Close(ctx)

	// Set a short timeout for the ping
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Attempt to ping the database
	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("error pinging PostgreSQL database: %w", err)
	}

	return nil
}
