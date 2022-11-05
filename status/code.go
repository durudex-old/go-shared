/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package status

// Status code key used in error codes.
type StatusCode int

const (
	// Status code used for internal server errors.
	CodeInternal StatusCode = iota

	// Status code used for not found errors.
	CodeNotFound

	// Status code used for already exists errors.
	CodeAlreadyExists

	// Status code used for invalid argument errors.
	CodeInvalidArgument
)
