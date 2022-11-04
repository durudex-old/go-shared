/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

// Postgres driver interface.
type Driver interface {
	// Begin acquires a connection from the Pool and starts a transaction. Unlike database/sql, the
	// context only affects the begin command. i.e. there is no auto-rollback on context cancellation.
	// Begin initiates a transaction block without explicitly setting a transaction mode for the block
	// (see BeginTx with TxOptions if transaction mode is required). *pgxpool.Tx is returned, which
	// implements the pgx.Tx interface. Commit or Rollback must be called on the returned transaction
	// to finalize the transaction block.
	Begin(ctx context.Context) (pgx.Tx, error)

	// Query acquires a connection and executes a query that returns pgx.Rows. Arguments should be
	// referenced positionally from the SQL string as $1, $2, etc. See pgx.Rows documentation to
	// close the returned Rows and return the acquired connection to the Pool.
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	// QueryRow acquires a connection and executes a query that is expected to return at most one
	// row (pgx.Row). Errors are deferred until pgx.Row's Scan method is called. If the query selects
	//  no rows, pgx.Row's Scan will return ErrNoRows. Otherwise, pgx.Row's Scan scans the first
	// selected row and discards the rest. The acquired connection is returned to the Pool when
	// pgx.Row's Scan method is called.
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	// Exec acquires a connection from the Pool and executes the given SQL. SQL can be either a
	// prepared statement name or an SQL string. Arguments should be referenced positionally from
	// the SQL string as $1, $2, etc. The acquired connection is returned to the pool when the Exec
	// function returns.
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)

	// Close closes all connections in the pool and rejects future Acquire calls. Blocks until all
	// connections are returned to pool and closed.
	Close()
}

// Postgres pool driver config structure.
type PoolConfig struct {
	// The full URL to connect to the database. It may also contain some configurations to create
	// a connection.
	URL string

	// The maximum and minimum number of pool connections to the database. Without a specified
	// value, the pgx settings will be used.
	MaxConns, MinConns int32

	// Implementation of the logger that will be used by the driver.
	Logger tracelog.Logger
}

// Configuring the postgres driver to connect to the database.
func (c *PoolConfig) Configure(cfg *pgxpool.Config) {
	// Setting the logger for the driver.
	if c.Logger != nil {
		cfg.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger: c.Logger, LogLevel: tracelog.LogLevelInfo,
		}
	}

	// Setting the max and min number postgres driver connections.
	if cfg.MaxConns > 1 {
		cfg.MaxConns, cfg.MinConns = c.MaxConns, c.MinConns
	}
}

// Creating a new postgres pool connection.
func NewPool(cfg *PoolConfig) (*pgxpool.Pool, error) {
	// Parsing database url from configuration.
	config, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}

	// Configuring the postgres driver to connect to the database.
	cfg.Configure(config)

	// Create a new pool connections by config.
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Ping a database connection.
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
