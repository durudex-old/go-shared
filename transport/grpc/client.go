/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package grpc

import (
	"github.com/durudex/go-shared/crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnectionConfig stores the configurations needed to create a new connection to the gRPC service.
type ConnectionConfig struct {
	// Addr specifies the URL to which the connection will be created.
	Addr string

	// TLS setup configuration for creating TLS config.
	TLS tls.PathConfig
}

// Connection implements the connection and closure methods and stores the service implementation.
type Connection[T any] struct {
	// Service client implementation.
	Service T

	// Connecting the client to the gRPC service.
	conn *grpc.ClientConn
}

// Connect creates a new connection to the gRPC service using the service constructor and configuration.
func Connect[T any](fc func(grpc.ClientConnInterface) T, cfg ConnectionConfig) *Connection[T] {
	var client Connection[T]

	// Creating a new connection to the gRPC service.
	conn, err := client.connect(cfg)
	if err != nil {
		panic("error creating a new connection to the service: " + err.Error())
	}

	// Using the constructor function to create a service.
	client.Service = fc(conn)

	return &client
}

// Connect creates a connection to the gRPC service using configuration.
func (c *Connection[T]) connect(cfg ConnectionConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	// Check if TLS is enabled.
	if cfg.TLS.Enable {
		// Loading TLS configuration from certificate files.
		creds, err := tls.LoadConfig(cfg.TLS)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(creds)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	var err error

	// Creating a new connection to the gRPC service.
	c.conn, err = grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		return nil, err
	}

	return c.conn, nil
}

// Close closes the gRPC connection to the service.
func (c *Connection[T]) Close() { c.conn.Close() }
