// Copyright 2024 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"errors"
	"testing"
)

func TestNewV8(t *testing.T) {
	var data [16]byte
	for i := range data {
		data[i] = 0xff
	}
	u := NewV8(data)
	if v := u.Version(); v != 8 {
		t.Errorf("version = %s, want 8", v)
	}
	if u.Variant() != RFC4122 {
		t.Errorf("variant = %s, want RFC4122", u.Variant())
	}
	if want := "ffffffff-ffff-8fff-bfff-ffffffffffff"; u.String() != want {
		t.Errorf("NewV8(all ones) = %s, want %s", u, want)
	}
	if _, ok := u.Timestamp(); ok {
		t.Error("Timestamp() returned ok = true for a v8 UUID")
	}
}

func TestNewV8KeepsCustomBits(t *testing.T) {
	data := [16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0x0d, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}
	u := NewV8(data)
	for i := range data {
		if i == 6 || i == 8 {
			continue
		}
		if u[i] != data[i] {
			t.Errorf("byte %d = %#x, want %#x", i, u[i], data[i])
		}
	}
	if u[6]&0x0f != data[6]&0x0f {
		t.Errorf("low nibble of byte 6 = %#x, want %#x", u[6]&0x0f, data[6]&0x0f)
	}
	if u[8]&0x3f != data[8]&0x3f {
		t.Errorf("low bits of byte 8 = %#x, want %#x", u[8]&0x3f, data[8]&0x3f)
	}
}

func TestNewV8FromBytes(t *testing.T) {
	if _, err := NewV8FromBytes(make([]byte, 15)); !errors.Is(err, ErrInvalidLength) {
		t.Errorf("NewV8FromBytes(15 bytes) error = %v, want ErrInvalidLength", err)
	}
	u, err := NewV8FromBytes(make([]byte, 16))
	if err != nil {
		t.Fatal(err)
	}
	if want := "00000000-0000-8000-8000-000000000000"; u.String() != want {
		t.Errorf("NewV8FromBytes(zeros) = %s, want %s", u, want)
	}
}
