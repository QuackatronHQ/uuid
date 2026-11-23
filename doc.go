// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uuid generates and inspects UUIDs.
//
// UUIDs are based on RFC 9562(obsoletes RFC 4122) and DCE 1.1: Authentication and Security
// Services.
//
// A UUID is a 16 byte (128 bit) array.  UUIDs may be used as keys to
// maps or compared directly.
//
// The package generates the following versions:
//
//	Version 1  NewUUID           time and node based
//	Version 2  NewDCESecurity    DCE security, time and node based
//	Version 3  NewMD5            name based, MD5
//	Version 4  New, NewRandom    random
//	Version 5  NewSHA1           name based, SHA-1
//	Version 6  NewV6             time and node based, sortable
//	Version 7  NewV7             Unix time based, sortable
//	Version 8  NewV8             custom, application defined bits
//
// New code should use Version 4 for opaque identifiers and Version 7 when the
// identifiers should sort by creation time.
package uuid
