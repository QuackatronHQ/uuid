// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import "bytes"

// xvalues returns the value of a byte as a hexadecimal digit or 255.
var xvalues = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}

// xtob converts hex characters x1 and x2 into a byte.
func xtob(x1, x2 byte) (byte, bool) {
	b1 := xvalues[x1]
	b2 := xvalues[x2]
	return (b1 << 4) | b2, b1 != 255 && b2 != 255
}

// hyphenatedOffsets holds the index of the first hex digit of each byte in
// the canonical xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx form.
var hyphenatedOffsets = [16]int{
	0, 2, 4, 6,
	9, 11,
	14, 16,
	19, 21,
	24, 26, 28, 30, 32, 34,
}

// decodeHyphenated decodes the canonical 36 byte form at the start of s.
// It returns ErrInvalidUUIDFormat if the hyphens are missing or any digit is
// not hex. The caller must ensure len(s) >= 36.
func decodeHyphenated[T string | []byte](s T) (UUID, error) {
	var uuid UUID
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return uuid, ErrInvalidUUIDFormat
	}
	for i, x := range hyphenatedOffsets {
		v, ok := xtob(s[x], s[x+1])
		if !ok {
			return uuid, ErrInvalidUUIDFormat
		}
		uuid[i] = v
	}
	return uuid, nil
}

// decodeRawHex decodes the 32 byte form without hyphens. The caller must
// ensure len(s) == 32.
func decodeRawHex[T string | []byte](s T) (UUID, error) {
	var uuid UUID
	for i := range uuid {
		v, ok := xtob(s[i*2], s[i*2+1])
		if !ok {
			return uuid, ErrInvalidUUIDFormat
		}
		uuid[i] = v
	}
	return uuid, nil
}

// setVersion stamps version v and the RFC 9562 variant onto uuid.
func setVersion(uuid *UUID, v Version) {
	uuid[6] = (uuid[6] & 0x0f) | byte(v)<<4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
}

// Compare returns an integer comparing two uuids lexicographically. The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func Compare(a, b UUID) int {
	return bytes.Compare(a[:], b[:])
}
