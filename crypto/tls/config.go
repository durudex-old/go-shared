/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// ErrAppendCerts is returned when certificates are incorrectly parsed and appended.
var ErrAppendCerts = errors.New("error append certificates from PEM files")

// PathConfig stores certificate paths used to create TLS configuration.
type PathConfig struct {
	// Enable indicates the usage status of this configuration.
	Enable bool

	// CA path to the PEM file.
	Ca string

	// Cert path to the PEM file.
	Cert string

	// Key path to the PEM file.
	Key string
}

// LoadConfig returns the TLS configuration that was extracted from the PEM certificate files.
func LoadConfig(cfg PathConfig) (*tls.Config, error) {
	// Reading CA certificate file.
	b, err := os.ReadFile(cfg.Ca)
	if err != nil {
		return nil, err
	}

	// Creating a new certification pool.
	pool := x509.NewCertPool()
	// Appends certificates from PEM files.
	if !pool.AppendCertsFromPEM(b) {
		return nil, ErrAppendCerts
	}

	// Analytic public/private key pair from a pair of files.
	certificate, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}, nil
}
