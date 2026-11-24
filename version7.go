// Copyright 2023 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"errors"
	"io"
	"time"
)

// UUID version 7 features a time-ordered value field derived from the widely
// implemented and well known Unix Epoch timestamp source,
// the number of milliseconds since midnight 1 Jan 1970 UTC, leap seconds excluded.
// As well as improved entropy characteristics over versions 1 or 6.
//
// see https://datatracker.ietf.org/doc/html/rfc9562#name-uuid-version-7
//
// Implementations SHOULD utilize UUID version 7 over UUID version 1 and 6 if possible.
//
// NewV7 returns a Version 7 UUID based on the current time(Unix Epoch).
// Uses the randomness pool if it was enabled with EnableRandPool.
// On error, NewV7 returns Nil and an error
func NewV7() (UUID, error) {
	uuid, err := NewRandom()
	if err != nil {
		return uuid, err
	}
	makeV7(uuid[:], getV7Time())
	return uuid, nil
}

// NewV7FromReader returns a Version 7 UUID based on the current time(Unix Epoch).
// It uses NewRandomFromReader to fill the random bits.
// On error, NewV7FromReader returns Nil and an error.
func NewV7FromReader(r io.Reader) (UUID, error) {
	uuid, err := NewRandomFromReader(r)
	if err != nil {
		return uuid, err
	}

	makeV7(uuid[:], getV7Time())
	return uuid, nil
}

// NewV7WithTime returns a Version 7 UUID whose timestamp is t instead of the
// current time. It is useful for backfilling records that already have a
// creation time, so the UUID sorts with the rest of the data.
//
// UUIDs from NewV7WithTime do not take part in the monotonic sequence used by
// NewV7. Two calls with the same t differ only in their random bits. Times
// before the Unix epoch return ErrInvalidV7Time.
func NewV7WithTime(t time.Time) (UUID, error) {
	if t.Before(time.Unix(0, 0)) {
		return Nil, ErrInvalidV7Time
	}
	uuid, err := NewRandom()
	if err != nil {
		return uuid, err
	}
	makeV7(uuid[:], v7TimeOf(t.UnixNano()))
	return uuid, nil
}

// ErrInvalidV7Time is returned by NewV7WithTime for times before the Unix
// epoch, which cannot be encoded in the 48 bit unix_ts_ms field.
var ErrInvalidV7Time = errors.New("uuid: v7 time is before the Unix epoch")

// makeV7 fills 48 bits of time (uuid[0] - uuid[5]) and sets version b0111
// (uuid[6]) from the packed value ts returned by getV7Time or v7TimeOf.
// uuid[8] already has the right variant (Variant is 10).
func makeV7(uuid []byte, ts int64) {
	/*
		 0                   1                   2                   3
		 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                           unix_ts_ms                          |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|          unix_ts_ms           |  ver  |  rand_a (12 bit seq)  |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|var|                        rand_b                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                            rand_b                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/
	_ = uuid[15] // bounds check

	t, s := ts>>12, ts&0xfff

	uuid[0] = byte(t >> 40)
	uuid[1] = byte(t >> 32)
	uuid[2] = byte(t >> 24)
	uuid[3] = byte(t >> 16)
	uuid[4] = byte(t >> 8)
	uuid[5] = byte(t)

	uuid[6] = 0x70 | (0x0F & byte(s>>8))
	uuid[7] = byte(s)
}

// lastV7time is the last time we returned stored as:
//
//	52 bits of time in milliseconds since epoch
//	12 bits of (fractional nanoseconds) >> 8
var lastV7time int64

const nanoPerMilli = 1000000

// v7TimeOf packs nano as (milliseconds << 12 | sequence), where the 12 bit
// sequence is the sub-millisecond remainder divided by 256.
func v7TimeOf(nano int64) int64 {
	milli := nano / nanoPerMilli
	// Sequence number is between 0 and 3906 (nanoPerMilli>>8)
	seq := (nano - milli*nanoPerMilli) >> 8
	return milli<<12 + seq
}

// getV7Time returns the current time packed as (milli << 12 + seq).
// The returned value is guaranteed to be greater than the value returned by
// any previous call to getV7Time.
func getV7Time() int64 {
	timeMu.Lock()
	defer timeMu.Unlock()

	now := v7TimeOf(timeNow().UnixNano())
	if now <= lastV7time {
		now = lastV7time + 1
	}
	lastV7time = now
	return now
}
