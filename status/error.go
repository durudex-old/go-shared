/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package status

import "fmt"

// Custom error structure.
type Error struct {
	// Error status code.
	Code StatusCode

	// Error message.
	Message string
}

// Getting error message in string type.
func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}
