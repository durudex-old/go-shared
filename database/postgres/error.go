/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package postgres

import (
	"errors"

	"github.com/durudex/go-shared/status"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorHandler handles of errors that may occur when working with the database.
func ErrorHandler(err error, object string) error {
	// Checking for pgx driver error.
	if errors.Is(err, pgx.ErrNoRows) {
		return &status.Error{Code: status.CodeNotFound, Message: object + " not found"}
	}

	var pgErr *pgconn.PgError

	// Checking for pgconn driver error.
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return &status.Error{
				Code: status.CodeAlreadyExists, Message: object + " already exists",
			}
		}
	}

	return &status.Error{Code: status.CodeInternal, Message: "Internal Server Error"}
}
