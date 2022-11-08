/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package status

// StatusCode is used as a key for error codes.
type StatusCode int

const (
	// CodeInternal indicates an internal server error.
	CodeInternal StatusCode = iota

	// CodeNotFound indicates a not found error.
	CodeNotFound

	// CodeAlreadyExists indicates errors that already exist.
	CodeAlreadyExists

	// CodeInvalidArgument indicates an invalid argument error.
	CodeInvalidArgument
)
