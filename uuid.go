// Copyright 2018 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
// 9562.
type UUID [16]byte

// A Version represents a UUID's version.
type Version byte

// A Variant represents a UUID's variant.
type Variant byte

// Constants returned by Variant.
const (
	Invalid   = Variant(iota) // Invalid UUID
	RFC4122                   // The variant specified in RFC9562(obsoletes RFC4122).
	Reserved                  // Reserved, NCS backward compatibility.
	Microsoft                 // Reserved, Microsoft Corporation backward compatibility.
	Future                    // Reserved for future definition.
)

// RFC9562 added V6 and V7 of UUID, but did not change specification of V1 and V4
// implemented in this module. To avoid creating new major module version,
// we still use RFC4122 for constant name.
const Standard = RFC4122

// Parse decodes s into a UUID or returns an error if it cannot be parsed.  Both
// the standard UUID forms defined in RFC 9562
// (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx and
// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx) are decoded.  In addition,
// Parse accepts non-standard strings such as the raw hex encoding
// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx and 38 byte "Microsoft style" encodings,
// e.g.  {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}.  Only the middle 36 bytes are
// examined in the latter case.  Parse should not be used to validate strings as
// it parses non-standard encodings as indicated above.
func Parse(s string) (UUID, error) {
	switch len(s) {
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36:

	// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + 9:
		if !strings.EqualFold(s[:9], "urn:uuid:") {
			return Nil, URNPrefixError{s[:9]}
		}
		s = s[9:]

	// {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
	case 36 + 2:
		s = s[1:]

	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	case 32:
		return decodeRawHex(s)

	default:
		return Nil, invalidLengthError{len(s)}
	}
	return decodeHyphenated(s)
}

// ParseBytes is like Parse, except it parses a byte slice instead of a string.
func ParseBytes(b []byte) (UUID, error) {
	switch len(b) {
	case 36: // xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + 9: // urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		if !bytes.EqualFold(b[:9], []byte("urn:uuid:")) {
			return Nil, URNPrefixError{string(b[:9])}
		}
		b = b[9:]
	case 36 + 2: // {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
		b = b[1:]
	case 32: // xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		return decodeRawHex(b)
	default:
		return Nil, invalidLengthError{len(b)}
	}
	return decodeHyphenated(b)
}

// MustParse is like Parse but panics if the string cannot be parsed.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustParse(s string) UUID {
	uuid, err := Parse(s)
	if err != nil {
		panic(`uuid: Parse(` + s + `): ` + err.Error())
	}
	return uuid
}

// FromBytes creates a new UUID from a byte slice. Returns an error if the slice
// does not have a length of 16. The bytes are copied from the slice.
func FromBytes(b []byte) (uuid UUID, err error) {
	err = uuid.UnmarshalBinary(b)
	return uuid, err
}

// Must returns uuid if err is nil and panics otherwise.
func Must(uuid UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return uuid
}

// Validate returns an error if s is not a properly formatted UUID in one of the following formats:
//
//	xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//	urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//	{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
//
// It returns an error if the format is invalid, otherwise nil.
func Validate(s string) error {
	switch len(s) {
	// Standard UUID format
	case 36:

	// UUID with "urn:uuid:" prefix
	case 36 + 9:
		if !strings.EqualFold(s[:9], "urn:uuid:") {
			return URNPrefixError{s[:9]}
		}
		s = s[9:]

	// UUID enclosed in braces
	case 36 + 2:
		if s[0] != '{' || s[len(s)-1] != '}' {
			return ErrInvalidBracketedFormat
		}
		s = s[1 : len(s)-1]

	// UUID without hyphens
	case 32:
		_, err := decodeRawHex(s)
		return err

	default:
		return invalidLengthError{len(s)}
	}

	_, err := decodeHyphenated(s)
	return err
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (uuid UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

// URN returns the RFC 2141 URN form of uuid,
// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx,  or "" if uuid is invalid.
func (uuid UUID) URN() string {
	var buf [36 + 9]byte
	copy(buf[:], "urn:uuid:")
	encodeHex(buf[9:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// Variant returns the variant encoded in uuid.
func (uuid UUID) Variant() Variant {
	switch {
	case (uuid[8] & 0xc0) == 0x80:
		return RFC4122
	case (uuid[8] & 0xe0) == 0xc0:
		return Microsoft
	case (uuid[8] & 0xe0) == 0xe0:
		return Future
	default:
		return Reserved
	}
}

// Version returns the version of uuid.
func (uuid UUID) Version() Version {
	return Version(uuid[6] >> 4)
}

func (v Version) String() string {
	if v > 15 {
		return fmt.Sprintf("BAD_VERSION_%d", v)
	}
	return fmt.Sprintf("VERSION_%d", v)
}

func (v Variant) String() string {
	switch v {
	case RFC4122:
		return "RFC4122"
	case Reserved:
		return "Reserved"
	case Microsoft:
		return "Microsoft"
	case Future:
		return "Future"
	case Invalid:
		return "Invalid"
	}
	return fmt.Sprintf("BadVariant%d", int(v))
}

// IsNil reports whether uuid is the Nil UUID (all zeros).
func (uuid UUID) IsNil() bool {
	return uuid == Nil
}

// IsMax reports whether uuid is the Max UUID (all bits set).
func (uuid UUID) IsMax() bool {
	return uuid == Max
}

// UUIDs is a slice of UUID types.
type UUIDs []UUID

// Strings returns a string slice containing the string form of each UUID in uuids.
func (uuids UUIDs) Strings() []string {
	var uuidStrs = make([]string, len(uuids))
	for i, uuid := range uuids {
		uuidStrs[i] = uuid.String()
	}
	return uuidStrs
}
