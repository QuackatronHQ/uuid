// Copyright 2024 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

// NewV8 returns a Version 8 UUID built from the 16 bytes in data.
//
// Version 8 leaves 122 bits to the application. NewV8 copies data and then
// overwrites the version nibble (data[6] bits 4-7) and the variant bits
// (data[8] bits 6-7). Callers that need all of their bits preserved should
// keep them out of those positions.
//
// see https://datatracker.ietf.org/doc/html/rfc9562#name-uuid-version-8
func NewV8(data [16]byte) UUID {
	uuid := UUID(data)
	setVersion(&uuid, 8)
	return uuid
}

// NewV8FromBytes is like NewV8 but takes a slice. It returns an error that
// matches ErrInvalidLength if b is not exactly 16 bytes long.
func NewV8FromBytes(b []byte) (UUID, error) {
	if len(b) != 16 {
		return Nil, invalidBinaryLengthError{len(b)}
	}
	var data [16]byte
	copy(data[:], b)
	return NewV8(data), nil
}
