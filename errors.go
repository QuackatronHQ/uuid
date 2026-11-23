// Copyright 2024 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidUUIDFormat is returned when the input has the right length
	// but is not made of hex digits and hyphens in the expected places.
	ErrInvalidUUIDFormat = errors.New("invalid UUID format")

	// ErrInvalidBracketedFormat is returned by Validate when a 38 byte input
	// is not enclosed in curly braces.
	ErrInvalidBracketedFormat = errors.New("invalid bracketed UUID format")

	// ErrInvalidURNPrefix matches any URNPrefixError with errors.Is.
	ErrInvalidURNPrefix = URNPrefixError{}

	// ErrInvalidLength matches any length error returned by this package,
	// for both text and binary inputs, with errors.Is.
	ErrInvalidLength = invalidLengthError{}
)

// URNPrefixError is returned when a 45 byte input does not start with
// "urn:uuid:".
type URNPrefixError struct{ prefix string }

func (e URNPrefixError) Error() string {
	return fmt.Sprintf("invalid urn prefix: %q", e.prefix)
}

func (e URNPrefixError) Is(target error) bool {
	_, ok := target.(URNPrefixError)
	return ok
}

type invalidLengthError struct{ len int }

func (err invalidLengthError) Error() string {
	return fmt.Sprintf("invalid UUID length: %d", err.len)
}

func (e invalidLengthError) Is(target error) bool {
	_, ok := target.(invalidLengthError)
	return ok
}

// invalidBinaryLengthError is returned when decoding a binary UUID that is
// not exactly 16 bytes long. It keeps the historical error message so callers
// matching on the string keep working, and matches ErrInvalidLength.
type invalidBinaryLengthError struct{ len int }

func (err invalidBinaryLengthError) Error() string {
	return fmt.Sprintf("invalid UUID (got %d bytes)", err.len)
}

func (err invalidBinaryLengthError) Is(target error) bool {
	switch target.(type) {
	case invalidBinaryLengthError, invalidLengthError:
		return true
	}
	return false
}

// IsInvalidLengthError is matcher function for custom error invalidLengthError
func IsInvalidLengthError(err error) bool {
	return errors.Is(err, ErrInvalidLength)
}
