/*
 * Copyright © 2022 Durudex
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package grpc

import (
	"errors"

	"github.com/durudex/go-shared/status"

	"google.golang.org/grpc/codes"
	gs "google.golang.org/grpc/status"
)

// ServerErrorHandler handles status.Error to gRPC status.Error.
func ServerErrorHandler(err error) error {
	var e *status.Error

	// Checking as status.Error.
	if errors.As(err, &e) {
		switch e.Code {
		case status.CodeInternal:
			return gs.Error(codes.Internal, "Internal Server Error")
		case status.CodeNotFound:
			return gs.Error(codes.NotFound, e.Message)
		case status.CodeAlreadyExists:
			return gs.Error(codes.AlreadyExists, e.Message)
		case status.CodeInvalidArgument:
			return gs.Error(codes.InvalidArgument, e.Message)
		}
	}

	return gs.Error(codes.Internal, "Internal Server Error")
}
