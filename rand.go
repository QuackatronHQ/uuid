// Copyright 2024 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"crypto/rand"
	"io"
	"sync"
	"sync/atomic"
)

const randPoolSize = 16 * 16

var (
	rander      = rand.Reader // random function
	poolEnabled atomic.Bool
	poolMu      sync.Mutex
	poolPos     = randPoolSize     // protected with poolMu
	pool        [randPoolSize]byte // protected with poolMu
)

// SetRand sets the random number generator to r, which implements io.Reader.
// If r.Read returns an error when the package requests random data then
// a panic will be issued.
//
// Calling SetRand with nil sets the random number generator to the default
// generator.
func SetRand(r io.Reader) {
	if r == nil {
		rander = rand.Reader
		return
	}
	rander = r
}

// EnableRandPool enables internal randomness pool used for Random
// (Version 4) UUID generation. The pool contains random bytes read from
// the random number generator on demand in batches. Enabling the pool
// may improve the UUID generation throughput significantly.
//
// Since the pool is stored on the Go heap, this feature may be a bad fit
// for security sensitive applications.
//
// EnableRandPool and DisableRandPool are safe to call concurrently with
// UUID generation. UUIDs generated while the pool is being toggled may come
// from either source.
func EnableRandPool() {
	poolEnabled.Store(true)
}

// DisableRandPool disables the randomness pool if it was previously
// enabled with EnableRandPool. Any buffered random bytes are discarded.
func DisableRandPool() {
	poolEnabled.Store(false)
	poolMu.Lock()
	defer poolMu.Unlock()
	poolPos = randPoolSize
}

// readRandom fills b from the pool when it is enabled, or from rander
// otherwise. len(b) must not exceed randPoolSize.
func readRandom(b []byte) error {
	if !poolEnabled.Load() {
		_, err := io.ReadFull(rander, b)
		return err
	}

	poolMu.Lock()
	defer poolMu.Unlock()
	if poolPos+len(b) > randPoolSize {
		if _, err := io.ReadFull(rander, pool[:]); err != nil {
			return err
		}
		poolPos = 0
	}
	copy(b, pool[poolPos:poolPos+len(b)])
	poolPos += len(b)
	return nil
}

// randomBits completely fills slice b with random data.
func randomBits(b []byte) {
	if _, err := io.ReadFull(rander, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
}
