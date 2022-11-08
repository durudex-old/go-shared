/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package grpc

import (
	"context"
	"net"

	"github.com/durudex/go-shared/crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Handler is used in gRPC server.
type Handler interface {
	// Register is intended for registration of gRPC server handlers.
	Register(srv *grpc.Server)
}

// ServerConfig stores the configurations needed to run the gRPC server.
type ServerConfig struct {
	// Host indicates on which host the server will be started.
	Host string

	// Port indicates on which port the server will be started.
	Port string

	// TLS setup configuration for creating TLS config.
	TLS tls.PathConfig
}

// Server that implements the run method.
type Server struct {
	// Server is an internal grpc.Server structure.
	Server *grpc.Server

	// Server handler.
	handler Handler

	// Server configurations.
	cfg ServerConfig
}

// NewServer returns the new gRPC server.
func NewServer(cfg ServerConfig, handler Handler) *Server {
	// Getting gRPC server options.
	options, err := GetServerOptions(cfg.TLS)
	if err != nil {
		panic("error getting gRPC server options: " + err.Error())
	}

	return &Server{
		Server:  grpc.NewServer(options...),
		handler: handler,
		cfg:     cfg,
	}
}

// Run starts the gRPC server.
func (s *Server) Run() {
	address := s.cfg.Host + ":" + s.cfg.Port

	// Create a new TCP listener.
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic("error creating a new listener: " + err.Error())
	}

	// Registration gRPC server handlers.
	s.handler.Register(s.Server)

	// Start gRPC server listening.
	if err := s.Server.Serve(lis); err != nil {
		panic("error starting gRPC server: " + err.Error())
	}
}

// GetServerOptions returns the options to be used by the gRPC server.
func GetServerOptions(cfg tls.PathConfig) ([]grpc.ServerOption, error) {
	var opts []grpc.ServerOption

	// Sets default options.
	opts = append(opts, grpc.UnaryInterceptor(UnaryInterceptor))

	// Check if TLS is enabled.
	if cfg.Enable {
		// Loading TLS configuration from certificate files.
		creds, err := tls.LoadConfig(cfg)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.Creds(credentials.NewTLS(creds)))
	}

	return opts, nil
}

// UnaryInterceptor is called as a wrapper between the request and the handler.
func UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Using a handler to complete a unary request.
	h, err := handler(ctx, req)
	if err != nil {
		return h, ServerErrorHandler(err)
	}

	return h, nil
}
