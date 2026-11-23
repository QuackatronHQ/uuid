// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

// MarshalText implements encoding.TextMarshaler.
func (uuid UUID) MarshalText() ([]byte, error) {
	return uuid.AppendText(make([]byte, 0, 36))
}

// AppendText appends the canonical string form of uuid to b. It implements
// encoding.TextAppender and avoids an allocation when b has spare capacity.
func (uuid UUID) AppendText(b []byte) ([]byte, error) {
	n := len(b)
	b = append(b, make([]byte, 36)...)
	encodeHex(b[n:], uuid)
	return b, nil
}

// AppendBinary appends the 16 raw bytes of uuid to b. It implements
// encoding.BinaryAppender.
func (uuid UUID) AppendBinary(b []byte) ([]byte, error) {
	return append(b, uuid[:]...), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (uuid *UUID) UnmarshalText(data []byte) error {
	id, err := ParseBytes(data)
	if err != nil {
		return err
	}
	*uuid = id
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (uuid UUID) MarshalBinary() ([]byte, error) {
	return uuid[:], nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (uuid *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return invalidBinaryLengthError{len(data)}
	}
	copy(uuid[:], data)
	return nil
}
