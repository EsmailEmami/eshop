//go:build !test && !coverage
// +build !test,!coverage

package db

import (
	"context"
	"fmt"
	"net/url"

	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type PgxIface interface {
	// Acquire returns a connection (*Conn) from the Pool
	Acquire(ctx context.Context) (*pgxpool.Conn, error)

	// AcquireAllIdle atomically acquires all currently idle connections. Its intended use is for health check and
	// keep-alive functionality. It does not update pool statistics.
	// AcquireAllIdle(ctx context.Context) []*pgxpool.Conn

	// Begin acquires a connection from the Pool and starts a transaction. Unlike database/sql, the context only affects the begin command. i.e. there is no
	// auto-rollback on context cancellation. Begin initiates a transaction block without explicitly setting a transaction mode for the block (see BeginTx with TxOptions if transaction mode is required).
	// *pgxpool.Tx is returned, which implements the pgx.Tx interface.
	// Commit or Rollback must be called on the returned transaction to finalize the transaction block.
	Begin(ctx context.Context) (pgx.Tx, error)

	// BeginTx acquires a connection from the Pool and starts a transaction with pgx.TxOptions determining the transaction mode.
	// Unlike database/sql, the context only affects the begin command. i.e. there is no auto-rollback on context cancellation.
	// *pgxpool.Tx is returned, which implements the pgx.Tx interface.
	// Commit or Rollback must be called on the returned transaction to finalize the transaction block.
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)

	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)

	// Exec acquires a connection from the Pool and executes the given SQL.
	// SQL can be either a prepared statement name or an SQL string.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	// The acquired connection is returned to the pool when the Exec function returns.
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)

	// QueryRow acquires a connection and executes a query that is expected
	// to return at most one row (pgx.Row). Errors are deferred until pgx.Row's
	// Scan method is called. If the query selects no rows, pgx.Row's Scan will
	// return ErrNoRows. Otherwise, pgx.Row's Scan scans the first selected row
	// and discards the rest. The acquired connection is returned to the Pool when
	// pgx.Row's Scan method is called.
	//
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	//
	// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
	// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
	// needed. See the documentation for those types for details.
	QueryRow(context.Context, string, ...interface{}) pgx.Row

	// Query acquires a connection and executes a query that returns pgx.Rows.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	// See pgx.Rows documentation to close the returned Rows and return the acquired connection to the Pool.
	//
	// If there is an error, the returned pgx.Rows will be returned in an error state.
	// If preferred, ignore the error returned from Query and handle errors using the returned pgx.Rows.
	//
	// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
	// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
	// needed. See the documentation for those types for details.
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)

	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	// Ping acquires a connection from the Pool and executes an empty sql statement against it.
	// If the sql returns without error, the database Ping is considered successful, otherwise, the error is returned.
	Ping(context.Context) error

	// Close closes all connections in the pool and rejects future Acquire calls. Blocks until all connections are returned
	// to pool and closed.
	Close()

	// Reset closes all connections, but leaves the pool open. It is intended for use when an error is detected that would
	// disrupt all connections (such as a network interruption or a server state change).
	//
	// It is safe to reset a pool while connections are checked out. Those connections will be closed when they are returned
	// to the pool.
	// Reset()

	// Stat returns a pgxpool.Stat struct with a snapshot of Pool statistics.
	// Stat() *pgxpool.Stat
}

type SqlDB = PgxIface
type SqlDBTx = pgx.Tx

var sqlDbConn SqlDB

// hold connection of type pgxpool.Pool to have move functionality
var pgxPoolConn *pgxpool.Pool

// Connection returns a database connection
func Connection() (SqlDB, error) {
	if sqlDbConn != nil {
		return sqlDbConn, nil
	}

	return PgxPoolConn()
}

// PgxPoolConn returns pgxpool.Pool to have the full functionality of the library
func PgxPoolConn() (*pgxpool.Pool, error) {
	if pgxPoolConn != nil {
		return pgxPoolConn, nil
	}

	username := viper.GetString("appdb.username")
	password := viper.GetString("appdb.password")
	userpass := url.UserPassword(username, password)
	host := viper.GetString("appdb.host")
	dbname := viper.GetString("appdb.db_name")
	port := viper.GetString("appdb.port")

	conURL := url.URL{
		Scheme: "postgres",
		User:   userpass,
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbname,
	}

	dbURL := conURL.String()

	var err error
	pgxPoolConn, err = pgxpool.New(context.Background(), dbURL)

	if err != nil {
		return nil, err
	}
	return pgxPoolConn, nil
}
