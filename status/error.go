/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package status

import "fmt"

// Error stores the error message and code.
type Error struct {
	// Code stores the error code.
	Code StatusCode

	// Message stores the error message.
	Message string
}

// Error returns error message in string type.
func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}
