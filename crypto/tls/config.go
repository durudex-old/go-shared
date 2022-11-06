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

// LoadConfig returns the TLS configuration that was extracted from the PEM certificate files.
func LoadConfig(ca, cert, key string) (*tls.Config, error) {
	// Reading CA certificate file.
	f, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	// Creating a new certification pool.
	pool := x509.NewCertPool()
	// Appends certificates from PEM files.
	if !pool.AppendCertsFromPEM(f) {
		return nil, ErrAppendCerts
	}

	// Analytic public/private key pair from a pair of files.
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}, nil
}
